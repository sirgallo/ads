package node

import "github.com/sirgallo/ads/pkg/counter"


func NewLFNodePool[T comparable](poolSize int) *LFNodePool[T] {
	intertnalCounter, _ := counter.NewCounter(0)

	return &LFNodePool[T]{ 
		PoolSize: *intertnalCounter,
		Pool: make(chan *LFNode[T], poolSize),
	}
}

func (np *LFNodePool[T]) GetLFNode() *LFNode[T] {
	select {
		case node := <- np.Pool:
			np.PoolSize.Decrement(1)
			node.Next = nil
			return node
		default:
			return &LFNode[T]{}
	}
}

func (np *LFNodePool[T]) PutLFNode(node *LFNode[T]) {
	select {
		case np.Pool <- node:
			np.PoolSize.Increment(1)
		default: // do nothing if pool is full
	}
}