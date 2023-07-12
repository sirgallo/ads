package lfstacktests

import "testing"

import "github.com/sirgallo/ads/pkg/stack"
import "github.com/sirgallo/ads/pkg/utils"


func TestStackOperation(t *testing.T) {
	maxRetries := 10
  sExpBackoffOpts := utils.ExpBackoffOpts{ MaxRetries: &maxRetries, TimeoutInMicroseconds: 10 }
  sOpts := stack.LFStackOpts{ MaxStackSize: 10000, ExpBackoffOpts: sExpBackoffOpts }
	
  lfStack := stack.NewLFStack[string](sOpts)

	pushVals := []string{"hi", "random", "hello!", "new"}
	popVals := []string{"new", "hello!", "random", "hi"}
	t.Log("pushVals:", pushVals)
	t.Log("popVals:", popVals)

	t.Log("pushing test values")
	for _, val := range pushVals {
		_, err := lfStack.Push(val)
		if err != nil {
			t.Error("error enqueuing values into queue")
		}
	}

	t.Log("test peek")
	val := lfStack.Peek()
	t.Logf("actual: %s, expected: %s", val, popVals[0])
	if val != popVals[0] {
		t.Error("top of stack not the expected value")
	}

	t.Log("test pop")
	for idx := range popVals {
		val, err := lfStack.Pop()
		if err != nil {
			t.Error("error dequeuing element")
		}

		t.Logf("actual: %s, expected: %s", val, popVals[idx])
		if val != popVals[idx] {
			t.Errorf("actual value not equal to expected: actual(%s), expected(%s)", val, popVals[idx])
		}
	}
}