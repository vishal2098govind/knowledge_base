#k8s #volumes
- Volumes are special directories that are mounted in containers
- volumes can many different use cases

## Use cases of volumes
### Share data between multiple containers
- e.g. if we have a container that needs to work on some reference dataset
- and if we have two containers to do that
	- one to download the data
	- another to process the downloaded data
- that's a really common scenario
### Share data between container and the host machine
- e.g. if we have log files on the machine, we can use volumes to access these log files from a container
- this gives a way to run a logging agent in a container even though it's going to access the log files that are on the host machine

## Exposing configurations and credentials securely to containers

## Volumes != Persistent Volumes
- Volumes and Persisted Volumes are related, but are very different
- Volumes:
	- Volumes don't exist directly as APIs in k8s
	- we can't do `kubectl get volumes`
	- they exist in pod specifications
	- just like for containers, we can't do `kubectl get containers` as there are no resource types `containers`, however we have pods and within the pods we've containers
	- similarly volumes exists within pods
- Persistent volumes:
	- they exist directly as APIs in k8s
	- we can do `kubectl get persistentvolumes`
	- they represent the actual volume or the disk or the partitions or any kind of storage that we can consume on the cluster

## Adding a volume to a pod
### `nginx-1-without-volume`
```sh
$ cat k8s/nginx-1-without-volume.yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx-without-volume
spec:
  containers:
  - name: nginx
    image: nginx
$ kubectl create -f k8s/nginx-1-without-volume.yaml
pod/nginx-without-volume created
```
```sh
$  kubectl get pods --watch -o wide --output-watch-events
EVENT      NAME                   READY   STATUS    RESTARTS   AGE   IP       NODE     NOMINATED NODE   READINESS GATES
ADDED      nginx-without-volume   0/1     Pending   0          0s    <none>   <none>   <none>           <none>
MODIFIED   nginx-without-volume   0/1     Pending   0          0s    <none>   docker-desktop   <none>           <none>
MODIFIED   nginx-without-volume   0/1     ContainerCreating   0          1s    <none>   docker-desktop   <none>           <none>
MODIFIED   nginx-without-volume   1/1     Running             0          18s   10.1.1.68   docker-desktop   <none>           <none>
```
- now, if we hit the `nginx` pod:
```sh
kubectl port-forward pods/nginx-without-volume 8080:80
Forwarding from 127.0.0.1:8080 -> 80
Forwarding from [::1]:8080 -> 80
Handling connection for 8080
```
```sh
-----------------
$ curl localhost:8080
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
html { color-scheme: light dark; }
body { width: 35em; margin: 0 auto;
font-family: Tahoma, Verdana, Arial, sans-serif; }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>
```
- get into the `nginx-without-volume` pod
```sh
$ kubectl exec -it pods/nginx-without-volume -- bash
root@nginx-without-volume:/# cat usr/share/nginx/html/index.html
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
html { color-scheme: light dark; }
body { width: 35em; margin: 0 auto;
font-family: Tahoma, Verdana, Arial, sans-serif; }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>
root@nginx-without-volume:/# cd /usr/share/nginx/html/
root@nginx-without-volume:/usr/share/nginx/html# echo 'Vishal was here' > index.html
root@nginx-without-volume:/usr/share/nginx/html# exit
exit
$ curl localhost:8080
Vishal was here
```
### `nginx-with-volume`

```sh
$ cat k8s/nginx-2-with-volume.yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx-with-volume
spec:
  volumes:
  - name: www
  containers:
  - name: nginx
    image: nginx
    volumeMounts:
    - name: www
      mountPath: /usr/share/nginx/html/
$ kubectl apply -f k8s/nginx-2-with-volume.yaml
pod/nginx-with-volume created
$ kubectl get all
NAME                       READY   STATUS    RESTARTS   AGE
pod/nginx-with-volume      1/1     Running   0          6s
pod/nginx-without-volume   1/1     Running   0          8m45s
```
```sh
$ kubectl port-forward pods/nginx-with-volume 8081:80
Forwarding from 127.0.0.1:8081 -> 80
Forwarding from [::1]:8081 -> 80
Handling connection for 8081

```
- now, we will find empty directory of `/usr/share/nginx/html`
```sh
$ cat k8s/nginx-2-with-volume.yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx-with-volume
spec:
  volumes:
  - name: www
  containers:
  - name: nginx
    image: nginx
    volumeMounts:
    - name: www
      mountPath: /usr/share/nginx/html/
$ curl localhost:8081
<html>
<head><title>403 Forbidden</title></head>
<body>
<center><h1>403 Forbidden</h1></center>
<hr><center>nginx/1.29.4</center>
</body>
</html>
$ kubectl exec -it pods/nginx-with-volume -- bash
root@nginx-with-volume:/# cd /usr/share/nginx/html/
root@nginx-with-volume:/usr/share/nginx/html# cat index.html
cat: index.html: No such file or directory
root@nginx-with-volume:/usr/share/nginx/html# ls
root@nginx-with-volume:/usr/share/nginx/html#
```
- this is because we mounted a volume on top of that directory and the **default type** of volume that we mounted is called `EmptyDir`
```sh
$ echo "Vishal is here again" > index.html
root@nginx-with-volume:/usr/share/nginx/html# exit
exit
$ curl localhost:8081
Vishal is here again
```
- volumes are special directories, meaning that if we go on the node that runs the pod, we can find that volume

### `nginx-with-git`
- multiple containers in the same pod (within a node) using the same volume
```sh
$ cat k8s/nginx-3-with-git.yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx-with-git
spec:
  volumes:
  - name: www
  containers:
  - name: nginx
    image: nginx
    volumeMounts:
    - name: www
      mountPath: /usr/share/nginx/html/
  - name: git
    image: alpine
    command: [ "sh", "-c", "apk add git && git clone https://github.com/octocat/Spoon-Knife /www" ]
    volumeMounts:
    - name: www
      mountPath: /www/
  restartPolicy: OnFailure
```
- the `nginx` container is just serving the content within the volume
- the `git` container with `Spoon-Knife` is going to download (git pull) the content on the volume
	- once the git container finishes cloning, the container will exit
	- and the pod will be just left with `nginx` container serving the content of the volume
```sh
$ kubectl apply -f k8s/nginx-3-with-git.yaml
pod/nginx-with-git created
```
```sh
$ kubectl get pods --watch -o wide --output-watch-events
EVENT      NAME                   READY   STATUS    RESTARTS      AGE   IP           NODE             NOMINATED NODE   READINESS GATES
ADDED      nginx-with-volume      1/1     Running   1 (10h ago)   14h   10.1.1.167   docker-desktop   <none>           <none>
ADDED      nginx-without-volume   1/1     Running   1 (10h ago)   14h   10.1.1.151   docker-desktop   <none>           <none>
ADDED      nginx-with-git         0/2     Pending   0             0s    <none>       <none>           <none>           <none>
MODIFIED   nginx-with-git         0/2     Pending   0             0s    <none>       docker-desktop   <none>           <none>
MODIFIED   nginx-with-git         0/2     ContainerCreating   0             0s    <none>       docker-desktop   <none>           <none>
MODIFIED   nginx-with-git         2/2     Running             0             7s    10.1.1.180   docker-desktop   <none>           <none>
MODIFIED   nginx-with-git         1/2     NotReady            0             13s   10.1.1.180   docker-desktop   <none>           <none>

```
- here, the pod ends up showing NotReady because git container got killed after it's done cloning, but the pod is still expecting 2 containers
- thus, now if we put such a pod behind a service or Load Balancer, since it's not ready, it's not going to be added to the load balancer's endpoints
- we can
	- either use an init container - **recommended**
	- or make the service send the traffic to the pod even if the pod says it's not ready
		- setting `spec.publishNotReadyAddresses` to true in the service's YAML manifest
> `kubectl explain` is a way to view the documentation of the k8s API without leaving the terminal
> `kubectl explain` doesn't do any online fetching of documentation, rather it does introspection on the API. 
   Helpful if we are working in a specific version of k8s if they are not yet in the online documentation
```sh
$ kubectl explain service.spec
KIND:       Service
VERSION:    v1

FIELD: spec <ServiceSpec>


DESCRIPTION:
    Spec defines the behavior of a service.
    https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
    ServiceSpec describes the attributes that a user creates on a service.

FIELDS:
  ....
  ....
  publishNotReadyAddresses	<boolean>
    publishNotReadyAddresses indicates that any agent which deals with endpoints
    for this Service should disregard any indications of ready/not-ready. The
    primary use case for setting this field is for a StatefulSets Headless
    Service to propagate SRV DNS records for its Pods for the purpose of peer
    discovery. The Kubernetes controllers that generate Endpoints and
    EndpointSlice resources for Services interpret this to mean that all
    endpoints are considered "ready" even if the Pods themselves are not. Agents
    which consume only Kubernetes generated endpoints through the Endpoints or
    EndpointSlice resources can safely assume this behavior.
.....
.....
.....

```

### `nginx-with-init`
```sh
$ cat k8s/nginx-4-with-init.yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx-with-init
spec:
  volumes:
  - name: www
  containers:
  - name: nginx
    image: nginx
    volumeMounts:
    - name: www
      mountPath: /usr/share/nginx/html/
  initContainers:
  - name: git
    image: alpine
    command: [ "sh", "-c", "apk add git && sleep 5 && git clone https://github.com/octocat/Spoon-Knife /www" ]
    volumeMounts:
    - name: www
      mountPath: /www/

$ kubectl apply -f k8s/nginx-4-with-init.yaml
pod/nginx-with-init created
```

```
$ kubectl get pods --watch -o wide --output-watch-events
EVENT      NAME                   READY   STATUS    RESTARTS      AGE   IP           NODE             NOMINATED NODE   READINESS GATES
ADDED      nginx-with-volume      1/1     Running   1 (10h ago)   14h   10.1.1.167   docker-desktop   <none>           <none>
ADDED      nginx-without-volume   1/1     Running   1 (10h ago)   14h   10.1.1.151   docker-desktop   <none>           <none>
ADDED      nginx-with-git         0/2     Pending   0             0s    <none>       <none>           <none>           <none>
MODIFIED   nginx-with-git         0/2     Pending   0             0s    <none>       docker-desktop   <none>           <none>
MODIFIED   nginx-with-git         0/2     ContainerCreating   0             0s    <none>       docker-desktop   <none>           <none>
MODIFIED   nginx-with-git         2/2     Running             0             7s    10.1.1.180   docker-desktop   <none>           <none>
MODIFIED   nginx-with-git         1/2     NotReady            0             13s   10.1.1.180   docker-desktop   <none>           <none>
ADDED      nginx-with-init        0/1     Pending             0             0s    <none>       <none>           <none>           <none>
MODIFIED   nginx-with-init        0/1     Pending             0             0s    <none>       docker-desktop   <none>           <none>
MODIFIED   nginx-with-init        0/1     Init:0/1            0             0s    <none>       docker-desktop   <none>           <none>
MODIFIED   nginx-with-init        0/1     Init:0/1            0             6s    10.1.1.181   docker-desktop   <none>           <none>
MODIFIED   nginx-with-init        0/1     PodInitializing     0             20s   10.1.1.181   docker-desktop   <none>           <none>
MODIFIED   nginx-with-init        1/1     Running             0             26s   10.1.1.181   docker-desktop   <none>           <none>
```
- `initContainers` run before other containers
- normal containers are meant to be running continuously forever, and if they crash the pod restarts them and keeps them up running
- the **init-containers** are meant to run just once, 
	- and it's supposed to be like a finite process 
#### Use cases for init-containers
- like downloading some content
- or generating some certificate
- or generating some configuration file
- or some one time or one shot initialization task and then that container exits and we don't see it again
- something like database migrations 
	- not always best solutions if we end up scaling pods later, where we might end up running database migrations multiple times, unless we properly acquire lock on the database before migration starts
- waiting for other services to be up
	- e.g. - we have a web server that needs to connect with the database
		- the main container would run only if the init-container finds the database to be ready and gets killed

## Volume lifecycle
- Lifecycle of volume is linked to the pod's lifecycle
- as long as the pod exists, the volume exists
- thus, we should not put persistent data in such volumes, since on loosing a pod, we loose the data
- volumes survive container restarts