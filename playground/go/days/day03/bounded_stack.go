package day03

// Day 3 -- The bounded stack, shared.
//
// Day 2's bounded stack, with one change: many goroutines use it at the
// same time. The contract is the same:
//
//   - Push(v): only allowed when the stack is not full. If it is already
//     full, return ErrFull and do not change anything.
//   - Pop(): only allowed when the stack is not empty. If it is empty,
//     return ErrEmpty. Otherwise remove and return the top value.
//   - Top(): same rule as Pop(), but the value stays on the stack.
//   - Full(), Empty(), Size(): report the current state.
//
// Day 2 said: the invariant 0 <= Size() <= capacity may be broken for a
// moment *inside* a method, as long as it is restored before the method
// returns. With two goroutines, another goroutine can arrive in exactly
// that moment. No goroutine may ever see or cause a half-done push or
// pop: no value lost, no value returned twice, no value made up.

import "errors"

var ErrFull = errors.New("full")
var ErrEmpty = errors.New("empty")

type BoundedStack struct {
	// TODO: choose your own representation.
}

func NewBoundedStack(capacity int) *BoundedStack {
	// TODO
	return &BoundedStack{}
}

func (s *BoundedStack) Full() bool {
	// TODO
	return false
}

func (s *BoundedStack) Empty() bool {
	// TODO
	return false
}

func (s *BoundedStack) Size() int {
	// TODO
	return 0
}

func (s *BoundedStack) Push(v int) error {
	// TODO
	return nil
}

func (s *BoundedStack) Pop() (int, error) {
	// TODO
	return 0, nil
}

func (s *BoundedStack) Top() (int, error) {
	// TODO
	return 0, nil
}
