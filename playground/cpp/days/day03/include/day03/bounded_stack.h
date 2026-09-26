#pragma once
// Day 3 -- The bounded stack, shared
//
// Day 2's bounded stack, with one change: many threads use it at the
// same time. The contract is the same:
//
//   - push(v): only allowed when the stack is not full. If it is already
//     full, throw std::runtime_error and do not change anything.
//   - pop(): only allowed when the stack is not empty. If it is empty,
//     throw std::runtime_error. Otherwise remove and return the top value.
//   - top(): same rule as pop(), but the value stays on the stack.
//   - full(), empty(), size(): report the current state.
//
// Day 2 said: the invariant 0 <= size() <= capacity may be broken for a
// moment *inside* a method, as long as it is restored before the method
// returns. With two threads, another thread can arrive in exactly that
// moment. No thread may ever see or cause a half-done push or pop: no
// value lost, no value returned twice, no value made up.

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
