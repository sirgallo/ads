package lfmap

import "encoding/binary"


const (
	c32_1 = 0x85ebca6b
	c32_2 = 0xc2b2ae35
	c32_3 = 0xe6546b64
	c32_4 = 0x1b873593
	c32_5 = 0x5c4bcea9

	c64_1 = 0xff51afd7ed558ccd
	c64_2 = 0xc4ceb9fe1a85ec53
	c64_3 = 0x7b6d5f86d192eaa1
	c64_4 = 0x4cf5ad432745937f
	c64_5 = 0x8a7d3eef7b5ea2e1
)


//=================================== 32 bit

func Murmur32(data string, seed uint32) uint32 {
	dataAsBytes := []byte(data)
	hash := seed
	
	length := uint32(len(dataAsBytes))
	total4ByteChunks := len(dataAsBytes) / 4
	
	for idx := range make([]int, total4ByteChunks) {
		startIdxOfChunk := idx * 4 
		endIdxOfChunk := (idx + 1) * 4
		chunk := binary.LittleEndian.Uint32(dataAsBytes[startIdxOfChunk:endIdxOfChunk])

		rotateRight32(&hash, chunk)
	}

	handleRemainingBytes32(&hash, dataAsBytes)

	hash ^= length
	hash ^= hash >> 16
	hash *= c32_4
	hash ^= hash >> 13
	hash *= c32_5
	hash ^= hash >> 16

	return hash
}

func rotateRight32(hash *uint32, chunk uint32) {
	chunk *= c32_1
	chunk = (chunk << 15) | (chunk >> 17) // Rotate right by 15
	chunk *= c32_2

	*hash ^= chunk
	*hash = (*hash << 13) | (*hash >> 19) // Rotate right by 13
	*hash = *hash * 5 + c32_3
}

func handleRemainingBytes32(hash *uint32, dataAsBytes []byte) {
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
				chunk *= c32_1
				chunk = (chunk << 15) | (chunk >> 17) // Rotate right by 15
				chunk *= c32_2
				*hash ^= chunk
			}
	}
}


//=================================== 64 bit

func Murmur64(data string, seed uint64) uint64 {
	dataAsBytes := []byte(data)
	hash := seed

	length := uint64(len(dataAsBytes))
	total8ByteChunks := len(dataAsBytes) / 8

	for idx := range make([]int, total8ByteChunks) {
		startIdxOfChunk := idx * 8
		endIdxOfChunk := (idx + 1) * 8
		chunk := binary.LittleEndian.Uint64(dataAsBytes[startIdxOfChunk:endIdxOfChunk])

		rotateRight64(&hash, chunk)
	}

	handleRemainingBytes64(&hash, dataAsBytes)

	hash ^= length
	hash ^= hash >> 33
	hash *= c64_4
	hash ^= hash >> 29
	hash *= c64_5
	hash ^= hash >> 32

	return hash
}

func rotateRight64(hash *uint64, chunk uint64) {
	chunk *= c64_1
	chunk = (chunk << 31) | (chunk >> 33) // Rotate right by 31
	chunk *= c64_2

	*hash ^= chunk
	*hash = (*hash << 27) | (*hash >> 37) // Rotate right by 27
	*hash = *hash * 5 + c64_3
}

func handleRemainingBytes64(hash *uint64, dataAsBytes []byte) {
	remaining := dataAsBytes[len(dataAsBytes)-len(dataAsBytes)%8:]

	if len(remaining) > 0 {
		var chunk uint64

		switch len(remaining) {
		case 7:
			chunk |= uint64(remaining[6]) << 48
			fallthrough
		case 6:
			chunk |= uint64(remaining[5]) << 40
			fallthrough
		case 5:
			chunk |= uint64(remaining[4]) << 32
			fallthrough
		case 4:
			chunk |= uint64(remaining[3]) << 24
			fallthrough
		case 3:
			chunk |= uint64(remaining[2]) << 16
			fallthrough
		case 2:
			chunk |= uint64(remaining[1]) << 8
			fallthrough
		case 1:
			chunk |= uint64(remaining[0])
			chunk *= c64_1
			chunk = (chunk << 31) | (chunk >> 33) // Rotate right by 31
			chunk *= c64_2
			*hash ^= chunk
		}
	}
}
