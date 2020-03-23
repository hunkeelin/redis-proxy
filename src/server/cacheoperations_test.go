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
		t.Errorf("cacheUpdate failed, foo's time didn't get update")
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
	if !ok {
		t.Errorf("Cache create test failed foo77 didn't get created")
	}
}
func TestCacheget(t *testing.T) {
	fmt.Println("Testing cacheGet()")
	cacheMap := make(map[string]cacheInfo)
	c := conn{
		cache: cacheMap,
	}
	c.cache["foo7"] = cacheInfo{
		item: "bar7",
	}
	_, ok := c.cacheGet("foo7")
	if !ok {
		t.Errorf("Cache get test failed foo7 didn't return a value")
	}
}
