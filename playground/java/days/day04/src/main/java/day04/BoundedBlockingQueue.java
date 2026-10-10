package day04;

import java.util.Optional;
import java.util.concurrent.TimeUnit;

/**
 * Day 4 -- The bounded blocking queue (many producers, many consumers).
 *
 * The queue from the lecture: a first-in-first-out queue that holds at
 * most `capacity` items, used by many threads at once. Build it from ONE
 * lock and TWO conditions ("not full" and "not empty") -- `synchronized`
 * with wait/notifyAll, or ReentrantLock with two Conditions. Do not use
 * anything from java.util.concurrent that is already a queue
 * (ArrayBlockingQueue, LinkedBlockingQueue, ...), and do not busy-wait:
 * a thread that cannot go on must sleep.
 *
 * Constructor:
 *   - BoundedBlockingQueue(capacity): throws IllegalArgumentException if
 *     capacity <= 0.
 *
 * Blocking calls (they declare InterruptedException; let it propagate):
 *   - push(x): waits as long as the queue is full, then adds x at the
 *     back and returns true. If the queue is closed -- or becomes closed
 *     while push is waiting -- returns false and adds nothing.
 *   - pop(): waits as long as the queue is empty, then removes and returns
 *     the front item. Once the queue is closed AND empty, returns
 *     Optional.empty() -- immediately, and also to every thread that was
 *     asleep in pop() when the last item was taken or close() was called.
 *
 * Calls that never wait:
 *   - tryPush(x): adds x and returns true if there is room and the queue
 *     is open; otherwise returns false and changes nothing.
 *   - tryPop(): removes and returns the front item; Optional.empty() if
 *     the queue is empty.
 *
 * Calls that wait at most a limited time:
 *   - pushFor(x, timeout, unit): like push, but gives up and returns
 *     false if there is still no room after the timeout.
 *   - popFor(timeout, unit): like pop, but returns Optional.empty() if
 *     there is still nothing after the timeout.
 *   The timeout is a deadline for the WHOLE call: being woken up and
 *   finding the condition still false must not restart the clock.
 *
 * Shutdown:
 *   - close(): no more items are accepted (every push variant returns
 *     false from now on), but the items already inside can still be
 *     popped -- consumers drain the queue, then get Optional.empty().
 *     Every thread sleeping in push, pop, pushFor or popFor must wake up.
 *     Calling close() again does nothing.
 *   - isClosed(), size(), capacity(): report the current state.
 *
 * Whatever happens, 0 <= size() <= capacity, no item is lost, none is
 * returned twice, and items come out in the order they went in.
 */
public class BoundedBlockingQueue<T> {
    // TODO: choose your own representation.

    public BoundedBlockingQueue(int capacity) {
        // TODO
    }

    public boolean push(T item) throws InterruptedException {
        // TODO
        return false;
    }

    public Optional<T> pop() throws InterruptedException {
        // TODO
        return Optional.empty();
    }

    public boolean tryPush(T item) {
        // TODO
        return false;
    }

    public Optional<T> tryPop() {
        // TODO
        return Optional.empty();
    }

    public boolean pushFor(T item, long timeout, TimeUnit unit) throws InterruptedException {
        // TODO
        return false;
    }

    public Optional<T> popFor(long timeout, TimeUnit unit) throws InterruptedException {
        // TODO
        return Optional.empty();
    }

    public void close() {
        // TODO
    }

    public boolean isClosed() {
        // TODO
        return false;
    }

    public int size() {
        // TODO
        return 0;
    }

    public int capacity() {
        // TODO
        return 0;
    }
}
