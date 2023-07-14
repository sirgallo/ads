package lfmap

import "github.com/sirgallo/ads/pkg/counter"
import "github.com/sirgallo/ads/pkg/utils"


func NewLFMapNodePool[T comparable](poolSize int) *LFMapNodePool[T] {
	intertnalCounter, _ := counter.NewCounter(0)

	return &LFMapNodePool[T]{ 
		PoolSize: *intertnalCounter,
		Pool: make(chan *LFMapNode[T], poolSize),
	}
}

func (np *LFMapNodePool[T]) GetLFMapNode() *LFMapNode[T] {
	select {
		case node := <- np.Pool:
			np.PoolSize.Decrement(1)
			return node
		default:
			return &LFMapNode[T]{}
	}
}

func (np *LFMapNodePool[T]) PutLFMapNode(node *LFMapNode[T]) {
	// reset node
	node.Key = utils.GetZero[string]()
	node.Value = utils.GetZero[T]()
	node.IsLeafNode = false
	node.BitMap = 0
	node.Children = []*LFMapNode[T]{}

	select {
		case np.Pool <- node:
			np.PoolSize.Increment(1)
		default: // do nothing if pool is full
	}
}