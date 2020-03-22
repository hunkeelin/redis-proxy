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
	"sync"
	"sync/atomic"
	"time"
)

type conn struct {
	cache       map[string]cacheInfo
	cacheMu     sync.Mutex    // Mutex for any cache related operation
	curateMu    sync.Mutex    // Mutex for any curate related operation
	redisClient *redis.Client // The redis client
	// configurations
	redisHost     string  // the backing redis host
	redisPort     string  // the backing redis port
	cachettl      float64 // time to live for each cache item
	redisPassword string
	redisDb       int
	hostPort      string
	cacheCapacity int
	curateCycle   int
}
type cacheInfo struct {
	modifiedAt time.Time
	item       string
}

func (c *conn) setRedisConfig() error {
	if os.Getenv("REDISHOST") == "" {
		c.redisHost = "localhost"
	} else {
		c.redisHost = os.Getenv("REDISHOST")
	}
	if os.Getenv("REDISPORT") == "" {
		c.redisPort = "6379"
	} else {
		c.redisPort = os.Getenv("REDISPORT")
	}
	c.redisPassword = os.Getenv("REDISPASSWORD")
	if os.Getenv("REDISDB") == "" {
		c.redisDb = 0
	} else {
		db, err := strconv.Atoi(os.Getenv("REDISDB"))
		if err != nil {
			return (err)
		}
		c.redisDb = db
	}
	return nil
}
func (c *conn) setServerConfig() error {
	if os.Getenv("CACHETTL") == "" {
		c.cachettl = 30
	} else {
		cachettl, err := strconv.Atoi(os.Getenv("CACHETTL"))
		if err != nil {
			return (err)
		}
		if cachettl < 1 {
			return (fmt.Errorf("Please give a c.cachettl for each cache of over 1 second"))
		}
		c.cachettl = float64(cachettl)
	}
	if os.Getenv("CACHECAPACITY") == "" {
		c.cacheCapacity = 10
	} else {
		cachecap, err := strconv.Atoi(os.Getenv("CACHECAPACITY"))
		if err != nil {
			return (err)
		}
		if cachecap < 2 {
			return (fmt.Errorf("Please give a cache capcity of at least 2"))
		}
		c.cacheCapacity = cachecap
	}
	if os.Getenv("HOSTPORT") == "" {
		c.hostPort = "2020"
	} else {
		_, err := strconv.Atoi(os.Getenv("HOSTPORT"))
		if err != nil {
			return (err)
		}
		c.hostPort = os.Getenv("HOSTPORT")
	}
	if os.Getenv("CURATECYCLE") == "" {
		c.curateCycle = 30
	} else {
		cycle, err := strconv.Atoi(os.Getenv("CURATECYCLE"))
		if err != nil {
			return (err)
		}
		if cycle < 30 || cycle > 600 {
			return (fmt.Errorf("Please have a curate cycle between 30 to 600 seconds"))
		}
		c.curateCycle = cycle
	}

	return nil
}

func (c *conn) showServerConfig() {
	log.Info("Starting proxy with the following configuration")
	log.Info(fmt.Sprintf("Redis Backing Host: %v", c.redisHost))
	log.Info(fmt.Sprintf("Redis Backing port: %v", c.redisPort))
	log.Info(fmt.Sprintf("Cache TTL: %v seconds", c.cachettl))
	log.Info(fmt.Sprintf("Number of Cache:  %v", c.cacheCapacity))
	log.Info(fmt.Sprintf("Host port: %v", c.hostPort))
	log.Info(fmt.Sprintf("Curate Cycle: %v seconds", c.curateCycle))
}

// Server Function to start the server
func Server() error {
	c := conn{}
	err := c.setServerConfig()
	if err != nil {
		return err
	}
	err = c.setRedisConfig()
	if err != nil {
		return err
	}
	client := redis.NewClient(&redis.Options{
		Addr:     c.redisHost + ":" + c.redisPort,
		Password: c.redisPassword,
		DB:       c.redisDb,
	})
	_, err = client.Ping().Result()
	if err != nil {
		return err
	}

	cacheMap := make(map[string]cacheInfo)
	c.cache = cacheMap
	c.redisClient = client
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c.mainHandler(w, r)
	})
	mux.Handle("/metrics", promhttp.InstrumentMetricHandler(
		prometheus.DefaultRegisterer, promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}),
	))
	j := &httpserver.ServerConfig{
		BindPort: c.hostPort,
		BindAddr: "",
		ServeMux: mux,
		Https:    false,
	}

	// Curating every 30 seconds
	var curateIsRunning uint32
	go func() {
		for {
			time.Sleep(time.Duration(c.curateCycle) * time.Second)
			if atomic.CompareAndSwapUint32(&curateIsRunning, 0, 1) {
				c.curate()
				atomic.StoreUint32(&curateIsRunning, 0)
			} else {
				curatequeue.Inc()
			}
		}
	}()
	c.showServerConfig()
	return httpserver.Server(j)
}
