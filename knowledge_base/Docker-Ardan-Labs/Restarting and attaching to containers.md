#docker #docker-compose
- if we just stop a container using `docker stop`, it's still there in the machine, can be found using `docker ps -a`
- to remove them completely, we can use `docker rm`
- the containers need not necessarily take up a lot of disk space
- to know size occupied by the containers:
```sh
$ docker ps --all --size
```
- example:
```sh
➜  wordsmith docker ps -a -s
CONTAINER ID   IMAGE                 COMMAND                  CREATED        STATUS                      PORTS                    NAMES                   SIZE
a270b8c9ec95   client-1-db           "docker-entrypoint.s…"   10 hours ago   Up 10 hours                 5432/tcp                 client-1-db-1           20.5kB (virtual 481MB)
90e73d0b9ce4   client-1-web          "./dispatcher"           10 hours ago   Up 10 hours                 0.0.0.0:9001->80/tcp     client-1-web-1          4.1kB (virtual 1.04GB)
03f98ef11b8a   client-1-words        "java -Xmx8m -Xms8m …"   10 hours ago   Up 10 hours                                          client-1-words-1        49.2kB (virtual 476MB)
c8c5742ddb50   wordsmith-web         "./dispatcher"           10 hours ago   Exited (2) 10 hours ago                              wordsmith-web-2         4.1kB (virtual 1.04GB)
a68d974298f7   wordsmith-db          "docker-entrypoint.s…"   11 hours ago   Exited (0) 10 hours ago                              wordsmith-db-1          16.4kB (virtual 481MB)
ef1382045dc6   wordsmith-words       "java -Xmx8m -Xms8m …"   11 hours ago   Exited (143) 10 hours ago                            wordsmith-words-1       16.4kB (virtual 476MB)
44a03fc8aacb   wordsmith-new-db      "docker-entrypoint.s…"   11 hours ago   Up 11 hours                 5432/tcp                 wordsmith-new-db-1      20.5kB (virtual 481MB)
e3f399c308f3   wordsmith-new-web     "./dispatcher"           11 hours ago   Up 11 hours                 0.0.0.0:9002->80/tcp     wordsmith-new-web-1     12.3kB (virtual 1.04GB)
8c0765d8fbff   wordsmith-new-words   "java -Xmx8m -Xms8m …"   11 hours ago   Up 11 hours                                          wordsmith-new-words-1   49.2kB (virtual 476MB)
4a912fb4e9b1   dev-db                "docker-entrypoint.s…"   12 hours ago   Up 12 hours                 5432/tcp                 dev-db-1                20.5kB (virtual 481MB)
0f6c4a82cb33   dev-web               "./dispatcher"           12 hours ago   Created                                              dev-web-1               4.1kB (virtual 1.04GB)
7def84b6923e   dev-words             "java -Xmx8m -Xms8m …"   12 hours ago   Up 12 hours                                          dev-words-1             49.2kB (virtual 476MB)
8823bd1b0e9d   alpine                "/bin/sh"                35 hours ago   Exited (0) 35 hours ago                              naughty_mahavira        12.3kB (virtual 9.14MB)
1bd15e22ff2d   redis                 "docker-entrypoint.s…"   36 hours ago   Exited (255) 15 hours ago   6379/tcp                 prod-redis-1            4.1kB (virtual 149MB)
e58edac8daf4   prod-www              "python counter.py"      36 hours ago   Created                                              prod-www-1              4.1kB (virtual 79.5MB)
5121ea02af2f   redis                 "docker-entrypoint.s…"   36 hours ago   Exited (255) 15 hours ago   6379/tcp                 dev-redis-1             4.1kB (virtual 149MB)
b221cf679fa5   dev-www               "python counter.py"      36 hours ago   Exited (255) 15 hours ago   0.0.0.0:8000->5000/tcp   dev-www-1               4.1kB (virtual 79.5MB)
2da55f35030f   redis                 "docker-entrypoint.s…"   39 hours ago   Exited (0) 36 hours ago                              quirky_engelbart        4.1kB (virtual 149MB)
93fcba79fb8a   redis                 "docker-entrypoint.s…"   39 hours ago   Exited (0) 39 hours ago                              jolly_nash              4.1kB (virtual 149MB)
c70df7897490   trainingwheels        "/bin/sh -c 'gunicor…"   39 hours ago   Exited (0) 36 hours ago                              trusting_kare           975kB (virtual 80.4MB)
b0691cf8c903   trainingwheels        "/bin/sh -c 'gunicor…"   39 hours ago   Exited (0) 39 hours ago                              upbeat_mcnulty          975kB (virtual 80.4MB)
b7d03973268a   trainingwheels        "/bin/sh -c 'gunicor…"   39 hours ago   Created                                              sad_albattani           4.1kB (virtual 79.5MB)
728d2245965a   trainingwheels        "/bin/sh -c 'gunicor…"   39 hours ago   Exited (0) 39 hours ago                              infallible_goodall      975kB (virtual 80.4MB)
9b8d13ae5543   trainingwheels        "/bin/sh -c 'gunicor…"   39 hours ago   Exited (137) 39 hours ago                            competent_wu            975kB (virtual 80.4MB)
f072abee0993   trainingwheels        "/bin/sh -c 'gunicor…"   39 hours ago   Exited (0) 39 hours ago                              reverent_shamir         975kB (virtual 80.4MB)
9b1b76bd1641   alpine                "/bin/sh"                2 days ago     Exited (255) 40 hours ago                            strange_noether         12.3kB (virtual 9.14MB)
1b03e4fc728e   nginx                 "/docker-entrypoint.…"   2 days ago     Exited (255) 40 hours ago   80/tcp                   peaceful_raman          81.9kB (virtual 165MB)
15132d3ece64   nginx                 "/docker-entrypoint.…"   2 days ago     Exited (255) 40 hours ago   80/tcp                   beautiful_allen         81.9kB (virtual 165MB)
7fcbf03075a1   alpine                "/bin/sh"                2 days ago     Exited (255) 40 hours ago                            beautiful_blackwell     8.48MB (virtual 17.6MB)
91b258c64660   nginx                 "/docker-entrypoint.…"   2 days ago     Exited (255) 40 hours ago   80/tcp                   awesome_tu              81.9kB (virtual 165MB)
45b65a50342c   nginx                 "/docker-entrypoint.…"   2 days ago     Exited (255) 40 hours ago   80/tcp                   quirky_grothendieck     81.9kB (virtual 165MB)
df1b22032689   alpine                "/bin/sh"                2 days ago     Exited (0) 2 days ago                                nervous_raman           12.3kB (virtual 9.14MB)
7c35c51b64c4   alpine                "/bin/sh"                2 days ago     Exited (130) 2 days ago                              sad_kirch               12.3kB (virtual 9.14MB)
3dcf5581a0e1   nginx                 "/docker-entrypoint.…"   2 days ago     Exited (255) 40 hours ago   0.0.0.0:55005->80/tcp    hopeful_darwin          81.9kB (virtual 165MB)
eb0f7cebaa95   nginx                 "/docker-entrypoint.…"   2 days ago     Exited (255) 40 hours ago   0.0.0.0:1234->80/tcp     funny_beaver            81.9kB (virtual 165MB)
492e44ee91ce   nginx                 "/docker-entrypoint.…"   2 days ago     Exited (255) 40 hours ago   0.0.0.0:55004->80/tcp    gallant_raman           81.9kB (virtual 165MB)
6aff8e37c57e   nginx                 "/docker-entrypoint.…"   2 days ago     Exited (255) 40 hours ago   0.0.0.0:55003->80/tcp    sweet_feynman           81.9kB (virtual 165MB)
f6cf45485ab5   nginx                 "/docker-entrypoint.…"   2 days ago     Exited (255) 40 hours ago   0.0.0.0:55002->80/tcp    keen_lumiere            81.9kB (virtual 165MB)
293a3bc23037   nginx                 "/docker-entrypoint.…"   2 days ago     Exited (255) 40 hours ago   0.0.0.0:55001->80/tcp    gallant_cori            81.9kB (virtual 165MB)
8de605178012   nginx                 "/docker-entrypoint.…"   2 days ago     Exited (0) 2 days ago                                distracted_merkle       77.8kB (virtual 165MB)
044a049f0f7b   db                    "docker-entrypoint.s…"   2 days ago     Exited (0) 2 days ago                                epic_solomon            16.4kB (virtual 481MB)
f8f03e2c7091   words                 "java -Xmx8m -Xms8m …"   4 days ago     Exited (129) 3 days ago                              nifty_napier            16.4kB (virtual 476MB)
992538a3049a   web                   "./dispatcher"           4 days ago     Exited (2) 4 days ago                                pedantic_napier         4.1kB (virtual 1.04GB)
694ed0bca97f   6ebda73726e4          "/bin/sh -c ./dispat…"   4 days ago     Exited (137) 4 days ago                              crazy_roentgen          4.1kB (virtual 1.04GB)
e3a2824f7200   6a39d807b7ec          "./dispatcher"           4 days ago     Exited (2) 4 days ago                                tender_kowalevski       4.1kB (virtual 1.04GB)
3eb3f8768a06   2ef8594b5a08          "/bin/sh -c ./dispat…"   4 days ago     Exited (137) 4 days ago                              epic_gates              4.1kB (virtual 1.04GB)
8c987640164f   figlet:cmd            "bash"                   5 days ago     Exited (0) 4 days ago                                peaceful_lalande        12.3kB (virtual 146MB)
c646e8fd0d2e   figlet:cmd            "/bin/sh -c 'figlet …"   5 days ago     Exited (0) 5 days ago                                fervent_villani         4.1kB (virtual 146MB)
5e100280cf3f   figlet:json           "figlet -f script 'H…"   5 days ago     Exited (0) 5 days ago                                laughing_mayer          4.1kB (virtual 146MB)
621fd317b0a3   figlet:json           "figlet -f script 'H…"   5 days ago     Exited (0) 5 days ago                                silly_maxwell           4.1kB (virtual 146MB)
9bac9460498c   figlet:json           "figlet -f script 'H…"   5 days ago     Exited (0) 5 days ago                                stoic_williams          4.1kB (virtual 146MB)
ef8e8485b60c   ubuntu                "/bin/bash"              5 days ago     Exited (127) 5 days ago                              confident_mirzakhani    4.1kB (virtual 87.6MB)
01e800e9e293   figlet:build          "/bin/bash"              5 days ago     Exited (127) 5 days ago                              exciting_galileo        16.4kB (virtual 146MB)
c734a94943b5   ubuntu                "/bin/bash"              5 days ago     Exited (0) 5 days ago                                nostalgic_poincare      59.4MB (virtual 147MB)
```
 - virtual = size of image
 - other size = size of container (read/write layer) + image

#### Recovering disk space:
```sh
$ docker system prune
```
- it removes all the stopped containers and all the dangling images

#### Starting old stopped container
```sh
$ docker start <container-name>
$ docker attach <container-name>
```
- docker run = docker create + docker start + docker attach

### Naming containers
- what's special about name: names are **unique** thus can't reuse names unless remove completely
- can set name using `--name` flag while running a new container
- compose names automatically

### Labels
- key-value pairs for labels
- rare to set labels ourselves
- compose sets labels automatically


### Getting inside a running container
```sh
$ docker exec <container-id> <command>
```
e.g.: only possible if the container has bash installed
```sh
$ docker exec -it <container-id> bash
root@<container-id>:#
```
- can use docker cp to copy or bring in the programs from local machine to a container inside
```sh
$ docker cp /bin/busybox <container-id>:/busybox
```
### Resource usage by containers
- to see instantaneous resource usage among containers:
```
$ docker stats
CONTAINER ID   NAME                    CPU %     MEM USAGE / LIMIT     MEM %     NET I/O           BLOCK I/O         PIDS
f57842307c96   client-1-web-1          0.00%     3.184MiB / 7.655GiB   0.04%     872B / 126B       0B / 0B           6
a270b8c9ec95   client-1-db-1           0.00%     23.74MiB / 7.655GiB   0.30%     1.46kB / 126B     0B / 418kB        9
03f98ef11b8a   client-1-words-1        0.15%     44.94MiB / 7.655GiB   0.57%     1.25kB / 126B     0B / 967kB        20
44a03fc8aacb   wordsmith-new-db-1      0.00%     24.21MiB / 7.655GiB   0.31%     4.78kB / 3.7kB    0B / 418kB        9
e3f399c308f3   wordsmith-new-web-1     0.00%     9.832MiB / 7.655GiB   0.13%     28.4kB / 318kB    2.87MB / 4.1kB    8
8c0765d8fbff   wordsmith-new-words-1   0.16%     53.96MiB / 7.655GiB   0.69%     7.44kB / 6.03kB   0B / 4.61MB       26
4a912fb4e9b1   dev-db-1                0.00%     28.71MiB / 7.655GiB   0.37%     1.91kB / 126B     4.1kB / 43.4MB    9
7def84b6923e   dev-words-1             0.22%     43.05MiB / 7.655GiB   0.55%     1.91kB / 126B     94.2kB / 5.41MB   23
```
- to see the resource usage **over time** and not just the resource usage at the moment, we can install a metric system on the **physical machine**,  like
	- prometheus
	- cacti
	- datadog

- for containers, we can monitor them from outside and need not install these metric systems on every container

### Q/A
- docker api doesn't have any notion of permissions or access control
	- what if we need access control?
	- option-1: use orchestration platform
		- in real world, it's rare to just have one giant machine running all of the containers
		- generally we have clusters with 100s or 1000s of machines, needing to have stuff like batch schedulers, to run workloads on throughout the cluster
	- option-2: if we really don't want orchestration and just have one machine, instead of letting people access directly the docker API, have them go through a wrapper between docker and user
- opposite of `docker cp <host-path> <container-id>:<container-path>` to bring from container to host:
	- `docker cp <container-id>:container-path <host-path>`