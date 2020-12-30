package storage

import (
	"sync"
)

// ValueHashtable the set of Items
type ValueHashtable struct {
	items map[string]MetricItem
	lock  sync.RWMutex
}

// Put item with value v and key k into the hashtable
func (ht *ValueHashtable) Put(k string, v MetricItem) {
	ht.lock.Lock()
	defer ht.lock.Unlock()
	if ht.items == nil {
		ht.items = make(map[string]MetricItem)
	}
	ht.items[k] = v
}

// Remove item with key k from hashtable
func (ht *ValueHashtable) Remove(k string) {
	ht.lock.Lock()
	defer ht.lock.Unlock()
	delete(ht.items, k)
}

// Get item with key k from the hashtable
func (ht *ValueHashtable) Get(k string) (MetricItem, bool) {
	ht.lock.RLock()
	defer ht.lock.RUnlock()
	value, ok := ht.items[k]
	return value, ok
}

// Size returns the number of the hashtable elements
func (ht *ValueHashtable) Size() int {
	ht.lock.RLock()
	defer ht.lock.RUnlock()
	return len(ht.items)
}

// GetMap from table
func (ht *ValueHashtable) GetMap() map[string]MetricItem {
	return ht.items
}
