package lwwset

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLWWSetLogic(t *testing.T) {
	set1 := NewSet()

	// element1 for testing
	element1 := 1
	ts1 := int64(1)   // one ts
	ts12 := int64(12) // later ts
	ts13 := int64(13) // latest

	// element2 for testing
	element2 := 2
	ts2 := int64(2)
	ts22 := int64(22) // later ts

	// add element1.
	set1.Add(element1, ts1)
	assert.Equal(t, true, set1.Lookup(element1), "Should have element1")
	assert.Equal(t, false, set1.Lookup(element2), "Should not have element2")

	// the bias case, should the element should still exist
	set1.Remove(element1, ts1)
	assert.Equal(t, true, set1.Lookup(element1), "Should still have element1")

	// element1 should not exist anymore
	set1.Remove(element1, ts12)
	assert.Equal(t, false, set1.Lookup(element1), "Should not have element1")

	// create a set2
	set2 := NewSet()

	set2.Add(element2, ts2)
	assert.Equal(t, true, set2.Lookup(element2), "Should have element2")

	// set1 merge set2
	set1.Merge(set2)
	assert.Equal(t, false, set1.Lookup(element1), "Should not have element1")
	assert.Equal(t, true, set1.Lookup(element2), "Should have element2")

	// set2 add the latest element1, after merging, element1 should be back again to set1
	set2.Add(element1, ts13)
	set1.Merge(set2)
	assert.Equal(t, true, set1.Lookup(element1), "Should have element1 again")

	// remove the element from set2, after merging, element2 should not exist anymore
	set2.Remove(element2, ts22)
	set1.Merge(set2)
	assert.Equal(t, false, set1.Lookup(element2), "Should not have element2 any more")
}

func BenchmarkLWWSetAdd(b *testing.B) {
	var waitGroup sync.WaitGroup

	num := 0
	// banckmark with multi-thread
	set := NewSet()
	for i := 0; i < b.N; i++ {
		num++
		waitGroup.Add(1)
		go func(i int) {
			set.Add(i, time.Now().Unix())
			waitGroup.Done()
		}(i)
	}

	waitGroup.Wait()

	// timer stop here, no need to measure the following code
	b.StopTimer()

	// check the correctness
	for i := 0; i < num; i++ {
		if !set.Lookup(i) {
			b.Fatal("element not exists:", i)
		}
	}
}

func BenchmarkLWWSetLookup(b *testing.B) {
	var waitGroup sync.WaitGroup

	num := 0

	// Add some data.
	set := NewSet()
	for i := 0; i < b.N; i++ {
		num++
		waitGroup.Add(1)
		go func(i int) {
			set.Add(i, time.Now().Unix())
			waitGroup.Done()
		}(i)
	}

	waitGroup.Wait()

	// Just measuring the following lookup code
	b.ResetTimer()

	// banckmark with multi-thread
	correct := true
	for i := 0; i < num; i++ {
		waitGroup.Add(1)
		go func(i int) {
			if !set.Lookup(i) {
				correct = false
			}
			waitGroup.Done()
		}(i)
	}
	waitGroup.Wait()
	if !correct {
		b.Fatal("wrong result")
	}
}

func BenchmarkLWWSetRemove(b *testing.B) {
	var waitGroup sync.WaitGroup

	num := 0
	// To add some data.
	set := NewSet()
	for i := 0; i < b.N; i++ {
		num++
		waitGroup.Add(1)
		go func(i int) {
			set.Add(i, time.Now().Unix())
			waitGroup.Done()
		}(i)
	}
	waitGroup.Wait()

	// Just measuring the following remove method
	b.ResetTimer()

	// banckmark with multi-thread
	for i := 0; i < num; i++ {
		waitGroup.Add(1)
		go func(i int) {
			set.Remove(i, time.Now().Unix()+2)
			waitGroup.Done()
		}(i)
	}
	waitGroup.Wait()

	// timer stop here, no need to measure the following code
	b.StopTimer()

	// check the correctness
	for i := 0; i < num; i++ {
		if set.Lookup(i) {
			b.Fatal("element exists:", i)
		}
	}
}

func BenchmarkLWWSetMerge(b *testing.B) {
	// Prepare the merging sets
	sets := make([]*Set, 0)
	for i := 0; i < b.N; i++ {
		set := NewSet()
		set.Add(i, time.Now().Unix())
		sets = append(sets, set)
	}

	// Just measuring the following merge method
	b.ResetTimer()

	var waitGroup sync.WaitGroup
	set := NewSet()
	// banckmark with multi-thread
	for i := 0; i < len(sets); i++ {
		waitGroup.Add(1)
		go func(i int) {
			set.Merge(sets[i])
			waitGroup.Done()
		}(i)
	}
	waitGroup.Wait()

	// timer stop here, no need to measure the following code
	b.StopTimer()

	// check the correctness
	for i := 0; i < len(sets); i++ {
		if !set.Lookup(i) {
			b.Fatal("element not exists:", i)
		}
	}
}
