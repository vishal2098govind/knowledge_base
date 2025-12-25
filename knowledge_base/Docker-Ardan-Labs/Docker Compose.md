#docker #docker-compose #yaml

## Yaml - Yet another markup language
- YAML is used to represent data
- List of things in yaml:
```yaml
- paris
- berlin
- nairobi
- tokyo
  
[ paris, berlin, seattle, naorobi, tokyo]

[ "paris", "berlin", "seatle", "nairobi", "tokyo" ] # it's JSON but also valid
```
- YAML is superset of JSON
- for strings, the quotes are optional in yaml so that it's easier to type
- for lists, we can use bullets (i.e. `-`) in each line one item
- the following is object in json as well as yaml:
```json
{"firstname": "vishal", "lastname": "govind", "job": "sde"}
```
```yaml
{"firstname": "vishal", "lastname": "govind", "job": "sde"}
```
- another possible way of writing object in yaml:
```yaml
firstname: vishal
lastname: govind
job: sde
```
- list of people:
```json
[
	{"firstname": "vishal", "lastname": "govind", "job": "sde"},
	{"firstname": "smruti", "lastname": "rath", "job": "doctor"}
]
```
- corresponding yaml. yaml also supports comments
```yaml
# they are a couple
- 
  firstname: vishal
  lastname: govind
  job: software engineer
- 
  firstname: smruti
  lastname: rath
  job: doctor
```
- in a single yaml file, we can have multiple nodes which we can separate them using `---`
```yaml
- paris
- berlin
- nairobi
- tokyo
  
---
  
[ paris, berlin, seattle, naorobi, tokyo ]

---
[ "paris", "berlin", "seatle", "nairobi", "tokyo" ] # it's JSON but also valid
---
- 
  firstname: vishal
  lastname: govind
  job: software engineer
- 
  firstname: smruti
  lastname: rath
  job: doctor
```


- `docker-compose.yaml`
```yaml
version: "2"

services:
	rng:
		build: rng
		ports:
		  - "8001:80"
	
	hasher:
		build: hasher
		ports:
		  - "8002:80"
	
	webui:
		build: webui
		ports:
		  - "8000:80"
		volumes:
		  - "./webui/files/:./files/"
		    
	redis:
		image: redis
		
	worker:
		build: worker
```
- we have a services object where for each service/container, if "image" field is specified, pull it from Docker Hub and if "build" field is specified, build the image from a Dockerfile
- basic docker compose file for trainingwheels system
```sh
➜  trainingwheels git:(master) tree
.
├── compose.yml
├── docker-compose.ecs.yml
├── docker-compose.simple.yml
└── www
    ├── assets
    │   ├── css
    │   │   ├── bootstrap-responsive.min.css
    │   │   └── bootstrap.min.css
    │   └── js
    │       └── bootstrap.min.js
    ├── counter.py
    ├── Dockerfile
    └── templates
        ├── error.html
        └── index.html
```

```yaml
version: "3"

services:
	www:
		# this build field indicates docker engine to build the image using the ./www/Dockerfile 
		build: www
		# EXPOSE port
		ports:
		  - 5000
	
	redis:
		image: redis
```
- a more advanced compose would look like:
```yaml
version: "3"

services:
	www:
		build: www
		ports:
		  - ${PORT-8000}:5000
		user: nobody
		environment:
		  DEBUG: 1
		volumes:
		  - ./www:/src
	
	redis:
		image: redis
```
- the options that we have in the compose file are a way to represent the options for `docker run` command
- the above compose for www container is equivalent to:
```sh
$ ~ docker run -v $PWD/www:/src -e DEBUG=1 -u nobody -p 8000:5000 www python counter.py
```
- Docker compose file for wordsmith system
```yaml
version: "3"

services:
	web:
		build: web
		or
		image: the-name-of-an-image-on-some-registry
	
	words:
	
	db:
	  build: db
	  ports:
	    - port-on-machine:port-in-container
```

### docker compose commands:
```sh
# to start all containers in compose:
$ docker compose up
# start all containers in background
$ docker compose up -d
# see all containers in stack of the compose
➜  trainingwheels git:(master) docker compose ps
NAME      IMAGE     COMMAND   SERVICE   CREATED   STATUS    PORTS
➜  trainingwheels git:(master) docker compose ps -a
NAME                     IMAGE                COMMAND                  SERVICE   CREATED              STATUS                          PORTS
trainingwheels-redis-1   redis                "docker-entrypoint.s…"   redis     About a minute ago   Exited (0) About a minute ago   
trainingwheels-www-1     trainingwheels-www   "python counter.py"      www       About a minute ago   Exited (0) About a minute ago   
# to stop all containers:
$ docker compose stop
# to see logs of all containers:
$ docker compose logs
# to kill all containers:
$ docker compose kill
# to redo or remove everything created by docker compose up
$ docker compose down
➜  trainingwheels git:(master) docker compose down
[+] Running 3/3
 ✔ Container trainingwheels-www-1    Removed 0.1s 
 ✔ Container trainingwheels-redis-1  Removed 0.1s 
 ✔ Network trainingwheels_default    Removed 0.2s 
```
- docker compose creates a new network for each compose file. 
- thus we can spin up a container (say of a simple alpine image) in that network to debug these containers
```sh
$ docker run --net trainingwheels_default -it alpine 
/ # ping www
PING www (172.21.0.3): 56 data bytes
64 bytes from 172.21.0.3: seq=0 ttl=64 time=0.080 ms
64 bytes from 172.21.0.3: seq=1 ttl=64 time=0.088 ms
64 bytes from 172.21.0.3: seq=2 ttl=64 time=0.086 ms
^C
--- www ping statistics ---
3 packets transmitted, 3 packets received, 0% packet loss
round-trip min/avg/max = 0.080/0.084/0.088 ms
/ # ping db
ping: bad address 'db'
/ # ping redis
PING redis (172.21.0.2): 56 data bytes
64 bytes from 172.21.0.2: seq=0 ttl=64 time=0.081 ms
64 bytes from 172.21.0.2: seq=1 ttl=64 time=0.090 ms
64 bytes from 172.21.0.2: seq=2 ttl=64 time=0.094 ms
^C
--- redis ping statistics ---
3 packets transmitted, 3 packets received, 0% packet loss
round-trip min/avg/max = 0.081/0.088/0.094 ms
```
- docker compose has the ability to run same stack multiple times (each in difference **namespace** in a way i.e. "dev" namespace or "prod" namespace)
```sh
$ docker compose --project-name dev up -d
$ docker compose --project-name dev ps
NAME          IMAGE     COMMAND                  SERVICE   CREATED          STATUS          PORTS
dev-redis-1   redis     "docker-entrypoint.s…"   redis     35 seconds ago   Up 34 seconds   6379/tcp
dev-www-1     dev-www   "python counter.py"      www       35 seconds ago   Up 34 seconds   0.0.0.0:8000->5000/tcp
$ docker compose --project-name prod up -d
[+] Running 3/4
 ✔ www                     Built                                                                                                                         0.0s 
 ✔ Network prod_default    Created                                                                                                                       0.1s 
 ✔ Container prod-redis-1  Started                                                                                                                       0.3s 
 ⠸ Container prod-www-1    Starting                                                                                                                      0.4s 
Error response from daemon: failed to set up container networking: driver failed programming external connectivity on endpoint prod-www-1 (3121f5ea3bcb742274c686d4e9101383a74bd9219623a23d24f2dbb49d892502): Bind for 0.0.0.0:8000 failed: port is already allocated
$ docker compose --project-name prod ps
NAME           IMAGE     COMMAND                  SERVICE   CREATED          STATUS          PORTS
prod-redis-1   redis     "docker-ent
```
- here, the port 8000 is trying to be exposed for both prod and dev containers. which is not possible. How to address this?
	- can use ENV variables
	-   