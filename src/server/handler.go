package redisproxy

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (c *conn) mainHandler(w http.ResponseWriter, r *http.Request) {
	towrite, statusHeader, _ := c.handle(w, r)
	w.WriteHeader(statusHeader)
	w.Write([]byte(towrite))
}
func (c *conn) handle(w http.ResponseWriter, r *http.Request) (string, int, error) {
	requestTotal.Inc()
	var towrite, requestkey string
	if r.Method != http.MethodGet {
		statusMethodNotAllowed.Inc()
		return "", http.StatusMethodNotAllowed, nil
	}
	if r.Header.Get("rediskey") == "" {
		statusBadRequest.Inc()
		return "", http.StatusBadRequest, nil
	}
	requestkey = r.Header.Get("rediskey")
	cacheitem, ok := c.cacheGet(requestkey)
	towrite = cacheitem.item
	// This is a cache hit
	if ok {
		cachehit.Inc()
		c.cacheUpdate(requestkey)
		return towrite, http.StatusOK, nil
	}
	towrite, err := c.redisClient.Get(requestkey).Result()
	// the requested key is not in cache nor in redis server
	if err != nil {
		log.Info(fmt.Sprintf("Client requesting key value that doesn't exist %s", requestkey))
		invalidkey.Inc()
		statusNotFound.Inc()
		return "", http.StatusNotFound, nil
	}

	cachemiss.Inc()
	// cache over capacity need to curate one
	c.curateLeastUse()
	c.cacheCreate(requestkey, towrite)
	statusOK.Inc()
	return towrite, http.StatusOK, nil
}
