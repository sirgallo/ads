package queue

import "errors"
import "fmt"
import "sync/atomic"
import "unsafe"

import "github.com/sirgallo/ads/pkg/counter"
import "github.com/sirgallo/ads/pkg/node"
import "github.com/sirgallo/ads/pkg/utils"


func NewLFQueue(opts LFQueueOpts) *LFQueue {
	fmt.Println("Initializing Lock Free Queue")
	
	qCounter, _ := counter.NewCounter(0)
	nodePool := node.NewLFNodePool(opts.MaxQueueSize)
	node := unsafe.Pointer(nodePool.GetLFNode())
	
	return &LFQueue {
		head: node,
		tail: node,
		nodePool: nodePool,
		length: qCounter,
		maxQueueSize: opts.MaxQueueSize,
		expBackoffOpts: opts.ExpBackoffOpts,
	}
}

func (queue *LFQueue) Enqueue(incoming interface{}) (bool, error) {
	if queue.Size() >= int64(queue.maxQueueSize) { 
		return false, errors.New("max queue size reached, unable to enqueue") 
	}

	expBackoffStrat := utils.NewExponentialBackoffStrat(queue.expBackoffOpts)

	newNode := queue.nodePool.GetLFNode()
	newNode.Value = incoming

	var tail, next unsafe.Pointer
	var tag uintptr

	for {
		tail = atomic.LoadPointer(&queue.tail)
		next = atomic.LoadPointer(&(*node.LFNode)(tail).Next)
		tag = atomic.LoadUintptr(&(*node.LFNode)(tail).Tag)

		if tail == atomic.LoadPointer(&queue.tail) && tag == atomic.LoadUintptr(&(*node.LFNode)(tail).Tag) {
			if next == nil {
				newNode.Tag = tag + 1

				if atomic.CompareAndSwapPointer(&(*node.LFNode)(tail).Next, nil, unsafe.Pointer(newNode)) {
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

func (queue *LFQueue) Dequeue() (interface{}, error) {
	expBackoffStrat := utils.NewExponentialBackoffStrat(queue.expBackoffOpts)

	// var NilQueueEntry = QueueEntry{} 
	var head, tail, next unsafe.Pointer
	var tag uintptr

	for {
		head = atomic.LoadPointer(&queue.head)
		tail = atomic.LoadPointer(&queue.tail)
		next = atomic.LoadPointer(&(*node.LFNode)(head).Next)
		tag = atomic.LoadUintptr(&(*node.LFNode)(head).Tag)

		if head == atomic.LoadPointer(&queue.head) && tag == atomic.LoadUintptr(&(*node.LFNode)(head).Tag) {
			if head == tail { 
				if head == nil { return nil, nil }
				atomic.CompareAndSwapPointer(&queue.tail, tail, next)
			} else {
				if next != nil {
					value := (*node.LFNode)(next).Value
					if atomic.CompareAndSwapPointer(&queue.head, head, next) {
						atomic.AddUintptr(&(*node.LFNode)(head).Tag, 1)

						queue.nodePool.PutLFNode((*node.LFNode)(head))
						queue.length.Decrement(1)

						return value, nil
					}
				}
			}
		}

		err := expBackoffStrat.PerformBackoff()
		if err != nil { return nil, err }
	}
}

func (queue *LFQueue) Peek() interface{} {
	head := atomic.LoadPointer(&queue.head) // head is dummy value, get value of next pointer
	next := atomic.LoadPointer(&(*node.LFNode)(head).Next)
	
	return (*node.LFNode)(next).Value.(QueueEntry)
}

func (queue *LFQueue) Size() int64 {
	return queue.length.GetValue()
}

func (queue *LFQueue) MaxSize() int {
	return queue.maxQueueSize
}

func (queue *LFQueue) Clear() bool {
	close(queue.nodePool.Pool)
	return true
}

func (entry QueueEntry) String() string {
	return fmt.Sprintf("{ timestamp: %s, message: %s, producerId: %s }", entry.Timestamp.Format("2006-01-02 15:04:05"), entry.Message, entry.ProducerId)
}