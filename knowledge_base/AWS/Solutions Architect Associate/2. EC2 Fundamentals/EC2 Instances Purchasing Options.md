#aws #aws-ec2 #cloud-billing

## On Demand Instances
- Allow us to run instances on demand
- Good for **short workloads**, get predictable pricing
- Pay for what we use:
	- For Linux or Windows : Pay by seconds, after first minute
	- All other OS : billing per hour
- Has **highest cost** but no upfront payment
- No long-term commitment
- Recommended for 
	- short-term 
	- and un-interrupted workloads
	- with unpredictable application behavior

If we have different kind of workloads, we can optimize the discounts and pricing, by specifying it to AWS

For example:

### Reserved Instances (1 & 3 years)
- **Reserved Instances** - Long workloads
	- If we know we're going to run a database for a long time, then a reserved instance is great
- **Convertible Reserved Instances** - 
	- If we want to have a flexible instance type, for e.g. if we want to change the instance type over time, use convertible reserved instances
- We reserve a specific instance attribute 
	- Instance Type 
	- or Region 
	- or Tenancy 
	- or OS
- Specify Reservation Period
	- 1 year (+ discounts)
	- or 3 years (+++ discounts)
- Payment Options
	- No Upfront payment (+ discounts)
	- Partial Upfront (++ discounts)
	- All Upfront (+++ discounts)
- Reserved Instance's Scope
	- Regional
	- or Zonal (reserve capacity in a particular AZ)
- Recommended for steady-state usage applications (think **database**)
- Can by and sell our reserved instances in the **Reserved-Instance-Marketplace** if we don't need them anymore

#### Convertible Reserved Instances
- Can change the EC2 instance type or instance family or OS or scope or tenancy
- Thus bit less discounts than pure reserved instances

### Savings Plan (1 & 3 years)
- More modern
- Instead of committing to a specific instance type, we commit to a specific amount of usage in dollars thus **long workload**
- Discounts based on long-term usages
- Commit to certain type of usage like $10 per hour for 1 to 3 years
- Any usage beyond EC2 savings plans is billed at On-demand price
- Savings plans are **locked to a specific instance family and AWS region**
	- Flexible across:
		- Instance size
		- OS
		- Tenancy (Dedicated Host or Instances or Default)
- Meant for very short workloads
- Very cheap, 
- Unreliable, can loose instances at any time
- Most aggressive discounts of up to 90% compared to On-demand pricing
- Instances that we can loose at any point of time
	- when our max price set for the spot instance is less than the current spot instance price, then we loose that instance
- The **most cost-efficient** instances in AWS
- Very helpful for workloads that are resilient (able to withstand) to failure
	- Batch Jobs
	- Data analysis
	- Image processing
	- Any kind of distributed workloads
	- Workloads with a flexible start and end time
- Not suitable for critical jobs like databases

### Dedicated Host
- Book an entire physical server
- Control instance placement
- We get EC2 instance capacity fully dedicated to our use case
- Allows us to address **compliance requirements**
	- Use our existing server bound software licenses
- Purchasing Options:
	- **On-demand** - pay per second for active Dedicated Host
	- **Reserved** - 1 or 3 years (No Upfront, Partial Upfront, All Upfront)
- The **most expensive** option
- Useful when we have core/socket specific **licensee**
- [[Default Tenancy vs Dedicated Instances vs Dedicated Hosts]]

### Dedicated Instances
- No other customer will share our hardware
- Instances run on hardware that's dedicated to our account
- May share hardware with other instances of same account
- No control over instance placement
- Useful when we have **compliances to address**
- [[Default Tenancy vs Dedicated Instances vs Dedicated Hosts]]

### Capacity Reservation
- Reserve **On-Demand** instances capacity in specific AZ for any duration
- You always have access to EC2 capacity when we need it
- **No time commitments** (create/cancel anytime), no **billing discounts**
- To get billing discounts, can combine with regional-reserved instances and savings plan
- Charged at On-Demand rate, whether or not we run instances
- Suitable for short-term, uninterrupted workloads that needs to be in a specific AZ

## Which purchasing option is right for me
### On Demand
Coming and staying in resort, whenever we like, we pay the full price

### Reserved
Plan ahead 
and if we plan to stay for a long time, we may get a good discount

### Savings Plan
If we are sure we are going to spend a specific amount in the resort, maybe $300 a month for the next 12 months and stay in any room type (king or suite or sea view etc) i.e. use any instance size

### Spot instances
When the hotel runs very last-minute discounts because they have empty rooms and they want to attract people
People bid on getting this empty room as they get very high discounts
People may get kicked out at any time if someone else is willing to pay more for the room that we wanted or are using

### Dedicated Hosts
Book entire building

### Capacity reservations
Book a room for a period with full price even if we don't end up staying in it

## EC2 Spot and Spot Fleet Instances
### Spot Instances
- Can get discount of up to 90% compared to On-demand
- Define max spot price and get instance as long as current-spot-price < max-spot-price
- If current spot price > max spot price, you can choose to stop or terminate spot instance with a 2 minutes grace period
	- **choose stop when**: may choose to restart spot instance when spot price goes below max spot price, and continue where we left off
	- **choose terminate when**: if we don't need the state of EC2 instance on which we were on, we can choose to terminate the instance and let it go, assuming anytime we would restart, we will start on a fresh new EC2 instance
- **Spot Block**: If we don't want spot instances to be reclaimed by AWS. Block the spot instance for a specified time-frame (1 to 6 hrs)
	- In rare situations the instance may be reclaimed though
	- It is no longer available in AWS

### Persistent Vs One time Spot Instances
![[Pasted image 20260411211502.png]]
- If we want to cancel a spot instance request, it needs to be either in **open** state or **active** state or **disabled** state
- Cancelling a spot request does not terminate any spot instances that are running. It is still our responsibility to terminate them
- To terminate spot instances for good, we first have to cancel the spot request, or else if the spot request was persistent, it would go ahead and launch another spot instance

### Spot Fleet
- Ultimate way to save money
- Spot fleet is a way to define a set of spot instances and (optionally) On-Demand instances
- A Spot Fleet will try to meet the target capacity with price constraints
	- We can define possible launch pools : instance type (`m5.large`), OS, AZ
	- Can have multiple launch pools, so that fleet can choose the best and most appropriate launch pool for us
	- When spot fleet reaches max capacity or max cost (our budget we set), it will stop launching instances
- Strategies to allocate spot instances:
	- **lowest price**: from pool with lowest price (**cost** optimization, **short** workload)
	- **diversified**: distributed across all pools (great for **availability**, **long** workload)
		- since if one pool instance(s) goes away, other pool instance(s) are still active
	- **capacity optimized**: pool with optimal capacity for no.of instances
	- **price capacity optimized**: first select pools with highest capacity available, then within that, select pool with lowest price (best for most work loads)
### Spot Vs Spot Fleet
- Spot Request – where u know exactly the type of instance u want and AZ u want
- Spot Fleet – where we give AWS to choose from all these instance types and all these AZs, and we need AWS to choose such that we get lowest price



| Purchasing Option                  | Savings vs OD        | Commitment                           | Interruptible      | Ideal Use Cases                                                                                                                                                                                                                                                                                                  | Watch Out For                                                                                           |
| ---------------------------------- | -------------------- | ------------------------------------ | ------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------- |
| On Demand<br>(Flexible)            | Baseline (0%)        | None — pay per second                | No                 | - Unpredictable / spiky workloads<br>- Dev & test environments<br>- Short-term experiments<br>- First-time sizing a new app                                                                                                                                                                                      | Most expensive model. Easy to leave instances running and forget — use billing alerts.                  |
| Reserved · 1yr<br>(Committed)      | ~40%                 | 1 year (all/partial/no upfront)      | No                 | - Steady-state production workloads<br>- Databases (RDS/EC2 hosted)<br>- Backend services with predictable load                                                                                                                                                                                                  | Locked to specific instance family + region. Wasted if workload changes.                                |
| Reserved · 3yr<br>(Committed)      | ~60%                 | 3 years (all/partial/no upfront)     | No                 | - Core infrastructure (auth, data plane)<br>- Highly stable, long-lived services<br>- Significant cost reduction goals                                                                                                                                                                                           | 3-year lock-in is risky if architecture evolves. All-upfront saves most but ties up cash.               |
| Savings Plans<br>(Flexible Commit) | Up to 66%            | 1 or 3 years ($/hr spend commitment) | No                 | - Orgs that want RI-level savings with flexibility<br>- Mixed instance families or regions<br>- Compute SP also covers Lambda & Fargate<br>- Savings plans are **locked to a specific instance family and AWS region**<br>- Flexible across: Instance size, OS, Tenancy (Dedicated Host or Instances or Default) | Commit to a $/hr floor — usage beyond that bills at OD rates. Easy to under-commit.                     |
| Spot Instances<br>Interruptible    | Up to 90%            | None                                 | Yes — 2 min notice | - Batch processing & data pipelines<br>- ML training jobs (checkpointed)<br>- CI/CD runners<br>- Stateless workers behind a queue<br>- Load test fleet                                                                                                                                                           | Can be reclaimed anytime by AWS. App must handle interruption gracefully (checkpoint, drain, retry).    |
| Dedicated Hosts<br>Compliance      | ≈ OD or higher       | On-demand or 1/3yr reservation       | No                 | - Bring Your Own License (BYOL) — Windows Server, SQL Server, Oracle<br>- Regulatory / compliance mandates (physical isolation)<br>- Full visibility into socket/core/host topology                                                                                                                              | Most expensive option. Overkill unless BYOL or compliance explicitly requires it.                       |
| Dedicated Instances<br>Isolation   | ~10% premium over OD | None (or Reserved pricing)           | No                 | - Hardware-level tenant isolation (no shared hardware with other AWS accounts)<br>- Compliance requiring dedicated hardware without needing full host control                                                                                                                                                    | No visibility into physical host. Can't use per-socket/core BYOL licenses — that needs Dedicated Hosts. |
| Capacity Reservation<br>Guaranteed | 0% (billed at OD)    | None — cancel anytime                | No                 | - Guaranteed capacity in a specific AZ for disaster recovery<br>- Planned traffic spikes (product launches, live events)<br>- Combine with Savings Plans / RIs for cost savings on reserved capacity                                                                                                             | You pay OD rate whether or not instances are running. Pure capacity guarantee, not a discount tool.     |