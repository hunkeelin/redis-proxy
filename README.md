# RedisProxy 
[![CircleCI](https://circleci.com/gh/hunkeelin/redis-proxy.svg?style=shield)](https://circleci.com/gh/hunkeelin/redis-proxy)
[![Go Report Card](https://goreportcard.com/badge/github.com/hunkeelin/redis-proxy)](https://goreportcard.com/report/github.com/hunkeelin/redis-proxy)
[![GoDoc](https://godoc.org/github.com/hunkeelin/redis-proxy/server?status.svg)](https://godoc.org/github.com/hunkeelin/redis-proxy/server)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/hunkeelin/redis-proxy/master/LICENSE)


## Motivations

This project implements additional features on top of redis

## Golang version

`redis-proxy` is currently compatible with golang version from 1.12+.

## Usage
```go
package main

import (
    redisproxy "github.com/hunkeelin/redis-proxy/server"
)

func main() {
    panic(redisproxy.Server())
}
```

## Features
* HTTP webservice: Clients interface to the Redis proxy through HTTP, with the
Redis `GET` command mapped to the HTTP `GET` method.
* Single backing instance: Each instance of the proxy service is associated with a single Redis service instance. The address of the backing Redis is configured at proxy startup via an env variable. 
* Cached GET: This proxy have a caching mechanism with 5 minutes of expiration for each cache. The duration can be set via env variable as well. With size limitation/ 
* Included `/metrics` URI for prometheus. 
* configuratable server via env variables. 

## Design decisions 
* Since this is just an exercise it will be up to me to decide on things and I got no way to gather requirements during my offline coding session. Sequential concurrent processing doesn't make sense in this case.(A request from the second request only starts processing after the first request has completed and a response has been returned first client) This basically means adding a mutex to the entire handler which is against `go` philosphy. If race condition is an issue we can use threadsafe techiques such as `sync/atomic` and `sync/mutex` on write operations. An example of Sequential concurrent processing is not optimal is when redis server is timing out, causing one request to halt. This will cause every subsequence request to halt, even though they could be served from cache while hoping redis comes back up. 
* Configurations on the proxy is set via envirnoment variables because I expect this to be deployed in docker-like envirnoment. 

