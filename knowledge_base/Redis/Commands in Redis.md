#redis

## Basic Strings
- `SET` and `GET` commands
```redis
SET message 'hi there!'
---
GET message // 'hi there!'
```
- `GET` and `SET` are simple and are intended to store only plain strings and plain numbers
- Redis is capable of storing many other types of data as well
- Each type of data has [diff set of commands](https://redis.io/docs/latest/commands/) supported

| Data Types  | Purpose                                | Some Commands                    |
| ----------- | -------------------------------------- | -------------------------------- |
| Strings     | Store plain string or number           | GET, SET, APPEND                 |
| List        | List of strings                        | LINDEX, LLEN                     |
| Hashes      | Collection of key-value pairs          | HGET, HSET, HDEL                 |
| Sets        | Set of strings (each string is unique) | SADD, SCARD, SDIFF               |
| Sorted Sets | Set of strings in a particular order   | ZADD, ZCOUNT, ZDIFF              |
| Bitmap      | Kind of a collection of booleans       | BITOP, BITCOUNT, BITPOS          |
| Hyperlog    | Kind of a collection of booleans       | PFADD, PFCOUNT, PFMERGE          |
| JSON        | Nested JSON structure                  | JSON.SET, JSON.GET, JSON.DEL     |
| Index       | Internal data structure for searching  | FT.SEARCH, FT.CREATE, FT.PROFILE |
### Reading/Understanding Command Docs
Capital words are keywords
'I' pipe symbols mean 'or'
Square brackets means it is optional
```
SET key value
	[
		  EX seconds
		| PX miliseconds
		| EXAT unix-time-seconds
		| PXAT unix-time-miliseconds
		| KEEPTTL
	]
	[NX|XX]
	[GET]
```

## SET
 SET commands **work with only strings and numbers**
```
SET 
	key // Key we are trying to set
	value // Value we want to store
	[EX s| PX ms| EXAT unix time s| PXAT unix time ms| KEEPTTL] // Options for when this value should **expire**
	[NX | XX] // either set only if key not exists before or only if key already exists 
	[GET] // Return the previous value stored at this key
```

- Use cases of **expiration** options in `SET`
```
SET
	color
	red
	EX 2 // automatically delete this value after 2 seconds!
```
- Why Expiration / Use cases of expiration?
	- Redis was originally designed as a caching server
	- Thus, for caching purposes, and not server old data from redis, we use expiration

### SETNX, SETEX, MSET
- similar commands
- SETEX is exact same thing as SET with EX
- SETNX is exact same thing as SET with NX
```
SETNX color 2 // only writes if key color does not exist
```
- MSET is used to set multiple key values in single command
```
MSET 
	color // Set key 'color'
	red   // Store value 'red' for 'color'
	car   // Set key 'car'
	toyota // Store value 'toyota' for car 'car'
```
- MSETNX is similar to SETNX allowing multiple key-value pairs to set

## GET
- MGET - allows to get values of multiple keys. 
- **GET and MGET work with strings and numbers only**
```
MGET
	color
	model

// output:
["red", "toyota"]
```


## DEL
DEL command works with any data type, not just strings or numbers
```
DEL
	key // deletes value associated with the key
```

## GETRANGE and SETRANGE
- GETRANGE
	- Returns a sequence of characters from an existing string stored at that key
	- `GETRANGE color 0 3`
```
SET model toyota
GETRANGE model 0 2 // toy
```
- SETRANGE
	- Update portion of existing string
	- `SETRANGE color 2 blue`
```
SET model toyota
GET model // toyota
SETRANGE model 2 blue
GET model // toblue
```

### Use case of GETRANGE and SETRAGE
- Why would we ever want to replace part of a string?
	- one of the big adv of redis is it's performance. 
	- Thus, we can use these basic operations to make much more complicated things
		- a lot of commands, at their face value, might appear to not be very useful, 
		- but we can kind of twist these commands a little but
		- we can use these commands in **very creative ways**
		- and doing so, we can make **systems that are extremely fast in nature**

### Numbers
- Redis stores numbers as strings
```
SET age 20 // Redis stores 20 as a string "20"
GET age // "20"

// these commands can be implemented using SET and GET as well, 
// but these still exist to achieve for consistent updates and handling concurrency
// just the sheer existence of these commands tell us synchronous nature of redis
INCR age // "21"
DECR age // "20"
INCRBY age 10 // "30"
DECRBY age 10 // "20"
INCRBYFLOAT age 6.400145 // "26.400145"
INCRBYFLOAT age -6.400145 // "20"
```
- **REDIS is SINGLE THREADED and SYNCHRONOUS - one command at a time**
- i.e. even if a tremendous number of commands are coming in exact same time, it doesn't make a diff to redis
- Redis only processes one command at a time in order that they are received.