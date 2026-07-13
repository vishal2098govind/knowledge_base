## Key Terms

- Volume Size - The storage capacity of an EBS volume, measured in GiB (Gibibytes).
---
- **IOPS** (Input/Output Operations Per Second)
    - Measures how many read and write operations per second a volume can handle
    - Indicates the speed and responsiveness of storage.
    - Important for transaction-heavy workloads (e.g. databases).
  - **Baseline IOPS**
      - The default guaranteed performance for your volume type (before bursting or provisioning).
  - **Provisioned IOPS** (PIOPS)
      - A feature where you explicitly set desired IOPS (up to 256,000).
---
- **Burst Performance / Burst Credits**
    - Temporary performance boost above baseline IOPS when needed.
    - **Volumes accumulate burst credits when idle and spend them under load.**
    - Helps small volumes handle occasional heavy I/O spikes.
---
- **Throughput (MB/s)**
    - Measures the amount of data transferred per second, typically in MB/s.
    - IOPS = “how many”; Throughput = “how fast”.
---
## Amazon EBS
- EBS stands for Elastic Block Store
- A type of a block storage like the Amazon EC2 Instance Store
- Its data is **more persistent** and will not get lost even if the EC2 instance was stopped, restarted, or terminated
- **Zonal in scope**, which means it only exists in a single Availability Zone
- Can be attached to any EC2 instances in the **same Availability Zone only**
- Can be encrypted at rest using AWS KMS
- You can attach one or more Amazon EBS volumes to a single EC2 instance

---
- Suitable for a variety of workloads such as 
  - databases, enterprise applications, 
  - big data analytics engines, 
  - file systems, media workflows, and others
- Allows you to store and retrieve your data with **high throughput and low latency**
- The Amazon EC2 **instance and its attached EBS volumes** are logically attached together and are both located **within a single Availability Zone**, which significantly **reduces latency**
- Since the underlying physical resources that power your Amazon EC2 instance and EBS volumes are located within the **same city or geographic area**, Amazon EBS is capable of providing **low latency read or write access to your data**

---
### Block Storage Vs File Storage Vs Object Storage
- Block Storage is an integral technology that mainly operates at hardware level

Meaning of "Block" in Block Storage:
- A block is a **sequence** of bytes or bits that represent our data
- When we store a file, our system splits it into multiple data blocks, each with same maximum length
- The maximum data length of a block is called "**block size**"
- If block size if 4KB, and our file is of 8KB, then we need two blocks to store the file
- This storage technology is present in almost all computers today
- We can know block size on a mac by using:
```bash
$ diskutil info / | grep "Device Block Size"
   Device Block Size:         4096 Bytes
```
 ![[Pasted image 20260712154756.png|400]] ![[Pasted image 20260712155000.png|400]] ![[Pasted image 20260712155030.png|400]]
- So a file is split into small 4KB blocks
---
**Storage Devices**
- Blocks operate on hardware components of storage device
- **Hard-Disk Drives (HDD)**, 
    - **the data blocks are scattered** on a **physical disk** - platters and sectors of HDD.
    ![[Pasted image 20260712155507.png|300]]
- **Solid-State Drives (SSD)**
    - the blocks are stored on the underlying **flash-memory chips**
	![[Pasted image 20260712155942.png|300]]
---
### RAID
- We can further improve performance of block storage by using multiple volumes and joining them using a RAID configuration - Redundant-Array-of-Independent-Disks
- This is a data-storage virtualization technology allowing to improve the performance and availability of storage devices
- Two types of RAID configurations - RAID-0 and RAID-1
    ![[Pasted image 20260712160959.png|517]]
---
### Storage Types in EBS
![[Pasted image 20260712161547.png]]

---
### EBS Snapshots
- An incremental backup that internally uses Amazon S3 to persist your data
- It only saves the data blocks that have changed after your most recent snapshot
- Allows you to restore the state of your EBS volume in the event of data loss
- Enables you to copy your EBS volume to another AWS Region for your **data migration**, **disaster recovery activities**
- **Can be used to encrypt an unencrypted Amazon EBS volume.**
- **Automate** the **creation**, **retention**, and **deletion** of your EBS snapshots and EBS-backed AMIs using the **Amazon Data Lifecycle Manager (Amazon DLM) service**
---
### EBS Encryption
- Uses **AWS KMS** Keys to encrypt EBS volumes and snapshots
- Feature "**Amazon EBS Encryption by Default**" can be enabled at a region-specific level manually.
    - This enables to automatically set the encryption settings of all of new EBS volumes that will be created and also copies of AWS EBS snapshots, there by eliminating the burden of encrypted all the volumes manually, ensuring **data security compliance all the time**
---
### EBS Volume Types
- Can attach multiple EBS volumes to a single EC2 instance, one of them being root volume with system image to boot instance
- Can mix types of EBS volumes attached to EC2 instance
- EBS volume types fall into one of two categories
    - SSDs 
        - optimized for **transactional workloads** involving **frequent I/O operations** with small I/O sizes
        - Performance Attribute = IOPS
    - HDDs - optimized for 
- Root volume should always be of SSD category and not HDD category

### SSDs Category
- This category has two types of EBS Volumes

#### General purpose SSDs (gp2/gp3)
- Provides **balance of price and performance**
- Recommended for most workloads
- Also suitable for apps with **unpredictable or unknown access patterns**
- Provides a **configurable and consistent** IOPS to allow you to accommodate the changes in your data storage requirements
- Suitable for **low-latency interactive** apps in production as well as your development and test environments
- For your infrequently accessed applications or systems that: 
    - Only peaks during certain times of the day
    - Has a varying Disk I/O operations
- Provides ample/enough IOPS for your applications but **not on par** with what a **Provisioned IOPS** type can give
- Still recommended for the **most cost-effective** storage option that **does NOT sacrifice performance**

#### Provisioned IOPS SSDs
- Primarily used for **mission-critical**, **low latency**, or **high-throughput workloads**
- Can also configure IOPS settings like gp2/gp3 SSD
- Provides **sub-millisecond** latency and **consistent IOPS** performance
- Allows you to set the amount of available IOPS of your EBS volume
- For hosting data to your applications that makes **small reads and writes** to a small file system
- For applications that require a number of **high read and write IOPS performance**
- For fixing latency issues
- For scenarios where your database storage performance is the bottleneck
- For storage systems that require a **configurable and consistent IOPS**
- Allows **EBS Multi-Attach**, allowing to attach same volume to multiple EC2 instances. Note that this is only supported for provisioned IOPS EBS Volume type, not on any other EBS Volume types
    - Although, all attached EC2 instances **cannot** still **concurrently** modify file content
- Considerably **more expensive** than general purpose volume type
![[Pasted image 20260712172548.png|400]]

---
### HDD Category
- Optimized for **large streaming workloads**
- For various types of applications and systems with large, **sequential** I/O operations
- Performance Attribute: **Throughput** (MB/s)
- Cannot be used as your boot (root device) volume

#### Throughput Optimized HDD
- A **low-cost** HDD designed for **frequently accessed**, throughput-intensive workloads
- Can be used for your Big data applications, Data Warehouses, and Log Processing

#### Cold HDD
- Lowest-cost HDD storage meant for **storing less-frequently accessed** data
- Most **cost-effective storage type** among all EBS volume types
- Suitable for throughput oriented data that is **infrequently accessed**
- Scenarios where we don't frequently access such data, but when accessed, we need it with high throughput
- Perfect for scenarios where the **lowest storage cost is of the utmost importance**
---
## EBS Anti Patterns
- If you just need a **temporary storage** for your data, use **EC2 Instance Store instead**
- If you have to store your application or system data in a POSIX-compliant **hierarchical directory structure** (use **Amazon EFS** instead)
- If you have multiple applications that are **concurrently accessing the same files** at the same time, it is better to use the **Amazon EFS or Amazon FSx** service instead
- If you need to store your static data in the **most cost-effective way to store any static data**, it’s more appropriate and cheaper to store them in **Amazon S3**


## Amazon Data-Lifecycle Manager (DLM) to manage EBS Snapshots and EBS-backed AMIs
- Can be done via customised serverless architecture using AWS Lambda and AWS EBS APIs, allowing to customize steps and rules for managing EBS Snapshots
- This will take time to develop/implement the architecture because we'll need to create lambda functions we need, set it up, and monitor
- Amazon Data Lifecycle Manager helps automate this
- Can enforce **regular backup schedule** as per organization's internal policy
- Reduce storage costs by **deleting outdated backups**
- implement disaster recovery backup policies that back up data to different AWS accounts

### DLM Steps
1. Create a policy in DLM
      - EBS snapshot policy
      - EBS-backed AMI policy
          - **DLM does not help backup AMIs which are backed by EC2 instance-store, but only those AMIs which are backed by EBS volumes**
      - Cross-account event policy
          - For automatically copying shared snapshots **across AWS accounts**
2. **Define tags** of EBS volumes we want to backup
    - DLM uses resource tags to identify resources backup
3. **Define schedule and frequency** to create backups
4. **Define retention period** for policy