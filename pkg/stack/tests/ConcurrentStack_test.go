package lfstacktests

import "sync"
import "testing"

import "github.com/sirgallo/ads/pkg/stack"
import "github.com/sirgallo/ads/pkg/utils"


func TestStackConcurrentOperations(t *testing.T) {
	maxRetries := 10
  sExpBackoffOpts := utils.ExpBackoffOpts{ MaxRetries: &maxRetries, TimeoutInMicroseconds: 10 }
  sOpts := stack.LFStackOpts{ MaxStackSize: 10000, ExpBackoffOpts: sExpBackoffOpts }
	
  lfStack := stack.NewLFStack[string](sOpts)

	pushVals := []string{"hi", "random", "hello!", "new"}

	var pushWG sync.WaitGroup

	for _, val := range pushVals {
		pushWG.Add(1)
		go func (val string) {
			defer pushWG.Done()
			
			t.Log("val pushed:", val)
			lfStack.Push(val)
		}(val)
	}

	pushWG.Wait()

	var popWG sync.WaitGroup

	for range pushVals {
		popWG.Add(1)
		go func () {
			defer popWG.Done()

			val, err := lfStack.Pop()
			t.Log("val popped:", val)
			if err != nil {
				t.Error("error popping value")
			}
		}()
	}

	popWG.Wait()

	expectedSize := 0
	stackSize := lfStack.Size()
	if int64(expectedSize) != stackSize {
		t.Errorf("stack not expected size: actual(%d), expected(%d)", stackSize, expectedSize)
	}
}