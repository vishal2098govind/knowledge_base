#k8s

- k8s is a container management system
- It runs and manages containerized applications on a cluster

## Basic things we can ask k8s to do:
- Start 5 containers using image `atseashop/api:v1.3`
- Place an internal load balancer in front of these containers
- Start 10 containers using `atseashop/webfront:v1.3`
- Place a public load balancer in front of these containers
- It's Black Friday (or Christmas), traffic spikes, grow our cluster and add containers
- New release! Replace my containers with the new image `atseashop/webfront:v1.4`
- Keep processing requests during upgrade, update my containers one at a time

## Advanced things we can ask k8s to do:
- Autoscaling
	- instead of we manually decide count of containers, ask k8s to make sure the **average CPU usage** in the containers is below a threshold
	- k8s **automatically scales** up containers on **spikes** of traffic
	- also on spike goes down, k8s automatically scales down
	- straightforward to autoscale on CPU
	- more complex to autoscale on other metrics
- Resource management and scheduling
	- reserve CPU/RAM for containers
		- k8s also plays the tetris game
			- the pieces are basically our containers
			- k8s is trying to fit our containers on our cluster as efficiently as possible
			- if we have a container needing a lot of CPU, a container needing a lot of RAM, it tries to put those containers together to maximise resource utilization on our cluster
	- placement constraints
		- can tell k8s about where to place our containers
- Advanced rollout patterns
	- Blue/green deployment
		- while we want to replace with a new version, instead of doing one by one, 
		- existing stack is the blue stack
		- create green stack with new version
		- after testing green stack, switch traffic from blue stack to green stack
		- adv:
			- on realising any mess after release, instead panicking or rolling back, switch it back to blue stack
			- in real world switching is done by re-configuring load balancers, where k8s helps
	- Canary deployment