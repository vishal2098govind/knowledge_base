#jwt #cryptography #asymmetric-encryption #hashing

Three parts of JWT separated by period "." are:
{base-64 url encoded header}.{base-64 url encoded payload}.{base-64 url encoded signature}
- base-64-url-encoded header
- base-64-url-encoded payload
- base-64 encoded signature

A sample JWT token:
```
aonfaeoeaobraegafojnaeojjff.oaenfoejajfeaojbaeoggaegae.ajegoaegaebaejjbaejgjbaejjg
```

decoded header:
```
{
  "typ": "JWT",
  "alg": "RS256",
  "kid": "PcX98GX420T1X6sBDkzhQmqgwMU"
}
```
decoded payload:
```
{
  "aud": "3dc3f710-d4ff-474c-8cb9-c9937cbf80c9",
  "iss": "https://login.microsoftonline.com/0242cf3f-b856-4261-9ea6-484/v2.0",
  "iat": 1770185836,
  "nbf": 1770185836,
  "exp": 1770236536,
  "aio": "AbQAS/8bAAAAdzsD6IBM/vkYT5roOmV5uNEFUNEA1s3zELv85+73Nru/4Impar6oJ6e1NCPxkW4bfeimIB7c6E/DDgSxrZF3viuV22Te3pmUvrqGrGPdX5n4FFPYnWRerFgklOIirkd5XRiLXtQ4LgXNJbl0hrqu9+SHXzmTXpAorSPy82VbzhA62eC3GhDsH/+EgvGzjGArbpA05V1TxWV4CyX+HqVsiYcfAUV0X6kff2mFUu3sZe0=",
  "email": "vgovind@in.h.com",
  "groups": [
    "789b6fe7-544d-49a8-ac28-165502393f8e"
  ],
  "name": "Vishal Govind",
  "oid": "b22d4542-e827-4991-9849-46c18938aee7",
  "preferred_username": "vgovind@in.h.com",
  "rh": "1.ATcAP89CAla4YUKeppHsCwLUhBD3wz3_1ExHjLnJk3y_gMnQAbY3AA.",
  "roles": [
    "allmysons"
  ],
  "sid": "001f28aa-6de1-89d8-4157-cc7b47a76ccd",
  "sub": "JMNLTeruBYRos2XD3-YUd1nOIzrh4mOX31pgpwCb4sU",
  "tid": "0242cf3f-b856-4261-9ea6-91ec0b02d484",
  "uti": "rlQSYsp1O06MExE1zbEZAA",
  "ver": "2.0"
}
```
encoded signature:
```
K_B5PurlDKt6-12r-F0LVXrD5Y0cH_nPIrUf2P_AOXyAp_LCd9gsbJQDMyUK-oOWDAptVtTTDkUVgVu7VIgT6mRTXQRVi2S2ycKw6pcmBXg3qxMpB_qTYk4CO84c62k2Yo5C997hQtRAGaRYnMaFSIdxWfIMoIV01tQ4nTH_ZIS13wZZ0BO3sWYEjkn2HcsFkvZ9b5t_rjHiLqTxdNseCkdxrktgS9caEf6HMsPGz2roX2mHTJR8Kracb6UcUdV0ZVGkx73Y1VWHgJqKYYhGaaWwZ4dkvQeX0OYuzCLZu59CX2UtM7h0w3budL1Jh7Wf7Ecf9IES2n1JdR0j-Y4u4g
```
verified public key used:
```
{
  "e": "AQAB",
  "kty": "RSA",
  "n": "hTfqnXJAT9zSM4Xm5QHD3kIW1tXHj5-1xj35D3ht0s9WogK2nzcwlSDWivqnt1aED3njfI9V6V8rOI5Sde8KB87QfLH6SmRtG4dLixDfFoulHPamQ9lEp5i8xVZeypRM03O30AGCfx0G83JiiHVAVLqpjc8Ervl6yMNHyg8hyBkbV_IyeZVVJW2nu0ljkQwyTLr1GH1h2I7a_ioK6WlW7lNK8OYLSEQ_B9ecaKf8dPDy1_Zt5f7I6RdJPRicSgsBUeX34CuJJuBiOB0k4TihJVhc43YexUDo-Sd_e2P3BgpdK3I0ksX5c58yO2z0OvpFHYSg2CYOdzEOj_mLKbOAzQ"
}
```

The authorization server first creates the header and payload JSON objects

Base-64 URL encoding of these json objects
{base-64-url-encoded-header}.{base-64-url-encoded-payload}
- base-64 encoding is done to make sure any binary data is converted to text
- url encoding is done to make sure the tokens can be sent as part of a URL

Next, from the `{base-64 url encoded header}.{base-64 url encoded payload}`, the signature is created by signing with private key of the authorization server

The header json object mentions that it is using "RS256" algorithm which is short-hand for RSA Async-Key-Encryption Algorithm and SHA-256 Hashing Algorithm

**Why are two algorithms used?**
- first a hash is created the data using a hashing algorithm, here SHA-256
- the hashed output is signed using the RSA private key of the authorization server
- the string obtained by encoded header and encoded payload together separated by a dot, is the string that is signed
- the signature is then base-64 url encoded and then added as the third component of the JWT
- any entity (here authorization server) which creates a signature will use it's own private key to sign it.
- the receiver on the other hand will verify the signature using the public key of the entity (here authorization server)

**If the algorithm specified was RS512 instead of RS256**
- for that, the hashing algorithm used would be SHA-512

**What are the signing algorithms that could be used by an authorization server**
- we can find out from the `openid-configuration` of that authorization server from the `https://{authorization-server-host-path}/.well-known/openid-configuration`
- supported algorithms can be found from `id_token_signing_alg_values_supported` key of the object
- possible examples:
	- **HS class** - HS256, HS384, HS512
		- HMAC with SHA-NNN
		- in HMAC the signature is created using a secret key, 
		- this is less secure than public-private key
	- **RS class** - RS256, RS384, RS512
		- RSA with SHA-NNN
		- RSA is a public-private-key encryption algorithm
		- RSA is very popular and is widely used
	- **ES class** - ES256, ES384, ES512
		- ECDSA with SHA-NNN
		- ECDSA - Elliptical curve digital signature algorithm 
		- ECDSA is a public-private-key encryption algorithm
		- it is newer and is supposed to give better performance than RSA
- In general, using a public and private key for signing is more secure than using a secret key

The application after receiving the encoded JWT token, it's going to do the exact reverse
- i.e. first it will do url-decode
- then base-64 decode
- then verify the signature using the public key of the authorization server

Verification of a digitally signed JWT token gives the application an assurance that 
- the JWT token was sent from the authorization server itself and no one else could have signed the token except for the authorization server
- and no one else has tampered the token on it's way

There's no one stopping from proceeding ahead to read the header and payload, without verifying the signature

Thus, confidential information like passwords should not be put inside a JWT token payload

## JWT can also be used for OAuth client authentication, not just for resource owner or end user

- Client authentication happens using the Client Secret
- When we hit the 
	- OAuth Token endpoint (`/token`) during the OAuth flow to get the access token
	- or the OAuth Introspection endpoint (`/tokeninfo`) to check the status of the access token
	- or the OAuth Revoke endpoint
- we send a client secret along with client id in the POST request payload
- The client secret is passed along so that the authorization server would be able to verify for sure that the request is actually coming from the valid registered client
- The client id and client secret are given to the client itself by the authorization server
- the authorization server uses these **client credentials** (client-id and client-secret combined) to authenticate the client

