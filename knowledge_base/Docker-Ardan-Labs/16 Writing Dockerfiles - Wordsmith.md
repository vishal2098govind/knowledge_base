#docker

```Dockerfile
FROM golang

COPY . .

RUN go build dispatcher.go

CMD ./dispatcher
```
- running a container with the above `CMD ./dispatcher` would start the go web server but we won't be able to exit out of it by entering anything like ^C or something. Why?
	- when we don't use Json String Arguments to CMD, docker uses `/bin/sh` to parse and further run the commands.
	- due to which the top most process (PID = 1) running in the container would be `/bin/sh` and not the go program
	- so, hitting ^C, docker sends the signal to the /bin/sh process since it has PID = 1
	- can find that by doing `docker exec 694e ps faux` 
```sh
➜  web git:(main) ✗ docker ps
CONTAINER ID   IMAGE     COMMAND                  CREATED          STATUS          PORTS     NAMES
694ed0bca97f   web       "/bin/sh -c ./dispat…"   56 seconds ago   Up 55 seconds             crazy_roentgen
➜  web git:(main) ✗ docker exec 694e ps faux
USER       PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
root        13  0.0  0.0   6396  3540 ?        Rs   11:53   0:00 ps faux
root         1  0.0  0.0   2680  1216 ?        Ss   11:52   0:00 /bin/sh -c ./dispatcher # <----(top most process i.e. PID = 1)
root         7  0.0  0.1 1602696 8084 ?        Sl   11:52   0:00 ./dispatcher
```
- `docker exec` is a way to run commands inside a running container 
- on hitting ^C, docker sends the OS signal to the first process in the container
- the default behavior of the shell (/bin/sh) on receiving ^C is to ignore it
```Dockerfile
FROM golang

COPY . .

RUN go build dispatcher.go

CMD ["./dispatcher"]
```
- running a container with the above `CMD ["./dispatcher"]` would start the go web server and also we can stop the server on hitting ^C.
	- this is because, the top most process running the container this time is the go program, and not the shell
	- so, hitting ^C, docker sends the signal to the ./dispatcher process since it has PID = 1
	- default behavior of go program on receiving ^C is to stop
```sh
➜  web git:(main) ✗ docker ps
CONTAINER ID   IMAGE     COMMAND          CREATED         STATUS         PORTS     NAMES
992538a3049a   web       "./dispatcher"   9 seconds ago   Up 9 seconds             pedantic_napier
➜  web git:(main) ✗ docker exec 992538a3049a ps faux
USER       PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
root        13  0.0  0.0   6396  3576 ?        Rs   12:01   0:00 ps faux
root         1  0.0  0.0 1602952 7808 ?        Ssl  11:59   0:00 ./dispatcher # <---- (top most process i.e. PID = 1)
```

- Generally, it's a good idea to have our program to get **(PID = 1)**
-  we can also use `CMD exec` to get PID = 1
```sh
CMD exec java -Xmx8m -Xms8m -jar words.jar
```

### Optimizing Dockerfile for faster build time
```Dockerfile
FROM ubuntu

RUN apt-get update
RUN apt-get install maven -y

WORKDIR /app

COPY . .

RUN mvn verify

WORKDIR /app/target

CMD ["java", "-Xmx8m", "-Xms8m", "-jar", "words.jar"]
```
Optimised:
```Dockerfile
FROM ubuntu

RUN apt-get update
RUN apt-get install maven -y

WORKDIR /app

COPY pom.xml . # <----- this allows to avoid re-runing `mvn verify` on every build unless pom.xml changes
RUN mvn verify

COPY . .

RUN mvn verify

WORKDIR /app/target

CMD ["java", "-Xmx8m", "-Xms8m", "-jar", "words.jar"]
```