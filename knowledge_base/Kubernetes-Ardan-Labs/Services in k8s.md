#k8s
- exposing a deployment creates a k8s **service** acting as a **static IP address** as well as a **load balancer**
```sh
$ kubectl expose <resource> <name> --port [--type]
```
- services can also be created using
```sh
$ kubectl create service <type> <name> [--tcp port]
```

```
## `ClusterIP` service
```sh
$ kubectl create deployment blue --image jpetazzo/color
deployment.apps/blue created
$ kubectl expose deployment blue --port 80
service/blue exposed
$ kubectl get service
NAME         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
blue         ClusterIP   10.96.184.245   <none>        80/TCP    63s
```
- now, if we use this Cluster-IP address, we will be able to survive pods getting down and new pods getting created
- thus, the **service** here acts like a **static ip address**
-  on scaling:
```sh
$ kubectl scale deployment blue --replicas=2
deployment.apps/blue scaled
$ kubectl get pods
NAME                    READY   STATUS    RESTARTS   AGE
blue-5c986bd7bf-bv8q6   1/1     Running   0          7m18s
blue-5c986bd7bf-gz5w2   1/1     Running   0          8s
```
- now, if we use this cluster-ip address to hit our pod or application, we will be able to hit to one of the pods, may not be in a perfectly round-robin fashion, but some here and some there, but on average around equal distribution of load being balanced
- thus, the **service** here acts like a **load balancer** among the pods
- if we do-not want to hard code cluster ip-addresses, we can use **DNS service**
	- we can directly do `curl blue.default.svc`
	- if we are in the same name space, we can just use the service name
![[Pasted image 20251228175222.png]]
- we have something that balances the load to the backend pods
	- that something could be `kube-proxy`, `kube-router`, `Cilium...`
	- by default, it's `kube-proxy` on most of the clusters
![[Pasted image 20251228180231.png]]
- instead of Cluster IP, we can connect using the service name
	- do a DNS resolution using something that could be `CoreDNS` (formerly `kube-dns`)
	- in most clusters it would be `CoreDNS`, but can be replaced as per specific requirements


## `LoadBalancer` service
- for traffic coming from outside the cluster, we cannot use cluster-ip as it is private within the cluster
```sh
$ kubectl get service
NAME         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
blue         ClusterIP   10.96.184.245   <none>        80/TCP    63s
```
- we need to use the EXTERNAL-IP to be able to connect from outside of the cluster
- here, it's \<none>
- thus, we need to get an external ip for the **blue** service
- this is how we will make the difference between internal and external services
- for an internal service, we just do `kubectl expose` and we get a `ClusterIP` service
- for external or public facing service, we need to have have an external IP
- we do that using load balancer
```sh
$ kubectl get deployment
No resources found in default namespace.
$ kubectl create deployment blue --image jpetazzo/color
deployment.apps/blue created
$ kubectl expose deployment blue --port 80
service/blue exposed
$ kubectl get services
NAME         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
blue         ClusterIP   10.99.247.202   <none>        80/TCP    8s
kubernetes   ClusterIP   10.96.0.1       <none>        443/TCP   3h13m
$ kubectl create deployment red --image jpetazzo/color
deployment.apps/red created
➜  dockercoins git:(main) kubectl expose deployment red --port 80 --type --help
The Service "red" is invalid: spec.type: Unsupported value: "--help": supported values: "ClusterIP", "ExternalName", "LoadBalancer", "NodePort"
$ kubectl expose deployment red --port 80 --type LoadBalancer
service/red exposed
$ kubectl get services
NAME         TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
blue         ClusterIP      10.99.247.202   <none>        80/TCP         3m7s
kubernetes   ClusterIP      10.96.0.1       <none>        443/TCP        3h16m
red          LoadBalancer   10.102.183.34   localhost     80:32680/TCP   7s
$ curl localhost
🔴This is pod default/red-77f6d65f98-pxlz4 on linux/amd64, serving / for 192.168.65.3:63048.
$ kubectl scale deployment red --replicas=3
deployment.apps/red scaled
$ kubectl get pods
NAME                    READY   STATUS              RESTARTS   AGE
blue-5c986bd7bf-tgh7s   1/1     Running             0          10m
red-77f6d65f98-454qr    1/1     Running             0          6s
red-77f6d65f98-pxlz4    1/1     Running             0          9m34s
red-77f6d65f98-spv8s    0/1     ContainerCreating   0          6s
$ curl localhost
🔴This is pod default/red-77f6d65f98-spv8s on linux/amd64, serving / for 192.168.65.3:63736.
$ curl localhost
🔴This is pod default/red-77f6d65f98-pxlz4 on linux/amd64, serving / for 192.168.65.3:63750.
$ curl localhost
🔴This is pod default/red-77f6d65f98-spv8s on linux/amd64, serving / for 192.168.65.3:63762.
$ curl localhost
🔴This is pod default/red-77f6d65f98-pxlz4 on linux/amd64, serving / for 192.168.65.3:63778.
$ curl localhost
🔴This is pod default/red-77f6d65f98-454qr on linux/amd64, serving / for 192.168.65.3:63780.
$ curl localhost
🔴This is pod default/red-77f6d65f98-spv8s on linux/amd64, serving / for 192.168.65.3:63784.
```
![[Pasted image 20251228182941.png]]

### Why we have to explicitly mention `--type=LoadBalancer`
consider this cluster:
 ![[Pasted image 20251228201142.png]]
 ![[Pasted image 20251228201214.png]]
- we have 
	- couple of internal services (green and purple)
	- couple of external services (blue and red)
- to connect to the external services from outside of the cluster, we create services of type load balancer
- at first, the EXTERNAL-IP is `<pending>
- then, a **special controller** in the control plane kicks in. it can be any of
	- CCM - cloud controller manager
	- MetalLB - typically for on premise deployment
- that controller notices the pin-board (`etcd`) and notices that we want a service of type `LoadBalancer`
- this controller is going to a cloud load balancer
- this load balancer is exactly at the edge of the cluster (one leg inside and another leg outside of the cluster)
- though, it's extremely rare to have load balancer at the edge of the cluster
 ![[Pasted image 20251228201240.png]]
- generally load balancers are outside of the cluster
- and thus it cannot get traffic inside the cluster
- how to solve this? - `NodePort`
![[Pasted image 20251228202153.png]]
## `NodePort` service
![[Pasted image 20251228202222.png]]
```sh
$ kubectl get svc
NAME         TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
blue         ClusterIP      10.99.247.202   <none>        80/TCP         126m
kubernetes   ClusterIP      10.96.0.1       <none>        443/TCP        5h20m
red          LoadBalancer   10.102.183.34   localhost     80:32680/TCP   123m
```
- for the `LoadBalancer` service, there is an extra port number of 32680
- on creation of the load balancer service, that port will be opened on every node in the cluster - thus called NodePort
![[Pasted image 20251228203021.png]]
![[Pasted image 20251228203440.png]]
- now we can configure the external load balancer to connect to those `NodePort`s
- and the external world can connect to the load balancer on port 80
- the full path - 80:32680
	- users connect to the load balancer on port 80
	- the load balancer connects to the NodePort
	- the NodePort then connects to the individual Pods

## Why we can't just have load balancer at the edge of the cluster (one leg inside and another leg outside the cluster)? or Why we need NodePort?
![[Pasted image 20251228201214.png]]
- There's a pretty good reason why most k8s clusters won't have load balancers at the edge of the cluster (one leg inside and another leg outside)
- The reason is that k8s lets us replace the network components
	- this means that, the load balancer needs to be able to communicate with the pods
	- the internal communication may or may not be possible depending on the network mechanism we have inside the cluster
	- sometimes, we can bring in load balancer but sometimes we may not, depending on what we have in the cluster
	- so to have something kind of universal, something that will work everywhere, regardless of what combination of things we use, we have this concept of **NodePort** so that we can have the load balancer sitting outside of the cluster and the NodePort serves as the interface at the edge or between inside and outside of the cluster