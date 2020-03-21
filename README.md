# RedisProxy 
[![CircleCI](https://circleci.com/gh/hunkeelin/redis-proxy.svg?style=shield)](https://circleci.com/gh/hunkeelin/redis-proxy)
[![Go Report Card](https://goreportcard.com/badge/github.com/hunkeelin/redis-proxy)](https://goreportcard.com/report/github.com/hunkeelin/redis-proxy)
[![GoDoc](https://godoc.org/github.com/hunkeelin/redis-proxy/server?status.svg)](https://godoc.org/github.com/hunkeelin/redis-proxy/server)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/hunkeelin/redis-proxy/master/LICENSE)


## Motivations

This project implements additional features on top of redis

## Golang version

`redis-proxy` is currently compatible with golang version from 1.12+.

## Features
* HTTP webservice: Clients interface to the Redis proxy through HTTP, with the
Redis “GET” command mapped to the HTTP “GET” method.
Note that the proxy still uses the Redis protocol to
communicate with the backend Redis server.
* Single backing instance: Each instance of the proxy service is associated with a single Redis service instance. The address of the backing Redis is configured at proxy startup via an env variable. 
* Cached GET: This proxy have a caching mechanism with 5 minutes of expiration for each cache. The duration can be set via env variable as well. With size limitation
