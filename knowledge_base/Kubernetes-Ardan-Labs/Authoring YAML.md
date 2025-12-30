#k8s

## How to come up with the correct YAML manifests
- Even if we don't want to do GitOps, can't we just rely on `kubectl`?
	- no, not always
	- there are some features of k8s, that require YAML manifests
	- for example
		- resource limits
		- health checks
		- pods with multiple containers
		- DaemonSets, StatefulSets
		- and more!
- There are multiple ways to do that
	- we can search for examples in k8s docs
	- there are bunch of tools and websites to help with coming up with YAML files
	- use `kubectl create --dry-run`

### `kubectl create --dry-run`
```sh
$ kubectl create deployment purple --image jpetazzo/color -o yaml --dry-run
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: purple
  name: purple
spec:
  replicas: 1
  selector:
    matchLabels:
      app: purple
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: purple
    spec:
      containers:
      - image: jpetazzo/color
        name: color
        resources: {}
status: {}

$ kubectl create deployment purple --image jpetazzo/color -o yaml --dry-run=client > purple.yaml
$ cat purple.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: purple
  name: purple
spec:
  replicas: 1
  selector:
    matchLabels:
      app: purple
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: purple
    spec:
      containers:
      - image: jpetazzo/color
        name: color
        resources: {}
status: {}
```
- this gives a pretty good deployment YAML manifest
- we can still clean things here like
	- remove null fields like `creationTimestamp` or empty fields like `status`, `containers[].resources`
```sh
$ cat purple.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: purple
  name: purple
spec:
  replicas: 1
  selector:
    matchLabels:
      app: purple
  template:
    metadata:
      labels:
        app: purple
    spec:
      containers:
      - image: jpetazzo/color
        name: color
        
$ kubectl apply -f purple.yaml
deployment.apps/purple created
$ kubectl get all --selector app=purple
NAME                          READY   STATUS    RESTARTS   AGE
pod/purple-65bb9bc655-9vfd2   1/1     Running   0          24s

NAME                     READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/purple   1/1     1            1           24s

NAME                                DESIRED   CURRENT   READY   AGE
replicaset.apps/purple-65bb9bc655   1         1         1       24s
```