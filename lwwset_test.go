package lwwset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLWWSetMethods(t *testing.T) {
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
	assert.Equal(t, false, set1.Lookup(element2), "Should have element1")
}
