#redis 

- Allows us to model more interesting relationships between different pieces of data
- Set is a collection of strings, where each string is unique

### SADD, SMEMBERS
```
SADD colors red
// 1
SADD colors red // 1
SADD colors green // 1
SADD colors blue orange

SMEMBERS colors
// ["red", "green"]
```
### SUNION, SINTER, SDIFF
Union, Intersection and Set-difference

```
SADD colors:1 red blue orange
SADD colors:2 blue green purple
SADD colors:3 blue red purple

SUNION colors:1 colors:2 colors:3 // ["red", "blue", "green", "purple"]
SINTER colors:1 colors:2 colors:3 // ["blue"]
SDIFF colors:1 colors:2 colors:3 // ["orange"]
```

### Store Variations
- SUNION, SUNIONSTORE
- SDIFF, SDIFFSTORE
- SINTER, SINTERSTORE
They are similar just a very small variation that they save the output in Redis as a seperate key
```
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SADD colors:1 red blue orange
(integer) 0
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SADD colors:2 blue green purple
(integer) 0
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SMEMBERS colors:1
[
    "red",
    "blue",
    "orange"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SMEMBERS colors:2
[
    "blue",
    "green",
    "purple"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SADD colors:3 blue red purple
(integer) 3
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SMEMBERS colors:3
[
    "blue",
    "red",
    "purple"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SINTER colors:1 colors:2 colors:3
[
    "blue"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SDIFF colors:1 colors:2 colors:3
[
    "orange"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SINTERSTORE colors:123 colors:1 colors:2 colors:3
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SMEMBERS colors:123
[
    "blue"
]
```

### Membership: SISMEMBER, SMISMEMBER
```
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SMEMBERS colors:1
[
    "red",
    "blue",
    "orange"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SISMEMBER colors:1 red
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SMISMEMBER colors:1 red blue
[
    1,
    1
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SMISMEMBER colors:1 red blue purple
[
    1,
    1,
    0
]
```
### Cardinality: SCARD
```
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SCARD colors:1
(integer) 3
```
### Remove: SREM
```
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SMEMBERS colors:2
[
    "blue",
    "green",
    "purple"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SREM colors:2 blue
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SISMEMBER colors:2 blue
(integer) 0
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SMEMBERS colors:2
[
    "green",
    "purple"
]
```

### SSCAN vs SMEMBERS
- Sometimes we might have a set with many many items inside of them
- `SMEMBERS` gives all the strings with one single command, no matter how many are present in the set
- `SSCAN` allows to mention # of elements to return
```
SSCAN colors:1 0 COUNT 100 
```