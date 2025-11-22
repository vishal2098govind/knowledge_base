#redis 

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
```