# Pub Sub Example

This is a simple publisher-subscriber using multiple go routines on the lock free queue as a test.

```go
package main

import "fmt"
import "math/rand"
import "runtime"
import "sync"
import "time"
import "github.com/google/uuid"

import "github.com/sirgallo/ads/pkg/queue"
import "github.com/sirgallo/ads/pkg/utils"


func main() {
	numCpu := runtime.NumCPU()

	fmt.Println("number of cpus", numCpu)
	fmt.Println("let's test how many messages subscribers can process on go routines ranging from 1 - 10")
	fmt.Println("messages in queue after publishing at each pass: 20000000")

	for idx := range make([]int, numCpu) {
		goRoutineCount := idx + 1
		fmt.Println("\ngo routines for subscribers:", goRoutineCount)

		maxRetries := 10
		expBackoffOpts := utils.ExpBackoffOpts{ MaxRetries: &maxRetries, TimeoutInMicroseconds: 1 }
		
		qOpts := queue.LFQueueOpts{ MaxQueueSize: 100000000, ExpBackoffOpts: expBackoffOpts }
		q := queue.NewLFQueue[queue.QueueEntry[string]](qOpts)

		var publisherWG sync.WaitGroup
		
		for range make([]int, 10) {
			publisherWG.Add(1)
			
			go func () bool {
				defer publisherWG.Done()
				
				pOpts := queue.PublisherOpts[string]{ LFQueue: q }
				publisher := queue.NewPublisher(pOpts)
				messageNumber := 200000

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
		qSizeAfterPublishing := q.Size()

		terminateSubscriber := make(chan bool, 1)

		dequeueHandler := func (subscriberId uuid.UUID, dequeued queue.QueueEntry[string]) bool{
			// deqMessage := fmt.Sprintf("Dequeued: %s on subscriber %s", dequeued, subscriberId.String())
			// fmt.Println(deqMessage)

			// let's simulate some work, just sleep for up to 5 microseconds (random)
			randNum := rand.Intn(5) + 1
			time.Sleep(time.Duration(randNum) * time.Microsecond)

			return true
		}

		for range make([]int, goRoutineCount) {
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
		
		fmt.Println("messages in queue after subscribing for 1 sec:", q.Size())
		fmt.Println("messages processed in 1 sec:", qSizeAfterPublishing - q.Size())
	}
}
```