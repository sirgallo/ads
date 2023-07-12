package lfmaptests

import "sync"
import "testing"

import "github.com/sirgallo/ads/pkg/map"


type KeyVal struct {
	Key string
	Value string
}
func TestMapConcurrentOperations(t *testing.T) {
	lfMap := lfmap.NewLFMap[string]()

	keyValPairs := []KeyVal{
		{ Key: "hello", Value: "world" },
		{ Key: "new", Value: "wow!" },
		{ Key: "again", Value: "test!" },
		{ Key: "asdf", Value: "hello" },
		{ Key: "key", Value: "Saturday!" },
	}

	var insertWG sync.WaitGroup

	for _, val := range keyValPairs {
		insertWG.Add(1)
		go func (val KeyVal) {
			defer insertWG.Done()

			lfMap.Insert(val.Key, val.Value)
		}(val)
	}

	insertWG.Wait()

	var retrieveWG sync.WaitGroup

	for _, val := range keyValPairs {
		retrieveWG.Add(1)
		go func (val KeyVal) {
			defer retrieveWG.Done()

			value := lfMap.Retrieve(val.Key)
			t.Logf("actual: %s, expected: %s", value, val.Value)
			if value != val.Value {
				t.Errorf("actual value not equal to expected: actual(%s), expected(%s)", value, val.Value)
			}
		}(val)
	}

	retrieveWG.Wait()
}