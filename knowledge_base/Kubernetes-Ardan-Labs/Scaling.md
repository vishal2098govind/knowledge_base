#k8s

## Scaling DockerCoin App
### Adding more workers for more loops per second
```sh
$ kubectl scale deployment worker --replicas 2
```
- if we want to scale `rng` as well such that we get one `rng` per machine or per node
	- if we have 4 nodes in our cluster, we need 4 `rng` replicas, one per node
	- This, one replica (pod) per node, is possible to achieve by using daemon sets

## Daemon Sets
#daemonsets
- a way in k8s to get one pod per node
- we could have just scale up `rng` deployment to `--replicas=4`, but it would not necessarily be distributed as one pod per node
	- imagine, if the cluster is very full, 
	- there are 4 nodes in the cluster, and they are completely full
	- then, if we realize the need to add a 5th node, and if we provision node-5 and if that node-5 is almost empty
	- now, if we scale up `rng`
	- where will the `rng` pods go
		- all the new `rng` pods are going to land on node-5, since all other nodes are full
- if we use daemon sets, we can guarantee that we will get one pod per node
- if we scale a daemon-set, in fact we do not scale a daemon-set
- we don't say we want a specific number of pods, daemon-sets give exactly one pod per node
- default daemon-set pods in every k8s
```sh
$ kubectl get pods --namespace kube-system -o wide
NAME                                     READY   STATUS    RESTARTS      AGE    IP             NODE             NOMINATED NODE   READINESS GATES
coredns-668d6bf9bc-6jqg7                 1/1     Running   0             110m   10.1.0.124     docker-desktop   <none>           <none>
coredns-668d6bf9bc-zwdqk                 1/1     Running   0             110m   10.1.0.123     docker-desktop   <none>           <none>
etcd-docker-desktop                      1/1     Running   4 (36h ago)   2d7h   192.168.65.3   docker-desktop   <none>           <none>
kube-apiserver-docker-desktop            1/1     Running   4 (36h ago)   2d7h   192.168.65.3   docker-desktop   <none>           <none>
kube-controller-manager-docker-desktop   1/1     Running   5 (36h ago)   2d7h   192.168.65.3   docker-desktop   <none>           <none>
kube-proxy-qgwbh                         1/1     Running   1 (36h ago)   2d7h   192.168.65.3   docker-desktop   <none>           <none>
kube-scheduler-docker-desktop            1/1     Running   7 (33h ago)   2d7h   192.168.65.3   docker-desktop   <none>           <none>
$ kubectl get daemonsets --namespace kube-system -o wide
NAME         DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR            AGE    CONTAINERS   IMAGES                               SELECTOR
kube-proxy   1         1         1       1            1           kubernetes.io/os=linux   2d7h   kube-proxy   registry.k8s.io/kube-proxy:v1.32.2   k8s-app=kube-proxy
```

## Creating `daemonset`
```sh
$ kubectl create daemonset rng --image dockercoins/rng:v0.1
error: unknown flag: --image
See 'kubectl create --help' for usage.
```
- we cannot create `daemonset` with `kubectl` or CLI
- we can create with YAML
- to write the YAML, we can try and adapt the YAML manifest of a deployment
- this is because, `daemonsets` and `deployments` are pretty close
- in both cases, we have 
	- a pod template
	- and then, in the deployment, we've something indicating how many replicas we want
	- and in the `daemonset`, we don't have the number of replicas, but otherwise it should be pretty similar in the sense that both needs multiple pods
```sh
$ kubectl create deployment rn --image dockercoins/rng:v0.1 --dry-run=client -o yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: rn
  name: rn
spec:
  replicas: 1
  selector:
    matchLabels:
      app: rn
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: rn
    spec:
      containers:
      - image: dockercoins/rng:v0.1
        name: rng
        resources: {}
status: {}
kubectl create deployment rng --image dockercoins/rng:v0.1 --dry-run=client -o yaml > rng.yaml
$ code rng.yaml
# replace deployment with daemonsets in Kind field
$ cat rng.yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  creationTimestamp: null
  labels:
    app: rng
  name: rng
spec:
  replicas: 1
  selector:
    matchLabels:
      app: rn
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: rn
    spec:
      containers:
      - image: dockercoins/rng:v0.1
        name: rng
        resources: {}
status: {}
$ kubectl create -f rng.yaml
Error from server (BadRequest): error when creating "rng.yaml": DaemonSet in version "v1" cannot be handled as a DaemonSet: strict decoding error: unknown field "spec.replicas", unknown field "spec.strategy"
$ code rng.yaml
# try removing spec.replicas field and spec.stratergy
$  kubectl create -f rng.yaml
daemonset.apps/rng created
$ kubectl get daemonsets
NAME   DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR   AGE
rng    1         1         1       1            1           <none>          4s
```
- we need not do `kubectl expose daemonset` like we had to do for `kubectl expose deployment` explicitly
	- If a service (load balancer or cluster-ip) already exists for a deployment and we are creating `daemonset` for the same pods used that deployment, then, we can use at least one matching **label** of the same deployment for the endpoints in the load-balancer of that deployment to get updated with the new pods added by the `daemonset`
- `daemonset` while creating multiple copies of pods, one per node, it doesn't consider the node with control-plane as it is reserved for control-plane
- the `daemonset` doesn't create a pod in the control-plane node
	- not because the role of that node had `control-plane` on it
	- but because this has to do with a mechanism called **taints-and-tolerations**
	- because the node with `control-plane` had **taints** saying `NoSchedule` 
		- and the pod in the `daemonset` which we wanted to create **didn't have** that same toleration saying `NoSchedule`
	- 
## Taints & Tolerations
#taints #tolerations
- **Taints** and **tolerations** are two features that kind of go together
- **Taints** are the things that we put on **nodes**
- **Tolerations** are the things that we put on **pods**
- **Taints** on the nodes will be used to prevent pods from running on that node unless the pods carry those **tolerations** to **cancel out the taints** in the node
- **Taint** can be imagined as some kind of a **police tape** **at a crime scene** saying do-not-cross or maybe a **VIP area**
	- and only the pods **with the VIP pass** or **who are police** can enter that VIP area or crime scene
	- and thus, here the **VIP pass** is the **Toleration**
- the metaphor being used here is a **Taint** is something that stinks and by default, the pods don't want to go into those nodes unless they can **Tolerate**
```sh
$ kubectl get nodes docker-desktop
NAME             STATUS   ROLES           AGE    VERSION
docker-desktop   Ready    control-plane   2d9h   v1.32.2
$ kubectl describe nodes docker-desktop | grep Taints
Taints:             node-role.kubernetes.io/control-plane:NoSchedule
$ kubectl get pods --namespace kube-system
NAME                                     READY   STATUS    RESTARTS      AGE
coredns-668d6bf9bc-6jqg7                 1/1     Running   0             3h45m
coredns-668d6bf9bc-zwdqk                 1/1     Running   0             3h45m
etcd-docker-desktop                      1/1     Running   4 (38h ago)   2d9h
kube-apiserver-docker-desktop            1/1     Running   4 (38h ago)   2d9h
kube-controller-manager-docker-desktop   1/1     Running   5 (38h ago)   2d9h
kube-proxy-qgwbh                         1/1     Running   1 (38h ago)   2d9h
kube-scheduler-docker-desktop            1/1     Running   7 (35h ago)   2d9h
$ kubectl describe pods coredns-668d6bf9bc-6jqg7 --namespace kube-system
Name:                 coredns-668d6bf9bc-6jqg7
Namespace:            kube-system
Priority:             2000000000
Priority Class Name:  system-cluster-critical
Service Account:      coredns
Node:                 docker-desktop/192.168.65.3
Start Time:           Thu, 01 Jan 2026 22:40:37 +0530
Labels:               k8s-app=kube-dns
                      pod-template-hash=668d6bf9bc
Annotations:          <none>
Status:               Running
IP:                   10.1.0.124
IPs:
  IP:           10.1.0.124
Controlled By:  ReplicaSet/coredns-668d6bf9bc
Containers:
  coredns:
    Container ID:  docker://3f2b497f379e72985d1b6b96a701f128af29d7678588aa3fb6f0e6cca0b87723
    Image:         registry.k8s.io/coredns/coredns:v1.11.3
    Image ID:      docker-pullable://registry.k8s.io/coredns/coredns@sha256:9caabbf6238b189a65d0d6e6ac138de60d6a1c419e5a341fbbb7c78382559c6e
    Ports:         53/UDP, 53/TCP, 9153/TCP
    Host Ports:    0/UDP, 0/TCP, 0/TCP
    Args:
      -conf
      /etc/coredns/Corefile
    State:          Running
      Started:      Thu, 01 Jan 2026 22:40:49 +0530
    Ready:          True
    Restart Count:  0
    Limits:
      memory:  170Mi
    Requests:
      cpu:        100m
      memory:     70Mi
    Liveness:     http-get http://:8080/health delay=60s timeout=5s period=10s #success=1 #failure=5
    Readiness:    http-get http://:8181/ready delay=0s timeout=1s period=10s #success=1 #failure=3
    Environment:  <none>
    Mounts:
      /etc/coredns from config-volume (ro)
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-twwqs (ro)
Conditions:
  Type                        Status
  PodReadyToStartContainers   True
  Initialized                 True
  Ready                       True
  ContainersReady             True
  PodScheduled                True
Volumes:
  config-volume:
    Type:      ConfigMap (a volume populated by a ConfigMap)
    Name:      coredns
    Optional:  false
  kube-api-access-twwqs:
    Type:                    Projected (a volume that contains injected data from multiple sources)
    TokenExpirationSeconds:  3607
    ConfigMapName:           kube-root-ca.crt
    ConfigMapOptional:       <nil>
    DownwardAPI:             true
QoS Class:                   Burstable
Node-Selectors:              kubernetes.io/os=linux
Tolerations:                 CriticalAddonsOnly op=Exists
                             node-role.kubernetes.io/control-plane:NoSchedule # this is the VIP pass this pod has to be able to be scheduled on a node with `NoSchedule` taint
                             node.kubernetes.io/not-ready:NoExecute op=Exists for 300s
                             node.kubernetes.io/unreachable:NoExecute op=Exists for 300s
Events:                      <none>
```
- A **Taint** is not like a binary thing
	- It's not like if there's a taint or if there's no taint
	- The taint has a **name** and **effect**
		- Here, the name of the taint is `node-role.kubernetes.io/control-plane`
		- and the effect is `NoSchedule`
	- There's also a taint effect called `NoExecute`
### Taint effects
- `NoSchedule` taint effect if found on a node
	- New pods cannot go there
	- If a pod, by any how, is found in the node, that's ok for this effect
- `NoExecute` taint effect if found on a node
	- no pod cannot go there
	- and if, by any how, a pod is found, it's going to get kicked out of the node
- `PreferNoSchedule` taint effect if found on a node
	- it means, we would prefer if pods would not be scheduled on that node, but if there's really no other choice, that's ok for this effect
	- it is used for the cluster auto-scaling when scaling down
	- imagine
		- if we have a big spike of traffic
		- we've lots of pods
		- automatically k8s is adding nodes to the cluster
		- then the traffic spike is gone, and we scale down the cluster
		- now, we have bunch of nodes that are empty
		- k8s is then going to put a `PreferNoSchedule` taint effect on those empty nodes, to indicate the scheduler or anyone who bothers, that it's considering to shut down this node, so try to avoid starting new pods on this node, if possible put that pod somewhere else so that we can shut down the node and save some money
### A common use case for taints-and-tolerations 
- when using GPU nodes
- if we do something that leverages GPU, we typically want to reserve those machines or nodes with GPUs for the GPU workload
- because, GPUs are pretty expensive, either we own or rent from Cloud providers, 
- if we have a bunch of containers running and out of them some of them need GPU, 
- we want to make sure that the containers which need GPU should be having GPU available, 
- and the containers do not need GPUs, they should not be on machines with GPUs because we would otherwise be wasting money
- one way to ensure that is put a **taint** like `GPUNoSchedule` on the nodes that have GPUs and having a **toleration** `GPUNoSchedule` on the pods running containers that need GPU
- this will prevent normal pods who don't have this toleration of `GPUNoSchedule` running on such nodes
## How we need not had to expose the `daemonset` separately, like we had to do for `deployment`
- **Labels** and **selectors** help to do that if possible
- here, we already had a service for `rng`
```sh
$ kubectl describe service rng
Name:                     rng
Namespace:                dev
Labels:                   app=rng
                          release=dev
Annotations:              <none>
Selector:                 app=rng # this indicates to which pods should this service be balancing the load among
Type:                     ClusterIP
IP Family Policy:         SingleStack
IP Families:              IPv4
IP:                       10.97.48.98
IPs:                      10.97.48.98
Port:                     <unset>  80/TCP
TargetPort:               80/TCP
Endpoints:                10.1.0.132:80,10.1.0.141:80 # this has the list of IP addresses of all the pods this service/cluster-ip can load balance the traffic to
Session Affinity:         None
Internal Traffic Policy:  Cluster
Events:                   <none>
```
- in a service, there exists a selector which indicates to which pods should the service be balancing the load among
```sh
$ kubectl get pods --selector app=rng -o wide --show-labels
NAME                   READY   STATUS    RESTARTS   AGE     IP           NODE             NOMINATED NODE   READINESS GATES   LABELS
rng-65d885d498-sgnhm   1/1     Running   0          3h55m   10.1.0.132   docker-desktop   <none>           <none>            app=rng,pod-template-hash=65d885d498
rng-zrklh              1/1     Running   0          100m    10.1.0.141   docker-desktop   <none>           <none>            app=rng,controller-revision-hash=65d885d498,pod-template-generation=1
```
- thus, the IP addresses of the pods match the Endpoints list of the `rng` service
- this means that, a service sends connections to a bunch of pods, and these pods are defined with a selector
- this is pretty powerful because it means that, when we have a service with multiple pods, we don't need to explicitly list these pods or something
- conceptually the service **continuously evaluates** that selector and so each time a pod with that selector is added or removed, the pod will also be added or removed from the service or load balancer or the endpoints list
- and, since while creating `rng` `daemonset` we used the same selector of `app=rng` as we used while creating the `rng` deployment, 
- the `rng` `daemonset` created pods with the same selector of `app=rng`
- and thus, the `rng` service was able to detect the pods created by the `daemonsets` since the `rng` service has the selector field of `app=rng`
```sh
$ kubectl create deployment rng --image dockercoins/rng:v0.1 --dry-run=client -o yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: rng
  name: rng
spec:
  replicas: 1
  selector:
    matchLabels:
      app: rng
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: rng # the `kubectl create` by default uses the name of the deployment as the value of the app label. we used this yaml as our base YAML while creating the YAML manifest for the daemonset, where we kept the app=rng label same
    spec:
      containers:
      - image: dockercoins/rng:v0.1
        name: rng
        resources: {}
status: {}


$ kubectl expose deployment rng --name lb_rng --port 80 --dry-run=client -o yaml
apiVersion: v1
kind: Service
metadata:
  creationTimestamp: null
  labels:
    app: rng
    release: dev
  name: lb_rng
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: rng # this selector ensures that it continuously looks for pods with label of app=rng
status:
  loadBalancer: {}

$ cat rng.yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  creationTimestamp: null
  labels:
    app: rng
  name: rng
spec:
  selector:
    matchLabels:
      app: rng
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: rng # this label ensures the pods created by this daemonset have this label, so that the rng service can discover and balance load to those pods
    spec:
      containers:
      - image: dockercoins/rng:v0.1
        name: rng
        resources: {}
status: {}
```