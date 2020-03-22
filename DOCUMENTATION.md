## Configurations 
Configurations are set during hte proxy startup. The proxy will grab the following environment if none are set it will use the default. 

## Environment Variables
* `REDISHOST`: The redis backing host. Default `localhost`
* `REDISPORT`: The redis backing host port. Default `6379`
* `REDISPASSWORD`: The password of the redis backing host. Default `""`
* `REDISDB`: The db of the redis backing host. Default `0`
* `CACHETTL`:  The duration the cache should live before getting curated. Default `30` seconds
* `CACHECAPACITY`: The size of the map[key]val for the cache. Default `10`
* `HOSTPORT`: The port the redix proxy will host on: Default `2020`
* `CURATECYCLE`: How long to wait before the next curate happens. Default `30` seconds

## Methods 
As of now `redis-proxy` will only accept `GET` request. Any other method will result in `StatusMethodNotAllowed`

### Method `GET`
* Headers: `rediskey`. The key to retrieve the value from redis-proxy (required)

## Monitoring 
 redis-proxy implements `/metrics` URI for prometheus to grab metrics. 
* `cachehit`: The total number of cache hit
* `cachemiss`: The total number of cache miss 
* `curatequeue`: The total number of curate being queued behind. Note this is an important metrics, this indicates that the curating the cache is taking longer than the curate cycle. Either the cache size will need to be lower or increase the time between the curate cycle. 
* `invalidekey`: The total number of times when client request an invalid key
* http_status_200: The total number of `StatusOK` response.
* http_status_400: The total number of `StatusBadRequest` response. Usually from client didn't specify `rediskey` as header.
* http_status_404: The total number of `StatusNotFound` response. Usually happen when client requested an invalid key.
* http_status_405: The toaly number of `StatusMethodNotAllowed`. Usually when client send request with invalid http method. 

