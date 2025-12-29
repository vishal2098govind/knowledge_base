#k8s 
## `kubectl describe` and events

> in case of any issues, `kubectl logs` might not always show anything useful, however, the events section in the `kubectl describe` should show more useful stuffs in such situations

>thus, prefer `kubectl describe` to start with, during any investigation with k8s issues

```sh
$ kubectl describe pod worker-5c6f84b477-bjg5w
Name:             worker-5c6f84b477-bjg5w
Namespace:        default
Priority:         0
Service Account:  default
Node:             docker-desktop/192.168.65.3
Start Time:       Mon, 29 Dec 2025 01:23:17 +0530
Labels:           app=worker
                  pod-template-hash=5c6f84b477
Annotations:      <none>
Status:           Running
IP:               10.1.0.28
IPs:
  IP:           10.1.0.28
Controlled By:  ReplicaSet/worker-5c6f84b477
Containers:
  worker:
    Container ID:   docker://8ef8939332ee1c072fecc6d8e3568ffb4cd66e21b4f4583b6bff5d75c044b65f
    Image:          dockercoins/worker:v0.1
    Image ID:       docker-pullable://dockercoins/worker@sha256:cd4651fd7b077fbcfe7afc677e46fdef6c8cbd0ea0ef27a4f8a5b6c1e2d9da13
    Port:           <none>
    Host Port:      <none>
    State:          Running
      Started:      Mon, 29 Dec 2025 01:23:18 +0530
    Ready:          True
    Restart Count:  0
    Environment:    <none>
    Mounts:
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-9gq9z (ro)
Conditions:
  Type                        Status
  PodReadyToStartContainers   True
  Initialized                 True
  Ready                       True
  ContainersReady             True
  PodScheduled                True
Volumes:
  kube-api-access-9gq9z:
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
  Type    Reason     Age   From               Message
  ----    ------     ----  ----               -------
  Normal  Scheduled  56s   default-scheduler  Successfully assigned default/worker-5c6f84b477-bjg5w to docker-desktop
  Normal  Pulled     56s   kubelet            Container image "dockercoins/worker:v0.1" already present on machine
  Normal  Created    55s   kubelet            Created container: worker
  Normal  Started    55s   kubelet            Started container worker
```
- these events in the `kubectl describe` are the internal log of events, which means these are stored within the control plane
	- Good news:
		- this is a great news because, even if we don't use or setup our own logging system with Loki or Elasticsearch or Datadog, we will still be able to see these events
		- thus, even on a very basic bare bones k8s cluster, we can see these events
	- Bad news:
		- these events are only stored **for one hour** within the control plane, within the `etcd`
		- if same event happens more than once within that hour, most recent event along with oldest time, stay in the `etcd`
		- we can change TTL or expiration time of these events in `etcd`
```sh
Events:
  Type     Reason   Age                   From     Message
  ----     ------   ----                  ----     -------
  Warning  Failed   49m (x100 over 9h)    kubelet  Error: ErrImagePull
  Normal   BackOff  5m37s (x605 over 9h)  kubelet  Back-off pulling image "dockercoins/workeret:v0.1"
  Warning  Failed   5m37s (x605 over 9h)  kubelet  Error: ImagePullBackOff
  Normal   Pulling  38s (x113 over 9h)    kubelet  Pulling image "dockercoins/workeret:v0.1"
  Warning  Failed   35s (x110 over 9h)    kubelet  Failed to pull image "dockercoins/workeret:v0.1": Error response from daemon: pull access denied for dockercoins/workeret, repository does not exist or may require 'docker login'
```

## `ImagePullBackOff` and Retries
```sh
$ kubectl create deployment workeret --image dockercoins/workeret:v0.1
$ kubectl get pods
NAME                       READY   STATUS             RESTARTS   AGE
blue-5c986bd7bf-tgh7s      1/1     Running            0          7h36m
hasher-99bbd4bb-c5q7j      1/1     Running            0          104m
red-77f6d65f98-454qr       1/1     Running            0          7h26m
red-77f6d65f98-pxlz4       1/1     Running            0          7h35m
red-77f6d65f98-spv8s       1/1     Running            0          7h26m
redis-7b47f84cc4-w9fcl     1/1     Running            0          102m
rng-65d885d498-sqjrp       1/1     Running            0          110m
webui-74bb6bbc59-m5n8m     1/1     Running            0          110m
worker-5c6f84b477-pwzll    1/1     Running            0          13m
workeret-8d57757c8-2ztpr   0/1     ImagePullBackOff   0          7m
$ kubectl describe pod workeret-8d57757c8-2ztpr
Name:             workeret-8d57757c8-2ztpr
Namespace:        default
Priority:         0
Service Account:  default
Node:             docker-desktop/192.168.65.3
Start Time:       Mon, 29 Dec 2025 01:49:47 +0530
Labels:           app=workeret
                  pod-template-hash=8d57757c8
Annotations:      <none>
Status:           Pending
IP:               10.1.0.30
IPs:
  IP:           10.1.0.30
Controlled By:  ReplicaSet/workeret-8d57757c8
Containers:
  workeret:
    Container ID:
    Image:          dockercoins/workeret:v0.1
    Image ID:
    Port:           <none>
    Host Port:      <none>
    State:          Waiting
      Reason:       ImagePullBackOff
    Ready:          False
    Restart Count:  0
    Environment:    <none>
    Mounts:
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-j2b6w (ro)
Conditions:
  Type                        Status
  PodReadyToStartContainers   True
  Initialized                 True
  Ready                       False
  ContainersReady             False
  PodScheduled                True
Volumes:
  kube-api-access-j2b6w:
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
  Type     Reason     Age                     From               Message
  ----     ------     ----                    ----               -------
  Normal   Scheduled  7m33s                   default-scheduler  Successfully assigned default/workeret-8d57757c8-2ztpr to docker-desktop
  Normal   Pulling    4m29s (x5 over 7m32s)   kubelet            Pulling image "dockercoins/workeret:v0.1"
  Warning  Failed     4m27s (x5 over 7m30s)   kubelet            Failed to pull image "dockercoins/workeret:v0.1": Error response from daemon: pull access denied for dockercoins/workeret, repository does not exist or may require 'docker login'
  Warning  Failed     4m27s (x5 over 7m30s)   kubelet            Error: ErrImagePull
  Warning  Failed     2m17s (x20 over 7m29s)  kubelet            Error: ImagePullBackOff
  Normal   BackOff    2m5s (x21 over 7m29s)   kubelet            Back-off pulling image "dockercoins/workeret:v0.1"
```
- the error in the pod was going back and forth between `ImagePullBackOff` and `ErrImagePull`
- for almost every error condition in k8s, at least for those which are not permanent, we are going to have two error conditions
	- the first one: `ErrImagePull` 
		- that means there was an error pulling the image
		- the most basic information about the issue
	- the second one: `ImagePullBackOff`
		- on receiving an error (`ErrImagePull` here), k8s retries
		- if k8s fails again, it retries again
		- if fails one more time, it retries one more time
		- k8s keeps retrying **forever**
		- we have many many systems like this in k8s, 
			- where the container crashes, we restart it
			- if we fail to provision something now, we retry later
			- and so on

## Why do we have retries forever in k8s
- This idea comes from the fact that, if we have a system that is big enough with enough components, 
	- **statistically**, we will always have somethings somewhere having problems or errors
	- think for instance about a data center, with 1000s of servers
	- in each server, we have maybe 10 to 20 disks
	- if we think about how often do we need to go to the data center to replace the disks that have failed, it is very very rare 
	- if we look at the numbers claimed by the cloud providers, it would say something like 1M hours of average time between failures of those disks
	- at first we could think, well, in our whole lifetime, we might never see a disk failure!! i.e. maybe I'll never have a disk that fails in my computer
		- although that's not an excuse to not do backups btw
	- However, if we get some stats, and we have not just one disk in our machine, but we are in charge of a data center with 100s and 1000s of disks, how often will disks fail?
		- multiple times a day
		- multiple times a week
		- i.e. if average failure time is 1M hours and we have 100,000 disks, statistically we have one disk failing every 10 hours
		- i.e. more than 10 in a week
		- every week we literally go to the data center with shopping cart full of disks and change failed disks all around the servers
	- Let's say our data center is running a big distributed application, like twitter or whatever, and think that if one disk being down
		- takes down the whole service, 
		- our application would be down all the time, 
		- will have outages all the time
	- so of course, we need to put systems in place to have **Redundancy** ( #redundancy )
		- for redundancy in disks, we have RAID (Redundant Array of Independent Disks)
		- Simplest form of RAID is RAID-1, is called **Disk Mirroring**
			- there are two disks, they are exact copy of each other
			- if one fails, we still have other one as backup
			- until we replace the failing disk and re-sync the data
		- in case of API backend servers, we have around 20 replicas (redundancy) of the same API, so that  if one of these API servers fails, we can replace or take it out of the load balancer
			- similar to if we go to some administration to renew our passport, we might have 10 people working there and if one person is not present on that day, the remaining 9 can still take care of the work, it doesn't stop everyone else
			- same idea of redundancy
			- we want to design our services to be resilient so that if there is a problem, it doesn't stop our service and everything else
- thus, one part of this strategy is, when something doesn't work, instead of completely giving up immediately, k8s says, that thing didn't work just now, but, k8s will rely on the fact that, maybe in 10 seconds, 10 minutes, 10 hours, but at some point in future, it's going to work again
	- e.g. if we call a friend for hanging out, and they don't pickup
	- we are not going to just give up after first attempt
	- we retry in sometime again
		- first retry after maybe 5 min
		- if still failed, retry after maybe couple of hours
		- if still failed, retry after maybe couple of days or a week
		- if still failed, maybe escalate further, get the correct phone number and try to reach them in other ways
- this is exactly what k8s is doing here, where if pulling an image doesn't work, it retries assuming that 
	- maybe there's a small glitch with the registry
	- maybe the registry has something more important to do than serving an image, at that moment of time when k8s tried pulling
	- so we might retry after 20 seconds
	- and so on...
	- so k8s double the wait time at each failure - **exponential back-off**
	- eventually, at some point, it would reach a point where it feels it's been really waiting for long time, and thinks there is something weird going on, and that's what k8s reports as **`ImagePullBackOff`**