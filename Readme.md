# Atomic Data Structures

A collection of atomic data structures meant to be used in concurrent environments. These aim to be thread safe and non blocking, utilizing atomic operations like Compare-and-Swap to update the following data structures.


## Installation

In your `Go` project main directory (where the `go.mod` file is located)
```bash
go get github.com/sirgallo/ads
go mod tidy
```

Make sure to run go mod tidy to install dependencies.


## Data Structures

[LockFreeMap](./docs/LockFreeMap.md)

to use:
```go
package main

import "github.com/sirgallo/ads/pkg/map"

func main() {
  // initialize lock free map

  // for 32 bit hash and bitmap
  lfMap := lfmap.NewLFMap[string, uint32]()

  // for 64 bit hash and bitmap
  lfMap := lfmap.NewLFMap[string, uint64]()

  // insert key/val pair
  lfMap.Insert("hi", "world")

  // retrieve value for key
  val := lfMap.Retrieve("hi")

  // delete key/val pair
  lfMap.Delete("hi")
}
```

to test:
```bash
go test -v ./pkg/map/tests
```


[LockFreeQueue](./docs/LockFreeQueue.md)

to use:
```go
package main

import "github.com/sirgallo/ads/pkg/queue"
import "github.com/sirgallo/ads/pkg/utils"

func main() {
  // define max queue size and exponential backoff on CAS failure
  maxRetries := 10
  expBackoffOpts := utils.ExpBackoffOpts{ MaxRetries: &maxRetries, TimeoutInMicroseconds: 1 }
  qOpts := lfqueue.LFQueueOpts{ MaxQueueSize: 10000, ExpBackoffOpts: expBackoffOpts }
	
  // initialize queue with opts
  q := lfqueue.NewLFQueue[string](qOpts)

  // enqueue
  q.Enqueue("hi")

  // dequeue
  val, err := q.Dequeue()
  if err != nil { // handle error }
}
```

to test:
```bash
go test -v ./pkg/queue/tests
```


[LockFreeStack](./docs/LockFreeStack.md)

to use:
```go
package main

import "github.com/sirgallo/ads/pkg/stack"
import "github.com/sirgallo/ads/pkg/utils"

func main() {
  // define max stack size and exponential backoff on CAS failure
  maxRetries := 10
  sExpBackoffOpts := utils.ExpBackoffOpts{ MaxRetries: &maxRetries, TimeoutInMicroseconds: 10 }
  sOpts := lfstack.LFStackOpts{ MaxStackSize: 10000, ExpBackoffOpts: sExpBackoffOpts }
	
  // initialize stack with opts
  lfStack := lfstack.NewLFStack[string](sOpts)

  // push
  lfStack.Push("hi")

  // pop
  val, err := lfStack.Pop()
  if err != nil { // handle error }
}
```

to test:
```bash
go test -v ./pkg/stack/tests
```


[Counter](./pkg/counter/Counter.go)

to use:
```go
package main

import "github.com/sirgallo/ads/pkg/counter"

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

to test:
```bash
go test -v ./pkg/counter/tests
```