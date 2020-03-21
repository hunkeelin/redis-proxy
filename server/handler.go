package redisproxy

import (
	"fmt"
	"net/http"
	"time"
)

func (c *conn) mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("rediskey") == "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	diff := time.Now().Sub(c.cache["foo"].modifiedAt)
	fmt.Println(diff.Seconds(), diff.Minutes())
	return
}
