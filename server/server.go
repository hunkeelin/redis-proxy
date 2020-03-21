package redisproxy

import (
	"github.com/go-redis/redis"
	"github.com/hunkeelin/mtls/v2/klinserver"
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

// Function to start the server
func Server() error {
	cacheMap := make(map[string]cacheInfo)
	cacheMap["foo"] = cacheInfo{
		modifiedAt: time.Now(),
		item:       "bar",
	}
	c := serverInit{
		cache: cacheMap,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c.mainHandler(w, r)
	})
	mux.Handle("/metrics", promhttp.InstrumentMetricHandler(
		prometheus.DefaultRegisterer, promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}),
	))
	j := &klinserver.ServerConfig{
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
	return klinserver.Server(j)
}
