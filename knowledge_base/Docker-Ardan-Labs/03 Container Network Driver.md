#docker

```sh
$ docker network ls
NETWORK ID     NAME                    DRIVER    SCOPE
a16e1f5c7656   bridge                  bridge    local
417b95ec21c1   client-1_default        bridge    local
5339818d18dd   dev                     bridge    local
ef89e4e3bcc9   dev_default             bridge    local
b5c0bd226597   host                    host      local
576331861b93   kind                    bridge    local
5fe14ea25a27   none                    null      local # this is a none network present by default in docker engine
5de06379d13f   prod                    bridge    local
5d7a13b72ff4   prod_default            bridge    local
8dcac28be493   wordsmith-new_default   bridge    local
1b54a9d3a18c   wordsmith_default       bridge    local
```
## `null` network driver
- a `null` network doesn't have internet access, and can only talk locally
```sh
$ docker run -it --net none alpine
/ # ifconfig
lo        Link encap:Local Loopback
          inet addr:127.0.0.1  Mask:255.0.0.0
          inet6 addr: ::1/128 Scope:Host
          UP LOOPBACK RUNNING  MTU:65536  Metric:1
          RX packets:0 errors:0 dropped:0 overruns:0 frame:0
          TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:1000
          RX bytes:0 (0.0 B)  TX bytes:0 (0.0 B)

/ # ping google.com
ping: bad address 'google.com'
```
- any other network has internet access
```sh
$ docker run -it --net dev alpine
/ # ifconfig
eth0      Link encap:Ethernet  HWaddr 3A:29:CE:EE:EE:2D
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

/ # ping google.com
PING google.com (142.250.182.110): 56 data bytes
64 bytes from 142.250.182.110: seq=0 ttl=63 time=55.444 ms
64 bytes from 142.250.182.110: seq=1 ttl=63 time=62.439 ms
^C
--- google.com ping statistics ---
3 packets transmitted, 2 packets received, 33% packet loss
round-trip min/avg/max = 55.444/58.941/62.439 ms
```

### Use cases of `none` network
- when we don't trust what the container does, we might run it from within a `none` network to see if it is trying to bring or send something from or to the internet
## `host` network driver
- when we run a container with `--net host`, then we run the container that has access to the network interfaces of the host machine.
- running `ifconfig` on the host machine and from within a container with `--net host` will return exact same set of network interfaces
- not just network interfaces are same, but anything related to network, like routes , IP table rules, etc would be visible in that case
### Use cases of `host` network
- performance reasons - for those containers who regularly push around 10-100 GiB per sec or 100s - 1000s of packets per sec
- to have direct connection to the outside network from within my program, and avoid receiving the network packets passing through a port mappings and reach the container or leave the container.
- using `--net host` means the container is sitting directly on the network interfaces of the host machine
- for containers that need more than one port and the ports the container uses are dynamic
	- like a port for http, another port for https
	- some protocols for video conferencing or VoIP might need a range of dynamically open ports
- for agent monitoring agents to collect some metrics on the network interfaces
	- store them in some time series database, Prometheus, Datadog, etc
- as an application developer, we might not worry much about using `--net host`, this is more for the Ops guys, especially building and managing k8s clusters

## container network driver

```sh
➜  db docker run -d nginx
482fc99bfa145f8e8ac082f9c66c6036be7cbbadc60cbff0ee8b1cbf7850c7cc
➜  db docker run --net container:482fc99bfa145f8e8ac082f9c66c6036be7cbbadc60cbff0ee8b1cbf7850c7cc -it alpine
/ # ifconfig
eth0      Link encap:Ethernet  HWaddr A2:0B:DF:99:2C:29
          inet addr:172.17.0.2  Bcast:172.17.255.255  Mask:255.255.0.0
          UP BROADCAST RUNNING MULTICAST  MTU:65535  Metric:1
          RX packets:12 errors:0 dropped:0 overruns:0 frame:0
          TX packets:3 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:0
          RX bytes:1172 (1.1 KiB)  TX bytes:126 (126.0 B)

lo        Link encap:Local Loopback
          inet addr:127.0.0.1  Mask:255.0.0.0
          inet6 addr: ::1/128 Scope:Host
          UP LOOPBACK RUNNING  MTU:65536  Metric:1
          RX packets:0 errors:0 dropped:0 overruns:0 frame:0
          TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:1000
          RX bytes:0 (0.0 B)  TX bytes:0 (0.0 B)

/ # ping google.com
PING google.com (142.250.182.14): 56 data bytes
64 bytes from 142.250.182.14: seq=0 ttl=63 time=38.445 ms
64 bytes from 142.250.182.14: seq=1 ttl=63 time=36.694 ms
^C
--- google.com ping statistics ---
2 packets transmitted, 2 packets received, 0% packet loss
round-trip min/avg/max = 36.694/37.569/38.445 ms
/ # netstat -ntlp
Active Internet connections (only servers)
Proto Recv-Q Send-Q Local Address           Foreign Address         State       PID/Program name
tcp        0      0 0.0.0.0:80              0.0.0.0:*               LISTEN      -
tcp        0      0 :::80                   :::*                    LISTEN      -
/ # curl
/bin/sh: curl: not found
/ # apk add curl
( 1/10) Installing brotli-libs (1.2.0-r0)
( 2/10) Installing c-ares (1.34.6-r0)
( 3/10) Installing libunistring (1.4.1-r0)
( 4/10) Installing libidn2 (2.3.8-r0)
( 5/10) Installing nghttp2-libs (1.68.0-r0)
( 6/10) Installing nghttp3 (1.13.1-r0)
( 7/10) Installing libpsl (0.21.5-r3)
( 8/10) Installing zstd-libs (1.5.7-r2)
( 9/10) Installing libcurl (8.17.0-r1)
(10/10) Installing curl (8.17.0-r1)
Executing busybox-1.37.0-r30.trigger
OK: 13.2 MiB in 26 packages
/ # curl localhost:80
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
html { color-scheme: light dark; }
body { width: 35em; margin: 0 auto;
font-family: Tahoma, Verdana, Arial, sans-serif; }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>
/ # ps faux
PID   USER     TIME  COMMAND
    1 root      0:00 /bin/sh
   14 root      0:00 ps faux
/ #
```
- `--net container:<container-id>` means asking to use the network stack of other container
- meaning, the containers share same every networking thing
	- same ip address, routes, ip table rules

### Use cases
- let's say we have a code that needs to connect to a banking API, and for security reasons, we need a VPN to connect to that banking API
	- classic approach would be have a vpn client in one container
	- another container would have the banking api client
- another e.g. when we have a web server and a cache like redis/memcache and want to communicate among them really very fast for super high performance
	- can get high throughput and use less CPU
- in k8s, we have pods, a pod is a bunch of containers sharing the same network stack

### using with docker compose
```yaml
services:
	pause:
		ports:
			- 8080:8080
		image: k8s.gcr.io/pause
	
	etcd:
		network_mode: "service:pause" # this is equivalent to running etcd container with the same network as pause container i.e. $ docker run --net container:pause etcd
```