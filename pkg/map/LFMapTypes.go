package lfmap

import "sync/atomic"


type LFMapNode [T comparable] struct {
	Key string
	Value T
	IsLeafNode bool
	BitMap uint32
	Children []*LFMapNode[T]
}

type LFMapNodePool [T comparable] struct {
	Pool chan *LFMapNode[T]
}

type LFMap [T comparable] struct {
	Root *atomic.Value
	BitChunkSize int
	TotalChildren int
}