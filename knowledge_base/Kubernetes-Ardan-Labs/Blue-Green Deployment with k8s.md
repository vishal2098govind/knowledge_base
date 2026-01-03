#k8s #blue-green-deployment
```sh
$ kubectl create namespace rainbow
namespace/rainbow created

$ kubectl config get-contexts
CURRENT   NAME             CLUSTER          AUTHINFO         NAMESPACE
*         docker-desktop   docker-desktop   docker-desktop   dev

$ kubectl config set-context --namespace rainbow --current
Context "docker-desktop" modified.

$ kubectl config get-contexts
CURRENT   NAME             CLUSTER          AUTHINFO         NAMESPACE
*         docker-desktop   docker-desktop   docker-desktop   rainbow
```
- create blue and green stacks
```sh
$ kubectl get pods
No resources found in rainbow namespace.

$ kubectl create deployment blue --image jpetazzo/color --replicas 2
deployment.apps/blue created
$ kubectl create deployment green --image jpetazzo/color --replicas 2
deployment.apps/green created
```
- create a frontend service to switch from blue to green
```sh
$ kubectl create service nodeport frontend --tcp 80
service/frontend created

# as of now, frontend service doesn't point to any pod
$ kubectl describe service frontend
Name:                     frontend
Namespace:                rainbow
Labels:                   app=frontend
Annotations:              <none>
Selector:                 app=frontend
Type:                     NodePort
IP Family Policy:         SingleStack
IP Families:              IPv4
IP:                       10.110.110.193
IPs:                      10.110.110.193
Port:                     80  80/TCP
TargetPort:               80/TCP
NodePort:                 80  30333/TCP
Endpoints:
Session Affinity:         None
External Traffic Policy:  Cluster
Internal Traffic Policy:  Cluster
Events:                   <none>
```
- to make the load balancer point to any of the pods, we need to update the selectors in the service
```sh
# currently the load balancer selects pods which match label of app=frontend
$ kubectl get service frontend -o yaml
apiVersion: v1
kind: Service
metadata:
  creationTimestamp: "2026-01-03T07:41:22Z"
  labels:
    app: frontend
  name: frontend
  namespace: rainbow
  resourceVersion: "152227"
  uid: 2d8c37e2-377f-4b21-8427-ac3e0a9a9380
spec:
  clusterIP: 10.110.110.193
  clusterIPs:
  - 10.110.110.193
  externalTrafficPolicy: Cluster
  internalTrafficPolicy: Cluster
  ipFamilies:
  - IPv4
  ipFamilyPolicy: SingleStack
  ports:
  - name: "80"
    nodePort: 30333
    port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: frontend
  sessionAffinity: None
  type: NodePort
status:
  loadBalancer: {}
  
# update the load balancer selector to match for label app=blue
$ kubectl edit service frontend
# spec.selector = app: blue
service/frontend edited
$ kubectl get service frontend -o yaml
apiVersion: v1
kind: Service
metadata:
  creationTimestamp: "2026-01-03T07:41:22Z"
  labels:
    app: frontend
  name: frontend
  namespace: rainbow
  resourceVersion: "152948"
  uid: 2d8c37e2-377f-4b21-8427-ac3e0a9a9380
spec:
  clusterIP: 10.110.110.193
  clusterIPs:
  - 10.110.110.193
  externalTrafficPolicy: Cluster
  internalTrafficPolicy: Cluster
  ipFamilies:
  - IPv4
  ipFamilyPolicy: SingleStack
  ports:
  - name: "80"
    nodePort: 30333
    port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: blue
  sessionAffinity: None
  type: NodePort
status:
  loadBalancer: {}
$ kubectl describe service frontend
Name:                     frontend
Namespace:                rainbow
Labels:                   app=frontend
Annotations:              <none>
Selector:                 app=blue
Type:                     NodePort
IP Family Policy:         SingleStack
IP Families:              IPv4
IP:                       10.110.110.193
IPs:                      10.110.110.193
Port:                     80  80/TCP
TargetPort:               80/TCP
NodePort:                 80  30333/TCP
Endpoints:                10.1.0.147:80,10.1.0.148:80 # now the load balancer points to the two pods of blue stack
Session Affinity:         None
External Traffic Policy:  Cluster
Internal Traffic Policy:  Cluster
Events:                   <none>

$ kubectl get pods --show-labels -o wide
NAME                     READY   STATUS    RESTARTS   AGE   IP           NODE             NOMINATED NODE   READINESS GATES   LABELS
blue-5c986bd7bf-5z6dz    1/1     Running   0          32m   10.1.0.148   docker-desktop   <none>           <none>            app=blue,pod-template-hash=5c986bd7bf
blue-5c986bd7bf-7qql7    1/1     Running   0          32m   10.1.0.147   docker-desktop   <none>           <none>            app=blue,pod-template-hash=5c986bd7bf
green-556754bb7d-7cg4b   1/1     Running   0          32m   10.1.0.149   docker-desktop   <none>           <none>            app=green,pod-template-hash=556754bb7d
green-556754bb7d-g8rrd   1/1     Running   0          32m   10.1.0.150   docker-desktop   <none>           <none>            app=green,pod-template-hash=556754bb7d
```
- for load balancer to point to green stack
```sh
$ kubectl edit service frontend
# spec.selector = app: green
service/frontend edited

$ kubectl describe service frontend
Name:                     frontend
Namespace:                rainbow
Labels:                   app=frontend
Annotations:              <none>
Selector:                 app=green
Type:                     LoadBalancer
IP Family Policy:         SingleStack
IP Families:              IPv4
IP:                       10.110.110.193
IPs:                      10.110.110.193
Port:                     80  80/TCP
TargetPort:               80/TCP
NodePort:                 80  30333/TCP
Endpoints:                10.1.0.150:80,10.1.0.149:80 # now the load balancer points to green pods
Session Affinity:         None
External Traffic Policy:  Cluster
Internal Traffic Policy:  Cluster
Events:                   <none>

$ kubectl get pods --show-labels -o wide
NAME                     READY   STATUS    RESTARTS   AGE   IP           NODE             NOMINATED NODE   READINESS GATES   LABELS
blue-5c986bd7bf-5z6dz    1/1     Running   0          38m   10.1.0.148   docker-desktop   <none>           <none>            app=blue,pod-template-hash=5c986bd7bf
blue-5c986bd7bf-7qql7    1/1     Running   0          38m   10.1.0.147   docker-desktop   <none>           <none>            app=blue,pod-template-hash=5c986bd7bf
green-556754bb7d-7cg4b   1/1     Running   0          38m   10.1.0.149   docker-desktop   <none>           <none>            app=green,pod-template-hash=556754bb7d
green-556754bb7d-g8rrd   1/1     Running   0          38m   10.1.0.150   docker-desktop   <none>           <none>            app=green,pod-template-hash=556754bb7d
```
- instead of editing the selector in service yaml file manually, we can edit it using the `kubectl patch` as well, in a CI/CD pipeline and automation
```sh
$ kubectl patch service frontend --patch '{"spec": {"selector": {"app": "blue"}}}'
service/frontend patched
➜  container.training git:(main) ✗ kubectl get service frontend -o yaml
apiVersion: v1
kind: Service
metadata:
  creationTimestamp: "2026-01-03T07:41:22Z"
  labels:
    app: frontend
  name: frontend
  namespace: rainbow
  resourceVersion: "153935"
  uid: 2d8c37e2-377f-4b21-8427-ac3e0a9a9380
spec:
  allocateLoadBalancerNodePorts: true
  clusterIP: 10.110.110.193
  clusterIPs:
  - 10.110.110.193
  externalTrafficPolicy: Cluster
  internalTrafficPolicy: Cluster
  ipFamilies:
  - IPv4
  ipFamilyPolicy: SingleStack
  ports:
  - name: "80"
    nodePort: 30333
    port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: blue
  sessionAffinity: None
  type: LoadBalancer
status:
  loadBalancer: {}

$ kubectl describe service frontend
Name:                     frontend
Namespace:                rainbow
Labels:                   app=frontend
Annotations:              <none>
Selector:                 app=blue
Type:                     LoadBalancer
IP Family Policy:         SingleStack
IP Families:              IPv4
IP:                       10.110.110.193
IPs:                      10.110.110.193
Port:                     80  80/TCP
TargetPort:               80/TCP
NodePort:                 80  30333/TCP
Endpoints:                10.1.0.148:80,10.1.0.147:80
Session Affinity:         None
External Traffic Policy:  Cluster
Internal Traffic Policy:  Cluster
Events:                   <none>

$ kubectl get pods --show-labels -o wide
NAME                     READY   STATUS    RESTARTS   AGE   IP           NODE             NOMINATED NODE   READINESS GATES   LABELS
blue-5c986bd7bf-5z6dz    1/1     Running   0          46m   10.1.0.148   docker-desktop   <none>           <none>            app=blue,pod-template-hash=5c986bd7bf
blue-5c986bd7bf-7qql7    1/1     Running   0          46m   10.1.0.147   docker-desktop   <none>           <none>            app=blue,pod-template-hash=5c986bd7bf
green-556754bb7d-7cg4b   1/1     Running   0          46m   10.1.0.149   docker-desktop   <none>           <none>            app=green,pod-template-hash=556754bb7d
green-556754bb7d-g8rrd   1/1     Running   0          46m   10.1.0.150   docker-desktop   <none>           <none>            app=green,pod-template-hash=556754bb7d
```