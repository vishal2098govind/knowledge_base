#k8s

- A k8s cluster, nodes and pods, is one big flat IP network
- This is pretty different from docker
- In docker we have multiple networks that are isolated from each other
- In k8s, we have one network and that's it
- All the pods, all the nodes, everything is in the same network and can communicate directly, 
- there is no NAT - network address translation
- there is no port mapping
- we don't have new protocols coming in to encapsulate
- it's primitively simple
- the only implementation detail is that, we cannot decide that the IP addresses of the pods
	- they are set by the **network layer implementation**
	- the k8s doesn't mandate any particular network implementation
## Good news
- we have lots of freedom on how we want to implement the network layer
## Bad news
- network layer has any and many implementations
- for e.g. we are setting up k8s on a cluster of servers or `Rassberrypi` or VMs
	- then, it tells us to set up our control plane
	- then, setup `kubelet`, `tls` etc
	- then it says, now we need to setup networking, and for that we need to pick a CNI Plugin
- the [CNI](https://github.com/containernetworking/cni?tab=readme-ov-file#who-is-using-cni) (Container Network Interface) Plugin have many implementations
- which one to pick? - It depends
	- If we are running in a specific cloud, we might want to use the plugin specific to that cloud
	- If we are on EC2, we might want to pick the CNI plugin specific to EC2, for example
	- But this is not a hard requirement
	- In fact, if we setup a Amazon EKS cluster, it will by default use Amazon ECS CNI Plugin, but we can replace that with another CNI plugin
	- In fact, for a k8s implementation to be compliant, we need to be able to change the network implementation
		- i.e. if a provider wants to announce themselves as a k8s provider, they need to let us change the CNI Plugin if we want to.
		- if we cannot change the CNI plugin while using a k8s provider, then it's not k8s implementation, it's just another container orchestration system

## CNI
- CNI is just responsible for the communication
	- between the pods
	- and between the pods and nodes
- but, we have multiple layers in the k8s networking
![[Pasted image 20251231174619.png]]
### Networking Layers
#### Pod networking layer
- pod-to-pod networking layer
- this part is managed by CNI plugin
- this part is responsible for communication between pods and nodes
- the circle in the diagram is just logical implementation of the network
	- just says it connects all these together
	- it could, under the hood, be using
		- routing
		- bridging
		- encapsulation
	- it is just an implementation detail
	- as long as we've something that gives us IP protocol or Layer-3 (OSI model) communication between pods and nodes, we're good
![[Pasted image 20251231174719.png]]
#### Pod-to-Service networking layer
- responsible for `ClusterIP` and `NodePort`
- on most clusters, this is implemented by `kube-proxy`
![[Pasted image 20260101003834.png]]
#### Network Policies
- namespaces don't give us network isolation
- if we want network isolation, we need network policies
- this layer is independent of others, and can be replaced if need
![[Pasted image 20260101003949.png]]
### All layers seeing together
![[Pasted image 20251231174619.png]]
- the engineers who designed this system are very smart 
- we can replace each of these layers independently
- without touching other layers
- that's a sign of a very good design
- because, even while a cluster is running, we can change the implementation at any layer
	- Example: If we're on AWS and we're using AWS EKS, by default they will be using AWS CNI Plugin in these layers