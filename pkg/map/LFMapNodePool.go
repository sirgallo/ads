package lfmap

import "github.com/sirgallo/ads/pkg/counter"
import "github.com/sirgallo/ads/pkg/utils"


func NewLFMapNodePool[T comparable, V uint32 | uint64](poolSize int) *LFMapNodePool[T, V] {
	intertnalCounter, _ := counter.NewCounter(0)

	return &LFMapNodePool[T, V]{ 
		PoolSize: *intertnalCounter,
		Pool: make(chan *LFMapNode[T, V], poolSize),
	}
}

func (np *LFMapNodePool[T, V]) GetLFMapNode() *LFMapNode[T, V] {
	select {
		case node := <- np.Pool:
			np.PoolSize.Decrement(1)
			return node
		default:
			return &LFMapNode[T, V]{}
	}
}

func (np *LFMapNodePool[T, V]) PutLFMapNode(node *LFMapNode[T, V]) {
	// reset node
	node.Key = utils.GetZero[string]()
	node.Value = utils.GetZero[T]()
	node.IsLeafNode = false
	node.BitMap = 0
	node.Children = []*LFMapNode[T, V]{}

	select {
		case np.Pool <- node:
			np.PoolSize.Increment(1)
		default: // do nothing if pool is full
	}
}