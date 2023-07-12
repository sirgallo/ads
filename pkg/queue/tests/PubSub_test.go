package lfqueuetests

import "fmt"
import "math/rand"
import "sync"
import "testing"
import "time"
import "github.com/google/uuid"

import "github.com/sirgallo/ads/pkg/queue"
import "github.com/sirgallo/ads/pkg/utils"


func TestPubSub(t *testing.T) {
	maxRetries := 10
  expBackoffOpts := utils.ExpBackoffOpts{ MaxRetries: &maxRetries, TimeoutInMicroseconds: 1 }
  qOpts := queue.LFQueueOpts{ MaxPoolSize: 10000, ExpBackoffOpts: expBackoffOpts }

	q := queue.NewLFQueue[queue.QueueEntry[string]](qOpts)

	var publisherWG sync.WaitGroup
		
	for range make([]int, 10) {
		publisherWG.Add(1)
		
		go func () bool {
			defer publisherWG.Done()
			
			pOpts := queue.PublisherOpts[string]{ LFQueue: q }
			publisher := queue.NewPublisher(pOpts)
			messageNumber := 1000

			for range make([]int, messageNumber) {
				_, err := publisher.Publish("hello world")
				if err != nil { fmt.Println(err) }
				
				// we'll randomly publish between every 1 to 5 microseconds
				randNum := rand.Intn(5) + 1
				time.Sleep(time.Duration(randNum) * time.Microsecond) // simulate timeout
			}

			return true
		}()
	}

	publisherWG.Wait()

	expectedSize := 10000
	qSizeAfterPublishing := q.Size()

	t.Logf("actual queue size: %d, expected queue size: %d", qSizeAfterPublishing, expectedSize)
	if int64(expectedSize) != qSizeAfterPublishing {
		t.Errorf("actual queue size does not match expected: actual(%d), expected(%d)", qSizeAfterPublishing, expectedSize)
	}

	terminateSubscriber := make(chan bool, 1)

	dequeueHandler := func (subscriberId uuid.UUID, dequeued queue.QueueEntry[string]) bool{
		// deqMessage := fmt.Sprintf("Dequeued: %s on subscriber %s", dequeued, subscriberId.String())
		// fmt.Println(deqMessage)

		// let's simulate some work, just sleep for up to 5 microseconds (random)
		randNum := rand.Intn(5) + 1
		time.Sleep(time.Duration(randNum) * time.Microsecond)

		return true
	}

	for range make([]int, 10) {
		go func () {
				sOpts := queue.SubscriberOpts[string]{ 
				LFQueue: q, 
				DequeueHandler: dequeueHandler, 
				StackSize: 100, 
				TerminationSignal: terminateSubscriber,
			}

			subscriber := queue.NewSubscriber(sOpts)
			_, err := subscriber.Subscribe()
			if err != nil { fmt.Println(err) }
		}()
	}

	time.Sleep(1 * time.Second) // lets just consume for 1 second and see how many we can dequeue
	close(terminateSubscriber)
	
	t.Log("messages in queue after subscribing for 1 sec:", q.Size())
	t.Log("messages processed in 1 sec:", qSizeAfterPublishing - q.Size())

	expectedSizeAfterDequeue := 0
	qSizeAfterDequeuing := q.Size()
	if int64(expectedSizeAfterDequeue) != qSizeAfterDequeuing {
		t.Errorf("actual queue size does not match expected: actual(%d), expected(%d)", qSizeAfterPublishing, expectedSize)
	}
}