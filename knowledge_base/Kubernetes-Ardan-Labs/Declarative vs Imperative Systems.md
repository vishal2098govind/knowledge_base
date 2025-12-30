#k8s

## Imperative vs Declarative
- The key thing here is to understand this idea that we are going to 
	- define resources with YAML manifests, 
	- and then we can version that YAML, keep successive commits
	- instead of we connecting to the cluster always and doing `kubectl scale` or `kubectl delete`, we would rather 
		- update the YAML manifests, 
		- commit
		- submit pull request, etc
	- and even if we don't go through a pull request, the fact that we are committing our changes to our manifests and the fact that we can revert back to our earlier manifests and the fact that we are able to maintain a history of the changes made in the manifests
	- so, instead of arriving to the cluster and realizing that a worker deployment just disappeared from the cluster, why? 
		- to know why, we can look at the YAML manifest and it's git log and know 
			- what happened to the worker deployment
			- who removed that worker deployment
			- why did they remove the worker deployment - can be known if there exists a detailed PR comment
	- and we could easily roll back to the point where we had that worker deployment and `kubectl apply -f` and we have our worker deployment back
```sh
$ code k8s/dockercoins.yaml
# include the worker deployment back
$ kubectl apply -f k8s/dockercoins.yaml
deployment.apps/hasher unchanged
service/hasher unchanged
deployment.apps/redis unchanged
service/redis unchanged
deployment.apps/rng unchanged
service/rng unchanged
deployment.apps/webui unchanged
deployment.apps/worker created
```
- k8s has a declarative API
	- Declarative: `I would like a cup of tea.` 
		- write a manifest asking for a tea
		- put that manifest in `etcd`
		- the controllers in the control plane would wake up and find the manifest in the `etcd`
		- the controllers start working on fulfilling the manifest requirements
	- Imperative:
		- here, we give a list of instructions
			- Boil some water.
			- Pour it in a tea pot.
			- Add tea leaves.
			- Sleep for a while.
			- Serve in a cup.
- At first Declarative seems more simpler, as long as we know how to brew tea
	- What declarative would really be:
	    _I want a cup of tea, obtained by pouring an infusion¹ of tea leaves in a cup._
	    _¹An infusion is obtained by letting the object steep a few minutes in hot² water._
	    _²Hot liquid is obtained by pouring it in an appropriate container³ and setting it on a stove._
	    _³Ah, finally, containers! Something we know about. Let's get to work, shall we?_
- so, the declarative instructions are recursive and thus have mystery surprise while following along
- imperative instructions are specified in a specific order, without any recursion or mystery surprises
### Which one is better? declarative or imperative? and what does it have to do with deploying containers
- imperative model is like a shell script or a code/program
	- where we write all the instructions that need to happen in the correct order that need to happen
- declarative model is when we indicate what we want
	- I want a replica set with 3 alpine pods
	- and then there's something working behind the scenes 
		- to realize that, 
		- to materialize that
### Imperative System
- imperative systems are simpler, but
	- if we have an interruption, we need to restart everything all over
	- e.g. brewing tea
		- if while boiling water, if we get a phone call and we talk for a while and by the time we are done talking, and go back to the tea operation, multiple things might have happened by that time
			- **best case scenario**: maybe it was as a short call, so the water is still warm
			- if, we stayed on the call for long enough, 
				- but we have a really fancy electric kettle that keeps the water warn, no matter how long, then also everything is fine
				- or, if we do not have a fancy electric kettle, just a normal electric kettle, and the water got warm and got cold again, we need to boil the water again
				- if we rather were boiling water on a stove top, instead of an electric kettle, all of the water boiled and then completely evaporated, we need to put more water again and boil
				- or maybe the water boiled completely and ruined the utensil in which we were boiling the water, and we now need to use another utensil and boil water again
	- so, in an imperative system, when there is an interruption, we have two options: 
		- just discard everything so far and start again
		- or add a lot of error checking code on each interruption
			- if we are interrupted during this particular phase, check this, check that, and if all looks good, proceed, or else start again
### Declarative System
- the declarative system is going to be more resilient in 
	- handling errors
	- interruptions
	- etc, 
- in a much nicer way
- if we are interrupted in a declarative system, we can figure out what's missing and do only what's necessary
- we need to be able to observe the system
- and compute a diff between 
	- what we have so far
	- and what we want
### Declarative vs Imperative in k8s
- at first k8s looks like an imperative system, looking at the `kubectl run`, `kubectl create`, `kubectl expose` etc
- these commands look imperative to the user
- but behind the scenes, we have a declarative API
- each time we want something to happen, 
	- we don't give instructions directly to a component
	- rather, write a manifest
	- put the manifest in `etcd`
	- a bunch of controllers observing what's going on


## Exact sequence of tasks that happen when we create a deployment
```sh
kubectl create deployment web --image nginx
```

1. Consider this cluster on which we run the command
![[Pasted image 20251230220411.png]]
2. The API server receives the request and is going to create the deployment
	- This means, it's going to put the deployment manifest in `etcd`
![[Pasted image 20251230220645.png]]
3. Then the API server will tell us "deployment created"
	- this doesn't yet mean that the containers are running already by now
	- it just means, yes we've put the deployment in `etcd` 
	- ~ to how when we go to a restaurant and order food, they are going to tell us, 
		- yes we've got it
		- it doesn't mean that your meal is ready yet
		- it means, the kitchen would be working on it soon
![[Pasted image 20251230221458.png]]
4. Now, the kitchen will start working on the requested meal
	- Since the meal asked for a deployment, the deployment controller that's going to wake up
		- In the controller manager, we don't have just one or two controllers
		- There are dozens of controllers, grouped together
		- each of the controller has the responsibility for a specific type of object
	- here, the deployment controller wakes up and it sees the deployment manifest in the `etcd`
![[Pasted image 20251230221921.png]]
5. The deployment controller then understands the need for a replica set
	- so the deployment controller creates a replica set manifest and puts it in `etcd`
	- then it goes back to sleep
![[Pasted image 20251230222056.png]]
![[Pasted image 20251230222105.png]]
6. Then, the replica set controller is going to wake up
![[Pasted image 20251230222148.png]]
7. The replica set controller finds a new replica set manifest in the `etcd`
	- a replica set means that we need a number of pods for the replica set
	- here, we need only one `nginx` pod as per the request
	- and we currently have 0 of them
	- so, the replica set controller understands the need to create a `nginx` pod
	- so the replica set controller creates a `nginx` pod manifest and puts it in the `etcd`
	- then the replica set controller goes back to sleep
![[Pasted image 20251230225054.png]]
![[Pasted image 20251230232516.png]]
8. Then, the scheduler wakes up on seeing a pending pod, since the responsibility of a scheduler is to assign pods to node
	- a scheduler constantly looks for pending pod manifests in the `etcd`
	- on seeing a pending pod manifest, it tries to understand any constraints
		- how many CPUs the pod needs
		- how much RAM the pod needs
		- any placement constraints for the pod
			- does the pod needs to be placed in that data center in that region
			- or whatever
	- thus, the scheduler makes a decision of which node to choose, based on those available constraints, if any
	- say, here scheduler decides for node-1
	- and then the scheduler goes back to sleep
![[Pasted image 20251230232526.png]]
![[Pasted image 20251230232535.png]]
![[Pasted image 20251230234023.png]]
9. at that point, `kubelet` on node-1, notices that pod, because the responsibility of `kubelet` is to register itself with the control plane, and basically say, "Hello, I am not node-1, I have"
	- "that many cores"
	- "that much RAM"
	- "do you have any pods or containers for me?"
- at that point, the control plane will the `kubelet` about this pod
- `kubelet` will update the status of the pod in `etcd` to "creating"
- then, it would 
	- pull the `nginx` image, 
	- allocate the IP address
	- setup the container
- once, the container is running, it would report back to the control plane and update the pod status to "Running" in the `etcd`
![[Pasted image 20251230235641.png]]
![[Pasted image 20251230235702.png]]
![[Pasted image 20251230235736.png]]
![[Pasted image 20251230235752.png]]