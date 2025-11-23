#redis #sorted-sets

- A mix of hash and set
- There are no keys and values
- There are members (key) and scores (value) .
- Members are unique
- Scores are always numbers
- All members are ordered by score (ascending)
```shell
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZADD products 20 monitor
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZSCORE products monitor
20
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZREM products monitor
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZSCAN products monitor
(error) ERR invalid cursor
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZSCORE products monitor
(nil)
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZADD products 45 cpu
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZADD products 10 keyboard
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZADD products 55 power
(integer) 1

# ZCARD key
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZCARD products
(integer) 3

# ZCOUNT key min max: to get how many products have score within >= min and <= max (range query)
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZCOUNT products 3 20
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZCOUNT products 40 50
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> 

# ZCOUNT inclusive and exclusive range query
# min <= score <= max
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZCOUNT products 40 55
(integer) 2 # cpu and power
# min <= score < max
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZCOUNT products 40 (55
(integer) 1 # only cpu, power is not included as it is = 55

# -inf < score < +inf
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZCOUNT products -inf +inf
(integer) 3
# score >= 40
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZCOUNT products 40 inf
(integer) 2

# Removing max and min member(s) using ZPOPMAX and ZPOPMIN
# ZPOPMIN key [count] : removes least [count] number of members
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZPOPMIN products 2 
[
    "keyboard",
    "10",
    "cpu",
    "45"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZADD products 45 cpu
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZADD products 10 keyboard
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZADD products 55 power
(integer) 0
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZPOPMAX products 2
[
    "power",
    "55",
    "cpu",
    "45"
]

# increment score of a member
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZINCRBY products 2 keyboard
12
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZSCORE products keyboard
12
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZINCRBY products -2 keyboard
10
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZSCORE products keyboard
10

# Query a sorted set
# ZRANGE key starting-index stopping-index [WITHSCORES]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZRANGE products 0 3
[
    "keyboard",
    "cpu",
    "power"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZRANGE products 0 3 WITHSCORES
[
    "keyboard",
    "10",
    "cpu",
    "45",
    "power",
    "55"
]


# ZRANGE key min-score max-score BYSCORE [WITHSCORES]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZRANGE products 30 50 BYSCORE WITHSCORES
[
    "cpu",
    "45"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZRANGE products 30 60 BYSCORE WITHSCORES
[
    "cpu",
    "45",
    "power",
    "55"
]
# score >= 20
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZRANGE products 20 inf BYSCORE WITHSCORES
[
    "cpu",
    "45",
    "power",
    "55"
]
# score >= 0
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZRANGE products 0 inf BYSCORE WITHSCORES
[
    "keyboard",
    "10",
    "cpu",
    "45",
    "power",
    "55"
]
# score >= 10
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZRANGE products 10 inf BYSCORE WITHSCORES
[
    "keyboard",
    "10",
    "cpu",
    "45",
    "power",
    "55"
]
# score > 10
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZRANGE products (10 inf BYSCORE WITHSCORES
[
    "cpu",
    "45",
    "power",
    "55"
]

# ZRANGE with REV : reverses the order of members before doing any comparison
# ZRANGE key min-index max-index REV
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZRANGE products 1 2 REV
[
    "cpu",
    "keyboard"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZRANGE products 1 2 REV WITHSCORES
[
    "cpu",
    "45",
    "keyboard",
    "10"
]

# ZRANGE with LIMIT : used for pagination
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZRANGE products -inf inf BYSCORE WITHSCORES LIMIT 0  2
[
    "keyboard",
    "10",
    "cpu",
    "45"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZRANGE products -inf inf BYSCORE WITHSCORES LIMIT 2  2
[
    "power",
    "55"
]
```

## Use cases of Sorted Sets
- Tabulate 'the most' or 'the least' of some attribute of a collection of hashes
![[Pasted image 20251122163656.png]]
- Create relationships between records and sorted by some criteria
```
Authors: Samantha, Alexa, Jimantha
Books: "A Biography", "History Book"

Let's say, "Samantha" and "Alex" both have worked on "A Biography" and "History Book" together

# Inside Redis:
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HSET authors:4 name samantha
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HGET authors:4 name
samantha
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HSET authors:14 name alex
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HSET books:5 name "A Biography"
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HSET books:43 name "History Book"
(integer) 1
----
This stores two author objects:
authors:4 -> {name: "samantha"}
authors:14 -> {name: "alex"}
also stores two books:
books:5 -> {name: "A Biography"}
books:4 -> {name: "History Book"}

Also we need to add some data to represent the relation between these objects/hashes with each other and also we need to order/sort this relationship based on the no.of copies sold for the books have recieved

authors:books:4 (for samantha) -> {member=5: score=560} (score represents no.of copies sold for the book with id 5 ("A Biography"))
authors:books:4 (for samantha) -> {member=43: score=4600} (score represents no.of copies sold for the book with id 5 ("History Book"))
authors:books:14 (for alex) -> {member=5: score=560} (score represents no.of copies sold for the book with id 5 ("A Biography"))
authors:books:14 (for alex) -> {member=43: score=4600} (score represents no.of copies sold for the book with id 5 ("History Book"))

Inside Redis:
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZADD authors:books:4 560 5
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZADD authors:books:4 4600 43
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZADD authors:books:14 560 5
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZADD authors:books:14 4600 43
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZRANGE authors:books:4 -inf +inf BYSCORE WITHSCORES
[
    "5",
    "560",
    "43",
    "4600"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZRANGE authors:books:14 -inf +inf BYSCORE WITHSCORES
[
    "5",
    "560",
    "43",
    "4600"
]

For each author, we can find the most popular books by author "samantha"

Also we can choose to represent the **author-date** property of each book in a sorted-set to be able to fetch most recently authored books by an author. thus here the score could be UNIX timestamp

```