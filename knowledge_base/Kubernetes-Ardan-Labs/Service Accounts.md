#service-accounts
- A service account is a user that exists in k8s API, for robots
- each namespace has default service account created automatically
```sh
$ kubectl get serviceaccounts
NAME      SECRETS   AGE
default   0         8d
```
- `kube-system` namespace has lot of service accounts corresponding to many controllers on the k8s control plane
- service accounts are little bit like unix users or services
- if we look into them:
```sh
$ cat /etc/passwd
##
# User Database
#
# Note that this file is consulted directly only when the system is running
# in single-user mode.  At other times this information is provided by
# Open Directory.
#
# See the opendirectoryd(8) man page for additional information about
# Open Directory.
##
nobody:*:-2:-2:Unprivileged User:/var/empty:/usr/bin/false
root:*:0:0:System Administrator:/var/root:/bin/sh
daemon:*:1:1:System Services:/var/root:/usr/bin/false
_uucp:*:4:4:Unix to Unix Copy Protocol:/var/spool/uucp:/usr/sbin/uucico
_taskgated:*:13:13:Task Gate Daemon:/var/empty:/usr/bin/false
_networkd:*:24:24:Network Services:/var/networkd:/usr/bin/false
_installassistant:*:25:25:Install Assistant:/var/empty:/usr/bin/false
_lp:*:26:26:Printing Services:/var/spool/cups:/usr/bin/false
_postfix:*:27:27:Postfix Mail Server:/var/spool/postfix:/usr/bin/false
_scsd:*:31:31:Service Configuration Service:/var/empty:/usr/bin/false
_ces:*:32:32:Certificate Enrollment Service:/var/empty:/usr/bin/false
_appstore:*:33:33:Mac App Store Service:/var/db/appstore:/usr/bin/false
_mcxalr:*:54:54:MCX AppLaunch:/var/empty:/usr/bin/false
_appleevents:*:55:55:AppleEvents Daemon:/var/empty:/usr/bin/false
_geod:*:56:56:Geo Services Daemon:/var/db/geod:/usr/bin/false
_devdocs:*:59:59:Developer Documentation:/var/empty:/usr/bin/false
....
```
- in UNIX systems, we have lot of users for different services that could run on the system
	- this is a common practice on UNIX system, when we want to run a service, like `nginx` or `apache` or `mysql` etc, we don't want to run them as `root` user
		- because running as root user could be dangerous 
		- if there is any security problem in the service, and if somebody can exploit that security problem, that vulnerability, 
		- and if that service runs as root user, then the attacker get's root access on the machine
		- so to contain things a little bit, we run services with non-privileged users, so for instance, `apache` web server could run with user `apache` or `http` or `wwwdata` or whatever
		- so if somebody manages to exploit the vulnerability, it would just be in the web server, instead of being root, 
- similar idea with `serviceaccounts` in k8s
- when we want to run something on k8s, like an ingress controller, or some application that needs to communicate with k8s API, we create a service account for that application, and we setup the application to use that service account, giving fine-grained narrow permissions to the service account, instead of giving them cluster admin permissions

#### Creating and using service accounts
- creating a application called autoscaler, that will have permission only to scale up and scale down a specific deployment (worker), and no other permission
```sh
$ kubectl create namespace securitydemo
namespace/securitydemo created
$ kubectl create serviceaccount autoscaler
serviceaccount/autoscaler created
$ kubectl apply -f k8s/dockercoins.yaml
deployment.apps/hasher created
service/hasher created
deployment.apps/redis created
service/redis created
deployment.apps/rng created
service/rng created
deployment.apps/webui created
deployment.apps/worker created
service/webui created
```
- to run something with that `serviceaccount`, we have two options
	- specify in yaml using `pod.spec.serviceAccountName` in case of pod
	- specify in `kubectl` by passing `--overrides='{"spec": {"serviceAccountName": "autoscaler"}}'` - exactly like `kubectl patch`
```sh
$ kubectl run autoscaler-pod -it --rm --image alpine --overrides='{"spec": {"serviceAccountName": "autoscaler"}}' -- sh
If you dont see a command prompt, try pressing enter.
/ #

...
$ kubectl describe pod autoscaler-pod
Name:             autoscaler-pod
Namespace:        securitydemo
Priority:         0
Service Account:  autoscaler
Node:             docker-desktop/192.168.65.3
Start Time:       Mon, 12 Jan 2026 01:29:21 +0530
Labels:           run=autoscaler-pod
Annotations:      <none>
Status:           Running
IP:               10.1.2.8
IPs:
  IP:  10.1.2.8
Containers:
  autoscaler-pod:
    Container ID:  docker://8a8d8f691b7c6138827c68d7f26b90bebf0e52a9e6ec1da988c6ea9e8b32c4eb
    Image:         alpine
    Image ID:      docker-pullable://alpine@sha256:865b95f46d98cf867a156fe4a135ad3fe50d2056aa3f25ed31662dff6da4eb62
    Port:          <none>
    Host Port:     <none>
    Args:
      sh
    State:          Running
      Started:      Mon, 12 Jan 2026 01:29:25 +0530
    Ready:          True
    Restart Count:  0
    Environment:    <none>
    Mounts:
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-r5mj2 (ro)
Conditions:
  Type                        Status
  PodReadyToStartContainers   True
  Initialized                 True
  Ready                       True
  ContainersReady             True
  PodScheduled                True
Volumes:
  kube-api-access-r5mj2:
    Type:                    Projected (a volume that contains injected data from multiple sources)
    TokenExpirationSeconds:  3607
    ConfigMapName:           kube-root-ca.crt
    ConfigMapOptional:       <nil>
    DownwardAPI:             true
QoS Class:                   BestEffort
Node-Selectors:              <none>
Tolerations:                 node.kubernetes.io/not-ready:NoExecute op=Exists for 300s
                             node.kubernetes.io/unreachable:NoExecute op=Exists for 300s
Events:
  Type    Reason     Age   From               Message
  ----    ------     ----  ----               -------
  Normal  Scheduled  25s   default-scheduler  Successfully assigned securitydemo/autoscaler-pod to docker-desktop
  Normal  Pulling    23s   kubelet            Pulling image "alpine"
  Normal  Pulled     21s   kubelet            Successfully pulled image "alpine" in 2.682s (2.682s including waiting). Image size: 3870955 bytes.
  Normal  Created    20s   kubelet            Created container: autoscaler-pod
  Normal  Started    20s   kubelet            Started container autoscaler-pod
```
- now, for the pod to be able to talk to k8s API
```sh
$ kubectl run autoscaler-pod -it --rm --image alpine --overrides='{"spec": {"serviceAccountName": "autoscaler"}}' -- sh
If you dont see a command prompt, try pressing enter.
/ # kubectl
sh: kubectl: not found
```
- since we do-not have `kubectl` in alpine, we can create a custom image on the fly using `nixery.dev` with the required programs/packages
```sh
$ docker run -it 'nixery.dev/shell/kubectl/curl/jq:latest'
$ docker images 'nixery.dev/shell/kubectl/curl/jq:latest'
REPOSITORY                         TAG       IMAGE ID       CREATED        SIZE
nixery.dev/shell/kubectl/curl/jq   latest    1adce3207b8b   56 years ago   306MB
```
- now, we can create a pod using that image
```sh
$ kubectl run autoscaler-pod -it --rm --image nixery.dev/shell/kubectl/curl/jq --overrides='{"spec": {"serviceAccountName": "autoscaler"}}' -- bash
If you don't see a command prompt, try pressing enter.'
bash-5.3# kubectl get pods
Error from server (Forbidden): pods is forbidden: User "system:serviceaccount:securitydemo:autoscaler" cannot list resource "pods" in API group "" in the namespace "securitydemo"
```
- where does those details in the account or right identity about the service account, come from - no magic
- whenever we run a pod in k8s, we have a file in `/var/run/secrets/kubernetes.io/serviceaccount/token`
```sh
$ kubectl run autoscaler-pod -it --rm --image nixery.dev/shell/kubectl/curl/jq --overrides='{"spec": {"serviceAccountName": "autoscaler"}}' -- bash
If you don't see a command prompt, try pressing enter.'
bash-5.3# cat /var/run/secrets/kubernetes.io/serviceaccount/
..2026_01_11_20_12_51.3365482807/ ca.crt                            token
..data/                           namespace
bash-5.3# cat /var/run/secrets/kubernetes.io/serviceaccount/token
eyJhbGciOiJSUzI1NiIsImtpZCI6Ikx1akJQa2EyVjRsMUlIUmZBMTR3Y0VfSDN5RmNneWxDUWtKU2lFM043aHMifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzk5Njk4MzcxLCJpYXQiOjE3NjgxNjIzNzEsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwianRpIjoiZTMzN2I2NDAtMmMzMy00MTE3LTkyMTgtYTQ1YjZiODU0ZjlhIiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJzZWN1cml0eWRlbW8iLCJub2RlIjp7Im5hbWUiOiJkb2NrZXItZGVza3RvcCIsInVpZCI6IjYwYmFlNWU0LWY1NTctNGMzZC1hNmEyLWE1MjQ3ZTU5ODYxOCJ9LCJwb2QiOnsibmFtZSI6ImF1dG9zY2FsZXItcG9kIiwidWlkIjoiZmFlMDBmZjYtMGJmNC00ZTM0LTlhYjEtYzRjN2ZhMTU1NDhkIn0sInNlcnZpY2VhY2NvdW50Ijp7Im5hbWUiOiJhdXRvc2NhbGVyIiwidWlkIjoiYjgwZTNiOTgtODI0YS00YTFjLTlkOWYtODc1ODU2MzNiNjRiIn0sIndhcm5hZnRlciI6MTc2ODE2NTk3OH0sIm5iZiI6MTc2ODE2MjM3MSwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50OnNlY3VyaXR5ZGVtbzphdXRvc2NhbGVyIn0.iybGgs3nzdWF5qefrM8uxgCWiq0XfIrsxVKneIgWkxtxCS2JEUYZctGo7Q3HGPqmfSOuS6SGkjEBuWmcy5Oc-NmlnDEr38KafHXLbQjdD0zeiWkzFoYE9AToGXVyVdK-8ij_V2fbHqI6H7nNYLyUxOZVLZKXnKAirgz_-2jlsVu2cpjZaJL_Euq3bUF6lIaLRh8mk0Xut3GUbzPGQtXPvMzq57eUJ62NGD_6w4Hw61OSVlEZJ1-GSv0kfxff9Zkv5PHZJOCMEhSxDEnksvebXXHWMQTN9Fo49pIUTYEQ0EAo7uD40YE3I-ZjpuislE2BxTQt-BkNW7Rz-ndgK_m4xg
```
- this is a serviceAccount token
- this is the account k8s uses automatically to communicate with k8s API
```sh
TOKEN=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)
bash-5.3# TOKEN
bash: TOKEN: command not found
bash-5.3# echo $TOKEN
eyJhbGciOiJSUzI1NiIsImtpZCI6Ikx1akJQa2EyVjRsMUlIUmZBMTR3Y0VfSDN5RmNneWxDUWtKU2lFM043aHMifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzk5Njk4MzcxLCJpYXQiOjE3NjgxNjIzNzEsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwianRpIjoiZTMzN2I2NDAtMmMzMy00MTE3LTkyMTgtYTQ1YjZiODU0ZjlhIiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJzZWN1cml0eWRlbW8iLCJub2RlIjp7Im5hbWUiOiJkb2NrZXItZGVza3RvcCIsInVpZCI6IjYwYmFlNWU0LWY1NTctNGMzZC1hNmEyLWE1MjQ3ZTU5ODYxOCJ9LCJwb2QiOnsibmFtZSI6ImF1dG9zY2FsZXItcG9kIiwidWlkIjoiZmFlMDBmZjYtMGJmNC00ZTM0LTlhYjEtYzRjN2ZhMTU1NDhkIn0sInNlcnZpY2VhY2NvdW50Ijp7Im5hbWUiOiJhdXRvc2NhbGVyIiwidWlkIjoiYjgwZTNiOTgtODI0YS00YTFjLTlkOWYtODc1ODU2MzNiNjRiIn0sIndhcm5hZnRlciI6MTc2ODE2NTk3OH0sIm5iZiI6MTc2ODE2MjM3MSwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50OnNlY3VyaXR5ZGVtbzphdXRvc2NhbGVyIn0.iybGgs3nzdWF5qefrM8uxgCWiq0XfIrsxVKneIgWkxtxCS2JEUYZctGo7Q3HGPqmfSOuS6SGkjEBuWmcy5Oc-NmlnDEr38KafHXLbQjdD0zeiWkzFoYE9AToGXVyVdK-8ij_V2fbHqI6H7nNYLyUxOZVLZKXnKAirgz_-2jlsVu2cpjZaJL_Euq3bUF6lIaLRh8mk0Xut3GUbzPGQtXPvMzq57eUJ62NGD_6w4Hw61OSVlEZJ1-GSv0kfxff9Zkv5PHZJOCMEhSxDEnksvebXXHWMQTN9Fo49pIUTYEQ0EAo7uD40YE3I-ZjpuislE2BxTQt-BkNW7Rz-ndgK_m4xg
bash-5.3# kubectl get pods -v6
I0111 20:21:18.229239      61 merged_client_builder.go:163] Using in-cluster namespace
I0111 20:21:18.229674      61 merged_client_builder.go:121] Using in-cluster configuration
I0111 20:21:18.230207      61 envvar.go:172] "Feature gate default state" feature="ClientsAllowCBOR" enabled=false
I0111 20:21:18.230290      61 envvar.go:172] "Feature gate default state" feature="ClientsPreferCBOR" enabled=false
I0111 20:21:18.230301      61 envvar.go:172] "Feature gate default state" feature="InOrderInformers" enabled=true
I0111 20:21:18.230308      61 envvar.go:172] "Feature gate default state" feature="InformerResourceVersion" enabled=false
I0111 20:21:18.230315      61 envvar.go:172] "Feature gate default state" feature="WatchListClient" enabled=false
I0111 20:21:18.249295      61 merged_client_builder.go:121] Using in-cluster configuration
I0111 20:21:18.261189      61 round_trippers.go:632] "Response" verb="GET" url="https://10.96.0.1:443/api/v1/namespaces/securitydemo/pods?limit=500" status="403 Forbidden" milliseconds=11
I0111 20:21:18.261837      61 helpers.go:246] server response object: %s[{
  "kind": "Status",
  "apiVersion": "v1",
  "metadata": {},
  "status": "Failure",
  "message": "pods is forbidden: User \"system:serviceaccount:securitydemo:autoscaler\" cannot list resource \"pods\" in API group \"\" in the namespace \"securitydemo\"",
  "reason": "Forbidden",
  "details": {
    "kind": "pods"
  },
  "code": 403
}]
Error from server (Forbidden): pods is forbidden: User "system:serviceaccount:securitydemo:autoscaler" cannot list resource "pods" in API group "" in the namespace "securitydemo"
bash-5.3# curl https://10.96.0.1:443/api/v1/namespaces/securitydemo/pods -H 'Authorization: Bearer $TOKEN'
curl: (60) SSL certificate OpenSSL verify result: unable to get local issuer certificate (20)
More details here: https://curl.se/docs/sslcerts.html

curl failed to verify the legitimacy of the server and therefore could not
establish a secure connection to it. To learn more about this situation and
how to fix it, please visit the webpage mentioned above.
bash-5.3#
bash-5.3# curl https://10.96.0.1:443/api/v1/namespaces/securitydemo/pods -H 'Authorization: Bearer $TOKEN'
curl: (60) SSL certificate OpenSSL verify result: unable to get local issuer certificate (20)
More details here: https://curl.se/docs/sslcerts.html

curl failed to verify the legitimacy of the server and therefore could not
establish a secure connection to it. To learn more about this situation and
how to fix it, please visit the webpage mentioned above.
bash-5.3# curl -k https://10.96.0.1:443/api/v1/namespaces/securitydemo/pods -H 'Authorization: Bearer $TOKEN'
{
  "kind": "Status",
  "apiVersion": "v1",
  "metadata": {},
  "status": "Failure",
  "message": "Unauthorized",
  "reason": "Unauthorized",
  "code": 401
bash-5.3# echo "curl -k https://10.96.0.1:443/api/v1/namespaces/securitydemo/pods -H 'Authorization: Bearer $TOKEN'"
curl -k https://10.96.0.1:443/api/v1/namespaces/securitydemo/pods -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6Ikx1akJQa2EyVjRsMUlIUmZBMTR3Y0VfSDN5RmNneWxDUWtKU2lFM043aHMifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzk5Njk4MzcxLCJpYXQiOjE3NjgxNjIzNzEsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwianRpIjoiZTMzN2I2NDAtMmMzMy00MTE3LTkyMTgtYTQ1YjZiODU0ZjlhIiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJzZWN1cml0eWRlbW8iLCJub2RlIjp7Im5hbWUiOiJkb2NrZXItZGVza3RvcCIsInVpZCI6IjYwYmFlNWU0LWY1NTctNGMzZC1hNmEyLWE1MjQ3ZTU5ODYxOCJ9LCJwb2QiOnsibmFtZSI6ImF1dG9zY2FsZXItcG9kIiwidWlkIjoiZmFlMDBmZjYtMGJmNC00ZTM0LTlhYjEtYzRjN2ZhMTU1NDhkIn0sInNlcnZpY2VhY2NvdW50Ijp7Im5hbWUiOiJhdXRvc2NhbGVyIiwidWlkIjoiYjgwZTNiOTgtODI0YS00YTFjLTlkOWYtODc1ODU2MzNiNjRiIn0sIndhcm5hZnRlciI6MTc2ODE2NTk3OH0sIm5iZiI6MTc2ODE2MjM3MSwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50OnNlY3VyaXR5ZGVtbzphdXRvc2NhbGVyIn0.iybGgs3nzdWF5qefrM8uxgCWiq0XfIrsxVKneIgWkxtxCS2JEUYZctGo7Q3HGPqmfSOuS6SGkjEBuWmcy5Oc-NmlnDEr38KafHXLbQjdD0zeiWkzFoYE9AToGXVyVdK-8ij_V2fbHqI6H7nNYLyUxOZVLZKXnKAirgz_-2jlsVu2cpjZaJL_Euq3bUF6lIaLRh8mk0Xut3GUbzPGQtXPvMzq57eUJ62NGD_6w4Hw61OSVlEZJ1-GSv0kfxff9Zkv5PHZJOCMEhSxDEnksvebXXHWMQTN9Fo49pIUTYEQ0EAo7uD40YE3I-ZjpuislE2BxTQt-BkNW7Rz-ndgK_m4xg'
bash-5.3# curl -k https://10.96.0.1:443/api/v1/namespaces/securitydemo/pods -H 'Authorization: Bearer $TOKEN'
{
  "kind": "Status",
  "apiVersion": "v1",
  "metadata": {},
  "status": "Failure",
  "message": "Unauthorized",
  "reason": "Unauthorized",
  "code": 401
}
bash-5.3# curl -k https://10.96.0.1:443/api/v1/namespaces/securitydemo/pods
{
  "kind": "Status",
  "apiVersion": "v1",
  "metadata": {},
  "status": "Failure",
  "message": "pods is forbidden: User \"system:anonymous\" cannot list resource \"pods\" in API group \"\" in the namespace \"securitydemo\"",
  "reason": "Forbidden",
  "details": {
    "kind": "pods"
  },
  "code": 403
}
```
- to curl without using `-k` (i.e. insecure), we can pass `cacert` option to the curl with the certificate of the API server
- every pod has the certificate of the API Server in `/var/run/secrets/kubernetes.io/serviceaccount/ca.crt`
```sh
$ cat /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
-----BEGIN CERTIFICATE-----
MIIDBTCCAe2gAwIBAgIIRUtSRIEKVDowDQYJKoZIhvcNAQELBQAwFTETMBEGA1UE
AxMKa3ViZXJuZXRlczAeFw0yNjAxMDMxMzUzMzZaFw0zNjAxMDExMzUzMzZaMBUx
EzARBgNVBAMTCmt1YmVybmV0ZXMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK
AoIBAQDc7nyO0Ltfu6yJn+qGxB7BECa63hx7NeS7XlUpeGDcuuqMkC62ltKE24QU
K4evmuqYYMcnbon33lnus/9aOEzEUlfpzORlELVz8V5829HFLmh36Q9YOg7ToBaI
AZ/f7wYRTpuK4OzERcdyNh5+9/lZiLooaxkzdIbrymcrrtR2Yto8OcoD2aXS4BIl
vpptOJ78RmIdCbKiECjQOKq9yY1nVzH66lQveQQNAo2e8+qqGD4w8KT3gQFQrDom
ZgyPpyazCqmqYpZ+Z12333Uxy/fwNN6bN40+5WaFGoItC4MR0qw7Lf7t93hj/X0i
gfdIFNOt++zrx4MxDiN+9qbB4RSnAgMBAAGjWTBXMA4GA1UdDwEB/wQEAwICpDAP
BgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSUzlvw9j3ELn+Nb+fKPp5lTGI/9TAV
BgNVHREEDjAMggprdWJlcm5ldGVzMA0GCSqGSIb3DQEBCwUAA4IBAQBR+dNi/Nvr
s+U+AoeJfVTY3SBEkoogtzTZr+dc72A3hynwBM+VLrNarLSfC2CTqSlfsF4QYKmN
7w4b9OSF0d9DwPcuX+9NfskSm5jAB8QWHgn3/3gsTJQsrd/s7Jw0i/q4TMwcmkQj
ikYvRAo6jM+LWzY+2CBg3Ym4YBZ6zRNVjS7zfh0Mn04ctCTdysSH7l4wgzNA99yl
xuXS1JBDvj+VhEDTl4cYCB/oOEiazCz8akpS4by/IhwQy35XW9ohBH9n4EUx8UEf
35RH8meqi6ryvHnmRDPzzMLrdapLiOsBu4WIcu4Hk7GxPjPM5XNPVTvACY4koJDb
hpAvemnHMaLC
-----END CERTIFICATE-----
```
- thus, now without having to do `-k` with curl, we can do
```sh
$ curl https://$KUBERNETES_SERVICE_HOST:$KUBERNETES_SERVICE_PORT -H "Authorization: Bearer $TOKEN" --cacert /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
{
  "kind": "Status",
  "apiVersion": "v1",
  "metadata": {},
  "status": "Failure",
  "message": "forbidden: User \"system:serviceaccount:securitydemo:autoscaler\" cannot get path \"/\"",
  "reason": "Forbidden",
  "details": {},
  "code": 403
}
```

## RBAC
#rbac
- RBAC - Role based access control
```sh
$ kubectl get pods
Error from server (Forbidden): pods is forbidden: User "system:serviceaccount:securitydemo:autoscaler" cannot list resource "pods" in API group "" in the namespace "securitydemo"
```
- this is where, the RBAC will kick in
- there's no such thing like `kubeclt grant list pods to user system:serviceaccount:securitydemo:autoscaler`
- at first we need to create a `role`, and then we need to associate or bind that role with the user
- we don't directly give permissions to users, we put the permissions in roles, and then we associate the role with users
```sh
# 1. create role
$ kubectl create role list-pods-role --verb=list --resource=pods
role.rbac.authorization.k8s.io/list-pods-role created
$ kubectl get roles
NAME             CREATED AT
list-pods-role   2026-01-11T20:58:08Z
```
- associate role with the user or `serviceaccount` by creating a `rolebinding`
```sh
$ kubectl create rolebinding -h
...
Usage:
  kubectl create rolebinding NAME --clusterrole=NAME|--role=NAME
[--user=username] [--group=groupname]
[--serviceaccount=namespace:serviceaccountname] [--dry-run=server|client|none]
[options]
...
$ kubectl create rolebinding autoscaler-can-list-pods --role=list-pods-role --serviceaccount=securitydemo:autoscaler
rolebinding.rbac.authorization.k8s.io/autoscaler-can-list-pods created
$ kubectl get rolebinding
NAME                       ROLE                  AGE
autoscaler-can-list-pods   Role/list-pods-role   6m23s
```
- now, inside the pod, we can do
```sh
bash-5.3# kubectl get pods
NAME                      READY   STATUS    RESTARTS   AGE
autoscaler-pod            1/1     Running   0          53m
hasher-99bbd4bb-vfdlp     1/1     Running   0          76m
redis-7b47f84cc4-9b7s5    1/1     Running   0          76m
rng-65d885d498-f667w      1/1     Running   0          76m
webui-74bb6bbc59-478mc    1/1     Running   0          76m
worker-5c6f84b477-dnncf   1/1     Running   0          76m
worker-5c6f84b477-fxwtg   1/1     Running   0          76m
worker-5c6f84b477-jmxv2   1/1     Running   0          76m
worker-5c6f84b477-n4wmd   1/1     Running   0          76m
worker-5c6f84b477-nvddl   1/1     Running   0          76m
worker-5c6f84b477-qkrj9   1/1     Running   0          76m
worker-5c6f84b477-qr45k   1/1     Running   0          76m
worker-5c6f84b477-wv972   1/1     Running   0          76m
worker-5c6f84b477-xhdn6   1/1     Running   0          76m
worker-5c6f84b477-zkgnk   1/1     Running   0          76m
```
- now, for the pod to be able to scale up or down worker deployment
```sh
bash-5.3# kubectl scale deployment worker --replicas=2
error: no objects passed to scale deployments.apps "worker" is forbidden: User "system:serviceaccount:securitydemo:autoscaler" cannot get resource "deployments" in API group "apps" in the namespace "securitydemo"
```
- this asks for permissions to `get resource "deployments"`
```sh
# pod
bash-5.3# kubectl scale deployment worker --replicas=2
error: no objects passed to scale deployments.apps "worker" is forbidden: User "system:serviceaccount:securitydemo:autoscaler" cannot get resource "deployments" in API group "apps" in the namespace "securitydemo"
-----
# host machine
$ kubectl create role get-deployments-role --verb get --resource deployments

$ kubectl create rolebinding autoscaler-can-get-deployments --role get-deployments-role --serviceaccount securitydemo:autoscaler
----
# pod
$ kubectl scale deployment worker --replicas=2
Error from server (Forbidden): deployments.apps "worker" is forbidden: User "system:serviceaccount:securitydemo:autoscaler" cannot patch resource "deployments/scale" in API group "apps" in the namespace "securitydemo"
```
- this little detail will let us have permissions to scale up and scale down, without having permissions on everything else
- the fact that we can give permission for the scale sub-resources, means that we can give access to that sub-resource and give access specifically to that scale, without giving access to, for instance, changing the image or metadata or deleting deployment etc etc.
- when we talk about k8s objects, sometimes we talk about
	- objects
	- or resources
	- or manifests
- and very often, we use all these in a kind of inter changeable way, when we say
	- "we need to create an object"
	- or "we need to create a resource"
	- or "we need to write a manifest"
- what's exactly the difference between these terms
	- object is the thing that exists in k8s API
		- typically it corresponds to some entry in `etcd`
	- the manifest is something that is present in the entry of `etcd`
		- it's YAML or JSON that is given to API Server to describe an object
		- internally, inside `etcd` we have `protobufs` for these objects
			- we might have more or less fields than what we have in manifest
	- resources are a way to access the objects through the API
		- sometimes we've multiple resources for the same object because we've different API endpoints
		- ~ social media app - users, friends, followers
			- /users endpoint, /followers endpoint, /friends endpoint
		- similarly multiple resources and api routes giving access to same object
	- sub-resources
		- special resource that gives access to reduced object
		- we have multiple different sub-resources depending on the type of object
```sh
----
# host machine
$ kubectl create role patch-deployments-scale-role --verb patch --resource deployments/scale

$ kubectl create rolebinding autoscaler-can-patch-scale --role patch-deployments-scale-role --serviceaccount securitydemo:autoscaler
----
bash-5.3# kubectl scale deployment worker --replicas=4
deployment.apps/worker scaled
---
$  kubectl get all
NAME                          READY   STATUS              RESTARTS   AGE
pod/autoscaler-pod            1/1     Running             0          70m
pod/hasher-99bbd4bb-vfdlp     1/1     Running             0          93m
pod/redis-7b47f84cc4-9b7s5    1/1     Running             0          93m
pod/rng-65d885d498-f667w      1/1     Running             0          93m
pod/webui-74bb6bbc59-478mc    1/1     Running             0          93m
pod/worker-5c6f84b477-mmrdn   0/1     ContainerCreating   0          2s
pod/worker-5c6f84b477-pgbfq   0/1     ContainerCreating   0          2s
pod/worker-5c6f84b477-qkrj9   1/1     Running             0          93m
pod/worker-5c6f84b477-qr45k   1/1     Running             0          93m

NAME             TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
service/hasher   ClusterIP   10.98.77.203     <none>        80/TCP         93m
service/redis    ClusterIP   10.111.75.228    <none>        6379/TCP       93m
service/rng      ClusterIP   10.100.121.121   <none>        80/TCP         93m
service/webui    NodePort    10.98.235.93     <none>        80:30227/TCP   93m

NAME                     READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/hasher   1/1     1            1           93m
deployment.apps/redis    1/1     1            1           93m
deployment.apps/rng      1/1     1            1           93m
deployment.apps/webui    1/1     1            1           93m
deployment.apps/worker   2/4     4            2           93m

NAME                                DESIRED   CURRENT   READY   AGE
replicaset.apps/hasher-99bbd4bb     1         1         1       93m
replicaset.apps/redis-7b47f84cc4    1         1         1       93m
replicaset.apps/rng-65d885d498      1         1         1       93m
replicaset.apps/webui-74bb6bbc59    1         1         1       93m
replicaset.apps/worker-5c6f84b477   4         4         2       93m
```
- creating `role` and `rolebindings` using YAML
```sh
$ kubectl create role get-deployment --verb get --resource deployment -o json --dry-run=client
{
    "kind": "Role",
    "apiVersion": "rbac.authorization.k8s.io/v1",
    "metadata": {
        "name": "get-deployment",
        "creationTimestamp": null
    },
    "rules": [
        {
            "verbs": [
                "get"
            ],
            "apiGroups": [
                "apps"
            ],
            "resources": [
                "deployments"
            ]
        }
    ]
}
```
- for much more fine-grained access control, we can mention `--resourceName` so that the permissions are only for the mentioned resource names of the resource types mentioned in `resources`
```sh
$ kubectl create role autoscaler-role --verb get --resource deployments --resource-name worker --dry-run=client -o yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  creationTimestamp: null
  name: autoscaler-role
rules:
- apiGroups:
  - apps
  resourceNames:
  - worker
  resources:
  - deployments
  verbs:
  - get
```
- can extend it to allow scale up and down
```sh
$ cat autoscaler-role.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  creationTimestamp: null
  name: autoscaler-role
rules:
- apiGroups:
  - apps
  resourceNames:
  - worker
  resources:
  - deployments
  verbs:
  - get
- apiGroups:
  - apps
  resourceNames:
  - worker
  resources:
  - deployments/scale
  verbs:
  - patch
```
- add `rolebinding`
```sh
$ kubectl create rolebinding autoscaler-can --role autoscaler-role --serviceaccount securitydemo:autoscaler --dry-run=client -o yaml
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  creationTimestamp: null
  name: autoscaler-can
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: autoscaler-role
subjects:
- kind: ServiceAccount
  name: autoscaler
  namespace: securitydemo
$ cat autoscaler-role.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: autoscaler-role
rules:
- apiGroups:
  - apps
  resourceNames:
  - worker
  resources:
  - deployments
  verbs:
  - get
- apiGroups:
  - apps
  resourceNames:
  - worker
  resources:
  - deployments/scale
  verbs:
  - patch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: autoscaler-can
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: autoscaler-role
subjects:
- kind: ServiceAccount
  name: autoscaler
  namespace: securitydemo
---
$ kubectl apply -f autoscaler-role.yaml
role.rbac.authorization.k8s.io/autoscaler-role created
rolebinding.rbac.authorization.k8s.io/autoscaler-can created
```
- now, if we try to scale in the pod
```sh
bash-5.3# kubectl scale deployment worker --replicas=4
deployment.apps/worker scaled
bash-5.3# kubectl scale deployment webui --replicas=3
error: no objects passed to scale deployments.apps "webui" is forbidden: User "system:serviceaccount:securitydemo:autoscaler" cannot get resource "deployments" in API group "apps" in the namespace "securitydemo"
```
