// Package memcache provides a generic in-memory key-value cache.
//
// The capacity of a cache is the total number of keys it is allowed to hold.
package memcache

// TODO:
// - Create dedicated underlying data structures for caches depending on their
//   needs? This can likely replace full storer implementations.
// - Consider renaming package to "inmem" or something like that to avoid
//   stutters.
