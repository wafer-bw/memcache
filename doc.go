// Package memcache provides a generic in-memory key-value cache.
package memcache

// TODO:
// - Determine if cache & closer need to be constructed as pointers or not.
// - Separate concerns of store vs top level cache. Is the cache just a wrapper
//   or does it have its own behavior to add?
