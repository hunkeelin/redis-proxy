package redisproxy

import (
	"fmt"
	"os"
	"testing"
)

func TestSetserverconfig(t *testing.T) {
	fmt.Println("testing setServerConfig()")
	os.Setenv("CACHETTL", "31")
	c := conn{}
	err := c.setServerConfig()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if c.cachettl != 31 {
		t.Errorf("setServerConfig failed, c.cachettl wasn't 31 it was %v", c.cachettl)
	}
}
func TestSetredisconfig(t *testing.T) {
	fmt.Println("testing setRedisConfig()")
	os.Setenv("REDISPORT", "1234")
	c := conn{}
	err := c.setRedisConfig()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if c.redisPort != "1234" {
		t.Errorf("setRedisConfig failed, redisport wasn't 1234")
	}
}
