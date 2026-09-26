package day03;

import java.util.OptionalInt;

/**
 * Day 3 -- The bounded stack, shared.
 *
 * Day 2's bounded stack, with one change: many threads use it at the
 * same time. The contract is the same:
 *   - push(v): only allowed when the stack is not full. If it is already
 *     full, throw IllegalStateException and do not change anything.
 *   - pop(): only allowed when the stack is not empty. If it is empty,
 *     throw NoSuchElementException. Otherwise remove and return the top
 *     value.
 *   - top(): same rule as pop(), but the value stays on the stack.
 *   - isFull(), isEmpty(), size(): report the current state.
 *
 * Day 2 said: the invariant 0 <= size() <= capacity may be broken for a
 * moment *inside* a method, as long as it is restored before the method
 * returns. With two threads, another thread can arrive in exactly that
 * moment. No thread may ever see or cause a half-done push or pop: no
 * value lost, no value returned twice, no value made up.
 *
 * Even with every method locked, a caller who writes
 *
 *     if (!s.isEmpty()) { int v = s.top(); s.pop(); ... }
 *
 * has an API race: another thread can run between the calls. So the
 * stack also offers calls that check and act as ONE step:
 *   - tryPush(v): if there is room, push v and return true. If the stack
 *     is full, return false and change nothing. Never throws.
 *   - tryPop(): if the stack is not empty, remove the top value and
 *     return it. If it is empty, return OptionalInt.empty(). Never
 *     throws.
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

    public boolean tryPush(int value) {
        // TODO
        return false;
    }

    public OptionalInt tryPop() {
        // TODO
        return OptionalInt.empty();
    }
}
