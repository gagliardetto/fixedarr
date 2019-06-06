// Package fixedarr can be used when you want an array of elements that never
// expands over a certain limit; if the array reached its limit capacity, and a
// new element is pushed to it, then the oldest element is removed.
package fixedarr

import "sync"

// Array is a fixed size array; is the current size reached max,
// old elements will be dropped when new elements are added.
type Array struct {
	mu         *sync.RWMutex
	array      []interface{}
	maxSize    int
	atCapacity bool
}

// New returns a new Array; maxSize MUST be a positive number.
func New(maxSize int) *Array {
	if maxSize < 0 {
		panic("fixedarr.New: maxSize cannot be less than 0")
	}
	return &Array{
		mu:      &sync.RWMutex{},
		array:   make([]interface{}, 0),
		maxSize: maxSize,
	}
}

// Push pushes (appends) an element to the array; if the array has reached
// its limit capacity, the oldest element will be removed.
func (a *Array) Push(el interface{}) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.atCapacity || len(a.array) >= a.maxSize && len(a.array) > 0 {

		if !a.atCapacity {
			a.atCapacity = true
		}
		i := 0
		copy(a.array[i:], a.array[i+1:])
		a.array[len(a.array)-1] = nil
		a.array = a.array[:len(a.array)-1]

	}

	a.array = append(a.array, el)
}

// Len returns the current length of the array
func (a *Array) Len() int {
	return len(a.Value())
}

// Max returns the limit size of the array
func (a *Array) Max() int {
	return a.maxSize
}

// Value returns the current array
func (a *Array) Value() []interface{} {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.array
}

// Reset resets the array
func (a *Array) Reset() {
	a.mu.RLock()
	defer a.mu.RUnlock()

	a.array = make([]interface{}, 0)
}

// GetAndReset returns the current array, and resets it
func (a *Array) GetAndReset() []interface{} {
	a.mu.RLock()
	defer a.mu.RUnlock()

	clone := make([]interface{}, 0)
	for i := range a.array {
		clone = append(clone, a.array[i])
	}

	a.array = make([]interface{}, 0)

	return clone
}
