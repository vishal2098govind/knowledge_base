#k8s #volumes #configuration-management 
- sometime our code needs sensitive information:
	- passwords
	- API tokens
	- TLS keys
	- ...
- It's almost like configuration information, except that it's security sensitive data
-  k8s offers a special resource called `secrets`, which works exactly like `configmaps`
	- this indicates different intents
		- when we see a `configmap` we know that it's just configuration information, and probably not too security sensitive
		- however when we see a `secret` we would not want that to leak or accessible to random folks
	- it can often be useful to have very open permissions on `configmaps` for troubleshooting
	- for secrets, we don't want any open permissions or be readable
- all secrets created will be **base-64 encoded**, and **not encrypted**
- there are multiple types of secrets
	- `docker-registry`
	- `generic`
	- `tls`
```sh
$ kubectl create secret generic codes --from-literal front-door=1234
secret/codes created
$ kubectl get secrets
NAME    TYPE     DATA   AGE
codes   Opaque   1      4s
$ kubectl get secrets
NAME    TYPE     DATA   AGE
codes   Opaque   1      4s
$ kubectl get secrets -o yaml
apiVersion: v1
items:
- apiVersion: v1
  data:
    front-door: MTIzNA==
  kind: Secret
  metadata:
    creationTimestamp: "2026-01-10T20:29:53Z"
    name: codes
    namespace: rainbow
    resourceVersion: "514895"
    uid: 5f32f44b-215a-421b-9cfe-385a6af02424
  type: Opaque
kind: List
metadata:
  resourceVersion: ""
$ base64 -d <<< MTIzNA==
1234%
```