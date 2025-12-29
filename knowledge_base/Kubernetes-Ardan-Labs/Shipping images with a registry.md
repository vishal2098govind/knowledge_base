#k8s
![[Pasted image 20251228205029.png]]
To be able to create deployment and exposing `dockercoins`  with services, we would need to do
```sh
$ kubectl create deployment --image ???
```
- what would we give the image name? 
- when running on a k8s cluster we have multiple nodes and thus we must ship our images
- we ship images using registry
	- build images
	- push images to registry
	- then specify those images in the `kubectl create` commands
- we can use pre-built and pre-shipped images for dockercoins
	- `hasher` → `dockercoins/hasher:v0.1`
	- `redis` → `redis`
	- `rng` → `dockercoins/rng:v0.1`
	- `webui` → `dockercoins/webui:v0.1`
	- `worker` → `dockercoins/worker:v0.1`
- All services should be internal services, except the web UI
	- since we want to be able to connect to the web UI from outside
```sh
# deployments
$ kubectl create deployment --image  dockercoins/worker:v0.1
$ kubectl create deployment --image dockercoins/hasher:v0.1
$ kubectl create deployment --image dockercoins/rng:v0.1
$ kubectl create deployment --image dockercoins/webui:v0.1
$ kubectl create deployment --image redis

# services
$ kubectl expose deployment hasher --port 80
$ kubectl expose deployment rng --port 80
$ kubectl expose deployment redis --port 6379
$ kubectl expose deployment webui --port 80 --type LoadBalancer

# all
$ kubectl get all
NAME                          READY   STATUS    RESTARTS   AGE
pod/hasher-99bbd4bb-c5q7j     1/1     Running   0          54m
pod/redis-7b47f84cc4-w9fcl    1/1     Running   0          53m
pod/rng-65d885d498-sqjrp      1/1     Running   0          60m
pod/webui-74bb6bbc59-m5n8m    1/1     Running   0          60m
pod/worker-5c6f84b477-zmlcf   1/1     Running   0          101m

NAME                 TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
service/hasher       ClusterIP      10.98.92.208     <none>        80/TCP         50m
service/kubernetes   ClusterIP      10.96.0.1        <none>        443/TCP        9h
service/redis        ClusterIP      10.104.206.191   <none>        6379/TCP       50m
service/rng          ClusterIP      10.99.100.55     <none>        80/TCP         49m
service/webui        LoadBalancer   10.96.228.206    localhost     80:32274/TCP   47m

NAME                     READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/hasher   1/1     1            1           54m
deployment.apps/redis    1/1     1            1           53m
deployment.apps/rng      1/1     1            1           60m
deployment.apps/webui    1/1     1            1           60m
deployment.apps/worker   1/1     1            1           101m

NAME                                DESIRED   CURRENT   READY   AGE
replicaset.apps/hasher-99bbd4bb     1         1         1       54m
replicaset.apps/redis-7b47f84cc4    1         1         1       53m
replicaset.apps/rng-65d885d498      1         1         1       60m
replicaset.apps/webui-74bb6bbc59    1         1         1       60m
replicaset.apps/worker-5c6f84b477   1         1         1       101m
```