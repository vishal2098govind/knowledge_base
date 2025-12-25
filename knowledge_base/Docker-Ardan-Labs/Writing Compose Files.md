#docker #docker-compose

### composing wordsmith app
```
$ pwd
/Users/govind/Dev/personel/docker-ardan-labs/wordsmith
$ tree
.
├── db
│   ├── Dockerfile
│   └── words.sql
├── LICENSE
├── README.en.md
├── README.fr.md
├── README.md
├── web
│   ├── dispatcher.go
│   ├── Dockerfile
│   └── static
│       ├── angular.min.js
│       ├── app.js
│       ├── favicon.ico
│       ├── fonts
│       │   ├── font1.woff2
│       │   └── font2.woff2
│       ├── images
│       │   ├── homes.png
│       │   ├── lego_blue.png
│       │   ├── lego_light_blue.png
│       │   ├── lego_yellow.png
│       │   └── logo.svg
│       ├── index.html
│       └── style.css
└── words
    ├── Dockerfile
    ├── pom.xml
    └── src
        └── main
            └── java
                └── Main.java

10 directories, 23 files
```
- compose file:
- version field in compose file:
	- in compose v1, we need not even had to put version field
	- after docker added networking and volumes etc, to make the difference between old format and new format, they added version field of `version: "2"`
	- then, Docker Swarm came, an orchestration system for docker containers
	- in swarm, we can have a cluster of docker machines, scaling, etc
	- thus, maintainers of compose added the swarm specific options and decided it to be `version: "3"`
	- so, unless we use swarm features, we can use either v2 or v3 
```yaml
version: "3"

services:
	web: 
	  build: web
	  ports:
	    - 8888:80
	words:
	  build: words

	db:
	  build: db
```

- force re-build of images using compose:
```sh
$ docker compose up --build
```
- scaling containers:
```sh
$ docker compose up --build --scale words=7
```
- here, if we try to run multiple `web` containers we will get error because all of the containers try to map to same host port of 8888
```sh
$ docker compose up --scale web=2 -d
[+] Running 3/4
 ✔ Container wordsmith-db-1     Running                                           0.0s 
 ✔ Container wordsmith-web-1    Running                                           0.0s 
 ✔ Container wordsmith-words-1  Running                                           0.0s 
 ⠹ Container wordsmith-web-2    Starting                                          0.2s 
Error response from daemon: failed to set up container networking: driver failed programming external connectivity on endpoint wordsmith-web-2 (09ccb293feed007410634cc88ce4cce08c08b941aabb6b64caef92c135761a48): Bind for 0.0.0.0:8888 failed: port is already allocated
```
- how to fix this:
- a possible solution could be using environment variable in compose file:
```sh
version: "3"

services:
  web:
    build: web
    ports:
      - $PORT:80

  db:
    build: db

  words:
    build: words
```
- can also use .env file
```.env
PORT=9001
```
- compose automatically reads .env file
- volumes in compose:
```
$ cd ./web/static
```
- to easily do `docker exec <container-id>` while using the compose file, can directly name the service instead of the container-id
```sh
$ docker compose exec web bash
```
- adding volume:
```yaml
version: "3"

services:
  web:
    build: web
    ports:
      - $PORT:80
    volumes:
      - ./web/static:/go/static

  db:
    build: db

  words:
    build: words
```
- the beauty of compose file:
- even if we run `docker compose up` within a nested folder, it uses the compose file in the root dir
- the special env variable for compose : `COMPOSE_PROJECT_NAME` this is useful to overwrite the default container name taken by the compose from the project directory where Dockerfile is present.
```sh
➜  wordsmith mkdir customers/client-1 customers/client-2 customers/client-3
➜  wordsmith cd customers/client-1
➜  client-1 echo "COMPOSE_PROJECT_NAME=client-1" > .env
➜  client-1 docker compose up -d
➜  client-1 docker compose ps
NAME               IMAGE            COMMAND                  SERVICE   CREATED          STATUS          PORTS
client-1-db-1      client-1-db      "docker-entrypoint.s…"   db        11 seconds ago   Up 10 seconds   5432/tcp
client-1-web-1     client-1-web     "./dispatcher"           web       11 seconds ago   Up 10 seconds   0.0.0.0:9001->80/tcp
client-1-words-1   client-1-words   "java -Xmx8m -Xms8m …"   words     11 seconds ago   Up 10 seconds
```
- to overwrite some environment variables
```sh
➜  customers cat /client-1/.env ./client-2/.env ./client-3/.env
PORT=10001
COMPOSE_PROJECT_NAME=client-1
PORT=10002
COMPOSE_PROJECT_NAME=client-2
PORT=10003
COMPOSE_PROJECT_NAME=client-3
%                          
```
```sh
➜  wordsmith tree -L 3 -a
.
├── .env
├── customers
│   ├── client-1
│   │   └── .env
│   ├── client-2
│   │   └── .env
│   └── client-3
│       └── .env
├── db
│   ├── Dockerfile
│   └── words.sql
├── docker-compose.yaml
├── LICENSE
├── README.en.md
├── README.fr.md
├── README.md
├── web
│   ├── dispatcher.go
│   ├── Dockerfile
│   └── static
│       ├── angular.min.js
│       ├── app.js
│       ├── favicon.ico
│       ├── fonts
│       ├── images
│       ├── index.html
│       └── style.css
└── words
    ├── Dockerfile
    ├── pom.xml
    └── src
        └── main

13 directories, 20 files
```

### Special handling of volumes:
- when an image gets updated, compose automatically creates a new container
- the data in the old container is lost ...
- ... Except if the container is using a **volume**
- Compose will then re-attach that **volume** to the new container
	- and data is then retained across database upgrade
- all good database images use volumes
	- e.g. all official images
- to wipe out existing data in the volume, use `docker compose down` and `docker compose up` as it would create completely fresh new volume

### Docker compose in production scenarios
- would we actually use docker and compose? - it depends
- we can use k8s or docker/compose
- if we use k8s, can we still use compose?
- **generally** it's seen that people use docker/compose when doing local development because it's relatively simple
	- git clone
	- docker compose up
	- wait for some time
	- app is up and running locally
- it's really hard to do something simpler than docker compose locally
- in k8s world, it's going to be way more complicated to get started
- so **generally**, k8s is used in production 
- if our entire application can fit in a single server or machine
	- if application small and simple, without a lot of users and load, use docker/compose
- if application needs multiple servers or machines to manage lot of load and lot of users
	- use k8s