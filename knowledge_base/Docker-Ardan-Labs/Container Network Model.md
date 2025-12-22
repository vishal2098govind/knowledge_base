#docker #container-network-model #service-discovery
-  **CNM** (Container Network Model) is a set of concepts introduced by Docker, when manipulating groups of containers 
- list available networks within docker: `docker network ls`
```
NETWORK ID     NAME      DRIVER    SCOPE
dd41f0e4f02d   bridge    bridge    local
b5c0bd226597   host      host      local
5fe14ea25a27   none      null      local
```
- can create networks: `docker network create dev`
```
$ docker network create dev
5339818d18dd196e7cb7ee9072d833681a8e79edde042783876854107104855d
$ docker network create prod
5de06379d13ff62623044c75f4c8415d5541ec890ad6770c638a5cfc11782fdc
$ docker network ls
NETWORK ID     NAME      DRIVER    SCOPE
dd41f0e4f02d   bridge    bridge    local
5339818d18dd   dev       bridge    local
b5c0bd226597   host      host      local
5fe14ea25a27   none      null      local
5de06379d13f   prod      bridge    local
```
- The scope attribute if "local", it means that it is a local network that only exists on that machine
- we can also have "global" scoped networks, to connect multiple machines with docker
- in practice, "global" is not used, other than docker swarm
- run a container in a specific network:
	- `docker run --net dev -it web`

```
➜  db git:(main) ✗ docker run --net dev -it alpine   
Unable to find image 'alpine:latest' locally
latest: Pulling from library/alpine
1074353eec0d: Pull complete 
Digest: sha256:865b95f46d98cf867a156fe4a135ad3fe50d2056aa3f25ed31662dff6da4eb62
Status: Downloaded newer image for alpine:latest
/ # ifconfig
eth0      Link encap:Ethernet  HWaddr 36:01:99:A9:BC:EB  
          inet addr:172.19.0.2  Bcast:172.19.255.255  Mask:255.255.0.0
          UP BROADCAST RUNNING MULTICAST  MTU:1500  Metric:1
          RX packets:17 errors:0 dropped:0 overruns:0 frame:0
          TX packets:3 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:0 
          RX bytes:1698 (1.6 KiB)  TX bytes:126 (126.0 B)

lo        Link encap:Local Loopback  
          inet addr:127.0.0.1  Mask:255.0.0.0
          inet6 addr: ::1/128 Scope:Host
          UP LOOPBACK RUNNING  MTU:65536  Metric:1
          RX packets:0 errors:0 dropped:0 overruns:0 frame:0
          TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:1000 
          RX bytes:0 (0.0 B)  TX bytes:0 (0.0 B)

```
```
✗ docker run --net prod -it alpine
/ # ifconfig
eth0      Link encap:Ethernet  HWaddr 1A:68:68:42:83:27  
          inet addr:172.20.0.2  Bcast:172.20.255.255  Mask:255.255.0.0
          UP BROADCAST RUNNING MULTICAST  MTU:1500  Metric:1
          RX packets:17 errors:0 dropped:0 overruns:0 frame:0
          TX packets:3 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:0 
          RX bytes:1698 (1.6 KiB)  TX bytes:126 (126.0 B)

lo        Link encap:Local Loopback  
          inet addr:127.0.0.1  Mask:255.0.0.0
          inet6 addr: ::1/128 Scope:Host
          UP LOOPBACK RUNNING  MTU:65536  Metric:1
          RX packets:0 errors:0 dropped:0 overruns:0 frame:0
          TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:1000 
          RX bytes:0 (0.0 B)  TX bytes:0 (0.0 B)

/ # 
```

```sh
$ docker run -d --net dev nginx
45b65a50342c
$ docker run -d --net prod nginx
91b258c64660
$  docker inspect 45b65a50342c 91b258c64660 | grep IPA
            "SecondaryIPAddresses": null,
            "IPAddress": "",
                    "IPAMConfig": null,
                    "IPAddress": "172.19.0.2",
            "SecondaryIPAddresses": null,
            "IPAddress": "",
                    "IPAMConfig": null,
                    "IPAddress": "172.20.0.2",
$ docker run -it --net dev alpine                       | $ docker run -it --net prod alpine
$ ping 172.19.0.2                                       | $ ping 172.20.0.2
$ PING 172.19.0.2 (172.19.0.2): 56 data bytes           | # works
64 bytes from 172.19.0.2: seq=0 ttl=64 time=0.171 ms    |
...
$ ping 172.20.0.2
# does not work
```
- By default, while within a network, we can only access a containers that are in the same network

## Service discovery
- Imagine in our code, at some time, we need to connect to our database
```python
DB_ADDRESS = env["DB_ADDRESS"] # how to find the address of the db is called service discovery. here we chose env
db_connection = pg_connect(DB_ADDRESS, username, pass) 
```
- what do containers provide here with **network-aliases**
- the idea is, we can just mention the name of the image of the container running db
```python
db_connection = pg_connect("db") # "db" is the name of our db image
```
- docker will transform the name to the address of the container running the Postgres database
```sh
$ docker run --net my_little_network --net-alias db my-db-image
```
- there is a dynamic DNS server embedded in the docker engine to manage domain name to address translation

```
$  docker run --net dev --net-alias dev-api -d nginx
15132d3ece640896b38e314420141d1563c887d812ee69ddecec45003daf41e4
$ docker run -it --net dev alpine
/ # ping dev-api
PING dev-api (172.20.0.3): 56 data bytes
64 bytes from 172.20.0.3: seq=0 ttl=64 time=0.146 ms
64 bytes from 172.20.0.3: seq=1 ttl=64 time=0.073 ms
64 bytes from 172.20.0.3: seq=2 ttl=64 time=0.123 ms
^C
--- dev-api ping statistics ---
3 packets transmitted, 3 packets received, 0% packet loss
round-trip min/avg/max = 0.073/0.114/0.146 ms
/ # ping prod-api
ping: bad address 'prod-api'
/ # exit
$ docker run --net prod --net-alias 
prod-api -d nginx
1b03e4fc728ec15ebf2a42a39a2ad32495f101eb7c6c157b6d89981c3320fed1
$ docker run -it --net prod alpine
/ # ping prod-api
PING prod-api (172.20.0.3): 56 data bytes
64 bytes from 172.20.0.3: seq=0 ttl=64 time=0.146 ms
64 bytes from 172.20.0.3: seq=1 ttl=64 time=0.073 ms
64 bytes from 172.20.0.3: seq=2 ttl=64 time=0.123 ms
^C
--- prod-api ping statistics ---
3 packets transmitted, 3 packets received, 0% packet loss
round-trip min/avg/max = 0.073/0.114/0.146 ms
/ # 
```

- This service discovery concept also maps to Kubernetes, where instead of a name corresponding to a container, we have a name corresponding to a cluster-ip or load balancer
- CNM (container network model) is the model used by docker
- Kubernetes uses a different model, architectured around CNI
- the two models are very much different
- concept of CNM/docker:
	- have different networks independent of each other
	- each network is isolated
	- here, we need to do extra steps to communicate among different containers
	- per-network service discovery
- concept of CNI/k8s:
	- we have one big network
	- all the containers (called pods k8s) are present in the same network
	- so by default, on k8s, everything can communicate with everything else
	- so, arguably, it's less secure than docker
	- but we add isolations with mechanism called Network Policies, which requires extra steps
	- in docker, we have isolation by default and need extra steps to communicate, 
	- in k8s, it's the other way around, we have communication by default and need extra steps to isolate
	- in k8s, it's per-namespace service discovery

### Service discovery in practice