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
	_, ok := c.cache["foo"]
	if ok {
		t.Error("cureate test failed foo didn't get curaeted")
	}
}
func TestCurateleastuse(t *testing.T) {
	fmt.Println("testing curateLeastUse()")
	cacheMap := make(map[string]cacheInfo)
	c := conn{
		cache: cacheMap,
	}
	c.cache["foo"] = cacheInfo{
		item:       "bar",
		modifiedAt: time.Now().Add(time.Duration(-20) * time.Second),
	}
	c.cache["foo1"] = cacheInfo{
		item:       "bar1",
		modifiedAt: time.Now().Add(time.Duration(-30) * time.Second),
	}
	cacheCapacity = 1
	ttl = 25
	c.curateLeastUse()
	_, ok := c.cache["foo1"]
	if ok {
		t.Error("cureateleastuse test failed foo1 didn't get curaeted")
	}
}