# Storage

## Cache

By design `terraformctl` will use some sort of caching layer as it's source of truth for operations. 
This layer is a stateless part of the system that works as a key/value storage cache.
The layer is abstracted out behind an interface so any component (`etcd`, `memcache`, `redis`, etc..) should be able to implement this interface.

## Persist

All data will ultimately persist to a more resilient storage location. 
This layer is also abstracted our behind an interface so any storage component (`Cosmos DB`, `Unix filesystem`, etc) should be able to implement this interface. 
This layer is the stateful part of the system, and will attempt to persist all data in the **cache** layer to some storage implementation.
