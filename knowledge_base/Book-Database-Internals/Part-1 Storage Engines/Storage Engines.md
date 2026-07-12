Databases are modular systems and consist of multiple parts:
- Transport Layer accepting requests
- Query Processor determining most efficient way to run queries
- Execution engine carrying out operations
- Storage Engine for 
	- storing,
	- retrieving,
	- managing data in memory and on disk
	- Designed to capture a persistent, long-term memory of each node

DBMS are applications built on top of storage engines
Storage engines offer simple data manipulation APIs to create, update, delete and retrieve records

For storage engine, any key and value is an arbitrary sequence of bytes with no prescribed form
From storage engine perspective, an int32 type key and ascii type key in tables are just serialized entries

Many storage engines were developed independently from DBMSs they're embedded into
Using pluggable storage engines has enabled Database Developers to bootstrap DB systems using existing storage engines, and concentrate on other sub-systems

Clear separation between DB system components opens up an opportunity to switch between different engines, potentially better suited for particular use cases.

For example, MySQL has several storage engines
- InnoDB
- MyISAM
- RocksDB

MongoDB allows switching between storage engines
- WiredTiger
- In-Memory
- MMAPv1

