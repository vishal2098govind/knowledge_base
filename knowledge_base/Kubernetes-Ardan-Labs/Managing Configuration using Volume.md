#k8s #volumes #configuration-management
## Configuration of code
- There are many ways we can configure our code:
	- command-line arguments
	- environment variables
	- configuration files
	- configuration servers
		- getting configuration from a database/zookeeper/APIs etc

### Passing configuration to containers
There are many ways to do that:
#### Bake in the custom image
- bake the configuration in the image - probably the worst way
```Dockerfile
# ....
CMD my-app --port=8000 --threads=4
# ...
ENV PORT=8000 CONCURRENCY=4

COPY app.conf /app/app.conf
```
- why a bad idea:
	- we have often have diff configuration in dev vs in prod
	- if we need to deploy same image but multiple times
		- deployment in multiple different data centers
		- if we are making SaaS application, we need dynamic configuration based on customer, and each customer has their own configuration
	- risk to keep secret information in Dockerfiles
#### Using YAML
```sh
$ kubectl create deployment ping --dry-run=client -o yaml > k8s/ping.yaml --image alpine -- ping 127.0.0.1
$ cat k8s/ping.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: ping
  name: ping
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ping
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: ping
    spec:
      containers:
      - command:
        - ping
        - 127.0.0.1
        image: alpine
        name: ping
        resources: {}
status: {}
```
- we can add environment variables to the yaml so that the container can pick them up
```sh
$ code k8s/ping.yaml
$ cat k8s/ping.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: ping
  name: ping
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ping
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: ping
    spec:
      containers:
      - command:
        - ping
        - $TARGET
        image: alpine
        env:
        - name: TARGET
          value: 127.0.0.3
        - name: MOOD
          value: 😎
        name: ping
        resources: {}
status: {}
```
- this will not work and will give `bad address $TARGET error`
- the way command is passed to the container here is equivalent to passing the JSON Array syntax in the docker file where the content of the array is given as it is without any parsing
- thus, we can either choose to use 
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: ping
  name: ping
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ping
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: ping
    spec:
      containers:
      - command:
        - sh
        - -c
        - ping $TARGET
        image: alpine
        env:
        - name: TARGET
          value: 127.0.0.3
        - name: MOOD
          value: 😎
        name: ping
        resources: {}
status: {}
```
- or we can use a special syntax which k8s supports for such use-cases (using $(`VAR_NAME`))
```sh
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: ping
  name: ping
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ping
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: ping
    spec:
      containers:
      - command:
        - ping
        - $(TARGET)
        image: alpine
        env:
        - name: TARGET
          value: 127.0.0.3
        - name: MOOD
          value: 😎
        name: ping
        resources: {}
status: {}
```
### Other ways to pass environment variables
```yaml
spec:
	containers:
	- command:
	  env:
		  - name: MY_POD_NAMESPACE
		    valueFrom: 
			    fieldRef:
				    fieldPath: metadata.namespace
```

## Configuration files and config maps
- we get a bunch of files and put in config map
```sh
$ code target.conf
$ cat target.conf
127.0.0.2
$ code poetry.txt
$ cat poetry.txt
Roses are red 🔴
Violets are blue 🔵
```
- can create config maps from
	- `--from-env-file` - use when we have bunch of key=value pairs
	- `--from-file`
	- `--from-literal`

```sh
$ kubectl create configmap pingconfig \
--from-file=target.conf \
--from-file=poetry.txt \
--from-literal mood=😁
configmap/pingconfig created

$ kubectl get configmaps
NAME               DATA   AGE
kube-root-ca.crt   1      3h28m
pingconfig         3      26s

$ kubectl describe configmap pingconfig
Name:         pingconfig
Namespace:    configdemo
Labels:       <none>
Annotations:  <none>

Data
====
mood:
----
😁

poetry.txt:
----
Roses are red 🔴
Violets are blue 🔵

target.conf:
----
127.0.0.2


BinaryData
====

Events:  <none>
$ kubectl get configmap pingconfig -o yaml
apiVersion: v1
data:
  mood: "\U0001F601"
  poetry.txt: "Roses are red \U0001F534\nViolets are blue \U0001F535"
  target.conf: 127.0.0.2
kind: ConfigMap
metadata:
  creationTimestamp: "2026-01-10T09:51:33Z"
  name: pingconfig
  namespace: configdemo
  resourceVersion: "469941"
  uid: a18393ca-4228-48da-a458-8440b088a84d
```

## Using config maps with volumes
```sh
$ code k8s/ping.yaml
$ cat k8s/ping.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: ping
  name: ping
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ping
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: ping
    spec:
      volumes:
      - name: conf
        configMap:
          name: pingconfig
      containers:
      - command:
        - ping
        - $(TARGET)
        image: alpine
        env:
        - name: TARGET
          value: 127.0.0.3
        - name: MOOD
          value: 😎
        name: ping
        resources: {}
        volumeMounts:
        - name: conf
          mountPath: /conf
status: {}
$ kubectl apply -f k8s/ping.yaml
$ kubectl exec -it ping-77b49b4887-vkk26 -- sh
/ # ls
bin    conf   dev    etc    home   lib    media  mnt    opt    proc   root   run    sbin   srv    sys    tmp    usr    var
/ # ls conf
mood         poetry.txt   target.conf
/ # ping $(cat /conf/target.conf)
PING 127.0.0.2 (127.0.0.2): 56 data bytes
64 bytes from 127.0.0.2: seq=0 ttl=64 time=0.143 ms
64 bytes from 127.0.0.2: seq=1 ttl=64 time=0.133 ms
64 bytes from 127.0.0.2: seq=2 ttl=64 time=0.377 ms
64 bytes from 127.0.0.2: seq=3 ttl=64 time=0.065 ms
^C
--- 127.0.0.2 ping statistics ---
4 packets transmitted, 4 packets received, 0% packet loss
round-trip min/avg/max = 0.065/0.179/0.377 ms
/ # exit
$ code k8s/ping.yaml
$ cat k8s/ping.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: ping
  name: ping
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ping
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: ping
    spec:
      volumes:
      - name: conf
        configMap:
          name: pingconfig
      containers:
      - command:
        - sh
        - -c
        - ping $(cat /conf/target.conf)
        image: alpine
        name: ping
        resources: {}
        volumeMounts:
        - name: conf
          mountPath: /conf
status: {}
$ kubectl exec ping-6d79bbcf76-tngt7 -it -- sh
/ # ls
bin    conf   dev    etc    home   lib    media  mnt    opt    proc   root   run    sbin   srv    sys    tmp    usr    var
/ # ps
PID   USER     TIME  COMMAND
    1 root      0:00 ping 127.0.0.2
    8 root      0:00 sh
   15 root      0:00 ps
/ #
```

### Editing config maps
```sh
$ code target.conf
$ cat target.conf
127.0.2.10
$ kubectl apply -f k8s/ping.yaml
deployment.apps/ping configured
$ kubectl get all
NAME                       READY   STATUS    RESTARTS   AGE
pod/ping-9d5cf448c-7nhzd   1/1     Running   0          111s

NAME                   READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/ping   1/1     1            1           13m

NAME                              DESIRED   CURRENT   READY   AGE
replicaset.apps/ping-6d79bbcf76   0         0         0       3m46s
replicaset.apps/ping-77b49b4887   0         0         0       13m
replicaset.apps/ping-9d5cf448c    1         1         1       111s
```
- on changing the target.conf, we didn't find any change in the pod/container
- this is because, we just changed the config file locally and didn't yet update the config map
```sh
$ kubectl edit configmap pingconfig
configmap/pingconfig edited
$ kubectl get configmap pingconfig -o yaml
apiVersion: v1
data:
  mood: "\U0001F601"
  poetry.txt: "Roses are red \U0001F534\nViolets are blue \U0001F535"
  target.conf: 127.0.2.10
kind: ConfigMap
metadata:
  creationTimestamp: "2026-01-10T09:51:33Z"
  name: pingconfig
  namespace: configdemo
  resourceVersion: "472610"
  uid: a18393ca-4228-48da-a458-8440b088a84d
$ kubectl apply -f k8s/ping.yaml
deployment.apps/ping configured
➜  container.training git:(main) ✗ kubectl get all
NAME                       READY   STATUS    RESTARTS   AGE
pod/ping-9d5cf448c-7nhzd   1/1     Running   0          9m5s

NAME                   READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/ping   1/1     1            1           20m

NAME                              DESIRED   CURRENT   READY   AGE
replicaset.apps/ping-6d79bbcf76   0         0         0       11m
replicaset.apps/ping-77b49b4887   0         0         0       20m
replicaset.apps/ping-9d5cf448c    1         1         1       9m5s
```
- even after changing the `configmap`, there's no change or rolling update in the pod/container
- although the config file in the container should have changed
```sh
$ kubectl exec -it pod/ping-9d5cf448c-7nhzd -- sh
/ # ls
bin    conf   dev    etc    home   lib    media  mnt    opt    proc   root   run    sbin   srv    sys    tmp    usr    var
/ # cat conf/target.conf
127.0.2.10/ #
/ # ls -al conf
total 12
drwxrwxrwx    3 root     root          4096 Jan 10 10:24 .
drwxr-xr-x    1 root     root          4096 Jan 10 10:16 ..
drwxr-xr-x    2 root     root          4096 Jan 10 10:24 ..2026_01_10_10_24_18.2830330291
lrwxrwxrwx    1 root     root            32 Jan 10 10:24 ..data -> ..2026_01_10_10_24_18.2830330291
lrwxrwxrwx    1 root     root            11 Jan 10 10:16 mood -> ..data/mood
lrwxrwxrwx    1 root     root            17 Jan 10 10:16 poetry.txt -> ..data/poetry.txt
lrwxrwxrwx    1 root     root            18 Jan 10 10:16 target.conf -> ..data/target.conf
/ #
```
- but our program is still is pinging 127.0.0.2
```sh
/ # ps
PID   USER     TIME  COMMAND
    1 root      0:00 ping 127.0.0.2
    8 root      0:00 sh
   19 root      0:00 ps
```
- this is because, our program only reads the configuration file just when it starts and keeps running with that configuration
- the updates in the config maps are going to **take a minute to land up and propagate** into the volume in the pods/containers
	- technically what happens is that `kubelet` is checking the changes in config-maps regularly
	- when the next time that it's going to check it, it will realize
- even after the updates in config maps propagate the container, we need to do something in our program so that it can reflect on the changes in the config file
	- we can either use `kill` command to send signal to the process running in the container
	- or rollout restart the pod/deployment
		- `$ kubectl rollout restart deployment ping`
```sh
$ kubectl rollout restart deployment ping
deployment.apps/ping restarted
$ kubectl get pods
NAME                    READY   STATUS        RESTARTS   AGE
ping-5669d5bccb-zlp9q   1/1     Running       0          7s
ping-9d5cf448c-7nhzd    1/1     Terminating   0          21m
$ kubectl exec -it ping-5669d5bccb-zlp9q -- sh
/ # ps
PID   USER     TIME  COMMAND
    1 root      0:00 ping 127.0.2.10
    8 root      0:00 sh
   14 root      0:00 ps
/ #
```

### Automating rollout restart on changes in configuration
- on hitting the `kubectl rollout restart deployment`, all it does is updates the `template.metadata.annotations.kubectl.kubernetes.io/restartedAt` timestamp, which thus changes the pod template, causing the deployment to do the rolling update
```sh
$ kubectl get deployment ping -o yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    deployment.kubernetes.io/revision: "4"
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"apps/v1","kind":"Deployment","metadata":{"annotations":{},"creationTimestamp":null,"labels":{"app":"ping"},"name":"ping","namespace":"configdemo"},"spec":{"replicas":1,"selector":{"matchLabels":{"app":"ping"}},"strategy":{},"template":{"metadata":{"creationTimestamp":null,"labels":{"app":"ping"}},"spec":{"containers":[{"command":["sh","-c","ping $(cat /conf/target.conf)"],"image":"alpine","name":"ping","resources":{},"volumeMounts":[{"mountPath":"/conf","name":"conf"}]}],"volumes":[{"configMap":{"name":"pingconfig"},"name":"conf"}]}}},"status":{}}
  creationTimestamp: "2026-01-10T10:04:19Z"
  generation: 4
  labels:
    app: ping
  name: ping
  namespace: configdemo
  resourceVersion: "473705"
  uid: f3d795a7-aca8-4a54-8559-a5dacfd7a37d
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10

  selector:
    matchLabels:
      app: ping
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      annotations:
        kubectl.kubernetes.io/restartedAt: "2026-01-10T16:07:21+05:30" # changing this timestamp triggers restart/rollout
      creationTimestamp: null
      labels:
        app: ping
    spec:
      containers:
      - command:
        - sh
        - -c
        - ping $(cat /conf/target.conf)
        image: alpine
        imagePullPolicy: Always
        name: ping
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /conf
          name: conf
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
      volumes:
      - configMap:
          defaultMode: 420
          name: pingconfig
        name: conf
status:
  availableReplicas: 1
  conditions:
  - lastTransitionTime: "2026-01-10T10:04:25Z"
    lastUpdateTime: "2026-01-10T10:04:25Z"
    message: Deployment has minimum availability.
    reason: MinimumReplicasAvailable
    status: "True"
    type: Available
  - lastTransitionTime: "2026-01-10T10:04:19Z"
    lastUpdateTime: "2026-01-10T10:37:28Z"
    message: ReplicaSet "ping-5669d5bccb" has successfully progressed.
    reason: NewReplicaSetAvailable
    status: "True"
    type: Progressing
  observedGeneration: 4
  readyReplicas: 1
  replicas: 1
  updatedReplicas: 1
```

### Automating Workflow
```sh
$ code ping.sh
$ kubectl create configmap pingconfig -o yaml --dry-run=client \
--from-file target.conf \
--from-file poetry.txt \
--from-literal mood=😁
apiVersion: v1
data:
  mood: "\U0001F601"
  poetry.txt: "Roses are red \U0001F534\nViolets are blue \U0001F535\n"
  target.conf: |
    127.0.2.10
kind: ConfigMap
metadata:
  creationTimestamp: null
  name: pingconfig
# can use that and put in ping.sh
$ kubectl create configmap pingconfig -o yaml --dry-run=client \
--from-file target.conf \
--from-file poetry.txt \
--from-literal mood=😁 >> k8s/ping.yaml # a single > overwrites, a double >> appends to existing file
$ cat k8s/ping.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: ping
  name: ping
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ping
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: ping
    spec:
      volumes:
      - name: conf
        configMap:
          name: pingconfig
      containers:
      - command:
        - sh
        - -c
        - ping $(cat /conf/target.conf)
        image: alpine
        name: ping
        resources: {}
        volumeMounts:
        - name: conf
          mountPath: /conf
status: {}
apiVersion: v1
data:
  mood: "\U0001F601"
  poetry.txt: "Roses are red \U0001F534\nViolets are blue \U0001F535\n"
  target.conf: |
    127.0.2.10
kind: ConfigMap
metadata:
  creationTimestamp: null
  name: pingconfig
# now, it's pretty important to put --- between the deployment and configmap specs
$ code k8s/ping.yaml
$ cat k8s/ping.yaml
cat k8s/ping.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: ping
  name: ping
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ping
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: ping
    spec:
      volumes:
      - name: conf
        configMap:
          name: pingconfig
      containers:
      - command:
        - sh
        - -c
        - ping $(cat /conf/target.conf)
        image: alpine
        name: ping
        resources: {}
        volumeMounts:
        - name: conf
          mountPath: /conf
status: {}

---

apiVersion: v1
data:
  mood: "\U0001F601"
  poetry.txt: "Roses are red \U0001F534\nViolets are blue \U0001F535\n"
  target.conf: |
    127.0.2.10
kind: ConfigMap
metadata:
  name: pingconfig
```
- now, even if we apply this, it won't trigger a rolling update
```sh
$ kubectl apply -f k8s/ping.yaml
deployment.apps/ping configured
Warning: resource configmaps/pingconfig is missing the kubectl.kubernetes.io/last-applied-configuration annotation which is required by kubectl apply. kubectl apply should only be used on resources created declaratively by either kubectl create --save-config or kubectl apply. The missing annotation will be patched automatically.
configmap/pingconfig configured
$ kubectl get all
NAME                        READY   STATUS    RESTARTS   AGE
pod/ping-5669d5bccb-zlp9q   1/1     Running   0          53m

NAME                   READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/ping   1/1     1            1           86m

NAME                              DESIRED   CURRENT   READY   AGE
replicaset.apps/ping-5669d5bccb   1         1         1       53m
replicaset.apps/ping-6d79bbcf76   0         0         0       76m
replicaset.apps/ping-77b49b4887   0         0         0       86m
replicaset.apps/ping-9d5cf448c    0         0         0       74m
```
- thus, we can do this
```sh
$ code deploy.sh
$ cat deploy.sh
#!/bin/sh
kubectl apply -f k8s/ping.yaml
kubectl rollout restart deployment ping
$ chmod +x deploy.sh
$ ./deploy.sh
deployment.apps/ping configured
configmap/pingconfig configured
deployment.apps/ping restarted
```
- one problem is every time we run `deploy.sh`, it will do the rolling update even if there are was no real change in the config file
### using config map hash
- take the hash of the `configmap` and use it as the annotation in the deployment pod template
- thus, if the config map doesn't change, the hash doesn't change and we don't need to do rolling update
```sh
$ code deploy.sh
$ kubectl get configmap pingconfig -o yaml | sha256sum | cut -c 1-64
9bf42a430fd585214962b0deede064783f4817df18cca8876218d33ea0f6f051
$ code deploy.sh
$ kubectl patch deployment ping --patch "
spec:
	template:
		metadata:
			annotations:
				configmapHash: $HASH
"
deployment.apps/ping patched
$ kubectl get deployment ping -o yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    deployment.kubernetes.io/revision: "6"
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"apps/v1","kind":"Deployment","metadata":{"annotations":{},"creationTimestamp":null,"labels":{"app":"ping"},"name":"ping","namespace":"configdemo"},"spec":{"replicas":1,"selector":{"matchLabels":{"app":"ping"}},"strategy":{},"template":{"metadata":{"creationTimestamp":null,"labels":{"app":"ping"}},"spec":{"containers":[{"command":["sh","-c","ping $(cat /conf/target.conf)"],"image":"alpine","name":"ping","resources":{},"volumeMounts":[{"mountPath":"/conf","name":"conf"}]}],"volumes":[{"configMap":{"name":"pingconfig"},"name":"conf"}]}}},"status":{}}
  creationTimestamp: "2026-01-10T10:04:19Z"
  generation: 6
  labels:
    app: ping
  name: ping
  namespace: configdemo
  resourceVersion: "479089"
  uid: f3d795a7-aca8-4a54-8559-a5dacfd7a37d
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10

  selector:
    matchLabels:
      app: ping
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      annotations:
        configmapHash: 9bf42a430fd585214962b0deede064783f4817df18cca8876218d33ea0f6f051
        kubectl.kubernetes.io/restartedAt: "2026-01-10T17:05:21+05:30"
      creationTimestamp: null
      labels:
        app: ping
    spec:
      containers:
      - command:
        - sh
        - -c
        - ping $(cat /conf/target.conf)
        image: alpine
        imagePullPolicy: Always
        name: ping
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /conf
          name: conf
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
      volumes:
      - configMap:
          defaultMode: 420
          name: pingconfig
        name: conf
status:
  availableReplicas: 1
  conditions:
  - lastTransitionTime: "2026-01-10T10:04:25Z"
    lastUpdateTime: "2026-01-10T10:04:25Z"
    message: Deployment has minimum availability.
    reason: MinimumReplicasAvailable
    status: "True"
    type: Available
  - lastTransitionTime: "2026-01-10T10:04:19Z"
    lastUpdateTime: "2026-01-10T11:56:18Z"
    message: ReplicaSet "ping-5bb845ddf" has successfully progressed.
    reason: NewReplicaSetAvailable
    status: "True"
    type: Progressing
  observedGeneration: 6
  readyReplicas: 1
  replicas: 1
  updatedReplicas: 1

$ code deploy.sh
$ cat deploy.sh
#!/bin/sh
kubectl apply -f k8s/ping.yaml
HASH=$(kubectl get configmap pingconfig -o yaml | sha256sum | cut -c 1-64)
kubectl patch deployment ping --patch "
spec:
    template:
        metadata:
            annotations:
                configmapHash: $HASH
"
# kubectl rollout restart deployment ping
$ ./deploy.sh
deployment.apps/ping configured
configmap/pingconfig configured
deployment.apps/ping patched (no change)

# change in the ping.yaml
$ code ping.yaml
$ cat ping.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: ping
  name: ping
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ping
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: ping
    spec:
      volumes:
      - name: conf
        configMap:
          name: pingconfig
      containers:
      - command:
        - sh
        - -c
        - ping $(cat /conf/target.conf)
        image: alpine
        name: ping
        resources: {}
        volumeMounts:
        - name: conf
          mountPath: /conf
status: {}

---

apiVersion: v1
data:
  mood: "\U0001F601"
  poetry.txt: "Roses are red \U0001F534\nViolets are blue \U0001F535\n"
  target.conf: |
    127.0.10.10
kind: ConfigMap
metadata:
  creationTimestamp: null
  name: pingconfig

$ ./deploy.sh
deployment.apps/ping configured
configmap/pingconfig configured
deployment.apps/ping patched

$ kubectl get all
NAME                        READY   STATUS    RESTARTS   AGE
pod/ping-67bb8fb46b-jrmh2   1/1     Running   0          5m38s

NAME                   READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/ping   1/1     1            1           173m

NAME                              DESIRED   CURRENT   READY   AGE
replicaset.apps/ping-5669d5bccb   0         0         0       140m
replicaset.apps/ping-56b8d89f49   0         0         0       82m
replicaset.apps/ping-5bb845ddf    0         0         0       62m
replicaset.apps/ping-67bb8fb46b   1         1         1       5m39s
replicaset.apps/ping-6d79bbcf76   0         0         0       164m
replicaset.apps/ping-77b49b4887   0         0         0       173m
replicaset.apps/ping-9d5cf448c    0         0         0       162m 
```

### Propagating updated config without rolling update
- create rainbow namespace and run blue and green containers
```sh
$ kubectl create namespace rainbow
$ kubectl config set-context --namespace rainbow --current
$ kubectl config get-contexts
CURRENT   NAME             CLUSTER          AUTHINFO         NAMESPACE
*         docker-desktop   docker-desktop   docker-desktop   rainbow
```
- see `haproxy.cfg` and `haproxy.yaml`
```sh
$ cat k8s/haproxy.cfg
global
  daemon

defaults
  mode tcp
  timeout connect 5s
  timeout client 50s
  timeout server 50s

listen very-basic-load-balancer
  bind *:80
  server blue blue:80
  server green green:80

$ cat k8s/haproxy.yaml
apiVersion: v1
kind: Pod
metadata:
  name: haproxy
spec:
  volumes:
  - name: config
    configMap:
      name: haproxy
  containers:
  - name: haproxy
    image: haproxy:1
    volumeMounts:
    - name: config
      mountPath: /usr/local/etc/haproxy/
```
- first step, for applying the `haproxy.yaml` file, we need to create the `haproxy` config-map
```sh
$ kubectl create configmap haproxy -o yaml --dry-run=client --from-file k8s/haproxy.cfg
apiVersion: v1
data:
  haproxy.cfg: |
    global
      daemon

    defaults
      mode tcp
      timeout connect 5s
      timeout client 50s
      timeout server 50s

    listen very-basic-load-balancer
      bind *:80
      server blue blue:80
      server green green:80

    # Note: the services above must exist,
    # otherwise HAproxy won't start.
kind: ConfigMap
metadata:
  creationTimestamp: null
  name: haproxy
$ kubectl create configmap haproxy --from-file k8s/haproxy.cfg
$ kubectl get configmaps
NAME               DATA   AGE
haproxy            1      11s
kube-root-ca.crt   1      92m
$ kubectl apply -f k8s/haproxy.yaml
$ kubectl get pods
NAME                     READY   STATUS    RESTARTS   AGE
blue-5c986bd7bf-8927k    1/1     Running   0          121m
green-556754bb7d-k4zmk   1/1     Running   0          121m
haproxy                  1/1     Running   0          56s
red-77f6d65f98-jmmnf     1/1     Running   0          121m
$ kubectl port-forward pod/haproxy 8080:80
Forwarding from 127.0.0.1:8080 -> 80
Forwarding from [::1]:8080 -> 80
Handling connection for 8080
Handling connection for 8080
Handling connection for 8080

...
$ curl localhost:8080
🔵This is pod rainbow/blue-5c986bd7bf-8927k on linux/amd64, serving / for 10.1.1.205:46892.
$ curl localhost:8080
🟢This is pod rainbow/green-556754bb7d-k4zmk on linux/amd64, serving / for 10.1.1.205:38246.
$ curl localhost:8080
🔵This is pod rainbow/blue-5c986bd7bf-8927k on linux/amd64, serving / for 10.1.1.205:38076.
$ curl localhost:8080
🟢This is pod rainbow/green-556754bb7d-k4zmk on linux/amd64, serving / for 10.1.1.205:38262.
$ curl localhost:8080
🔵This is pod rainbow/blue-5c986bd7bf-8927k on linux/amd64, serving / for 10.1.1.205:38086.
```
- the `haproxy` does a perfect round-robin load balancing between blue and green
- now, if we want to change that configuration ( maybe to do blue-green-red instead of blue-green )
```sh
$ kubectl edit configmap haproxy
configmap/haproxy edited
$ kubectl get configmap haproxy -o yaml
apiVersion: v1
data:
  haproxy.cfg: |
    global
      daemon

    defaults
      mode tcp
      timeout connect 5s
      timeout client 50s
      timeout server 50s

    listen very-basic-load-balancer
      bind *:80
      server blue blue:80
      server green green:80
      server red red:80

    # Note: the services above must exist,
    # otherwise HAproxy won't start.
kind: ConfigMap
metadata:
  creationTimestamp: "2026-01-10T15:20:03Z"
  name: haproxy
  namespace: rainbow
  resourceVersion: "491501"
  uid: 06cea999-252e-4b9f-9d75-19ace2adb1fc
```
- after editing the config map, we still find the same load balancing
```sh
$ curl localhost:8080
🟢This is pod rainbow/green-556754bb7d-k4zmk on linux/amd64, serving / for 10.1.1.205:45896.
$ curl localhost:8080
🔵This is pod rainbow/blue-5c986bd7bf-8927k on linux/amd64, serving / for 10.1.1.205:45294.
$ curl localhost:8080
🟢This is pod rainbow/green-556754bb7d-k4zmk on linux/amd64, serving / for 10.1.1.205:45898.
$ curl localhost:8080
🔵This is pod rainbow/blue-5c986bd7bf-8927k on linux/amd64, serving / for 10.1.1.205:45306
```
- although, if we see inside the HAProxy container, we can find the updated `haproxy.cfg`
```sh
$ kubectl exec -it haproxy -- bash
root@haproxy:/#
root@haproxy:/# ls
bin   docker-entrypoint.sh  lib    mnt	 root  srv  usr
boot  etc		    lib64  opt	 run   sys  var
dev   home		    media  proc  sbin  tmp
root@haproxy:/#
root@haproxy:/# ps
bash: ps: command not found
root@haproxy:/#
root@haproxy:/# cd /usr/local/
bin/     games/   lib/     sbin/    src/
etc/     include/ man/     share/
root@haproxy:/# cd /usr/local/etc/haproxy/
root@haproxy:/usr/local/etc/haproxy# ls
haproxy.cfg
root@haproxy:/usr/local/etc/haproxy# cat haproxy.cfg
global
  daemon

defaults
  mode tcp
  timeout connect 5s
  timeout client 50s
  timeout server 50s

listen very-basic-load-balancer
  bind *:80
  server blue blue:80
  server green green:80
  server red red:80

# Note: the services above must exist,
# otherwise HAproxy won't start.
```
- thus, we need to tell HAProxy to reload it's configuration
- we can send signals to HAProxy to reload it's configuration
	- specific to HAProxy, we can send `HUP` or `USR1` signal
	- to send a signal on a normal VM or server, we can do that by doing `systemctl reload` which eventually sends a signal to HAProxy
	- to send a signal to HAProxy, we need to know the PID of HAProxy to be able to pass to `kill`
	- `kill -HUP <PID>`
```sh
$ cd /proc
root@haproxy:/proc# ls
1	   crypto	irq	     misc	   swaps
29	   devices	kallsyms     modules	   sys
40	   diskstats	kcore	     mounts	   sysrq-trigger
7	   dma		key-users    mtrr	   sysvipc
acpi	   docker	keys	     net	   thread-self
buddyinfo  driver	kmsg	     pagetypeinfo  timer_list
bus	   execdomains	kpagecgroup  partitions    tty
cgroups    filesystems	kpagecount   pressure	   uptime
cmdline    fs		kpageflags   self	   version
config.gz  interrupts	loadavg      slabinfo	   vmallocinfo
consoles   iomem	locks	     softirqs	   vmstat
cpuinfo    ioports	meminfo      stat	   zoneinfo
root@haproxy:/proc#
root@haproxy:/proc# ls -l
total 0
dr-xr-xr-x  9 root root     0 Jan 10 15:20 1
dr-xr-xr-x  9 root root     0 Jan 10 16:34 29
dr-xr-xr-x  9 root root     0 Jan 10 16:49 41
dr-xr-xr-x  9 root root     0 Jan 10 16:46 7
drwxrwxrwt  2 root root    40 Jan 10 15:20 acpi
-r--r--r--  1 root root     0 Jan 10 16:49 buddyinfo
dr-xr-xr-x  4 root root     0 Jan 10 15:20 bus
-r--r--r--  1 root root     0 Jan 10 16:49 cgroups
-r--r--r--  1 root root   275 Jan 10 16:49 cmdline
-r--r--r--  1 root root 28049 Jan 10 16:49 config.gz
-r--r--r--  1 root root     0 Jan 10 16:49 consoles
-r--r--r--  1 root root     0 Jan 10 16:49 cpuinfo
-r--r--r--  1 root root     0 Jan 10 16:49 crypto
-r--r--r--  1 root root     0 Jan 10 16:49 devices
-r--r--r--  1 root root     0 Jan 10 16:49 diskstats
-r--r--r--  1 root root     0 Jan 10 16:49 dma
drwxr-xr-x  7 root root     0 Jan 10 16:49 docker
dr-xr-xr-x  4 root root     0 Jan 10 16:49 driver
-r--r--r--  1 root root     0 Jan 10 16:49 execdomains
-r--r--r--  1 root root     0 Jan 10 16:32 filesystems
dr-xr-xr-x 14 root root     0 Jan 10 15:20 fs
-r--r--r--  1 root root     0 Jan 10 16:49 interrupts
-r--r--r--  1 root root     0 Jan 10 16:49 iomem
-r--r--r--  1 root root     0 Jan 10 16:49 ioports
dr-xr-xr-x 53 root root     0 Jan 10 15:20 irq
-r--r--r--  1 root root     0 Jan 10 16:49 kallsyms
crw-rw-rw-  1 root root  1, 3 Jan 10 15:20 kcore
-r--r--r--  1 root root     0 Jan 10 16:49 key-users
crw-rw-rw-  1 root root  1, 3 Jan 10 15:20 keys
-r--------  1 root root     0 Jan 10 16:49 kmsg
-r--------  1 root root     0 Jan 10 16:49 kpagecgroup
-r--------  1 root root     0 Jan 10 16:49 kpagecount
-r--------  1 root root     0 Jan 10 16:49 kpageflags
-r--r--r--  1 root root     0 Jan 10 16:49 loadavg
-r--r--r--  1 root root     0 Jan 10 16:49 locks
-r--r--r--  1 root root     0 Jan 10 16:49 meminfo
-r--r--r--  1 root root     0 Jan 10 16:49 misc
-r--r--r--  1 root root     0 Jan 10 16:49 modules
lrwxrwxrwx  1 root root    11 Jan 10 16:49 mounts -> self/mounts
-rw-r--r--  1 root root     0 Jan 10 16:49 mtrr
lrwxrwxrwx  1 root root     8 Jan 10 16:49 net -> self/net
-r--------  1 root root     0 Jan 10 16:49 pagetypeinfo
-r--r--r--  1 root root     0 Jan 10 16:49 partitions
dr-xr-xr-x  5 root root     0 Jan 10 16:49 pressure
lrwxrwxrwx  1 root root     0 Jan 10 15:20 self -> 41
-r--------  1 root root     0 Jan 10 16:49 slabinfo
-r--r--r--  1 root root     0 Jan 10 16:49 softirqs
-r--r--r--  1 root root     0 Jan 10 16:49 stat
-r--r--r--  1 root root     0 Jan 10 16:49 swaps
dr-xr-xr-x  1 root root     0 Jan 10 15:20 sys
--w-------  1 root root     0 Jan 10 15:20 sysrq-trigger
dr-xr-xr-x  5 root root     0 Jan 10 16:49 sysvipc
lrwxrwxrwx  1 root root     0 Jan 10 15:20 thread-self -> 41/task/41
crw-rw-rw-  1 root root  1, 3 Jan 10 15:20 timer_list
dr-xr-xr-x  6 root root     0 Jan 10 16:49 tty
-r--r--r--  1 root root     0 Jan 10 16:49 uptime
-r--r--r--  1 root root     0 Jan 10 16:49 version
-r--------  1 root root     0 Jan 10 16:49 vmallocinfo
-r--r--r--  1 root root     0 Jan 10 16:49 vmstat
-r--r--r--  1 root root     0 Jan 10 16:49 zoneinfo
root@haproxy:/proc#
root@haproxy:/proc# ls -l */exe
lrwxrwxrwx 1 root root 0 Jan 10 16:49 1/exe -> /usr/local/sbin/haproxy
lrwxrwxrwx 1 root root 0 Jan 10 16:34 29/exe -> /bin/bash
lrwxrwxrwx 1 root root 0 Jan 10 16:49 7/exe -> /usr/local/sbin/haproxy
lrwxrwxrwx 1 root root 0 Jan 10 16:49 self/exe -> /bin/ls
lrwxrwxrwx 1 root root 0 Jan 10 16:49 thread-self/exe -> /bin/ls
root@haproxy:/proc# kill -h
bash: kill: h: invalid signal specification
root@haproxy:/proc# kill --help
kill: kill [-s sigspec | -n signum | -sigspec] pid | jobspec ... or kill -l [sigspec]
    Send a signal to a job.

    Send the processes identified by PID or JOBSPEC the signal named by
    SIGSPEC or SIGNUM.  If neither SIGSPEC nor SIGNUM is present, then
    SIGTERM is assumed.

    Options:
      -s sig	SIG is a signal name
      -n sig	SIG is a signal number
      -l	list the signal names; if arguments follow `-l` they are
    		assumed to be signal numbers for which names should be listed
      -L	synonym for -l

    Kill is a shell builtin for two reasons: it allows job IDs to be used
    instead of process IDs, and allows processes to be killed if the limit
    on processes that you can create is reached.

    Exit Status:
    Returns success unless an invalid option is given or an error occurs.
root@haproxy:/proc# kill -HUP 1
root@haproxy:/proc# exit
$ curl localhost:8080
🔵This is pod rainbow/blue-5c986bd7bf-8927k on linux/amd64, serving / for 10.1.1.205:52788.
$ curl localhost:8080
🟢This is pod rainbow/green-556754bb7d-k4zmk on linux/amd64, serving / for 10.1.1.205:54706.
$ curl localhost:8080
🔵This is pod rainbow/blue-5c986bd7bf-8927k on linux/amd64, serving / for 10.1.1.205:40296.
$ curl localhost:8080
🟢This is pod rainbow/green-556754bb7d-k4zmk on linux/amd64, serving / for 10.1.1.205:40234.
$ curl localhost:8080
🔴This is pod rainbow/red-77f6d65f98-jmmnf on linux/amd64, serving / for 10.1.1.205:59374.
$ curl localhost:8080
🔵This is pod rainbow/blue-5c986bd7bf-8927k on linux/amd64, serving / for 10.1.1.205:36200.
$ curl localhost:8080
🟢This is pod rainbow/green-556754bb7d-k4zmk on linux/amd64, serving / for 10.1.1.205:42664.
$ curl localhost:8080
🔴This is pod rainbow/red-77f6d65f98-jmmnf on linux/amd64, serving / for 10.1.1.205:58700.
$ curl localhost:8080
🔵This is pod rainbow/blue-5c986bd7bf-8927k on linux/amd64, serving / for 10.1.1.205:36214.
$ curl localhost:8080
🟢This is pod rainbow/green-556754bb7d-k4zmk on linux/amd64, serving / for 10.1.1.205:42676.
$ curl localhost:8080
🔴This is pod rainbow/red-77f6d65f98-jmmnf on linux/amd64, serving / for 10.1.1.205:58710.
```
- thus, this way, we could reload HAProxy, without having to destroy the container
- to be able to include this workflow in a CI/CD pipeline, we can have two containers in same pod (side-car pattern #side-car )
- the side-car container's goal will be to reload HAProxy configuration
- this pattern is implemented by many applications like
	- **prometheus** if installed on k8s, most of the prometheus deployments in the prometheus pod, we have two containers
		- prometheus itself
		- another is config-reloader which doesn't use signal but an API call
```sh
$ code k8s/haproxy.yaml
apiVersion: v1
kind: Pod
metadata:
  name: haproxy
spec:
  volumes:
  - name: config
    configMap:
      name: haproxy
  containers:
  - name: haproxy
    image: haproxy:1
    volumeMounts:
    - name: config
      mountPath: /usr/local/etc/haproxy/
  - name: config-reloader # side-car container
    image: alpine
    volumeMounts:
    - name: config
      mountPath: /conf
    command:
      - sh
      - -c
      - |
        cp /conf/haproxy.cfg /tmp/haproxy.cfg
        while sleep 1; do
          if ! diff /tmp/haproxy.cfg /conf/haproxy.cfg; then 
            echo "Reloading HAProxy."
            # FIXME
            cp /conf/haproxy.cfg /tmp/haproxy.cfg
          fi
        done

$ kubectl apply -f k8s/haproxy.yaml
The Pod "haproxy" is invalid: spec.containers: Forbidden: pod updates may not add or remove containers
$ # since we don't have a deployment and replica-sets backing the pod, we need to force applying the specs
$ kubectl apply -f k8s/haproxy.yaml -
-force
error: error when applying patch:
{"apiVersion":"v1","kind":"Pod","metadata":{"annotations":{"kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"v1\",\"kind\":\"Pod\",\"metadata\":{\"annotations\":{},\"name\":\"haproxy\",\"namespace\":\"rainbow\"},\"spec\":{\"containers\":[{\"image\":\"haproxy:1\",\"name\":\"haproxy\",\"volumeMounts\":[{\"mountPath\":\"/usr/local/etc/haproxy/\",\"name\":\"config\"}]},{\"command\":[\"sh\",\"-c\",\"cp /conf/haproxy.conf /tmp/haproxy.conf\\nwhile sleep 1; do\\n  if ! diff /tmp/haproxy.conf /conf/haproxy.conf; then \\n    echo \\\"Reloading HAProxy.\\\"\\n    # FIXME\\n    cp /conf/haproxy.conf /tmp/haproxy.conf\\n  fi\\ndone\\n\"],\"image\":\"alpine\",\"name\":\"config-reloader\",\"volumeMounts\":[{\"mountPath\":\"/conf\",\"name\":\"config\"}]}],\"volumes\":[{\"configMap\":{\"name\":\"haproxy\"},\"name\":\"config\"}]}}\n"},"name":"haproxy","namespace":"rainbow"},"spec":{"containers":[{"image":"haproxy:1","name":"haproxy","volumeMounts":[{"mountPath":"/usr/local/etc/haproxy/","name":"config"}]},{"command":["sh","-c","cp /conf/haproxy.conf /tmp/haproxy.conf\nwhile sleep 1; do\n  if ! diff /tmp/haproxy.conf /conf/haproxy.conf; then \n    echo \"Reloading HAProxy.\"\n    # FIXME\n    cp /conf/haproxy.conf /tmp/haproxy.conf\n  fi\ndone\n"],"image":"alpine","name":"config-reloader","volumeMounts":[{"mountPath":"/conf","name":"config"}]}],"volumes":[{"configMap":{"name":"haproxy"},"name":"config"}]}}

to:
Resource: "/v1, Resource=pods", GroupVersionKind: "/v1, Kind=Pod"
Name: "haproxy", Namespace: "rainbow"
for: "k8s/haproxy.yaml": context deadline exceeded
$ kubectl get pods
NAME                     READY   STATUS    RESTARTS   AGE
blue-5c986bd7bf-8927k    1/1     Running   0          3h52m
green-556754bb7d-k4zmk   1/1     Running   0          3h52m
red-77f6d65f98-jmmnf     1/1     Running   0          3h52m
$ kubectl apply -f k8s/haproxy.yaml --force
pod/haproxy created
$ kubectl get pods
NAME                     READY   STATUS    RESTARTS   AGE
blue-5c986bd7bf-8927k    1/1     Running   0          3h52m
green-556754bb7d-k4zmk   1/1     Running   0          3h52m
haproxy                  2/2     Running   0          11s
red-77f6d65f98-jmmnf     1/1     Running   0          3h52m
$ kubectl exec -it pod/haproxy -- bash
Defaulted container "haproxy" out of: haproxy, config-reloader

^[[Aroot@haproxy:/#
root@haproxy:/# exit
exit
➜  container.training git:(main) ✗ kubectl exec -it pod/haproxy --container=config-reloader -- sh
/ # ps
PID   USER     TIME  COMMAND
    1 root      0:00 sh -c cp /conf/haproxy.conf /tmp/haproxy.conf wh
  841 root      0:00 sh
  912 root      0:00 sleep 1
  913 root      0:00 ps
/ # ps
PID   USER     TIME  COMMAND
    1 root      0:00 sh -c cp /conf/haproxy.conf /tmp/haproxy.conf while sleep 1; do   if ! diff /t
  841 root      0:00 sh
  937 root      0:00 sleep 1
  938 root      0:00 ps
/ # exit
$ kubectl create deployment yellow --image jpetazzo/color
deployment.apps/yellow created
$ kubectl expose deployment yellow --port 80
service/yellow exposed
$ kubectl edit configmap haproxy
$ kubectl get configmap haproxy -o yaml
apiVersion: v1
data:
  haproxy.cfg: |
    global
      daemon

    defaults
      mode tcp
      timeout connect 5s
      timeout client 50s
      timeout server 50s

    listen very-basic-load-balancer
      bind *:80
      server yellow yellow:80
      server green green:80
      server red red:80

    # Note: the services above must exist,
    # otherwise HAproxy won't start.
kind: ConfigMap
metadata:
  creationTimestamp: "2026-01-10T15:20:03Z"
  name: haproxy
  namespace: rainbow
  resourceVersion: "502811"
  uid: 06cea999-252e-4b9f-9d75-19ace2adb1fc
$ kubectl get pods
NAME                      READY   STATUS    RESTARTS   AGE
blue-5c986bd7bf-tgj97     1/1     Running   0          85s
green-556754bb7d-k4zmk    1/1     Running   0          4h40m
haproxy                   2/2     Running   0          2m58s
red-77f6d65f98-jmmnf      1/1     Running   0          4h40m
yellow-69d94b589c-rbrqk   1/1     Running   0          33m
$ kubectl logs haproxy -c config-reloader
--- /tmp/haproxy.cfg
+++ /conf/haproxy.cfg
@@ -9,7 +9,7 @@

 listen very-basic-load-balancer
   bind *:80
-  server blue blue:80
+  server yellow yellow:80
   server green green:80
   server red red:80

Reloading HAProxy.
```
- now, in the `#FIXME` we need a way to reach the HAProxy container and send a HUF signal to the `haproxy` process (PID=1)
- to be able to see `haproxy` container from within this `config-reloader` container, we can use shared namespace pid #shared-pid-namespace
- we will tell k8s to put these containers together
- by default each container gets it's own pid namespace, i.e. each container has it's own PID=1, but we are going to convey to k8s to put the containers together
```sh
$ kubectl explain pod.spec.shareProcessNamespace
KIND:       Pod
VERSION:    v1

FIELD: shareProcessNamespace <boolean>


DESCRIPTION:
    Share a single process namespace between all of the containers in a pod.
    When this is set containers will be able to view and signal processes from
    other containers in the same pod, and the first process in each container
    will not be assigned PID 1. HostPID and ShareProcessNamespace cannot both be
    set. Optional: Default to false.
```
> `pod.spec.hostPID <boolean>` can be used to share host's pid namespace to be able to share processes and send signals to across different pods in the same node
- thus, update the `haproxy.yaml`
```sh
$ code k8s/haproxy.yaml
$ cat k8s/haproxy.yaml
apiVersion: v1
kind: Pod
metadata:
  name: haproxy
spec:
  shareProcessNamespace: true
  volumes:
  - name: config
    configMap:
      name: haproxy
  containers:
  - name: haproxy
    image: haproxy:1
    volumeMounts:
    - name: config
      mountPath: /usr/local/etc/haproxy/
  - name: config-reloader # side-car container
    image: alpine
    volumeMounts:
    - name: config
      mountPath: /conf
    command:
      - sh
      - -c
      - |
        cp /conf/haproxy.cfg /tmp/haproxy.cfg
        while sleep 1; do
          if ! diff /tmp/haproxy.cfg /conf/haproxy.cfg; then
            echo "Reloading HAProxy."
            # FIXME
            cp /conf/haproxy.cfg /tmp/haproxy.cfg
          fi
        done%
$ kubectl delete pod haproxy
pod "haproxy" deleted
$ kubectl apply -f k8s/haproxy.yaml
pod/haproxy created
$ kubectl exec -it haproxy -c config-reloader -- sh
/ #  ps faux
PID   USER     TIME  COMMAND
    1 65535     0:00 /pause
    7 root      0:00 haproxy -W -db -f /usr/local/etc/haproxy/haproxy.cfg
   13 root      0:00 haproxy -W -db -f /usr/local/etc/haproxy/haproxy.cfg
   14 root      0:00 sh -c cp /conf/haproxy.cfg /tmp/haproxy.cfg while sleep 1; do   if ! diff /tmp/haproxy.cfg /conf/haprox
  262 root      0:00 sh
  294 root      0:00 sleep 1
  295 root      0:00 ps faux
/ # exit
$ # we need to send HUF signal to all the haproxy processes
$ code k8s/haproxy.yaml
$ cat k8s/haproxy.yaml$
apiVersion: v1
kind: Pod
metadata:
  name: haproxy
spec:
  shareProcessNamespace: true
  volumes:
  - name: config
    configMap:
      name: haproxy
  containers:
  - name: haproxy
    image: haproxy:1
    volumeMounts:
    - name: config
      mountPath: /usr/local/etc/haproxy/
  - name: config-reloader # side-car container
    image: alpine
    volumeMounts:
    - name: config
      mountPath: /conf
    command:
      - sh
      - -c
      - |
        cp /conf/haproxy.cfg /tmp/haproxy.cfg
        while sleep 1; do
          if ! diff /tmp/haproxy.cfg /conf/haproxy.cfg; then 
            echo "Reloading HAProxy."
            killall -HUP haproxy
            cp /conf/haproxy.cfg /tmp/haproxy.cfg
          fi
        done
$ kubectl replace -f k8s/haproxy.yaml --force
pod "haproxy" deleted
pod/haproxy replaced
$ kubectl edit configmap haproxy
$ kubectl get configmap haproxy -o yaml
apiVersion: v1
data:
  haproxy.cfg: |
    global
      daemon

    defaults
      mode tcp
      timeout connect 5s
      timeout client 50s
      timeout server 50s

    listen very-basic-load-balancer
      bind *:80
      server yellow yellow:80
      server green green:80
      server red red:80

    # Note: the services above must exist,
    # otherwise HAproxy won't start.
kind: ConfigMap
metadata:
  creationTimestamp: "2026-01-10T15:20:03Z"
  name: haproxy
  namespace: rainbow
  resourceVersion: "510340"
  uid: 06cea999-252e-4b9f-9d75-19ace2adb1fc
$ kubectl port-forward pod/haproxy 8080:80
Forwarding from 127.0.0.1:8080 -> 80
Forwarding from [::1]:8080 -> 80
Handling connection for 8080
Handling connection for 8080
Handling connection for 8080
Handling connection for 8080
Handling connection for 8080

...
$ curl localhost:8080
🟡This is pod rainbow/yellow-69d94b589c-rbrqk on linux/amd64, serving / for 10.1.1.218:43052.
$ curl localhost:8080
🟢This is pod rainbow/green-556754bb7d-k4zmk on linux/amd64, serving / for 10.1.1.218:48956.
$ curl localhost:8080
🔴This is pod rainbow/red-77f6d65f98-jmmnf on linux/amd64, serving / for 10.1.1.218:42734.
$ curl localhost:8080
🟡This is pod rainbow/yellow-69d94b589c-rbrqk on linux/amd64, serving / for 10.1.1.218:54180.
$ curl localhost:8080
🟢This is pod rainbow/green-556754bb7d-k4zmk on linux/amd64, serving / for 10.1.1.218:54514.
$ curl localhost:8080
🔴This is pod rainbow/red-77f6d65f98-jmmnf on linux/amd64, serving / for 10.1.1.218:55320.
$ curl localhost:8080
🟡This is pod rainbow/yellow-69d94b589c-rbrqk on linux/amd64, serving / for 10.1.1.218:54192.
```

### Should we need to do all these to keep containers running during configuration reload?
- If we need to reload a configuration without stopping running process, without having to stop the world, 
	- if we have some long-running connections running in the container, like web-sockets or video or gaming or VoIP server, we can't just restart the running container midway, we might want to let the existing connections drain, waiting until all the clients have finished their session or whatever
	- then, this method tells it's possible