- Amazon DynamoDB is a fully managed, serverless NoSQL (key-value) database service
- **Consistent**, 
	- **single-digit millisecond read and write performance** at any scale 
	- **millions** of requests/seconds, 
	- **trillions of row**, 100s of TB of storage
- **Serverless**: No provisioning or capacity management, Encryption for data at rest, 99.999% availability
- Support different data types i.e. Scalar (Number, String etc.), Document (JSON) and Set types
---
## DynamoDB Features
- DynamoDB Components - Table, Item, Partition keys, Secondary Indexes (GSI and LSI)
- DynamoDB Table class - Standard, Standard-IA
- DynamoDB Global Tables
    - Can configure DynamoDB as a multi-region database
- DynamoDB Accelerator (DAX)
    - Can reduce response time from miliseconds to microseconds through caching (DAX)
- DynamoDB Streams
- DynamoDB TTL
    - Automatically expire items(rows) based on timestamp
- DynamoDB Transactions - Can add ACID properties to a DynamoDB table
- DynamoDB Backup
- DynamoDB export / import
- DynamoDB Scaling - Read/Write Capacity (RCU/WCU)
    - Modify read/write throughput capacity
- DynamoDB Security
---
### DynamoDB Component - Table, Item
![[Pasted image 20260712101037.png|309]]
DynamoDB has Table as basic entity to store data
- Table contains items (like rows)
- Items contains Partition Key, Sort Key (optional), and attributes (key-value pairs).
- **Partition Key = Simple Primary Key**
- **Partition Key + Sort Key = Composite Primary Key**
    - Sort key is used to sort items with same partition values
---
### DynamoDB Component - Partition Key
- Acts as the primary index that uniquely identifies each item in your DynamoDB table
- Provides the ability to search for a particular item in your table
- Used an an **input to the internal hash function** in DynamoDB. The output from that function determines the physical internal storage in which the item will be stored
- The primary key attribute **must be a scalar**
    - string/number/binary 
---
### DynamoDB Component - Query and Scan
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
### DynamoDB Component - Global & Local Secondary Index
**Local Secondary Index (LSI)**
- Queries data over a **single partition** only (localized)
- Supports both **eventual consistency** or **strong consistency**
- **Can only be added at the same time that you create the base table (unlike GSI)**
- Must use same partition key as table's partition key, sort key can differ
- Use cases:
    - Query all items with single partition key value, like all moves where director = "Christopher Nolan", where director is partition key of the base table, and thus also LSI.
- When queried (reads from or writes to LSI), it **consumes provisioned RCU/WCU throughput of base table**. Unlike GSI, where RCU/WCU of GSI can be seperately or independently provisoned
    - Thus, using **eventual consistency** consumes lesser WCU than choosing **strong consistency**
- Can include un-projected attributes (columns) in the read queries. They are automatically fetched from base table. This concept is called "fetching". Although will causes additional query latencies in query responses and incurs higher provisioned throughput cost (RCU)
    - Thus, avoid "fetching" by carefully planning projected attributes of LSI during base table creation itself

**Global Local Secondary Index (GSI)**
- Queries data across **all partitions** of the entire table
- Only supports eventual consistency only, **no strong consistency**
- Can be added or deleted at **any time**
- Can have different partition key and sort key than table's partition key
- Can have a default limit of upto 20 GSIs
- **Only projected attributes** of GSIs are accessible while querying from GSI table, unlike LSI where un-projected attributes are "fetched" automatically from base table
- Do not consume RCU/WCU of base table when GSIs are queried. So RCU/WCU of GSI is separate from that of base table.
- Applications do not write directly to GSI. Any updates or deletes to DynamoDB table happen on base table first, and then asynchronously reflected to GSI table, using an **eventual consistency model** only
- Thus single update to base table corresponds to 2 or more write actions depending on how many GSIs we have
- **Recommended** to provison **WCU of GSI which is either greater or equal to WCU of base table**, since we can provison WCU/RCU of GSI seperately than that of base table, to avoid potential throttling
![[Pasted image 20260712102914.png|295]]
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
### DynamoDB Scaling - Read/Write Capacity (RCU/WCU)
- We can define the DynamoDB Read and Write capacity as per the application need.
- There are two modes to set the capacity:

**Provisioned Capacity Mode (Default)**
- You pre-define RCU/WCU for the table.
- Best when you have **predictable** or steady traffic.
- Cheaper than on-demand if your workload is consistent.
- **Risk of throttling if traffic suddenly spikes above provisioned limits.**
- Suitable if your application has predictable traffic that doesn’t vary over time
- Allows you to manually set or provision the RCU and WCU of your DynamoDB table
- Has an **Auto Scaling feature** that you can configure
    - Can set the target utilization, 
    - minimum provisioned capacity, 
    - and maximum provisioned capacity values in the Auto Scaling settings
    - By default auto-scaling is not enabled when created using CLI
- At risk of **over-provisioning** and having **unnecessary costs** when the incoming traffic is way lower than expected
- DynamoDB auto scaling uses the AWS Application Auto Scaling service to dynamically adjust provisioned throughput capacity on your behalf, in response to actual traffic patterns. 
    - This enables a table or a global secondary index to increase its provisioned read and write capacity to handle sudden increases in traffic, without throttling. When the workload decreases, Application Auto Scaling decreases the throughput so that you don’t pay for unused provisioned capacity.


**On-Demand Capacity Mode**
- No need to set RCU/WCU - DynamoDB **scales automatically**.
- Ideal for **unpredictable** or **spiky** workloads.
- You pay per request, which may be **costlier for steady high throughput**.
- Virtually no throttling unless **you hit account-level limits**. 
- For applications with inconsistent traffic or has varying access patterns
- Suitable if you expect that there’ll be more traffic with **sharp spikes** in the future
- **No manual Auto Scaling setting that you can configure**. 
    - The RCU & WCU are automatically scaled without any intervention
- Can be used if your application has a combination of predictable and variable traffic
    - Suitable if you have clearly defined access patterns throughout the year but with variable amounts of traffic on certain days only (flash sales or product announcements)
---
### DynamoDB Streams
- A data stream that **captures each and every data change made to the items**
    - Captures the item level modifications in time-ordered sequence
  - Stores the changes in the logs for 24 hours, exactly once and strictly ordered
- If an item was added, modified, or deleted, then that item will be included in the DynamoDB stream
- **Can be associated with AWS Lambda**. The function can poll the stream and execute a set of actions whenever it detects new stream records
- Can also be integrated with **Kinesis Data Streams**
- You can enable/disable a stream on a new or existing table
    - Important component that **needs to be enabled when using Amazon DynamoDB Global Tables**
- DynamoDB Streams operates asynchronously, so there is no performance impact on a table if you enable a stream.
![[20260712105301.png|134]]
**Use cases:**
- Maintain time-series data
- Real-time monitoring e.g. connected vehicles, sensor data, Notifying everyone on an activity e.g. friend creates a post on social media
- Backup/Change data capture (CDC) of DynamoDB table
---
### Multi Region Database - DynamoDB Global Tables
- provide multi-region replication across all replicas/tables
- Can read/write to any replica
- includes ongoing data changes across all the tables during replication

- Multi-region, multi-active (multi-master), serverless tables across regions
- completely automated

- underlying infra is entirely managed by AWS themselves
- **doesn't reside within a custom VPC**
- 99.999% availability

**Use cases:**
- Global application requiring low latency access for users
- Can handle region level failure (DR - Disaster Recovery)
![[Pasted image 20260712121323.png]]
---
### DynamoDB Accelerator - DAX
- Fully managed highly available **in-memory cache for DynamoDB**
- 10x performance improvement with **single digit millisecond to microsecond level latency**
- API-compatible with DynamoDB. Only client & endpoint needs to change to use with an existing application.
- DAX provides access to **eventually consistent** data from DynamoDB tables
- An in-memory cache for Amazon DynamoDB that is fully managed and highly available
- Launches a DAX cluster that can be run in your default or custom Amazon **VPC**
- Provides response time in **microseconds** and not just in milliseconds
- Delivers fast response times for accessing eventually consistent data
- Significantly reduces the response times of your DynamoDB database

**Use cases:**
- Real-time bidding, social gaming, and trading applications
- Read intensive or high throughput reads
**Anti-pattern:**
  - Strongly consistent read requirements
  - Write intensive, not having repeated reads
  

---
### DynamoDB TTL (Time To Live)
- TTL lets you automatically delete expired items (per item) based on a **timestamp attribute**.
- You specify a numeric attribute (e.g. expiresAt) containing a Unix epoch time (in seconds).
- DynamoDB checks items in the background and removes expired ones **asynchronously**
- No extra cost for TTL deletion; only standard read/write charges apply when you insert/update items.
- **Helps reduce storage cost and keeps tables clean.**

**Use cases:**
- User sessions/tokens - Automatically remove expired login sessions or JWT tokens.
- Rate-limiting counters - Keep count of API hits per minute and auto-reset them via TTL.
- OTP / verification codes - Remove expired OTP codes, email verification tokens, etc
---
### DynamoDB Transactions
- Provides ACID properties to your DynamoDB table for your transactional workloads
- Provides an **all-or-nothing change** to multiple items both within and **across DynamoDB tables**
- Consists of DynamoDB transactional read and write APIs
	- `TransactWriteItems`
	- `TransactGetItems`
- Empowers you to manage complex business workflows that require adding, updating, or deleting multiple items as an atomic operation
---
### DynamoDB Backup (DR)

**Point In Time Restore (PITR)**
- PITR continuously backs up the DynamoDB table
- PITR allows restore the DynamoDB table to **any second** in the last 35 days. - (Recovery-point-objective can be 1 second)
- The recovery creates a new table
- Automated backup process
- Enables continuous backups to your table
- Allows you to restore your table at a point in time that you specify
- **Entails additional costs**
- For low RPO

**On-demand Backup**
- For acceptable low RPO (recovery point object)
- On-demand backup creates a full manual snapshot that you can keep indefinitely.
- Can use **AWS Backup service** to create the backups including cross-region copy
- Restoring a backup always creates a new table.
- **Manual** backup process
- **No continuous** backups
- Can only restore to a particular backup that you’ve taken
- A **cost-effective** yet limited backup option feature for your data
- **When restoring from a on-demand backup**
    - When you create an on-demand backup, a time marker of the request is cataloged. 
    - The backup is created asynchronously by applying all changes until the time of the request to the last full table snapshot. 
    - Backup requests are processed instantaneously and become available for restore within minutes.
    - However, some settings are not carried over on the restored table and you must manually configure them after restoring.
    - You must manually set up the following on the restored table:
		- Auto scaling policies
		- AWS Identity and Access Management (IAM) policies
		- Amazon CloudWatch metrics and alarms
		- Tags
		- Stream settings
		- Time to Live (TTL) settings


---
### DynamoDB Export to / Import from S3
**Export to S3**
- Uses **PITR snapshots**, so you can export data from a specific point in time.
- Export **does not consume** DynamoDB Read/Write capacity
- Useful for 
    - **analytics, data pipelines**
        - DynamoDB -> S3 -> Athena
    - and **long-term data storage**.
**Import from S3**
- Allows you to import data directly from Amazon S3 into a new DynamoDB table.
- Supports CSV, DynamoDB JSON, and ION formats.
- **No write capacity consumed** during import (doesn't use WCUs).
- Import creates a new table — it cannot overwrite an existing one.
- Useful for 
    - **bulk loading, migrations**, 
    - and rebuilding tables from exported data.
---
### DynamoDB Security
- Protects your data both in transit and at rest
- All data stored in Amazon DynamoDB is fully **encrypted at rest by default**
- The API calls from your private Amazon EC2 instances that go to DynamoDB can be **configured to not traverse the public Internet** by 
	- creating a **VPC Gateway Endpoint** 
	- and adding a new route table entry
---
