// Package memcache provides a generic in-memory key-value cache.
//
// The capacity of a cache is the total number of keys it is allowed to hold.
package memcache

// TODO:
// - Embed dupe eviction policy store methods using common implementation.
// - Add tests covering remaining eviction policy store methods.
// - Consider renaming package to "inmem" to avoid stutters.
// - Update WithActiveExpiration to use an interface and make expirers public.
