## Security Groups
Security Groups are the fundamental of network security in AWS
• They control how traffic is allowed into or out of our EC2 Machines.
• It is the most fundamental skill to learn to troubleshoot networking issues
![[Pasted image 20260614201144.png]]

They regulate: 
	• Access to Ports 
	• Authorized IP ranges – IPv4 and IPv6
	• Control of inbound network (from other to the instance) • Control of outbound network (from the instance to other)
• Security groups are **stateful** 
- If an outbound traffic is allowed, it's corresponding response will also be allowed automatically without any explicit mention in the inbound traffic rule
- vice versa as well
• You can reference another Security group as source
![[Pasted image 20260614201251.png]]

**Security Group Good To Know**
• Can be attached to multiple instances
• Locked down to a Region / VPC combination
• Does live “outside” the EC2 – if traffic is blocked the EC2 instance won’t see it
• If your application is not accessible (time out), then it’s a security group issue
• If your application gives a “connection refused“ error, then it’s an application error or it’s not yet in running state
• All inbound traffic is blocked by default
• All outbound traffic is authorized by default

## Network Access Control List (NACL)
• Works at Subnet level – Hence automatically applied to all instances
• **Stateless** – We need to explicitly open outbound traffic
- Here, if an outbound request is coming into our EC2 instance http server, we need to mention both inbound and outbound rule
	- Inbound: Source port = 80, destination port = 1.2.3.4 (if we want to accept from a specific IP)
	- Outbound: Destination port = ephemeral-port-range, Source port = 80
		- Ephemeral port range because here the client uses ephemeral port to connect
		- if our EC2 instance was originating request to outside web server, then it would use ephemeral port to 80
![[Pasted image 20260614202616.png]]
• Contains **both Allow and Deny rules**
• Rules are evaluated in the order of rule number (1 to 32766) and **stops at first rule match**
• Default NACL allows all inbound and outbound traffic
• **NACL are a great way of blocking a specific IP at the subnet level**
- With SG, we cannot **Deny** inbound from a particular IP if we accept inbound from all sources (0.0.0.0/0)
- **SG is only an Allow List**
- with NACL we can have an entry for Deny with rule number lesser than the rule for Allow
- Thus, it is preferred to start rule numbers not from 0 but from something like 100 and try to have rule numbers with orders of 10s or 100s so that we have enough buffer in between to add dedicated Deny rules in between

| \#Rule | Type <br>   | Protocol | Port | Source            | Allow/Deny |
| ------ | ----------- | -------- | ---- | ----------------- | ---------- |
| 100    | ALL Traffic | ALL      | ALL  | 180.151.138.43/32 | DENY       |
| 101    | HTTPS       | TCP      | 443  | 0.0.0.0/0         | ALLOW      |
| *      | ALL Traffic | ALL      | ALL  | 0.0.0.0/0         | DENY       |

