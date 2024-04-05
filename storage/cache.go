package storage

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	cacheOnce     sync.Once
	cacheInstance *cache.Cache
)

func PieceCache() *cache.Cache {
	cacheOnce.Do(func() {
		cacheInstance = cache.New(time.Minute, time.Second)
	})
	return cacheInstance
}
