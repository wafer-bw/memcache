// Package memcache provides a generic in-memory key-value cache.
package memcache

// TODO:
// - Come up with a nicer way of handling capacity constraints across the
//   different stores.
// - Verify correctness of each store so far.
// - Use interface for active expiration.
// - Add more eviction policies.
