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


func TestMapRandomSmallConcurrentOperations(t *testing.T) {
	opts := lfmap.LFMapOpts{ PoolSize: 10000000 }
	lfMap := lfmap.NewLFMap[string](opts)

	inputSize := 1000000
	keyValPairs := make([]KeyVal, inputSize)

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

	// lfMap.PrintChildren()

	var retrieveWG sync.WaitGroup

	for _, val := range keyValPairs {
		retrieveWG.Add(1)
		go func (val KeyVal) {
			defer retrieveWG.Done()

			value := lfMap.Retrieve(val.Key)
			// t.Logf("actual: %s, expected: %s", value, val.Value)
			if value != val.Value {
				t.Errorf("actual value not equal to expected: actual(%s), expected(%s)", value, val.Value)
			}
		}(val)
	}

	retrieveWG.Wait()
}

func TestMapRandomLargeConcurrentOperations(t *testing.T) {
	opts := lfmap.LFMapOpts{ PoolSize: 10000000 }
	lfMap := lfmap.NewLFMap[string](opts)

	inputSize := 10000000

	keyValPairs := make([]KeyVal, inputSize)
	keyValChan := make(chan KeyVal, inputSize)
	
	var fillArrWG sync.WaitGroup

	for range keyValPairs {
		fillArrWG.Add(1)
		go func () {
			defer fillArrWG.Done()

			randomString, _ := GenerateRandomStringCrypto(32)
			keyValChan <- KeyVal{ Key: randomString, Value: randomString }
		}()
	}

	fillArrWG.Wait()
	t.Log("filled random key val pairs chan with size:", inputSize)

	for idx := range keyValPairs {
		keyVal :=<- keyValChan
		keyValPairs[idx] = keyVal
	}

	t.Log("seeded keyValPairs array:", inputSize)

	t.Log("inserting values -->")
	var insertWG sync.WaitGroup

	for _, val := range keyValPairs {
		insertWG.Add(1)
		go func (val KeyVal) {
			defer insertWG.Done()
			// we'll randomly publish between every 1 to 5 microseconds
			// randNum := mathRand.Intn(5) + 1
			//time.Sleep(time.Duration(randNum) * time.Microsecond) // simulate timeout
			
			lfMap.Insert(val.Key, val.Value)
		}(val)
	}

	insertWG.Wait()

	// lfMap.PrintChildren()

	t.Log("retrieving values -->")
	var retrieveWG sync.WaitGroup

	for _, val := range keyValPairs {
		retrieveWG.Add(1)
		go func (val KeyVal) {
			defer retrieveWG.Done()

			value := lfMap.Retrieve(val.Key)
			// t.Logf("actual: %s, expected: %s", value, val.Value)
			if value != val.Value {
				t.Errorf("actual value not equal to expected: actual(%s), expected(%s)", value, val.Value)
			}
		}(val)
	}

	retrieveWG.Wait()

	t.Log("done")
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