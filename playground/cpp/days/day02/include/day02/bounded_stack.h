#pragma once
// Day 2 -- Contracts & Invariants
//
// A stack of ints with a fixed capacity, given once in the constructor.
// The stack must protect its own rules:
//
//   - push(v): only allowed when the stack is not full. If it is already
//     full, throw std::runtime_error and do not change anything.
//   - pop(): only allowed when the stack is not empty. If it is empty,
//     throw std::runtime_error. Otherwise remove and return the top value.
//   - top(): same rule as pop(), but the value stays on the stack.
//   - full(), empty(), size(): report the current state.
//
// At any moment between calls, 0 <= size() <= capacity must hold.

#include <cstddef>
#include <stdexcept>

class BoundedStack {
public:
    explicit BoundedStack(std::size_t capacity) {
        (void)capacity;
        // TODO
    }

    bool full() const {
        // TODO
        return false;
    }

    bool empty() const {
        // TODO
        return false;
    }

    std::size_t size() const {
        // TODO
        return 0;
    }

    void push(int v) {
        (void)v;
        // TODO
    }

    int pop() {
        // TODO
        return 0;
    }

    int top() const {
        // TODO
        return 0;
    }

private:
    // TODO: choose your own representation.
};
