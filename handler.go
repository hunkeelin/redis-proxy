package main

import (
	"fmt"
	"net/http"
	"time"
)

func (c *serverInit) mainHandler(w http.ResponseWriter, r *http.Request) {
	//	if r.GetHeader("rediskey") == "" {
	//		w.WriteHeader(http.StatusBadRequest)
	//	}
	diff := time.Now().Sub(c.cache["foo"].modifiedAt)
	fmt.Println(diff.Seconds(), diff.Minutes())
	return
}
