package lfmap

import "unsafe"

import "github.com/sirgallo/ads/pkg/counter"


type LFMapOpts struct {
	PoolSize int
}

type LFMapNode [T comparable, V uint32 | uint64] struct {
	Key string
	Value T
	IsLeafNode bool
	BitMap V
	Children []*LFMapNode[T, V]
}

type LFMapNodePool [T comparable, V uint32 | uint64] struct {
	PoolSize counter.Counter
	Pool chan *LFMapNode[T, V]
}

type LFMap [T comparable, V uint32 | uint64] struct {
	BitChunkSize int
	HashChunks int
	Is64Bit bool
	NodePool *LFMapNodePool[T, V]
	Root unsafe.Pointer
}