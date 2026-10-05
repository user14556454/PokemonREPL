package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	Cachemap map[string]CacheEntry
	Mu       sync.Mutex
}

type CacheEntry struct {
	CreatedAt time.Time
	Val       []byte
}

func (ch *Cache) Add(key string, value []byte) {
	newmapentry := CacheEntry{
		CreatedAt: time.Now(),
		Val:       value,
	}
	ch.Mu.Lock()
	defer ch.Mu.Unlock()
	ch.Cachemap[key] = newmapentry
}

func (ch *Cache) Get(key string) ([]byte, bool) {
	ch.Mu.Lock()
	defer ch.Mu.Unlock()
	for k := range ch.Cachemap {
		if k == key {
			return ch.Cachemap[key].Val, true
		}
	}
	return []byte{}, false
}

func (ch *Cache) reapLoop(interval time.Duration) {
	for {
		time.Sleep(interval)
		t := time.Now()
		ch.Mu.Lock()
		for k := range ch.Cachemap {
			if t.Sub(ch.Cachemap[k].CreatedAt) >= interval {
				delete(ch.Cachemap, k)
			}
		}
		ch.Mu.Unlock()
	}
}

func NewCache(interval time.Duration) *Cache {
	newcache := Cache{
		Cachemap: map[string]CacheEntry{},
		Mu:       sync.Mutex{},
	}
	go newcache.reapLoop(interval)
	return &newcache
}
