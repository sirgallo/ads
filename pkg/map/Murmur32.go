package lfmap

import "encoding/binary"


const constant1 = 0x85ebca6b
const constant2 = 0xc2b2ae35
const constant3 = 0xe6546b64
const constant4 = 0x1b873593
const	constant5 = 0x5c4bcea9


func Murmur32(data string, seed uint32) uint32 {
	dataAsBytes := []byte(data)
	hash := seed
	
	length := uint32(len(dataAsBytes))
	total4ByteChunks := len(dataAsBytes) / 4
	for idx, _ := range make([]int, total4ByteChunks) {
		startIdxOfChunk := idx * 4 
		endIdxOfChunk := (idx + 1) * 4
		chunk := binary.LittleEndian.Uint32(dataAsBytes[startIdxOfChunk:endIdxOfChunk])

		rotateRight(&hash, chunk)
	}

	handleRemainingBytes(&hash, dataAsBytes)

	hash ^= length
	hash ^= hash >> 16
	hash *= constant4
	hash ^= hash >> 13
	hash *= constant5
	hash ^= hash >> 16

	return hash
}

func rotateRight(hash *uint32, chunk uint32) {
	chunk *= constant1
	chunk = (chunk << 15) | (chunk >> 17) // Rotate right by 15
	chunk *= constant2

	*hash ^= chunk
	*hash = (*hash << 13) | (*hash >> 19) // Rotate right by 13
	*hash = *hash * 5 + constant3
}

func handleRemainingBytes(hash *uint32, dataAsBytes []byte) {
	remaining := dataAsBytes[len(dataAsBytes)-len(dataAsBytes) % 4:]
	
	if len(remaining) > 0 {
		var chunk uint32
		
		switch len(remaining) {
			case 3:
				chunk |= uint32(remaining[2]) << 16
				fallthrough
			case 2:
				chunk |= uint32(remaining[1]) << 8
				fallthrough
			case 1:
				chunk |= uint32(remaining[0])
				chunk *= constant1
				chunk = (chunk << 15) | (chunk >> 17) // Rotate right by 15
				chunk *= constant2
				*hash ^= chunk
			}
	}
}