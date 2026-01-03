#k8s #daemon-sets #rolling-updates

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
Endpoints:                10.1.0.145:80,10.1.0.151:80
Session Affinity:         None
Internal Traffic Policy:  Cluster
Events:                   <none>

$ kubectl get pods --selector app=rng --show-labels -o wide
NAME                   READY   STATUS    RESTARTS   AGE   IP           NODE             NOMINATED NODE   READINESS GATES   LABELS
rng-65d885d498-zl8qq   1/1     Running   0          18h   10.1.0.145   docker-desktop   <none>           <none>            app=rng,pod-template-hash=65d885d498
rng-q7274              1/1     Running   0          14m   10.1.0.151   docker-desktop   <none>           <none>            app=rng,controller-revision-hash=65d885d498,pod-template-generation=1

# let's say the pod rng-q7274 is sick
# let's say we want to treat it
$ kubectl label pods active=yes --selector app=rng
pod/rng-65d885d498-zl8qq labeled
pod/rng-q7274 labeled

$ kubectl get pods -o wide --show-labels --selector app=rng
NAME                   READY   STATUS    RESTARTS   AGE   IP           NODE             NOMINATED NODE   READINESS GATES   LABELS
rng-65d885d498-zl8qq   1/1     Running   0          18h   10.1.0.145   docker-desktop   <none>           <none>            active=yes,app=rng,pod-template-hash=65d885d498
rng-q7274              1/1     Running   0          26m   10.1.0.151   docker-desktop   <none>           <none>            active=yes,app=rng,controller-revision-hash=65d885d498,pod-template-generation=1

# update the selector of service/rng to match labels with active=yes instead of matching labels with app=rng so that the sick pod can be taken out of load balancer
$ kubectl edit service rng
service/rng edited
$ kubectl describe service rng
Name:                     rng
Namespace:                dev
Labels:                   app=rng
                          release=dev
Annotations:              <none>
Selector:                 active=yes
Type:                     ClusterIP
IP Family Policy:         SingleStack
IP Families:              IPv4
IP:                       10.97.48.98
IPs:                      10.97.48.98
Port:                     <unset>  80/TCP
TargetPort:               80/TCP
Endpoints:                10.1.0.145:80,10.1.0.151:80
Session Affinity:         None
Internal Traffic Policy:  Cluster
Events:                   <none>

# now, for treating the sick pod, change it's active label value to "no" so that the load balancer cannot reach it while we treat it
$ kubectl get pods --selector app=rng --show-labels -o wide
NAME                   READY   STATUS    RESTARTS   AGE   IP           NODE             NOMINATED NODE   READINESS GATES   LABELS
rng-65d885d498-zl8qq   1/1     Running   0          19h   10.1.0.145   docker-desktop   <none>           <none>            active=yes,app=rng,pod-template-hash=65d885d498
rng-q7274              1/1     Running   0          64m   10.1.0.151   docker-desktop   <none>           <none>            active=yes,app=rng,controller-revision-hash=65d885d498,pod-template-generation=1

$ kubectl label pod rng-q7274 active=no --overwrite
pod/rng-q7274 labeled

$ kubectl get pods --selector app=rng --show-labels -o wide
NAME                   READY   STATUS    RESTARTS   AGE   IP           NODE             NOMINATED NODE   READINESS GATES   LABELS
rng-65d885d498-zl8qq   1/1     Running   0          19h   10.1.0.145   docker-desktop   <none>           <none>            active=yes,app=rng,pod-template-hash=65d885d498
rng-q7274              1/1     Running   0          68m   10.1.0.151   docker-desktop   <none>           <none>            active=no,app=rng,controller-revision-hash=65d885d498,pod-template-generation=1

$ kubectl describe service rng
Name:                     rng
Namespace:                dev
Labels:                   app=rng
                          release=dev
Annotations:              <none>
Selector:                 active=yes
Type:                     ClusterIP
IP Family Policy:         SingleStack
IP Families:              IPv4
IP:                       10.97.48.98
IPs:                      10.97.48.98
Port:                     <unset>  80/TCP
TargetPort:               80/TCP
Endpoints:                10.1.0.145:80
Session Affinity:         None
Internal Traffic Policy:  Cluster
Events:                   <none>

# now, after treating the sick pod, make it active again
$ kubectl label pod rng-q7274 active=yes --overwrite
pod/rng-q7274 labeled

$ kubectl describe service rng
Name:                     rng
Namespace:                dev
Labels:                   app=rng
                          release=dev
Annotations:              <none>
Selector:                 active=yes
Type:                     ClusterIP
IP Family Policy:         SingleStack
IP Families:              IPv4
IP:                       10.97.48.98
IPs:                      10.97.48.98
Port:                     <unset>  80/TCP
TargetPort:               80/TCP
Endpoints:                10.1.0.145:80,10.1.0.151:80
Session Affinity:         None
Internal Traffic Policy:  Cluster
Events:                   <none>
```
- now, if a pod gets deleted, for any reason, either the replica-set or the daemon-set immediately starts a replacement pod, but it will not add the active=yes label, since we added it via `kubectl` and it's not present in the manifest YAML of the replica-set or daemon-set
```
$ kubectl get pods --selector app=rng --show-labels -o wide
NAME                   READY   STATUS    RESTARTS   AGE   IP           NODE             NOMINATED NODE   READINESS GATES   LABELS
rng-65d885d498-zl8qq   1/1     Running   0          19h   10.1.0.145   docker-desktop   <none>           <none>            active=yes,app=rng,pod-template-hash=65d885d498
rng-q7274              1/1     Running   0          74m   10.1.0.151   docker-desktop   <none>           <none>            active=yes,app=rng,controller-revision-hash=65d885d498,pod-template-generation=1

$ kubectl delete pod rng-65d885d498-zl8qq
pods/rng-65d885d498-zl8qq deleted

$ kubectl get pods --selector app=rng --show-labels -o wide
NAME                   READY   STATUS        RESTARTS   AGE   IP           NODE             NOMINATED NODE   READINESS GATES   LABELS
rng-65d885d498-54mwv   1/1     Running       0          9s    10.1.0.152   docker-desktop   <none>           <none>            app=rng,pod-template-hash=65d885d498
rng-65d885d498-zl8qq   1/1     Terminating   0          19h   10.1.0.145   docker-desktop   <none>           <none>            active=yes,app=rng,pod-template-hash=65d885d498
rng-q7274              1/1     Running       0          74m   10.1.0.151   docker-desktop   <none>           <none>            active=yes,app=rng,controller-revision-hash=65d885d498,pod-template-generation=1

$ kubectl describe service rng
Name:                     rng
Namespace:                dev
Labels:                   app=rng
                          release=dev
Annotations:              <none>
Selector:                 active=yes
Type:                     ClusterIP
IP Family Policy:         SingleStack
IP Families:              IPv4
IP:                       10.97.48.98
IPs:                      10.97.48.98
Port:                     <unset>  80/TCP
TargetPort:               80/TCP
Endpoints:                10.1.0.151:80
Session Affinity:         None
Internal Traffic Policy:  Cluster
Events:                   <none>
```
- thus, even if the replacement pods are running, load-balancer is not able able to point to them and thus we are loosing capacity
- to solve this, the right thing to do is to update the label in the daemon-set and/or replica-set to include the `active=yes` labels in the manifest YAML, instead of just having it in the pods
- thus, after treating the sick pod, add the `active=yes` label to the replica-set and/or daemon-set as well
```sh
$ kubectl edit daemonset rng
daemonset.apps/rng edited
$ kubectl get daemonset rng -o yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  annotations:
    deprecated.daemonset.template.generation: "2"
  creationTimestamp: "2026-01-01T19:25:30Z"
  generation: 2
  labels:
    active: "yes" # adding here just adds or decorates the daemonset, and does nothing to the pods being created with this
    app: rng
  name: rng
  namespace: dev
  resourceVersion: "161149"
  uid: ec4cc3ba-abba-4b1c-bc72-561c11b7db1c
spec:
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: rng
  template:
    metadata:
      creationTimestamp: null
      labels:
        active: "yes" # add this to the pod template
        app: rng
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
  updateStrategy:
    rollingUpdate:
      maxSurge: 0
      maxUnavailable: 1
    type: RollingUpdate
status:
  currentNumberScheduled: 1
  desiredNumberScheduled: 1
  numberAvailable: 1
  numberMisscheduled: 0
  numberReady: 1
  observedGeneration: 2
  updatedNumberScheduled: 1
  
$ kubectl get pods --selector app=rng --show-labels -o wide
NAME                   READY   STATUS        RESTARTS   AGE     IP           NODE             NOMINATED NODE   READINESS GATES   LABELS
rng-65d885d498-n2vgv   1/1     Running       0          3m30s   10.1.0.153   docker-desktop   <none>           <none>            app=rng,pod-template-hash=65d885d498
rng-qvtcv              1/1     Terminating   0          2m7s    10.1.0.154   docker-desktop   <none>           <none>            app=rng,controller-revision-hash=65d885d498,pod-template-generation=1

$ kubectl get pods --selector app=rng --show-labels -o wide
NAME                   READY   STATUS              RESTARTS   AGE     IP           NODE             NOMINATED NODE   READINESS GATES   LABELS
rng-65d885d498-n2vgv   1/1     Running             0          3m36s   10.1.0.153   docker-desktop   <none>           <none>            app=rng,pod-template-hash=65d885d498
rng-lf2kk              0/1     ContainerCreating   0          0s      <none>       docker-desktop   <none>           <none>            active=yes,app=rng,controller-revision-hash=6887d4bfbc,pod-template-generation=2

$ kubectl get pods --selector app=rng --show-labels -o wide
NAME                   READY   STATUS    RESTARTS   AGE     IP           NODE             NOMINATED NODE   READINESS GATES   LABELS
rng-65d885d498-n2vgv   1/1     Running   0          8m51s   10.1.0.153   docker-desktop   <none>           <none>            app=rng,pod-template-hash=65d885d498
rng-lf2kk              1/1     Running   0          5m15s   10.1.0.155   docker-desktop   <none>           <none>            active=yes,app=rng,controller-revision-hash=6887d4bfbc,pod-template-generation=2
```
- the replacement pods are taking time for being created in case of the `daemonset` , once the existing running pod is being deleted as it is waiting for the termination to finish, and then to start creating replacement  pod
- the `rng` container is not handling signals in a good way and thus is termination is taking 30 s
- this is a relatively nice rolling update in case of `daemonset`
```sh
$ kubectl get daeomonset rng -o yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  annotations:
    deprecated.daemonset.template.generation: "2"
  creationTimestamp: "2026-01-01T19:25:30Z"
  generation: 2
  labels:
    active: "yes"
    app: rng
  name: rng
  namespace: dev
  resourceVersion: "161149"
  uid: ec4cc3ba-abba-4b1c-bc72-561c11b7db1c
spec:
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: rng
  template:
    metadata:
      creationTimestamp: null
      labels:
        active: "yes"
        app: rng
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
  updateStrategy: # the update stratergy used is RollingUpdate here
    rollingUpdate:
      maxSurge: 0
      maxUnavailable: 1
    type: RollingUpdate
status:
  currentNumberScheduled: 1
  desiredNumberScheduled: 1
  numberAvailable: 1
  numberMisscheduled: 0
  numberReady: 1
  observedGeneration: 2
  updatedNumberScheduled: 1
```
- it is nice because it avoids completely loosing the capacity all at once
- but if we have lot of nodes, it will take a lot of time
- in case of deployment, the type of rolling update used is different

## Rolling update strategies on Deployment
How should we update a running application?
### Strategy 1: Delete old version,  then deploy new version
- provokes downtime
### Strategy 2: Deploy new version, then delete old version
- uses lot of resources - double footprint
- how do we shift traffic?
### Strategy 3: Replace running pods one at a time
- k8s does it for us
- when we make a change to a deployment, 
	- the deployment controller detects that
	- and scales up the new version, while simultaneously scaling down the old version
	- this is done by creating a new replica-set which starts running new pods. 
		- once all the new pods start running, then k8s will continue the rolling updates i.e. terminating the old pods
```sh
$ watch kubectl get pods,replicasets,deployments --selector app=worker
Every 2.0s: kubectl get pods,replicas… vishal-govind-laptop: 20:14:57
                                                        in 0.110s (0)
NAME                          READY   STATUS    RESTARTS   AGE
pod/worker-5c6f84b477-4fhdc   1/1     Running   0          4m28s
pod/worker-5c6f84b477-56g9l   1/1     Running   0          4m28s
pod/worker-5c6f84b477-5fvhc   1/1     Running   0          5m15s
pod/worker-5c6f84b477-8djj8   1/1     Running   0          4m28s
pod/worker-5c6f84b477-cv2xk   1/1     Running   0          4m28s
pod/worker-5c6f84b477-fk9zt   1/1     Running   0          4m28s
pod/worker-5c6f84b477-gpqr4   1/1     Running   0          4m28s
pod/worker-5c6f84b477-sfcqq   1/1     Running   0          4m28s
pod/worker-5c6f84b477-xgflc   1/1     Running   0          4m28s
pod/worker-5c6f84b477-zwz7c   1/1     Running   0          4m28s

NAME                                DESIRED   CURRENT   READY   AGE
replicaset.apps/worker-5c6f84b477   10        10        10      5m15s

NAME                     READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/worker   10/10   10           10          5m15s
```
- update version of worker to `dockercoins/worker:v0.2`
```sh
$ kubectl get dep
loyment worker -o yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    deployment.kubernetes.io/revision: "2"
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"apps/v1","kind":"Deployment","metadata":{"annotations":{},"labels":{"app":"worker"},"name":"worker","namespace":"dev"},"spec":{"replicas":1,"selector":{"matchLabels":{"app":"worker"}},"template":{"metadata":{"labels":{"app":"worker"}},"spec":{"containers":[{"image":"dockercoins/worker:v0.1","name":"worker"}]}}}}
  creationTimestamp: "2026-01-03T14:39:42Z"
  generation: 3
  labels:
    app: worker
  name: worker
  namespace: dev
  resourceVersion: "2206"
  uid: 2bbfbb4c-ca90-45ff-a2aa-5118121acbce
spec:
  progressDeadlineSeconds: 600
  replicas: 10
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: worker
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: worker
    spec:
      containers:
      - image: dockercoins/worker:v0.2 # updated the image version
        imagePullPolicy: IfNotPresent
        name: worker
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
status:
  availableReplicas: 10
  conditions:
  - lastTransitionTime: "2026-01-03T14:40:33Z"
    lastUpdateTime: "2026-01-03T14:40:33Z"
    message: Deployment has minimum availability.
    reason: MinimumReplicasAvailable
    status: "True"
    type: Available
  - lastTransitionTime: "2026-01-03T14:39:42Z"
    lastUpdateTime: "2026-01-03T14:47:33Z"
    message: ReplicaSet "worker-5bd89bd7fd" has successfully progressed.
    reason: NewReplicaSetAvailable
    status: "True"
    type: Progressing
  observedGeneration: 3
  readyReplicas: 10
  replicas: 10
  updatedReplicas: 10
```

```sh
Every 2.0s: kubectl get pods,replicasets,deployments -… vishal-govind-laptop: 20:17:20

NAME                          READY   STATUS              RESTARTS   AGE
pod/worker-5bd89bd7fd-2j7rl   0/1     ContainerCreating   0          10s
pod/worker-5bd89bd7fd-hn4zs   0/1     ContainerCreating   0          10s
pod/worker-5bd89bd7fd-ppl6w   0/1     ContainerCreating   0          10s
pod/worker-5bd89bd7fd-sbzrp   0/1     ContainerCreating   0          10s
pod/worker-5bd89bd7fd-wq6gt   0/1     ContainerCreating   0          10s
pod/worker-5c6f84b477-4fhdc   1/1     Running             0          6m51s
pod/worker-5c6f84b477-56g9l   1/1     Running             0          6m51s
pod/worker-5c6f84b477-5fvhc   1/1     Terminating         0          7m38s
pod/worker-5c6f84b477-8djj8   1/1     Running             0          6m51s
pod/worker-5c6f84b477-cv2xk   1/1     Running             0          6m51s
pod/worker-5c6f84b477-fk9zt   1/1     Running             0          6m51s
pod/worker-5c6f84b477-gpqr4   1/1     Running             0          6m51s
pod/worker-5c6f84b477-sfcqq   1/1     Running             0          6m51s
pod/worker-5c6f84b477-xgflc   1/1     Terminating         0          6m51s
pod/worker-5c6f84b477-zwz7c   1/1     Running             0          6m51s

NAME                                DESIRED   CURRENT   READY   AGE
replicaset.apps/worker-5bd89bd7fd   5         5         0       10s
replicaset.apps/worker-5c6f84b477   8         8         8       7m38s

NAME                     READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/worker   8/10    5            8           7m38s

```

```sh
Every 2.0s: kubectl get pods,replicasets,deployments -… vishal-govind-laptop: 20:17:37
                                                                         in 0.140s (0)
NAME                          READY   STATUS        RESTARTS   AGE
pod/worker-5bd89bd7fd-2j7rl   1/1     Running       0          27s
pod/worker-5bd89bd7fd-5zr6m   1/1     Running       0          8s
pod/worker-5bd89bd7fd-c5xfg   1/1     Running       0          10s
pod/worker-5bd89bd7fd-hn4zs   1/1     Running       0          27s
pod/worker-5bd89bd7fd-lfv7h   1/1     Running       0          10s
pod/worker-5bd89bd7fd-mqs9r   1/1     Running       0          12s
pod/worker-5bd89bd7fd-ppl6w   1/1     Running       0          27s
pod/worker-5bd89bd7fd-qd2qr   1/1     Running       0          7s
pod/worker-5bd89bd7fd-sbzrp   1/1     Running       0          27s
pod/worker-5bd89bd7fd-wq6gt   1/1     Running       0          27s
pod/worker-5c6f84b477-4fhdc   1/1     Terminating   0          7m8s
pod/worker-5c6f84b477-56g9l   1/1     Terminating   0          7m8s
pod/worker-5c6f84b477-5fvhc   1/1     Terminating   0          7m55s
pod/worker-5c6f84b477-8djj8   1/1     Terminating   0          7m8s
pod/worker-5c6f84b477-cv2xk   1/1     Terminating   0          7m8s
pod/worker-5c6f84b477-fk9zt   1/1     Terminating   0          7m8s
pod/worker-5c6f84b477-gpqr4   1/1     Terminating   0          7m8s
pod/worker-5c6f84b477-sfcqq   1/1     Terminating   0          7m8s
pod/worker-5c6f84b477-xgflc   1/1     Terminating   0          7m8s
pod/worker-5c6f84b477-zwz7c   1/1     Terminating   0          7m8s

NAME                                DESIRED   CURRENT   READY   AGE
replicaset.apps/worker-5bd89bd7fd   10        10        10      27s
replicaset.apps/worker-5c6f84b477   0         0         0       7m55s

NAME                     READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/worker   10/10   10           10          7m55s
```
- edit to v0.3
```sh
Every 2.0s: kubectl get deployments,replicasets,pods --selector ap… vishal-govind-laptop: 20:39:19
                                                                                     in 0.110s (0)
NAME                     READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/worker   8/10    5            8           29m

NAME                                DESIRED   CURRENT   READY   AGE
replicaset.apps/worker-5bd89bd7fd   8         8         8       22m
replicaset.apps/worker-5c6f84b477   0         0         0       29m
replicaset.apps/worker-68f978b56c   5         5         0       86s

NAME                          READY   STATUS             RESTARTS   AGE
pod/worker-5bd89bd7fd-2j7rl   1/1     Running            0          22m
pod/worker-5bd89bd7fd-5zr6m   1/1     Running            0          21m
pod/worker-5bd89bd7fd-c5xfg   1/1     Running            0          21m
pod/worker-5bd89bd7fd-hn4zs   1/1     Running            0          22m
pod/worker-5bd89bd7fd-lfv7h   1/1     Running            0          21m
pod/worker-5bd89bd7fd-mqs9r   1/1     Running            0          21m
pod/worker-5bd89bd7fd-ppl6w   1/1     Running            0          22m
pod/worker-5bd89bd7fd-qd2qr   1/1     Running            0          21m
pod/worker-68f978b56c-8dc76   0/1     ImagePullBackOff   0          86s
pod/worker-68f978b56c-jjzq4   0/1     ImagePullBackOff   0          86s
pod/worker-68f978b56c-qqwkr   0/1     ImagePullBackOff   0          86s
pod/worker-68f978b56c-vhpcc   0/1     ImagePullBackOff   0          86s
pod/worker-68f978b56c-z7hcx   0/1     ImagePullBackOff   0          85s
```
- here our rolling update doesn't seem to proceed beyond this to v0.3
- investigating issue
```sh
$ kubectl describe pod worker-68f978b56c-8dc76
Name:             worker-68f978b56c-8dc76
Namespace:        dev
Priority:         0
Service Account:  default
Node:             docker-desktop/192.168.65.3
Start Time:       Sat, 03 Jan 2026 20:37:54 +0530
Labels:           app=worker
                  pod-template-hash=68f978b56c
Annotations:      <none>
Status:           Pending
IP:               10.1.0.221
IPs:
  IP:           10.1.0.221
Controlled By:  ReplicaSet/worker-68f978b56c
Containers:
  worker:
    Container ID:
    Image:          dockercoins/worker:v0.3
    Image ID:
    Port:           <none>
    Host Port:      <none>
    State:          Waiting
      Reason:       ErrImagePull
    Ready:          False
    Restart Count:  0
    Environment:    <none>
    Mounts:
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-k7s22 (ro)
Conditions:
  Type                        Status
  PodReadyToStartContainers   True
  Initialized                 True
  Ready                       False
  ContainersReady             False
  PodScheduled                True
Volumes:
  kube-api-access-k7s22:
    Type:                    Projected (a volume that contains injected data from multiple sources)
    TokenExpirationSeconds:  3607
    ConfigMapName:           kube-root-ca.crt
    ConfigMapOptional:       <nil>
    DownwardAPI:             true
QoS Class:                   BestEffort
Node-Selectors:              <none>
Tolerations:                 node.kubernetes.io/not-ready:NoExecute op=Exists for 300s
                             node.kubernetes.io/unreachable:NoExecute op=Exists for 300s
Events:
  Type     Reason     Age                   From               Message
  ----     ------     ----                  ----               -------
  Normal   Scheduled  6m26s                 default-scheduler  Successfully assigned dev/worker-68f978b56c-8dc76 to docker-desktop
  Normal   Pulling    3m8s (x5 over 6m25s)  kubelet            Pulling image "dockercoins/worker:v0.3"
  Warning  Failed     3m7s (x5 over 6m18s)  kubelet            Failed to pull image "dockercoins/worker:v0.3": Error response from daemon: failed to resolve reference "docker.io/dockercoins/worker:v0.3": docker.io/dockercoins/worker:v0.3: not found
  Warning  Failed     3m7s (x5 over 6m18s)  kubelet            Error: ErrImagePull
  Warning  Failed     75s (x19 over 6m17s)  kubelet            Error: ImagePullBackOff
  Normal   BackOff    47s (x21 over 6m17s)  kubelet            Back-off pulling image "dockercoins/worker:v0.3"
```
- Rolling back to v0.1
```sh
$ git checkout v0.1
$ kubectl apply -f k8s/dockercoins.yaml
$ watch kubectl get deployments,replicasets,pods --selector app=worker
Every 2.0s: kubectl get deployments,rep… vishal-govind-laptop: 20:51:37
                                                          in 0.168s (0)
NAME                     READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/worker   10/10   10           10          41m

NAME                                DESIRED   CURRENT   READY   AGE
replicaset.apps/worker-5bd89bd7fd   0         0         0       34m
replicaset.apps/worker-5c6f84b477   10        10        10      41m
replicaset.apps/worker-68f978b56c   0         0         0       13m

NAME                          READY   STATUS        RESTARTS   AGE
pod/worker-5bd89bd7fd-2j7rl   1/1     Terminating   0          34m
pod/worker-5bd89bd7fd-5zr6m   1/1     Terminating   0          34m
pod/worker-5bd89bd7fd-c5xfg   1/1     Terminating   0          34m
pod/worker-5bd89bd7fd-hn4zs   1/1     Terminating   0          34m
pod/worker-5bd89bd7fd-lfv7h   1/1     Terminating   0          34m
pod/worker-5bd89bd7fd-mqs9r   1/1     Terminating   0          34m
pod/worker-5bd89bd7fd-ppl6w   1/1     Terminating   0          34m
pod/worker-5bd89bd7fd-qd2qr   1/1     Terminating   0          34m
pod/worker-5c6f84b477-24ghp   1/1     Running       0          21s
pod/worker-5c6f84b477-2k2sm   1/1     Running       0          21s
pod/worker-5c6f84b477-88x78   1/1     Running       0          17s
pod/worker-5c6f84b477-jk6hx   1/1     Running       0          21s
pod/worker-5c6f84b477-kplkh   1/1     Running       0          18s
pod/worker-5c6f84b477-mjsd6   1/1     Running       0          17s
pod/worker-5c6f84b477-rchj6   1/1     Running       0          17s
pod/worker-5c6f84b477-sd7m4   1/1     Running       0          21s
pod/worker-5c6f84b477-t5r2s   1/1     Running       0          21s
pod/worker-5c6f84b477-tz2g7   1/1     Running       0          18s
```

```sh
Every 2.0s: kubectl get deployments,rep… vishal-govind-laptop: 20:52:16

NAME                     READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/worker   10/10   10           10          42m

NAME                                DESIRED   CURRENT   READY   AGE
replicaset.apps/worker-5bd89bd7fd   0         0         0       35m
replicaset.apps/worker-5c6f84b477   10        10        10      42m
replicaset.apps/worker-68f978b56c   0         0         0       14m

NAME                          READY   STATUS    RESTARTS   AGE
pod/worker-5c6f84b477-24ghp   1/1     Running   0          60s
pod/worker-5c6f84b477-2k2sm   1/1     Running   0          60s
pod/worker-5c6f84b477-88x78   1/1     Running   0          56s
pod/worker-5c6f84b477-jk6hx   1/1     Running   0          60s
pod/worker-5c6f84b477-kplkh   1/1     Running   0          57s
pod/worker-5c6f84b477-mjsd6   1/1     Running   0          56s
pod/worker-5c6f84b477-rchj6   1/1     Running   0          56s
pod/worker-5c6f84b477-sd7m4   1/1     Running   0          60s
pod/worker-5c6f84b477-t5r2s   1/1     Running   0          60s
pod/worker-5c6f84b477-tz2g7   1/1     Running   0          57s
```

### `kubectl rollout`
```sh
$ kubectl rollout
Manage the rollout of one or many resources.

 Valid resource types include:

  *  deployments
  *  daemonsets
  *  statefulsets

Examples:
  # Rollback to the previous deployment
  kubectl rollout undo deployment/abc

  # Check the rollout status of a daemonset
  kubectl rollout status daemonset/foo

  # Restart a deployment
  kubectl rollout restart deployment/abc

  # Restart deployments with the 'app=nginx' label
  kubectl rollout restart deployment --selector=app=nginx

Available Commands:
  history       View rollout history
  pause         Mark the provided resource as paused
  restart       Restart a resource
  resume        Resume a paused resource
  status        Show the status of the rollout
  undo          Undo a previous rollout

Usage:
  kubectl rollout SUBCOMMAND [options]

Use "kubectl rollout <command> --help" for more information about a given command.
Use "kubectl options" for a list of global command-line options (applies to all commands).
```
- we can do `kubectl rollout undo deployment worker`
```sh
$ 
Every 2.0s: kubectl get deployments,replica… vishal-govind-laptop: 20:58:03
                                                              in 0.129s (0)
NAME                     READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/worker   8/10    5            8           48m

NAME                                DESIRED   CURRENT   READY   AGE
replicaset.apps/worker-5bd89bd7fd   0         0         0       40m
replicaset.apps/worker-5c6f84b477   8         8         8       48m
replicaset.apps/worker-68f978b56c   5         5         0       20m

NAME                          READY   STATUS              RESTARTS   AGE
pod/worker-5c6f84b477-24ghp   1/1     Running             0          6m47s
pod/worker-5c6f84b477-2k2sm   1/1     Running             0          6m47s
pod/worker-5c6f84b477-6lkwb   1/1     Terminating         0          49s
pod/worker-5c6f84b477-88x78   1/1     Running             0          6m43s
pod/worker-5c6f84b477-jk6hx   1/1     Running             0          6m47s
pod/worker-5c6f84b477-kplkh   1/1     Running             0          6m44s
pod/worker-5c6f84b477-mhsms   1/1     Terminating         0          49s
pod/worker-5c6f84b477-sd7m4   1/1     Running             0          6m47s
pod/worker-5c6f84b477-t5r2s   1/1     Running             0          6m47s
pod/worker-5c6f84b477-tz2g7   1/1     Running             0          6m44s
pod/worker-68f978b56c-96494   0/1     ErrImagePull        0          7s
pod/worker-68f978b56c-fvg2r   0/1     ImagePullBackOff    0          7s
pod/worker-68f978b56c-rdbzw   0/1     ContainerCreating   0          7s
pod/worker-68f978b56c-tzn8f   0/1     ErrImagePull        0          7s
pod/worker-68f978b56c-vhllw   0/1     ContainerCreating   0          7s


$ kubectl rollout undo deployment worker
deployment.apps/worker rolled back

$ watch kubectl get deployments,replicasets,pods --selector app=worker
Every 2.0s: kubectl get deployments,replicasets,pods --selector… vishal-govind-laptop: 21:02:25

NAME                     READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/worker   8/10    10           8           52m

NAME                                DESIRED   CURRENT   READY   AGE
replicaset.apps/worker-5bd89bd7fd   0         0         0       45m
replicaset.apps/worker-5c6f84b477   10        10        8       52m
replicaset.apps/worker-68f978b56c   0         0         0       24m

NAME                          READY   STATUS              RESTARTS   AGE
pod/worker-5c6f84b477-24ghp   1/1     Running             0          11m
pod/worker-5c6f84b477-2k2sm   1/1     Running             0          11m
pod/worker-5c6f84b477-799dm   0/1     ContainerCreating   0          3s
pod/worker-5c6f84b477-88x78   1/1     Running             0          11m
pod/worker-5c6f84b477-bh4t2   0/1     ContainerCreating   0          3s
pod/worker-5c6f84b477-jk6hx   1/1     Running             0          11m
pod/worker-5c6f84b477-kplkh   1/1     Running             0          11m
pod/worker-5c6f84b477-kw8b4   1/1     Terminating         0          56s
pod/worker-5c6f84b477-ltl86   1/1     Terminating         0          16s
pod/worker-5c6f84b477-sd7m4   1/1     Running             0          11m
pod/worker-5c6f84b477-t5r2s   1/1     Running             0          11m
pod/worker-5c6f84b477-t9jxs   1/1     Terminating         0          16s
pod/worker-5c6f84b477-twtqr   1/1     Terminating         0          57s
pod/worker-5c6f84b477-tz2g7   1/1     Running             0          11m
pod/worker-68f978b56c-5mms2   0/1     Terminating         0          9s
pod/worker-68f978b56c-68777   0/1     Terminating         0          9s
pod/worker-68f978b56c-gdw4z   0/1     Terminating         0          9s
pod/worker-68f978b56c-jm6j6   0/1     Terminating         0          9s
pod/worker-68f978b56c-wmt5n   0/1     Terminating         0          9s
```

## Update strategy in deployment - `maxSurge` and `maxUnavailable`
```sh
$ kubectl get deployment worker -o yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    deployment.kubernetes.io/revision: "18"
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"apps/v1","kind":"Deployment","metadata":{"annotations":{},"labels":{"app":"worker"},"name":"worker","namespace":"dev"},"spec":{"replicas":10,"selector":{"matchLabels":{"app":"worker"}},"template":{"metadata":{"labels":{"app":"worker"}},"spec":{"containers":[{"image":"dockercoins/worker:v0.1","name":"worker"}]}}}}
  creationTimestamp: "2026-01-03T14:39:42Z"
  generation: 19
  labels:
    app: worker
  name: worker
  namespace: dev
  resourceVersion: "7576"
  uid: 2bbfbb4c-ca90-45ff-a2aa-5118121acbce
spec:
  progressDeadlineSeconds: 600
  replicas: 10
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: worker
  strategy: # talks about the update stratergy during rolling update, which is different than that used in daemon-set
    rollingUpdate:
      maxSurge: 25% # daemon set has maxSurge = 0
      maxUnavailable: 25% # daemon set has maxUnavailable = 1
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: worker
    spec:
      containers:
      - image: dockercoins/worker:v0.1
        imagePullPolicy: IfNotPresent
        name: worker
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
status:
  availableReplicas: 10
  conditions:
  - lastTransitionTime: "2026-01-03T14:40:33Z"
    lastUpdateTime: "2026-01-03T14:40:33Z"
    message: Deployment has minimum availability.
    reason: MinimumReplicasAvailable
    status: "True"
    type: Available
  - lastTransitionTime: "2026-01-03T15:21:16Z"
    lastUpdateTime: "2026-01-03T15:32:27Z"
    message: ReplicaSet "worker-5c6f84b477" has successfully progressed.
    reason: NewReplicaSetAvailable
    status: "True"
    type: Progressing
  observedGeneration: 19
  readyReplicas: 10
  replicas: 10
  updatedReplicas: 10
```
- Two parameters determine the pace of rollout: `maxUnavailable` and `maxSurge`
- They can be specified in 
	- absolute number of pods
	- or percentage of replicas count
- At any give time
	- there will always be at least `replicas - maxUnavailable` pods available
	- there will never be more than `replicas + maxSurge` pods in total
	- there will therefore be up to `maxUnavailable + maxSurge` pods being updated
- if `replicas = 10` and `maxUnavailable = 25%`
	- it means we can loose up to 25% of the capacity 
		- 25% of 10 = 2.5 i.e. floor of 2.5 i.e. 2
		- i.e. we can loose 2 pods
		- thus, during rolling update of a deployment, we have 8 pods available, and 2 pods unavailable
- `maxSurge` indicates how much extra resources we can use during the rolling update
	- for daemon-sets, maxSurge is 0, thus, it waits for entire old version pod to terminate and then starts new version pod
- if `maxSurge = 25%` and `replicas = 10`
	- i.e. 25% of 10 = 2.5 i.e. ceil of 2.5 is 3
	- i.e. we can use 3 extra pods to the rolling update
- examples of strategies
	- if we don't want to loose any capacity during the rolling update, we can put `maxUnavailable` to 0
	- if we want rolling update to be really fast, we can put `maxSurge` to 100%
		- i.e when we start the rolling update, it's going to create all the new 10 (`replicas`) pods right away
		- it's going to be very fast but is going to use a lot of extra resources
		- if it's a deployment with few 3-4 pods, that's ok, but not if there are 1000s of pods
- what exactly will happen if we have lot of pods and have high `maxSurge`, during a rolling update
	- there would not be exactly a lot of load on the control plane
	- at that scale, most likely we will be using cluster auto-scaling
	- and in case of such sudden need of extra capacity, the cluster auto-scaler will kick in and provision a bunch of extra nodes
	- the problem is we don't need these extra nodes for long term
	- we only need them for a few minutes, while we do the rolling update
	- so, the new nodes will come up, which will be filled up with the new pods
	- then, the old pods are shut down
	- now, we end up with a cluster which is highly **fragmented**, i.e. not all nodes are fully utilized by capacity
	- in a perfect scenario, when there's no rolling update, the scheduler in the control panel maintains the pods among the clusters with as much less fragmentation as possible utilizing node's full or max capacity, so that we use as lesser nodes as possible and we don't pay for unused node capacity
	- now, after the rolling update, we end up with a cluster that looks like a cheese with wholes
	- thus, after a rolling update with high `maxSurge`, we've lot of nodes that are 60-70% full and we end up paying for a lot of extra capacity for nothing
	- thus, it's not a good idea to go with high `maxSurge` on really big deployments
- default values for `maxSurge` and `maxUnavailable` for a deployment is 25% each
	- for no loss of capacity during rolling update, we can go with 0 `maxUnavailable`
	- for `maxSurge` during rolling update, it really depends on the kind of work load we have


### Older replica-sets during rolling updates - `revisionHistoryLimit`
- during rolling updates of deployments, the older replica-sets are kept
- by default last 10 revisions of replica-sets are kept
- this default 10 can be found in the deployment spec in the `spec.revisionHistoryLimit`
```sh
$ kubectl get deployment worker -o yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    deployment.kubernetes.io/revision: "18"
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"apps/v1","kind":"Deployment","metadata":{"annotations":{},"labels":{"app":"worker"},"name":"worker","namespace":"dev"},"spec":{"replicas":10,"selector":{"matchLabels":{"app":"worker"}},"template":{"metadata":{"labels":{"app":"worker"}},"spec":{"containers":[{"image":"dockercoins/worker:v0.1","name":"worker"}]}}}}
  creationTimestamp: "2026-01-03T14:39:42Z"
  generation: 19
  labels:
    app: worker
  name: worker
  namespace: dev
  resourceVersion: "7576"
  uid: 2bbfbb4c-ca90-45ff-a2aa-5118121acbce
spec:
  progressDeadlineSeconds: 600
  replicas: 10
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: worker
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: worker
    spec:
      containers:
      - image: dockercoins/worker:v0.1
        imagePullPolicy: IfNotPresent
        name: worker
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
status:
  availableReplicas: 10
  conditions:
  - lastTransitionTime: "2026-01-03T14:40:33Z"
    lastUpdateTime: "2026-01-03T14:40:33Z"
    message: Deployment has minimum availability.
    reason: MinimumReplicasAvailable
    status: "True"
    type: Available
  - lastTransitionTime: "2026-01-03T15:21:16Z"
    lastUpdateTime: "2026-01-03T15:32:27Z"
    message: ReplicaSet "worker-5c6f84b477" has successfully progressed.
    reason: NewReplicaSetAvailable
    status: "True"
    type: Progressing
  observedGeneration: 19
  readyReplicas: 10
  replicas: 10
  updatedReplicas: 10
```