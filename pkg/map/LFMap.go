package lfmap

import "sync/atomic"
import "unsafe"

import "github.com/sirgallo/ads/pkg/utils"


func NewLFMap[T comparable]() *LFMap[T] {
	bitChunkSize := 5
	
	rootNode := &LFMapNode[T]{
		BitMap: 0,
		Children: []*LFMapNode[T]{},
	}

	return &LFMap[T]{
		BitChunkSize: bitChunkSize,
		Root: unsafe.Pointer(rootNode),
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

func CopyNode[T comparable](node *LFMapNode[T]) *LFMapNode[T] {
	nodeCopy := &LFMapNode[T]{
		Key: node.Key,
		Value: node.Value,
		IsLeafNode: node.IsLeafNode,
		BitMap: node.BitMap,
		Children: make([]*LFMapNode[T], len(node.Children)),
	}

	copy(nodeCopy.Children, node.Children)

	return nodeCopy
}

func (lfMap *LFMap[T]) Insert(key string, value T) bool {
	for {
		completed := lfMap.insertRecursive(&lfMap.Root, key, value, 0)
		if completed { return true }
	}
}

func (lfMap *LFMap[T]) insertRecursive(node *unsafe.Pointer, key string, value T, level int) bool {
	hash := utils.FnvHash(key)
	index := lfMap.getSparseIndex(hash, level)
	
	currNode := (*LFMapNode[T])(atomic.LoadPointer(node))
	nodeCopy := CopyNode[T](currNode)

	if ! IsBitSet(nodeCopy.BitMap, index) {
		newLeaf := NewLeafNode(key, value)
		nodeCopy.BitMap = SetBit(nodeCopy.BitMap, index)
		pos := lfMap.getPosition(nodeCopy.BitMap, hash, level)
		nodeCopy.Children = ExtendTable(nodeCopy.Children, nodeCopy.BitMap, pos, newLeaf)
		
		if atomic.CompareAndSwapPointer(node, unsafe.Pointer(currNode), unsafe.Pointer(nodeCopy)) {
			return true
		} else { return false }
	} else {
		pos := lfMap.getPosition(nodeCopy.BitMap, hash, level)
		childNode := nodeCopy.Children[pos]

		if childNode.IsLeafNode {
			if key == childNode.Key {
				nodeCopy.Children[pos].Value = value

				if atomic.CompareAndSwapPointer(node, unsafe.Pointer(currNode), unsafe.Pointer(nodeCopy)) {
					return true
				} else { return false }
			} else {
				newINode := NewInternalNode[T]()
				ncPtr := unsafe.Pointer(newINode)
				
				lfMap.insertRecursive(&ncPtr, childNode.Key, childNode.Value, level + 1)
				lfMap.insertRecursive(&ncPtr, key, value, level + 1)

				nodeCopy.Children[pos] = (*LFMapNode[T])(atomic.LoadPointer(&ncPtr))

				if atomic.CompareAndSwapPointer(node, unsafe.Pointer(currNode), unsafe.Pointer(nodeCopy)) {
					return true
				} else { return false }
			}
		} else {			
			childPtr := unsafe.Pointer(nodeCopy.Children[pos])
			lfMap.insertRecursive(&childPtr, key, value, level + 1) 

			nodeCopy.Children[pos] = (*LFMapNode[T])(atomic.LoadPointer(&childPtr))
			
			if atomic.CompareAndSwapPointer(node, unsafe.Pointer(currNode), unsafe.Pointer(nodeCopy)) {
				return true
			} else { return false }
		}
	}
}

func (lfMap *LFMap[T]) Retrieve(key string) T {
	hash := utils.FnvHash(key)
	return lfMap.retrieveRecursive(&lfMap.Root, key, hash, 0)
}

func (lfMap *LFMap[T]) retrieveRecursive(node *unsafe.Pointer, key string, hash uint32, level int) T {
	index := lfMap.getSparseIndex(hash, level)
	currNode := (*LFMapNode[T])(atomic.LoadPointer(node))
	
	if ! IsBitSet(currNode.BitMap, index) { 
		return utils.GetZero[T]() 
	} else {
		pos := lfMap.getPosition(currNode.BitMap, hash, level)
		childNode := currNode.Children[pos]

		if childNode.IsLeafNode && key == childNode.Key {
			if childNode.Value == (*LFMapNode[T])(atomic.LoadPointer(node)).Children[pos].Value {
				return childNode.Value
			} else { return utils.GetZero[T]() }
 		} else { 
			childPtr := unsafe.Pointer(currNode.Children[pos])
			return lfMap.retrieveRecursive(&childPtr, key, hash, level + 1) 
		}
	}
}

func (lfMap *LFMap[T]) Delete(key string) bool {
	hash := utils.FnvHash(key)
	for {
		completed := lfMap.deleteRecursive(&lfMap.Root, key, hash, 0)
		if completed { return true }
	}
}

func (lfMap *LFMap[T]) deleteRecursive(node *unsafe.Pointer, key string, hash uint32, level int) bool {
	index := lfMap.getSparseIndex(hash, level)
	currNode := (*LFMapNode[T])(atomic.LoadPointer(node))
	nodeCopy := CopyNode[T](currNode)

	if ! IsBitSet(nodeCopy.BitMap, index) { 
		return true 
	} else {
		pos := lfMap.getPosition(nodeCopy.BitMap, hash, level)
		childNode := nodeCopy.Children[pos]
		
		if childNode.IsLeafNode {
			if key == childNode.Key {
				nodeCopy.BitMap = SetBit(nodeCopy.BitMap, index)
				nodeCopy.Children = ShrinkTable(nodeCopy.Children, nodeCopy.BitMap, pos)
				
				if atomic.CompareAndSwapPointer(node, unsafe.Pointer(currNode), unsafe.Pointer(nodeCopy)) {
					return true
				} else { return false }
			}
			
			return false
		} else { 
			childPtr := unsafe.Pointer(nodeCopy.Children[pos])
			
			lfMap.deleteRecursive(&childPtr, key, hash, level + 1)

			popCount := calculateHammingWeight(nodeCopy.BitMap)
			if popCount == 0 { // if empty internal node, remove from the mapped array
				nodeCopy.BitMap = SetBit(nodeCopy.BitMap, index)
				nodeCopy.Children = ShrinkTable(nodeCopy.Children, nodeCopy.BitMap, pos)
			} 

			if atomic.CompareAndSwapPointer(node, unsafe.Pointer(currNode), unsafe.Pointer(nodeCopy)) {
				return true
			} else { return false }
		}
	}
}