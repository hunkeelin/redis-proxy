package redisproxy

import (
	"fmt"
	"github.com/go-redis/redis"
	httpserver "github.com/hunkeelin/server"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"
	//	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	redisHost     string  // the backing redis host
	redisPort     string  // the backing redis port
	ttl           float64 // time to live for each cache item
	cacheSize     int
	redisPassword string
	redisDb       int
	hostPort      string
	cacheCapacity int
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
		if cachettl < 1 {
			panic(fmt.Errorf("Please give a ttl for each cache of over 1 second"))
		}
		ttl = float64(cachettl)
	}
	if os.Getenv("CACHECAPACITY") == "" {
		cacheCapacity = 10
	} else {
		cachecap, err := strconv.Atoi(os.Getenv("CACHECAPACITY"))
		if err != nil {
			panic(err)
		}
		if cachecap < 10 {
			panic(fmt.Errorf("Please give a cache capcity of at least 10"))
		}
		cacheCapacity = cachecap
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

// Function to start the server
func Server() error {
	cacheMap := make(map[string]cacheInfo)
	client := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       redisDb,
	})
	c := conn{
		cache:       cacheMap,
		redisClient: client,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c.mainHandler(w, r)
	})
	mux.Handle("/metrics", promhttp.InstrumentMetricHandler(
		prometheus.DefaultRegisterer, promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}),
	))
	j := &httpserver.ServerConfig{
		BindPort: hostPort,
		BindAddr: "",
		ServeMux: mux,
		Https:    false,
	}

	// Curating every 30 seconds
	var curateIsRunning uint32
	go func() {
		for {
			time.Sleep(30 * time.Second)
			if atomic.CompareAndSwapUint32(&curateIsRunning, 0, 1) {
				c.curate()
				atomic.StoreUint32(&curateIsRunning, 0)
			}
		}
	}()
	return httpserver.Server(j)
}
