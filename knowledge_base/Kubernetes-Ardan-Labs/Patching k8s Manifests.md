#k8s

## Strategic merge patch
- we can edit the manifests of k8s objects directly by `kubectl edit deployment worker`
- k8s supports many other ways as well to accept patches of manifests to be updated
```sh
$ kubectl get service webui -o yaml > service-webui.yaml
$ cp service-webui.yaml service-webui-new.yaml
$ code service-webui-new.yaml
# make changes to the yaml
$ diff service-webui.yaml service-webui-new.yaml
29c29
< type: LoadBalancer
---
> type: NodePort
```
- if we are working with source code, that's typically how we would create a patch and then apply that patch
- this is a text diff which works great for source code, or documentation, or any kind of text, except here, we are not working with text.
- the yaml is just the representation of the data that we have on `etcd` on the control plane
- but in reality, we're working with data
- thus, if we want to just update few parts of the manifest, say type of service to be changed to `NodePort`, we can do a patch like this:
```sh
$ kubectl get service worker-lb -o yaml
apiVersion: v1
kind: Service
metadata:
  creationTimestamp: "2025-12-29T07:22:22Z"
  labels:
    app: worker
  name: worker-lb
  namespace: default
  resourceVersion: "59784"
  uid: 48fb9338-a21b-4533-8df9-c363161e0e2f
spec:
  allocateLoadBalancerNodePorts: true
  clusterIP: 10.109.2.204
  clusterIPs:
  - 10.109.2.204
  externalTrafficPolicy: Cluster
  internalTrafficPolicy: Cluster
  ipFamilies:
  - IPv4
  ipFamilyPolicy: SingleStack
  ports:
  - nodePort: 31502
    port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: worker
  sessionAffinity: None
  type: LoadBalancer
status:
  loadBalancer: {}
$ code node-port-patch.yaml
spec:
	type: NodePort
$ kubectl patch service webui --patch-file node-port-patch.yaml
service/webui patched
```
- can also do inline patch:
```sh
$ kubectl patch service blue --patch "
> spec:
> 	type: ClusterIP
> "
service/blue patched
```