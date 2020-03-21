package main

import (
	"github.com/hunkeelin/mtls/v2/klinserver"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"sync/atomic"
	"time"
	//	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

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
