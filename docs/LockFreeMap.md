# Lock Free Map


## Data Structure 

### CTrie

CTries are a non-blocking concurrent `Hash Array Mapped Tries (HAMT)` that utilize atomic Compare and Swap (CAS) operations.

To learn more about the Hash Array Mapped Trie algorithm, check out (hamt)(https://github.com/sirgallo/hamt/blob/main/docs/HashArrayMappedTrie.md).

It supports `insert`, `retrieve`, and `delete` operations.


## Design

The design takes the basic algorithm for `HAMT`, and adds in `CAS` to insert/delete new values. A thread will modify an element at the point in time it loads it, and if the compare and swap operation fails, the update is discarded and it will start back at the root and traverse the path through the trie and reattempt to add/delete the element.