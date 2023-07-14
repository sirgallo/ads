package lfmaptests

import "testing"

import "github.com/sirgallo/ads/pkg/map"


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

	opts := lfmap.LFMapOpts{ PoolSize: 10000000 }
	lfMap := lfmap.NewLFMap[string](opts)

	for idx := range levels {
		hash := lfMap.CalculateHashForCurrentLevel(key, idx)
		index := lfmap.GetIndexForLevel(hash, chunkSize, idx, totalLevels)
		t.Logf("hash: %d, index: %d", hash, index)
	}
}