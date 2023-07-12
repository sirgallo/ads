package countertests

import "testing"

import "github.com/sirgallo/ads/pkg/counter"


func TestCounterOperations(t *testing.T) {
	counter, _ := counter.NewCounter(0)
  
	counter.Increment(1)

  valIncr := counter.GetValue()
	t.Log("val after incr:", valIncr)

	if valIncr != 1 {
		t.Error("val set incorrectly")
	}

  counter.Decrement(1)

	valDec := counter.GetValue()
	t.Log("val after decr", valDec)

	if valDec != 0 {
		t.Error("val set incorrectly")
	}
}