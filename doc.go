// Package memcache provides a generic in-memory key-value cache.
//
// The capacity of a cache is the total number of keys it is allowed to hold.
package memcache

// TODO:
// - Rename package to "inmem" or something like that to avoid stutters.
// - Verify correctness of each store so far.
// - Use interface for active expiration.
// - Add more eviction policies.
