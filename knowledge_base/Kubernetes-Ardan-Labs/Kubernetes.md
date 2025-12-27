#k8s

Sample application
```sh
➜  dockercoins git:(main) tree
.
├── compose.yml
├── docker-compose.images.yml
├── docker-compose.logging.yml
├── hasher
│   ├── Dockerfile
│   └── hasher.rb
├── rng
│   ├── Dockerfile
│   └── rng.py
├── Tiltfile
├── webui
│   ├── Dockerfile
│   ├── files
│   │   ├── d3.min.js
│   │   ├── index.html
│   │   ├── jquery-1.11.3.min.js
│   │   ├── rickshaw.min.css
│   │   └── rickshaw.min.js
│   └── webui.js
└── worker
    ├── Dockerfile
    └── worker.py
```
- compose yaml
```yaml
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
      - "./webui/files/:/files/"

  redis:
    image: redis

  worker:
    build: worker
```
![[Pasted image 20251227182708.png]]
![[Pasted image 20251227183040.png]]
![[Pasted image 20251227183054.png]]
### Service Discovery
```python
# worker/worker.py
redis = Redis("redis")
def get_random_bytes():
	r = requests.get("http://rng/32")
	return r.content

def hash_bytes(data):
	r = requests.post(
		"http://hasher/", 
		data=data, 
		headers={"Content-Type":"application/octet-stream"}
	)
```
- in k8s, containers don't directly have network aliases
- instead we have something called **service**
- when we stop a container, we send a single to the first process of the container, on hitting ^C.
- on sending the signal, some containers are well behaved who stop right away on receiving the signal
- some containers keep running for sometime even on receiving stop signals
- after 10 seconds of hitting stop (^C) signal, we hit timeout
- after 10 seconds, docker will kill the container
- 10 seconds might not always seem a big deal
- say while in production, there are 100 containers running our application, and our containers are handling signals properly and take around 3 seconds top stop on receiving a container
- let's say at that time we are doing a rolling update, in which we replace each container one by one
	- stop a container, wait for it to stop and start replacement container
	- do this one container at a time
- thus, 100 containers x 3s = 300 s = 5 min for rolling update
- if we don't handle signals properly, we wait for 10 seconds in just docker container or 30 seconds in k8s
	- 100 containers x 33 s = 3300 s = 55 min ~ 1hr