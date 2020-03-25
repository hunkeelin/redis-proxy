package redisproxy

import (
	log "github.com/sirupsen/logrus"
	"time"
)

// curate: This is the periodic curate function
func (c *conn) curate() error {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()
	log.Info("Curating cache")
	for key, val := range c.cache {
		if time.Now().Sub(val.modifiedAt).Seconds() > c.cachettl {
			delete(c.cache, key)
		}
	}
	return nil
}

// curateLeastUse: Curates the least use cache key or over cachettl
func (c *conn) curateLeastUse() error {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()
	if len(c.cache) < c.cacheCapacity {
		return nil
	}
	log.Info("The cache is full, curating least use")
	var leastusekey string
	var longestliving float64
	now := time.Now()
	for key, val := range c.cache {
		// If the item is close to c.cachettl just delete it.  it's going to curate by the periodic curate anyways, this logic save time.
		if now.Sub(val.modifiedAt).Seconds() > c.cachettl {
			delete(c.cache, key)
			return nil
		}
		if now.Sub(val.modifiedAt).Seconds() > longestliving {
			longestliving = now.Sub(val.modifiedAt).Seconds()
			leastusekey = key
		}
	}
	_, ok := c.cache[leastusekey]
	if ok {
		delete(c.cache, leastusekey)
	} else {
		log.Info("The least use key was deleted: " + leastusekey)
	}
	return nil
}
