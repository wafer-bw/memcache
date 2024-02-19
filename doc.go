// Package memcache provides a generic in-memory key-value cache.
//
// The capacity of a cache is the total number of keys it is allowed to hold.
package memcache

// TODO:
// - Consider renaming package to "inmem" to avoid stutters.
// - Update WithActiveExpiration to use an interface and make expirers public.
