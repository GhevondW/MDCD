#pragma once
// Day 2 -- Contracts & Invariants
//
// Implement BoundedStack so it enforces its own contract:
//   - push(v): precondition "not full" -- if the stack IS full, throw
//     std::runtime_error instead of silently corrupting state.
//   - pop() / top(): precondition "not empty" -- same idea.
//   - Invariant to hold at all times between calls: 0 <= size() <= capacity.

#include <cstddef>
#include <stdexcept>
#include <vector>

class BoundedStack {
public:
    explicit BoundedStack(std::size_t capacity) : cap_(capacity) {}

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
        // TODO: throw std::runtime_error("full") if full(), otherwise store v.
    }

    int pop() {
        // TODO: throw std::runtime_error("empty") if empty(), otherwise
        // remove and return the top value.
        return 0;
    }

    int top() const {
        // TODO: same precondition as pop(), but don't remove anything.
        return 0;
    }

private:
    std::vector<int> data_;
    std::size_t cap_;
};
