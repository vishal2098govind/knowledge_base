#k8s #helm
- `helm` is kind of a package manager for k8s
- thanks to helm, we can install packages called `helm charts`, which can be used in two workflows
	- install `prometheus` or `grafana` or use `ingress` and use `helm` to manage those installations - package manager use case
	- also can be used manage the lifecycle of the applications we are working on
		- installing and upgrading the applications and so on

## why do we need help charts, and can't just use YAML files
- as long as we have simple applications, YAML manifests are fine
- if we need to configure things in the application, maybe
	- if we want to set the number of replicas
	- if we have particular rules about scaling
	- if there's a cluster where we need to have an odd number of replicas
	- install a web application and generate an ingress resource and we want to put the domain name in that ingress resource
		- tell users 
			- to download yaml file, 
			- make some changes manually
			- then apply `kubectl` to k8s
		- or give the users a helm chart which is going to expose parameters which people can set when deploying application

## beyond yaml
- instead of having `app-prod.yaml`, `app-staging.yaml`, `app-dev.yaml`
- install helm
```sh
$ helm install app ... --set this.paramter=that.value
# now, this paramter can be used by the yaml, and we can further apply yaml files to k8s
```
- helm is not the only tool that can do this
	- there's `kustomize`
	- there's `ytt`
	- etc
- helm gives pretty powerful tools to generate the YAML
- helm keeps track of what we have installed and comes with a really big library of helm charts

## charts vs packages
- package contains binaries, libraries, etc
	- they are pretty big usually
- chart is just a bunch of yaml files
- `gitlab` has really big and complex chart
- the chart itself doesn't contain the programs or libraries etc
- the charts contain several YAML files
	- the YAML files contain references to the images
- normally we install a package only once
- each installation of a chart is called a release
```sh
$ kubectl create namespace helmdemo
namespace/helmdemo created
$ kubectl config set-context --namespace helmdemo --current
Context "docker-desktop" modified.
$ kubectl config get-contexts
CURRENT   NAME             CLUSTER          AUTHINFO         NAMESPACE
*         docker-desktop   docker-desktop   docker-desktop   helmdemo
$ helm install my-juice-shop securecodebox/juice-shop --version 5.4.0
NAME: my-juice-shop
LAST DEPLOYED: Sun Jan 11 12:10:07 2026
NAMESPACE: helmdemo
STATUS: deployed
REVISION: 1
DESCRIPTION: Install complete
NOTES:
1. Get the application URL by running these commands:
  echo "Visit http://127.0.0.1:3000 to use your application"
  kubectl --namespace helmdemo port-forward service/my-juice-shop 3000:3000
$ kubectl get all
NAME                                 READY   STATUS    RESTARTS   AGE
pod/my-juice-shop-5485c85d44-9dg2q   0/1     Running   0          12s

NAME                    TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
service/my-juice-shop   ClusterIP   10.108.82.165   <none>        3000/TCP   12s

NAME                            READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/my-juice-shop   0/1     1            0           12s

NAME                                       DESIRED   CURRENT   READY   AGE
replicaset.apps/my-juice-shop-5485c85d44   1         1         0       12s

$ kubectl get all
NAME                                 READY   STATUS    RESTARTS   AGE
pod/my-juice-shop-5485c85d44-9dg2q   1/1     Running   0          2m18s

NAME                    TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
service/my-juice-shop   ClusterIP   10.108.82.165   <none>        3000/TCP   2m18s

NAME                            READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/my-juice-shop   1/1     1            1           2m18s

NAME                                       DESIRED   CURRENT   READY   AGE
replicaset.apps/my-juice-shop-5485c85d44   1         1         1       2m18s

$ kubectl port-forward service/my-juice-shop 3000:3000

$ curl localhost:3000
```
- to change the service type to be NodePort from ClusterIP, we can either choose to use `kubectl` or choose to install another specific version of the helm chart specifying this
```sh
$ helm install public-juice-shop securecodebox/juice-shop --version 5.4.0 --set service.type=NodePort
NAME: public-juice-shop
LAST DEPLOYED: Sun Jan 11 12:16:40 2026
NAMESPACE: helmdemo
STATUS: deployed
REVISION: 1
DESCRIPTION: Install complete
NOTES:
1. Get the application URL by running these commands:
  export NODE_PORT=$(kubectl get --namespace helmdemo -o jsonpath="{.spec.ports[0].nodePort}" services public-juice-shop)
  export NODE_IP=$(kubectl get nodes --namespace helmdemo -o jsonpath="{.items[0].status.addresses[0].address}")
  echo http://$NODE_IP:$NODE_PORT
$ kubectl get services
NAME                TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
my-juice-shop       ClusterIP   10.108.82.165   <none>        3000/TCP         6m38s
public-juice-shop   NodePort    10.107.69.124   <none>        3000:30114/TCP   5s
```
- now, to instead use service.type of LoadBalancer, we can also re-configure the installed helm chart using `helm upgrade`
```sh
$ helm upgrade public-juice-shop securecodebox/juice-shop --version 5.4.0 --set service.type=LoadBalancer
Release "public-juice-shop" has been upgraded. Happy Helming!
NAME: public-juice-shop
LAST DEPLOYED: Sun Jan 11 12:19:26 2026
NAMESPACE: helmdemo
STATUS: deployed
REVISION: 2
DESCRIPTION: Upgrade complete
NOTES:
1. Get the application URL by running these commands:
     NOTE: It may take a few minutes for the LoadBalancer IP to be available.
           You can watch the status of by running 'kubectl get --namespace helmdemo svc -w public-juice-shop'
  export SERVICE_IP=$(kubectl get svc --namespace helmdemo public-juice-shop --template "{{ range (index .status.loadBalancer.ingress 0) }}{{.}}{{ end }}")
  echo http://$SERVICE_IP:3000
$ kubectl get services
NAME                TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
my-juice-shop       ClusterIP      10.108.82.165   <none>        3000/TCP         9m27s
public-juice-shop   LoadBalancer   10.107.69.124   localhost     3000:30114/TCP   2m54s
```
- we can scale as well. we can use `--reuse-values` to not loose the load-balancer back to cluster-ip
```sh
$ helm upgrade public-juice-shop securecodebox/juice-shop --version 5.4.0 --set replicaCount=3 --reuse-values
Release "public-juice-shop" has been upgraded. Happy Helming!
NAME: public-juice-shop
LAST DEPLOYED: Sun Jan 11 12:22:03 2026
NAMESPACE: helmdemo
STATUS: deployed
REVISION: 4
DESCRIPTION: Upgrade complete
NOTES:
1. Get the application URL by running these commands:
  echo "Visit http://127.0.0.1:3000 to use your application"
  kubectl --namespace helmdemo port-forward service/public-juice-shop 3000:3000
$ kubectl get all
NAME                                    READY   STATUS    RESTARTS   AGE
pod/my-juice-shop-5485c85d44-9dg2q      1/1     Running   0          12m
pod/public-juice-shop-774ff5449-grnxx   0/1     Pending   0          6s
pod/public-juice-shop-774ff5449-p79gq   0/1     Pending   0          6s
pod/public-juice-shop-774ff5449-zctj5   1/1     Running   0          5m28s

NAME                        TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
service/my-juice-shop       ClusterIP   10.108.82.165   <none>        3000/TCP   12m
service/public-juice-shop   ClusterIP   10.107.69.124   <none>        3000/TCP   5m29s

NAME                                READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/my-juice-shop       1/1     1            1           12m
deployment.apps/public-juice-shop   1/3     3            1           5m29s

NAME                                          DESIRED   CURRENT   READY   AGE
replicaset.apps/my-juice-shop-5485c85d44      1         1         1       12m
replicaset.apps/public-juice-shop-774ff5449   3         3         1       5m29s
```
- using YAML file
```sh
$ code juice.yaml
$ cat juice.yaml
replicaCount: 5

service:
  type: LoadBalancer
  port: 80

$ helm upgrade public-juice-shop securecodebox/juice-shop --version 5.4.0 --values juice.yaml

Release "public-juice-shop" has been upgraded. Happy Helming!
NAME: public-juice-shop
LAST DEPLOYED: Sun Jan 11 12:28:49 2026
NAMESPACE: helmdemo
STATUS: deployed
REVISION: 7
DESCRIPTION: Upgrade complete
NOTES:
1. Get the application URL by running these commands:
     NOTE: It may take a few minutes for the LoadBalancer IP to be available.
           You can watch the status of by running 'kubectl get --namespace helmdemo svc -w public-juice-shop'
  export SERVICE_IP=$(kubectl get svc --namespace helmdemo public-juice-shop --template "{{ range (index .status.loadBalancer.ingress 0) }}{{.}}{{ end }}")
  echo http://$SERVICE_IP:80

$ kubectl get all
NAME                                    READY   STATUS    RESTARTS   AGE
pod/my-juice-shop-5485c85d44-9dg2q      1/1     Running   0          18m
pod/public-juice-shop-774ff5449-7tf2w   0/1     Pending   0          3s
pod/public-juice-shop-774ff5449-ksc95   0/1     Pending   0          3s
pod/public-juice-shop-774ff5449-zctj5   1/1     Running   0          12m
pod/public-juice-shop-774ff5449-zdcpb   1/1     Running   0          4m5s
pod/public-juice-shop-774ff5449-zw5kp   1/1     Running   0          4m5s

NAME                        TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
service/my-juice-shop       ClusterIP      10.108.82.165   <none>        3000/TCP       18m
service/public-juice-shop   LoadBalancer   10.107.69.124   localhost     80:31426/TCP   12m

NAME                                READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/my-juice-shop       1/1     1            1           18m
deployment.apps/public-juice-shop   3/5     5            3           12m

NAME                                          DESIRED   CURRENT   READY   AGE
replicaset.apps/my-juice-shop-5485c85d44      1         1         1       18m
replicaset.apps/public-juice-shop-774ff5449   5         5         3       12m
```
- helm list
```sh
$ helm list
NAME             	NAMESPACE	REVISION	UPDATED                             	STATUS  	CHART           	APP VERSION
my-juice-shop    	helmdemo 	1       	2026-01-11 12:10:07.50271 +0530 IST 	deployed	juice-shop-5.4.0	v19.1.1
public-juice-shop	helmdemo 	7       	2026-01-11 12:28:49.975556 +0530 IST	deployed	juice-shop-5.4.0	v19.1.1
```

## Installing Ingress Controller using Helm
```sh
$ helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
"ingress-nginx" has been added to your repositories
$ helm install ingress-nginx-demo ingress-nginx/ingress-nginx --version 4.14.1
$ kubectl create ingress red --rule=/red=red:80
```

### classes and Ingress classes
#k8s-classes
```sh
$ kubectl get ingressclasses
NAME      CONTROLLER                      PARAMETERS   AGE
nginx     k8s.io/ingress-nginx            <none>       25m
traefik   traefik.io/ingress-controller   <none>       4d17h
$ kubectl get ingressclass
NAME      CONTROLLER                      PARAMETERS   AGE
nginx     k8s.io/ingress-nginx            <none>       25m
traefik   traefik.io/ingress-controller   <none>       4d17h
```
- the classes concept can be found in multiple places in k8s, when we have the choice between multiple systems or implementations, we're going to see a `blahblahblahclass`
- since potentially we can have multiple ingress controllers, we are going to see `ingressclass`
- similarly we have `storageclasses`, `priorityclasses`, `runtimeclasses`
- while working with a `traefik` ingress controller, if we don't put a class while creating ingress resources, it would still work
- for `nginx` ingress controller, we must indicate class while creating ingress resources
```sh
$ kubectl delete ingress red
ingress.networking.k8s.io "red" deleted
$ kubectl create ingress red --class nginx --rule /red=red:80
$  kubectl get all
NAME                                                 READY   STATUS    RESTARTS   AGE
pod/blue-5c986bd7bf-zdb9r                            1/1     Running   0          22m
pod/ingress-nginx-demo-controller-75f4856bd7-fnmrb   1/1     Running   0          34m
pod/red-77f6d65f98-qmxts                             1/1     Running   0          22m

NAME                                              TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
service/blue                                      ClusterIP      10.102.126.22    <none>        80/TCP                       16m
service/ingress-nginx-demo-controller             LoadBalancer   10.96.23.225     localhost     80:30205/TCP,443:32195/TCP   34m
service/ingress-nginx-demo-controller-admission   ClusterIP      10.100.161.127   <none>        443/TCP                      34m
service/red                                       ClusterIP      10.110.230.21    <none>        80/TCP                       17m

NAME                                            READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/blue                            1/1     1            1           22m
deployment.apps/ingress-nginx-demo-controller   1/1     1            1           34m
deployment.apps/red                             1/1     1            1           22m

NAME                                                       DESIRED   CURRENT   READY   AGE
replicaset.apps/blue-5c986bd7bf                            1         1         1       22m
replicaset.apps/ingress-nginx-demo-controller-75f4856bd7   1         1         1       34m
replicaset.apps/red-77f6d65f98                             1         1         1       22m
$ kubectl port-forward service/ingress-nginx-demo-controller 808
0:80
Forwarding from 127.0.0.1:8080 -> 80
Forwarding from [::1]:8080 -> 80
Handling connection for 8080
Handling connection for 8080

...
$ curl localhost:8080
<html>
<head><title>404 Not Found</title></head>
<body>
<center><h1>404 Not Found</h1></center>
<hr><center>nginx</center>
</body>
</html>
➜  ~ curl localhost:8080/red
🔴This is pod nginx-helm-demo/red-77f6d65f98-qmxts on linux/amd64, serving /red for 10.1.1.232:59874.
```