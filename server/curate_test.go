package redisproxy

import (
	"fmt"
	"testing"
	"time"
)

func TestCurate(t *testing.T) {
	cacheMap := make(map[string]cacheInfo)
	c := conn{
		cache: cacheMap,
	}
	c.cache["foo"] = cacheInfo{
		item:       "bar",
		modifiedAt: time.Now().Add(time.Duration(-60) * time.Minute),
	}
	_, ok := c.cache["foo"]
	if ok {
		fmt.Println("cache with key value foo currently exist")
		fmt.Println(c.cache["foo"].item)
	}
	ttl = 30 // setting ttl for the cache to be 30 seconds
	fmt.Println("curating all items with more then 30 seconds old")
	c.curate()
	_, ok = c.cache["foo"]
	if ok {
		fmt.Println("foo didn't get curated, curate() is bugged")
	} else {
		fmt.Println("foo is curated test pass")
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
	_, ok := c.cache["foo"]
	fmt.Println("key foo exist?", ok)
	_, ok = c.cache["foo1"]
	fmt.Println("key foo1 exist?", ok)
	fmt.Println("Running curateLeastUse(), foo1 should get curated because it's least used")
	c.curateLeastUse()
	_, ok = c.cache["foo"]
	fmt.Println("key foo exist?", ok)
	_, ok = c.cache["foo1"]
	fmt.Println("key foo1 exist?", ok)
}
