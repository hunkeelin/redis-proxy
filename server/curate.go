package redisproxy

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func (c *conn) curate() error {
	log.Info("Curating cache")
	for key, val := range c.cache {
		if time.Now().Sub(val.modifiedAt).Seconds() > ttl {
			c.cacheMu.Lock()
			delete(c.cache, key)
			c.cacheMu.Unlock()
		}
	}
	return nil
}
func (c *conn) curateLeastUse() {
	log.Info("The cache is full, curating least use")
	var leastusekey string
	var longestliving float64
	now := time.Now()
	for key, val := range c.cache {
		// If the item is close to ttl just delete it.  it's going to curate by the periodic curate anyways, this logic save time.
		if now.Sub(val.modifiedAt).Seconds() > ttl {
			c.cacheMu.Lock()
			delete(c.cache, key)
			c.cacheMu.Unlock()
			return
		}
		if now.Sub(val.modifiedAt).Seconds() > longestliving {
			longestliving = now.Sub(val.modifiedAt).Seconds()
			leastusekey = key
		}
	}
	c.cacheMu.Lock()
	_, ok := c.cache[leastusekey]
	if ok {
		delete(c.cache, leastusekey)
	} else {
		log.Info("The least use key was deleted: " + leastusekey)
	}
	c.cacheMu.Unlock()
}
