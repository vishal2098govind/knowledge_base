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