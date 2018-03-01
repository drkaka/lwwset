package lwwset

import "sync"

// Set for LWW-Element
type Set struct {
	l         sync.RWMutex
	addSet    map[interface{}]int64
	removeSet map[interface{}]int64
}

// NewSet create a new LWW-Element-Set
func NewSet() *Set {
	return &Set{
		addSet:    make(map[interface{}]int64, 0),
		removeSet: make(map[interface{}]int64, 0),
	}
}

// Add an element with its unix timestamp.
func (s *Set) Add(element interface{}, ts int64) {
	// add the element to the addSet
	s.l.Lock()
	s.addSet[element] = ts
	s.l.Unlock()
}

// Lookup an element.
func (s *Set) Lookup(element interface{}) bool {
	// get the element info in both add/remove set
	s.l.RLock()
	tsAdd := s.addSet[element]
	tsRemove := s.removeSet[element]
	s.l.RUnlock()

	// if the element exists in addSet and the timestamp is greater or equal to the timestamp in removeSet, return true
	if tsAdd > 0 && tsAdd >= tsRemove {
		return true
	}

	return false
}

// Remove an element with its unix timestamp
func (s *Set) Remove(element interface{}, ts int64) {
	if s.Lookup(element) {
		// The element can be looked up, add the element to the removeSet
		s.l.Lock()
		s.removeSet[element] = ts
		s.l.Unlock()
	}
}

// Merge the set t to set s
func (s *Set) Merge(t *Set) {
	t.l.RLock()
	defer t.l.RUnlock()

	// merging the addSet of t to s with keeping the latest timestamp
	for element, ts := range t.addSet {

		// get the element and its ts in addSet of s
		s.l.RLock()
		tsInS, inS := s.addSet[element]
		s.l.RUnlock()

		if !inS || (inS && tsInS < ts) {
			// if not in s or in s but the ts is earlier, add the element to addSet of s
			s.l.Lock()
			s.addSet[element] = ts
			s.l.Unlock()
		}
	}

	// merging the removeSet of t to s with keeping the latest timestamp
	for element, ts := range t.removeSet {

		// get the element and its ts in removeSet of s
		s.l.RLock()
		tsInS, inS := s.removeSet[element]
		s.l.RUnlock()

		if !inS || (inS && tsInS < ts) {
			// if not in s or in s but the ts is earlier, add the element to removeSet of s
			s.l.Lock()
			s.removeSet[element] = ts
			s.l.Unlock()
		}
	}
}
