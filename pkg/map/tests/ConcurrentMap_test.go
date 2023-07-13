package lfmaptests

import "crypto/rand"
import "encoding/base64"
import "sync"
import "testing"

import "github.com/sirgallo/ads/pkg/map"


type KeyVal struct {
	Key string
	Value string
}


func TestMapConcurrentOperations(t *testing.T) {
	lfMap := lfmap.NewLFMap[string]()

	keyValPairs := make([]KeyVal, 10000)
	for idx := range keyValPairs {
		randomString, _ := GenerateRandomStringCrypto(32)
		keyValPairs[idx] = KeyVal{ Key: randomString, Value: randomString }
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

	lfMap.PrintChildren()

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

func GenerateRandomStringCrypto(length int) (string, error) {
	randomBytes := make([]byte, length)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	randomString := base64.RawURLEncoding.EncodeToString(randomBytes)
	return randomString[:length], nil
}