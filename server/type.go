package redisproxy

import (
	"sync"
	"time"
)

type serverInit struct {
	cache    map[string]cacheInfo
	cacheMu  sync.Mutex // Mutex for any cache related operation
	curateMu sync.Mutex // Mutex for any curate related operation
}
type cacheInfo struct {
	modifiedAt time.Time
	item       string
}
