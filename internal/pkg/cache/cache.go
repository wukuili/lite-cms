package cache

import (
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/lite-cms/cms/internal/config"
)

// Cache 全局缓存实例
var Cache *ristretto.Cache

// Init 初始化缓存
func Init(cfg *config.CacheConfig) (*ristretto.Cache, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: cfg.NumCounters, // 跟踪key数量的10倍
		MaxCost:     cfg.MaxCost,     // 最大内存占用
		BufferItems: cfg.BufferItems, // 每个Get缓冲区大小
	})
	if err != nil {
		return nil, err
	}

	Cache = cache
	return cache, nil
}

// Get 获取缓存
func Get(key string) (interface{}, bool) {
	if Cache == nil {
		return nil, false
	}
	return Cache.Get(key)
}

// Set 设置缓存
func Set(key string, value interface{}, cost int64) bool {
	if Cache == nil {
		return false
	}
	return Cache.Set(key, value, cost)
}

// SetWithTTL 设置带过期时间的缓存
func SetWithTTL(key string, value interface{}, cost int64, ttl time.Duration) bool {
	if Cache == nil {
		return false
	}
	return Cache.SetWithTTL(key, value, cost, ttl)
}

// Del 删除缓存
func Del(key string) {
	if Cache == nil {
		return
	}
	Cache.Del(key)
}

// Clear 清空缓存
func Clear() {
	if Cache == nil {
		return
	}
	Cache.Clear()
}
