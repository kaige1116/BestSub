package node

import "sync"

type exist struct {
	mu   sync.RWMutex
	data map[uint64]struct{}
}

func NewExist(size int) *exist {
	return &exist{data: make(map[uint64]struct{}, size)}
}

func (k *exist) Exist(key uint64) bool {
	k.mu.RLock()
	_, exists := k.data[key]
	k.mu.RUnlock()
	return exists
}

func (k *exist) Add(key uint64) {
	k.mu.Lock()
	k.data[key] = struct{}{}
	k.mu.Unlock()
}
func (k *exist) Remove(key uint64) {
	k.mu.Lock()
	delete(k.data, key)
	k.mu.Unlock()
}
