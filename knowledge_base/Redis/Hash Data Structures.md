#redis 

**key** : any string
**value** : map or python dict where the values can only be strings or numbers, no nested map or arrays supported

### HSET, HGET, HGETALL
```
// "company" -> {"name": "Company Co", "age": 1915}
HSET company name 'Company Co' age 1915

HGET company name
Company Co
```

Example:
```
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HGET company
(error) ERR wrong number of arguments for 'hget' command
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HGET company name
Company Co
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HGETALL company
[
    "name",
    "Company Co",
    "age",
    "1915",
    "industry",
    "materials",
    "revenue",
    "5.3"
]
```

### HEXISTS, HDEL
```
// check if the key exists (not truthyness of the key)
>HEXISTS company age
1

redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> DEL company
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HEXISTS company age
(integer) 0
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HGET company age
(nil)
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HSET company name 'Company Co' age 1915
(integer) 2
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HGET company age
1915
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> 
1915
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HDEL company age
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> 
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HGET company age
(nil)
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HGET company
(error) ERR wrong number of arguments for 'hget' command
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HGETALL company
[
    "name",
    "Company Co"
]
```

### HINCRBY, HINCRBYFLOAT, HSTRLEN
```
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HSET company age 10
(integer) 1
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HGETALL company
[
    "name",
    "Company Co",
    "age",
    "10"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HINCRBY company age 10
(integer) 20
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HGETALL company
[
    "name",
    "Company Co",
    "age",
    "20"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HINCRBY company revenue 10
(integer) 10
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HGETALL company
[
    "name",
    "Company Co",
    "age",
    "20",
    "revenue",
    "10"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HINCRBY company age -2
(integer) 18
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HGETALL company
[
    "name",
    "Company Co",
    "age",
    "18",
    "revenue",
    "10"
]
```

### HSTRLEN
```
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HINCRBYFLOAT company age 1.001
19.001
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HSTRLEN company name
(integer) 10
```

### HKEYS, HVALS
```
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HKEYS company
[
    "name",
    "age",
    "revenue"
]
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HVALUES company
(error) ERR unknown command 'HVALUES'
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HVALS company
[
    "Company Co",
    "19.001",
    "10"
]
```

### Hashes with Code
#### Issue with `HSET`
```
import 'dotenv/config';
import { client } from '../src/services/redis';

const run = async () => {
    await client.hSet('car', {
        color: 'red',
        year: 1950,
        engine: { cylinder: 1 },
    })
    // HSET car color red year 1950

    const car = await client.hGetAll('car')
    // HGETALL car

    console.log(car);
};
run();
```

```
redis-19789.c326.us-east-1-3.ec2.redns.redis-cloud.com:19789> HGETALL car
[
    "color",
    "red",
    "year",
    "1950",
    "engine",
    "[object Object]"
]
```

- Issue with HGETALL
```
import 'dotenv/config';
import { client } from '../src/services/redis';

const run = async () => {
    await client.hSet('car', {
        color: 'red',
        year: 1950,
        engine: { cylinder: 1 },
        owner: null,
        service: undefined
    })
    // HSET car color red year 1950

    const car = await client.hGetAll('car')

    console.log(car);
};
run();

// OUTPUT:
[ERROR] 00:39:30 TypeError: Cannot read properties of null (reading 'toString')
```
- Why error while HSET?:
	- The node-redis client converts all values to string using .toString() method before putting in to the HSET command
```
HSET car 
	color 'red'.toString()
	year 1950.toString()
	engine: {cylinder: 1}.toString() // '[object Object]'
	owner: null.toString() // results into error
	service: undefined.toString() // results into error
```
- Thus, we are not able to convey that the car has no owner yet by setting owner to null. So to work around this, we can do:
```
import 'dotenv/config';
import { client } from '../src/services/redis';

const run = async () => {
    await client.hSet('car', {
        color: 'red',
        year: 1950,
        engine: { cylinder: 1 },
        owner: null || '',
        service: undefined || ''
    })
    // HSET car color red year 1950

    const car = await client.hGetAll('car')

    console.log(car);
};
run();
```
- output:
```
[Object: null prototype] {
  color: 'red',
  year: '1950',
  engine: '[object Object]',
  owner: '',
  service: ''
}
```

#### Issues with `HGETALL`
in code, the `HGETALL` always returns an object (empty if null)
```
HGETALL car553
[]
------------------------------
const car553 = await client.hGetAll('car#553')
console.log(car553)
// [Object: null prototype] {}

// thus a simple existense check will not work
if (!car553) {
	// this block is never reached
	console.log('respond with 404')
	return
}
// thus solution:
if (Object.keys(car553).length == 0) {
	console.log('respond with 404')
	return
}
```


### When to use Hash and when not to
- When to use
	- Record has many attributes 
	- Collection of these records need to be sorted by many different ways
		- Tabular view allowing to sort by different columns
	- Often need to access a single record at a time
		- Item details
- When **NOT** to use
	- Record is only for counting or enforcing uniqueness: 
		- Like/View count
	- Record with only one or two attributes
	- **Relational data**: when record is only used for creating relations between different records
	- **Time series data**: Snapshots of single or few values changing over time