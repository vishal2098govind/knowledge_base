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
		- There are tools like `kube-score` and `kube-linter` that help in checking our YAML
		- for example, most linters are going to scream at us if we don't include resource requests in our YAML manifests.
		- Resource requests are used to define
			- no.of CPUs
			- amount of RAM
			- that our containers need
			- this is very important to do that in production
	- use `kubectl create --dry-run`

## `kubectl create --dry-run`

### `--dry-run=client`
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

### `--dry-run=server`
- we can also do `--dry-run=server`
```sh
$ kubectl create deployment orange --image jpetazzo/color -o yaml --dry-run=server
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: "2025-12-30T21:13:41Z"
  generation: 1
  labels:
    app: orange
  name: orange
  namespace: dev
  uid: 4b098ce5-f287-467b-9692-5f769ff0e203
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: orange
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: orange
    spec:
      containers:
      - image: jpetazzo/color
        imagePullPolicy: Always
        name: color
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
status: {}
```
- this is much longer output
- what happens here is, we do a round trip to the API server
	- here, it generates the client side YAML
	- sends it to the server
	- the server sends us back the YAML but after adding all the default fields and values, that we don't need, but that are provided at runtime

### When to use what, `--dry-run=client` vs `--dry-run=server`
- usually we want to use `--dry-run=client`
- the purpose `--dry-run=server` is 
	- not to give use all those values that are provided at runtime, because honestly we don't care about those runtime added values 
	- but when we have things like admission control or web-hooks