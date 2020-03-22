package redisproxy

import (
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"testing"
)

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
func TestSetserverconfig(t *testing.T) {
	fmt.Println("testing setServerConfig()")
	err := setServerConfig()
	if err != nil {
		fmt.Println(err)
	}
	showServerConfig()
}
func TestSetredisconfig(t *testing.T) {
	fmt.Println("testing setRedisConfig()")
	os.Setenv("REDISPORT", "1234")
	err := setRedisConfig()
	if err != nil {
		fmt.Println(err)
	}
	showServerConfig()
}
