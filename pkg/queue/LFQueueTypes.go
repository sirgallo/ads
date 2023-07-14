package lfqueue

import "time"
import "unsafe"
import "github.com/google/uuid"

import "github.com/sirgallo/ads/pkg/counter"
import "github.com/sirgallo/ads/pkg/node"
import "github.com/sirgallo/ads/pkg/stack"
import "github.com/sirgallo/ads/pkg/utils"


type LFQueueOpts struct {
	ExpBackoffOpts utils.ExpBackoffOpts
	MaxPoolSize int
}

type LFQueue [T comparable] struct {
	head unsafe.Pointer
	tail unsafe.Pointer
	nodePool *node.LFNodePool[T]
	length *counter.Counter
	maxPoolSize int
	expBackoffOpts utils.ExpBackoffOpts
}

type QueueEntry [T comparable] struct {
  Timestamp time.Time
	Message T
	ProducerId string
}

type PublisherOpts [T comparable] struct {
	LFQueue *LFQueue[QueueEntry[T]]
}

type Publisher [T comparable] struct {
	publisherId uuid.UUID
	lfQueue *LFQueue[QueueEntry[T]]
}

type SubscriberOpts [T comparable] struct {
	LFQueue *LFQueue[QueueEntry[T]]
	DequeueHandler func(subscriberId uuid.UUID, dequeued QueueEntry[T]) bool
	StackSize int
	TerminationSignal chan bool
}

type Subscriber [T comparable] struct {
	subscriberId uuid.UUID
	lfQueue *LFQueue[QueueEntry[T]]
	lfStack *lfstack.LFStack[QueueEntry[T]]
	dequeueHandler func(subscriberId uuid.UUID, dequeued QueueEntry[T]) bool
	terminationSignal chan bool
}