package day02

// Day 2 -- Contracts & Invariants.
//
// Implement BoundedStack so it enforces its own contract:
//   - Push(v): precondition "not full" -- if the stack IS full, return
//     ErrFull instead of silently corrupting state.
//   - Pop() / Top(): precondition "not empty" -- return ErrEmpty instead.
//   - Invariant to hold at all times between calls: 0 <= Size() <= capacity.

import "errors"

var ErrFull = errors.New("full")
var ErrEmpty = errors.New("empty")

type BoundedStack struct {
	data     []int
	capacity int
}

func NewBoundedStack(capacity int) *BoundedStack {
	return &BoundedStack{capacity: capacity}
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
	// TODO: return ErrFull if Full(), otherwise store v and return nil.
	return nil
}

func (s *BoundedStack) Pop() (int, error) {
	// TODO: return 0, ErrEmpty if Empty(), otherwise remove and return
	// the top value with a nil error.
	return 0, nil
}

func (s *BoundedStack) Top() (int, error) {
	// TODO: same precondition as Pop(), but don't remove anything.
	return 0, nil
}
