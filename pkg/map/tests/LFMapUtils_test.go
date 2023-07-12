package lfmaptests

import "testing"
import "fmt"

import "github.com/sirgallo/ads/pkg/map"
import "github.com/sirgallo/ads/pkg/utils"


func TestGetIndex(t *testing.T) {
	chunkSize := 5

	key1 := "hello"
	hash1 := utils.FnvHash(key1)

	fmt.Printf("hash 1: %032b\n:", hash1)

	expectedValues1 := []int{9, 30, 15, 18, 25, 10}

	for idx, val := range expectedValues1 {
		index := lfmap.GetIndex(hash1, chunkSize, idx)
		t.Logf("index: %d, expected: %d", index, val)
		if index != val {
			t.Error("index produced does not match expected value")
		}
	}

	key2 := "new"
	hash2 := utils.FnvHash(key2)

	fmt.Printf("hash 2: %032b\n:", hash2)

	expectedValues2 := []int{5, 2, 12, 25, 12, 4}

	for idx, val := range expectedValues2 {
		index := lfmap.GetIndex(hash2, chunkSize, idx)
		t.Logf("index: %d, expected: %d", index, val)
		if index != val {
			t.Error("index produced does not match expected value")
		}
	}
}

/*
          0     1     2     3     4     5    extra
hello = 01001 11110 01111 10010 11001 01010 11 
ignore last bit

level 0 = 01001 = 9 
level 1 = 11110 = 30
level 2 = 01111 = 15
level 3 = 10010 = 18
level 4 = 11001 = 25
level 5 = 01010 = 10

so at each shift 

shift 27
level 0 = 00000 00000 00000 00000 00000 00010 01  --> shifted 27

shift 22
level 1 = 00000 00000 00000 00000 00010 01111 10 --> shifted 22

shift 17
level 2 = 00000 00000 00000 00010 01111 10011 11 --> shifted 17

shift 12
level 3 = 00000 00000 00010 01111 10011 11100 10 --> shifted 12

shift 7
level 4 = 00000 00010 01111 10011 11100 10110 01 --> shifted 7

shift 2
level 5 = 00010 01111 10011 11100 10110 01010 10 --> shifted 2
*/

/*
			  0     1     2     3     4     5    extra
new = 00101 00010 01100 11001 01100 00100 01

level 0 = 00101 = 5
level 1 = 00010 = 2
level 2 = 01100 = 12
level 3 = 11001 = 25
level 4 = 01100 = 12
level 5 = 00100 = 4
*/


func TestSetBitmap(t *testing.T) {
	bitmap := uint32(0)
	index1 := 1

	bitmap = lfmap.SetBit(bitmap, index1)
	fmt.Printf("current bitmap: %032b\n", bitmap)

	isBitSet1 := lfmap.IsBitSet(bitmap, index1)
	if ! isBitSet1 {
		t.Error("bit at index 1 is not set")
	}

	index5 := 5

	bitmap = lfmap.SetBit(bitmap, index5)
	fmt.Printf("current bitmap: %032b\n", bitmap)
	
	isBitSet5 := lfmap.IsBitSet(bitmap, index5)
	if ! isBitSet5 {
		t.Error("bit at index 5 is not set")
	}
}