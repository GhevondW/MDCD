package day02;

import java.util.ArrayList;
import java.util.List;
import java.util.NoSuchElementException;

/**
 * Day 2 -- Contracts & Invariants.
 *
 * Implement BoundedStack so it enforces its own contract:
 *   - push(v): precondition "not full" -- if the stack IS full, throw
 *     IllegalStateException instead of silently corrupting state.
 *   - pop() / top(): precondition "not empty" -- throw
 *     NoSuchElementException instead.
 *   - Invariant to hold at all times between calls: 0 <= size() <= capacity.
 */
public class BoundedStack {
    private final List<Integer> data = new ArrayList<>();
    private final int capacity;

    public BoundedStack(int capacity) {
        this.capacity = capacity;
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
        // TODO: throw new IllegalStateException("full") if isFull(),
        // otherwise store value.
    }

    public int pop() {
        // TODO: throw new NoSuchElementException("empty") if isEmpty(),
        // otherwise remove and return the top value.
        return 0;
    }

    public int top() {
        // TODO: same precondition as pop(), but don't remove anything.
        return 0;
    }
}
