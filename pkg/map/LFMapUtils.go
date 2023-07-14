package lfmap

import "sync/atomic"
import "fmt"
import "math/bits"
import "math"
import "unsafe"


func (lfMap *LFMap[T]) CalculateHashForCurrentLevel(key string, level int) uint32 {
	currChunk := level / lfMap.TotalLevels
	seed := uint32(currChunk + 1)
	return Murmur32(key, seed)
}

func (lfMap *LFMap[T]) getSparseIndex(hash uint32, level int) int {
	return GetIndexForLevel(hash, lfMap.BitChunkSize, level, lfMap.TotalLevels)
}

func (lfMap *LFMap[T]) getPosition(bitMap uint32, hash uint32, level int) int {
	sparseIdx := GetIndexForLevel(hash, lfMap.BitChunkSize, level, lfMap.TotalLevels)
	mask := uint32((1 << sparseIdx) - 1)
	isolatedBits := bitMap & mask
	
	return calculateHammingWeight(isolatedBits)
}

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

func calculateHammingWeight(bitmap uint32) int {
	return bits.OnesCount32(bitmap)
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
	currNode := (*LFMapNode[T])(atomic.LoadPointer(node))
	if currNode == nil { return }
	for i, child := range currNode.Children {
		if child != nil {
			fmt.Printf("Level: %d, Index: %d, Key: %s, Value: %v\n", level, i, child.Key, child.Value)
			
			childPtr := unsafe.Pointer(child)
			lfMap.printChildrenRecursive(&childPtr, level+1)
		}
	}
}