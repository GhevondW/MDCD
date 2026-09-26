package day03

// Day 3 -- Challenge: Compute Once per Key.
//
// Day 1's Singleton had a line to remember:
//
//	if instance == nil { instance = newSingleton() }
//
// Two goroutines can both see nil, and both create an instance. This
// task is the same problem for many keys: a cache where the value for
// each key is computed at most once.
//
//	Get(key, compute)  if key already has a value, return it.
//	                   Otherwise call compute() to make one, remember
//	                   it, and return it.
//	Size()             how many keys have a remembered value. A
//	                   computation that is still running does not
//	                   count yet.
//
// The rules that make it hard:
//
//  1. compute() runs AT MOST ONCE per key -- even when many goroutines
//     ask for the same missing key at the same moment. One of them runs
//     compute(); the others wait for it and return the same value.
//  2. Different keys never wait for each other: while compute() runs
//     for key "a", a Get() for key "b" can run its own compute() at the
//     same time. So you cannot hold one lock while compute() runs.
//  3. If compute() returns an error, nothing is remembered for that key,
//     and Get returns that error (with the zero value of V) to the
//     caller that ran compute(). Goroutines that were waiting for that
//     same computation get the same error. The next Get() for that key
//     calls compute() again.
//
// Every method may be called by many goroutines at the same time.

type OnceCache[V any] struct {
	// TODO: choose your own representation.
}

func NewOnceCache[V any]() *OnceCache[V] {
	// TODO
	return &OnceCache[V]{}
}

func (c *OnceCache[V]) Get(key string, compute func() (V, error)) (V, error) {
	// TODO
	var zero V
	return zero, nil
}

func (c *OnceCache[V]) Size() int {
	// TODO
	return 0
}
