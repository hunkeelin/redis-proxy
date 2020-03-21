package redisproxy

import (
	"github.com/go-redis/redis"
	"sync"
	"time"
)

type conn struct {
	cache       map[string]cacheInfo
	cacheMu     sync.Mutex    // Mutex for any cache related operation
	curateMu    sync.Mutex    // Mutex for any curate related operation
	redisClient *redis.Client // The redis client
}
type cacheInfo struct {
	modifiedAt time.Time
	item       string
}
