package counter

import "errors"
import "sync/atomic"


func NewCounter(initialValue int64) (*Counter, error) {
	if initialValue < 0 { return nil, errors.New("initial value must be positive") }
	return &Counter{ value: 0 }, nil
}

func (counter *Counter) Increment(step int64) (int64, error) {
	incremented := atomic.AddInt64(&counter.value, step)
	if incremented < 0 { return counter.value, errors.New("increment returned negative value") }
	
	return incremented, nil
}

func (counter *Counter) Decrement(step int64) (int64, error) {
	decremented := atomic.AddInt64(&counter.value, -step)
	if decremented < 0 { return counter.value, errors.New("decrement returned negative value") }
	
	return decremented, nil
}

func (counter *Counter) GetValue() int64 {
	return atomic.LoadInt64(&counter.value)
}