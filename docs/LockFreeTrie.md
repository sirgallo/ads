# Lock Free Trie


## CTrie

CTries are a non-blocking concurrent `Hash Array Mapped Trie (HAMT)` that utilizes atomic Compare and Swap (CAS) operations.

It supports `insert`, `lookup`, and `remove`.


## HAMT Algorithm

### Node Structure 

```go

type HAMTNode struct {
  key string
  value interface{}
  children [2 ^ C]HAMTNode
}
```

### Insert

1. Compute a hash of the incoming key.

```
The key will be hashed using SHA1, and will return the [20]byte representation of the hash.
```

2. Convert the [20]byte representation to byte array representation 

3. If current trie node is empty, we can insert the key and value here

4. If current trie node is not empty, but incoming key == existing key, overwrite the current value with the incoming value

5. 