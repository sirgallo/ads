package lfmap

import "github.com/sirgallo/ads/pkg/counter"


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
	select {
		case np.Pool <- node:
			np.PoolSize.Increment(1)
		default: // do nothing if pool is full
	}
}