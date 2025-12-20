#docker #docker-container 

## Hello World
```sh
$ docker run busybox echo hello world
```
- `busybox` is like a Swiss Army knife of unix commands
- `busybox` is a program with many classic unix commands

### running container in interactive mode (`-it`)
```zsh
$ docker run -it ubuntu #<----- -i => Input/Output -t => Terminal
```
- each time we run `docker run -it ubuntu`, we create a brand new container, why?
	- we can reuse the container but it's not the default workflow with Docker
	- What's the default workflow with Docker?
		- always start with a fresh container
		- if we need something **pre-installed** in our container, build a **custom image**
	- Why?
		- This puts a strong emphasis on automation and repeatability
		- Pets vs Cattle metaphors for two kinds of servers
			- **Pet Servers**: 
				- have distinct names and unique configurations
				- when they have an outage, we do everything we can to fix them
			- **Cattle Servers**:
				- have generic names and generic configurations
				- **configuration** is **enforced** by configuration management, golden images ...
				- when they have an outage, we can replace them immediately with a new server

## Background Containers
- Difference between foreground and background containers:
	- to run a container in background: use `-d` option (d stands detach or daemonise)
		- in UNIX systems, stuffs running in the background are called DAEMON 
			- **D**isk **A**synchronous **E**xecution **Mon**itor
```zsh
$ docker run -d jpetazzo/clock
aefaobgebg # <------ just returns the container id
```
- Docker took inspiration of background containers from the `&` command in UNIX shell
```zsh
$ sleep 10 &
[1] 12596 # <-------- just returns process id
$ ps # <----- lists processes currently running
```
in docker:
```zsh
$ docker run -d jpetazzo/clock
aefaobgebg # <------ just returns the container id
$ docker ps # <------- lists containers currently running
```
- To see what the containers are running within, use `docker logs <cid>`
- `docker logs <cid> -f` for **following** the log stream (~ to `tail -f` on UNIX systems to see files)
- docker by default uses some randomly generated name for docker containers. 
	- the default random names are of some well known scientists, mathematicians, inventors, mixed with some randomly generated mood
```zsh
$ docker logs <cname> -f #< ---- also works
```

- To **stop containers**, use `docker stop <cid/cname>`
	- might take a while to stop container
	- there are usually two ways to shut down a program
		- **grace full** shutdown
			- allowing to finish whatever it was doing at the moment
			- after that's done, shut down
			- `docker stop <cid>`  is a grace full shutdown
			- `docker stop` allows the container 10 or some seconds seconds after which, it forcefully shuts down the container if not stopped yet
			- `docker stop` sends signal term (`SIGTERM`) to the first process in that container
		- **force** shutdown
			- immediate, kill
			- `SIGKILL`
- the processes in the container are actual processes in the underlying machine
- if a program can run on a machine, it should be able to run in a container

### docker ps commands
`docker ps -l`: just shows the details of last container started
`docker ps -q`: "quick" i.e. just shows container ids, not whole table
`docker ps -ql`: just show container id of last container started

### stopping/kill container(s)
`docker kill <cid>`
`docker kill $(docker ps -q)`: kills all containers with ids returned by `docker ps -q`

### remove container completely
`docker rm <cid>`