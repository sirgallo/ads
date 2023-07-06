package node

import "unsafe"


type LFNode struct {
	Value interface{}
	Next unsafe.Pointer
	Tag uintptr
}

type LFNodePool struct {
	Pool chan *LFNode
}