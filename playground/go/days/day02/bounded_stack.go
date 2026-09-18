package day02

// Day 2 -- Contracts & Invariants.
//
// A stack of ints with a fixed capacity, given once to NewBoundedStack.
// The stack must protect its own rules:
//
//   - Push(v): only allowed when the stack is not full. If it is already
//     full, return ErrFull and do not change anything.
//   - Pop(): only allowed when the stack is not empty. If it is empty,
//     return ErrEmpty. Otherwise remove and return the top value.
//   - Top(): same rule as Pop(), but the value stays on the stack.
//   - Full(), Empty(), Size(): report the current state.
//
// At any moment between calls, 0 <= Size() <= capacity must hold.

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
