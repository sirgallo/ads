# Lock Free Map


## Data Structure 

### CTrie

A `Concurrent Trie` is a non-blocking implementation of a `Hash Array Mapped Trie (HAMT)` that utilizes atomic Compare and Swap (CAS) operations.

To learn more about the `Hash Array Mapped Trie` algorithm, check out [hamt](https://github.com/sirgallo/hamt/blob/main/docs/HashArrayMappedTrie.md).


## Design

The design takes the basic algorithm for `HAMT`, and adds in `CAS` to insert/delete new values. A thread will modify an element at the point in time it loads it, and if the compare and swap operation fails, the update is discarded and the operation will start back at the root of the trie and traverse the path through to reattempt to add/delete the element.