package node

import "unsafe"


type LFNode [T comparable] struct {
	Value T
	Next unsafe.Pointer
	Tag uintptr
}

type LFNodePool [T comparable] struct {
	Pool chan *LFNode[T]
}