#k8s 

## `kubectl explain API`
- If u want to quickly get documentation of how to specify a specific field in yaml, just quickly go to terminal and leverage `kubectl explain` API
```sh
$ kubectl explain pod
$ kubectl explain pod.spec
$ kubectl explain pod.spec.containers
$ kubectl explain pod.spec.containers.env
$ kubectl explain pod.spec.containers.env.valueFrom
$ kubectl explain pod.spec.containers.env.valueFrom.resourceFieldRef
```
- now u can easily understand how mention resourceFieldRef
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: sales

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sales
  namespace: sales
spec:
  selector:
    matchLabels:
      app: sales
  template:
    metadata:
      labels:
        app: sales
    spec:
      containers:
        - name: sales
          image: vishalgovind/sales
          imagePullPolicy: Never
          resources:
            requests:
              cpu: "250m"
              memory: "36Mi"
            limits:
              cpu: "250m"
              memory: "36Mi"
          env:
            - name: GOMAXPROCS
              valueFrom:
                resourceFieldRef:
                  resource: limits.cpu

```