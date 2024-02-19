# redesign notes

## Removal types
- [ ] Active Expiration: Delete keys that are expired periodically.
    - Information needed:
        - Item.Expire
        - Item.TTL
        - Item.OnExpire
    - Implementation:
        - Logic lives in custom method & called in base cache.
- [x] Passive Expiration: Delete keys that are expired when they are read.
    - Information needed:
        - Item.Expire
        - Item.TTL
        - Item.OnExpire
    - Implementation:
        - Logic lives in base cache & called in base cache.
- [x] Eviction: Delete key based on policy when adding a new key would breach cap.
    - Implementation:
        - Information needed:
            - Item.Expire (volatile-*)
            - Item.TTL (volatile-ttl)
            - Item.Frequency (*-lfu)
            - Item.OnEvict (*)
        - Datastructure.Evict
            - Called from base cache.
        - Logic lives in datastructure.

## Store Duplicates
These methods are the same across stores:
- Get
- Items
- Len

Thesse methods are the same across stores but may not be later:
- Keys

These methods could technically be the same if they were passed a function to
call:
- Flush
- Delete
