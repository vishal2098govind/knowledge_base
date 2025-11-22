#redis 

### Serialisation
- We don't really need to store ID as a key value pair inside the hash map when the id is already part of our `redis` key
- By including ID inside of the hash map, we would use up a little bit of extra memory (since redis stores everything in memory)
	- memory costs more in general than hard drive space, so we usually want to pare down the records with redis to be as small as possibly can
	- removing the ID field is generally desirable

- ![[Pasted image 20251111234909.png]]


## Deserialisation
- When now we try to get company by doing `HGET company company#1234`, we would want to include the ID of the company while returning the company object to the rest of the app.


```ts
export const getUserById = async (id: string) => {
    const user = await client.hGetAll(usersKey(id))
    return deserialize(id, user)
};

export const createUser = async (attrs: CreateUserAttrs) => {
    const id = genId()
    await client.hSet(usersKey(id), serialize(attrs))
    return id;
};


const serialize = (user: CreateUserAttrs) => {
    return {
        username: user.username,
        password: user.password,
    }
}

const deserialize = (id: string, user: {[x: string]: string}) => {
    return {
        id,
        username: user.username,
        password: user.password,
    }
}
```

## Serialising Date Times
- Whenever we store a date of any kind in Redis, we have a little decision to make
	- Redis can't (by default) query/search for or sort certain kind of dates, like
		- 1994-11-05T08:15:30-05:00
		- 'Thu Mar 17 1994 10:56:16 GMT-0500 (Central Daylight Time)'
	- Thus, we would prefer to store dates in
		- Unix time seconds
		- or Unix time as milliseconds
		- Unix time represents number of seconds elapsed since Jan 1 1970
```ts
import type { CreateItemAttrs } from '$services/types';

export interface CreateItemAttrs {
	name: string;
	imageUrl: string;
	description: string;
	createdAt: DateTime;
	endingAt: DateTime;
	ownerId: string;
	highestBidUserId: string;
	status: string;
	price: number;
	views: number;
	likes: number;
	bids: number;
}

export const serialize = (attrs: CreateItemAttrs) => {
        return {
            ...attrs,
            createdAt: attrs.createdAt.toMillis(),
            endingAt: attrs.endingAt.toMillis(),
    }
};

```