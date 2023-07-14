package lfmap

import "math"
import "sync/atomic"
import "unsafe"

import "github.com/sirgallo/ads/pkg/utils"


func NewLFMap[T comparable, V uint32 | uint64](opts LFMapOpts) *LFMap[T, V] {
	var v V
	var bitChunkSize int
	
	switch any(v).(type) {
		case uint32: 
			bitChunkSize = 5
		case uint64:
			bitChunkSize = 6
	}

	hashChunks := int(math.Pow(float64(2), float64(bitChunkSize))) / bitChunkSize
	
	nodePool := NewLFMapNodePool[T, V](opts.PoolSize)
	rootNode := nodePool.GetLFMapNode()
	
	rootNode.IsLeafNode = false
	rootNode.BitMap = 0
	rootNode.Children = []*LFMapNode[T, V]{}

	return &LFMap[T, V]{
		BitChunkSize: bitChunkSize,
		HashChunks: hashChunks,
		Root: unsafe.Pointer(rootNode),
		NodePool: nodePool,
	}
}

func (lfMap *LFMap[T, V]) NewLeafNode(key string, value T) *LFMapNode[T, V] {
	leafNode := lfMap.NodePool.GetLFMapNode()
	
	leafNode.Key = key
	leafNode.Value = value
	leafNode.IsLeafNode = true

	return leafNode
}

func (lfMap *LFMap[T, V]) NewInternalNode() *LFMapNode[T, V] {
	iNode := lfMap.NodePool.GetLFMapNode()

	iNode.IsLeafNode = false
	iNode.BitMap = 0
	iNode.Children = []*LFMapNode[T, V]{}

	return iNode
}

func (lfMap *LFMap[T, V]) CopyNode(node *LFMapNode[T, V]) *LFMapNode[T, V] {
	nodeCopy := lfMap.NodePool.GetLFMapNode()
	
	nodeCopy.Key = node.Key
	nodeCopy.Value = node.Value
	nodeCopy.IsLeafNode = node.IsLeafNode
	nodeCopy.BitMap = node.BitMap
	nodeCopy.Children = make([]*LFMapNode[T, V], len(node.Children))

	copy(nodeCopy.Children, node.Children)

	return nodeCopy
}

func (lfMap *LFMap[T, V]) Insert(key string, value T) bool {
	for {
		completed := lfMap.insertRecursive(&lfMap.Root, key, value, 0)
		if completed { return true }
	}
}

func (lfMap *LFMap[T, V]) insertRecursive(node *unsafe.Pointer, key string, value T, level int) bool {
	hash := lfMap.CalculateHashForCurrentLevel(key, level)
	index := lfMap.getSparseIndex(hash, level)
	
	currNode := (*LFMapNode[T, V])(atomic.LoadPointer(node))
	nodeCopy := lfMap.CopyNode(currNode)

	if ! IsBitSet(nodeCopy.BitMap, index) {
		newLeaf := lfMap.NewLeafNode(key, value)
		nodeCopy.BitMap = SetBit(nodeCopy.BitMap, index)
		pos := lfMap.getPosition(nodeCopy.BitMap, hash, level)
		nodeCopy.Children = ExtendTable(nodeCopy.Children, nodeCopy.BitMap, pos, newLeaf)
		
		return lfMap.compareAndSwap(node, currNode, nodeCopy)
	} else {
		pos := lfMap.getPosition(nodeCopy.BitMap, hash, level)
		childNode := nodeCopy.Children[pos]

		if childNode.IsLeafNode {
			if key == childNode.Key {
				nodeCopy.Children[pos].Value = value
				return lfMap.compareAndSwap(node, currNode, nodeCopy)
			} else {
				newINode := lfMap.NewInternalNode()
				iNodePtr := unsafe.Pointer(newINode)
				
				lfMap.insertRecursive(&iNodePtr, childNode.Key, childNode.Value, level + 1)
				lfMap.insertRecursive(&iNodePtr, key, value, level + 1)

				nodeCopy.Children[pos] = (*LFMapNode[T, V])(atomic.LoadPointer(&iNodePtr))
				return lfMap.compareAndSwap(node, currNode, nodeCopy)
			}
		} else {			
			childPtr := unsafe.Pointer(nodeCopy.Children[pos])
			lfMap.insertRecursive(&childPtr, key, value, level + 1) 

			nodeCopy.Children[pos] = (*LFMapNode[T, V])(atomic.LoadPointer(&childPtr))
			return lfMap.compareAndSwap(node, currNode, nodeCopy)
		}
	}
}

func (lfMap *LFMap[T, V]) Retrieve(key string) T {
	return lfMap.retrieveRecursive(&lfMap.Root, key, 0)
}

func (lfMap *LFMap[T, V]) retrieveRecursive(node *unsafe.Pointer, key string, level int) T {
	hash := lfMap.CalculateHashForCurrentLevel(key, level)
	index := lfMap.getSparseIndex(hash, level)
	currNode := (*LFMapNode[T, V])(atomic.LoadPointer(node))
	
	if ! IsBitSet(currNode.BitMap, index) { 
		return utils.GetZero[T]() 
	} else {
		pos := lfMap.getPosition(currNode.BitMap, hash, level)
		childNode := currNode.Children[pos]

		if childNode.IsLeafNode && key == childNode.Key {
			if childNode.Value == (*LFMapNode[T, V])(atomic.LoadPointer(node)).Children[pos].Value {
				return childNode.Value
			} else { return utils.GetZero[T]() }
 		} else { 
			childPtr := unsafe.Pointer(currNode.Children[pos])
			return lfMap.retrieveRecursive(&childPtr, key, level + 1) 
		}
	}
}

func (lfMap *LFMap[T, V]) Delete(key string) bool {
	for {
		completed := lfMap.deleteRecursive(&lfMap.Root, key, 0)
		if completed { return true }
	}
}

func (lfMap *LFMap[T, V]) deleteRecursive(node *unsafe.Pointer, key string, level int) bool {
	hash := lfMap.CalculateHashForCurrentLevel(key, level)
	index := lfMap.getSparseIndex(hash, level)
	
	currNode := (*LFMapNode[T, V])(atomic.LoadPointer(node))
	nodeCopy := lfMap.CopyNode(currNode)

	if ! IsBitSet(nodeCopy.BitMap, index) { 
		return true 
	} else {
		pos := lfMap.getPosition(nodeCopy.BitMap, hash, level)
		childNode := nodeCopy.Children[pos]
		
		if childNode.IsLeafNode {
			if key == childNode.Key {
				nodeCopy.BitMap = SetBit(nodeCopy.BitMap, index)
				nodeCopy.Children = ShrinkTable(nodeCopy.Children, nodeCopy.BitMap, pos)
				
				return lfMap.compareAndSwap(node, currNode, nodeCopy)
			}
			
			return false
		} else { 
			childPtr := unsafe.Pointer(nodeCopy.Children[pos])
			lfMap.deleteRecursive(&childPtr, key, level + 1)

			popCount := calculateHammingWeight(nodeCopy.BitMap)
			if popCount == 0 { // if empty internal node, remove from the mapped array
				nodeCopy.BitMap = SetBit(nodeCopy.BitMap, index)
				nodeCopy.Children = ShrinkTable(nodeCopy.Children, nodeCopy.BitMap, pos)
			} 

			return lfMap.compareAndSwap(node, currNode, nodeCopy)
		}
	}
}

func (lfMap *LFMap[T, V]) compareAndSwap(node *unsafe.Pointer, currNode *LFMapNode[T, V], nodeCopy *LFMapNode[T, V]) bool {
	if atomic.CompareAndSwapPointer(node, unsafe.Pointer(currNode), unsafe.Pointer(nodeCopy)) {
		if currNode.IsLeafNode { lfMap.NodePool.PutLFMapNode(currNode) }
		return true
	} else { 
		if nodeCopy.IsLeafNode { lfMap.NodePool.PutLFMapNode(nodeCopy) }
		return false 
	}
}