#k8s #health-checks

- k8s uses `healthchecks` to realize there's something wrong with the new update and it should pause doing the rolling update
- there are multiple `healthchecks` in k8s
- there's one particular type of `healthcheck` which helps realizing to pause a rolling update in case of a buggy update
- there are 3 kinds of `healthchecks` (aka "probes")
	- `startupProbe`
	- `readinessProbe`
	- `livenessProbe`
- `healthchecks` are completely optional to put
	- can put none of them, 
	- all of them 
	- or some of them
- These `healthchecks` can then be implemented with different techniques
	- can have a HTTP Request to check if a web server responds correctly
	- can be a script or a program
- if `healthchecks` are not defined, it will default to "success" result always

## Use-cases of health-checks
- Why do we have different kinds of health-checks
### `startupProbe`
- Used for containers which take long time to start
### `readinessProbe`
- for containers that sometimes are overloaded and we need to do some re-balancing
### `livenessProbe`
- for broken containers which can only be fixed by a restart
## Liveness probes
- The container is dead, we don't know how to fix it, other than restarting it
- goal of liveness probe is to detect if a container is dead
- if a container is dead, k8s will restart it
- with default parameters, k8s takes 
	- up to 30 seconds to detect if a container is dead
	- up to 30 seconds to terminate the dead container

### Liveness probes gotchas
- Should only be used for problems that actually go away on restart
	- If we have a problem that is completely un-related, like disk-full, even on restart we end up getting the same issue, and we will be restarting in a loop which will not help at all
- Service dependencies
	- When we have a bunch of web-frontends and a DB server
	- If the liveness probe makes a request to a database to conclude the web server being down if the database is down
	- this is a big mistake
		- if the database really fails, all the liveness probes on the web server are going to fail, and k8s will restart all the web servers or frontends
		- and if the database really fails, restarting the web server will not help
		- thus, again end up restarting in a loop, which will not help at all
			- in the best case scenario, restarting in a loop doesn't change anything
			- in the worst case scenario, we are going to add a lots of extra load on the platform because of all the containers being restarted
			- if we're really unlucky, it might prevent the recovery of the dependent service just because of repeated restarts and liveness prob repeatedly talking to the dependent service
		- thus, `liveness` probes should only take care of the local containers or pods, without external dependencies in those `liveness` probes, otherwise we have **cascading failures**
- `Liveness` probe needs to respond quickly
	- Not just `liveness` probe, but all probes should respond quickly
	- but it's particularly important for the liveness probes, to respond quickly
	- quickly means
		- default probe timeout is 1 second
		- if we have a liveness probe that takes more than one second, that's failure for k8s, and k8s will eventually restart the container
			- e.g. if we implement liveness probe using a HTTP request, and it takes more than 1 second to respond, that's a failure, even it responds with 200 OK

## Readiness probes
- Sometimes our container "needs a break"
- When the container needs a break, it can use the readiness probe to indicate "stop, I am not ready, stop sending requests to my way"
### Use cases when a container needs a break:
#### Ambitious SLO with GC language
- E.g. 1.: A service where the SLO (Service Level Object) on the latency is very low, very ambitious, for instance, this is maybe a service, that is part of a NAT platform or search engine or whatever, and we want to serve every request within 10ms, which is pretty low. 
	- If the service is written in maybe Go, or Java or any garbage collected language, once in a while we will notice that there's a huge spike in latency because of the garbage collection (GC) happens
	- one way to work around this is to trigger GC deliberately
	- instead of letting the GC happen whenever the runtime decides, we can kind of 
		- orchestrate 
		- and stop receiving requests (denial of service)
		- trigger the GC
		- and when we're done, we can accept requests again
#### Cache reload
- E.g. 2: Reload a bunch of data
	- If the container is serving data from an in-memory cache
	- once in a while we need to reload the cache
	- when we reload the cache, it's going to interrupt the service for few seconds
	- we can use readiness probe and reload cache during that time, asking the load balancer to send requests to somewhere else to stay available
#### Mark itself ready/not ready
- E.g. 3: If we have connected service for serving videos or video conferencing or gaming, something where we have connection from client to the server
	- we have benchmarked our service and realized that we can do about 100 connections per pod
	- and if we expected to have around 250 clients, and we have 3-4 pods to be on safe side, the connections are going to be load-balanced among pods
	- imagine, if we have spike in traffic
	- if we have 4 pods, but now we have more clients connecting around 300-350, we realize to scale up real quick
	- so now, we have 5 pods, out of which 4 are super busy with almost 100 connections each, and the 5th pod is completely empty and thus we want to forward new connections to 5th pod. It can be achieved through readiness probe
	- we can set a readiness probe, so that when the pod has more than 100 connections, it's going to **mark itself not ready**, indicating the load balancer to stop sending more requests to my way and send somewhere else 
	- when the pod is not ready, it will be removed from the load balancer, so that the new connections will be sent to the other pods, rebalancing the load and connections
	- over time, if the number of connections drop below 100, the container or pod is going to mark itself as ready again, and is going to accept connections again
	- ~ to going to an administration with around 10 desks serving users or customers
		- can imagine a queue in front of each desk
		- and when we have more than 10 people in a queue
		- then, we can imagine a little traffic signal at the entrance of each queue or desk with 
			- a green light, lit when we have less than 10 people in the queue, indicating can accept more people
			- a red light, lit when we have more than 10 people in the queue, indicating not to send more people this way
	- when the container marks itself as not-ready, what happens to the existing connections?
		- they stay as they are, everything is fine for them
		- no existing connections are broken

## Startup Probes
- When containers take a long time to start and that time is hard predict exactly how long
- startup probe checks if the container is done booting
- the container starts and then the startup probe is going to continuously check on the container, and as soon as the startup probe works, k8s is now going to start the liveness probes and readiness probes etc
- when we define the startup probes, we should adjust the timings for the startup probe
- if we don't adjust timing and the container fails to start within 30 seconds, k8s terminates the container and restarts it
- when we define all these probes, for each probe, we have bunch of parameters that can be customized
	- `periodSeconds` - How often - by default it's every 10s
	- `timeoutSeconds` - by default it's 1 second
		- if it takes more than 1 second, it's a failure
		- if we have a web-server that takes a bit of time while booting up, it could be a good idea to increase the `timeoutSeconds` parameter of the liveness probe so that the container is not terminated
	- `successThreshold` or `failureThreshold` 
		- indicates how many successes or failures in a row do we need to consider that the probe succeeded or failed
		- for `successThreshold` - default is 1
			- which means that as soon as the probe succeeds, k8s concludes that everything is fine
		- for `failureThreshold` - default is 3
			- which means that we need 3 failures in a row, for k8s to conclude that there's a problem and k8s restarts the container gracefully
				- i.e. we don't burn the container during restarts and restart it
				- we shut down the container nicely, with a **signal**
				- if we have a container that is well behaved, i.e. sends signal during shut down, and not completely crashes, it might react pretty quickly to the signal and shut down and everything is fine
				- if we have 
					- either a container that doesn't handle signals properly
					- or the container is in such a bad state that it doesn't react to the signal, which could happen, then we have to wait for the **grace period** (30s by default) before k8s forcefully deletes the container
	- when the startup probe fails, we have the same behavior as the liveness probe
	- if the startup probe fails three times, k8s is going to restart the container, concluding that it's taking too long to start
	- thus, when we define the startupProbe, we need to increase the failureThreshold accordingly

### Note
- We don't always need to use all of the three available probes together
- in fact, most often times we will be fine with just a readiness probe, and sometimes also a liveness probe


## Different types of probes
### exec
- arbitrary program execution
- Runs a program inside the container
- `kubectl exec` or `docker exec`
- the script for the program that we want to use needs to be present in the image
- the only thing k8s looks at is the status code
- it doesn't care about the message or output like std:out or std:err etc
- a good example is when we have a container with something like a worker, unlike a HTTP/Web server where we would often think of having a health check endpoint and very often that's the easiest thing to do with a web server for implementing `healthchecks`
	- when we have a worker, to check things if the worker is really working or if it's crashed
	- we can use a script to be executed to indicate what's going
### httpGet
- HTTP GET request
- we indicate the port and optionally the path to which we want to make the request
- k8s only looks at the status code of the response
- if it's something in 200s or 300s, it's a success, else it's a failure, so 404 is a failure
- if we want to implement something fancy like we want to make some request to the API and the result is going to be some JSON, and we want to extract a particular field in the response JSON, to know if the service is healthy, in that case, we are going to an `exec` probe
- something like
```sh
# here, we need to have curl, jq and grep available inside the container image, so that it can be executed using exec probe
exec:
	command:
	- sh
	- -c
	- "curl http://localhost:5000/status | jq .ready | grep true"
```
### tcpSocket
- check if a TCP port is accepting connections
- There's no additional check
- It's pretty easy to have a service that accepts connections, but is still broken
### `grpc`
- standard gRPC Health Checking Protocol


### Best Practices for `Healthchecks`
- Readiness probes are almost always beneficial
	- It almost never hurts the application
	- it's almost always a good idea to put a readiness probe
	- the only case where it would be a bad idea is when the readiness probe is expensive
		- e.g. if we have a content management system or a news site
		- in a particular article page, it would just load that particular article, but in a home page it might make a bunch of requests to display a lot of or a variety of content.
		- so that might make a bunch of requests to the database and that could be pretty expensive, although mostly such home page is going to be cached by the CDN or some other caching mechanism
		- in such case, having a specific api route for health is suggested like `/ping` or `/health`
- Be more careful with liveness and startup probes as we don't always necessarily need them
	- we should first get some production experience with that service before adding the liveness probe
		- ahead of time, we might not know exactly what are the failure thresholds or failure modes of our service
		- instead we want to run it in production, get some experience there, and maybe within the first week, we will notice, almost of every night, the service crashes, we need to restart it manually and when we troubleshoot it, we notice a particular thing doesn't work
		- we will thus use that particular thing that we noticed, in our liveness probe, where we want our liveness probe to identify exactly the problem that we have, we don't want it to be overly broad, because we don't want it to accidentally restart our containers or services
	- if we have a container that starts within 10-15-20s we don't need a startup probe
- Readiness and liveness probes should be different
	- if they are same, it's not good
	- we want to have a two step response when there is a problem
	- in first level, we might realize the container is a bit overloaded - readiness probe
	- in next level, we might realize that the container is completely broken and destroyed and we need to create a new one - liveness probe
		- thus, the threshold for liveness probe should be higher