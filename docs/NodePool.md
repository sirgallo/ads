# Node Pool

The node pool utilizes a Go `object pool` to create a channel that contains a buffer of the max defined size where nodes can be recycled. This helps to limit the number of objects in circulation and reduces the load on the Go garbage collector, making the queue more memory efficient. This should also optimize performance by reducing the frequency of memory allocations and deallocations.


## Source

[LockFreeNodePool](../src/node/LFNodePool.go)