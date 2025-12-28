#k8s

- First things first: we cannot run a container
- we can only run a pod, which will have container

## `docker run` vs `kubectl run`
```sh
$ docker run [--name] [--detach] <image> [command...]
$ kubectl run <name> [--attach] <--image ...> [command] 
```
- `docker run`
	- name is optional
	- by default it runs attach mode
	- image name is mandatory
	- optional commands to run on the container
- `kubectl run`
	- name is mandatory
	- by default it runs in detach mode, thus specify `--attach` to attach
	- image is mandatory and needs `--image` flag
	- optional command to run on the pod
```sh
$ docker run --name pingpong --detach alpine ping localhost
$ kubectl run pingpong --image alpine ping localhost
pod/pingpong created
$ kubectl get pods
NAME       READY   STATUS    RESTARTS   AGE
pingpong   1/1     Running   0          42s
$ kubectl logs pingpong --tail 3 --follow
64 bytes from ::1: seq=121 ttl=64 time=0.045 ms
64 bytes from ::1: seq=122 ttl=64 time=0.044 ms
64 bytes from ::1: seq=123 ttl=64 time=0.052 ms
```

## Scaling
- `kubectl` gives a simple command to scale a workload:
- `kubectl scale TYPE NAME --replicas=HOWMANY`
- e.g.
	- `kubectl scale pod pingpong --replicas=3`
- there are many `kubectl` commands that will look like:
	- `kubectl VERB_ACTION TYPE_KIND_OF_OBJECT NAME_OF_OBJECT --FLAGS...`
```sh
$ kubectl scale pods pingpong
error: required flag(s) "replicas" not set
$ kubectl scale pods pingpong --replicas=3
Error from server (NotFound): the server could not find the requested resource
```
- to see the exact requests made by `kubectl`, use `-v6`
	- The `**-v=6**` flag in `kubectl` is used to display the specific API requests being made to the Kubernetes API server, rather than just general "verbose" logging in the traditional sense. 

	- The logging level system in Kubernetes uses a numeric scale, and each level corresponds to different types of information: 
		- **`--v=0`**: Generally useful information that should always be visible.
		- **`--v=2`**: The recommended default log level, providing useful steady-state information.
		- **`--v=4`**: Debug-level verbosity.
		- **`--v=6`**: Displays the **requested resources** and the actual API paths used, including the HTTP status codes and response times, which helps you see the underlying API calls for a given `kubectl` command.
		- **`--v=7`**: Displays HTTP request headers.
		- **`--v=8`**: Displays HTTP request contents (body).
		- **`--v=9`**: Displays the full HTTP request contents without truncation.
	- Therefore, using `-v=6` (or higher) is crucial for **debugging** as it allows administrators and developers to look "under the hood" and see exactly how `kubectl` is interacting with the Kubernetes API, which is helpful for troubleshooting authentication, authorization, or networking issues. The number '6' doesn't have a special meaning beyond its place in this specific, predefined log level hierarchy.
```sh
$ kubectl scale pods pingpong --replicas=3 -v6
I1228 15:32:02.413303   89322 loader.go:402] Config loaded from file:  /Users/govind/.kube/config
I1228 15:32:02.416868   89322 envvar.go:172] "Feature gate default state" feature="ClientsAllowCBOR" enabled=false
I1228 15:32:02.416887   89322 envvar.go:172] "Feature gate default state" feature="ClientsPreferCBOR" enabled=false
I1228 15:32:02.416893   89322 envvar.go:172] "Feature gate default state" feature="InformerResourceVersion" enabled=false
I1228 15:32:02.416898   89322 envvar.go:172] "Feature gate default state" feature="WatchListClient" enabled=false
I1228 15:32:02.464896   89322 round_trippers.go:560] GET https://127.0.0.1:6443/api/v1/namespaces/default/pods/pingpong 200 OK in 18 milliseconds
I1228 15:32:02.467543   89322 round_trippers.go:560] PATCH https://127.0.0.1:6443/api/v1/namespaces/default/pods/pingpong/scale 404 Not Found in 1 milliseconds
I1228 15:32:02.469304   89322 helpers.go:246] server response object: %s[{
  "metadata": {},
  "status": "Failure",
  "message": "the server could not find the requested resource",
  "reason": "NotFound",
  "details": {},
  "code": 404
}]
Error from server (NotFound): the server could not find the requested resource
```
- here, internally, `kubectl` first tried to hit `GET /api/v1/namespaces/default/pods/pingpong` which responded with 200 OK
- then `kubectl` tries to hit `PATCH /api/v1/namespaces/default/pods/pingpong/scale` which responded with 404 Not Found
- in `kubectl scale`, conceptually, we tell k8s, to go that object, there'd be a setting called "scale", get that scaling thing and put it to 3, but when it goes the it fails to find `/scale` resource on the pod
- in k8s, a pod is the basic unit of scaling
- when we want to scale, we don't put multiple containers in the pod, we put multiple pods
- so, if we want to scale the pingpong, we need to create multiple pods
- we can't go to one pod and say, scale that by creating multiple copies of itself
- we need to create multiple pods

### Creating multiple pods
```sh
$ kubectl run pingpong --image alpine ping localhost
Error from server (AlreadyExists): pods "pingpong" already exists
$ kubectl run pingpong2 --image alpine ping localhost
pod/pingpong2 created
$ kubectl run pingpong3 --image alpine ping localhost
pod/pingpong3 created
$ kubectl get pods
NAME        READY   STATUS    RESTARTS   AGE
pingpong    1/1     Running   0          24m
pingpong2   1/1     Running   0          72s
pingpong3   1/1     Running   0          67s
```
- it kind of works but not done in practice, not manually
- better way is to use `controller`
	- an umbrella above my pods or a baby-sitter taking care of my pods
- when we do `kubectl run pingpong --image alpine ping localhost`, the `kubectl` creates that manifest and sticks on the pin-board (`etcd`) in the control plane
- then, the controller is acting on that manifest
- thus, for scaling,  instead of creating a pod, let's create a `ReplicaSet` where we can mention replicas