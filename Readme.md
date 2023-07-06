# Atomic Data Structures


A collection of atomic data structures meant to be used in concurrent environments


## Installation


Import this repository to utilize different atomic datastructure, implemented in `Go`.

in your `Go` project main directory (where the `go.mod` file is located)
```bash
go get github.com/sirgallo/ads
go mod tidy
```

make sure to run go mod tidy to install dependencies


## Data Structures


[LockFreeQueue](./docs/LockFreeQueue.md)

to use:
```go
import "github.com/sirgallo/ads/queue"

func main() {
  // define max queue size and exponential backoff on CAS failure
  maxRetries := 10
  expBackoffOpts := utils.ExpBackoffOpts{ MaxRetries: &maxRetries, TimeoutInMicroseconds: 1 }
  qOpts := queue.LFQueueOpts{ MaxQueueSize: 10000, ExpBackoffOpts: expBackoffOpts }
	
  // intantiate queue with opts
  q := queue.NewLFQueue(qOpts)

  // enqueue
  q.Enqueue("hi")

  // dequeue
  val, err := q.Dequeue()
  if err != nil { // handle error }
}
```


[LockFreeStack](./docs/LockFreeStack.md)

to use:
```go
import "github.com/sirgallo/ads/stack"

func main() {
  // define max stack size and exponential backoff on CAS failure
  maxRetries := 10
  sExpBackoffOpts := utils.ExpBackoffOpts{ MaxRetries: &maxRetries, TimeoutInMicroseconds: 10 }
  sOpts := stack.LFStackOpts{ MaxStackSize: 10000, ExpBackoffOpts: sExpBackoffOpts }
	
  // intantiate stack with opts
  lfStack := stack.NewLFStack(sOpts)

  // push
  lfStack.Push("hi")

  // pop
  val, err := lfStack.Pop()
  if err != nil { // handle error }
```

[Counter](./counter/Counter.go)

to use:
```go
import "github.com/sirgallo/ads/counter"

func main() {
  // intantiate the counter
  counter, _ := counter.NewCounter(0)

  // add to the counter
  counter.Increment(1)

  // get the current value
  val := counter.getValue()

  // subtract from the counter
  counter.Decrement(1)
}
```
