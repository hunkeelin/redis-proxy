package redisproxy

import (
	"fmt"
	"os"
	"testing"
)

func TestSetserverconfig(t *testing.T) {
	fmt.Println("testing setServerConfig()")
	os.Setenv("CACHETTL", "31")
	err := setServerConfig()
	if err != nil {
		t.Errorf(err.Error())
	}
	if ttl != 31 {
		t.Errorf("setServerConfig failed, ttl wasn't 31")
	}
}
func TestSetredisconfig(t *testing.T) {
	fmt.Println("testing setRedisConfig()")
	os.Setenv("REDISPORT", "1234")
	err := setRedisConfig()
	if err != nil {
		t.Errorf(err.Error())
	}
	if redisPort != "1234" {
		t.Errorf("setRedisConfig failed, redisport wasn't 1234")
	}
}
