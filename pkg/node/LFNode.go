package node

import "unsafe"

import "github.com/sirgallo/ads/pkg/counter"


type LFNode [T comparable] struct {
	Value T
	Next unsafe.Pointer
	Tag uintptr
}

type LFNodePool [T comparable] struct {
	PoolSize counter.Counter
	Pool chan *LFNode[T]
}