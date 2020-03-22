package redisproxy

import (
	"fmt"
	"github.com/go-redis/redis"
	httpserver "github.com/hunkeelin/server"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	redisHost     string  // the backing redis host
	redisPort     string  // the backing redis port
	ttl           float64 // time to live for each cache item
	redisPassword string
	redisDb       int
	hostPort      string
	cacheCapacity int
	curateCycle   int
)

func setRedisConfig() error {
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
	redisPassword = os.Getenv("REDISPASSWORD")
	if os.Getenv("REDISDB") == "" {
		redisDb = 0
	} else {
		db, err := strconv.Atoi(os.Getenv("REDISDB"))
		if err != nil {
			return (err)
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
		return (err)
	}
	return nil
}
func setServerConfig() error {

	if os.Getenv("CACHETTL") == "" {
		ttl = 30
	} else {
		cachettl, err := strconv.Atoi(os.Getenv("CACHETTL"))
		if err != nil {
			return (err)
		}
		if cachettl < 1 {
			return (fmt.Errorf("Please give a ttl for each cache of over 1 second"))
		}
		ttl = float64(cachettl)
	}
	if os.Getenv("CACHECAPACITY") == "" {
		cacheCapacity = 10
	} else {
		cachecap, err := strconv.Atoi(os.Getenv("CACHECAPACITY"))
		if err != nil {
			return (err)
		}
		if cachecap < 2 {
			return (fmt.Errorf("Please give a cache capcity of at least 2"))
		}
		cacheCapacity = cachecap
	}
	if os.Getenv("HOSTPORT") == "" {
		hostPort = "2020"
	} else {
		_, err := strconv.Atoi(os.Getenv("HOSTPORT"))
		if err != nil {
			return (err)
		}
		hostPort = os.Getenv("HOSTPORT")
	}
	if os.Getenv("CURATECYCLE") == "" {
		curateCycle = 30
	} else {
		cycle, err := strconv.Atoi(os.Getenv("CURATECYCLE"))
		if err != nil {
			return (err)
		}
		if cycle < 30 || cycle > 600 {
			return (fmt.Errorf("Please have a curate cycle between 30 to 600 seconds"))
		}
		curateCycle = cycle
	}

	return nil
}

func showServerConfig() {
	log.Info("Starting proxy with the following configuration")
	log.Info(fmt.Sprintf("Redis Backing Host: %v", redisHost))
	log.Info(fmt.Sprintf("Redis Backing port: %v", redisPort))
	log.Info(fmt.Sprintf("Cache TTL: %v seconds", ttl))
	log.Info(fmt.Sprintf("Number of Cache:  %v", cacheCapacity))
	log.Info(fmt.Sprintf("Host port: %v", hostPort))
	log.Info(fmt.Sprintf("Curate Cycle: %v seconds", curateCycle))
}

// Server Function to start the server
func Server() error {
	err := setServerConfig()
	if err != nil {
		return err
	}
	err = setRedisConfig()
	if err != nil {
		return err
	}
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
			time.Sleep(time.Duration(curateCycle) * time.Second)
			if atomic.CompareAndSwapUint32(&curateIsRunning, 0, 1) {
				c.curate()
				atomic.StoreUint32(&curateIsRunning, 0)
			} else {
				curatequeue.Inc()
			}
		}
	}()
	showServerConfig()
	return httpserver.Server(j)
}
