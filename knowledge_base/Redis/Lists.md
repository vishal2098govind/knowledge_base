#redis #time-series-data

- Implemented as a doubly-linked-list
- Not arrays
- Unless we only deal with head or tail of the list, using lists would be expensive
- thus, mostly used for data like **time-series** data, where we only add at one end of the list

```
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> LPUSH temps 12 23 45
(integer) 3
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> LLEN temps
(integer) 3
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> LINDEX temps 1
23
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> LINDEX temps -1
12
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> LINDEX temps 0
45
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> RPUSH temps 34
(integer) 4
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> LINDEX temps -1
34
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> LPUSH temps 0
(integer) 5
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> LINDEX temps 0
0
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> LPUSH temps 0 1 2 3 4 5
(integer) 11
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> LINDEX temps -1
34
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> LINDEX temps 0
5
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> LRANGE temps 0 -1
[
    "5",
    "4",
    "3",
    "2",
    "1",
    "0",
    "0",
    "45",
    "23",
    "12",
    "34"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> LPOP temps 
5
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> LPOP temps 2
[
    "4",
    "3"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> LRANGE temps 0 -1
[
    "2",
    "1",
    "0",
    "0",
    "45",
    "23",
    "12",
    "34"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> RPOP temps
34
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> RPOP temps 2
[
    "12",
    "23"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> LRANGE temps 0 -1
[
    "2",
    "1",
    "0",
    "0",
    "45"
]
```

## Use cases of lists
- Append-only or prepend-only data: 
	- temperature readings
	- stock values
	- time
- When we need only last/first N of something
- Our data has no other sort order besides the order of insertion
```
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HSET books:a1 title 'Good book'
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HSET books:a2 title 'Bad Book'
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> RPUSH reviews a1
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> RPUSH reviews a2
(integer) 2
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SORT reviews BY nosort GET books:*->title
[
    "Good book",
    "Bad Book"
]
```