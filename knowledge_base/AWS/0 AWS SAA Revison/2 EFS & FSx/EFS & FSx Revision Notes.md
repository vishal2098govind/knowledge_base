
## When to use EFS?
- For applications that require a **shared file system** that can be accessed by multiple computers at the same time
- For applications that access data using the standard file system interface provided by the underlying OS
- Amazon EFS provides file systems storage for **Linux workloads** only specifically
- EFS is completely managed by AWS
    - Thus removes complexity of
        - Deploying
        - Patching
        - and Maintaining the underlying infrastructure of the file system
- Can scale EFS without provisioning new physical servers
- Integrates with AWS KMS for data at rest and in-transit

### Availability and Durability modes
- Regional (EFS standard) - redundantly across multiple azs - low RPO (DR)
- One Zone - redundantly across single az - high RPO

### Performance and Throughput modes
**Performance mode**
- Can choose among 
    - **General purpose** 
    - and **Max I/O** - for Big data analytics and media processing workloads
    - **How to decide** when to use Max I/O performance mode:
        - Start with general purpose
        - run and test performance
        - If **PercentIOLimit** metric of file system is near 100% for extended amount of time, then we should switch to Max I/O performance mode
    - **How to switch** from general purpose to Max I/O performance mode?
        - Create a new Max I/O EFS System
        - Migrate all data from general purpose EFS to this using **AWS DataSync service**

**Throughput mode**
- Can choose among
    - **Bursting** 
        - default
        - Allows EFS to scale size of File system, in Regional (standard) or one-zone, grows
        - For **unpredicatble** access patterns and eventual EFS size
        - Throughput scales with size of underlying EFS
    - **Provisioned**
        - Fixed throughput
        - If **amount of data in EFS is less relative to throughput demands**, like in development environment, can use provisoned mode
        - mainly for **known access patterns**

## FSx for Lustre

### High Performance Computing (HPC)
- It's like your everyday computing, only it's more powerful
- **Aggregating computing power** to deliver much high "horsepower" than what traditional computers and servers can offer
- This solves the problem of **processing huge volumes** of data **at very high speeds** using **multiple computers and storage devices**
- Applications like following need HPC:
    - Machine learning
    - Analytics
    - Genomics
    - Video processing
    - Financial modelling
    - Electronic design automation

### Amazon FSx for Llustre
- Amazon FSx for Lustre is a fully managed service that provides **high-performance and scalable file system storage for HPC workloads**
- It's a managed service, thus eliminates complexity of setting up and managing Lustre file systems, allowing us to spin up resources within minutes from when we configured to deploy it
- It's basically a **High-performance-File-System**
- **Provides a native file system interfance** that works as any file system does with Linux operating system

### Amazon FSx for Lustre features:
- Offers 
    - **sub-millisecond latencies**
    - up to **100s of GBps of throughput**
    - and **millions of IOPS**
- It's POSIX compliant
    - We can continue using current linux-based applications without having to make any changes
![[Pasted image 20260713053758.png|400]]
- Can link to data repositories like **Amazon S3 or on-premises data store**.
    - i.e. we can burst our data processing workloads from on-premises data center into AWS by importing data using AWS Direct Connect (DX) or AWS VPN
![[Pasted image 20260713054139.png|400]] ![[Pasted image 20260713054207.png|400]]
- Can use SSD or HDDs to store data
- Integrates with AWS KMS to allow **encrypting data at rest** using customer managed keys
- **Data in transit** is automatically encrypted in **certain AWS Regions only**, when accessed from **supported EC2 instances**


## Amazon FSx for Windows File Server
- One of the most common setups in an organization is having an infrastructure where employees can **share files** through a server
- In an on-premises, we typically have a windows server provisioned, where users can read and write files via **SMB protocol**
- **SMB - Server-message-block**, a network protocol that allows users to communicate with **remote servers** to **share, open and edit files**
- Amazon FSx for Windows File Server brings this infrastructure to AWS Cloud
- It's a fully compatible **shared Windows File Service**
- Amazon FSx file system can be accessed from **Windows machine, Linux, macOS, EC2, VMWare, ECS, EKS** and even **Lambda functions**
- Can **increase size** of file storage as much as we want. But we **cannot decrease storage capacity** once increased
- Supports Single AZ and Multi-AZ deployments
    - **Single AZ1 deployment** 
        - Supports only SSD, no HDD
        - Can replicate data using Microsoft Distributed File System Replication **(DFSR)**
        - For High Availability provison Single AZ1 FSx for windows deployment in two seperate AZs and replicate data syncronously via DFSR
    - **Single AZ2 deployment**
        - Supports both SSD and HDD
        - **Can't use (DFSR)**
        - Can't get High Availability
    - **Multi-AZ deployment**
        - Supports both SSD and HDD
        - **Can't use (DFSR)**
        - **Highly Available**
        - Automatically replicate data across primary and secondary file servers in multiple AZs
- Other features:
    - Native Amazon FSx backups
    - Creating backups throughp **AWS Backup service**
    - Windows shadow copy for easily undo file changes and compare file versions
    - **Schedule replication on AWS DataSync** to second file system
    - Integrates with **Microsoft Active Directory**
    - Upon creating a file system with Amazon FSx Windows File Server, can join with on-premises existing **Microsoft Active Directory** by **using AWS Managed Microsoft Active Directory (AD)**
	- ![[Pasted image 20260713061444.png|403]]
    - For migrating Files from on-premises to Amazon FSx for Windows File Server using **AWS DataSync 
        - a fully managed online-data-transfer service**
        - Can copy data over **internet or AWS Direct Connect**
