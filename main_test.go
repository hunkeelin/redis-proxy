package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"testing"
)

func TestHello(t *testing.T) {
	fmt.Println("hello")
}
func TestRedis(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	val, err := client.Get("foo").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)
}
