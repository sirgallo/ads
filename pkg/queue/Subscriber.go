package lfqueue

import "github.com/google/uuid"

import "github.com/sirgallo/ads/pkg/stack"
import "github.com/sirgallo/ads/pkg/utils"


func NewSubscriber[T comparable](opts SubscriberOpts[T]) *Subscriber[T] {
	subscriberId := uuid.New()

	maxRetries := 10
	sExpBackoffOpts := utils.ExpBackoffOpts{ MaxRetries: &maxRetries, TimeoutInMicroseconds: 10 }
	sOpts := lfstack.LFStackOpts{ MaxStackSize: opts.StackSize, ExpBackoffOpts: sExpBackoffOpts }
	lfStack := lfstack.NewLFStack[QueueEntry[T]](sOpts)

	return &Subscriber[T] {
		subscriberId: subscriberId,
		lfQueue: opts.LFQueue,
		lfStack: lfStack,
		dequeueHandler: opts.DequeueHandler,
		terminationSignal: opts.TerminationSignal,
	}
}

func (subscriber *Subscriber[T]) Subscribe() (bool, error) {
	fillStackSignal := make(chan bool, 1)
	processStackSignal := make(chan bool, 1)

	select {
		case <- subscriber.terminationSignal:
			return true, nil
		default: 
			// here, we do not need max retries. Subscriber should stay alive until new elements enter the queue
			expBackoffStrat := utils.NewExponentialBackoffStrat(utils.ExpBackoffOpts{ TimeoutInMicroseconds: 10 })

			for {
				select {
					case <- processStackSignal: 
						if subscriber.lfStack.Size() > 0 { 
							dequeued, _ := subscriber.lfStack.Pop()
							subscriber.dequeueHandler(subscriber.subscriberId, dequeued)
							
							processStackSignal <- true 
						} else { fillStackSignal <- true }
					case <- fillStackSignal:
						if subscriber.lfStack.Size() < int64(subscriber.lfStack.MaxSize()) {
							dequeued, err := subscriber.lfQueue.Dequeue() 
							if err != nil { return false, err }
				
							if dequeued != utils.GetZero[QueueEntry[T]]() {
								subscriber.lfStack.Push(dequeued)
								expBackoffStrat.Reset() // successful dequeue, reset exp backoff
							} else { expBackoffStrat.PerformBackoff() }

							fillStackSignal <- true
						} else { processStackSignal <- true }
					default:
						fillStackSignal <- true
				} 
			}
	}
}