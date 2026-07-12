- Amazon DynamoDB is a fully managed, serverless NoSQL (key-value) database service
- **Consistent**, 
	- **single-digit millisecond read and write performance** at any scale 
	- **millions** of requests/seconds, 
	- **trillions of row**, 100s of TB of storage
- **Serverless**: No provisioning or capacity management, Encryption for data at rest, 99.999% availability
- Support different data types i.e. Scalar (Number, String etc.), Document (JSON) and Set types
---
DynamoDB has Table as basic entity to store data
- Table contains items (like rows)
- Items contains Partition Key, Sort Key (optional), and attributes (key-value pairs).
- **Partition Key + Sort Key = Primary Key**
---
### Query and Scan
**Query** - fetches specific items using Partition Key
- Retrieves items by Partition Key (and optional Sort Key conditions).
- Can use conditions like <, >, <=, >= etc. with Sort Key
- Fast and **efficient** — looks only in matching partitions.
- Can use filters, but filtering happens after reading the data.

**Scan** – fetches every item of the table
- Reads every item in the table or index.
- Slow and **expensive** for large tables.
- Can use filters to narrow down the results but still scans all items first.
---
### DynamoDB - Global & Local Secondary Index
![[Pasted image 20260712101037.png|309]]
**Local Secondary Index**
- Queries data over a **single partition** only (localized)
- Supports both **eventual consistency** or **strong consistency**
- Can only be added at the same time that you create the base table
**Global Local Secondary Index**
- Queries data across **all partitions** of the entire table
- Only supports eventual consistency only, no strong consistency
- Can be added or deleted at **any time**
![[Pasted image 20260712102914.png|295]]
---
### DynamoDB Features
- Can add ACID properties to a DynamoDB table
- DynamoDB Table class - Standard, Standard-IA
- DynamoDB Read/Write Capacity (RCU/WCU)
    - Modify read/write throughput capacity
- DynamoDB Global Tables
    - Can configure DynamoDB as a multi-region database
- DynamoDB Accelerator (DAX)
    - Can reduce response time from miliseconds to microseconds through caching (DAX)
- DynamoDB Streams
- DynamoDB TTL
    - Automatically expire items(rows) based on timestamp
- DynamoDB Backup
- DynamoDB export / import
---
### DynamoDB Table Class
DynamoDB supports Standard Table class and Infrequent Access (IA) Table class
**Standard (default)**
- Best for frequently accessed data.
- Higher storage cost, lower read/write cost.
- Ideal for hot or regularly used items.
**Standard-IA (Infrequent Access)**
- Lower storage cost, higher read/write cost.
- Designed for rarely accessed or cold data.
- Useful for archival, history tables, old session data, etc.
- We can change the table class from Standard to IA
---
### DynamoDB Read and Write capacity
- We can define the DynamoDB Read and Write capacity as per the application need.
- There are two modes to set the capacity:

**Provisioned Capacity Mode (Default)**
- You pre-define RCU/WCU for the table.
- Best when you have **predictable** or steady traffic.
- Cheaper than on-demand if your workload is consistent.
- **Risk of throttling if traffic suddenly spikes above provisioned limits.**

**On-Demand Capacity Mode**
- No need to set RCU/WCU - DynamoDB **scales automatically**.
- Ideal for **unpredictable** or **spiky** workloads.
- You pay per request, which may be **costlier for steady high throughput**.
- Virtually no throttling unless you hit account-level limits.

---
### Multi Region Database - DynamoDB Global Tables
- provide multi-region replication across all replicas/tables
- Can read/write to any replica
- includes ongoing data changes across all the tables during replication

- Multi-region, multi-active, serverless tables across regions
- completely automated

- underlying infra is entirely managed by AWS themselves
- doesn't reside within a custom VPC
- 99.999% availability

**Use cases:**
- Global application requiring low latency access for users
- Can handle region level failure (DR - Disaster Recovery)
---
### DynamoDB Streams
- A data stream that **captures each and every data change made to the items**
- If an item was added, modified, or deleted, then that item will be included in the DynamoDB stream
- **Can be associated with AWS Lambda**. The function can poll the stream and execute a set of actions whenever it detects new stream records
- Can also be integrated with **Kinesis Data Streams**
- Important component that **needs to be enabled when using Amazon DynamoDB Global Tables**
![[Pasted image 20260712105301.png|134]]
---
### DynamoDB Accelerator - DAX
- Fully managed highly available in-memory cache for DynamoDB
- 10x performance improvement with single digit millisecond to microsecond level latency
- API-compatible with DynamoDB. Only client & endpoint needs to change to use with an existing application.
- DAX provides access to eventually consistent data from DynamoDB tables