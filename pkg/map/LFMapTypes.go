package lfmap

import "sync/atomic"

import "github.com/sirgallo/ads/pkg/utils"


type LFMapOpts struct {
	ExpBackoffOpts utils.ExpBackoffOpts
	MaxPoolSize int
}

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
	// nodePool *LFMapNodePool
	// expBackoffOpts utils.ExpBackoffOpts
}