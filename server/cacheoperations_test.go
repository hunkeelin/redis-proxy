package redisproxy

import (
	"fmt"
	"testing"
	"time"
)

func TestCacheupdate(t *testing.T) {
	fmt.Println("Testing cacheUpdate()")
	cacheMap := make(map[string]cacheInfo)
	c := conn{
		cache: cacheMap,
	}
	c.cache["foo"] = cacheInfo{
		item:       "bar",
		modifiedAt: time.Now().Add(time.Duration(-60) * time.Minute),
	}
	old := c.cache["foo"].modifiedAt
	c.cacheUpdate("foo")
	new := c.cache["foo"].modifiedAt
	if new.Sub(old).Minutes() < 60 {
		fmt.Println("Cache update test failed")
	} else {
		fmt.Println("Cache update test passed")
	}

}
func TestCachecreate(t *testing.T) {
	fmt.Println("Testing cacheCreate()")
	cacheMap := make(map[string]cacheInfo)
	c := conn{
		cache: cacheMap,
	}
	c.cacheCreate("foo77", "bar88")
	_, ok := c.cache["foo77"]
	if ok {
		fmt.Println("Cache create test passed")
	} else {
		fmt.Println("Cache create test failed")
	}
}
