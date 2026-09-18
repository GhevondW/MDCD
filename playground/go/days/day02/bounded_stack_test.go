package day02

import (
	"errors"
	"testing"
)

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBoundedStackStartsEmpty(t *testing.T) {
	s := NewBoundedStack(2)
	if !s.Empty() {
		t.Error("expected Empty() true")
	}
	if s.Size() != 0 {
		t.Errorf("expected Size() 0, got %d", s.Size())
	}
	if s.Full() {
		t.Error("expected Full() false")
	}
}

func TestBoundedStackPushUpToCapacity(t *testing.T) {
	s := NewBoundedStack(2)
	must(t, s.Push(5))
	must(t, s.Push(7))
	if s.Size() != 2 {
		t.Errorf("expected Size() 2, got %d", s.Size())
	}
	if !s.Full() {
		t.Error("expected Full() true")
	}
}

func TestBoundedStackPushBeyondCapacityErrors(t *testing.T) {
	s := NewBoundedStack(1)
	must(t, s.Push(1))
	if err := s.Push(2); !errors.Is(err, ErrFull) {
		t.Errorf("expected ErrFull, got %v", err)
	}
	if s.Size() != 1 {
		t.Errorf("rejected push must not touch state, got size %d", s.Size())
	}
}

func TestBoundedStackPopReturnsLifoOrder(t *testing.T) {
	s := NewBoundedStack(3)
	must(t, s.Push(1))
	must(t, s.Push(2))
	must(t, s.Push(3))
	if v, err := s.Pop(); err != nil || v != 3 {
		t.Fatalf("expected (3, nil), got (%d, %v)", v, err)
	}
	if v, err := s.Pop(); err != nil || v != 2 {
		t.Fatalf("expected (2, nil), got (%d, %v)", v, err)
	}
	if s.Size() != 1 {
		t.Errorf("expected Size() 1, got %d", s.Size())
	}
}

func TestBoundedStackPopEmptyErrors(t *testing.T) {
	s := NewBoundedStack(1)
	if _, err := s.Pop(); !errors.Is(err, ErrEmpty) {
		t.Errorf("expected ErrEmpty, got %v", err)
	}
}

func TestBoundedStackTopDoesNotRemove(t *testing.T) {
	s := NewBoundedStack(2)
	must(t, s.Push(9))
	if v, err := s.Top(); err != nil || v != 9 {
		t.Fatalf("expected (9, nil), got (%d, %v)", v, err)
	}
	if s.Size() != 1 {
		t.Errorf("expected Size() 1, got %d", s.Size())
	}
}

func TestBoundedStackTopEmptyErrors(t *testing.T) {
	s := NewBoundedStack(1)
	if _, err := s.Top(); !errors.Is(err, ErrEmpty) {
		t.Errorf("expected ErrEmpty, got %v", err)
	}
}

func TestBoundedStackZeroCapacityIsAlwaysFull(t *testing.T) {
	s := NewBoundedStack(0)
	if !s.Full() {
		t.Error("expected Full() true")
	}
	if err := s.Push(1); !errors.Is(err, ErrFull) {
		t.Errorf("expected ErrFull, got %v", err)
	}
}

func TestBoundedStackPopFreesCapacityForAnotherPush(t *testing.T) {
	s := NewBoundedStack(1)
	must(t, s.Push(1))
	if !s.Full() {
		t.Fatal("expected Full() true")
	}
	if _, err := s.Pop(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Full() {
		t.Error("expected Full() false after Pop")
	}
	if !s.Empty() {
		t.Error("expected Empty() true after Pop")
	}
	must(t, s.Push(2)) // must succeed -- capacity was freed
	if v, _ := s.Top(); v != 2 {
		t.Errorf("expected Top() 2, got %d", v)
	}
}

func TestBoundedStackTopReflectsMostRecentPush(t *testing.T) {
	s := NewBoundedStack(3)
	must(t, s.Push(1))
	if v, _ := s.Top(); v != 1 {
		t.Errorf("expected Top() 1, got %d", v)
	}
	must(t, s.Push(2))
	if v, _ := s.Top(); v != 2 {
		t.Errorf("expected Top() 2, got %d", v)
	}
	must(t, s.Push(3))
	if v, _ := s.Top(); v != 3 {
		t.Errorf("expected Top() 3, got %d", v)
	}
}

func TestBoundedStackEmptyAfterPoppingEverything(t *testing.T) {
	s := NewBoundedStack(2)
	must(t, s.Push(1))
	must(t, s.Push(2))
	if _, err := s.Pop(); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Pop(); err != nil {
		t.Fatal(err)
	}
	if !s.Empty() {
		t.Error("expected Empty() true")
	}
	if s.Size() != 0 {
		t.Errorf("expected Size() 0, got %d", s.Size())
	}
	if _, err := s.Pop(); !errors.Is(err, ErrEmpty) {
		t.Errorf("expected ErrEmpty, got %v", err)
	}
}
