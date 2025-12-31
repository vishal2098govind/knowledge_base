#k8s

## Connecting to the services in k8s cluster
### `kubectl port-forward` - easiest option
- this works always, be it local cluster or remote cluster
```sh
$ kubectl port-forward service/redis 6379
Forwarding from 127.0.0.1:6379 -> 6379
Forwarding from [::1]:6379 -> 6379
Handling connection for 6379
```
```sh
$ redis-cli
127.0.0.1:6379> SET HELLO WORLD
OK
127.0.0.1:6379> GET HELLO
"WORLD"
```
- we can also serve other services
```sh
$ kubectl port-forward service/webui 80
Unable to listen on port 80: Listeners failed to create with the following errors: [unable to create listener: Error listen tcp4 127.0.0.1:80: bind: permission denied unable to create listener: Error listen tcp6 [::1]:80: bind: permission denied]
error: unable to listen on any of the requested ports: [{80 80}]
$ kubectl port-forward service/webui 1234:80 # now we can hit http://locahost:1234 in a browser and see the web UI of dockercoins
Forwarding from 127.0.0.1:1234 -> 80
Forwarding from [::1]:1234 -> 80
Handling connection for 1234
Handling connection for 1234
Handling connection for 1234
```
- port-forward works great with services, deployments, pods
- we can port-forward anything we want and k8s will try to make it work
- behind the scenes
	- if we port-forward to a pod, it connects to that pod
	- if we port-forward to a deployment or service or whatever else, its will use the `--selector` of that deployment or service or whatever else, to locate the pod that it needs to talk to
#### Connecting to a remote cluster locally
- if we can obtain the copy `~/.kube/config` file of the remote cluster, it is possible
- if the cluster.server address has a private IP address as the host, usually the case when it is from a cloud provider, then we have to replace the private IP address with the public IP address of the remote server as the cluster.server address
```sh
$ code kubeconfig.remote
# has the remote ~/.kube/config file
$ kubectl --kubeconfig kubeconfig.remote get nodes
# shows all the nodes in the remote cluster
$ kubectl --kubeconfig kubeconfig.remote port-forward service/webui 1234:80 # now we can hit http://locahost:1234 in a browser and see the web UI of the remote dockercoins
Forwarding from 127.0.0.1:1234 -> 80
Forwarding from [::1]:1234 -> 80
Handling connection for 1234
Handling connection for 1234
Handling connection for 1234
```
`kubectl port-forawrd` is the most common way to access internal services

### Checking connectivity to a bunch of services in a cluster
- If we want to check the health of a bunch of services
	- `rng`
	- hasher
	- `webui`
	- etc
- it would be a little bit annoying if we had to setup individual `port-forwards` for each of them
- a common option is to do 
```sh
# -it => interactive mode
# --rm => remove after exiting
 kubectl --kubeconfig kubeconfig.remote run testpod -it --rm --image alpine
If you don't see a command prompt, try pressing enter.
/ # apk add curl
( 1/10) Installing brotli-libs (1.2.0-r0)
( 2/10) Installing c-ares (1.34.6-r0)
( 3/10) Installing libunistring (1.4.1-r0)
( 4/10) Installing libidn2 (2.3.8-r0)
( 5/10) Installing nghttp2-libs (1.68.0-r0)
( 6/10) Installing nghttp3 (1.13.1-r0)
( 7/10) Installing libpsl (0.21.5-r3)
( 8/10) Installing zstd-libs (1.5.7-r2)
( 9/10) Installing libcurl (8.17.0-r1)
(10/10) Installing curl (8.17.0-r1)
Executing busybox-1.37.0-r30.trigger
OK: 13.2 MiB in 26 packages
/ # curl rng
RNG running on rng-65d885d498-vxb4l
/ # curl hasher
HASHER running on hasher-99bbd4bb-g7lln
/ # curl webui
Found. Redirecting to /index.html/ #
/ # curl worker
curl: (6) Could not resolve host: worker (Domain name not found)
/ #
```
- this is something like we are "in the place" so to speak
- like we are in the network and in the namespace and we can test a bunch of things from a convenient location