package lfmap

import "sync/atomic"

import "github.com/sirgallo/ads/pkg/utils"


func NewLFMap[T comparable]() *LFMap[T] {
	bitChunkSize := 5
	
	root := &atomic.Value{}
	rootNode := &LFMapNode[T]{
		BitMap: 0,
		Children: []*LFMapNode[T]{},
	}

	root.Store(rootNode)

	return &LFMap[T]{
		BitChunkSize: bitChunkSize,
		Root: root,
	}
}

func NewLeafNode[T comparable](key string, value T) *LFMapNode[T] {
	return &LFMapNode[T]{ 
		Key: key, 
		Value: value, 
		IsLeafNode: true,
	}
}

func NewInternalNode[T comparable]() *LFMapNode[T] {
	return &LFMapNode[T]{
		IsLeafNode: false, 
		BitMap: 0,
		Children: []*LFMapNode[T]{},
	}
}

func (lfMap *LFMap[T]) Insert(key string, value T) bool {
	for {
		completed := lfMap.insertRecursive(lfMap.Root, key, value, 0)
		if completed {
			return true
		}
	}
}

func (lfMap *LFMap[T]) insertRecursive(node *atomic.Value, key string, value T, level int) bool {
	hash := utils.FnvHash(key)
	index := lfMap.getSparseIndex(hash, level)
	currNode := node.Load().(*LFMapNode[T])

	if ! IsBitSet(currNode.BitMap, index) {
		newLeaf := NewLeafNode(key, value)
		currNode.BitMap = SetBit(currNode.BitMap, index)
		pos := lfMap.getPosition(currNode.BitMap, hash, level)
		currNode.Children = ExtendTable(currNode.Children, currNode.BitMap, pos, newLeaf)
		
		if node.CompareAndSwap(node.Load().(*LFMapNode[T]), currNode) {
			return true
		} else { return false }
	} else {
		pos := lfMap.getPosition(currNode.BitMap, hash, level)
		childNode := currNode.Children[pos]

		if childNode.IsLeafNode {
			if key == childNode.Key {
				currNode.Children[pos].Value = value

				if node.CompareAndSwap(node.Load().(*LFMapNode[T]), currNode) {
					return true
				} else { return false }
			} else {
				newInternalNode := NewInternalNode[T]()
				currNode.Children[pos] = newInternalNode

				atomicINode := &atomic.Value{}
				atomicINode.Store(currNode.Children[pos])

				completed := lfMap.insertRecursive(atomicINode, childNode.Key, childNode.Value, level + 1)
				if ! completed {
					return false
				}
				
				return lfMap.insertRecursive(atomicINode, key, value, level + 1)
			}
		} else {
			atomicChild := &atomic.Value{}
			atomicChild.Store(currNode.Children[pos])
			
			return lfMap.insertRecursive(atomicChild, key, value, level + 1) 
		}
	}
}

func (lfMap *LFMap[T]) Retrieve(key string) T {
	hash := utils.FnvHash(key)
	return lfMap.retrieveRecursive(lfMap.Root, key, hash, 0)
}

func (lfMap *LFMap[T]) retrieveRecursive(node *atomic.Value, key string, hash uint32, level int) T {
	index := lfMap.getSparseIndex(hash, level)
	currNode := node.Load().(*LFMapNode[T])
	
	if ! IsBitSet(currNode.BitMap, index) { 
		return utils.GetZero[T]() 
	} else {
		pos := lfMap.getPosition(currNode.BitMap, hash, level)
		childNode := currNode.Children[pos]

		if childNode.IsLeafNode && key == childNode.Key {
			if childNode.Value == (node.Load().(*LFMapNode[T])).Children[pos].Value {
				return childNode.Value
			} else { return utils.GetZero[T]() }
 		} else { 
			atomicChild := &atomic.Value{}
			atomicChild.Store(childNode)

			return lfMap.retrieveRecursive(atomicChild, key, hash, level + 1) 
		}
	}
}

func (lfMap *LFMap[T]) Delete(key string) bool {
	hash := utils.FnvHash(key)
	for {
		completed := lfMap.deleteRecursive(lfMap.Root, key, hash, 0)
		if completed {
			return true
		}
	}
}

func (lfMap *LFMap[T]) deleteRecursive(node *atomic.Value, key string, hash uint32, level int) bool {
	index := lfMap.getSparseIndex(hash, level)
	currNode := node.Load().(*LFMapNode[T])

	if ! IsBitSet(currNode.BitMap, index) { 
		return true 
	} else {
		pos := lfMap.getPosition(currNode.BitMap, hash, level)
		childNode := currNode.Children[pos]
		
		if childNode.IsLeafNode {
			if key == childNode.Key {
				currNode.BitMap = SetBit(currNode.BitMap, index)
				currNode.Children = ShrinkTable(currNode.Children, currNode.BitMap, pos)
				
				if node.CompareAndSwap(node.Load().(*LFMapNode[T]), currNode) {
					return true
				} else { return false }
			}
			
			return false
		} else { 
			atomicChild := &atomic.Value{}
			atomicChild.Store(childNode)
			
			completed := lfMap.deleteRecursive(atomicChild, key, hash, level + 1) 
			if ! completed {
				return false
			}

			popCount := calculateHammingWeight(currNode.BitMap)
			if popCount == 0 { // if empty internal node, remove from the mapped array
				currNode.BitMap = SetBit(currNode.BitMap, index)
				currNode.Children = ShrinkTable(currNode.Children, currNode.BitMap, pos)
			} 

			if node.CompareAndSwap(node.Load().(*LFMapNode[T]), currNode) {
				return true
			} else { return false }
		}
	}
}