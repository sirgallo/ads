# Lock Free Queue


## Algorithm

`Michael & Scott Queue Algorithm`

This lock free queue algorithm implements atomic compare and swap (CAS) to ensure synchronization of elements entering and leaving the queue.

Below are explainations of both the enqueue and dequeue operations. The queue is composed of a singly linked list composed of elements containing a value and a next value, which is the link to the next element in the queue. The queue is implemented in memory and each next value is a memory address (or pointer) to the next node. New elements are added to the tail of the queue and elements ready to be dequeued are at the head.

Operating on the queue involves a couple of pieces, which include go routines for multiple threads to operate on the queue, and go routines to communicate with the queue:

The pieces:
  
1. ) The Queue
2. ) Publisher (Enqueue function) --> enqueues elements on the queue
4. ) Subscriber (Dequeue function) --> listens for new available elements and dequeues them into an internal call stack, and once the stack is full, processes all of the elements in the stack

For more information on the stack algorithm, check out [LockFreeStack](LockFreeStack.md)

Multiple Publishers and Subscribers can be created to operate on the queue.


### Enqueue

```
  1.) create a new node with value to be enqueued
  2.) continuously attempt enqueue by:
    I.) atomically read tail node pointer and next pointer of the tail node
    II.) check that the tail is still the same
    III.) if the next node is null, then the tail node is at the end of the queue --> CAS next node to new node
    IV.) if the next node is not null, then the tail node has been updated by a different thread --> CAS tail node to next node
```


### Dequeue

```
1.) continuously attempt dequeue by:
  I.) atomically read the head, tail, and next of (of the head node) pointers
  II.) check if the head and tail nodes are the same
  III.) if the head and tail nodes are the same next is null, then the queue is empty --> return null
  IV.) if the head pointer is the same, it means the next pointer is a valid node --> CAS head and next node, return next value
  V.) if the head node is not the same, then another thread modified it --> CAS head to next node
```


## Additional Notes


### Optimistic Strategies on Contention

When contention is high, an `exponential backoff` strategy is applied, where, if the CAS operation fails, the thread operating on the queue will pause for longer and longer periods of time until it is able to fulfill the operation, or fails after a max number of attempts

```go
// 2 ^ (depth - 1) * timeout
func ExpBackoffStrat(depth int, timeout int) (int, int) {
	time.Sleep(time.Duration(timeout) * time.Microsecond)

	timeout = int(math.Pow(float64(2), float64(depth - 1))) * timeout
	depth++

	return depth, timeout
}
```


## Sources

[LockFreeQueue](../src/queue/LFQueue.go)

[Publisher](../src/queue/Publisher.go)

[Subscriber](../src/queue/Subscriber.go)