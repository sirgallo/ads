package lfmap

import "sync/atomic"

import "github.com/sirgallo/ads/pkg/utils"


func NewLFMap() *LFMap {
	bitChunkSize := 5
	
	root := &atomic.Value{}
	rootNode := &LFMapNode{
		BitMap: 0,
		Children: []*LFMapNode{},
	}

	root.Store(rootNode)

	return &LFMap{
		BitChunkSize: bitChunkSize,
		Root: root,
	}
}

func NewLeafNode(key string, value interface{}) *LFMapNode {
	return &LFMapNode{ 
		Key: key, 
		Value: value, 
		IsLeafNode: true,
	}
}

func NewInternalNode() *LFMapNode {
	return &LFMapNode{
		IsLeafNode: false, 
		BitMap: 0,
		Children: []*LFMapNode{},
	}
}

func (lfMap *LFMap) Insert(key string, value interface{}) bool {
	for {
		completed := lfMap.insertRecursive(lfMap.Root, key, value, 0)
		if completed {
			return true
		}
	}
}

func (lfMap *LFMap) insertRecursive(node *atomic.Value, key string, value interface{}, level int) bool {
	hash := utils.FnvHash(key)
	index := lfMap.getSparseIndex(hash, level)
	currNode := node.Load().(*LFMapNode)

	if ! isBitSet(currNode.BitMap, index) {
		newLeaf := NewLeafNode(key, value)
		currNode.BitMap = setBit(currNode.BitMap, index)
		pos := lfMap.getPosition(currNode.BitMap, hash, level)
		currNode.Children = ExtendTable(currNode.Children, currNode.BitMap, pos, newLeaf)
		
		if node.CompareAndSwap(node.Load().(*LFMapNode), currNode) {
			return true
		} else { return false }
	} else {
		pos := lfMap.getPosition(currNode.BitMap, hash, level)
		childNode := currNode.Children[pos]

		if childNode.IsLeafNode {
			if key == childNode.Key {
				currNode.Children[pos].Value = value

				if node.CompareAndSwap(node.Load().(*LFMapNode), currNode) {
					return true
				} else { return false }
			} else {
				newInternalNode := NewInternalNode()
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

func (lfMap *LFMap) Retrieve(key string) interface{} {
	hash := utils.FnvHash(key)
	return lfMap.retrieveRecursive(lfMap.Root, key, hash, 0)
}

func (lfMap *LFMap) retrieveRecursive(node *atomic.Value, key string, hash uint32, level int) interface{} {
	index := lfMap.getSparseIndex(hash, level)
	currNode := node.Load().(*LFMapNode)
	
	if ! isBitSet(currNode.BitMap, index) { 
		return nil 
	} else {
		pos := lfMap.getPosition(currNode.BitMap, hash, level)
		childNode := currNode.Children[pos]

		if childNode.IsLeafNode && key == childNode.Key {
			if childNode.Value == (node.Load().(*LFMapNode)).Children[pos].Value {
				return childNode.Value
			} else { return nil }
 		} else { 
			atomicChild := atomic.Value{}
			atomicChild.Store(childNode)

			return lfMap.retrieveRecursive(&atomicChild, key, hash, level + 1) 
		}
	}
}

func (lfMap *LFMap) Delete(key string) bool {
	hash := utils.FnvHash(key)
	for {
		completed := lfMap.deleteRecursive(lfMap.Root, key, hash, 0)
		if completed {
			return true
		}
	}
}

func (lfMap *LFMap) deleteRecursive(node *atomic.Value, key string, hash uint32, level int) bool {
	index := lfMap.getSparseIndex(hash, level)
	currNode := node.Load().(*LFMapNode)

	if ! isBitSet(currNode.BitMap, index) { 
		return true 
	} else {
		pos := lfMap.getPosition(currNode.BitMap, hash, level)
		childNode := currNode.Children[pos]
		
		if childNode.IsLeafNode {
			if key == childNode.Key {
				currNode.BitMap = setBit(currNode.BitMap, index)
				currNode.Children = ShrinkTable(currNode.Children, currNode.BitMap, pos)
				
				if node.CompareAndSwap(node.Load().(*LFMapNode), currNode) {
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
				currNode.BitMap = setBit(currNode.BitMap, index)
				currNode.Children = ShrinkTable(currNode.Children, currNode.BitMap, pos)
			} 

			if node.CompareAndSwap(node.Load().(*LFMapNode), currNode) {
				return true
			} else { return false }
		}
	}
}