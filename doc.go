// Package memcache provides a generic in-memory key-value cache.
package memcache

// TODO:
// - Remove closer and use regular way to do it so it's easier to read and
//   understand the code.
// - Separate concerns of cache wrapper and active & passive expiration.
// - Add more eviction policies.
