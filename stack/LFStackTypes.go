package stack

import "unsafe"

import "github.com/sirgallo/ads/counter"
import "github.com/sirgallo/ads/node"
import "github.com/sirgallo/ads/utils"


type LFStackOpts struct {
	ExpBackoffOpts utils.ExpBackoffOpts
	MaxStackSize int
}

type LFStack struct {
	top unsafe.Pointer
	nodePool *node.LFNodePool
	expBackoffOpts utils.ExpBackoffOpts
	length *counter.Counter
	maxStackSize int
}