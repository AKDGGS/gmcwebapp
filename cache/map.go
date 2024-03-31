package cache

import (
	"sync"
	"time"
)

func NewMap() Map {
	return Map{entries: map[string]*Entry{}}
}

type Map struct {
	entries map[string]*Entry
	mutex   sync.RWMutex
}

func (m *Map) Get(name string) *Entry {
	m.mutex.RLock()
	entry, ok := m.entries[name]
	m.mutex.RUnlock()
	if !ok {
		return nil
	}

	if time.Now().After(entry.Expires) {
		m.Remove(name)
		return nil
	}
	return entry
}

func (m *Map) Put(name string, entry *Entry) {
	// Refuse to set nil or nil content items
	if entry == nil || entry.pl_content == nil {
		return
	}

	m.mutex.Lock()
	m.entries[name] = entry
	m.mutex.Unlock()
}

func (m *Map) Remove(name string) {
	m.mutex.Lock()
	delete(m.entries, name)
	m.mutex.Unlock()
}

func (m *Map) PurgeExpired() {
	now := time.Now()
	m.mutex.Lock()
	for k, v := range m.entries {
		if !v.Expires.IsZero() && now.After(v.Expires) {
			delete(m.entries, k)
		}
	}
	m.mutex.Unlock()
}
