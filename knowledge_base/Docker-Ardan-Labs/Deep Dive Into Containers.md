#docker #docker-container

## Namespaces, no containers
- Inside a linux, we don't really have containers, but everything is namespaces and c-groups (control groups)
- there are a number of namespaces
- namespaces are equivalent to different rooms of the building/house where the house is the entire container
### Namespaces are available in modern kernels:
#### pid
- makes each container have it's own processes
- we can also share those processes among containers #shared-pid-namespace
```sh
$ docker run --pid host -ti alpine # this shares all the processes of the host with the container
$ docker run --pid container:<container-id> -ti alpine # this shares the processes of the other container with this container
```
- the original container runs the code, the other container with which the processes are shared, just sits and looks at the configuration file, the original container is going to send a signal to the other container to reload the configuration
#### net 
- a way for each container to have it's own network stack
#### mnt (mount)
- let's each container have it's own file system or files
- mounting a volume in a container is only visible within that container
#### uts
- responsible for each container have it's own hostname
#### ipc
- inter-process-communication
- shared memory, semaphores, message boxes
- each container has it's own ipc resources
#### time
- to deal with time within container
#### user
- let's to re-map users from host machine to container to use different users
- by default docker doesn't use `user` namespace
```sh
# to start docker engine to use `user` namespace:
$ sudo dockerd --userns-remap <user-host>:<user-container> # starting docker daemon (dockerd)
```
```sh
$ docker run -ti alpine
/ # whoami
root
/ # top
^C
$ ~ ps faux | grep top
# the host machine will now have non root user running top

$ docker run -it --user 1 alpine
/ $ whoami
bin
```
## cgroup - Control Groups
- another building blocks of containers
- used to measure resources and limit resources
- to check how much resources used by the container
	- how much bytes of RAM
	- how many seconds of CPU
	- how many reads and writes on disk, etc
### Crowd control groups
- cgroups allow to group processes for special operations
	- freezer ~ mass - SIGSTOP/SIGCONT
		- atomic way to stop or continue - either all or none
	- perf_event (gather performance events on multiple processes)
	- cpuset (limit or pin processes to specific CPUs)
- cgroups can allow to limit amount RAM used for certain containers