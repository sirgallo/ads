package lfmap

import "sync/atomic"


type LFMapNode struct {
	Key string
	Value interface{}
	IsLeafNode bool
	BitMap uint32
	Children []*LFMapNode
}

type LFMapNodePool struct {
	Pool chan *LFMapNode
}

type LFMap struct {
	Root *atomic.Value
	BitChunkSize int
	TotalChildren int
}