#docker #docker-images

## What is an image?
- Image = files + metadata
- These files form the root filesystem of our container
- The metadata can indicate a number of things like:
	- author of the image - a **decorative** metadata
	- command to execute in container when starting - **active** metadata
	- environment variables to be set - an **active** metadata
	- etc
### Images have layers
- Images are made of **layers**, conceptually stacked on top of each other.
	- Each layer can add, change, and remove files and metadata
	- Example: a java webapp
		- CentOS base layer
		- Packages ad configuration files
		- JRE layer
		- Tomcat layer
		- App's dependencies
		- Application Code and assets
		- Application configurations 
- Why layers?
	- when we make changes, these layers can be updated independent of each other
	- on release a new version of the app, only the application-code-assets layer and application-configuration layers need to be updated, and other layers can remain the same
	- saves bytes to be transferred over the network while releasing the new version
	- saves disk space storing images
	- saves memory
- In the stack of layers, we can only write on the top layer, and not on layers below
- The **image layers are read-only layers** and cannot be updated
- the **container layers are read-write** layers and thus can be updated or removed

### Difference between containers and images
- An image is a read-only filesystem
- A container is an encapsulated set of processes,
	- running in  a read-write copy of that filesystem
### Docker uses `copy-on-write`
- `copy-on-write` is a mechanism can be found also in many other systems, like UNIX and others 
- metaphor for `copy-on-write`
	- consider going to a library where we're told we could take notes if we want in the book itself, it's totally fine, and those notes are going to be just for us
	- so, at the moment when we are about to write notes on the book, at the moment when our pen is going to hit the page of the book, imagine that the time stops, 
		- and the librarian shows up 
		- and gets the book, 
		- makes a copy of the page that was going to be written by us
		- thus, basically when we would write, we instead of writing on the original page, we would be writing on the copy of the page
- This is exactly what happens when we try to write a file, in the image, from within a container, instead of changing the original file of the image, it's going to make a copy of the file and writes on the copy
- if there was no `copy-on-write`, then we would have to make an entire copy of the image every time before starting the container.
- `docker run` starts a container from a given image

### If image layer is read-only, how to change an image?
- We don't change image
- We create a new container from that image
- Then we make changes to that container
- When we are satisfied with those changes, we transform those changes into a new layer
- A new image is created by stacking the new layer on top of the old image

### Creating images
- `docker commit`
- `docker build` i.e. `Dockerfile` (**used 99% of the time**)

#### Images namespace and tags
- There are 3 namespaces:
	- Official images **(Root namespace)**: 
		- one word names
		- e.g. `ubuntu`, `busybox`
		- usually very well designed with optimizations
	- User (and organizations) images: **(User namespace)**
		- `<username>/<name>` 
		- hosting images on Docker Hub
		- e.g. `vishal2098govind/clock`
	- Self-hosted images **(Registry namespace)**
		- If we don't want to put our private images on the Docker Hub
			- we can use our own repository of images, called **`registries`**
		- e.g. 
			- `registry.example.com:5000/my-private/image`
			- `quay.io/cores/etcd`
				- `quay.io` is a popular registry
			- `gcr.io/google-containers-hugo`
- docker images are versioned using tags
- default is `latest` tag

### Building images interactively
#### Docker Commit
- using `docker commit`
```
$ docker run -it ubuntu
$ root@c734a94943b5:/# apt-get install figlet
$ root@c734a94943b5:/# figlet yay it works
$ root@c734a94943b5:/# exit
$ docker ps -a
CONTAINER ID   IMAGE     COMMAND       CREATED         STATUS                      PORTS     NAMES
c734a94943b5   ubuntu    "/bin/bash"   2 minutes ago   Exited (0) 10 seconds ago             nostalgic_poincare

# -> ask docker to take a snapshot of the container and save as an image with name 'figlet' and a tag as 'build', i.e. `figlet:build`
$ docker commit c734 figlet:build
$ docker commit c734 figlet:build
sha256:e12731d81e20a3658c6aec52cb741f24109607f8ceac29491c0ed4fcee40daa1
$ docker images
REPOSITORY   TAG       IMAGE ID       CREATED         SIZE
figlet       build     e12731d81e20   3 seconds ago   216MB
ubuntu       latest    c35e29c94501   2 months ago    117MB
~ docker images
REPOSITORY   TAG       IMAGE ID       CREATED         SIZE
figlet       build     e12731d81e20   3 seconds ago   216MB
ubuntu       latest    c35e29c94501   2 months ago    117MB
$ docker run -it figlet:build
root@01e800e9e293:/# figlet i love you smruti
 _   _                                                               _   _
(_) | | _____   _____   _   _  ___  _   _   ___ _ __ ___  _ __ _   _| |_(_)
| | | |/ _ \ \ / / _ \ | | | |/ _ \| | | | / __| '_ ` _ \| '__| | | | __| |
| | | | (_) \ V /  __/ | |_| | (_) | |_| | \__ \ | | | | | |  | |_| | |_| |
|_| |_|\___/ \_/ \___|  \__, |\___/ \__,_| |___/_| |_| |_|_|   \__,_|\__|_|
                        |___/
```

- now, instead of again doing another `docker commit` to add some more changes on my container, before doing `docker commit` we can also do `docker diff` to see the differences between the image and the currently running container based on that image
	- i.e. it shows all the files that were added or removed or changed, between the image and container
```
root@01e800e9e293:/# exit
$ docker diff c734
C /root
A /root/.bash_history
C /etc
C /etc/alternatives
A /etc/alternatives/figlet
A /etc/emacs
A /etc/emacs/site-start.d
A /etc/emacs/site-start.d/50figlet.el
C /usr
C /usr/bin
A /usr/bin/showfigfonts
....
```

#### use cases for `docker diff`
- after running a script without knowing what changes it brought in
- security use-case
	- attackers have end goal of taking over the web server, by running some arbitrary code (**arbitrary execution vulnerability**)
	- there are a whole bunch of attacks on web services that happen in two steps, 
		- the first step is to finding a way to upload the code to the server
		- and then to execute that code in some way
	- with docker diff, we can detect if there are some files that got added in our container
- running read only containers
```
 $ docker run -it --read-only ubuntu
root@ef8e8485b60c:/# apt-get update
Reading package lists... Done
E: List directory /var/lib/apt/lists/partial is missing. - Acquire (30: Read-only file system)
root@ef8e8485b60c:/#
```
- read-only containers and `docker diff`. What's the link?
	- many services and programs that would not work in `read-only` mode
	- thus, running `docker diff` once in a while allows to see the files being created
- `docker diff` checks at the read-write layer of the container

- thus, `docker commit` is pretty manual way of creating images
- thus, use `docker build <docker-file>` with `Dockerfile`

### Dockerfile
- a recipe indicating how to build a container image
#### Steps:
1. Pretend writing a shell script
```Dockerfile
apt-get update
apt-get install figlet
```
2. add `RUN` in front of every command/line
```Dockerfile
RUN apt-get update
RUN apt-get install figlet
```
3. add base image name
```Dockerfile
FROM ubuntu
RUN apt-get update
RUN apt-get install figlet
```
4. `docker build` command 
```
$ docker build . -t figlet:dockerfile
```
- `docker build` command uses build cache to avoid re-building entire image on every time while we run `docker build .`
- we can get over cache using `--no-cache` flag
- we can also use `--pull` flag to update the base image
	- `docker build . -t figlet:dockerfile --pull`
	- if base image has newer version on pull, it rebuilds the next steps, discarding the cache for the next steps as well

#### Dockerfile helps in staying idempotent even when building multiple times
- writing a shell script for creating images is also an option though, but it need not be **idempotent**
	- idempotent meaning, applying a function any number of times, the result would not change
	- multiplying by 0 is idempotent
	- some commands are idempotent, like `apt-get update`
	- some commands are not **idempotent** like creating a user `useradd smruti`
		- to make it idempotent, in the shell script we should check if the user already exists and proceed creating only if not found

### What happens when we build the image?
- running `docker build` on a mac and windows machine will have different outputs. Why?
	- on docker, we have two generations of build systems
	- we have old style builder 
		- e.g. `docker build . -t figlet:dockerfile`
	- new style builder called `BUILDKIT`
		- `DOCKER_BUILDKIT=1 docker build . -t figlet:dockerfile`
		- `BUILDKIT` is a completely new re-implementation of the build system
		- feels almost same
		- ~ to going from a gas powered car to an elective vehicle (EV)
	- why different on mac and windows
		- while in mac or windows, by default it will use new `buildkit` builder
		- on linux, by default it will use old style
			- can enable it by adding environment variable `DOCKER_BUILDKIT=1` before running `docker build .`

### Docker history command
- shows layers of an image
```sh
 docker-ardan-labs docker build . -t figlet:dockerfile
[+] Building 16.1s (7/7) FINISHED                                                                  docker:desktop-linux
 => [internal] load build definition from Dockerfile                                                               0.0s
 => => transferring dockerfile: 132B                                                                               0.0s
 => [internal] load metadata for docker.io/library/ubuntu:latest                                                   0.0s
 => [internal] load .dockerignore                                                                                  0.0s
 => => transferring context: 2B                                                                                    0.0s
 => [1/3] FROM docker.io/library/ubuntu:latest@sha256:c35e29c9450151419d9448b0fd75374fec4fff364a27f176fb458d472df  2.3s
 => => resolve docker.io/library/ubuntu:latest@sha256:c35e29c9450151419d9448b0fd75374fec4fff364a27f176fb458d472df  2.3s
 => [2/3] RUN apt-get update                                                                                       7.2s
 => [3/3] RUN apt-get install figlet                                                                               3.7s
 => exporting to image                                                                                             2.8s
 => => exporting layers                                                                                            2.2s
 => => exporting manifest sha256:e8675194548218ccfd4c7fad61c6c0d929e56f3cf10fc635c88b2a41bd15e18f                  0.0s
 => => exporting config sha256:dd27f1127402eaa7a3917666a15fa812f74b8668df45be78df92b0810c346f73                    0.0s
 => => exporting attestation manifest sha256:2e34fd6fc75271c72d574ae02a2088cee1f05ea4c7275a0b21c8b1496d6b6e93      0.0s
 => => exporting manifest list sha256:cf39075dfcc7fb75533de43e18f90772b19685d763c89b9f6198c9357b606926             0.0s
 => => naming to docker.io/library/figlet:dockerfile                                                               0.0s
 => => unpacking to docker.io/library/figlet:dockerfile                                                            0.5s
➜  docker-ardan-labs docker history figlet:dockerfile
IMAGE          CREATED          CREATED BY                                      SIZE      COMMENT
cf39075dfcc7   18 seconds ago   RUN /bin/sh -c apt-get install figlet # buil…   1.34MB    buildkit.dockerfile.v0
<missing>      21 seconds ago   RUN /bin/sh -c apt-get update # buildkit        57.5MB    buildkit.dockerfile.v0
<missing>      2 months ago     /bin/sh -c #(nop)  CMD ["/bin/bash"]            0B
<missing>      2 months ago     /bin/sh -c #(nop) ADD file:ddf1aa62235de6657…   87.6MB
<missing>      2 months ago     /bin/sh -c #(nop)  LABEL org.opencontainers.…   0B
<missing>      2 months ago     /bin/sh -c #(nop)  LABEL org.opencontainers.…   0B
<missing>      2 months ago     /bin/sh -c #(nop)  ARG LAUNCHPAD_BUILD_ARCH     0B
<missing>      2 months ago     /bin/sh -c #(nop)  ARG RELEASE                  0B

$ docker-ardan-labs cat Dockerfile
FROM ubuntu
RUN apt-get update
RUN apt-get install figlet
```

#### sh -c: what is CREATED BY `/bin/sh -c apt-get install figlet`
- when running something on an UNIX system, like `ls -l /tmp`
	- the shell does a system call which conceptually would look like: `exec("ls", "-l", "/tmp")`
- `ls ~/docs` => `"ls" "/home/govind/docs"` 
- the shell (`sh`) is responsible for this conversion or parsing or transformation of the commands to system call inputs
- In Dockerfile, when we have `RUN apt-get install figlet`
	- when need something to transform or parse it into `"apt-get"` `"install"` and `"figlet"`
	- easiest way is to offload that job to shell
		- `"sh" "-c" "apt-get install figlet"`
		- "-c" is for "command" indicating shell to execute the following command
		- So, Docker doesn't know how to parse and interpret those commands and thus offload that job to shell
- asking docker to skip shell: use \[""]
```Dockerfile
...
RUN ["apt-get", "install", "figlet"]
...
```
- this tells docker builder that it need not further parse it by reaching the shell and doing the `sh -c`


### `CMD` and `ENTRYPOINT`

#### `CMD`
- let's us define what we want to execute on starting the container
```Dockerfile
...
CMD figlet hello world
...
```
#### `ENTRYPOINT`
- very similar to `CMD`
- let's us define a base command to be executed inside the container
```Dockerfile
...
ENTRYPOINT ["figlet", "-f", "script", "Hello, "]
```
- Running this:
```shell
 docker-ardan-labs docker run -it figlet:json
 ,          _   _
/|   |     | | | |
 |___|  _  | | | |  __
 |   |\|/  |/  |/  /  \_
 |   |/|__/|__/|__/\__/o
                       /

➜  docker-ardan-labs docker run -it figlet:json smruti
 ,          _   _
/|   |     | | | |                                         o
 |___|  _  | | | |  __       ,   _  _  _    ,_         _|_
 |   |\|/  |/  |/  /  \_    / \_/ |/ |/ |  /  |  |   |  |  |
 |   |/|__/|__/|__/\__/o     \/   |  |  |_/   |_/ \_/|_/|_/|_/
                       /

➜  docker-ardan-labs docker run -it figlet:json SMRUTI
 ,          _   _                ,__ __   , __  _        ______ _
/|   |     | | | |            ()/|  |  | /|/  \(_|    | (_) |  | |
 |___|  _  | | | |  __        /\ |  |  |  |___/  |    |     |  | |
 |   |\|/  |/  |/  /  \_     /  \|  |  |  | \    |    |   _ |_ |/
 |   |/|__/|__/|__/\__/o    /(__/|  |  |_/|  \_/  \__/\_/(_/ \_/\/
                       /
```


#### Having multiple `ENTRYPOINT` and `CMD` statements inside Dockerfile
- Multiple `CMD` in a Dockerfile doesn't work
- Each new `CMD` line replaces the earlier `CMD` line
- thus, only the last `CMD` sustains and gets executed on starting the container
