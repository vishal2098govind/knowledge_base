#k8s #capacity-planning

## Resource Limits
- we need to specify **limits** and/or **requests** for **CPU** and/or **memory** to containers
- can specify in all combinations
	- CPU-limit, CPU requests
	- Memory-limit, Memory-requests
	- CPU-limit, Memory-requests
	- CPU-limit and CPU requests, Memory Limits
	- etc
- or we can specify none of them

## CPU vs memory
- CPU is **compressible resource**
- Memory is **incompressible resource**
- here, **compressibility** means what happens if we have too much requests and not enough resources
- e.g. if we have a compute job, that needs 4 CPUs and 4GB RAM
	- if we run it on a machine that has plenty of RAM ~10GB or 20GB, and only 2 CPUs
		- then, the job will run twice as slow
		- i.e. if we don't have enough CPU, it just get's bit slower
		- but there is no harm done
			- unless we have an SLA (service-level-agreement) that we need to serve the request in that many milliseconds
	- with not enough memory on a machine, say only 2GB RAM but we need 4GB RAM
		- then the job just cannot run
		- end of the story!
		- one solution: use swapping
			- when we don't have enough RAM, we write the overflow to the disk
			- this will be very slow due to the involvement of disk, which is much slower than memory
		- another solution is terminate a process and reclaim all it's memory
			- OOM or out of memory killer on linux

## Mechanisms for overcoming termination on reaching memory limits
- **Limits** are "hard limits"
	- a container exceeding it's memory limit is killed
	- a container exceeding it's CPU limit is throttled
		- i.e. the container will not be able to increase the limit
- **Requests** are reservations made by containers
	- a container can use more than it's requested CPU or RAM amounts
	- a container using less than what it's requested should never be killed or throttled
- on a given node, the sum of pod requests cannot be higher than the node size

### Pod QoS - Quality of Service
- each pod is assigned a QoS class (visible in `status.qosClass`)
- if limits = requests:
	- as long as the container uses less than limit, it won't be affected
	- if all containers in a pod have limits=requests, QoS is considered "**Guaranteed**"
	- since memory is a incompressible resource, for memory it is **recommended to go** with this limits=requests
- if requests < limits:
	- as long as container uses the request, it won't be affected
	- otherwise it might be killed/evicted if the node gets overloaded
	- if at least one container has (requests < limit), QoS is considered "**Burst-able**"
	- since CPU is a compressible resource, for CPU **we can go** with requests < limits
- If a pod doesn't have any request nor limit, QoS is considered "**BestEffort**"
	- on reaching limits, the pods will be kicked out for other pods to run

## Specifying resources
- Memory is to be specified in bytes
	- 250M (upper-case M)= 250 x 1000 x 1000 = 250 Million or 250 Megabytes
	- 250Mi = 250 x 1024 x 1024 = 250 Mebibytes
	- 250m (lower-case m) = 250 mili = 0.25 bytes = 2 bits
- CPU is indicated using number of cores
	- this can be a decimal value as well
	- `resources.limits.cpu: 1` => 1 CPU cores
	- `resources.limits.cpu: 2` => 2 CPU cores
	- `resources.limits.cpu: 10` => 10 CPU cores
	- `resources.limits.cpu: 0.5` => half of a CPU
```yaml
resources:
	limits:
		memory: 8G
		cpu: 2
```
- there is **no physical reservation** of CPUs and Memory for the containers within the pod
- each time we put a pod that has 8GB of memory request on a node, we decrement that much amount from the available memory on that node
	- it doesn't matter if the container or pod is actually using the requested memory
- Default values used
	- it is recommended to either put nothing, or put everything
	- If we specify a limit without a request:
		- the request is set to the limit
	- If we specify a request without a limit:
		- there will be no limit
		- (which means that the limit will be the size of the node)
	- If we don't specify anything:
		- the request is zero and the limit is the size of the node
		- _Unless there are default values defined for our namespace!_
- We can create `LimitRange` objects to indicate any combination of:
    - min and/or max resources allowed per pod
    - default resource _limits_
    - default resource _requests_
    - maximal burst ratio (_limit/request_)
- LimitRange objects are `namespaced`
- They apply to their namespace only
- this is helpful if all pods in the namespace have same kind of applications or same resource footprint, which is very unusual


## Checking Node and Pod resource usage
```sh
$ kubectl top nodes
NAME             CPU(cores)   CPU(%)   MEMORY(bytes)   MEMORY(%)
docker-desktop   2446m        20%      5747Mi          74
$ kubectl top pods
NAME                      CPU(cores)   MEMORY(bytes)
hasher-99bbd4bb-vfdlp     39m          22Mi
redis-7b47f84cc4-9b7s5    5m           16Mi
rng-65d885d498-f667w      50m          25Mi
webui-74bb6bbc59-478mc    1m           41Mi
worker-5c6f84b477-mmrdn   27m          26Mi
worker-5c6f84b477-pgbfq   30m          26Mi
worker-5c6f84b477-qkrj9   29m          25Mi
worker-5c6f84b477-qr45k   32m          26Mi
```

## Cluster sizing
- what happens when cluster gets full
- when are we out of resources
	- `kubelet` sits on each node, and monitors resource usage
	- for memory and disk, we've a threshold of around 90%, which on reaching, `kubelet` needs to react and free up some memory
	- we've two thresholds
		- soft threshold - `kubelet` doesn't reach immediately when hit, but reacts when we stay too long on soft threshold
		- hard threshold - `kubelet` reacts immediately when hit
	- the pods with "BestEffort" as QoS are the first ones to be removed
	- then, the pods with "Burst-able" as QoS are the ones to be removed
	- by then, we have solved memory memory issues, because after that, we only are left with pods that use less than their memory request
	- if disk usage is too high, k8s will try to terminate pods and then try to **evict pods**
		- evict is technical word for removing the pod
		- which means the pod is **moved to another node**
			- when we evict the pod, the pod is terminated in that node
			- as if we did `kubectl delete pod pod-name` and we have a graceful shutdown
			- and the controller above will create a replacement pod in another node
	- the node is marked as "under pressure"
		- technically, this is a **taint** placed on the node
		- the node marks itself as not ready and asking not to be sent new pods here for a while, and thus the scheduler will avoid assigning pods to that node
	- sometimes a pod cannot be scheduled anywhere after eviction
		- if all nodes are under pressure
		- or the pod requests more resources than available
		- then, the pod remains in `Pending` state until the situation improves