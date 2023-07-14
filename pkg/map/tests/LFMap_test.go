package lfmaptests

import "testing"
import "sync/atomic"

import "github.com/sirgallo/ads/pkg/map"


//=================================== 32 bit

func TestMapOperations32(t *testing.T) {
	opts := lfmap.LFMapOpts{ PoolSize: 10000000 }
	lfMap := lfmap.NewLFMap[string, uint32](opts)

	lfMap.Insert("hello", "world")
	lfMap.Insert("new", "wow!")
	lfMap.Insert("again", "test!")
	lfMap.Insert("woah", "random entry")
	lfMap.Insert("key", "Saturday!")
	lfMap.Insert("sup", "6")
	lfMap.Insert("final", "the!")
	lfMap.Insert("6", "wow!")
	lfMap.Insert("asdfasdf", "add 10")
	lfMap.Insert("asdfasdf", "123123") // note same key, will update value
	lfMap.Insert("asd", "queue!")
	lfMap.Insert("fasdf", "interesting")
	lfMap.Insert("yup", "random again!")
	lfMap.Insert("asdf", "hello")
	lfMap.Insert("asdffasd", "uh oh!")
	lfMap.Insert("fasdfasdfasdfasdf", "error message")
	lfMap.Insert("fasdfasdf", "info!")
	lfMap.Insert("woah", "done")

	rootBitMap := (*lfmap.LFMapNode[string, uint32])(atomic.LoadPointer(&lfMap.Root)).BitMap

	t.Logf("lfMap after inserts")
	lfMap.PrintChildren()

	expectedBitMap := uint32(542198999)
	t.Logf("actual root bitmap: %d, expected root bitmap: %d\n", rootBitMap, expectedBitMap)
	t.Logf("actual root bitmap: %032b, expected root bitmap: %032b\n", rootBitMap, expectedBitMap)
	if expectedBitMap != rootBitMap {
		t.Errorf("actual bitmap does not match expected bitmap: actual(%032b), expected(%032b)\n", rootBitMap, expectedBitMap)
	}

	t.Log("retrieve values")

	val1 := lfMap.Retrieve("hello")
	expVal1 :=  "world"
	t.Logf("actual: %s, expected: %s", val1, expVal1)
	if val1 != expVal1 {
		t.Errorf("val 1 does not match expected val 1: actual(%s), expected(%s)\n", val1, expVal1)
	}

	val2 := lfMap.Retrieve("new")
	expVal2 :=  "wow!"
	t.Logf("actual: %s, expected: %s", val2, expVal2)
	if val2 != expVal2 {
		t.Errorf("val 1 does not match expected val 1: actual(%s), expected(%s)\n", val2, expVal2)
	}

	val3 := lfMap.Retrieve("asdf")
	expVal3 := "hello"
	t.Logf("actual: %s, expected: %s", val3, expVal3)
	if val3 != expVal3 {
		t.Errorf("val 1 does not match expected val 1: actual(%s), expected(%s)", val3, expVal3)
	}

	val4 := lfMap.Retrieve("asdfasdf")
	expVal4 := "123123"
	t.Logf("actual: %s, expected: %s", val4, expVal4)
	if val4 != expVal4 {
		t.Errorf("val 1 does not match expected val 1: actual(%s), expected(%s)", val4, expVal4)
	}

	lfMap.Delete("hello")
	lfMap.Delete("yup")
	lfMap.Delete("asdf")
	lfMap.Delete("asdfasdf")
	lfMap.Delete("new")
	lfMap.Delete("6")

	rootBitMapAfterDelete := (*lfmap.LFMapNode[string, uint32])(atomic.LoadPointer(&lfMap.Root)).BitMap
	t.Logf("bitmap of root after deletes: %032b\n", rootBitMapAfterDelete)
	t.Logf("bitmap of root after deletes: %d\n", rootBitMapAfterDelete)

	t.Log("hamt after deletes")
	lfMap.PrintChildren()

	expectedRootBitmapAfterDelete := uint32(536956102)
	t.Log("actual bitmap:", rootBitMapAfterDelete, "expected bitmap:", expectedRootBitmapAfterDelete)
	if expectedRootBitmapAfterDelete != rootBitMapAfterDelete {
		t.Errorf("actual bitmap does not match expected bitmap: actual(%032b), expected(%032b)\n", rootBitMapAfterDelete, expectedRootBitmapAfterDelete)
	}
}


//=================================== 64 bit

func TestMapOperations64(t *testing.T) {
	opts := lfmap.LFMapOpts{ PoolSize: 10000000 }
	lfMap := lfmap.NewLFMap[string, uint64](opts)

	lfMap.Insert("hello", "world")
	lfMap.Insert("new", "wow!")
	lfMap.Insert("again", "test!")
	lfMap.Insert("woah", "random entry")
	lfMap.Insert("key", "Saturday!")
	lfMap.Insert("sup", "6")
	lfMap.Insert("final", "the!")
	lfMap.Insert("6", "wow!")
	lfMap.Insert("asdfasdf", "add 10")
	lfMap.Insert("asdfasdf", "123123") // note same key, will update value
	lfMap.Insert("asd", "queue!")
	lfMap.Insert("fasdf", "interesting")
	lfMap.Insert("yup", "random again!")
	lfMap.Insert("asdf", "hello")
	lfMap.Insert("asdffasd", "uh oh!")
	lfMap.Insert("fasdfasdfasdfasdf", "error message")
	lfMap.Insert("fasdfasdf", "info!")
	lfMap.Insert("woah", "done")

	rootBitMap := (*lfmap.LFMapNode[string, uint64])(atomic.LoadPointer(&lfMap.Root)).BitMap

	t.Logf("lfMap after inserts")
	lfMap.PrintChildren()

	expectedBitMap := uint64(18084858599620633)
	t.Logf("actual root bitmap: %d, expected root bitmap: %d\n", rootBitMap, expectedBitMap)
	t.Logf("actual root bitmap: %032b, expected root bitmap: %032b\n", rootBitMap, expectedBitMap)
	if expectedBitMap != rootBitMap {
		t.Errorf("actual bitmap does not match expected bitmap: actual(%032b), expected(%032b)\n", rootBitMap, expectedBitMap)
	}

	t.Log("retrieve values")

	val1 := lfMap.Retrieve("hello")
	expVal1 :=  "world"
	t.Logf("actual: %s, expected: %s", val1, expVal1)
	if val1 != expVal1 {
		t.Errorf("val 1 does not match expected val 1: actual(%s), expected(%s)\n", val1, expVal1)
	}

	val2 := lfMap.Retrieve("new")
	expVal2 :=  "wow!"
	t.Logf("actual: %s, expected: %s", val2, expVal2)
	if val2 != expVal2 {
		t.Errorf("val 1 does not match expected val 1: actual(%s), expected(%s)\n", val2, expVal2)
	}

	val3 := lfMap.Retrieve("asdf")
	expVal3 := "hello"
	t.Logf("actual: %s, expected: %s", val3, expVal3)
	if val3 != expVal3 {
		t.Errorf("val 1 does not match expected val 1: actual(%s), expected(%s)", val3, expVal3)
	}

	val4 := lfMap.Retrieve("asdfasdf")
	expVal4 := "123123"
	t.Logf("actual: %s, expected: %s", val4, expVal4)
	if val4 != expVal4 {
		t.Errorf("val 1 does not match expected val 1: actual(%s), expected(%s)", val4, expVal4)
	}

	lfMap.Delete("hello")
	lfMap.Delete("yup")
	lfMap.Delete("asdf")
	lfMap.Delete("asdfasdf")
	lfMap.Delete("new")
	lfMap.Delete("6")

	rootBitMapAfterDelete := (*lfmap.LFMapNode[string, uint64])(atomic.LoadPointer(&lfMap.Root)).BitMap
	t.Logf("bitmap of root after deletes: %032b\n", rootBitMapAfterDelete)
	t.Logf("bitmap of root after deletes: %d\n", rootBitMapAfterDelete)

	t.Log("hamt after deletes")
	lfMap.PrintChildren()

	expectedRootBitmapAfterDelete := uint64(18014472667152401)
	t.Log("actual bitmap:", rootBitMapAfterDelete, "expected bitmap:", expectedRootBitmapAfterDelete)
	if expectedRootBitmapAfterDelete != rootBitMapAfterDelete {
		t.Errorf("actual bitmap does not match expected bitmap: actual(%032b), expected(%032b)\n", rootBitMapAfterDelete, expectedRootBitmapAfterDelete)
	}
}