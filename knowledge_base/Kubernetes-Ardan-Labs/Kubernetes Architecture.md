#k8s

## What does a k8s cluster look like?
![[Pasted image 20251228002528.png]]
- Control plane (master node)
- nodes
	- bunch of machines in the cluster
	- each node has 
		- containers
		- container engine - e.g. docker
		- operating system - linux, windows
		- infrastructure - physical machines
		- `kubelet` and `kubeproxy`
			- components to connect the node to the cluster

### k8s Arch
![[Pasted image 20251228002635.png]]
- the master node above is the control plane
- the control plane has several components
	- controller manager (controller loops)
	- etcd (key-value DB, SSOT)
	- API Server (REST API)
	- Scheduler (Bind Pod to Node)
- conceptually in k8s, we have a data store (**etcd**), can imagine it like a pin board on which we can pin some notes on that board
- when we want to do something in k8s, the key thing to understand is that, we cannot directly talk or tell a node or component to do something
- we can only write a manifest
	- "I want an nginx container please"
	- put that manifest on that pin board
- there's going to be what we call controllers that will notice the pin board and understand that we want to run an NGINX container
- thus, every communication goes through **etcd**, the pin board
- the user cannot directly talk to **etcd**, but can only through the **api server** in the control plane, as indicated by the blue arrow line
- the api server take the manifest, checks, validates permissions and then forwards to **etcd**
- in case of any issue in manifest, the **api server** rejects the request in manifest
- the **scheduler** is the tetris player, the one who decides which container goes on which node/machine
- the **controller manager** is the one actually doing the work
	- if we have in total 10 containers in our cluster
	- if suddenly one of the machine goes down, which might be having 2 containers of our application
	- now, we have 8 containers in total
	- somebody needs to notice this and react to this and do something about it
	- that somebody is the **controller manager**
	- controller manager, along with many other responsibilities will notice the situation at some point, so it's going to mark that machine as unreachable
	- then, there will be cascade of events
	- since the machine is unreachable, the containers in the machine would also not be unreachable, and redeploy missing containers
	- similar for **autoscaling**
- **kubelet** is basically the k8s agent. 
	- it runs on every node
	- on every node, it connects to the control plane
	- the job of the kubelet is to register the node with **control plane**
	- and for instance, say "hello, I am kubelet on node-1, i have 8 cores, 32 gigs of RAM, do you have any containers for me?"
	- the **API server** would reply "well, welcome onboard node-1, oh yes, you should these containers"
- the only component with persistense or state is "**etcd**"
- everything else is stateless

## Single-node cluster (for development)
- in k8s, the cluster can have very different shapes and forms
- e.g.:
	- development cluster
	- using mini-kube, docker-desktop, kind
	- idea is to put everything in a single node or machine or vm or container
- if we loose that single machine, we loose everything

## Managed k8s cluster (cloud)
![[Pasted image 20251228012234.png]]
- every modern cloud provider has a k8s implementation these days
- EKS in AWS, GKE in GCP, AKS in Azure
- cloud provider takes care of running the control plane, so we don't see the control plane as it would be running on the cloud provider's infra structure
- the nodes run on the normal machines (EC2) of that clone provider

## Self deployed k8s cluster
![[Pasted image 20251228012456.png]]
- here, we don't have high availability on the control plane
- if the control plane node crashes, we loose the entire control plane

## Stacked Control Plane k8s cluster
![[Pasted image 20251228012854.png]]
- control plane on multiple node
- on each node, we put the complete k8s component stack of 
	- `etcd`
	- api server
	- controller manager
	- scheduler
- there's an API Load balancer, which means when we connect to the API server, we actually connect to the load balancer, which takes care of forwarding it to one of the API servers on control plane nodes
- there's still a single point of failure here, when API Load balancer goes down
- but the big difference is that we know how to make a highly available TCP load balancers, as it is a solved problem
	- whether we are running in the cloud and then we can use the cloud provider's load balancers
	- or whether we are running in a data center, and we can use HAProxy etc
	- either way, we can have a highly available load balancer


## Most self hosted k8s clusters at scale would look like this instead
![[Pasted image 20251228013606.png]]
- different components of the control plane scale in very different manners
- `etcd`
	- highly available key-value store, which leverages a protocol called RAFT
	- in a RAFT cluster, we need to have a majority of nodes running all the time
	- for majority nodes, atleast 3 nodes is needed, so that if one node crashes, majority of 2 would be running
	- sometimes 5 nodes
	- very rarely more than 5
	- the more nodes for `etcd`, the slower is the performance
		- because on writing on one node of etcd, it has to replicate on all other nodes
		- thus, etcd cluster cannot run well on low latency as it has to replicate all writes on all other nodes
- for controller manager and scheduler, we typically see 2 nodes
	- they use an active-stand-by model
	- if we have 5 schedulers, one of them is going to be active, and others are just going to be watching and waiting, thus no need of more than 2, typically
- the API server uses active-active redundancy
	- can load balance traffic across multiple instances
	- since, API server also acts like a cache for `etcd`, it makes sense to add more API servers as the load on the control plane gets higher
- on big k8s clusters, something similar to this is what we would find
	- what does **big** cluster mean?
	- what's the threshold to start thinking about something like this?
	- we are talking about 100s of nodes
	- if we just have a cluster with say 10 nodes or so, we can stay at stacked-control-plane + workers or cloud provider managed k8s cluster
	- e.g. of big clusters running this arch at that scale:
		- Datadog - the metrics company

## Control plane and apps running on same nodes
![[Pasted image 20251228020159.png]]
- this is what we would end up with when running k8s on Raspberry Pi, at home
	- if we have 4 Raspberry PI and need highly available k8s


## How many nodes per clusters and how many clusters should we have

### how many nodes per cluster?
- no particular constraint
	- no need to have an odd number of nodes
- can have clusters with zero node, but it won't be able to start any pods/containers
- for testing and development, having a single node is fine
- for production, make sure to have extra capacity
- k8s is officially tested with upto 5000 notes
	- it's a purely artificial threshold
	- not a hard limit, we actually have instances with more than that, pretty rare
	- however running a cluster of that size requires a lot of tuning
	- this is similar to talking about the speed record for train
		- highest speeds of trains are just for records or stats, it doesn't mean that while common public travels in such high speed trains, that train would be traveling in that much speed
		- rather fastest it goes in "production" (with general public), would be way lesser than the records
		- the highest speeds are just like when one day the engineers decide on "let's break some speed records!!", when they find the right rail track to tune a train with just one passenger car and the engine, put special breaks etc, clocking that speed record

### how many clusters?
- if we have more clusters with less nodes in each, we might want to re-think as it might going to be very expensive
- if we have more nodes in few clusters, we might want to have more smaller clusters to reduce the **blast radius** ( #blast-radius )
	- number of applications or customers etc that are going to be affected if there's an outage = **blast radius**
	- if we pull all our eggs in the same basket, and if the basket has a problem, all the eggs in it would face the problem