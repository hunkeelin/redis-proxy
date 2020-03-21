package main

import (
	"github.com/go-redis/redis"
	"os"
	"strconv"
)

var (
	redisHost     string  // the backing redis host
	redisPort     string  // teh backing redis port
	ttl           float64 // time to live for each cache item
	cacheSize     int
	redisPassword string
	redisDb       int
	hostPort      string
)

func init() {
	if os.Getenv("REDISHOST") == "" {
		redisHost = "localhost"
	} else {
		redisHost = os.Getenv("REDISHOST")
	}
	if os.Getenv("REDISPORT") == "" {
		redisPort = "6379"
	} else {
		redisPort = os.Getenv("REDISPORT")
	}

	if os.Getenv("CACHETTL") == "" {
		ttl = 30
	} else {
		cachettl, err := strconv.Atoi(os.Getenv("CACHETTL"))
		if err != nil {
			panic(err)
		}
		ttl = float64(cachettl)
	}
	if os.Getenv("HOSTPORT") == "" {
		hostPort = "2020"
	} else {
		_, err := strconv.Atoi(os.Getenv("HOSTPORT"))
		if err != nil {
			panic(err)
		}
		hostPort = os.Getenv("HOSTPORT")
	}

	redisPassword = os.Getenv("REDISPASSWORD")
	if os.Getenv("REDISDB") == "" {
		redisDb = 0
	} else {
		db, err := strconv.Atoi(os.Getenv("REDISDB"))
		if err != nil {
			panic(err)
		}
		redisDb = db
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       redisDb,
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
}
func main() {
	panic(Server())
}
