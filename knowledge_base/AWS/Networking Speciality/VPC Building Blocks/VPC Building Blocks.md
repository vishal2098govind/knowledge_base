- CIDR
- Subnets
- Route tables
- Internet Gateway
- Firewalls
	- Security Groups - at instance level within a subnet
	- NACL - at subnet level
- DNS - Route53 Resolver

### Subnets, Route Table and Internet Gateway
Subnet has **AZ scope**
Every VPC comes with 
- Local Route Table
	- Responsible for communication between and across subnets within VPC
- Main Route Table

#### Communication between instances in different subnets
If we have two instances, each in different subnet
E.g. VPC CIDR = 10.10.0.0/16
Subnet 1 = 10.10.0.0/24
Subnet 2 = 10.10.1.0/24

Instance A in subnet 1 has private IP = 10.10.0.15/32
Instance B in subnet 2 has private IP = 10.10.1.11/32

There is a local route table in the VPC and main route table in VPC

Main Route Table:

| Destination  | Target           |
| ------------ | ---------------- |
| 10.10.0.0/16 | Local            |
| 0.0.0.0/0    | Internet Gateway |
The main routing table is responsible for routing any IP packet bearing destination IP address within this range of 10.10.0.0/16

**Any instance in any subnet can by default communicate with any other instance in any other subnet**

By Default, every subnets and instances within subnets use main route table by default

### Internet Connectivity
For instances to connect to internet from within a VPC subnet, 
- we need internet gateway at VPC level
- instance can be public subnet. if it is in private subnet, we need NAT Gateway
- instance can have public IP address. if it has only private IP, we need NAT Gateway
- main route table must have entry of **0.0.0.0/0 as destination** with target as **Internet Gateway**


If we want only one subnet (Subnet-A) to have internet connectivity and other subnet (Subnet-B) to not have internet connectivity, we can use Custom Route tables for each subnets.

| Subnet - A Route Table | a public subnet  |     | Subnet - B Route Table | a private subnet |
| ---------------------- | ---------------- | --- | ---------------------- | ---------------- |
| **Destination**        | **Target**       |     | **Destination**        | **Target**       |
| 10.10.0.0/16           | Local            |     | 10.0.1.0/16            | Local            |
| 0.0.0.0/0              | Internet Gateway |     |                        |                  |
As soon as we attach a custom route table to a subnet, it no longer uses main route table

Route Table of a **public subnet** has an entry to Internet Gateway. E.g. we have NAT Gateway, Web Servers, Load Balancers in public subnets
Route Table of a **private subnets** have no entry to Internet Gateway. E.g. Databases, application server
