package day02;

/**
 * Day 2 -- Contracts & Invariants.
 *
 * A stack of ints with a fixed capacity, given once in the constructor.
 * The stack must protect its own rules:
 *   - push(v): only allowed when the stack is not full. If it is already
 *     full, throw IllegalStateException and do not change anything.
 *   - pop(): only allowed when the stack is not empty. If it is empty,
 *     throw NoSuchElementException. Otherwise remove and return the top
 *     value.
 *   - top(): same rule as pop(), but the value stays on the stack.
 *   - isFull(), isEmpty(), size(): report the current state.
 *
 * At any moment between calls, 0 <= size() <= capacity must hold.
 */
public class BoundedStack {
    // TODO: choose your own representation.

    public BoundedStack(int capacity) {
        // TODO
    }

    public boolean isFull() {
        // TODO
        return false;
    }

    public boolean isEmpty() {
        // TODO
        return false;
    }

    public int size() {
        // TODO
        return 0;
    }

    public void push(int value) {
        // TODO
    }

    public int pop() {
        // TODO
        return 0;
    }

    public int top() {
        // TODO
        return 0;
    }
}
