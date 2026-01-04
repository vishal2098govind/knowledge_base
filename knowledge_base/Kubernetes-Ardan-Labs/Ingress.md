#k8s

- When we expose a service of type `LoadBalancer`, that's how we have a way to connect outside a cluster

![[Pasted image 20260104175606.png]]
- when we expose a load balancer, we find the external-ip address to stay in `<pending>` for most of the time
```sh
$ kubectl expose deployment red --port 80 --type LoadBalancer
service/red exposed
➜  container.training git:(main) ✗ kubectl get all --selector 'app in (blue, red)'
NAME                        READY   STATUS    RESTARTS   AGE
pod/blue-5c986bd7bf-p8mk2   1/1     Running   0          32m
pod/red-77f6d65f98-b4pgz    1/1     Running   0          12m

NAME           TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
service/blue   LoadBalancer   10.97.224.215    localhost     80:32732/TCP   30m
service/red    LoadBalancer   10.104.152.107   <pending>     80:30912/TCP   11s

NAME                   READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/blue   1/1     1            1           32m
deployment.apps/red    1/1     1            1           12m

NAME                              DESIRED   CURRENT   READY   AGE
replicaset.apps/blue-5c986bd7bf   1         1         1       32m
replicaset.apps/red-77f6d65f98    1         1         1       12m
```
- here, we cannot connect to the red container, but to the blue container from our localhost:80
- and can connect to the red container via the node port 30912
- usually we don't give that big container to a website url while sharing to someone, rather we just give the hostname as the url and the browser assumes it to be on port 80 by default for HTTP and port 443 for HTTPS
- thus, on these clusters, we need to solve an extra challenge to solve to expose all the containers that we want, on port 80
- even on a managed k8s cluster, we have a load-balancer services, and each of these LB services are going to cost money
- that's where the Ingress controllers are going to help
- Instead of having one LoadBalancer or one IP address per service, we can use a single load-balancer or single IP address for all our HTTP services
![[Pasted image 20260104203124.png]]
- Here, we only have one load-balancer service per cluster and can route requests to pods using that one load balancer. This saves IP addresses to be bought from one per service to one per cluster
- thus, we can expose any number of services using just one load balancer service
- The idea of ingress controller is
	- saving money is one aspect
	- another aspect is the ability to do content based routing
		- it means we can forward the requests to different pods depending on what the request is
		- i.e. we can decide to send requests for 
			- `/api` to a particular pod 
			- and `/static` to other pod
			- and `/graphql` to another pod
		- or, we can decide to send
			- `/api/v1` to a particular set of pods
			- `/api/v2` to another particular set of pods
	- that way, we have multiple services behind a single domain and then use different URIs or HTTP paths to host our services
		- this is convenient for SPAs (single page applications)
	- if we have different domains, we have to setup CORS (cross origin resource sharing) for allowing requests from one domain to the other domain