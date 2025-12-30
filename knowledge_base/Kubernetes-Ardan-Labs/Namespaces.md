#k8s
- Namespaces are a way to organize our resources on k8s
- They are a logical separation, not a physical separation
- ~ directories/folders on our computer
- The four basic namespaces that we have on pretty much every modern k8s cluster
```sh
$ kubectl get namespaces
NAME              STATUS   AGE
default           Active   33h
kube-node-lease   Active   33h
kube-public       Active   33h
kube-system       Active   33h
```
- we can create more namespaces

## Creating namespaces
### `kubectl create namespace`
```sh
$ kubectl create namespace dockercoin
namespace/dockercoin created
```
### Using YAML Manifest
```sh
$ kubectl create namespace blue
namespace/blue created
$ kubectl get namespace blue -o yaml
apiVersion: v1
kind: Namespace
metadata:
  creationTimestamp: "2025-12-29T18:42:23Z"
  labels:
    kubernetes.io/metadata.name: blue
  name: blue
  resourceVersion: "104378"
  uid: 79cf35f0-d9cd-4f02-a4bd-2fae29ab3c40
spec:
  finalizers:
  - kubernetes
status:
  phase: Active
```
- although while creating the namespaces using `kubectl create`, it created a yaml manifest with a bunch of stuff, we need not need all that to create a namespace
- many fields are present which are not strictly necessary
	- `metadata.creationTimestamp`: this is something that was automatically added by k8s
	- similarly for the `metadata.resourceVersion`, `uid`
	- even the `spec.finalizers[kubernetes]` and `status.phase=Active` stuff are added by k8s automatically later after we submitted just the `kubectl create namespace blue`
- thus, the bare minimum manifest for creating a namespace would be:
```sh
$ code green.yaml
```
```yaml
apiVersion: v1
kind: Namespace
metadata:
	name: green
```
```sh
# to create a namespace using the manifest green.yaml file:
$ kubectl create -f green.yaml
namespace/green created
$ kubectl get ns
NAME              STATUS   AGE
blue              Active   9m8s
default           Active   33h
green             Active   6s
kube-node-lease   Active   33h
kube-public       Active   33h
kube-system       Active   33h
$ kubectl get namespace green -o yaml
apiVersion: v1
kind: Namespace
metadata:
  creationTimestamp: "2025-12-29T18:51:25Z"
  labels:
    kubernetes.io/metadata.name: green
  name: green
  resourceVersion: "105103"
  uid: 9b37a9a6-b62e-4dde-9d3f-0d2d816a8551
spec:
  finalizers:
  - kubernetes
status:
  phase: Active
```

## Using Namespaces
```sh
$ kubectl create deployment purple --image jpetazzo/color --namespace blue
deployment.apps/purple created
$ kubectl get deployment
NAME       READY   UP-TO-DATE   AVAILABLE   AGE
hasher     1/1     1            1           6h51m
pingpong   4/4     4            4           6h9m
redis      1/1     1            1           6h51m
rng        1/1     1            1           6h51m
webui      1/1     1            1           6h50m
worker     1/1     1            1           6h51m
$ kubectl get deployment --namespace blue
NAME     READY   UP-TO-DATE   AVAILABLE   AGE
purple   1/1     1            1           31s
$ kubectl expose deployment purple --port 80
Error from server (NotFound): deployments.apps "purple" not found
$ kubectl expose deployment purple --port 80 --namespace blue
service/purple exposed
$ kubectl get services --namespace blue
NAME     TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
purple   ClusterIP   10.102.223.15   <none>        80/TCP    47s
kubectl edit service purple -n blue
service/purple edited # edited to change from cluster-ip to load-balancer to be able to access the server from my mac
$ kubectl get service -n blue
NAME     TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
purple   LoadBalancer   10.102.223.15   <pending>     80:30868/TCP   3m3s
$ kubectl get service -n blue
NAME     TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
purple   LoadBalancer   10.102.223.15   <pending>     80:30868/TCP   3m9s
$ kubectl get service -n blue
NAME     TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
purple   LoadBalancer   10.102.223.15   <pending>     80:30868/TCP   3m11s
$ kubectl get service
NAME         TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
hasher       ClusterIP      10.103.124.226   <none>        80/TCP         6h54m
kubernetes   ClusterIP      10.96.0.1        <none>        443/TCP        33h
redis        ClusterIP      10.109.22.65     <none>        6379/TCP       6h55m
rng          ClusterIP      10.106.8.225     <none>        80/TCP         6h55m
webui        LoadBalancer   10.96.255.64     localhost     80:32744/TCP   6h54m
➜  docker-ardan-labs kubectl delete service webui
service "webui" deleted
$ kubectl get service -n blue
NAME     TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
purple   LoadBalancer   10.102.223.15   <pending>     80:30868/TCP   3m30s
$ kubectl get service -n blue
NAME     TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
purple   LoadBalancer   10.102.223.15   localhost     80:30868/TCP   3m58s
$ curl localhost
🔵🟣This is pod blue/purple-65bb9bc655-mqg7f on linux/amd64, serving / for 192.168.65.3:60214.
```
- namespaces are useful to deploy other copies of our application

## Deploying `dockercoins` on `dev` namespace
### using `kubectl`
- we can just re-run the `kubectl create deployment <service>` commands for each of the 5 services with a `--namespace dev`
- similarly run `kubectl create service <service>` for each of the service with a `--namespace` for each command
- here, tedious or annoying task to do would to add `--namespace dev` every single time for each resource, without forgetting
- how can we change the default active namespace while using `kubectl`?
	- deploying `dockercoins` is not too hard
	- just 9 `kubectl` commands in total
		- 5 commands for `deployments`
		- 4 for `services`
	- but is there any way to deploy everything with a single command?
		- yes
```sh
$ kubectl create -f container.training/k8s/dockercoins.yaml --namespace dev
```
#### `dockercoins.yaml`
- in a yaml file, we can have multiple resources separated with `---`
- this `---` is not something specific to k8s, this is standard yaml syntax
- we may add an extra `---` at the beginning and/or at the end of the yaml file
- we don't have to but,
- it's a good idea to have `---` in the beginning and/or because
	- why put at the beginning -
		- it's a visual reminder, on opening the yaml file, that there are, or can be, multiple resources mentioned in this yaml file
	- why put at the end and beginning - 
		- if we have multiple resources in multiple yaml file and if we have to `cat` those or multiple yaml files together, it will appear separated, thanks to those `---` in each file at the beginning and end of each file, and thus will be less confusing
		- not only `cat`, but if we have to combine multiple yaml files later, for any reason, we have to put `---` at that time
```sh
$ cat container.training/k8s/dockercoins.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: hasher
  name: hasher
spec:
  replicas: 1
  selector:
    matchLabels:
      app: hasher
  template:
    metadata:
      labels:
        app: hasher
    spec:
      containers:
      - image: dockercoins/hasher:v0.1
        name: hasher
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: hasher
  name: hasher
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: hasher
  type: ClusterIP
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: redis
  name: redis
spec:
  replicas: 1
  selector:
    matchLabels:
      app: redis
  template:
    metadata:
      labels:
        app: redis
    spec:
      containers:
      - image: redis
        name: redis
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: redis
  name: redis
spec:
  ports:
  - port: 6379
    protocol: TCP
    targetPort: 6379
  selector:
    app: redis
  type: ClusterIP
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: rng
  name: rng
spec:
  replicas: 1
  selector:
    matchLabels:
      app: rng
  template:
    metadata:
      labels:
        app: rng
    spec:
      containers:
      - image: dockercoins/rng:v0.1
        name: rng
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: rng
  name: rng
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: rng
  type: ClusterIP
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: webui
  name: webui
spec:
  replicas: 1
  selector:
    matchLabels:
      app: webui
  template:
    metadata:
      labels:
        app: webui
    spec:
      containers:
      - image: dockercoins/webui:v0.1
        name: webui
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: webui
  name: webui
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: webui
  type: NodePort
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: worker
  name: worker
spec:
  replicas: 1
  selector:
    matchLabels:
      app: worker
  template:
    metadata:
      labels:
        app: worker
    spec:
      containers:
      - image: dockercoins/worker:v0.1
        name: worker
```

### `kubectl apply` vs `kubectl create`
#### `kubectl create`
```sh
$ kubectl create -f container.training/k8s/dockercoins.yaml --namespace dev
deployment.apps/hasher created
service/hasher created
deployment.apps/redis created
service/redis created
deployment.apps/rng created
service/rng created
deployment.apps/webui created
service/webui created
deployment.apps/worker created

$ kubectl create -f container.training/k8s/dockercoins.yaml --namespace dev
Error from server (AlreadyExists): error when creating "container.training/k8s/dockercoins.yaml": deployments.apps "hasher" already exists
Error from server (AlreadyExists): error when creating "container.training/k8s/dockercoins.yaml": services "hasher" already exists
Error from server (AlreadyExists): error when creating "container.training/k8s/dockercoins.yaml": deployments.apps "redis" already exists
Error from server (AlreadyExists): error when creating "container.training/k8s/dockercoins.yaml": services "redis" already exists
Error from server (AlreadyExists): error when creating "container.training/k8s/dockercoins.yaml": deployments.apps "rng" already exists
Error from server (AlreadyExists): error when creating "container.training/k8s/dockercoins.yaml": services "rng" already exists
Error from server (AlreadyExists): error when creating "container.training/k8s/dockercoins.yaml": deployments.apps "webui" already exists
Error from server (AlreadyExists): error when creating "container.training/k8s/dockercoins.yaml": services "webui" already exists
Error from server (AlreadyExists): error when creating "container.training/k8s/dockercoins.yaml": deployments.apps "worker" already exists
```
#### `kubectl apply`
- we can also use `kubectl apply -f dockercoins.yaml --namespace blue`
```sh
$ kubectl apply -f container.training/k8s/dockercoins.yaml --namespace blue
deployment.apps/hasher created
service/hasher created
deployment.apps/redis created
service/redis created
deployment.apps/rng created
service/rng created
deployment.apps/webui created
service/webui created
deployment.apps/worker created

# if we apply again, instead of giving bunch of errors, we get "unchanged" response
$ kubectl apply -f container.training/k8s/dockercoins.yaml --namespace blue
deployment.apps/hasher unchanged
service/hasher unchanged
deployment.apps/redis unchanged
service/redis unchanged
deployment.apps/rng unchanged
service/rng unchanged
deployment.apps/webui unchanged
service/webui unchanged
deployment.apps/worker unchanged

$ kubectl get all -n dev
NAME                          READY   STATUS    RESTARTS   AGE
pod/hasher-99bbd4bb-9vss6     1/1     Running   0          29m
pod/redis-7b47f84cc4-4w4wm    1/1     Running   0          29m
pod/rng-65d885d498-dqwb8      1/1     Running   0          29m
pod/webui-74bb6bbc59-xltz8    1/1     Running   0          29m
pod/worker-5c6f84b477-vjs6j   1/1     Running   0          29m

NAME             TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
service/hasher   ClusterIP   10.100.212.255   <none>        80/TCP         29m
service/redis    ClusterIP   10.99.42.196     <none>        6379/TCP       29m
service/rng      ClusterIP   10.110.244.252   <none>        80/TCP         29m
service/webui    NodePort    10.106.24.250    <none>        80:30106/TCP   29m

NAME                     READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/hasher   1/1     1            1           29m
deployment.apps/redis    1/1     1            1           29m
deployment.apps/rng      1/1     1            1           29m
deployment.apps/webui    1/1     1            1           29m
deployment.apps/worker   1/1     1            1           29m

NAME                                DESIRED   CURRENT   READY   AGE
replicaset.apps/hasher-99bbd4bb     1         1         1       29m
replicaset.apps/redis-7b47f84cc4    1         1         1       29m
replicaset.apps/rng-65d885d498      1         1         1       29m
replicaset.apps/webui-74bb6bbc59    1         1         1       29m
replicaset.apps/worker-5c6f84b477   1         1         1       29m
```
- `kubectl apply` is something like describing infrastructure to the API server - **Infrastructure-as-code IaC**
- if we change the `webui` service to be a `LoadBalancer` instead of `NodePort`, and do `kubectl apply`:
```sh
$ code container.training/k8s/dockercoins.yaml
...
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: webui
  name: webui
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: webui
  type: LoadBalancer
...

$ kubectl apply -f container.training/k8s/dockercoins.yaml -n dev
deployment.apps/hasher unchanged
service/hasher unchanged
deployment.apps/redis unchanged
service/redis unchanged
deployment.apps/rng unchanged
service/rng unchanged
deployment.apps/webui unchanged
service/webui configured # kubectl was able to detect that change
deployment.apps/worker unchanged

$ kubectl get all -n dev
NAME                          READY   STATUS    RESTARTS   AGE
pod/hasher-99bbd4bb-9vss6     1/1     Running   0          33m
pod/redis-7b47f84cc4-4w4wm    1/1     Running   0          33m
pod/rng-65d885d498-dqwb8      1/1     Running   0          33m
pod/webui-74bb6bbc59-xltz8    1/1     Running   0          33m
pod/worker-5c6f84b477-vjs6j   1/1     Running   0          33m

NAME             TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
service/hasher   ClusterIP      10.100.212.255   <none>        80/TCP         33m
service/redis    ClusterIP      10.99.42.196     <none>        6379/TCP       33m
service/rng      ClusterIP      10.110.244.252   <none>        80/TCP         33m
service/webui    LoadBalancer   10.106.24.250    <pending>     80:30106/TCP   33m

NAME                     READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/hasher   1/1     1            1           33m
deployment.apps/redis    1/1     1            1           33m
deployment.apps/rng      1/1     1            1           33m
deployment.apps/webui    1/1     1            1           33m
deployment.apps/worker   1/1     1            1           33m

NAME                                DESIRED   CURRENT   READY   AGE
replicaset.apps/hasher-99bbd4bb     1         1         1       33m
replicaset.apps/redis-7b47f84cc4    1         1         1       33m
replicaset.apps/rng-65d885d498      1         1         1       33m
replicaset.apps/webui-74bb6bbc59    1         1         1       33m
replicaset.apps/worker-5c6f84b477   1         1         1       33m
```
- `kubectl apply -f file.yaml` is extremely powerful because it will let us keep track of changes that we make, keep history of these changes
- for instance, if we made some changes and applied to the cluster, and then realize that it was a mistake, and go back to the previous version, we can revert our changes in the `file.yaml` and apply again
- so when making changes to file.yaml, 
	- first change the `file.yaml`
	- commit changes
	- push to a new branch
	- make a pull request
	- code review
	- merge pull request to main or prod branch
	- some automation would kick in
	- and the automation performs `kubectl apply -f file.yaml`
	- if something is not working, can revert back to previous commit
- this pattern or idea of having our manifests in a git repository instead of directly on the machine, we call it **GitOps**
	- keeping manifest yaml files in git repo
	- use git repo as the source of truth, to decide what should be on the cluster
	- then that's **GitOps**

## Firewalling and Network Policies across namespaces
- to see things in all namespaces, use `--all-namespaces` flag 
```sh
$ kubectl get deployments --all-namespaces
NAMESPACE     NAME       READY   UP-TO-DATE   AVAILABLE   AGE
blue          hasher     1/1     1            1           95m
blue          purple     1/1     1            1           132m
blue          redis      1/1     1            1           95m
blue          rng        1/1     1            1           95m
blue          webui      1/1     1            1           95m
blue          worker     1/1     1            1           95m
default       hasher     1/1     1            1           9h
default       pingpong   4/4     4            4           8h
default       redis      1/1     1            1           9h
default       rng        1/1     1            1           9h
default       webui      1/1     1            1           9h
default       worker     1/1     1            1           9h
dev           hasher     1/1     1            1           98m
dev           redis      1/1     1            1           98m
dev           rng        1/1     1            1           98m
dev           webui      1/1     1            1           98m
dev           worker     1/1     1            1           98m
kube-system   coredns    2/2     2            2           35h
```
- here, there are three copies of `dockercoins`, one in `blue` namespace and another in `dev` namespace and another in the `default` namespace
- these copies are **completely independent** of each other, however, they are **not isolated**
	- i.e. we don't have strict isolation or **`firewalling`** between these copies of `dockercoins`
	- if someone hostile actor manages to hack to the `dockercoins` running on the `blue` namespace, they will be able to connect to the other namespaces as well
	- if we want to prevent that, we will need to use another set of features in k8s, called **Network Policies**, which we can put in place, to define which network traffic is allowed, and which one is denied.

## Changing the default active namespace
- let's say if we want to work for a bit on the `dev` namespace
- and we don't want to put `--namespace` for every `kubectl` command
- there are three ways to change
### `kubens` or `kns`
- it's an extra tool, outside of k8s or `kubectl`
```sh
$ kns dev
```

### Edit `~/.kube/config` file or use `kubectx`
```sh
$ code ~/.kube/config
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: 
    server: https://127.0.0.1:6443
  name: docker-desktop
contexts:
- context:
    cluster: docker-desktop
    namespace: dev # <-------- sets default active namespace
    user: docker-desktop
  name: docker-desktop
current-context: docker-desktop
kind: Config
preferences: {}
users:
- name: docker-desktop
  user:
    client-certificate-data: 
    client-key-data: 
```
```sh
$ kubectl get pods
NAME                      READY   STATUS    RESTARTS   AGE
hasher-99bbd4bb-9vss6     1/1     Running   0          113m
redis-7b47f84cc4-4w4wm    1/1     Running   0          113m
rng-65d885d498-dqwb8      1/1     Running   0          113m
webui-74bb6bbc59-xltz8    1/1     Running   0          113m
worker-5c6f84b477-vjs6j   1/1     Running   0          113m
$ kubectl get pods --namespace default
NAME                        READY   STATUS    RESTARTS   AGE
hasher-99bbd4bb-xblfb       1/1     Running   0          9h
pingpong-5bc478477c-5s54j   1/1     Running   0          6h20m
pingpong-5bc478477c-666bp   1/1     Running   0          6h20m
pingpong-5bc478477c-cf6v6   1/1     Running   0          6h20m
pingpong-5bc478477c-hhpt2   1/1     Running   0          6h20m
redis-7b47f84cc4-6kbp2      1/1     Running   0          9h
rng-65d885d498-vq5h8        1/1     Running   0          9h
webui-74bb6bbc59-rmgng      1/1     Running   0          9h
worker-5c6f84b477-gw6jn     1/1     Running   0          9h
```
- this is simple, but only if we have a simple `.kube/config` file here
- if we have a complex `.kube/config` file like below
![[Pasted image 20251230113043.png]]
- here, there are
	- not one, but two clusters
	- not one, but two users
	- not one, but two **contexts**
- a context is a combination of 
	- cluster
	- user
	- namespace
- this means that, if we need to work on multiple clusters
	- a prod cluster
	- a staging cluster
	- a dev cluster
	- a local cluster
- then, we can have multiple contexts, and then we can switch between contexts with `kubectx`
	- `kubectx` is an external tool outside of k8s or `kubectl`
- when we have only one context in our `~/.kube/config` file, it was pretty easy to change the namespace by just editing the namespace field of that one context in the `~/.kube/config` file
- but, if we want to change current namespace in the current context, it would be difficult to do with a script as it would parse the config yaml and write some code and etc
### `kubectl config`
```sh
$ kubectl config
Modify kubeconfig files using subcommands like "kubectl config set current-context my-context".

 The loading order follows these rules:

  1.  If the --kubeconfig flag is set, then only that file is loaded. The flag may only be set once and no merging takes
place.
  2.  If $KUBECONFIG environment variable is set, then it is used as a list of paths (normal path delimiting rules for
your system). These paths are merged. When a value is modified, it is modified in the file that defines the stanza. When
a value is created, it is created in the first file that exists. If no files in the chain exist, then it creates the
last file in the list.
  3.  Otherwise, ${HOME}/.kube/config is used and no merging takes place.

Available Commands:
  current-context   Display the current-context
  delete-cluster    Delete the specified cluster from the kubeconfig
  delete-context    Delete the specified context from the kubeconfig
  delete-user       Delete the specified user from the kubeconfig
  get-clusters      Display clusters defined in the kubeconfig
  get-contexts      Describe one or many contexts
  get-users         Display users defined in the kubeconfig
  rename-context    Rename a context from the kubeconfig file
  set               Set an individual value in a kubeconfig file
  set-cluster       Set a cluster entry in kubeconfig
  set-context       Set a context entry in kubeconfig
  set-credentials   Set a user entry in kubeconfig
  unset             Unset an individual value in a kubeconfig file
  use-context       Set the current-context in a kubeconfig file
  view              Display merged kubeconfig settings or a specified kubeconfig file

Usage:
  kubectl config SUBCOMMAND [options]

Use "kubectl config <command> --help" for more information about a given command.
Use "kubectl options" for a list of global command-line options (applies to all commands).
```
- `kubectl config` is not even a command, it's a group of commands, to manipulate the `~/.kube/config` file
```sh
$ kubectl config set-context
error: you must specify a non-empty context name or --current
$ kubectl config set-context --help
Set a context entry in kubeconfig.

 Specifying a name that already exists will merge new fields on top of existing values for those fields.

Examples:
  # Set the user field on the gce context entry without touching other values
  kubectl config set-context gce --user=cluster-admin

Options:
    --cluster='':
	cluster for the context entry in kubeconfig

    --current=false:
	Modify the current context

    --namespace='':
	namespace for the context entry in kubeconfig

    --user='':
	user for the context entry in kubeconfig

Usage:
  kubectl config set-context [NAME | --current] [--cluster=cluster_nickname] [--user=user_nickname]
[--namespace=namespace] [options]

Use "kubectl options" for a list of global command-line options (applies to all commands).
```

```sh
# change the default active namespace for current context to dev
$ kubectl config set-context --namespace dev --current
```

### `kube-ps1`
- It's easy to lose track of our current cluster / context / namespace
- `kube-ps1` makes it easy to track these, by showing them in our shell prompt
- It gives us a prompt looking like this one:
```sh
[123.45.67.89] (kubernetes-admin@kubernetes:default) docker@node1 ~
```
- The highlighted part is `context:namespace`, managed by `kube-ps1`
- in the `(kubernetes-admin@kubernetes:default)` it covers the current **context**
	- `kubernetes-admin` is the **user**
	- `kubernetes` is the **cluster** name
	- `default` is the **namespace**
- Highly recommended if when working across multiple contexts or namespaces!