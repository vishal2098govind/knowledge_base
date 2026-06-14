VPC spans across multiple AZs, so VPC a **regional** level entity
Subnets are **zonal** entities
Internet Gateway is **regional** level entity scoped to a **VPC**
NAT Gateway is **zonal** level entity (updated later in 2025 Re:invent to **Regional NAT Gateway**)
ELB are regional level entity and have **VPC level** scope
DynamoDB and S3 have **regional** level scope
EC2 and RDS have subnet or **zonal** level scope
Route53, IAM, Billing Service have **global** scope. No particular **region**.

![[Pasted image 20260614092208.png]]
![[Pasted image 20260614092722.png]]
## Simplest way to host an application in AWS
![[Pasted image 20260614075918.png]]

### For Higher Availability
![[Pasted image 20260614080204.png]]

## Introducing Databases to architecture
DBs are not expected to need outbound traffic, thus better to have separate private subnets for them
![[Pasted image 20260614080511.png]]

## Allowing Outbound Traffic to application servers
![[Pasted image 20260614081018.png]]

### Accessing AWS Services from within Application Service within VPC
NAT Gateway charge will be seen if more data exchange happens
![[Pasted image 20260614081412.png]]

### VPC Endpoint (Gateway and Private-Link)
![[Pasted image 20260614081833.png]]
- Better than using NAT Gateway and Internet Gateway, if the AWS Service is in same region
- There are two types of VPC Endpoints
	- VPC Gateway Endpoint
		- Allows accessing S3 and DynamoDB privately
		- Need not use NAT Gateway
	- VPC Interface Endpoint (Private Link)
		- Can use different Interface Endpoints for different AWS Services privately
		- Can also be used to connect to applications in other VPCs. For that, other service has to expose the application through **Network Load Balancer**
![[Pasted image 20260614082223.png]]

## VPC beyond single VPC or region
- VPC Peering connection
- Transit Gateway
- Cloud WAN
- VPC Lattice
