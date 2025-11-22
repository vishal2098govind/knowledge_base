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
```python
# command doc:
SSCAN key cursor [MATCH pattern] [COUNT count]

SSCAN 
	colors:1     # key
	0            # cursor ID
	COUNT 100    # no.of elements to return
```


## Set Use Cases
- Enforcing uniqueness of any value
	- E.g. set of usernames
	- `SISMEMBER usernames vishal_a_geek` to check if `vishal_a_geek` is already taken as username
- Creating relationship between records
	- E.g. implement a like-system inside our app, where user can like certain items
	- maintain a separate set for every single user like `users:45:likes` is a set of liked items by user with id of 45
	- we can use this set to:
		- find items liked by the user: `SMEMBERS users:45:likes`
		- find how many items are liked by this user: `SCARD users:45:likes`
		- find if this user liked item with id 145: `SISMEMBER users:45:likes 145`
- Finding common attributes between different entities
	- Which items both user-45 and user-32
		- `SINTER users:45:likes users:32:likes`
```
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SADD users:45:likes 12 23 245 134 154 145
(integer) 6
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SADD users:46:likes 1 3 2 13 154 145
(integer) 6
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SINTER users:45:likes users:46:likes
[
    "145",
    "154"
]
```
- General list of elements where the order of elements doesn't matter
	- `SADD domains:banned ezmail.com freemail.com scammail.com`