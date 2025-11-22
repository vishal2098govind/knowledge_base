#redis 
- depends from language client of redis
- in most redis client
```python
client = redis.Redis(...)

pipe = client.pipeline()
pipe.set('foo', 'bar')
pipe.get('bing')
pipe.execute()
```
- in node-redis client
```ts
cosnt results = await Promise.all([
	client.get('color'),
	client.get('name'),
])
```
