#redis 

- Used on sets, sorted sets and lists
- Calling this command 'SORT' is misleading!!!!
- We do not always use this SORT command for sorting data, there are actually many other use cases where no sorting is done

```
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HSET books:good title 'Good book'
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HSET books:good year 1950
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HGETALL books:good
[
    "title",
    "Good book",
    "year",
    "1950"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HSET books:bad title 'Bad book' year 1930 
(integer) 2
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HSET books:ok title 'ok book' year 1940
(integer) 2
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HGETALL books:bad 
[
    "title",
    "Bad book",
    "year",
    "1930"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HGETALL books:ok
[
    "title",
    "ok book",
    "year",
    "1940"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZADD books:likes 999 good
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZADD books:likes 0 bad
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> ZADD books:likes 40 ok
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> zrange books:likes -inf inf BYSCORE WITHSCORES
[
    "bad",
    "0",
    "ok",
    "40",
    "good",
    "999"
]
```

- SORT command operates on **members** in case of **sorted sets**, **not** the scores
```sh
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SORT books:likes ALPHA
[
    "bad",
    "good",
    "ok"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789>  SORT books:likes ALPHA LIMIT 0 2
[
    "bad",
    "good"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789>  SORT books:likes ALPHA LIMIT 1 2
[
    "good",
    "ok"
]

# sort by year property of book
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SORT books:likes BY books:*->year
[
    "bad",
    "ok",
    "good"
]
```
### Joining data with SORT
```sh
# titles of books sorted by book
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SORT books:likes BY books:*->year GET books:*->title
[
    "Bad book",
    "ok book",
    "Good book"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SORT books:likes BY books:*->year GET books:*->title GET books:*->year
[
    "Bad book",
    "1930",
    "ok book",
    "1940",
    "Good book",
    "1950"
]
# using # to include the members of the sorted set
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SORT books:likes BY books:*->year GET # GET books:*->title GET books:*->year
[
    "bad",
    "Bad book",
    "1930",
    "ok",
    "ok book",
    "1940",
    "good",
    "Good book",
    "1950"
]
```

### Using SORT just for joining and not sorting anything explicitly

```
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SORT books:likes BY nosort GET # GET books:*->title GET books:*->year
[
    "bad",
    "Bad book",
    "1930",
    "ok",
    "ok book",
    "1940",
    "good",
    "Good book",
    "1950"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> SORT books:likes BY nosort GET # GET books:*->title GET books:*->year DESC
[
    "good",
    "Good book",
    "1950",
    "ok",
    "ok book",
    "1940",
    "bad",
    "Bad book",
    "1930"
]
```