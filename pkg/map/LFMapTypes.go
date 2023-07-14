package lfmap

import "unsafe"

import "github.com/sirgallo/ads/pkg/counter"
import "github.com/sirgallo/ads/pkg/utils"


type LFMapOpts struct {
	ExpBackoffOpts utils.ExpBackoffOpts
	PoolSize int
}

type LFMapNode [T comparable] struct {
	Key string
	Value T
	IsLeafNode bool
	BitMap uint32
	Children []*LFMapNode[T]
}

type LFMapNodePool [T comparable] struct {
	PoolSize counter.Counter
	Pool chan *LFMapNode[T]
}

type LFMap [T comparable] struct {
	BitChunkSize int
	TotalLevels int
	NodePool *LFMapNodePool[T]
	Root unsafe.Pointer
	expBackoffOpts utils.ExpBackoffOpts
}

type KeyHashState struct {
	Key string
	Hash uint32
	Level int
}