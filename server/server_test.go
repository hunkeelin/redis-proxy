package redisproxy

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
		fmt.Println("no local redis server")
		return
	}
	val, err := client.Get("foo7").Result()
	if err != nil {
		fmt.Println("no such value in redis server")
		return
	}
	fmt.Println("foo7", val)
}
func TestSetconfig(t *testing.T) {
	err := setServerConfig()
	if err != nil {
		fmt.Println(err)
	}
}
func TestShowconfig(t *testing.T) {
	showServerConfig()
}
