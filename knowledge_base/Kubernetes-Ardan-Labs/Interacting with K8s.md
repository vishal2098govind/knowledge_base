#k8s

- How do we talk to k8s clusters -> use k8s API
	- using `kubectl` tool
- using k8s API, we will manipulate **resources**
	- **nodes** - a machine, that belongs to the cluster
	- **pods** - group of containers running together on a node
	- **services** - a **network thing**
		- can be a **load balancer**
		- can be a **DNS entry**
		- can be a **static IP address**

## Node, Pod, Container
![[Pasted image 20251228024459.png]]
- in a pod, very often, we will see only one container

## Use cases of multiple containers in a pod
### Side car
- **side-car** patterns or side-kick containers ( #side-car )
- the idea is that we have the main container `nginx`
- and then, when we run in production, we add another container, here we have `logger`, which is here to collect the logs of `nginx` and send them to some centralized logging platform
- it's called `side-car`
- because, here, `nginx` can run on it's own, but we can add the side-car container to add some extra features
- we can have lot's of side cars for like
	- logging
	- metrics
	- service meshes
	- synchronizing configuration with a bigger system

## Pod
- all containers in a pod have same IP address
- thus, the IP address doesn't belong to the container, but pod

## Scaling in Pods
- How would we scale the pod?
- **Do** create additional pods
	- each pod can be on a different node
	- each pod will have it's own IP address
- **Do not** add more `nginx` containers in the pod, or else
	- all `nginx` controllers would be on same node
	- they would all have same IP address, resulting in **address already in use** errors

## Which containers to put together in same pods or separate pods
- e.g. should we have a web application and cache (Redis or Memcached) together?
- generally we put in different pods, because it will let us scale them separately
	- can have 10 web frontends and 3 caches - possible
- Putting them **in same pod** means:
	- they have to be scaled together
	- they can **communicate very efficiently** over localhost
	- this is extremely rare use case if done for communication efficiency
	- mostly multiple containers in same pods is just for implementing side-cars
- Putting them **in different pods** means:
	- they can be scaled separately
	- they must communicate over remote IP addresses
		- incurring **more latency**, lower performance
- both scenarios makes sense, depending on our goals