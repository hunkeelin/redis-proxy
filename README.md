# RedisProxy 
[![CircleCI](https://circleci.com/gh/hunkeelin/redis-proxy.svg?style=shield)](https://circleci.com/gh/hunkeelin/redis-proxy)
[![Go Report Card](https://goreportcard.com/badge/github.com/hunkeelin/redis-proxy)](https://goreportcard.com/report/github.com/hunkeelin/redis-proxy)
[![GoDoc](https://godoc.org/github.com/hunkeelin/redis-proxy/src/server?status.svg)](https://godoc.org/github.com/hunkeelin/redis-proxy/src/server)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/hunkeelin/redis-proxy/master/LICENSE)


## Overview
This repo implements the redis-proxy service in pods in the below diagram. It act as middleware that takes http `get` and map it to redis `get` protocol. 
![Architecture](arch.png)

## Golang version

`redis-proxy` is currently compatible with golang version from 1.12+.

## User Manual
* The server only allow `GET` and it requires uri `/key/${rediskey}`. For detail documentation on configurations please checkout [documentation](src/server/README.md).
## Usage
```go
package main

import (
    redisproxy "github.com/hunkeelin/redis-proxy/src/server"
)

func main() {
    panic(redisproxy.Server())
}
```

## Example
```bash
$ redis-cli
127.0.0.1:6379> set foo1 bar1
OK
127.0.0.1:6379> set foo2 bar2
OK
127.0.0.1:6379> set foo3 bar3
OK
127.0.0.1:6379> set foo4 bar4
OK
127.0.0.1:6379> set foo5 bar5
OK
127.0.0.1:6379> exit

$ make build
o build -o redis-proxy -v
$ export CACHECAPACITY=3
$ export REQUESTLIMIT=4
$ ./redis-proxy &
INFO[0000] Starting proxy with the following configuration
INFO[0000] Redis Backing Host: localhost
INFO[0000] Redis Backing port: 6379
INFO[0000] Cache TTL: 30 seconds
INFO[0000] Number of Cache:  10
INFO[0000] Host port: 2020
INFO[0000] Curate Cycle: 30 seconds
INFO[0000] Request limit: 4
INFO[0000] Request burst: 3
listening to :2020

// Should fail because require header 
$ curl -i localhost:2020
HTTP/1.1 400 Bad Request
Date: Sun, 22 Mar 2020 05:26:32 GMT
Content-Length: 0

$ curl -i localhost:2020 -H "rediskey: foo1"
HTTP/1.1 200 OK
Date: Sun, 22 Mar 2020 05:27:20 GMT
Content-Length: 4
Content-Type: text/plain; charset=utf-8

bar1

```

## Features
* HTTP webservice: Clients interface to the Redis proxy through HTTP, with the
Redis `GET` command mapped to the HTTP `GET` method.
* Single backing instance: Each instance of the proxy service is associated with a single Redis service instance. The address of the backing Redis is configured at proxy startup via an env variable. 
* Cached GET: This proxy have a caching mechanism with `x` seconds of expiration for each cache. The duration can be set via env variable as well. With size limitation.
* Included `/metrics` URI for prometheus. 
* Global expiry: Entries added to the proy cache are expired after being in cache for a time duration that is globally, configured by environment variable `CACHETTL`
* LRU eviction: Once the cache fills to capacity, the least recently used key is evicted each time a new key needs to be added to the cache.
* Fixed key size: Cache capacity in terms of number of keys retains. Controlled by environment variable `CACHECAPACITY`
* Sequential concurrent processing: Multiple clients are able to concurrently connect to the proxy (Up to configuratble maximum limit `REQUESTLIMIT`) without adversely impacting the functional behaviour of the proxy.
* Configuration to server is done via environment variables. 
* Parallel concurrent processing: The limit is controlled by `BURSTLIMIT`
* Uses redis backing protocol when calling the `get` method to redis when there's a cachemiss. 
* The cache is of form `map[string]struct` complexity. 

## Design decisions 
* Configurations on the proxy is set via environment variables because I expect this to be deployed in docker-like environment. 

## Todo
- Add helm.yaml for easier deployment.
- Dockerize the app for testing. Some functionality cannot be tested with unit test.
- Add terraform for to setup the infrastructure needed. 
