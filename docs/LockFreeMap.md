# Lock Free Map


## Data Structure 

### CTrie

A `Concurrent Trie` is a non-blocking implementation of a `Hash Array Mapped Trie (HAMT)` that utilizes atomic `Compare-and-Swap (CAS)` operations.

To learn more about the `Hash Array Mapped Trie` algorithm, check out [hamt](https://github.com/sirgallo/hamt/blob/main/docs/HashArrayMappedTrie.md).


## Design

The design takes the basic algorithm for `HAMT`, and adds in `CAS` to insert/delete new values. A thread will modify an element at the point in time it loads it, and if the compare and swap operation fails, the update is discarded and the operation will start back at the root of the trie and traverse the path through to reattempt to add/delete the element.


### Path Copying

This CTrie implements full path copying. As an operation traverses down the path to the key, on inserts/deletes it will make a copy of the current node and modify the copy instead of modifying the node in place. This makes the CTrie [persistant](https://en.wikipedia.org/wiki/Persistent_data_structure). The modified node causes all parent nodes to point to it by cascading the changes up the path back to the root of the trie. This is done by passing a copy of the node being looked at, and then performing compare and swap back up the path. If the compare and swap operation fails, the copy is discarded and the operation retries back at the root.


### Object Pool

This Ctrie has a hybrid approach to cleaning up nodes, where it utilizes both `Go's` garbage collection as well as an Object Pool. When copies of nodes are created, the compare and swap operation will recycle the failed node if it is a leaf nodes So, if the current node is replaced by the new copy, and the current node is a leaf node, it is recycled. Otherwise, the failed replacement of the new copy is recycled. On inserts, if there are available objects in the pool, a new node can be pulled from the pool. This ensures that memory is not being allocated/deallocated all the time, and should have an overall positive effect on performance of the trie. Internal nodes are destroyed by the garbage collection.


### Hash Exhaustion

Since the 32 bit hash only has 6 chunks of 5 bits, the Ctrie is capped at 6 levels (or around 1000000000 key val pairs), which is not optimal for a trie data strucutre. To circumvent this, we can re-seed our hash after every 6 levels, using [Murmur32](../pkg/map/Murmur32.go) as our hash function. To achieve this, we utilize the following functions:

```go
func (lfMap *LFMap[T]) CalculateHashForCurrentLevel(key string, level int) uint32 {
	currChunk := level / lfMap.TotalLevels
	seed := uint32(currChunk + 1)
	return Murmur32(key, seed)
}
```

```go
func GetIndexForLevel(hash uint32, chunkSize int, level int, totalLevels int) int {
	updatedLevel := level % totalLevels
	return GetIndex(hash, chunkSize, updatedLevel)
}

func GetIndex(hash uint32, chunkSize int, level int) int {
	slots := int(math.Pow(float64(2), float64(chunkSize)))
	mask := uint32(slots - 1)
	shiftSize := slots - (chunkSize * (level + 1))

	return int(hash >> shiftSize & mask)
}
```

this ensures we take steps of 6 levels, and at the start of the next 6 levels, re-seed the hash and start from the beginning of the new hash value for the key. Now we are no longer limited to just 6 levels. 

The seed value is just the `uint32` representation of the current chunk of levels + 1.


## Sources

[LockFreeMap](../pkg/map/LFMap.go)