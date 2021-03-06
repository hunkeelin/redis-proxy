package redisproxy

import (
	"time"
)

func (c *conn) cacheUpdate(requestkey string) {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()
	tmp := c.cache[requestkey]
	tmp.modifiedAt = time.Now()
	c.cache[requestkey] = tmp
}
func (c *conn) cacheCreate(key, val string) {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()
	c.cache[key] = cacheInfo{
		item:       val,
		modifiedAt: time.Now(),
	}
}
func (c *conn) cacheGet(key string) (cacheInfo, bool) {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()
	val, ok := c.cache[key]
	return val, ok
}
