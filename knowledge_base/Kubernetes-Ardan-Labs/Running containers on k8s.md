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
	- by **default it runs attach** mode
	- image name is mandatory
	- optional commands to run on the container
- `kubectl run`
	- name is mandatory
	- by **default it runs in detach** mode, thus specify `--attach` to attach
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

## Creating Replica Sets
- How to create replica sets?
	- to create a replica set, either we have to write some YAML
	- or create **deployment**
		- because if we create a deployment, it will create a replica-set
		- and deployment can be created through CLI, unlike replica set
```sh
$ kubectl create deployment pingpong --image alpine -- ping localhost
```
> the `--` towards the end, before `ping` is basically to indicate the end of the options for `kubectl create` and beginning of command line happening within the container

```sh
$ watch kubectl get pods
NAME                        READY   STATUS    RESTARTS   AGE
pingpong                    1/1     Running   0          52m
pingpong-86959f6599-z2rtk   1/1     Running   0          12s
pingpong2                   1/1     Running   0          29m
pingpong3                   1/1     Running   0          29m
$ kubectl create deployment pingpong --image alpine -- ping localhost
deployment.apps/pingpong created
$ kubectl get all
NAME                            READY   STATUS    RESTARTS   AGE
pod/pingpong                    1/1     Running   0          54m
pod/pingpong-86959f6599-z2rtk   1/1     Running   0          2m40s
pod/pingpong2                   1/1     Running   0          31m
pod/pingpong3                   1/1     Running   0          31m

NAME                 TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
service/kubernetes   ClusterIP   10.96.0.1    <none>        443/TCP   67m

NAME                       READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/pingpong   1/1     1            1           2m40s

NAME                                  DESIRED   CURRENT   READY   AGE
replicaset.apps/pingpong-86959f6599   1         1         1       2m40s
```
- the deployment created the `replicaset`, then the `replicaset` created that `pod`
- to notice the hierarchy between the objects, look at the naming:
	- pod name: `pingpong-86959f6599-z2rtk`
		- `z2rtk` is the unique suffix created for this particular pod
	- "`pingpong-86959f6599`" part in the pod name is the name of the replica-set `replicaset.apps/pingpong-86959f6599`
	- "pingpong" part in the `replicaset` name is the name of the deployment `deployment.apps/pingpong`

### Scaling deployment
```sh
$ kubectl scale deployment pingpong --replicas=3
deployment.apps/pingpong scaled
$ kubectl get pods
NAME                        READY   STATUS    RESTARTS   AGE
pingpong                    1/1     Running   0          65m
pingpong-86959f6599-42b8p   1/1     Running   0          18s
pingpong-86959f6599-9bbmg   1/1     Running   0          18s
pingpong-86959f6599-z2rtk   1/1     Running   0          13m
pingpong2                   1/1     Running   0          43m
pingpong3                   1/1     Running   0          43m
$ kubectl get all
NAME                            READY   STATUS    RESTARTS   AGE
pod/pingpong                    1/1     Running   0          66m
pod/pingpong-86959f6599-42b8p   1/1     Running   0          69s
pod/pingpong-86959f6599-9bbmg   1/1     Running   0          69s
pod/pingpong-86959f6599-z2rtk   1/1     Running   0          14m
pod/pingpong2                   1/1     Running   0          43m
pod/pingpong3                   1/1     Running   0          43m

NAME                 TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
service/kubernetes   ClusterIP   10.96.0.1    <none>        443/TCP   79m

NAME                       READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/pingpong   3/3     3            3           14m

NAME                                  DESIRED   CURRENT   READY   AGE
replicaset.apps/pingpong-86959f6599   3         3         3       14m
```
- if we disrupt this replica set by deleting one of the pod
```sh
$ kubectl delete pod pingpong
pod "pingpong" deleted

^C%
$ kubectl get all
NAME                            READY   STATUS        RESTARTS   AGE
pod/pingpong                    1/1     Terminating   0          69m
pod/pingpong-86959f6599-42b8p   1/1     Running       0          3m31s
pod/pingpong-86959f6599-9bbmg   1/1     Running       0          3m31s
pod/pingpong-86959f6599-z2rtk   1/1     Running       0          17m
pod/pingpong2                   1/1     Running       0          46m
pod/pingpong3                   1/1     Running       0          46m

NAME                 TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
service/kubernetes   ClusterIP   10.96.0.1    <none>        443/TCP   81m

NAME                       READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/pingpong   3/3     3            3           17m

NAME                                  DESIRED   CURRENT   READY   AGE
replicaset.apps/pingpong-86959f6599   3         3         3       17m
```
- if we try to delete one of the pod in the replica set, we immediately see another new pod starting to run to maintain the replica set
- this is taken care by the replica-set controller
- in real world, if a node goes down for instance, it would take a bit longer, 30 seconds later, it would show unreachable, and 5 minutes later, eviction strategy is used and the pods are marked as having a problem and start a replacement pod

## Why we have these layers of deployments/replica-sets/pods?
- in the initial versions of k8s, we had a replication controller
- it was kind of a mix between deployment and replica set, to take care of scaling and rolling updates
- the **replica set** has one very simple mission - run N similar pods 
- to do **rolling updates**, we use **deployment**
	- when we make a change to the deployment, it's going to create a separate replica-set and use that to do the roll over
	- then scale down the old version of replica-set
	- scale up the new version of replica-set

## Troubleshooting error in pods
```sh
$ kubectl create deployment --image alpine -- ping localhost
$ kubectl get pods
NAME                        READY   STATUS             RESTARTS     AGE
pingpong-75c985598-sk84l    0/1     CrashLoopBackOff   1 (4s ago)   9s
pingpong-86959f6599-lthk5   1/1     Terminating        0            22s
$ kubectl get all
NAME                            READY   STATUS        RESTARTS      AGE
pod/pingpong-75c985598-sk84l    0/1     Error         2 (20s ago)   25s
pod/pingpong-86959f6599-lthk5   1/1     Terminating   0             38s

NAME                 TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
service/kubernetes   ClusterIP   10.96.0.1    <none>        443/TCP   108m

NAME                       READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/pingpong   0/1     1            0           25s

NAME                                 DESIRED   CURRENT   READY   AGE
replicaset.apps/pingpong-75c985598   1         1         0       25s
$ kubectl create deployment pingpong --image alpine -- ping
$ kubectl get pods
$ kubectl get all
$ kubectl describe pingpong
$ kubectl describe deployment pingpong
$ kubectl logs deployment pingpong
$ kubectl logs pingpong
$ kubectl logs pods pingpong
$ kubectl get pods
$ kubectl logs pods pingpong-75
$ kubectl logs pods pingpong-75c985598-sk84l
$ kubectl logs pingpong-75c985598-sk84l
$ kubectl logs pingpong
$ kubectl edit deployment pingpong
$ kubectl describe deployment pingpong
$ kubectl get pods
```