#k8s

### Exposing HTTP services with Ingress resources
- Service = layer 4 (TCP, UDP, SCTP)
    - works with every TCP/UDP/SCTP protocol
    - doesn't "see" or interpret HTTP
- Ingress = layer 7 (HTTP)
    - only for HTTP
    - can route requests depending on URI or host header
    - can handle TLS


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


## Steps
### Step-1: Deploy an *ingress* controller
- one time setup per cluster
- this is because there are very few k8s clusters that ship with an ingress controller installed by default, except GKE (google cloud)
- why not shipped by default?
	- there are many ingress controllers available out there
	- none of them can be really be the default
	- similar to CNI plugins
#### `Traefik` ingress controller
- it's a cloud native load balancer
- what makes `Traefik` a cloud native LB
	- one of the first LB that primarily uses APIs to configure itself
	- for most LB, when we install it, we need to write a configuration file like
		- NGINX
		- HAProxy
		- Apache, etc
	- in which we define frontends and backends and start the LB, and when want to make changes, we update the configuration file and then we restart the LB or send a special to LB to reload the updated configuration
- `Traefik` is different in the sense that it's going to get it's configuration information using APIs
	- it can use the docker APIs to see available docker containers
	- it can use the k8s API to see available ingress resources
- it doesn't need or have to have configuration file anymore, although we can
- very often, `traefik` is going to be set up with just a few command line flags, where it's going to obtain all the configuration information **dynamically**, through an API
	- dynamically here means, we need not reload the configuration any more
	- it subscribes to changes in the API
- thus, whenever we create ingress resources or scale deployments up or down, `traefik` is going to immediately pickup these changes and re-configure itself
- 
### Step-2: Create *Ingress Resources*
- maps a domain or path to k8s service
- the ingress controller watches ingress resources and sets up a LB
### Step-3: set up DNS
- associate DNS entries with LB address

#### `nip.io` A magic domain
- http://nip.io
```sh
$ ping whateveryouwant.1.1.1.1.nip.io
PING whateveryouwant.1.1.1.1.nip.io (1.1.1.1): 56 data bytes
64 bytes from 1.1.1.1: icmp_seq=0 ttl=53 time=15.377 ms
64 bytes from 1.1.1.1: icmp_seq=1 ttl=53 time=12.292 ms
64 bytes from 1.1.1.1: icmp_seq=2 ttl=53 time=19.241 ms
64 bytes from 1.1.1.1: icmp_seq=3 ttl=53 time=20.681 ms
64 bytes from 1.1.1.1: icmp_seq=4 ttl=53 time=20.420 ms
64 bytes from 1.1.1.1: icmp_seq=5 ttl=53 time=21.254 ms
64 bytes from 1.1.1.1: icmp_seq=6 ttl=53 time=12.423 ms
^C
--- whateveryouwant.1.1.1.1.nip.io ping statistics ---
7 packets transmitted, 7 packets received, 0.0% packet loss
round-trip min/avg/max/stddev = 12.292/17.384/21.254/3.647 ms
```
- we can use 1.1.1.1nip.io as our domain 
- we pretend that we bought a domain name like "`cloudnative.party`" and we did set it up to point to the IP address of one node of the cluster
- to make our lives easier of not buying domain names and wait for DNS propagation etc, we can use nip.io
- we want to be able to 
```sh
$ curl 192.168.29.91.nip.io
curl: (7) Failed to connect to 192.168.29.91.nip.io port 80 after 3 ms: Couldn't connect to server
```


## Setting up `Traefix`
### Step-1: Create Ingress Controller
- We are going to use a daemon-set with a host network
- using host network means equivalent to
```sh
$ docker run --net host ...
```
- i.e. if we are listening on port 80 in that pod, then we are listening on port 80 in that host as well
```sh
$ kubectl -f apply k8s/traefik.yaml
$ kubectl -n traefik port-forward pod/traefik-sjwkf 8080:8080
Forwarding from 127.0.0.1:8080 -> 8080
Forwarding from [::1]:8080 -> 8080
Handling connection for 8080
Handling connection for 8080
Handling connection for 8080
Handling connection for 8080
Handling connection for 8080
Handling connection for 8080
Handling connection for 8080
Handling connection for 8080
Handling connection for 8080
Handling connection for 8080

---
$ curl localhost:8080/dashboard
<!DOCTYPE html>
<html lang="en">
  <head>


    <script>
      window.APIUrl = "/api/"
    </script>


    <meta charset="utf-8" />
    <link rel="icon" href="./favicon.ico" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta name="theme-color" content="#ddee6d" />
    <meta
      name="description"
      content="Traefik Proxy"
    />
    <!--
      manifest.json provides metadata used when your web app is installed on a
      user's mobile device or desktop. See https://developers.google.com/web/fundamentals/web-app-manifest/
    -->
    <link rel="manifest" href="./manifest.json" />
    <title>Traefik Proxy</title>
    <script type="module" crossorigin src="./assets/index-BNX0Ip2T.js"></script>
  </head>
  <body>
    <noscript>You need to enable JavaScript to run this app.</noscript>
    <div id="root"></div>
  </body>
</html>
```
### Next Step - Create Ingress resources
- we want to exposing `webui` deployment of the `dockercoins`
- for that we need to create an ingress
```sh
$ kubectl create ingress dockercoins --rule=dockercoins.127.0.0.1.nip.io/*=webui:80

$ kubectl port-forward ds/traefik 8081:80 # needed while running from docker-desktop on a mac

$ curl dockercoins.127.0.0.1.nip:8081
# redirects to webui index.html

$ kubectl get ingress dockercoins -o yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  creationTimestamp: "2026-01-06T17:37:10Z"
  generation: 7
  name: dockercoins
  namespace: traefik
  resourceVersion: "224667"
  uid: fcc59279-9ef5-48b2-837a-7f4c52df0712
spec:
  ingressClassName: traefik
  rules:
  - host: dockercoins.127.0.0.1.nip.io
    http:
      paths:
      - backend:
          service:
            name: webui
            port:
              number: 80
        path: /
        pathType: Prefix
status:
  loadBalancer: {}
```

#### Rainbow
```sh
$ kubectl create deployment blue --image jpetazzo/color
deployment.apps/blue created
$ kubectl create deployment red --image jpetazzo/color
deployment.apps/red created
$ kubectl create ingress rainbow \
> --rule=rainbow.127.0.0.1.nip.io/blue=blue:80 \
> --rule=rainbow.127.0.0.1.nip.io/red=red:80
ingress.networking.k8s.io/rainbow created
$ kubectl get ingress rainbow -o yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  creationTimestamp: "2026-01-06T18:57:42Z"
  generation: 1
  name: rainbow
  namespace: traefik
  resourceVersion: "225927"
  uid: 82febccb-8ced-4cc8-8418-1c0d5181133f
spec:
  ingressClassName: traefik
  rules:
  - host: rainbow.127.0.0.1.nip.io
    http:
      paths:
      - backend:
          service:
            name: blue
            port:
              number: 80
        path: /blue
        pathType: Exact
      - backend:
          service:
            name: red
            port:
              number: 80
        path: /red
        pathType: Exact
status:
  loadBalancer: {}
$ kubectl expose deployment red --port 80
service/red exposed
$ kubectl expose deployment blue --port 80
service/blue exposed
$  curl http://rainbow.127.0.0.1.nip.io:8081/blue
🔵This is pod traefik/blue-5c986bd7bf-w8wwc on linux/amd64, serving /blue for 10.1.0.1:65316.
$ curl http://rainbow.127.0.0.1.nip.io:8081/red
🔴This is pod traefik/red-77f6d65f98-rd8nq on linux/amd64, serving /red for 10.1.0.1:56930.
$ curl http://rainbow.127.0.0.1.nip.io:8081/green
404 page not found
```
![[Pasted image 20260107004326.png]]

## Single Ingress controller per cluster VS one per namespace
- If we want simplicity - single ingress controller per cluster
	- this has some security implications
	- global ingress controller we need a little bit too much permission, which we don't want and rather want to really isolate namespaces from each other
- If we want to be picky about the ingress controller they want to use for different projects
	- we want to use this particular ingress controller because it has support for 
		- A/B testing
		- Canary deployment
		- very advance like tracing etc
		- so we don't want the default global ingress controller in our project's namespace

## Ingress standard features
- Load Balancing
- SSL termination
- URI Routing
	- `/api` -> api-service
	- /static -> assets-service

## Ingress extended features
More advanced features that are supported by most ingress controllers, but are not part of the spec, and are done by vendor specific extensions
- Routing with other headers or cookies
- A/B testing
- Canary deployment
	- e.g. send 1% of the traffic to a specific different version or backend than the most of the traffic
- etc