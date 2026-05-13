package store

import "sync"

var (
	mu    sync.RWMutex
	store = map[string]bool{}
)

func Set(flag string, enabled bool) {
	mu.Lock()
	defer mu.Unlock()
	store[flag] = enabled
}

func IsEnabled(flag string) bool {
	mu.RLock()
	defer mu.RUnlock()
	val, exists := store[flag]
	if !exists {
		return true
	}
	return val
}

func LoadAll(incoming map[string]bool) {
	mu.Lock()
	defer mu.Unlock()
	for k, v := range incoming {
		store[k] = v
	}
}
