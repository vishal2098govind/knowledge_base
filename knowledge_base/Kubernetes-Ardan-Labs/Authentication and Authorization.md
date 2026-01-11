#k8s #k8s-authentication #rbac

## Definitions
- **Authentication** = verifying the identity of a person
	- on a UNIX system, we can authenticate with login+password, SSH keys, ...
- **Authorization** = listing what they are allowed to do
	- on a UNIX system, this can include file permissions, `sudoer` entries
- Sometimes abbreviated as `authn` and `authz`
- In good modular systems, these things are decoupled
	- e.g. change password or SSH key without having to reset access rights
	- we can imagine like an vertical hour glass, 
		- where on top we have authentication - we have many authentication systems
			- login-password
			- ssh key
			- face-recognition
			- biometrics
		- the authentication system gives us the user id
		- this is decoupled with authorization
			- i.e. if we change the password, it should not change the user id and thus reset permission lists

## Authentication in k8s
- each time we make request to k8s API server, we go through various authenticator and authorizers
	- ~ conveyor belt in a factory where the API requests arrive
	- when we do `kubectl get nodes`, we have bunch of `authn` and `authz` modules that will look into that request
	- **authentication phase** find out the user id - who is making the request
	- after finding the user id, the authentication module slaps a label containing user-id, on the request and the request continues further on the conveyor belt until it arrives at the authorization module/phase
	- in the **authorization phase**, bunch of modules are present that kind of look at the request and check if the request made is allowed by the user-id
### Authentication methods in k8s
### TLS client certificates
```sh
$ kubectl config view
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: DATA+OMITTED
    server: https://127.0.0.1:6443
  name: docker-desktop
contexts:
- context:
    cluster: docker-desktop
    namespace: nginx-helm-demo
    user: docker-desktop
  name: docker-desktop
current-context: docker-desktop
kind: Config
preferences: {}
users:
- name: docker-desktop
  user:
    client-certificate-data: DATA+OMITTED # here we have base-64 encoded literal of the client certificate data
    client-key-data: DATA+OMITTED
```
- #base-64 is an encoding which is often used when we have some binary data and we want to be able to safely copy-paste it in plain text
	- it's called base-64 because it uses 64 unique characters
		- upper-case letters - 26 characters
		- lower-case letters - 26 characters
		- numbers - 10 digits
		- 2 extra characters - depend on implementation of base-64
			- mostly it's going to be + or -
		- it very often ends with two '=' i.e. '\=='
	- it's relatively easy to recognize a base-64 encoded literal by checking if it ends with `==` and has mix of uppercase, lowercase, numbers
- the client certificate data after decoding the base-64 often looks like this
```sh
-----BEGIN CERTIFICATE-----
MIIEowIBAAKCAQEA41rkWIYth3eQGclDjGuIN1TvC3thh+9ywaYR3/FstypWLsAg
YEOzLidqmegG7JH5leE3zS185rEWBcvMA4zV3Gc5NGRHwCakujzIY09lspb9UUtS
3FkHDHm6ahnmBx0rgWv+C0MB1DedrUZN7o7gapyPbQ+KRTILk8h+olF5H+cHO8mo
dZOnkNR+iNtKvaefaGaxzfT3tIUfB9rSd0zscl/pvHL7/VHRDD664L4GtEZ1AYN0
OMMJK5b9Q1RQM9rQdgzD74cJyuopol+XqCL5qNH54GFWhkKB5FHtDiN//Zdam6k6
kQI+AQ9jXdp/0A05akekyWIIr0x9kKlBHeYs7QIDAQABAoIBAQCNiUUEcyue3OkB
lJEPttXLFH3oWPwDoyZ+tYMFNgfnE10pp4PASgUfI3yyqYa9+bY1/0o82ieaef3q
x8+PGjr9BenP2unsDmKYakSZxpDaGDNFS64e7PF9a76IVO0l2pv96xvxKNrXoXPw
rgbQD3RYCnUjgPwkTZ51duiW5j+QgpdVHD1oOzYzfTaQWzBWynB94kRKSkk8Kwkl
ufWB937IubuLtcmn+xIKwo++OS1/gJHjP4cTmzQytsSRWMyMkhtfIyBNCY8I04pv
YQ0viIXLufKXJOuDGdDsl4fwRXxRsK6NFMisEH07Ao95uzRiE7Pvx6pTSG+vquur
I0PtX9JhAoGBAPJGEX7lfiLibGJMYDZ+0QkKzBtIRVpQsfI3J+fVGIbskqIr2L9X
R9MIi8jrOmuxVZp2Ysn7ecRql6pOd5SPwK9cItqiRt47ywu00FGtjBZbSzojkuvM
XjMG2OXZ2GcTm0sFE8yzOLEvoMsNphJSQYfyqR4DMEWxYaebnakRK64lAoGBAPA8
cZzx1UnB/ltWQ9Hu8lsazjoveyYn60LyjM8dB/xu5wKfUh0WpvJQMieO/iuWiHWs
UWX1Kdil8o9C9e7DwO3V7+nXIYlrWjQqdvsH1kj4Ou4jWZW2tvxrdZLhKF0EohMN
pk2kAL8War7c7b4uae+ee239mAUjguTT6lzoMVUpAoGAbvgNF3SCW/Qd9MBK6WsY
z+10I6LssTt45hrBWnzOqS4+060FsE6IBW9Kp5KmxaEKNp+3DBD1azPvmAAs4Y0e
krO++ymNEEmO7SO6r06wdaUHRe/5YavQEcs3GcC4UC442RnQQdhtRSstmRP9VzL5
9Qz+zVJkj+d5dU3f3wqQWUkCgYAh2yLXz20Tm/dQv/OG3nnhRQHTSWW9ltqc8LnP
ZnP9q1N7FyDnYI3ekFbBJHHFis1oaASArfBby+dHMVSfdY729bN97QeY6jwoJePm
tNAKMJF9hhXV944KPMqxDHI3ybNdpArP3lEMnQgmJdSLfNuVQE+KCYVLFBoaTZM1
zLNHIQKBgDtEKWOdazjnQ30DqBFX0lrOYdwWtOjXMGcS+5eCROlngczpJRUfZ37Z
LBNqGnNmTxzY9bTSg4x2+vupV7V5DokTj6hyOu3X3Cm0Ofgz3Ci6b4br8yer59KR
hBJ4Y1aOqQZHUHtOgrTh4mba6Jsr8TWRCO8WmP2Om2DXL1hv1tUf
-----END CERTIFICATE-----%
```
- this is a format named PEM - Privacy-Enhanced-Mail
- this is a format that was designed to make it easy to copy binary stuff in email files in particular
- back in the days when attaching emails to files was risky
- when we needed to send key or certificate or etc, we would use this format which is
```sh
-----BEGIN <> -----
# base-64 encoded data
-----END <> -----
```
- we can have begin-certificate, begin-public-key, begin-private-key, begin-message-signature, begin-whatever etc
- in the middle we have base-64 encoded data
- although we cannot simply decode this base-64 by using `base64 -d <data>`
- to decode this thing is to use a special program like `openssl` that can read PEM certificates

### Bearer tokens
- a secret token in HTTP headers of the request
- tokens conceptually work just like certificate
- conceptually, a certificate is a kind of cryptographic document that says "this is user Vishal, and this is confirmed by this Authority"
- similarly with tokens, we can have a token that says "this is user Vishal, and this is confirmed by this Authority"
- the implementation or protocols or format is different but the concept is exactly the same
### Anonymous requests
- `system:anonymous` is used when we don't pass any authentication tokens or certificates along the request
```sh
$ kubectl get nodes -v6
I0111 23:58:00.533420   53652 loader.go:402] Config loaded from file:  /Users/govind/.kube/config
I0111 23:58:00.535957   53652 envvar.go:172] "Feature gate default state" feature="InformerResourceVersion" enabled=false
I0111 23:58:00.535984   53652 envvar.go:172] "Feature gate default state" feature="WatchListClient" enabled=false
I0111 23:58:00.535992   53652 envvar.go:172] "Feature gate default state" feature="ClientsAllowCBOR" enabled=false
I0111 23:58:00.535998   53652 envvar.go:172] "Feature gate default state" feature="ClientsPreferCBOR" enabled=false
I0111 23:58:00.558763   53652 round_trippers.go:560] GET https://127.0.0.1:6443/api/v1/nodes?limit=500 200 OK in 16 milliseconds
NAME             STATUS   ROLES           AGE   VERSION
docker-desktop   Ready    control-plane   8d    v1.32.2
$ curl https://127.0.0.1:6443/api/v1/nodes
curl: (60) SSL certificate problem: unable to get local issuer certificate
More details here: https://curl.se/docs/sslcerts.html

curl failed to verify the legitimacy of the server and therefore could not
establish a secure connection to it. To learn more about this situation and
how to fix it, please visit the web page mentioned above.
$ curl -k https://127.0.0.1:6443/api/v1/nodes
{
  "kind": "Status",
  "apiVersion": "v1",
  "metadata": {},
  "status": "Failure",
  "message": "nodes is forbidden: User \"system:anonymous\" cannot list resource \"nodes\" in API group \"\" at the cluster scope",
  "reason": "Forbidden",
  "details": {
    "kind": "nodes"
  },
  "code": 403
}%                                                                   $ curl -k https://127.0.0.1:6443/api/v1/nodes -H 'Authorization: Bearer aefjaofe'
{
  "kind": "Status",
  "apiVersion": "v1",
  "metadata": {},
  "status": "Failure",
  "message": "Unauthorized",
  "reason": "Unauthorized",
  "code": 401
}%
```
## JWT tokens
#jwt
```sh
$ kubectl create token default eyJhbGciOiJSUzI1NiIsImtpZCI6Ikx1akJQa2EyVjRsMUlIUmZBMTR3Y0VfSDN5RmNneWxDUWtKU2lFM043aHMifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzY4MTYxMTg2LCJpYXQiOjE3NjgxNTc1ODYsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwianRpIjoiY2MzODdlNjYtMTdjOC00ODIyLTg4YzctZmM0MmIwMzQ4NjY1Iiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJkZWZhdWx0Iiwic2VydmljZWFjY291bnQiOnsibmFtZSI6ImRlZmF1bHQiLCJ1aWQiOiJkY2IwZjBhZi02MDM5LTRmNWEtOTQ2Yy01ZGFjNjYyNTEzZTMifX0sIm5iZiI6MTc2ODE1NzU4Niwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50OmRlZmF1bHQ6ZGVmYXVsdCJ9.s7rT04YgJADVayhaf2yaY0wV-g9M32i8DkLq2AzU5JeV7mjbILJaiv_RXdvmkdKFyk2ujKW-Esou-_SFk3CQKjkneBpWe-dv2lvKnNPG0ah7GhVwz52wVQH5lnw9BTN-1t566uvEU8RCexrI1Lcxqb5J80lue78OUkV3NJFN0L9a6MR0Jl_CasKRC8EeOUrw6-dCyiBDq0A02PGLLfHw01Wk5vBE3aB_DOZBrF7KF5RFkx4wKyevhM_L5UqtlFa1RlEo0WOyxvqCfdTiIPgdCO6ZTL8sl8w8PKfWEkOk_aD4Oe4FpxqyFzfxXXaiixQ2ZtomSTo5W5SZ5zvKN-BkbQ
```
- a JWT is nothing but bunch of base-64 encoded strings joined by dot (.) i.e. header.payload.signature
	- header - JSON
	- payload - JSON
	- signature - not JSON
```sh
$ kubectl create token default | cut -d . -f2  | base64 -d | jq .
{
  "aud": [
    "https://kubernetes.default.svc.cluster.local"
  ],
  "exp": 1768161600,
  "iat": 1768158000,
  "iss": "https://kubernetes.default.svc.cluster.local",
  "jti": "34366ad6-007e-4726-b007-95f0989c1885",
  "kubernetes.io": {
    "namespace": "default",
    "serviceaccount": {
      "name": "default",
      "uid": "dcb0f0af-6039-4f5a-946c-5dac662513e3"
    }
  },
  "nbf": 1768158000,
  "sub": "system:serviceaccount:default:default"
}
```
- JWT token and SSL certificates are conceptually the exact same thing
- both of them say about the user identity and certificate authority (CA)
- these tokens are very often linked to **service accounts**