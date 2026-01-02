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
#daemon-sets
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

## Exercise - Treat a sick pod, while keeping other fit pods as it is
### Idea - 1 : Change the label of sick pod, keeping others as it is
- NOTE: this idea is when the label being changed is same label that is being used by the replica-set spec selector while matching labels
```sh
$ kubectl edit pod rng-65d885d498-sgnhm
pod/rng-65d885d498-sgnhm edited
$ kubectl get pods rng-65d885d498-sgnhm -o yaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: "2026-01-01T17:10:02Z"
  generateName: rng-65d885d498-
  labels:
    app: sick_rng
    pod-template-hash: 65d885d498
  name: rng-65d885d498-sgnhm
  namespace: dev
  resourceVersion: "121129"
  uid: 86139c4b-411c-4ede-b3bf-0a72d054f4a9
spec:
  containers:
  - image: dockercoins/rng:v0.1
    imagePullPolicy: IfNotPresent
    name: rng
    resources: {}
    terminationMessagePath: /dev/termination-log
    terminationMessagePolicy: File
    volumeMounts:
    - mountPath: /var/run/secrets/kubernetes.io/serviceaccount
      name: kube-api-access-qgp4j
      readOnly: true
  dnsPolicy: ClusterFirst
  enableServiceLinks: true
  nodeName: docker-desktop
  preemptionPolicy: PreemptLowerPriority
  priority: 0
  restartPolicy: Always
  schedulerName: default-scheduler
  securityContext: {}
  serviceAccount: default
  serviceAccountName: default
  terminationGracePeriodSeconds: 30
  tolerations:
  - effect: NoExecute
    key: node.kubernetes.io/not-ready
    operator: Exists
    tolerationSeconds: 300
  - effect: NoExecute
    key: node.kubernetes.io/unreachable
    operator: Exists
    tolerationSeconds: 300
  volumes:
  - name: kube-api-access-qgp4j
    projected:
      defaultMode: 420
      sources:
      - serviceAccountToken:
          expirationSeconds: 3607
          path: token
      - configMap:
          items:
          - key: ca.crt
            path: ca.crt
          name: kube-root-ca.crt
      - downwardAPI:
          items:
          - fieldRef:
              apiVersion: v1
              fieldPath: metadata.namespace
            path: namespace
status:
  conditions:
  - lastProbeTime: null
    lastTransitionTime: "2026-01-01T17:26:14Z"
    status: "True"
    type: PodReadyToStartContainers
  - lastProbeTime: null
    lastTransitionTime: "2026-01-01T17:10:38Z"
    status: "True"
    type: Initialized
  - lastProbeTime: null
    lastTransitionTime: "2026-01-01T17:26:14Z"
    status: "True"
    type: Ready
  - lastProbeTime: null
    lastTransitionTime: "2026-01-01T17:26:14Z"
    status: "True"
    type: ContainersReady
  - lastProbeTime: null
    lastTransitionTime: "2026-01-01T17:10:31Z"
    status: "True"
    type: PodScheduled
  containerStatuses:
  - containerID: docker://3effb0a6babc104db962ccdb3b5e4ee4c26401bdb774790209906f23af71ca0c
    image: dockercoins/rng:v0.1
    imageID: docker-pullable://dockercoins/rng@sha256:17f79daa5cbb38319519be0cd6c2d81b17327c118b944983d6f7eb26826607a9
    lastState: {}
    name: rng
    ready: true
    restartCount: 0
    started: true
    state:
      running:
        startedAt: "2026-01-01T17:10:51Z"
    volumeMounts:
    - mountPath: /var/run/secrets/kubernetes.io/serviceaccount
      name: kube-api-access-qgp4j
      readOnly: true
      recursiveReadOnly: Disabled
  hostIP: 192.168.65.3
  hostIPs:
  - ip: 192.168.65.3
  phase: Running
  podIP: 10.1.0.132
  podIPs:
  - ip: 10.1.0.132
  qosClass: BestEffort
  startTime: "2026-01-01T17:10:38Z"
```
- but if we look at the pods and services:
```sh
$  kubectl describe service rng
Name:                     rng
Namespace:                dev
Labels:                   app=rng
                          release=dev
Annotations:              <none>
Selector:                 app=rng
Type:                     ClusterIP
IP Family Policy:         SingleStack
IP Families:              IPv4
IP:                       10.97.48.98
IPs:                      10.97.48.98
Port:                     <unset>  80/TCP
TargetPort:               80/TCP
Endpoints:                10.1.0.141:80,10.1.0.142:80
Session Affinity:         None
Internal Traffic Policy:  Cluster
Events:                   <none>
$ kubectl get pods --selector app=rng -o wide --show-labels
NAME                   READY   STATUS    RESTARTS   AGE     IP           NODE             NOMINATED NODE   READINESS GATES   LABELS
rng-65d885d498-8vf27   1/1     Running   0          2m29s   10.1.0.142   docker-desktop   <none>           <none>            app=rng,pod-template-hash=65d885d498
rng-zrklh              1/1     Running   0          14h     10.1.0.141   docker-desktop   <none>           <none>            app=rng,controller-revision-hash=65d885d498,pod-template-generation=1
```
- changing the label lead to create another replacement pod with that label right away, instead of only removing the pod in that label
- this happens because of replica-set:
	- the replica-set that was created the deployment, has `--replicas=1` it's YAML manifest
	- a replica-set works with labels and selectors
	- in a replica-set we have a selector, and what a replica-set wants to do is make sure that we have N pods with given labels
```sh
$ kubectl get rs --selector app=rng -o yaml
apiVersion: v1
items:
- apiVersion: apps/v1
  kind: ReplicaSet
  metadata:
    annotations:
      deployment.kubernetes.io/desired-replicas: "1"
      deployment.kubernetes.io/max-replicas: "2"
      deployment.kubernetes.io/revision: "1"
    creationTimestamp: "2025-12-30T11:24:07Z"
    generation: 1
    labels:
      app: rng
      pod-template-hash: 65d885d498
      release: dev
    name: rng-65d885d498
    namespace: dev
    ownerReferences:
    - apiVersion: apps/v1
      blockOwnerDeletion: true
      controller: true
      kind: Deployment
      name: rng
      uid: e03f037b-d32a-46a5-84ea-90b9e5e6b613
    resourceVersion: "121146"
    uid: e150933d-9732-447b-b1bd-454a2726ef2e
  spec:
    replicas: 1
    selector:
      matchLabels:
        app: rng
        pod-template-hash: 65d885d498
    template:
      metadata:
        creationTimestamp: null
        labels:
          app: rng
          pod-template-hash: 65d885d498
      spec:
        containers:
        - image: dockercoins/rng:v0.1
          imagePullPolicy: IfNotPresent
          name: rng
          resources: {}
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
        dnsPolicy: ClusterFirst
        restartPolicy: Always
        schedulerName: default-scheduler
        securityContext: {}
        terminationGracePeriodSeconds: 30
  status:
    availableReplicas: 1
    fullyLabeledReplicas: 1
    observedGeneration: 1
    readyReplicas: 1
    replicas: 1
kind: List
metadata:
  resourceVersion: ""
```
- if we see the specification of replica set, what it really says is
```sh
spec:
    replicas: 1
    selector:
      matchLabels:
        app: rng
        pod-template-hash: 65d885d498
```
- it makes sure that we have 1 pod with the label matching the selector of `app=rng`
- the replica set just checks for pods with matching labels, and not really checks whats actually running inside the pod, it could be `color` image or `nginx` image or any image container running in the pod, it should have the specific label
- thus, the moment we change the label of the pod, we get other replacement pod right away
```sh
$ kubectl describe service rng
Name:                     rng
Namespace:                dev
Labels:                   app=rng
                          release=dev
Annotations:              <none>
Selector:                 app=rng
Type:                     ClusterIP
IP Family Policy:         SingleStack
IP Families:              IPv4
IP:                       10.97.48.98
IPs:                      10.97.48.98
Port:                     <unset>  80/TCP
TargetPort:               80/TCP
Endpoints:                10.1.0.141:80,10.1.0.142:80
Session Affinity:         None
Internal Traffic Policy:  Cluster
Events:                   <none>
$ kubectl get pods --show-labels -o wide | grep rng
NAME                      READY   STATUS    RESTARTS   AGE   IP           NODE             NOMINATED NODE   READINESS GATES   LABELS
rng-65d885d498-8vf27      1/1     Running   0          51m   10.1.0.142   docker-desktop   <none>           <none>            app=rng,pod-template-hash=65d885d498
rng-65d885d498-sgnhm      1/1     Running   0          17h   10.1.0.132   docker-desktop   <none>           <none>            app=sick_rng,pod-template-hash=65d885d498
rng-zrklh                 1/1     Running   0          15h   10.1.0.141   docker-desktop   <none>           <none>            app=rng,controller-revision-hash=65d885d498,pod-template-generation=1
```
- now, we can debug and trouble shoot the sick pod, while the replacement pod could serve the traffic
- now, if we are done without trouble shooting and want to put the label back to `app=rng` from `app=sick_rng`
```sh
$ kubectl label pod rng-65d885d498-sgnhm app=rng
error: 'app' already has a value (sick_rng), and --overwrite is false
$ kubectl label pod rng-65d885d498-sgnhm app=rng --overwrite
pod/rng-65d885d498-sgnhm labelled
# ...
# a few seconds later, that pod (rng-65d885d498-sgnhm) dissapears
$ kubectl get pods -o wide --selector app=rng --show-labels
NAME                   READY   STATUS    RESTARTS   AGE     IP           NODE             NOMINATED NODE   READINESS GATES   LABELS
rng-65d885d498-8vf27   1/1     Running   0          72m     10.1.0.142   docker-desktop   <none>           <none>            app=rng,pod-template-hash=65d885d498
rng-n9js7              1/1     Running   0          5m12s   10.1.0.143   docker-desktop   <none>           <none>            app=rng,controller-revision-hash=65d885d498,pod-template-generation=1
$  kubectl get pods rng-65d885d498-sgnhm
Error from server (NotFound): pods "rng-65d885d498-sgnhm" not found
```
- when we updated the label of pod `rng-65d885d498-sgnhm` back to `app=rng` from `app=sick_rng`, the replica-set found two pods matching the selector labels which the replica-set spec contain, and thus, replica-set decides to delete or remove one of the two as the replica-set spec only specifies 1 replica.

### Idea - 2 : add a common label to all the pods and then use that in the service/load-balancer selector
- Here, the idea is not to change the label that is used by the replica-set selector, to separate the sick pod and treat, and rather 
	- change the label which the load-balancer uses to find pods, if load balancer listens finds pods using a different selector than that of replica-set
	- or if load balancer the load balancer uses same selector as used by the replica-set, then make load balancer use a different or new selector to find pods and use that label in all pods, 
	- so here **replica-set doesn't know about the treatment that is being taken place of the sick pod**
```sh
$ kubectl label pods active="yes" --selector app=rng
pod/rng-65d885d498-8vf27 labeled
pod/rng-n9js7 labeled

$ kubectl get pods --selector app=rng --show-labels
NAME                   READY   STATUS    RESTARTS   AGE   LABELS
rng-65d885d498-8vf27   1/1     Running   0          88m   active=yes,app=rng,pod-template-hash=65d885d498
rng-n9js7              1/1     Running   0          21m   active=yes,app=rng,controller-revision-hash=65d885d498,pod-template-generation=1

$ kubectl edit service rng
# update the selector to be `active="yes"`

# label the sick pod active=no
$ kubectl label pods rng-65d885d498-8vf27 active=no --overwrite
pods/rng-65d885d498-8vf27 labelled

# treat the sick pod
# label the sick pod active=yes
$ kubectl label pods rng-65d885d498-8vf27 active=yes --overwrite
pods/rng-65d885d498-8vf27 labelled
```
- if we want to remove or drop a label, we can use `<label-name>-` i.e.
```sh
# after treating the sick pod

$ kubectl edit service rng
# update the selector to be `app=rng` instead of `active=yes`

# remove active label
$ kubectl label service rng active-
$ kubectl label pods active- --selector app=rng
```
- All steps:
```sh
$ kubectl label pods active="yes" --selector app=rng
pod/rng-65d885d498-zl8qq labeled
pod/rng-lh6ts labeled # this is running due to daemonset
$ kubectl label service rng active="yes"
service/rng labeled


$ kubectl get pods --selector app=rng --show-labels
NAME                   READY   STATUS    RESTARTS   AGE    LABELS
rng-65d885d498-zl8qq   1/1     Running   0          6h5m   active=yes,app=rng,pod-template-hash=65d885d498
rng-lh6ts              1/1     Running   0          33m    active=yes,app=rng,controller-revision-hash=65d885d498,pod-template-generation=1
$ kubectl get service rng --show-labels
NAME   TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)   AGE     LABELS
rng    ClusterIP   10.97.48.98   <none>        80/TCP    3d10h   active=yes,app=rng,release=dev
$ kubectl get replicasets --selector app=rng --show-labels
NAME             DESIRED   CURRENT   READY   AGE     LABELS
rng-65d885d498   1         1         1       3d10h   app=rng,pod-template-hash=65d885d498,release=dev

# set label of the sick pod `active="no"`
$ kubectl label pod rng-65d885d498-zl8qq active="no" --overwrite
pod/rng-65d885d498-zl8qq labeled
$ kubectl edit service rng
# update selector of rng service/load-balancer to match labels `active="yes"`
$ kubectl get service rng -o yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"v1","kind":"Service","metadata":{"annotations":{},"labels":{"app":"rng"},"name":"rng","namespace":"dev"},"spec":{"ports":[{"port":80,"protocol":"TCP","targetPort":80}],"selector":{"app":"rng"},"type":"ClusterIP"}}
  creationTimestamp: "2025-12-30T11:24:07Z"
  labels:
    active: "yes"
    app: rng
    release: dev
  name: rng
  namespace: dev
  resourceVersion: "145610"
  uid: 928a5b35-ec8b-4994-99ac-b5266364ce95
spec:
  clusterIP: 10.97.48.98
  clusterIPs:
  - 10.97.48.98
  internalTrafficPolicy: Cluster
  ipFamilies:
  - IPv4
  ipFamilyPolicy: SingleStack
  ports:
  - port: 80
    protocol: TCP
    targetPort: 80
  selector:
    active: "yes"
  sessionAffinity: None
  type: ClusterIP
status:
  loadBalancer: {}
  
# treat the sick pod

$ kubectl describe service rng
Name:                     rng
Namespace:                dev
Labels:                   active=yes
                          app=rng
                          release=dev
Annotations:              <none>
Selector:                 active=yes # selects pods which match this label, to balance load
Type:                     ClusterIP
IP Family Policy:         SingleStack
IP Families:              IPv4
IP:                       10.97.48.98
IPs:                      10.97.48.98
Port:                     <unset>  80/TCP
TargetPort:               80/TCP
Endpoints:                10.1.0.146:80
Session Affinity:         None
Internal Traffic Policy:  Cluster
Events:                   <none>
$ kubectl get pods --selector app=rng -o wide --show-labels
NAME                   READY   STATUS    RESTARTS   AGE     IP           NODE             NOMINATED NODE   READINESS GATES   LABELS
rng-65d885d498-zl8qq   1/1     Running   0          6h22m   10.1.0.145   docker-desktop   <none>           <none>            active=no,app=rng,pod-template-hash=65d885d498
rng-lh6ts              1/1     Running   0          50m     10.1.0.146   docker-desktop   <none>           <none>            active=yes,app=rng,controller-revision-hash=65d885d498,pod-template-generation=1

# change back the selector of load-balancer to match labels with `app=rng`
$ kubectl edit service rng
service/rng edited
$ kubectl describe service rng
Name:                     rng
Namespace:                dev
Labels:                   active=yes
                          app=rng
                          release=dev
Annotations:              <none>
Selector:                 app=rng
Type:                     ClusterIP
IP Family Policy:         SingleStack
IP Families:              IPv4
IP:                       10.97.48.98
IPs:                      10.97.48.98
Port:                     <unset>  80/TCP
TargetPort:               80/TCP
Endpoints:                10.1.0.146:80,10.1.0.145:80
Session Affinity:         None
Internal Traffic Policy:  Cluster
Events:                   <none>
$ kubectl get pods --selector app=rng -o wide --show-labels
NAME                   READY   STATUS    RESTARTS   AGE     IP           NODE             NOMINATED NODE   READINESS GATES   LABELS
rng-65d885d498-zl8qq   1/1     Running   0          6h26m   10.1.0.145   docker-desktop   <none>           <none>            active=no,app=rng,pod-template-hash=65d885d498
rng-lh6ts              1/1     Running   0          54m     10.1.0.146   docker-desktop   <none>           <none>            active=yes,app=rng,controller-revision-hash=65d885d498,pod-template-generation=1

# remove active label from the pods and load-balancer
$ kubectl label service rng active- --overwrite
service/rng labeled
$ kubectl describe service rng
Name:                     rng
Namespace:                dev
Labels:                   app=rng
                          release=dev
Annotations:              <none>
Selector:                 app=rng
Type:                     ClusterIP
IP Family Policy:         SingleStack
IP Families:              IPv4
IP:                       10.97.48.98
IPs:                      10.97.48.98
Port:                     <unset>  80/TCP
TargetPort:               80/TCP
Endpoints:                10.1.0.146:80,10.1.0.145:80
Session Affinity:         None
Internal Traffic Policy:  Cluster
Events:                   <none>
$ kubectl label pods --selector app=rng active-
pod/rng-65d885d498-zl8qq unlabeled
pod/rng-lh6ts unlabeled
$ kubectl get pods --selector app=rng -o wide --show-labels
NAME                   READY   STATUS    RESTARTS   AGE     IP           NODE             NOMINATED NODE   READINESS GATES   LABELS
rng-65d885d498-zl8qq   1/1     Running   0          6h29m   10.1.0.145   docker-desktop   <none>           <none>            app=rng,pod-template-hash=65d885d498
rng-lh6ts              1/1     Running   0          57m     10.1.0.146   docker-desktop   <none>           <none>            app=rng,controller-revision-hash=65d885d498,pod-template-generation=1
```

### Conclusion
- this seems pretty nice because this means now that we can manipulate the load balancer configuration just with labels and selectors, without bother about
	- how to configure IP tables
	- how to write configuration for HAProxy, or NGINX or Apache or whatever being used as the load balancer
- and thus, 
	- on changing the labels and selectors of pods and load-balancers respectively, 
	- immediately in realtime, 
	- the load balancers adds or removes backends from it's endpoints list