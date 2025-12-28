#k8s
![[Pasted image 20251228205029.png]]
To be able to create deployment and exposing `dockercoins`  with services, we would need to do
```sh
$ kubectl create deployment --image ???
```
- what would we give the image name? 
- when running on a k8s cluster we have multiple nodes and thus we must ship our images
- we ship images using registry
	- build images
	- push images to registry
	- then specify those images in the `kubectl create` commands
- we can use pre-built and pre-shipped images for dockercoins
	- `hasher` → `dockercoins/hasher:v0.1`
	- `redis` → `redis`
	- `rng` → `dockercoins/rng:v0.1`
	- `webui` → `dockercoins/webui:v0.1`
	- `worker` → `dockercoins/worker:v0.1`
- All services should be internal services, except the web UI
	- since we want to be able to connect to the web UI from outside
```sh
# deployments
$ kubectl create deployment --image  dockercoins/worker:v0.1
$ kubectl create deployment --image dockercoins/hasher:v0.1
$ kubectl create deployment --image dockercoins/rng:v0.1
$ kubectl create deployment --image dockercoins/webui:v0.1
$ kubectl create deployment --image redis

# services
$ kubectl expose deployment hasher --port 80
$ kubectl expose deployment rng --port 80
$ kubectl expose deployment redis --port 6379
$ kubectl expose deployment webui --port 80 --type LoadBalancer
```
