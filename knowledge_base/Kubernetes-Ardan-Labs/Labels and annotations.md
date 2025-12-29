#k8s

```sh
$ kubectl describe node
Name:               docker-desktop
Roles:              control-plane
Labels:             beta.kubernetes.io/arch=amd64
                    beta.kubernetes.io/os=linux
                    kubernetes.io/arch=amd64
                    kubernetes.io/hostname=docker-desktop
                    kubernetes.io/os=linux
                    node-role.kubernetes.io/control-plane=
                    node.kubernetes.io/exclude-from-external-load-balancers=
Annotations:        kubeadm.alpha.kubernetes.io/cri-socket: unix:///var/run/cri-dockerd.sock
                    node.alpha.kubernetes.io/ttl: 0
                    volumes.kubernetes.io/controller-managed-attach-detach: true
```

## What are labels and annotations
- arbitrary pieces of meta data that we can attach to our objects
- not just to nodes, but to pretty much everything in k8s
- if we run our cluster in a cloud provider, the nodes might have meta data about 
	- the region
	- instance type
	- etc
- if we want, we can put our own labels
	- let's say that we have our nodes in a data center in different racks
	- if we want to keep track of which rack a node has been kept
	- we can do
```sh
$ kubectl label node docker-desktop rack=2
node/docker-desktop labeled
➜  ~ kubectl describe node docker-desktop
Name:               docker-desktop
Roles:              control-plane
Labels:             beta.kubernetes.io/arch=amd64
                    beta.kubernetes.io/os=linux
                    kubernetes.io/arch=amd64
                    kubernetes.io/hostname=docker-desktop
                    kubernetes.io/os=linux
                    node-role.kubernetes.io/control-plane=
                    node.kubernetes.io/exclude-from-external-load-balancers=
                    rack=2
Annotations:        kubeadm.alpha.kubernetes.io/cri-socket: unix:///var/run/cri-dockerd.sock
                    node.alpha.kubernetes.io/ttl: 0
                    volumes.kubernetes.io/controller-managed-attach-detach: true
```


## Querying and Filtering with labels
- we can query by labels: (`--label-columns` or `-L`)
```sh
$ kubectl get nodes --label-columns rack
NAME             STATUS   ROLES           AGE   VERSION   RACK
docker-desktop   Ready    control-plane   25h   v1.32.2   2
$  kubectl get nodes -L rack
NAME             STATUS   ROLES           AGE   VERSION   RACK
docker-desktop   Ready    control-plane   25h   v1.32.2   2

# query by multiple labels (comma seperated, without space)
$ kubectl get nodes -L rack,kubernetes.io/arch
NAME             STATUS   ROLES           AGE   VERSION   RACK   ARCH
docker-desktop   Ready    control-plane   25h   v1.32.2   2      amd64
```
- selecting specific labels
```sh
$ kubectl get nodes -L rack,kubernetes.io/arch --selector rack=2
NAME             STATUS   ROLES           AGE   VERSION   RACK   ARCH
docker-desktop   Ready    control-plane   25h   v1.32.2   2      amd64

# query for nodes with rack not equal to 2
$ kubectl get nodes -L rack,kubernetes.io/arch --selector rack!=2
# query for nodes with any value of rack 
$ kubectl get nodes -L rack,kubernetes.io/arch --selector rack
# query for nodes with no rack label
$ kubectl get nodes -L rack,kubernetes.io/arch --selector '!rack'
# query for nodes with rack in (2, 4)
$ kubectl get nodes -L rack,kubernetes.io/arch --selector 'rack in (2, 4)'
# query for nodes with rack not in (2, 4)
$ kubectl get nodes -L rack,kubernetes.io/arch --selector 'rack notin (2, 4)'
```
- combine selectors
```sh
$ kubectl label node node1 deployed-when=2021
$ kubectl label node node2 deployed-when=2022
$ kubectl label node node3 deployed-when=2022
$ kubectl label node node4 deployed-when=2022

# query for nodes with rack = 2 and deployed-when=2022 (comma seperated without spaces)
$ kubectl get nodes -L rack,kubernetes.io/arch --selector rack=2,deployed-when=2022
$ kubectl get nodes -L rack,kubernetes.io/arch --selector rack=2,'deployed-when > 2021'
```

- we can see labels during `kubectl get`
```sh
$ kubectl get nodes --show-labels

# filter
$ kubectl get nodes --show-labels --selector rack=2,deployed-when=2022
```

### Why filtering using `--selector` is recommended
- applying `--selector` is equivalent to 
	- `kubectl get nodes --show-labels | grep <condition>`
	- but using grep is client side filtering
	- using `--selector` is server side filter, at the control-plane node
- the **filtering** using `--selector` (server side filtering) is helpful if we have 1000s of nodes in our clusters or we filtering over pods with our cluster having lots of pods,  that can be a huge save of load (both network and CPU) on control-plane from avoiding to transfer data about all the pods over the network
- this is similar to doing `SELECT * FROM table;`  without any `WHERE` filter in a database vs adding `WHERE` filter in the SQL  query

## Difference between labels and annotations
```sh
$ kubectl annotate node node1 admin="Vishal Govind"
node/node1 annotated
$ kubectl annotate node node1 mood="🥳"
```
- labels do not allow special characters like emojis or accents
- annotations allow emojis and accents
- so, with annotations we can have any value we want including fancy emojis and accents, while with labels we have to stick with letters, numbers, . , _ and -
- there is also size limit for values in labels
- no size limit for annotations
- **selectors** only work on labels
- if we want to put arbitrary information, we can use annotations
- if we want to put some information that we can use in a selector to filter our `kubectl get` results, we can use labels
- labels are used typically to group things, to keep track of which thing belongs to what
- annotations are used if we want to add a little bundle of data that we just want to keep with that object, it could be
	- a configuration file
	- png icon we could see in the supervision dashboard
	- name or email address or phone number of the person who is on core of that service

### See all labels
- if we create deployments using `kubectl`, it adds default labels of `app=<name>`
```sh
kubectl get all --show-labels
NAME                          READY   STATUS    RESTARTS   AGE    LABELS
pod/hasher-99bbd4bb-xblfb     1/1     Running   0          90s    app=hasher,pod-template-hash=99bbd4bb
pod/redis-7b47f84cc4-6kbp2    1/1     Running   0          107s   app=redis,pod-template-hash=7b47f84cc4
pod/rng-65d885d498-vq5h8      1/1     Running   0          79s    app=rng,pod-template-hash=65d885d498
pod/webui-74bb6bbc59-rmgng    1/1     Running   0          64s    app=webui,pod-template-hash=74bb6bbc59
pod/worker-5c6f84b477-gw6jn   1/1     Running   0          2m4s   app=worker,pod-template-hash=5c6f84b477

NAME                 TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE   LABELS
service/hasher       ClusterIP      10.103.124.226   <none>        80/TCP         25s   app=hasher
service/kubernetes   ClusterIP      10.96.0.1        <none>        443/TCP        26h   component=apiserver,provider=kubernetes
service/redis        ClusterIP      10.109.22.65     <none>        6379/TCP       44s   app=redis
service/rng          ClusterIP      10.106.8.225     <none>        80/TCP         37s   app=rng
service/webui        LoadBalancer   10.96.255.64     localhost     80:32744/TCP   9s    app=webui

NAME                     READY   UP-TO-DATE   AVAILABLE   AGE    LABELS
deployment.apps/hasher   1/1     1            1           90s    app=hasher
deployment.apps/redis    1/1     1            1           107s   app=redis
deployment.apps/rng      1/1     1            1           79s    app=rng
deployment.apps/webui    1/1     1            1           64s    app=webui
deployment.apps/worker   1/1     1            1           2m4s   app=worker

NAME                                DESIRED   CURRENT   READY   AGE    LABELS
replicaset.apps/hasher-99bbd4bb     1         1         1       90s    app=hasher,pod-template-hash=99bbd4bb
replicaset.apps/redis-7b47f84cc4    1         1         1       107s   app=redis,pod-template-hash=7b47f84cc4
replicaset.apps/rng-65d885d498      1         1         1       79s    app=rng,pod-template-hash=65d885d498
replicaset.apps/webui-74bb6bbc59    1         1         1       64s    app=webui,pod-template-hash=74bb6bbc59
replicaset.apps/worker-5c6f84b477   1         1         1       2m4s   app=worker,pod-template-hash=5c6f84b477
```
- the labels given to the deployment propagate all over
	- replica sets
	- services, if any
	- pods
- thus, we can filter by the `app` label using `--selector`
```sh
$ kubectl get all --show-labels --selector app=worker
NAME                          READY   STATUS    RESTARTS   AGE     LABELS
pod/worker-5c6f84b477-gw6jn   1/1     Running   0          5m33s   app=worker,pod-template-hash=5c6f84b477

NAME                     READY   UP-TO-DATE   AVAILABLE   AGE     LABELS
deployment.apps/worker   1/1     1            1           5m33s   app=worker

NAME                                DESIRED   CURRENT   READY   AGE     LABELS
replicaset.apps/worker-5c6f84b477   1         1         1       5m33s   app=worker,pod-template-hash=5c6f84b477

$ kubectl get all --show-labels --selector app=webui
NAME                         READY   STATUS    RESTARTS   AGE     LABELS
pod/webui-74bb6bbc59-rmgng   1/1     Running   0          5m34s   app=webui,pod-template-hash=74bb6bbc59

NAME            TYPE           CLUSTER-IP     EXTERNAL-IP   PORT(S)        AGE     LABELS
service/webui   LoadBalancer   10.96.255.64   localhost     80:32744/TCP   4m39s   app=webui

NAME                    READY   UP-TO-DATE   AVAILABLE   AGE     LABELS
deployment.apps/webui   1/1     1            1           5m34s   app=webui

NAME                               DESIRED   CURRENT   READY   AGE     LABELS
replicaset.apps/webui-74bb6bbc59   1         1         1       5m34s   app=webui,pod-template-hash=74bb6bbc59
```
- when we create a pod using `kubectl run`, they are given default label of `run` by `kubectl`