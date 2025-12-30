#k8s

## `-f` option can be used with many other commands with k8s
```sh
$ kubectl delete -f dockercoins.yaml --namespace dev
$ kubectl label -f dockercoins.yaml release=dev --namespace dev
```

- setting up k8s cluster
```sh
kubectl create namespace dev
namespace/dev created
$ kubectl config set-context --namespace dev --current
Context "docker-desktop" modified.
$ kubectl config get-contexts
CURRENT   NAME             CLUSTER          AUTHINFO         NAMESPACE
*         docker-desktop   docker-desktop   docker-desktop   dev
$ kubectl apply -f k8s/dockercoins.yaml
deployment.apps/hasher created
service/hasher created
deployment.apps/redis created
service/redis created
deployment.apps/rng created
service/rng created
deployment.apps/webui created
deployment.apps/worker created
$ kubectl get all
NAME                          READY   STATUS    RESTARTS   AGE
pod/hasher-99bbd4bb-g7lln     1/1     Running   0          7s
pod/redis-7b47f84cc4-c7w96    1/1     Running   0          7s
pod/rng-65d885d498-vxb4l      1/1     Running   0          7s
pod/webui-74bb6bbc59-zdww5    1/1     Running   0          7s
pod/worker-5c6f84b477-cpm2j   1/1     Running   0          7s

NAME             TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
service/hasher   ClusterIP   10.111.143.217   <none>        80/TCP     7s
service/redis    ClusterIP   10.110.43.0      <none>        6379/TCP   7s
service/rng      ClusterIP   10.97.48.98      <none>        80/TCP     7s

NAME                     READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/hasher   1/1     1            1           7s
deployment.apps/redis    1/1     1            1           7s
deployment.apps/rng      1/1     1            1           7s
deployment.apps/webui    1/1     1            1           7s
deployment.apps/worker   1/1     1            1           7s

NAME                                DESIRED   CURRENT   READY   AGE
replicaset.apps/hasher-99bbd4bb     1         1         1       7s
replicaset.apps/redis-7b47f84cc4    1         1         1       7s
replicaset.apps/rng-65d885d498      1         1         1       7s
replicaset.apps/webui-74bb6bbc59    1         1         1       7s
replicaset.apps/worker-5c6f84b477   1         1         1       7s
```
- If we remove a resource from the `dockercoins.yaml` file and apply
```sh
$ code dockercoins.yaml
# remove worker deployment section in the dockercoins.yaml
$  kubectl apply -f k8s/dockercoins.yaml
deployment.apps/hasher unchanged
service/hasher unchanged
deployment.apps/redis unchanged
service/redis unchanged
deployment.apps/rng unchanged
service/rng unchanged
deployment.apps/webui unchanged
$ kubectl get all
NAME                          READY   STATUS    RESTARTS   AGE
pod/hasher-99bbd4bb-g7lln     1/1     Running   0          2m41s
pod/redis-7b47f84cc4-c7w96    1/1     Running   0          2m41s
pod/rng-65d885d498-vxb4l      1/1     Running   0          2m41s
pod/webui-74bb6bbc59-zdww5    1/1     Running   0          2m41s
pod/worker-5c6f84b477-cpm2j   1/1     Running   0          2m41s

NAME             TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
service/hasher   ClusterIP   10.111.143.217   <none>        80/TCP     2m41s
service/redis    ClusterIP   10.110.43.0      <none>        6379/TCP   2m41s
service/rng      ClusterIP   10.97.48.98      <none>        80/TCP     2m41s

NAME                     READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/hasher   1/1     1            1           2m41s
deployment.apps/redis    1/1     1            1           2m41s
deployment.apps/rng      1/1     1            1           2m41s
deployment.apps/webui    1/1     1            1           2m41s
deployment.apps/worker   1/1     1            1           2m41s

NAME                                DESIRED   CURRENT   READY   AGE
replicaset.apps/hasher-99bbd4bb     1         1         1       2m41s
replicaset.apps/redis-7b47f84cc4    1         1         1       2m41s
replicaset.apps/rng-65d885d498      1         1         1       2m41s
replicaset.apps/webui-74bb6bbc59    1         1         1       2m41s
replicaset.apps/worker-5c6f84b477   1         1         1       2m41s
```
- thus, `kubectl apply -f` doesn't remove the stuff removed from the yaml file
- if we actually want to remove some of the resources, during our GitOps workflow, use the `--prune` flag
```sh
$ kubectl apply -f k8s/dockercoins.yaml --prune
error: all resources selected for prune without explicitly passing --all. To prune all resources, pass the --all flag. If you did not mean to prune all resources, specify a label selector
```
- the error is about:
	- maybe we deployed multiple things in that namespace, like other applications in that namespace
	- so the error is asking about which one we want to remove or prune
- we could either prune all, if we are super sure
```sh
$ kubectl apply -f dockercoins.yaml --prune --all
```
- or we can prune with a label selector filter
```sh
$ kubectl label deployments --all release=dev
deployment.apps/hasher labeled
deployment.apps/redis labeled
deployment.apps/rng labeled
deployment.apps/webui labeled
deployment.apps/worker labeled

$ kubectl label service --all release=dev
service/hasher labeled
service/redis labeled
service/rng labeled

$ kubectl label replicaset --all release=dev
replicaset.apps/hasher-99bbd4bb labeled
replicaset.apps/redis-7b47f84cc4 labeled
replicaset.apps/rng-65d885d498 labeled
replicaset.apps/webui-74bb6bbc59 labeled
replicaset.apps/worker-5c6f84b477 labeled

$ kubectl apply -f k8s/dockercoins.yaml --prune
error: all resources selected for prune without explicitly passing --all. To prune all resources, pass the --all flag. If you did not mean to prune all resources, specify a label selector

$ kubectl apply -f k8s/dockercoins.yaml --prune --selector release=dev
error: no objects passed to apply
# this didn't work because we need to have the label in the yaml file as well
```
- basically, if we have multiple applications and multiple groups of objects and multiple yaml files, together in a same single namespace, then we can use the label selector `--selector` to indicate which group we are targeting while we do the `--prune`
- the k8s documentation kind of shows a warning while mentioning about the `--prune` flag because, the idea here is that
	- if we want we can implement our own CI/CD pipeline with this `kubectl apply -f --prune` etc
	- BUT, we have many tools like Helm, ArgoCD, Flux etc, that are specifically designed to do that for us, and handle a lot of edge cases