package common

import (
	"encoding/json"
	"sync"
)

// SyncMap 是一个线程安全的泛型 map，可以被 JSON 序列化
type SyncMap[K comparable, V any] struct {
	mu sync.RWMutex
	m  map[K]V
}

// NewSyncMap 创建一个新的 SyncMap
func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		m: make(map[K]V),
	}
}

// Set 设置键值对
func (sm *SyncMap[K, V]) Set(key K, value V) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.m[key] = value
}

// Get 获取值
func (sm *SyncMap[K, V]) Get(key K) (V, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	value, ok := sm.m[key]
	return value, ok
}

// Delete 删除键值对
func (sm *SyncMap[K, V]) Delete(key K) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.m, key)
}

// MarshalJSON 实现 json.Marshaler 接口
func (sm *SyncMap[K, V]) MarshalJSON() ([]byte, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return json.Marshal(sm.m)
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (sm *SyncMap[K, V]) UnmarshalJSON(data []byte) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return json.Unmarshal(data, &sm.m)
}
