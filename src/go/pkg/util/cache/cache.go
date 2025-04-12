// Package cache provides general-purpose in-memory caches.
// Different caches provide different eviction policies suitable for
// specific use cases.
package cache

import "time"

// Stats returns usage statistics about an individual cache, useful to assess the
// efficiency of a cache.
//
// The values returned in this struct are approximations of the current state of the cache.
// For the sake of efficiency, certain edge cases in the implementation can lead to
// inaccuracies.
type Stats struct {
	Writes    uint64
	Hits      uint64
	Misses    uint64
	Evictions uint64
	Removals  uint64
}

// Cache defines the standard behavior of in-memory thread-safe caches.
//
// Different caches can have different eviction policies which determine
// when and how entries are automatically removed from the cache.
//
// Using a cache is very simple:
//
//	  c := NewLRU(5*time.Second,     // default per-entry ttl
//	              5*time.Second,     // eviction interval
//	              500)               // max # of entries tracked
//	  c.Set("foo", "bar")			// add an entry
//	  value, ok := c.Get("foo")		// try to retrieve the entry
//	  if ok {
//			fmt.Printf("Got value %v\n", value)
//	  } else {
//	     fmt.Printf("Value was not found, must have been evicted")
//	  }
type Cache interface {
	// Ideas for the future:
	//   - Return the number of entries in the cache in stats.
	//   - Provide an eviction callback to know when entries are evicted.
	//   - Have Set and Remove return the previous value for the key, if any.
	//   - Have Get return the expiration time for entries.

	// Set inserts an entry in the cache. This will replace any entry with
	// the same key that is already in the cache. The entry may be automatically
	// expunged from the cache at some point, depending on the eviction policies
	// of the cache and the options specified when the cache was created.
	Set(key any, value any)

	// Get retrieves the value associated with the supplied key if the key
	// is present in the cache.
	Get(key any) (value any, ok bool)

	// Remove synchronously deletes the given key from the cache. This has no effect if the key is not
	// currently in the cache.
	Remove(key any)

	// RemoveAll synchronously deletes all entries from the cache.
	RemoveAll()

	// Stats returns information about the efficiency of the cache.
	Stats() Stats
}

// ExpiringCache is a cache with entries that are evicted over time
type ExpiringCache interface {
	Cache

	// SetWithExpiration inserts an entry in the cache with a requested expiration time.
	// This will replace any entry with the same key that is already in the cache.
	// The entry will be automatically expunged from the cache at or slightly after the
	// requested expiration time.
	SetWithExpiration(key any, value any, expiration time.Duration)

	// EvictExpired() synchronously evicts all expired entries from the cache
	EvictExpired()
}
