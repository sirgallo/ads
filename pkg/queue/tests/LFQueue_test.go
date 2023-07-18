package lfqueuetests

import "testing"

import "github.com/sirgallo/ads/pkg/queue"
import "github.com/sirgallo/ads/pkg/utils"


func TestQueueOperations(t *testing.T) {
	maxRetries := 10
  expBackoffOpts := utils.ExpBackoffOpts{ MaxRetries: &maxRetries, TimeoutInMicroseconds: 1 }
  qOpts := lfqueue.LFQueueOpts{ MaxPoolSize: 10000, ExpBackoffOpts: expBackoffOpts }

	q := lfqueue.NewLFQueue[string](qOpts)

	enqueueVals := []string{"hi", "random", "hello!", "new"}
	
	t.Log("enqueuing test values")
	for _, val := range enqueueVals {
		_, err := q.Enqueue(val)
		if err != nil {
			t.Error("error enqueuing values into queue")
		}
	}

	t.Log("test peek")
	val := q.Peek()
	t.Logf("expected: %s, actual: %s", enqueueVals[0], val)
	if val != enqueueVals[0] {
		t.Error("head of queue not the expected value")
	}

	t.Log("test dequeue")
	for idx := range enqueueVals {
		val, err := q.Dequeue()
		if err != nil {
			t.Error("error dequeuing element")
		}

		t.Logf("actual: %s, expected: %s", val, enqueueVals[idx])
		if val != enqueueVals[idx] {
			t.Errorf("actual value not equal to expected: actual(%s), expected(%s)", val, enqueueVals[idx])
		}
	}
}

func TestClear(t *testing.T) {
	maxRetries := 10
  expBackoffOpts := utils.ExpBackoffOpts{ MaxRetries: &maxRetries, TimeoutInMicroseconds: 1 }
  qOpts := lfqueue.LFQueueOpts{ MaxPoolSize: 10000, ExpBackoffOpts: expBackoffOpts }

	q := lfqueue.NewLFQueue[string](qOpts)

	enqueueVals := []string{"hi", "random", "hello!", "new"}
	
	for _, val := range enqueueVals {
		q.Enqueue(val)
	}

	q.Clear()

	expectedSize := 0
	actualSize := q.Size()

	if int64(expectedSize) != actualSize {
		t.Errorf("actual queue size does not match expected: actual(%d), expected(%d)", actualSize, expectedSize)
	}
}

func TestEmpty(t *testing.T) {
	maxRetries := 10
  expBackoffOpts := utils.ExpBackoffOpts{ MaxRetries: &maxRetries, TimeoutInMicroseconds: 1 }
  qOpts := lfqueue.LFQueueOpts{ MaxPoolSize: 10000, ExpBackoffOpts: expBackoffOpts }

	q := lfqueue.NewLFQueue[string](qOpts)

	nilVal := utils.GetZero[string]()
	val, err := q.Dequeue()

	if err != nil {
		t.Error("error dequeuing element from queue")
	}

	if val != nilVal {
		t.Error("value is not equal to null value")
	}
}