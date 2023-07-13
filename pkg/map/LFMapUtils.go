package lfmap

import "sync/atomic"
import "fmt"
import "math/bits"
import "math"
import "unsafe"


func (lfMap *LFMap[T]) getSparseIndex(hash uint32, level int) int {
	return GetIndex(hash, lfMap.BitChunkSize, level)
}

func (lfMap *LFMap[T]) getPosition(bitMap uint32, hash uint32, level int) int {
	sparseIdx := GetIndex(hash, lfMap.BitChunkSize, level)
	mask := uint32((1 << sparseIdx) - 1)
	isolatedBits := bitMap & mask
	
	return calculateHammingWeight(isolatedBits)
}

func calculateHammingWeight(bitmap uint32) int {
	return bits.OnesCount32(bitmap)
}

func GetIndex(hash uint32, chunkSize int, level int) int {
	slots := int(math.Pow(float64(2), float64(chunkSize)))
	mask := uint32(slots - 1)
	shiftSize := slots - (chunkSize * (level + 1))

	return int(hash >> shiftSize & mask)
}

func SetBit(bitmap uint32, position int) uint32 {
	return bitmap ^ (1 <<  position)
}

func IsBitSet(bitmap uint32, position int) bool {
	return (bitmap & (1 << position)) != 0
}

func ExtendTable[T comparable](orig []*LFMapNode[T], bitMap uint32, pos int, newNode *LFMapNode[T]) []*LFMapNode[T] {
	tableSize := calculateHammingWeight(bitMap)
	newTable := make([]*LFMapNode[T], tableSize)
	
	copy(newTable[:pos], orig[:pos])
	newTable[pos] = newNode
	copy(newTable[pos + 1:], orig[pos:])
	
	return newTable
}

func ShrinkTable[T comparable](orig []*LFMapNode[T], bitMap uint32, pos int) []*LFMapNode[T] {
	tableSize := calculateHammingWeight(bitMap)
	newTable := make([]*LFMapNode[T], tableSize)
	
	copy(newTable[:pos], orig[:pos])
	copy(newTable[pos:], orig[pos + 1:])

	return newTable
}


// for debugging

func (lfMap *LFMap[T]) PrintChildren() {
	lfMap.printChildrenRecursive(&lfMap.Root, 0)
}

func (lfMap *LFMap[T]) printChildrenRecursive(node *unsafe.Pointer, level int) {
	// currNode := node.Load().(*LFMapNode[T])
	currNode := (*LFMapNode[T])(atomic.LoadPointer(node))
	if currNode == nil { return }
	// fmt.Println("currNode:", currNode)
	for i, child := range currNode.Children {
		if child != nil {
			fmt.Printf("Level: %d, Index: %d, Key: %s, Value: %v\n", level, i, child.Key, child.Value)
			
			// atomicChild := atomic.Value{}
			// atomicChild.Store(child)
			childPtr := unsafe.Pointer(child)
			lfMap.printChildrenRecursive(&childPtr, level+1)
		}
	}
}