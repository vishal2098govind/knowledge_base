#docker #multi-stage-builds

### Can't we remove superfluous files with `RUN`
What happens if we do one of the following commands?
- `RUN rm -rf ...`
- `RUN apt-get remove ...`
- `RUN make clean ...`

#### Removing files with an extra layer
- when downloading an image, all layers must be downloaded
- Thus, `RUN rm` does not reduce the size of the image or free up disk space

### Collapsing Layers
We frequently would see Dockerfiles like this:
```Dockerfile
FROM ubuntu
RUN apt-get update && apt-get install xxx && ... && apt-get remove xxx && ...
```
OR the more readable variant:
```Dockerfile
FROM ubuntu
RUN apt-get update \
	&& apt-get install <compiler> \
	&& ... \
	&& apt-get remove <compiler> \
	&& ...
```
- This `RUN` command gives a  single layer.
- The files that are added, then removed in the same layer do not grow the layer size.
- So, this works but, ruins the cache as any change in program, needs it to re-download (`apt-get`) the compiler , re-compile etc
- thus, size of image is less but at the cost of longer build time
- also not very readable

### Build binaries outside of Dockerfile
```Dockerfile
FROM ubuntu
COPY xxx /usr/local/bin
```
- requires extra build tools to be build outside of docker
- back in dependency hell and "works on my machine"
- breaks portability across different platforms
- grows code repository size if binary is updated frequently

### Multi Stage Builds
It's a good solution to reduce the image size
Multi-stage builds allows us to have multiple stages
Each stage is a separate image, and can copy files from **previous stages**

The idea of a multi stage build is that at the minimum, we have 
	a build stage - build or compile the program
	and a run stage - just take the compiled program and run

```Dockerfile
# stage 0: compile
FROM ubuntu

RUN apt-get update
RUN apt-get install build-essential -y

WORKDIR app

COPY . .

RUN gcc -o hello hello.c

# stage 1: run
FROM ubuntu
COPY --from=0 /app/hello .
CMD ./hello
```
- multi-stage docker file with stage aliases:
```Dockerfile
# stage 0: build
FROM ubuntu AS build

RUN apt-get update
RUN apt-get install build-essential -y

WORKDIR app

COPY . .

RUN gcc -o hello hello.c

# stage 1: run
FROM ubuntu
COPY --from=build /app/hello .
CMD ./hello
```
- often, using alpine versions of images reduces size
- e.g. using `FROM python:alpine` or `FROM alpine` etc