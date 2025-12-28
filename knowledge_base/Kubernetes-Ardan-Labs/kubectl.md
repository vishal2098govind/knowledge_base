#k8s
- `kubectl` is (almost) the only tool we'll need to talk to k8s
- `kubectl` is a CLI - basically a glorified curl
	- everything which we can do with `kubectl` we could also do with `curl` directly with the API
- on our machines, there's a `~/.kube/config` file with:
	- the k8s API address
	- the path to our TLS certificates used to authenticate
- can use `--kubeconfig` flag to pass a config file
- or directly `--server`, `--user` etc.
- `kubectl` can be pronounced "Cube C T L", "Cube cuttle", "Cube cuddle" ...

```sh
$ kubectl get nodes -v6
I1228 12:05:36.854904   82896 loader.go:402] Config loaded from file:  /Users/govind/.kube/config
I1228 12:05:36.855194   82896 envvar.go:172] "Feature gate default state" feature="ClientsPreferCBOR" enabled=false
I1228 12:05:36.855207   82896 envvar.go:172] "Feature gate default state" feature="InformerResourceVersion" enabled=false
I1228 12:05:36.855213   82896 envvar.go:172] "Feature gate default state" feature="WatchListClient" enabled=false
I1228 12:05:36.855217   82896 envvar.go:172] "Feature gate default state" feature="ClientsAllowCBOR" enabled=false
I1228 12:05:36.856758   82896 round_trippers.go:560] GET http://localhost:8080/api?timeout=32s  in 1 milliseconds
E1228 12:05:36.856879   82896 memcache.go:265] "Unhandled Error" err="couldn't get current server API group list: Get \"http://localhost:8080/api?timeout=32s\": dial tcp [::1]:8080: connect: connection refused"
I1228 12:05:36.857994   82896 cached_discovery.go:120] skipped caching discovery info due to Get "http://localhost:8080/api?timeout=32s": dial tcp [::1]:8080: connect: connection refused
I1228 12:05:36.858687   82896 round_trippers.go:560] GET http://localhost:8080/api?timeout=32s  in 0 milliseconds
E1228 12:05:36.858765   82896 memcache.go:265] "Unhandled Error" err="couldn't get current server API group list: Get \"http://localhost:8080/api?timeout=32s\": dial tcp [::1]:8080: connect: connection refused"
I1228 12:05:36.859972   82896 cached_discovery.go:120] skipped caching discovery info due to Get "http://localhost:8080/api?timeout=32s": dial tcp [::1]:8080: connect: connection refused
I1228 12:05:36.859995   82896 shortcut.go:103] Error loading discovery information: Get "http://localhost:8080/api?timeout=32s": dial tcp [::1]:8080: connect: connection refused
I1228 12:05:36.860528   82896 round_trippers.go:560] GET http://localhost:8080/api?timeout=32s  in 0 milliseconds
E1228 12:05:36.860565   82896 memcache.go:265] "Unhandled Error" err="couldn't get current server API group list: Get \"http://localhost:8080/api?timeout=32s\": dial tcp [::1]:8080: connect: connection refused"
I1228 12:05:36.861711   82896 cached_discovery.go:120] skipped caching discovery info due to Get "http://localhost:8080/api?timeout=32s": dial tcp [::1]:8080: connect: connection refused
I1228 12:05:36.862252   82896 round_trippers.go:560] GET http://localhost:8080/api?timeout=32s  in 0 milliseconds
E1228 12:05:36.862288   82896 memcache.go:265] "Unhandled Error" err="couldn't get current server API group list: Get \"http://localhost:8080/api?timeout=32s\": dial tcp [::1]:8080: connect: connection refused"
I1228 12:05:36.863436   82896 cached_discovery.go:120] skipped caching discovery info due to Get "http://localhost:8080/api?timeout=32s": dial tcp [::1]:8080: connect: connection refused
I1228 12:05:36.863936   82896 round_trippers.go:560] GET http://localhost:8080/api?timeout=32s  in 0 milliseconds
E1228 12:05:36.863972   82896 memcache.go:265] "Unhandled Error" err="couldn't get current server API group list: Get \"http://localhost:8080/api?timeout=32s\": dial tcp [::1]:8080: connect: connection refused"
I1228 12:05:36.865144   82896 cached_discovery.go:120] skipped caching discovery info due to Get "http://localhost:8080/api?timeout=32s": dial tcp [::1]:8080: connect: connection refused
I1228 12:05:36.865205   82896 helpers.go:264] Connection error: Get http://localhost:8080/api?timeout=32s: dial tcp [::1]:8080: connect: connection refused
The connection to the server localhost:8080 was refused - did you specify the right host or port?
$ cat ~/.kube/config
apiVersion: v1
kind: Config
```
- an usual kube config would look something like:
![[Pasted image 20251228120824.png]]

### `kubectl` is the new SSH
- When working with cloud, we often create few VMs, and then we are going to connect to the VM via ssh, install packages, check logs, metrics, etc with tools like top, vmstat, free etc
- this is ok if we have one or few servers
- if we have more servers, dozens, hundreds, thousands of servers, we can't use ssh as it doesn't scale very well
- in the beginning we can manually create k8s clusters using `kubectl`
- over time, we will put in place automation

## `kubectl get`
- to get any resources within k8s cluster
```sh
# kebectl get <resource>
$ kubectl get node # to get node resources in our cluster
# also same result for nodes (plural) i.e. kubectl get nodes
$ kubectl get no # short form for node
$ kubectl get node node1
```

- to change output format:
```sh
$ kubectl get nodes -o wide
$ kubectl get node node1 -o json
$ kubectl get node node1 -o yaml
```

- almost every k8s api request and response would be having these fields:
	- apiVersion
	- kind
	- metadata

## Type names
- the most common resource names have three forms:
	- **singular** - node, service, deployment, pod, customresourcedefinition, replicaset, statefulset, daemonset
	- **plural** - nodes, services, deployments, pods, customresourcedefinitions, replicasets, statefulsets, daemonsets
	- **short** - no, svc, deploy, po, crd, rs, sts, ds
- some resources do not have short names
- `Endpoints` only have a plural form
	- even a single `Endpoints` resource is actually a list of endpoints

## `kubectl describe`
- on doing `kubectl get node node1 -o yaml`, we get a lot of info, including images on the node, but it doesn't give list of pods and containers on that node
- to get pods and containers running in a node, we can use `kubectl describe` command
- `kubectl describe` is going to be our best friend when troubleshooting almost any kind of issue in k8s, like
	- a pod doesn't stop
	- a node is behaving in a weird way
	- a service that times out, don't know what is happening with load balancer
```sh
# kubectl describe <resource-type> <resource-identifier or name>
$ kubectl describe node node1
```
![[Pasted image 20251228130046.png]]
- the reason why `kubectl describe` would be very useful for trouble shooting is at the very end, we have `Events` section which shows lot of useful information

## Namespaces in k8s
- namespaces are themselves k8s subjects
- can do `kubectl get namespaces` or `kubectl get ns`
- the 4 basic namespaces that we would see on any k8s cluster
	- default
	- kube-node-lease
	- kube-public
	- kube-system
- directly using `kubectl get` would use `default` namespace
- to specify namespace, we need to use 
	- `kubectl get pods --namespace kube-system`
	- or `kubectl get pods -n kube-system`
- to see in all namespaces
	- `kubectl get pods --all-namespaces`
- namespaces are just a way to organise resources **logically**
- they do not reflect the topology or the organization of the cluster