package lfmap

import "unsafe"


type LFMapNode [T comparable, V uint32 | uint64] struct {
	Key string
	Value T
	IsLeafNode bool
	BitMap V
	Children []*LFMapNode[T, V]
}

type LFMap [T comparable, V uint32 | uint64] struct {
	BitChunkSize int
	HashChunks int
	Root unsafe.Pointer
}