# Lock Free Stack


## Algorithm

This lock free stack algorithm implements atomic compare and swap (CAS) to ensure synchronization of elements entering and leaving the stack.

Below are explainations of both the push and pop operations. The stack is composed of a singly linked list composed of elements containing a value and a next value, which is the link to the next element in the stack. The stack is implemented in memory and each next value is a memory address (or pointer) to the next node. New elements are added to the top of the stack.


### Push

```
  1.) create a new node with value to be pushed on the stack
  2.) continuously attempt enqueue by:
    I.) If stack has reached max size, return
    II.) atomically read top node pointer and next pointer of the top node
    II.) check that the tail is still the same
    III.) CAS new node with top and if success, increment stack length and return true
    IV.) CAS new node with top and if failure --> return to 2.)
```


### Pop

```
1.) continuously attempt dequeue by:
  I.) atomically read the top node
  II.) if the top node is null, then the stack is empty --> return null
  III.) get next top value from next pointer
  IV.) if CAS top and next node succeeds, return top value and decremented stack length
  V.) if CAS top and next node fails, then another thread modified it --> return to 1.)
```


## Additional Nodes


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


## Source

[LockFreeStack](../pkg/stack/LFStack.go)