package queue

import "time"
import "unsafe"
import "github.com/google/uuid"

import "github.com/sirgallo/ads/counter"
import "github.com/sirgallo/ads/node"
import "github.com/sirgallo/ads/stack"
import "github.com/sirgallo/ads/utils"


type LFQueueOpts struct {
	ExpBackoffOpts utils.ExpBackoffOpts
	MaxQueueSize int
}

type LFQueue struct {
	head unsafe.Pointer
	tail unsafe.Pointer
	nodePool *node.LFNodePool
	length *counter.Counter
	maxQueueSize int
	expBackoffOpts utils.ExpBackoffOpts
}

type QueueEntry struct {
  Timestamp time.Time
	Message interface{}
	ProducerId string
}

type PublisherOpts struct {
	LFQueue *LFQueue
}

type Publisher struct {
	publisherId uuid.UUID
	lfQueue *LFQueue
}

type SubscriberOpts struct {
	LFQueue *LFQueue
	DequeueHandler func(subscriberId uuid.UUID, dequeued interface{}) bool
	StackSize int
	TerminationSignal chan bool
}

type Subscriber struct {
	subscriberId uuid.UUID
	lfQueue *LFQueue
	lfStack *stack.LFStack
	dequeueHandler func(subscriberId uuid.UUID, dequeued interface{}) bool
	terminationSignal chan bool
}