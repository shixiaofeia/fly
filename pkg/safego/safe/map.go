package safe

import (
	"sync"
)

type Map struct {
	val map[interface{}]interface{}
	mu  sync.RWMutex
}

func NewMap() *Map {
	return &Map{val: make(map[interface{}]interface{})}
}

func (slf *Map) Add(k, v interface{}) {
	slf.mu.Lock()
	defer slf.mu.Unlock()

	slf.val[k] = v
}

func (slf *Map) Del(k interface{}) {
	slf.mu.Lock()
	defer slf.mu.Unlock()

	delete(slf.val, k)
}

func (slf *Map) Get(k interface{}) (interface{}, bool) {
	slf.mu.RLock()
	defer slf.mu.RUnlock()
	if val, ok := slf.val[k]; ok {
		return val, ok
	}
	return nil, false
}

func (slf *Map) Len() int {
	slf.mu.RLock()
	defer slf.mu.RUnlock()

	return len(slf.val)
}

func (slf *Map) Range(f func(k, v interface{}) bool) bool {
	slf.mu.RLock()
	defer slf.mu.RUnlock()

	for k, v := range slf.val {
		if !f(k, v) {
			return false
		}
	}

	return true
}
