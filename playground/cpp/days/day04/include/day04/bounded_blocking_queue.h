#pragma once
// Day 4 -- The Bounded Blocking Queue
//
// The structure from the slides, and the one to know cold: a FIFO queue
// with a fixed capacity, shared by many producers and many consumers
// (MPMC). Day 3's tryPush / tryPop give up when the queue is full or
// empty. This queue WAITS instead -- asleep, not spinning -- until the
// call can go on.
//
//   BoundedBlockingQueue(capacity)
//       capacity must be at least 1; otherwise throw std::invalid_argument.
//
//   push(v)            wait until there is room, then add v at the back
//                      and return true. If the queue is closed -- already,
//                      or while this call waits -- return false and add
//                      nothing (v is dropped).
//   pop()              wait until there is an item, then remove and return
//                      the one at the front. If the queue is closed AND
//                      empty, return std::nullopt. A closed queue still
//                      hands out the items it holds.
//
//   tryPush(v)         push without waiting: true if v was added; false if
//                      the queue is full or closed.
//   tryPop()           pop without waiting: std::nullopt if it is empty.
//
//   pushFor(v, t)      like push, but wait at most `t`; false on timeout.
//   popFor(t)          like pop, but wait at most `t`; std::nullopt on
//                      timeout, or when the queue is closed and empty.
//
//   close()            no more items will be added. Every thread that is
//                      asleep in push / pop / pushFor / popFor must wake
//                      up and return (false / the remaining items, then
//                      std::nullopt). Calling close() again does nothing.
//   closed()           has close() been called?
//   size(), capacity() the current number of items; the fixed maximum.
//
// The rules that make it hard:
//
//   1. size() never goes above capacity, and items come out in the order
//      they went in.
//   2. A waiting thread sleeps (condition_variable). No busy waiting.
//   3. Wake-ups can be spurious, and by the time a woken thread gets the
//      lock another thread may have taken the item or the room it was
//      woken for: always re-check the condition in a loop.
//   4. Two kinds of sleepers -- "waiting for room" and "waiting for an
//      item" -- share one mutex. Make sure a notification can never reach
//      only the wrong kind. (Two condition variables, as on the slides.)
//   5. A timed call has ONE deadline, fixed when the call starts. A wake-up
//      that finds nothing to do must not restart the timeout.
//   6. T may be move-only (std::unique_ptr<int>): move values in and out,
//      never copy them.
//
// Every method may be called by many threads at the same time.
// Don't use anything that is already a thread-safe queue.

#include <chrono>
#include <cstddef>
#include <optional>
#include <stdexcept>
#include <utility>

template <typename T>
class BoundedBlockingQueue {
public:
    explicit BoundedBlockingQueue(std::size_t capacity) {
        (void)capacity;
        // TODO
    }

    bool push(T value) {
        (void)value;
        // TODO
        return false;
    }

    std::optional<T> pop() {
        // TODO
        return std::nullopt;
    }

    bool tryPush(T value) {
        (void)value;
        // TODO
        return false;
    }

    std::optional<T> tryPop() {
        // TODO
        return std::nullopt;
    }

    bool pushFor(T value, std::chrono::milliseconds timeout) {
        (void)value;
        (void)timeout;
        // TODO
        return false;
    }

    std::optional<T> popFor(std::chrono::milliseconds timeout) {
        (void)timeout;
        // TODO
        return std::nullopt;
    }

    void close() {
        // TODO
    }

    bool closed() const {
        // TODO
        return false;
    }

    std::size_t size() const {
        // TODO
        return 0;
    }

    std::size_t capacity() const {
        // TODO
        return 0;
    }

private:
    // TODO: choose your own representation.
};
