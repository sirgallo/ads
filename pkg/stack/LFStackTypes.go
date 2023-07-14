package lfstack

import "unsafe"

import "github.com/sirgallo/ads/pkg/counter"
import "github.com/sirgallo/ads/pkg/node"
import "github.com/sirgallo/ads/pkg/utils"


type LFStackOpts struct {
	ExpBackoffOpts utils.ExpBackoffOpts
	MaxStackSize int
}

type LFStack [T comparable] struct {
	top unsafe.Pointer
	nodePool *node.LFNodePool[T]
	expBackoffOpts utils.ExpBackoffOpts
	length *counter.Counter
	maxStackSize int
}