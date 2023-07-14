package lfmaptests

import "testing"

import "github.com/sirgallo/ads/pkg/map"


//=================================== 32 bit

func TestMurmur32(t *testing.T) {
	key := "hello"
	seed := uint32(1)

	hash := lfmap.Murmur32(key, seed)
	t.Log("hash:", hash)
}

func TestMurmur32ReSeed(t *testing.T) {
	key := "hello"
	levels := make([]int, 17)
	totalLevels := 6
	chunkSize := 5

	lfMap := lfmap.NewLFMap[string, uint32]()

	for idx := range levels {
		hash := lfMap.CalculateHashForCurrentLevel(key, idx)
		index := lfmap.GetIndexForLevel(hash, chunkSize, idx, totalLevels)
		t.Logf("hash: %d, index: %d", hash, index)
	}
}


//=================================== 64 bit

func TestMurmur64(t *testing.T) {
	key := "hello"
	seed := uint64(1)

	hash := lfmap.Murmur64(key, seed)
	t.Log("hash:", hash)
}

func TestMurmur64ReSeed(t *testing.T) {
	key := "hello"
	levels := make([]int, 33)
	totalLevels := 10
	chunkSize := 6

	lfMap := lfmap.NewLFMap[string, uint64]()

	for idx := range levels {
		hash := lfMap.CalculateHashForCurrentLevel(key, idx)
		index := lfmap.GetIndexForLevel(hash, chunkSize, idx, totalLevels)
		t.Logf("hash: %d, index: %d", hash, index)
	}
}