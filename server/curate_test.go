package redisproxy

import (
	"fmt"
	"testing"
	"time"
)

func TestCurate(t *testing.T) {
	fmt.Println("testing curate()")
	cacheMap := make(map[string]cacheInfo)
	c := conn{
		cache: cacheMap,
	}
	c.cache["foo"] = cacheInfo{
		item:       "bar",
		modifiedAt: time.Now().Add(time.Duration(-60) * time.Minute),
	}
	ttl = 30 // setting ttl for the cache to be 30 seconds
	c.curate()
	_, ok = c.cache["foo"]
	if ok {
		fmt.Println("curate test failed")
		return
	} else {
		fmt.Println("curate() test pass")
	}
}
func TestCurateleastuse(t *testing.T) {
	cacheMap := make(map[string]cacheInfo)
	c := conn{
		cache: cacheMap,
	}
	c.cache["foo"] = cacheInfo{
		item:       "bar",
		modifiedAt: time.Now().Add(time.Duration(-20) * time.Minute),
	}
	c.cache["foo1"] = cacheInfo{
		item:       "bar1",
		modifiedAt: time.Now().Add(time.Duration(-30) * time.Minute),
	}
	c.curateLeastUse()
	_, ok := c.cache["foo1"]
	if ok {
		fmt.Println("curateLeastUse() test failed")
	} else {
		fmt.Println("curateLeastUse() test passed")
	}
}
