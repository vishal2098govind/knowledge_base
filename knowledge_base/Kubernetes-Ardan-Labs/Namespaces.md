#k8s
- Namespaces are a way to organize our resources on k8s
- They are a logical separation, not a physical separation
- ~ directories/folders on our computer
- The four basic namespaces that we have on pretty much every modern k8s cluster
```sh
$ kubectl get namespaces
NAME              STATUS   AGE
default           Active   33h
kube-node-lease   Active   33h
kube-public       Active   33h
kube-system       Active   33h
```
- we can create more namespaces

## Creating namespaces
### `kubectl create namespace`
```sh
$ kubectl create namespace dockercoin
namespace/dockercoin created
```
### Using YAML Manifest
```sh
$ kubectl create namespace blue
namespace/blue created
$ kubectl get namespace blue -o yaml
apiVersion: v1
kind: Namespace
metadata:
  creationTimestamp: "2025-12-29T18:42:23Z"
  labels:
    kubernetes.io/metadata.name: blue
  name: blue
  resourceVersion: "104378"
  uid: 79cf35f0-d9cd-4f02-a4bd-2fae29ab3c40
spec:
  finalizers:
  - kubernetes
status:
  phase: Active
```
- although while creating the namespaces using `kubectl create`, it created a yaml manifest with a bunch of stuff, we need not need all that to create a namespace
- many fields are present which are not strictly necessary
	- `metadata.creationTimestamp`: this is something that was automatically added by k8s
	- similarly for the `metadata.resourceVersion`, `uid`
	- even the `spec.finalizers[kubernetes]` and `status.phase=Active` stuff are added by k8s automatically later after we submitted just the `kubectl create namespace blue`
- thus, the bare minimum manifest for creating a namespace would be:
```sh
$ code green.yaml
```
```yaml
apiVersion: v1
kind: Namespace
metadata:
	name: green
```
```sh
# to create a namespace using the manifest green.yaml file:
$ kubectl create -f green.yaml
namespace/green created
$ kubectl get ns
NAME              STATUS   AGE
blue              Active   9m8s
default           Active   33h
green             Active   6s
kube-node-lease   Active   33h
kube-public       Active   33h
kube-system       Active   33h
$ kubectl get namespace green -o yaml
apiVersion: v1
kind: Namespace
metadata:
  creationTimestamp: "2025-12-29T18:51:25Z"
  labels:
    kubernetes.io/metadata.name: green
  name: green
  resourceVersion: "105103"
  uid: 9b37a9a6-b62e-4dde-9d3f-0d2d816a8551
spec:
  finalizers:
  - kubernetes
status:
  phase: Active
```

## Using Namespaces
```sh
$ kubectl create deployment purple --image jpetazzo/color --namespace blue
deployment.apps/purple created
$ kubectl get deployment
NAME       READY   UP-TO-DATE   AVAILABLE   AGE
hasher     1/1     1            1           6h51m
pingpong   4/4     4            4           6h9m
redis      1/1     1            1           6h51m
rng        1/1     1            1           6h51m
webui      1/1     1            1           6h50m
worker     1/1     1            1           6h51m
$ kubectl get deployment --namespace blue
NAME     READY   UP-TO-DATE   AVAILABLE   AGE
purple   1/1     1            1           31s
$ kubectl expose deployment purple --port 80
Error from server (NotFound): deployments.apps "purple" not found
$ kubectl expose deployment purple --port 80 --namespace blue
service/purple exposed
$ kubectl get services --namespace blue
NAME     TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
purple   ClusterIP   10.102.223.15   <none>        80/TCP    47s
kubectl edit service purple -n blue
service/purple edited # edited to change from cluster-ip to load-balancer to be able to access the server from my mac
$ kubectl get service -n blue
NAME     TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
purple   LoadBalancer   10.102.223.15   <pending>     80:30868/TCP   3m3s
$ kubectl get service -n blue
NAME     TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
purple   LoadBalancer   10.102.223.15   <pending>     80:30868/TCP   3m9s
$ kubectl get service -n blue
NAME     TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
purple   LoadBalancer   10.102.223.15   <pending>     80:30868/TCP   3m11s
$ kubectl get service
NAME         TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
hasher       ClusterIP      10.103.124.226   <none>        80/TCP         6h54m
kubernetes   ClusterIP      10.96.0.1        <none>        443/TCP        33h
redis        ClusterIP      10.109.22.65     <none>        6379/TCP       6h55m
rng          ClusterIP      10.106.8.225     <none>        80/TCP         6h55m
webui        LoadBalancer   10.96.255.64     localhost     80:32744/TCP   6h54m
➜  docker-ardan-labs kubectl delete service webui
service "webui" deleted
$ kubectl get service -n blue
NAME     TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
purple   LoadBalancer   10.102.223.15   <pending>     80:30868/TCP   3m30s
$ kubectl get service -n blue
NAME     TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
purple   LoadBalancer   10.102.223.15   localhost     80:30868/TCP   3m58s
$ curl localhost
🔵🟣This is pod blue/purple-65bb9bc655-mqg7f on linux/amd64, serving / for 192.168.65.3:60214.
```
- namespaces are useful to deploy other copies of our application

## Deploying `dockercoins` on `dev` namespace
### using `kubectl`
- we can just re-run the `kubectl create deployment <service>` commands for each of the 5 services with a `--namespace dev`
- similarly run `kubectl create service <service>` for each of the service with a `--namespace` for each command
- here, tedious or annoying task to do would to add `--namespace dev` every single time for each resource, without forgetting
- how can we change the default active namespace while using `kubectl`?
	- deploying `dockercoins` is not too hard
	- just 9 `kubectl` commands in total
		- 5 commands for `deployments`
		- 4 for `services`
	- but is there any way to deploy everything with a single command?
		- yes
```sh
$ kubectl create -f container.training/k8s/dockercoins.yaml --namespace dev
```
#### `dockercoins.yaml`
```sh
$ cat container.training/k8s/dockercoins.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: hasher
  name: hasher
spec:
  replicas: 1
  selector:
    matchLabels:
      app: hasher
  template:
    metadata:
      labels:
        app: hasher
    spec:
      containers:
      - image: dockercoins/hasher:v0.1
        name: hasher
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: hasher
  name: hasher
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: hasher
  type: ClusterIP
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: redis
  name: redis
spec:
  replicas: 1
  selector:
    matchLabels:
      app: redis
  template:
    metadata:
      labels:
        app: redis
    spec:
      containers:
      - image: redis
        name: redis
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: redis
  name: redis
spec:
  ports:
  - port: 6379
    protocol: TCP
    targetPort: 6379
  selector:
    app: redis
  type: ClusterIP
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: rng
  name: rng
spec:
  replicas: 1
  selector:
    matchLabels:
      app: rng
  template:
    metadata:
      labels:
        app: rng
    spec:
      containers:
      - image: dockercoins/rng:v0.1
        name: rng
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: rng
  name: rng
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: rng
  type: ClusterIP
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: webui
  name: webui
spec:
  replicas: 1
  selector:
    matchLabels:
      app: webui
  template:
    metadata:
      labels:
        app: webui
    spec:
      containers:
      - image: dockercoins/webui:v0.1
        name: webui
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: webui
  name: webui
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: webui
  type: NodePort
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: worker
  name: worker
spec:
  replicas: 1
  selector:
    matchLabels:
      app: worker
  template:
    metadata:
      labels:
        app: worker
    spec:
      containers:
      - image: dockercoins/worker:v0.1
        name: worker
```

### `kubectl apply` vs `kubectl create`
#### `kubectl create`
```sh
$ kubectl create -f container.training/k8s/dockercoins.yaml --namespace dev
deployment.apps/hasher created
service/hasher created
deployment.apps/redis created
service/redis created
deployment.apps/rng created
service/rng created
deployment.apps/webui created
service/webui created
deployment.apps/worker created

$ kubectl create -f container.training/k8s/dockercoins.yaml --namespace dev
Error from server (AlreadyExists): error when creating "container.training/k8s/dockercoins.yaml": deployments.apps "hasher" already exists
Error from server (AlreadyExists): error when creating "container.training/k8s/dockercoins.yaml": services "hasher" already exists
Error from server (AlreadyExists): error when creating "container.training/k8s/dockercoins.yaml": deployments.apps "redis" already exists
Error from server (AlreadyExists): error when creating "container.training/k8s/dockercoins.yaml": services "redis" already exists
Error from server (AlreadyExists): error when creating "container.training/k8s/dockercoins.yaml": deployments.apps "rng" already exists
Error from server (AlreadyExists): error when creating "container.training/k8s/dockercoins.yaml": services "rng" already exists
Error from server (AlreadyExists): error when creating "container.training/k8s/dockercoins.yaml": deployments.apps "webui" already exists
Error from server (AlreadyExists): error when creating "container.training/k8s/dockercoins.yaml": services "webui" already exists
Error from server (AlreadyExists): error when creating "container.training/k8s/dockercoins.yaml": deployments.apps "worker" already exists
```
#### `kubectl apply`
- we can also use `kubectl apply -f dockercoins.yaml --namespace blue`
```sh
$ kubectl apply -f container.training/k8s/dockercoins.yaml --namespace blue
deployment.apps/hasher created
service/hasher created
deployment.apps/redis created
service/redis created
deployment.apps/rng created
service/rng created
deployment.apps/webui created
service/webui created
deployment.apps/worker created

# if we apply again, instead of giving bunch of errors, we get "unchanged" response
$ kubectl apply -f container.training/k8s/dockercoins.yaml --namespace blue
deployment.apps/hasher unchanged
service/hasher unchanged
deployment.apps/redis unchanged
service/redis unchanged
deployment.apps/rng unchanged
service/rng unchanged
deployment.apps/webui unchanged
service/webui unchanged
deployment.apps/worker unchanged
```