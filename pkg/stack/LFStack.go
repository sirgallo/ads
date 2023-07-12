package stack

import "errors"
import "sync/atomic"
import "unsafe"

import "github.com/sirgallo/ads/pkg/counter"
import "github.com/sirgallo/ads/pkg/node"
import "github.com/sirgallo/ads/pkg/utils"


func NewLFStack[T comparable](opts LFStackOpts) *LFStack[T] {
	sCounter, _ := counter.NewCounter(0)
	nodePool := node.NewLFNodePool[T](opts.MaxStackSize)
	
	return &LFStack[T]{
		nodePool: nodePool,
		length: sCounter,
		maxStackSize: opts.MaxStackSize,
		expBackoffOpts: opts.ExpBackoffOpts,
	}
}

func (stack *LFStack[T]) Push(value T) (bool, error) {
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

func (stack *LFStack[T]) Pop() (T, error) {
	expBackoffStrat := utils.NewExponentialBackoffStrat(stack.expBackoffOpts)

	for {
		top := atomic.LoadPointer(&stack.top)
		if top == nil { return utils.GetZero[T](), nil }

		newTop := (*node.LFNode[T])(top).Next

		if atomic.CompareAndSwapPointer(&stack.top, top, newTop) {
			value := (*node.LFNode[T])(top).Value
			stack.nodePool.PutLFNode((*node.LFNode[T])(top))
			stack.length.Decrement(1)

			return value, nil
		}

		err := expBackoffStrat.PerformBackoff()
		if err != nil { return utils.GetZero[T](), err }
	}
}

func (stack *LFStack[T]) Peek() T {
	top := atomic.LoadPointer(&stack.top)
	if top == nil { return utils.GetZero[T]() }
	
	return (*node.LFNode[T])(top).Value
}

func (stack *LFStack[T]) Size() int64 {
	return stack.length.GetValue()
}

func (stack *LFStack[T]) MaxSize() int {
	return stack.maxStackSize
}

func (stack *LFStack[T]) Clear() bool {
	close(stack.nodePool.Pool)
	return true
}