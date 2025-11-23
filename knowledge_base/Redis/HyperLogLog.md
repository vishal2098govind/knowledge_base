#redis #hyperloglog

- This is an algorithm for **approximately** counting the number of unique elements
- Similar to a set, but doesn't store the elements
- Will seem useless at first glance
- Redis offers two commands `PFADD` and `PFCOUNT` for this

```
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> PFADD vegetables celery
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> PFADD vegetables celery
(integer) 0
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> PFADD vegetables potato
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> PFADD vegetables potato
(integer) 0
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> PFCOUNT vegetables
(integer) 2
```

## When to use HyperLogLogs
- get **approximate** number of **unique** views a particular resource has
```
> PFADD views:items:5 username1
```
- If this returned 1, then this user has viewed the item with id=5 for the first time. thus, can increment the resource's views count by 1
- If this returned 0, then this user has already viewed the item with id=5 earlier, and is viewing again this time. thus, need not increment resource's views count
```
> PFCOUNT views
```
- This gives **approximate** number of unique views of this resource
-  **HyperLogLog always has a constant size** no matter how many items have been added
	- about 12kb no matter what, always

### Approximate count
- In general, whenever we do a `PFCOUNT views:items:5` , the number seen is incorrect by the **error rate of 0.81%** 
	- If we do `PFADD views:items:5 <id>` with **1000** unique ids, 
		- then `PFCOUNT views:items:5` would return a number anywhere around **991** to **1008**
	- This is the **trade-off** me make for not actually storing unique usernames who have viewed the item in a data-structure like SET