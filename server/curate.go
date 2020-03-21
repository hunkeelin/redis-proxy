package redisproxy

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func (c *serverInit) curate() error {
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
