#docker

- This concept also answers "How can we invalidate build cache selectively for particular line in Dockerfile"
- Consider a C program:
```c
int main() {
	puts("Helloworld");
	
	return 0;
}
```
- if we want to put this program in a container, the thing with languages like C or Go or Rust or Java etc, these are compiled languages
- we can't just get the code, put into container and call it a day. Need to compile first.
- Ideally we want to compile the code within the container. Thus:
```Dockerfile
FROM ubuntu

# install C compiler
RUN apt-get update
RUN apt-get install build-essential -y

# copy hello.c to the container
# <WARNING: do not do this>
# if our hello.c is in a remote repository https://github.com/vishal2098govind/hello-c
RUN apt-get install git
RUN git clone https://github.com/vishal2098govind/hello-c
# then cd into that and further compile hello.c in that repository.
# but doing this would have us always commit and push to the remote repository and everytime we want to build docker image, we need to use --no-cache flag
# this would work but
# thus, use COPY
# </WARNING>

# cp src dest
COPY hello.c .

# compile hello.c to an executable
RUN gcc -o hello hello.c

# run the executable
CMD ./hello
```
- the Docker builder handles `COPY` specially where it will check if the files have changed to invalidate or consider build cache

### Changing working directories (`WORKDIR`)
- we might do think of doing something like:
```Dockerfile
FROM ubuntu

RUN apt-get update
RUN apt-get install build-essential -y

RUN mkdir /app
RUN cd /app

COPY hello.c

RUN gcc -o hello hello.c

CMD ./hello
```
- but this will not work as expected
```Dockerfile
# this does create app directory at the root as mkdir brings in some changes in the container, thus this will be saved in the image as well
RUN mkdir /app

# however, doing cd doesn't bring any change in the container, and thus would not be saved in the image
RUN cd /app
# thus, instead use WORKDIR. doing WORKDIR is enough, and we need not do a seperate `RUN mkdir`
WORKDIR /app
```
- while copying files, we can choose to copy file by file, or copy files in bulk (copy directories)
```Dockerfile
WORKDIR /app

# copy everything in current working directory from where we run docker build on the host machine, into the working directory of container (here set to /app in the line above)
COPY . .
```
- can also copy files matching a regex (or glob)
```Dockerfile
COPY *.c .
```
- can also copy everything except some specific files. use `.dockerignore` file
```dockerignore
Dockerfile
.dockerignore
```
- In a docker file, we cannot specify paths or things outside of the build context
```Dockerfile
...
COPY /tmp/readme.txt .
...
```
```sh
$ vim /tmp/readme.txt
```
```vi
it's cool
```
```sh
$ docker build . -t hello
ERROR: COPY failed: file not found in build context or excluded by .dockerignore: stat tmp/readme.txt: file does. not exist
```
```sh
$ ls -l /tmp/readme.txt
-rw-r--r-- 1 docker users 20 Dec ...
```
- so, while copying, the files need to be present in the same directory from where we run `docker build` command