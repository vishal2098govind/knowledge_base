- `-P` flag while running a container
	- short for `--publish-all`
	- publish/expose all the ports of the container to the outside world
	- this let's docker decide which port to choose
```sh
$docker run -d -P nginx
293a3bc230378fe23a89b26bd065b247ce31616bb1593387984605fa0a949c24

 docker ps -a
CONTAINER ID   IMAGE          COMMAND                  CREATED          STATUS                      PORTS                   NAMES
293a3bc23037   nginx          "/docker-entrypoint.…"   18 seconds ago   Up 17 seconds               0.0.0.0:55001->80/tcp   gallant_cori
8de605178012   nginx          "/docker-entrypoint.…"   46 seconds ago   Exited (0) 26 seconds ago                           distracted_merkle
044a049f0f7b   db             "docker-entrypoint.s…"   36 minutes ago   Exited (0) 25 minutes ago                           epic_solomon
f8f03e2c7091   words          "java -Xmx8m -Xms8m …"   43 hours ago     Exited (129) 8 hours ago                            nifty_napier
992538a3049a   web            "./dispatcher"           2 days ago       Exited (2) 43 hours ago                             pedantic_napier
694ed0bca97f   6ebda73726e4   "/bin/sh -c ./dispat…"   2 days ago       Exited (137) 2 days ago                             crazy_roentgen
e3a2824f7200   6a39d807b7ec   "./dispatcher"           2 days ago       Exited (2) 2 days ago                               tender_kowalevski
3eb3f8768a06   2ef8594b5a08   "/bin/sh -c ./dispat…"   2 days ago       Exited (137) 2 days ago                             epic_gates
8c987640164f   figlet:cmd     "bash"                   2 days ago       Exited (0) 2 days ago                               peaceful_lalande
c646e8fd0d2e   figlet:cmd     "/bin/sh -c 'figlet …"   2 days ago       Exited (0) 2 days ago                               fervent_villani
5e100280cf3f   figlet:json    "figlet -f script 'H…"   2 days ago       Exited (0) 2 days ago                               laughing_mayer
621fd317b0a3   figlet:json    "figlet -f script 'H…"   2 days ago       Exited (0) 2 days ago                               silly_maxwell
9bac9460498c   figlet:json    "figlet -f script 'H…"   2 days ago       Exited (0) 2 days ago                               stoic_williams
ef8e8485b60c   ubuntu         "/bin/bash"              2 days ago       Exited (127) 2 days ago                             confident_mirzakhani
01e800e9e293   figlet:build   "/bin/bash"              2 days ago       Exited (127) 2 days ago                             exciting_galileo
c734a94943b5   ubuntu         "/bin/bash"              2 days ago       Exited (0) 2 days ago                               nostalgic_poincare
```

- notice that only the nginx container has the PORTS column set as we launched the container with -P flag.
- 0.0.0.0:55001->80/tcp port mapping:
	- we can connect to the nginx web server within the container (on port 80) form the local machine via port 55001 (port mapping)
- Why docker chooses ports like 55001 by default and why not 80 directly?
	- so that multiple containers don't compete for a single port and can get some of the high port numbers
```
$ docker ps
CONTAINER ID   IMAGE     COMMAND                  CREATED         STATUS         PORTS                   NAMES
293a3bc23037   nginx     "/docker-entrypoint.…"   7 minutes ago   Up 7 minutes   0.0.0.0:55001->80/tcp   gallant_cori
➜  db git:(main) ✗ docker run -d -P nginx
f6cf45485ab52e3d4bbe3677a9a030d27e89d4ec1829c3277ce511bc7df48086
➜  db git:(main) ✗ docker run -d -P nginx
6aff8e37c57e1985333b52c285fd59ca8dfdb7a3b8531f1f5455fe01366d7047
➜  db git:(main) ✗ docker run -d -P nginx
492e44ee91ce35be095a8c4d0e7d65cefae4466925b64ac3b17f1247a17df95f
$ docker ps
CONTAINER ID   IMAGE     COMMAND                  CREATED         STATUS         PORTS                   NAMES
492e44ee91ce   nginx     "/docker-entrypoint.…"   2 seconds ago   Up 2 seconds   0.0.0.0:55004->80/tcp   gallant_raman
6aff8e37c57e   nginx     "/docker-entrypoint.…"   3 seconds ago   Up 2 seconds   0.0.0.0:55003->80/tcp   sweet_feynman
f6cf45485ab5   nginx     "/docker-entrypoint.…"   4 seconds ago   Up 4 seconds   0.0.0.0:55002->80/tcp   keen_lumiere
293a3bc23037   nginx     "/docker-entrypoint.…"   9 minutes ago   Up 9 minutes   0.0.0.0:55001->80/tcp   gallant_cori
 
```
- to specify port mapping:
```
$ docker run -d -p 1234:80 nginx
eb0f7cebaa95164f6da0b394bf7d81367eea46c7937a530482882829e454d852
$ docker ps                     
CONTAINER ID   IMAGE     COMMAND                  CREATED              STATUS              PORTS                   NAMES
eb0f7cebaa95   nginx     "/docker-entrypoint.…"   2 seconds ago        Up 1 second         0.0.0.0:1234->80/tcp    funny_beaver
492e44ee91ce   nginx     "/docker-entrypoint.…"   About a minute ago   Up About a minute   0.0.0.0:55004->80/tcp   gallant_raman
6aff8e37c57e   nginx     "/docker-entrypoint.…"   About a minute ago   Up About a minute   0.0.0.0:55003->80/tcp   sweet_feynman
f6cf45485ab5   nginx     "/docker-entrypoint.…"   About a minute ago   Up About a minute   0.0.0.0:55002->80/tcp   keen_lumiere
293a3bc23037   nginx     "/docker-entrypoint.…"   11 minutes ago       Up 11 minutes       0.0.0.0:55001->80/tcp   gallant_cori
```

- There are many ways to integrate containers in our network:
	- start container, letting docker allocate a public port for it - use `-P` flag. then retrieve that port number and feed it to our config
	- pick a fixed port number in advance, when we generate our config. start container specifying the mapping manually. use `-p` flag
	- use an orchestrator like Kubernetes or swarm. the orchestrator will provide its own networking facilities. - most  used
	- orchestrator typically provide mechanisms to enable direct container-to-container communication across hosts and publishing/load balancing for inbound traffic


- Containers have IP addresses
```sh
➜  db git:(main) ✗ docker ps
CONTAINER ID   IMAGE     COMMAND                  CREATED             STATUS             PORTS                   NAMES
eb0f7cebaa95   nginx     "/docker-entrypoint.…"   About an hour ago   Up About an hour   0.0.0.0:1234->80/tcp    funny_beaver
492e44ee91ce   nginx     "/docker-entrypoint.…"   About an hour ago   Up About an hour   0.0.0.0:55004->80/tcp   gallant_raman
6aff8e37c57e   nginx     "/docker-entrypoint.…"   About an hour ago   Up About an hour   0.0.0.0:55003->80/tcp   sweet_feynman
f6cf45485ab5   nginx     "/docker-entrypoint.…"   About an hour ago   Up About an hour   0.0.0.0:55002->80/tcp   keen_lumiere
293a3bc23037   nginx     "/docker-entrypoint.…"   About an hour ago   Up About an hour   0.0.0.0:55001->80/tcp   gallant_cori
➜  db git:(main) ✗ docker inspect eb0
[
    {
        "Id": "eb0f7cebaa95164f6da0b394bf7d81367eea46c7937a530482882829e454d852",
        "Created": "2025-12-22T13:30:04.147184113Z",
        "Path": "/docker-entrypoint.sh",
        "Args": [
            "nginx",
            "-g",
            "daemon off;"
        ],
        "State": {
            "Status": "running",
            "Running": true,
            "Paused": false,
            "Restarting": false,
            "OOMKilled": false,
            "Dead": false,
            "Pid": 1213,
            "ExitCode": 0,
            "Error": "",
            "StartedAt": "2025-12-22T13:30:04.192530475Z",
            "FinishedAt": "0001-01-01T00:00:00Z"
        },
        "Image": "sha256:fb01117203ff38c2f9af91db1a7409459182a37c87cced5cb442d1d8fcc66d19",
        "ResolvConfPath": "/var/lib/docker/containers/eb0f7cebaa95164f6da0b394bf7d81367eea46c7937a530482882829e454d852/resolv.conf",
        "HostnamePath": "/var/lib/docker/containers/eb0f7cebaa95164f6da0b394bf7d81367eea46c7937a530482882829e454d852/hostname",
        "HostsPath": "/var/lib/docker/containers/eb0f7cebaa95164f6da0b394bf7d81367eea46c7937a530482882829e454d852/hosts",
        "LogPath": "/var/lib/docker/containers/eb0f7cebaa95164f6da0b394bf7d81367eea46c7937a530482882829e454d852/eb0f7cebaa95164f6da0b394bf7d81367eea46c7937a530482882829e454d852-json.log",
        "Name": "/funny_beaver",
        "RestartCount": 0,
        "Driver": "overlayfs",
        "Platform": "linux",
        "MountLabel": "",
        "ProcessLabel": "",
        "AppArmorProfile": "",
        "ExecIDs": null,
        "HostConfig": {
            "Binds": null,
            "ContainerIDFile": "",
            "LogConfig": {
                "Type": "json-file",
                "Config": {}
            },
            "NetworkMode": "bridge",
            "PortBindings": {
                "80/tcp": [
                    {
                        "HostIp": "",
                        "HostPort": "1234"
                    }
                ]
            },
            "RestartPolicy": {
                "Name": "no",
                "MaximumRetryCount": 0
            },
            "AutoRemove": false,
            "VolumeDriver": "",
            "VolumesFrom": null,
            "ConsoleSize": [
                14,
                158
            ],
            "CapAdd": null,
            "CapDrop": null,
            "CgroupnsMode": "private",
            "Dns": [],
            "DnsOptions": [],
            "DnsSearch": [],
            "ExtraHosts": null,
            "GroupAdd": null,
            "IpcMode": "private",
            "Cgroup": "",
            "Links": null,
            "OomScoreAdj": 0,
            "PidMode": "",
            "Privileged": false,
            "PublishAllPorts": false,
            "ReadonlyRootfs": false,
            "SecurityOpt": null,
            "UTSMode": "",
            "UsernsMode": "",
            "ShmSize": 67108864,
            "Runtime": "runc",
            "Isolation": "",
            "CpuShares": 0,
            "Memory": 0,
            "NanoCpus": 0,
            "CgroupParent": "",
            "BlkioWeight": 0,
            "BlkioWeightDevice": [],
            "BlkioDeviceReadBps": [],
            "BlkioDeviceWriteBps": [],
            "BlkioDeviceReadIOps": [],
            "BlkioDeviceWriteIOps": [],
            "CpuPeriod": 0,
            "CpuQuota": 0,
            "CpuRealtimePeriod": 0,
            "CpuRealtimeRuntime": 0,
            "CpusetCpus": "",
            "CpusetMems": "",
            "Devices": [],
            "DeviceCgroupRules": null,
            "DeviceRequests": null,
            "MemoryReservation": 0,
            "MemorySwap": 0,
            "MemorySwappiness": null,
            "OomKillDisable": null,
            "PidsLimit": null,
            "Ulimits": [],
            "CpuCount": 0,
            "CpuPercent": 0,
            "IOMaximumIOps": 0,
            "IOMaximumBandwidth": 0,
            "MaskedPaths": [
                "/proc/asound",
                "/proc/acpi",
                "/proc/kcore",
                "/proc/keys",
                "/proc/latency_stats",
                "/proc/timer_list",
                "/proc/timer_stats",
                "/proc/sched_debug",
                "/proc/scsi",
                "/sys/firmware",
                "/sys/devices/virtual/powercap"
            ],
            "ReadonlyPaths": [
                "/proc/bus",
                "/proc/fs",
                "/proc/irq",
                "/proc/sys",
                "/proc/sysrq-trigger"
            ]
        },
        "GraphDriver": {
            "Data": null,
            "Name": "overlayfs"
        },
        "Mounts": [],
        "Config": {
            "Hostname": "eb0f7cebaa95",
            "Domainname": "",
            "User": "",
            "AttachStdin": false,
            "AttachStdout": false,
            "AttachStderr": false,
            "ExposedPorts": {
                "80/tcp": {}
            },
            "Tty": false,
            "OpenStdin": false,
            "StdinOnce": false,
            "Env": [
                "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
                "NGINX_VERSION=1.29.4",
                "NJS_VERSION=0.9.4",
                "NJS_RELEASE=1~trixie",
                "PKG_RELEASE=1~trixie",
                "DYNPKG_RELEASE=1~trixie"
            ],
            "Cmd": [
                "nginx",
                "-g",
                "daemon off;"
            ],
            "Image": "nginx",
            "Volumes": null,
            "WorkingDir": "",
            "Entrypoint": [
                "/docker-entrypoint.sh"
            ],
            "OnBuild": null,
            "Labels": {
                "maintainer": "NGINX Docker Maintainers \u003cdocker-maint@nginx.com\u003e"
            },
            "StopSignal": "SIGQUIT"
        },
        "NetworkSettings": {
            "Bridge": "",
            "SandboxID": "9f172f79a1bf6348efd086647eb2f1b79346fe8e6c470e743e00e6e8977da1d3",
            "SandboxKey": "/var/run/docker/netns/9f172f79a1bf",
            "Ports": {
                "80/tcp": [
                    {
                        "HostIp": "0.0.0.0",
                        "HostPort": "1234"
                    }
                ]
            },
            "HairpinMode": false,
            "LinkLocalIPv6Address": "",
            "LinkLocalIPv6PrefixLen": 0,
            "SecondaryIPAddresses": null,
            "SecondaryIPv6Addresses": null,
            "EndpointID": "ae978d9e672e3ed0cfaeafc77eabd97c3cf3c510a9d7437bdd8ce2b9ecd37e2a",
            "Gateway": "172.17.0.1",
            "GlobalIPv6Address": "",
            "GlobalIPv6PrefixLen": 0,
            "IPAddress": "172.17.0.6",
            "IPPrefixLen": 16,
            "IPv6Gateway": "",
            "MacAddress": "36:39:1e:f7:bc:e9",
            "Networks": {
                "bridge": {
                    "IPAMConfig": null,
                    "Links": null,
                    "Aliases": null,
                    "MacAddress": "36:39:1e:f7:bc:e9",
                    "DriverOpts": null,
                    "GwPriority": 0,
                    "NetworkID": "dd41f0e4f02d29731728dc0c40207f954c316d46ff25c77ea71ee085abe4d791",
                    "EndpointID": "ae978d9e672e3ed0cfaeafc77eabd97c3cf3c510a9d7437bdd8ce2b9ecd37e2a",
                    "Gateway": "172.17.0.1",
                    "IPAddress": "172.17.0.6",
                    "IPPrefixLen": 16,
                    "IPv6Gateway": "",
                    "GlobalIPv6Address": "",
                    "GlobalIPv6PrefixLen": 0,
                    "DNSNames": null
                }
            }
        },
        "ImageManifestDescriptor": {
            "mediaType": "application/vnd.oci.image.manifest.v1+json",
            "digest": "sha256:547c8c6863a88abd0c987779413489ff0e9a693c24d2d88ce9eb8515e0ae0335",
            "size": 2290,
            "annotations": {
                "com.docker.official-images.bashbrew.arch": "amd64",
                "org.opencontainers.image.base.digest": "sha256:335d26b968695837d1566c442de2a4026049b26ecb8567be2b00f7a900c13637",
                "org.opencontainers.image.base.name": "debian:trixie-slim",
                "org.opencontainers.image.created": "2025-12-09T22:50:25Z",
                "org.opencontainers.image.revision": "afa829ae8cd9e25cf539cb03167dff1162f852cb",
                "org.opencontainers.image.source": "https://github.com/nginx/docker-nginx.git#afa829ae8cd9e25cf539cb03167dff1162f852cb:mainline/debian",
                "org.opencontainers.image.url": "https://hub.docker.com/_/nginx",
                "org.opencontainers.image.version": "1.29.4"
            },
            "platform": {
                "architecture": "amd64",
                "os": "linux"
            }
        }
    }
]
```
here, the ip address of the container with id starting with `eb0` is `172.17.0.6` (can be found in the result of `docker inspect eb0`)

```sh
$ ping 172.17.0.6
```


## Docker networking just uses the linux networking
- There is nothing new in the whole docker networking stack
- Docker uses linux bridges, ip tables and all.... which are bunch of things that docker didn't invent
- docker is not the first to use that
- linux has been using that for decades
- docker, instead of inventing something completely new, it uses pre-existing mechanisms and constructs and protocols etc