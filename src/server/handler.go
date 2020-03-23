package redisproxy

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (c *conn) mainHandler(w http.ResponseWriter, r *http.Request) {
	requestTotal.Inc()
	var val, requestkey string
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		statusMethodNotAllowed.Inc()
		return
	}
	if r.Header.Get("rediskey") == "" {
		w.WriteHeader(http.StatusBadRequest)
		statusBadRequest.Inc()
		return
	}
	requestkey = r.Header.Get("rediskey")
	cacheitem, ok := c.cacheGet(requestkey)
	val = cacheitem.item
	// This is a cache hit
	if ok {
		cachehit.Inc()
		c.cacheUpdate(requestkey)
		w.Write([]byte(val))
		return
	}
	val, err := c.redisClient.Get(requestkey).Result()
	// the requested key is not in cache nor in redis server
	if err != nil {
		log.Info(fmt.Sprintf("Client requesting key value that doesn't exist %s", requestkey))
		invalidkey.Inc()
		w.WriteHeader(http.StatusNotFound)
		statusNotFound.Inc()
		return
	}

	cachemiss.Inc()
	// cache over capacity need to curate one
	c.curateLeastUse()
	w.Write([]byte(val))
	c.cacheCreate(requestkey, val)
	statusOK.Inc()
	return
}
