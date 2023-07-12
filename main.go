package main

import "fmt"
// import "github.com/sirgallo/ads/pkg/hamt"


import "github.com/sirgallo/ads/pkg/map"
// import "github.com/sirgallo/ads/pkg/utils"


func main() {
	/*
	maxRetries := 5
	opts := lfmap.LFMapOpts{
		ExpBackoffOpts: utils.ExpBackoffOpts{ MaxRetries: &maxRetries, TimeoutInMicroseconds: 1 },
		MaxPoolSize: 1000000,
	}
	*/
	hamt := lfmap.NewLFMap()

	// maxRetries := 10
  // expBackoffOpts := utils.ExpBackoffOpts{ MaxRetries: &maxRetries, TimeoutInMicroseconds: 1 }
  // tOpts := lftrie.LFTrieOpts{ BitChunkSize: 5, MaxPoolSize: 10000, ExpBackoffOpts: expBackoffOpts }
	// hamt := lftrie.NewLFTrie(tOpts)
	
	fmt.Println("insert values")
	hamt.Insert("key", "Saturday!")
	hamt.Insert("hi", "world")
	hamt.Insert("new", "wow!")
	hamt.Insert("again", "test!")
	hamt.Insert("woah", "random entry")
	hamt.Insert("sup", "6")
	hamt.Insert("final", "the!")
	hamt.Insert("6", "wow!")
	hamt.Insert("asdfasdf", "add 10")
	hamt.Insert("asdfasdf", "hi")
	hamt.Insert("asd", "queue!")
	hamt.Insert("fasdf", "interesting")
	hamt.Insert("yup", "random again!")
	hamt.Insert("asdf", "hello")
	hamt.Insert("asdffasd", "uh oh!")
	hamt.Insert("fasdfasdfasdfasdf", "error message")
	hamt.Insert("fasdfasdf", "info!")
	hamt.Insert("woah", "done")
	
	fmt.Println("retrieve values")
	val1 := hamt.Retrieve("hi")
	fmt.Println("val", val1)
	val2 := hamt.Retrieve("new")
	fmt.Println("val", val2)
	val3 := hamt.Retrieve("key")
	fmt.Println("val", val3)
	
	fmt.Println("print all children")
	hamt.PrintChildren()

	fmt.Println("delete key hi")
	hamt.Delete("hi")
	fmt.Println("delete key yup")
	hamt.Delete("yup")
	fmt.Println("delete key asdf")
	hamt.Delete("asdf")
	fmt.Println("delete key asdffasd")
	hamt.Delete("asdffasd")
	fmt.Println("delete key fasdfasdfasdfasdf")
	hamt.Delete("fasdfasdfasdfasdf")

	fmt.Println("print all children after delete")
	hamt.PrintChildren()
}