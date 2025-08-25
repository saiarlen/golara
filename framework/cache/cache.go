package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache interface defines caching operations
type Cache interface {
	Get(key string, dest interface{}) error
	Set(key string, value interface{}, ttl time.Duration) error
	Delete(key string) error
	Flush() error
	Remember(key string, ttl time.Duration, callback func() (interface{}, error), dest interface{}) error
}

// CacheManager manages different cache stores
type CacheManager struct {
	stores  map[string]Cache
	default_ string
}

// NewCacheManager creates a new cache manager
func NewCacheManager() *CacheManager {
	return &CacheManager{
		stores: make(map[string]Cache),
	}
}

// AddStore adds a cache store
func (cm *CacheManager) AddStore(name string, store Cache) {
	cm.stores[name] = store
	if cm.default_ == "" {
		cm.default_ = name
	}
}

// Store returns a cache store by name
func (cm *CacheManager) Store(name ...string) Cache {
	storeName := cm.default_
	if len(name) > 0 {
		storeName = name[0]
	}
	return cm.stores[storeName]
}

// RedisCache implements Redis caching
type RedisCache struct {
	client *redis.Client
	prefix string
}

// NewRedisCache creates a new Redis cache
func NewRedisCache(addr, password string, db int, prefix string) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	
	return &RedisCache{
		client: rdb,
		prefix: prefix,
	}
}

func (rc *RedisCache) key(key string) string {
	if rc.prefix != "" {
		return fmt.Sprintf("%s:%s", rc.prefix, key)
	}
	return key
}

func (rc *RedisCache) Get(key string, dest interface{}) error {
	ctx := context.Background()
	val, err := rc.client.Get(ctx, rc.key(key)).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return err
	}
	
	return json.Unmarshal([]byte(val), dest)
}

func (rc *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
	ctx := context.Background()
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	
	return rc.client.Set(ctx, rc.key(key), data, ttl).Err()
}

func (rc *RedisCache) Delete(key string) error {
	ctx := context.Background()
	return rc.client.Del(ctx, rc.key(key)).Err()
}

func (rc *RedisCache) Flush() error {
	ctx := context.Background()
	return rc.client.FlushDB(ctx).Err()
}

func (rc *RedisCache) Remember(key string, ttl time.Duration, callback func() (interface{}, error), dest interface{}) error {
	// Try to get from cache first
	err := rc.Get(key, dest)
	if err == nil {
		return nil
	}
	
	if err != ErrCacheMiss {
		return err
	}
	
	// Cache miss, execute callback
	value, err := callback()
	if err != nil {
		return err
	}
	
	// Store in cache
	if err := rc.Set(key, value, ttl); err != nil {
		return err
	}
	
	// Marshal and unmarshal to populate dest
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, dest)
}

// MemoryCache implements in-memory caching
type MemoryCache struct {
	data   map[string]cacheItem
	mutex  sync.RWMutex
	prefix string
}

type cacheItem struct {
	value     interface{}
	expiresAt time.Time
}

// NewMemoryCache creates a new memory cache
func NewMemoryCache(prefix string) *MemoryCache {
	mc := &MemoryCache{
		data:   make(map[string]cacheItem),
		prefix: prefix,
	}
	
	// Start cleanup goroutine
	go mc.cleanup()
	
	return mc
}

func (mc *MemoryCache) key(key string) string {
	if mc.prefix != "" {
		return fmt.Sprintf("%s:%s", mc.prefix, key)
	}
	return key
}

func (mc *MemoryCache) Get(key string, dest interface{}) error {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()
	
	item, exists := mc.data[mc.key(key)]
	if !exists || time.Now().After(item.expiresAt) {
		return ErrCacheMiss
	}
	
	// Use JSON marshal/unmarshal for type conversion
	data, err := json.Marshal(item.value)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, dest)
}

func (mc *MemoryCache) Set(key string, value interface{}, ttl time.Duration) error {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	
	expiresAt := time.Now().Add(ttl)
	if ttl == 0 {
		expiresAt = time.Now().Add(24 * time.Hour) // Default 24 hours
	}
	
	mc.data[mc.key(key)] = cacheItem{
		value:     value,
		expiresAt: expiresAt,
	}
	
	return nil
}

func (mc *MemoryCache) Delete(key string) error {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	
	delete(mc.data, mc.key(key))
	return nil
}

func (mc *MemoryCache) Flush() error {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	
	mc.data = make(map[string]cacheItem)
	return nil
}

func (mc *MemoryCache) Remember(key string, ttl time.Duration, callback func() (interface{}, error), dest interface{}) error {
	// Try to get from cache first
	err := mc.Get(key, dest)
	if err == nil {
		return nil
	}
	
	if err != ErrCacheMiss {
		return err
	}
	
	// Cache miss, execute callback
	value, err := callback()
	if err != nil {
		return err
	}
	
	// Store in cache
	if err := mc.Set(key, value, ttl); err != nil {
		return err
	}
	
	// Marshal and unmarshal to populate dest
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, dest)
}

func (mc *MemoryCache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		mc.mutex.Lock()
		now := time.Now()
		for key, item := range mc.data {
			if now.After(item.expiresAt) {
				delete(mc.data, key)
			}
		}
		mc.mutex.Unlock()
	}
}

// ErrCacheMiss indicates cache miss
var ErrCacheMiss = fmt.Errorf("cache miss")