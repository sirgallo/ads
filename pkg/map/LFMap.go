package lfmap

import "sync/atomic"
import "unsafe"

import "github.com/sirgallo/ads/pkg/utils"


func NewLFMap[T comparable](opts LFMapOpts) *LFMap[T] {
	bitChunkSize := 5
	
	nodePool := NewLFMapNodePool[T](opts.PoolSize)
	rootNode := nodePool.GetLFMapNode()
	
	rootNode.IsLeafNode = false
	rootNode.BitMap = 0
	rootNode.Children = []*LFMapNode[T]{}

	return &LFMap[T]{
		BitChunkSize: bitChunkSize,
		Root: unsafe.Pointer(rootNode),
		NodePool: nodePool,
	}
}

func (lfMap *LFMap[T]) NewLeafNode(key string, value T) *LFMapNode[T] {
	leafNode := lfMap.NodePool.GetLFMapNode()
	
	leafNode.Key = key
	leafNode.Value = value
	leafNode.IsLeafNode = true

	return leafNode
}

func (lfMap *LFMap[T]) NewInternalNode() *LFMapNode[T] {
	iNode := lfMap.NodePool.GetLFMapNode()

	iNode.IsLeafNode = false
	iNode.BitMap = 0
	iNode.Children = []*LFMapNode[T]{}

	return iNode
}

func (lfMap *LFMap[T]) CopyNode(node *LFMapNode[T]) *LFMapNode[T] {
	nodeCopy:= lfMap.NodePool.GetLFMapNode()
	
	nodeCopy.Key = node.Key
	nodeCopy.Value = node.Value
	nodeCopy.IsLeafNode = node.IsLeafNode
	nodeCopy.BitMap = node.BitMap
	nodeCopy.Children = make([]*LFMapNode[T], len(node.Children))

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

				nodeCopy.Children[pos] = (*LFMapNode[T])(atomic.LoadPointer(&iNodePtr))
				return lfMap.compareAndSwap(node, currNode, nodeCopy)
			}
		} else {			
			childPtr := unsafe.Pointer(nodeCopy.Children[pos])
			lfMap.insertRecursive(&childPtr, key, value, level + 1) 

			nodeCopy.Children[pos] = (*LFMapNode[T])(atomic.LoadPointer(&childPtr))
			return lfMap.compareAndSwap(node, currNode, nodeCopy)
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
			lfMap.deleteRecursive(&childPtr, key, hash, level + 1)

			popCount := calculateHammingWeight(nodeCopy.BitMap)
			if popCount == 0 { // if empty internal node, remove from the mapped array
				nodeCopy.BitMap = SetBit(nodeCopy.BitMap, index)
				nodeCopy.Children = ShrinkTable(nodeCopy.Children, nodeCopy.BitMap, pos)
			} 

			return lfMap.compareAndSwap(node, currNode, nodeCopy)
		}
	}
}

func (lfMap *LFMap[T]) compareAndSwap(node *unsafe.Pointer, currNode *LFMapNode[T], nodeCopy *LFMapNode[T]) bool {
	if atomic.CompareAndSwapPointer(node, unsafe.Pointer(currNode), unsafe.Pointer(nodeCopy)) {
		if currNode.IsLeafNode {
			lfMap.NodePool.PutLFMapNode(currNode)
		}

		return true
	} else { 
		if nodeCopy.IsLeafNode {
			lfMap.NodePool.PutLFMapNode(nodeCopy)
		}

		return false 
	}
}