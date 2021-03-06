package redisproxy

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func (c *conn) mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("yes")
	towrite, statusHeader, _ := c.handle(w, r)
	w.WriteHeader(statusHeader)
	w.Write([]byte(towrite))
}
func (c *conn) faultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("no")
	w.WriteHeader(http.StatusBadRequest)
}

// handle
// Note: For now error only returns nil, it can be expand later when it is needed, a go function should always have an error to return.
func (c *conn) handle(w http.ResponseWriter, r *http.Request) (string, int, error) {
	requestTotal.Inc()
	var towrite, requestkey string
	if r.Method != http.MethodGet {
		statusMethodNotAllowed.Inc()
		return "", http.StatusMethodNotAllowed, nil
	}
	if len(r.RequestURI) < 4 {
		log.Error("Request url have a length lower than 4 something is seriously wrong. " + r.RequestURI)
		return "", http.StatusBadRequest, nil
	}

	if !strings.Contains(r.RequestURI, "/key/") {
		log.Error("The inintial URI directory is not /key/ something is seriously wrong. " + r.RequestURI)
		return "", http.StatusBadRequest, nil
	}

	requestkey = strings.Trim(r.RequestURI, "/key/")
	if requestkey == "" {
		statusBadRequest.Inc()
		return "", http.StatusBadRequest, nil
	}
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
		log.Warning(fmt.Sprintf("Client requesting key value that doesn't exist %s", requestkey))
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
