package node


func NewLFNodePool(poolSize int) *LFNodePool {
	return &LFNodePool{ 
		Pool: make(chan *LFNode, poolSize),
	}
}

func (np *LFNodePool) GetLFNode() *LFNode {
	select {
		case node := <- np.Pool:
			node.Next = nil
			return node
		default:
			return &LFNode{}
	}
}

func (np *LFNodePool) PutLFNode(node *LFNode) {
	select {
		case np.Pool <- node:
		default: // do nothing if pool is full
	}
}