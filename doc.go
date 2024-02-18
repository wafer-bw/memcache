// Package memcache provides a generic in-memory key-value cache.
//
// The capacity of a cache is the total number of keys it is allowed to hold.
package memcache

// TODO:
// - Consider relocating errors back into package memcache.
// - Consider renaming package to "inmem" or something like that to avoid
//   stutters.
// - Use interface for active expiration like used with closer?
