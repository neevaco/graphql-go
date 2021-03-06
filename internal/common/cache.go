package common

import (
	"hash/maphash"
	"math"
	"sync"
)

// LRUKeyCache implements a cache for a set of uint64 keys.
// It has a least-recently-used eviction policy.
// This is designed for a smallish number of keys.
//
// The keys are uint64 rather than string to save memory; in our usage, our strings may be large,
// and we only need to check for existence.
type LRUKeyCache struct {
	// Initialized on construction and never mutated, safe for concurrent access.
	// https://golang.org/pkg/hash/maphash/#Hash
	keySeed maphash.Seed

	mu     sync.Mutex
	cache  map[uint64]int64 // GUARDED_BY(mu)
	clock  int64            // GUARDED_BY(mu)
	maxLen int              // GUARDED_BY(mu)
}

// NewLRUKeyCache returns a new LRUKeyCache.
func NewLRUKeyCache(maxLen int) *LRUKeyCache {
	return &LRUKeyCache{
		keySeed: maphash.MakeSeed(),
		cache:   make(map[uint64]int64, maxLen),
		maxLen:  maxLen,
	}
}

// MakeKey makes a key from s.
// Insertions and lookups into c must use keys generated by c.
// https://golang.org/pkg/hash/maphash/#Hash.SetSeed
func (c *LRUKeyCache) MakeKey(s string) uint64 {
	var h maphash.Hash
	h.SetSeed(c.keySeed)
	_, _ = h.WriteString(s) // always succeeds
	return h.Sum64()
}

// Lookup returns true if the key exists in the cache, or false otherwise.
// The returned int is the number of keys in the cache, and is valid on true or false.
func (c *LRUKeyCache) Lookup(key uint64) (int, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.cache[key]; ok {
		c.clock++
		c.cache[key] = c.clock
		return len(c.cache), true
	}
	return len(c.cache), false
}

// Insert inserts key into the cache.
// If the cache is already at its max len, the least recently used key is evicted.
func (c *LRUKeyCache) Insert(key uint64) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.clock++
	c.cache[key] = c.clock
	if len(c.cache) > c.maxLen {
		// Find and evict the LRU item.  A scan over the cache is fine, since it's small.
		minClock, lruKey := int64(math.MaxInt64), uint64(0)
		for k, clock := range c.cache {
			if clock < minClock {
				minClock, lruKey = clock, k
			}
		}
		delete(c.cache, lruKey)
	}
	return len(c.cache)
}
