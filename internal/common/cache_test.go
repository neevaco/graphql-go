package common_test

import (
	"testing"

	"github.com/graph-gophers/graphql-go/internal/common"
)

func TestLRUKeyCache(t *testing.T) {
	var (
		cache = common.NewLRUKeyCache(3)
		key1  = cache.MakeKey("key1")
		key2  = cache.MakeKey("key2")
		key3  = cache.MakeKey("key3")
		key4  = cache.MakeKey("key4")
	)

	// First test the empty cache.
	for i, key := range []uint64{key1, key2, key3, key4} {
		if len, ok := cache.Lookup(key); len != 0 || ok {
			t.Errorf("empty cache key%d lookup got (%v, %v)", i, len, ok)
		}
	}

	// Insert and lookup a key.
	if got, want := cache.Insert(key1), 1; got != want {
		t.Errorf("key1 insert got len %v, want %v", got, want)
	}
	if len, ok := cache.Lookup(key1); len != 1 || !ok {
		t.Errorf("key1 lookup got (%v, %v)", len, ok)
	}

	// Insert and lookup more keys.
	if got, want := cache.Insert(key2), 2; got != want {
		t.Errorf("key2 insert got len %v, want %v", got, want)
	}
	if len, ok := cache.Lookup(key2); len != 2 || !ok {
		t.Errorf("key2 lookup got (%v, %v)", len, ok)
	}
	if got, want := cache.Insert(key3), 3; got != want {
		t.Errorf("key3 insert got len %v, want %v", got, want)
	}
	if len, ok := cache.Lookup(key3); len != 3 || !ok {
		t.Errorf("key3 lookup got (%v, %v)", len, ok)
	}

	// Lookup key1 again, so that key2 becomes LRU.
	if len, ok := cache.Lookup(key1); len != 3 || !ok {
		t.Errorf("key1 lookup got (%v, %v)", len, ok)
	}

	// Inserting key4 causes an eviction.
	if got, want := cache.Insert(key4), 3; got != want {
		t.Errorf("key4 insert got len %v, want %v", got, want)
	}

	// Make the keys again, and make sure key2 was evicted, and all other keys are still there.
	// Making the keys again is a regression test for a tricky bug where maphash is auto-initialized
	// with a random seed.
	key1 = cache.MakeKey("key1")
	key2 = cache.MakeKey("key2")
	key3 = cache.MakeKey("key3")
	key4 = cache.MakeKey("key4")
	for i, key := range []uint64{key1, key2, key3, key4} {
		wantLen, wantOK := 3, key != key2
		if len, ok := cache.Lookup(key); len != wantLen || ok != wantOK {
			t.Errorf("full cache key%d lookup got (%v, %v) want (%v, %v)", i, len, ok, wantLen, wantOK)
		}
	}
}
