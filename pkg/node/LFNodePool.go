package node


func NewLFNodePool[T comparable](poolSize int) *LFNodePool[T] {
	return &LFNodePool[T]{ 
		Pool: make(chan *LFNode[T], poolSize),
	}
}

func (np *LFNodePool[T]) GetLFNode() *LFNode[T] {
	select {
		case node := <- np.Pool:
			node.Next = nil
			return node
		default:
			return &LFNode[T]{}
	}
}

func (np *LFNodePool[T]) PutLFNode(node *LFNode[T]) {
	select {
		case np.Pool <- node:
		default: // do nothing if pool is full
	}
}