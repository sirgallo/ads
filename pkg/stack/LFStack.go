package stack

import "errors"
import "sync/atomic"
import "unsafe"

import "github.com/sirgallo/ads/pkg/counter"
import "github.com/sirgallo/ads/pkg/node"
import "github.com/sirgallo/ads/pkg/utils"


func NewLFStack(opts LFStackOpts) *LFStack {
	sCounter, _ := counter.NewCounter(0)
	nodePool := node.NewLFNodePool(opts.MaxStackSize)
	
	return &LFStack{
		nodePool: nodePool,
		length: sCounter,
		maxStackSize: opts.MaxStackSize,
		expBackoffOpts: opts.ExpBackoffOpts,
	}
}

func (stack *LFStack) Push(value interface{}) (bool, error) {
	expBackoffStrat := utils.NewExponentialBackoffStrat(stack.expBackoffOpts)

	newNode := stack.nodePool.GetLFNode()
	newNode.Value = value

	for {
		if stack.Size() == int64(stack.maxStackSize) { return false, errors.New("max stack size reached") }

		top := atomic.LoadPointer(&stack.top)
		newNode.Next = top

		if atomic.CompareAndSwapPointer(&stack.top, top, unsafe.Pointer(newNode)) {
			stack.length.Increment(1)
			return true, nil
		} 

		err := expBackoffStrat.PerformBackoff()
		if err != nil { return false, err }
	}
}

func (stack *LFStack) Pop() (interface{}, error) {
	expBackoffStrat := utils.NewExponentialBackoffStrat(stack.expBackoffOpts)

	for {
		top := atomic.LoadPointer(&stack.top)
		if top == nil { return nil, nil }

		newTop := (*node.LFNode)(top).Next

		if atomic.CompareAndSwapPointer(&stack.top, top, newTop) {
			value := (*node.LFNode)(top).Value
			stack.nodePool.PutLFNode((*node.LFNode)(top))
			stack.length.Decrement(1)

			return value, nil
		}

		err := expBackoffStrat.PerformBackoff()
		if err != nil { return false, err }
	}
}

func (stack *LFStack) Peek() interface{} {
	top := atomic.LoadPointer(&stack.top)
	if top == nil { return nil }
	
	return (*node.LFNode)(top).Value
}

func (stack *LFStack) Size() int64 {
	return stack.length.GetValue()
}

func (stack *LFStack) MaxSize() int {
	return stack.maxStackSize
}

func (stack *LFStack) Clear() bool {
	close(stack.nodePool.Pool)
	return true
}