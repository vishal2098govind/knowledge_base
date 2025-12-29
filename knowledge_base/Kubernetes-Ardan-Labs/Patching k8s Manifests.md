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

## Rolling out Restarts of a container or Entering into a container in k8s
3 options
### Delete pod and create new
- delete the pod, and replica set would take care of creating a new one automatically
### Enter into the container and try restarting
- use `kubectl exec`, which is just like `docker exec`
```sh
kubectl get pods
NAME                      READY   STATUS    RESTARTS   AGE
blue-5c986bd7bf-tgh7s     1/1     Running   0          20h
hasher-99bbd4bb-c5q7j     1/1     Running   0          14h
red-77f6d65f98-454qr      1/1     Running   0          19h
red-77f6d65f98-pxlz4      1/1     Running   0          20h
red-77f6d65f98-spv8s      1/1     Running   0          19h
redis-7b47f84cc4-w9fcl    1/1     Running   0          14h
rng-65d885d498-sqjrp      1/1     Running   0          14h
webui-74bb6bbc59-m5n8m    1/1     Running   0          14h
worker-5c6f84b477-pwzll   1/1     Running   0          12h
➜  dockercoins git:(main) kubectl exec -it worker-5c6f84b477-pwzll -- sh
/app # ls
worker.py
/app # cd ..
/ # ls
app    bin    dev    etc    home   lib    media  mnt    opt    proc   root   run    sbin   srv    sys    tmp    usr    var
/ # exit
➜  dockercoins git:(main) kubectl exec -it webui-74bb6bbc59-m5n8m -- sh
/app # ls
Dockerfile         files              node_modules       package-lock.json  package.json       webui.js
/app # cat Dockerfile
FROM node:23-alpine
WORKDIR /app
RUN npm install express
RUN npm install morgan
RUN npm install redis@5
COPY . .
CMD ["node", "webui.js"]
EXPOSE 80
/app # ps aux
PID   USER     TIME  COMMAND
    1 root      0:06 {MainThread} node webui.js
   35 root      0:00 sh
   43 root      0:00 ps aux
```

### Rollout restart
- another option is using `kubectl rollout restart`
	- this creates the new pod and once the new pod is running, it then deletes or terminates the old pod
	- this is better than deleting old pod before starting the new pod
```sh
$ kubectl rollout restart deployment webui
deployment.apps/webui restarted
```