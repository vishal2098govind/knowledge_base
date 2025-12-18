#docker

## What is Docker?
- When we say "Docker", we refer to the "Docker Engine"
- Docker Engine is a daemon - service running in the background
- We control the Docker Engine with Docker CLI using docker command
- when we submit any docker command using Docker CLI, behind the scenes the Docker CLI will 
	- create an API request to the Docker Engine
	- Send that to Docker Engine
	- Get the response back from Docker Engine
	- Display the response back on the terminal

```zsh
➜  ~ docker version
Client:
 Version:           28.0.1
 API version:       1.48
 Go version:        go1.23.6
 Git commit:        068a01e
 Built:             Wed Feb 26 10:38:16 2025
 OS/Arch:           darwin/amd64 # <---------- this talks about the host machine on which Docker is running
 Context:           desktop-linux

Server: Docker Desktop 4.39.0 (184744)
 Engine:
  Version:          28.0.1
  API version:      1.48 (minimum version 1.24)
  Go version:       go1.23.6
  Git commit:       bbd0a17
  Built:            Wed Feb 26 10:41:16 2025
  OS/Arch:          linux/amd64 # <----------- this talks about the Docker Engine
  Experimental:     false
 containerd:
  Version:          1.7.25
  GitCommit:        bcc810d6b9066471b0b6fa75f557a15a1cbf31bb
 runc:
  Version:          1.2.4
  GitCommit:        v1.2.4-0-g6c52b3f
 docker-init:
  Version:          0.19.0
  GitCommit:        de40ad0
```
 - the client block in the output talks about the Docker CLI
 - the server block in the output talks about the Docker Engine 
 - The OS/Arch of "Docker Engine" is `linux/amd64` even when we run Docker on a Mac, because
	 - At the end of the day, **Docker only runs on Linux** (other than some special cases)
	 - How does it work on Windows/Mac?
		 - Docker VM is used