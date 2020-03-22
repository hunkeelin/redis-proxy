package redisproxy

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	cachehit = promauto.NewCounter(prometheus.CounterOpts{
		Name: "cachehit",
		Help: "The total number of cachehit",
	})
	cachemiss = promauto.NewCounter(prometheus.CounterOpts{
		Name: "cachemiss",
		Help: "The total number of cachemiss",
	})
	curatequeue = promauto.NewCounter(prometheus.CounterOpts{
		Name: "curatequeue",
		Help: "The total number of periodic getting queued",
	})
	invalidkey = promauto.NewCounter(prometheus.CounterOpts{
		Name: "invalidekey",
		Help: "The total number of invalid key client requested",
	})
	statusOK = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_status_200",
		Help: "The total number of 200 response",
	})
	statusNotFound = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_status_404",
		Help: "The total number of 404 response",
	})
	statusBadRequest = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_status_400",
		Help: "The total number of 400 response",
	})
	statusMethodNotAllowed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_status_405",
		Help: "The total number of 405 response",
	})
)
