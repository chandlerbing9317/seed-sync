package common

import (
	"encoding/json"
)

// Cache 全局缓存实例
// 只对外暴露函数
var globalCache = newCache()

// Cache 泛型缓存结构体
type Cache struct {
	store *SyncMap[string, string]
}

// newCache 创建新的缓存实例
func newCache() *Cache {
	return &Cache{
		store: NewSyncMap[string, string](),
	}
}

// Get 获取字符串值
func CacheGet(key string) string {
	value, _ := globalCache.store.Get(key)
	return value
}

// GetObject 获取并反序列化为指定类型对象
func CacheGetObject[T any](key string) (T, bool) {
	var result T
	value, exists := globalCache.store.Get(key)
	if !exists {
		return result, false
	}

	err := json.Unmarshal([]byte(value), &result)
	if err != nil {
		return result, false
	}
	return result, true
}

// GetList 获取并反序列化为指定类型的切片
func CacheGetList[T any](key string) ([]T, bool) {
	var result []T
	value, exists := globalCache.store.Get(key)
	if !exists {
		return result, false
	}

	err := json.Unmarshal([]byte(value), &result)
	if err != nil {
		return result, false
	}
	return result, true
}

// Set 设置字符串值
func CacheSet(key string, value string) {
	globalCache.store.Set(key, value)
}

// SetObject 序列化对象并存储
func CacheSetObject(key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	globalCache.store.Set(key, string(data))
	return nil
}

// Remove 删除缓存项
func CacheRemove(key string) {
	globalCache.store.Delete(key)
}
