package queue 

import "time"
import "github.com/google/uuid"


func NewPublisher(opts PublisherOpts) *Publisher {
	publisherId := uuid.New()
	return &Publisher{
		publisherId: publisherId,
		lfQueue: opts.LFQueue,
	}
}

func (publisher *Publisher) Publish(message interface{}) (bool, error) {
	queueEntry := QueueEntry{ Timestamp: time.Now(), Message: message, ProducerId: publisher.publisherId.String() }
	_, err := publisher.lfQueue.Enqueue(queueEntry)

	if err != nil { return false, err }

	return true, nil
}