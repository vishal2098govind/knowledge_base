#k8s

- Services are **network things**
- A service is a stable endpoint to connect to **something**
	- in the initial proposal, they were called "portals"
		- a portal is like a door that take us somewhere
		- also service is something that let's us connect to something
- can we be more specific about **something**?
	- no we can't
	- because, sometimes, that something is going to be 
		- a **pod** in the cluster
		- or something **outside of the cluster**
		- or something that can be identified as an **IP address**
		- or something that can be identified as a **DNS name**
- to get the kind of services present in our cluster
	- `kubectl get services`

## ClusterIP Services
```sh
$ kubectl get services
NAME         TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
kubernetes   ClusterIP   10.96.0.1    <none>        443/TCP   5m22s
```
- cluster ip means, it works purely internally
- it doesn't work outside of the cluster
- it is a private internal ip address
- cluster ip can be different among clusters
- A `ClusterIP` service is internal, available from within the cluster only
- This is useful for introspection from within containers
- Each `ClusterIP` service also gets `DNS` integration
- which is why, if we are within one of the pod we can also do
```sh
$ ping kubernetes.default.svc
# this resolves to 10.96.0.1, only if we are within one of the pod of the cluster
```
- the DNS resolves service names to their cluster IP