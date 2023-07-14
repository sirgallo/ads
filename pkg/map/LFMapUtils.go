package lfmap

import "sync/atomic"
import "fmt"
import "math/bits"
import "math"
import "unsafe"


func (lfMap *LFMap[T, V]) CalculateHashForCurrentLevel(key string, level int) V {
	currChunk := level / lfMap.HashChunks

	var v V 
	switch any(v).(type) {
		case uint64:
			seed := uint64(currChunk + 1)
			return (V)(Murmur64(key, seed))
		default:
			seed := uint32(currChunk + 1)
			return (V)(Murmur32(key, seed))
	}
}

func (lfMap *LFMap[T, V]) getSparseIndex(hash V, level int) int {
	return GetIndexForLevel(hash, lfMap.BitChunkSize, level, lfMap.HashChunks)
}

func (lfMap *LFMap[T, V]) getPosition(bitMap V, hash V, level int) int {
	sparseIdx := GetIndexForLevel(hash, lfMap.BitChunkSize, level, lfMap.HashChunks)
	
	switch any(bitMap).(type) {
		case uint64:
			mask := uint64((1 << sparseIdx) - 1)
			isolatedBits := (uint64)(bitMap) & mask
			return calculateHammingWeight(isolatedBits)
		default:
			mask := uint32((1 << sparseIdx) - 1)
			isolatedBits := (uint32)(bitMap) & mask
			return calculateHammingWeight(isolatedBits)
	}
}

func GetIndexForLevel[V uint32 | uint64](hash V, chunkSize int, level int, hashChunks int) int {
	updatedLevel := level % hashChunks
	return GetIndex(hash, chunkSize, updatedLevel)
}

func GetIndex[V uint32 | uint64](hash V, chunkSize int, level int) int {
	slots := int(math.Pow(float64(2), float64(chunkSize)))
	shiftSize := slots - (chunkSize * (level + 1))

	switch any(hash).(type) {
		case uint64:
			mask := uint64(slots - 1)
			return int((uint64)(hash) >> shiftSize & mask)
		default:
			mask := uint32(slots - 1)
			return int((uint32)(hash) >> shiftSize & mask)
	}
}

func calculateHammingWeight[V uint32 | uint64](bitmap V) int {
	switch any(bitmap).(type) {
		case uint64:
			return bits.OnesCount64((uint64)(bitmap))
		default:
			return bits.OnesCount32((uint32)(bitmap))
	}
}

func SetBit[V uint32 | uint64](bitmap V, position int) V {
	return bitmap ^ (1 <<  position)
}

func IsBitSet[V uint32 | uint64](bitmap V, position int) bool {
	return (bitmap & (1 << position)) != 0
}

func ExtendTable[T comparable, V uint32 | uint64](orig []*LFMapNode[T, V], bitMap V, pos int, newNode *LFMapNode[T, V]) []*LFMapNode[T, V] {
	tableSize := calculateHammingWeight(bitMap)
	newTable := make([]*LFMapNode[T, V], tableSize)
	
	copy(newTable[:pos], orig[:pos])
	newTable[pos] = newNode
	copy(newTable[pos + 1:], orig[pos:])
	
	return newTable
}

func ShrinkTable[T comparable, V uint32 | uint64](orig []*LFMapNode[T, V], bitMap V, pos int) []*LFMapNode[T, V] {
	tableSize := calculateHammingWeight(bitMap)
	newTable := make([]*LFMapNode[T, V], tableSize)
	
	copy(newTable[:pos], orig[:pos])
	copy(newTable[pos:], orig[pos + 1:])

	return newTable
}


// for debugging

func (lfMap *LFMap[T, V]) PrintChildren() {
	lfMap.printChildrenRecursive(&lfMap.Root, 0)
}

func (lfMap *LFMap[T, V]) printChildrenRecursive(node *unsafe.Pointer, level int) {
	currNode := (*LFMapNode[T, V])(atomic.LoadPointer(node))
	if currNode == nil { return }

	for idx, child := range currNode.Children {
		if child != nil {
			fmt.Printf("Level: %d, Index: %d, Key: %s, Value: %v\n", level, idx, child.Key, child.Value)
			
			childPtr := unsafe.Pointer(child)
			lfMap.printChildrenRecursive(&childPtr, level + 1)
		}
	}
}