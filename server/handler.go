package redisproxy

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
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
	cacheitem, ok := c.cache[requestkey]
	val = cacheitem.item
	if ok {
		cachehit.Inc()
		c.cacheMu.Lock()
		tmp := c.cache[requestkey]
		tmp.modifiedAt = time.Now()
		c.cache[requestkey] = tmp
		c.cacheMu.Unlock()
		w.Write([]byte(val))
	} else {
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
		if len(c.cache) > cacheCapacity {
			c.curateLeastUse()
		}
		w.Write([]byte(val))
		c.cacheMu.Lock()
		c.cache[requestkey] = cacheInfo{
			item:       val,
			modifiedAt: time.Now(),
		}
		c.cacheMu.Unlock()
	}
	statusOK.Inc()
	return
}
