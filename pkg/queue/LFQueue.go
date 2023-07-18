package lfqueue

import "errors"
import "fmt"
import "sync/atomic"
import "unsafe"

import "github.com/sirgallo/ads/pkg/counter"
import "github.com/sirgallo/ads/pkg/node"
import "github.com/sirgallo/ads/pkg/utils"


func NewLFQueue[T comparable](opts LFQueueOpts) *LFQueue[T] {
	fmt.Println("Initializing Lock Free Queue")
	
	qCounter, _ := counter.NewCounter(0)
	nodePool := node.NewLFNodePool[T](opts.MaxPoolSize)
	node := unsafe.Pointer(nodePool.GetLFNode())
	
	return &LFQueue[T] {
		head: node,
		tail: node,
		nodePool: nodePool,
		length: qCounter,
		maxPoolSize: opts.MaxPoolSize,
		expBackoffOpts: opts.ExpBackoffOpts,
	}
}

func (queue *LFQueue[T]) Enqueue(incoming T) (bool, error) {
	if queue.Size() >= int64(queue.maxPoolSize) { 
		return false, errors.New("max queue size reached, unable to enqueue") 
	}

	expBackoffStrat := utils.NewExponentialBackoffStrat(queue.expBackoffOpts)

	newNode := queue.nodePool.GetLFNode()
	newNode.Value = incoming

	var tail, next unsafe.Pointer
	var tag uintptr

	for {
		tail = atomic.LoadPointer(&queue.tail)
		next = atomic.LoadPointer(&(*node.LFNode[T])(tail).Next)
		tag = atomic.LoadUintptr(&(*node.LFNode[T])(tail).Tag)

		if tail == atomic.LoadPointer(&queue.tail) && tag == atomic.LoadUintptr(&(*node.LFNode[T])(tail).Tag) {
			if next == nil {
				newNode.Tag = tag + 1

				if atomic.CompareAndSwapPointer(&(*node.LFNode[T])(tail).Next, nil, unsafe.Pointer(newNode)) {
					atomic.CompareAndSwapPointer(&queue.tail, tail, unsafe.Pointer(newNode))
					queue.length.Increment(1)
					
					return true, nil
				}
			} else { atomic.CompareAndSwapPointer(&queue.tail, tail, next) }
		}

		err := expBackoffStrat.PerformBackoff()
		if err != nil { return false, err }
	}
}

func (queue *LFQueue[T]) Dequeue() (T, error) {
	expBackoffStrat := utils.NewExponentialBackoffStrat(queue.expBackoffOpts)

	var head, tail, next unsafe.Pointer
	var tag uintptr

	for {
		head = atomic.LoadPointer(&queue.head)
		tail = atomic.LoadPointer(&queue.tail)
		next = atomic.LoadPointer(&(*node.LFNode[T])(head).Next)
		tag = atomic.LoadUintptr(&(*node.LFNode[T])(head).Tag)

		if head == atomic.LoadPointer(&queue.head) && tag == atomic.LoadUintptr(&(*node.LFNode[T])(head).Tag) {
			if head == tail { 
				if (*node.LFNode[T])(head).Value == utils.GetZero[T]() { return utils.GetZero[T](), nil }
				atomic.CompareAndSwapPointer(&queue.tail, tail, next)
			} else {
				if next != nil {
					value := (*node.LFNode[T])(next).Value
					if atomic.CompareAndSwapPointer(&queue.head, head, next) {
						atomic.AddUintptr(&(*node.LFNode[T])(head).Tag, 1)

						queue.nodePool.PutLFNode((*node.LFNode[T])(head))
						queue.length.Decrement(1)

						return value, nil
					}
				}
			}
		}

		err := expBackoffStrat.PerformBackoff()
		if err != nil { return utils.GetZero[T](), err }
	}
}

func (queue *LFQueue[T]) Peek() T {
	head := atomic.LoadPointer(&queue.head) // head is dummy value, get value of next pointer
	next := atomic.LoadPointer(&(*node.LFNode[T])(head).Next)
	
	return (*node.LFNode[T])(next).Value
}

func (queue *LFQueue[T]) Size() int64 {
	return queue.length.GetValue()
}

func (queue *LFQueue[T]) MaxSize() int {
	return queue.maxPoolSize
}

func (queue *LFQueue[T]) Clear() bool {
	_, err := queue.length.Reset()
	if err != nil {
		return false
	}
	
	if atomic.CompareAndSwapPointer(&queue.head, unsafe.Pointer(&queue.head), nil) {
		close(queue.nodePool.Pool)
		return true
	}
	
	return false
}

func (entry QueueEntry[T]) String() string {
	return fmt.Sprintf("{ timestamp: %s, message: %T, producerId: %s }", entry.Timestamp.Format("2006-01-02 15:04:05"), entry.Message, entry.ProducerId)
}