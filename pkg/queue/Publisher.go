package lfqueue 

import "time"
import "github.com/google/uuid"


func NewPublisher[T comparable](opts PublisherOpts[T]) *Publisher[T] {
	publisherId := uuid.New()
	return &Publisher[T]{
		publisherId: publisherId,
		lfQueue: opts.LFQueue,
	}
}

func (publisher *Publisher[T]) Publish(message T) (bool, error) {
	queueEntry := QueueEntry[T]{ Timestamp: time.Now(), Message: message, ProducerId: publisher.publisherId.String() }
	_, err := publisher.lfQueue.Enqueue(queueEntry)

	if err != nil { return false, err }

	return true, nil
}