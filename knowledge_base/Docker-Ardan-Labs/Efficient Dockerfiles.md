#docker
### Reduce number of layers
- each line in a `Dockerfile` creates a new layer
- Build your `Dockerfile` to take advantage of Docker's caching system
- Combine commands by using `&&` to continue commands and `\` to wrap lines
> it is frequent to build a `Dockerfile` line by line:
```Dockerfile
RUN apt-get install thisthing
RUN apt-get install andthatthing andthatotherone
RUN apt-get install somemorestuff
```
And then refactor it trivially before shipping:
```Dockerfile
RUN apt-get install thisthing andthatthing andthatotherone somemorestuff
```

###  be careful with `chown`, `chmod`, `mv`
 - say if we want to bring in a big file (can be like a ML model) into a container image
- say if `model.pck` is a big 100MB file and ubuntu image is around 78.8MB
```Dockerfile
FROM ubuntu
COPY . .
WORKDIR /data
RUN mv /model.pck .
```
- `docker build . -t ml`
- the `ml` image is around 288MB. Why? Why not 100+78.8 = ~179MB?
- this is because moving is simply removing at source + copying at destination
- during that, the docker doesn't know it already copied that file in the above `COPY` step
- problem if changing permissions in the docker file
```Dockerfile
FROM ubuntu
COPY . .
WORKDIR /data
RUN mv /model.pck .
RUN chmod 600 model.pck
```
- even for changing permissions, the whole file is being copied again.
- the size of this image now shows ~392MB
- thus, better to put files on the right place
- Instead of using `mv`, directly put files at the right place
- When extracting archives (tar, zip...) merge operations in a single layer.
- E.g.
```Dockerfile
RUN wget http://.../foo.tar.gz \
	&& tar -zxf foo.tar.gz \
	&& mv foo/fooctl /usr/local/bin \
	&& rm -rf foo foo.tar.gz
```