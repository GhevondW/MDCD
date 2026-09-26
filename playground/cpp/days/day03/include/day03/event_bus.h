#pragma once
// Day 3 -- Challenge: A Thread-Safe Event Bus (Observer, again)
//
// Day 1 left a question to write down: "What happens if an observer
// subscribes at the exact moment the publisher is notifying?" This is
// Day 1's publisher, made safe for many threads -- and for callbacks
// that call back into the bus.
//
//   subscribe(callback)  register a callback. Returns a new id (> 0);
//                        no two calls ever return the same id.
//   unsubscribe(id)      remove that subscriber. Returns false if the id
//                        is unknown or was already removed.
//   publish(event)       call every subscriber with `event`, on the
//                        calling thread, before returning. Returns how
//                        many callbacks were called.
//   subscriberCount()    how many subscribers are registered right now.
//
// Rules:
//
//   1. Order: a publish calls subscribers in the order they subscribed.
//   2. Snapshot: a publish calls exactly the subscribers registered at
//      the moment it starts. Changes made while it runs -- by another
//      thread, or by a callback -- only affect later publishes: a
//      subscriber added mid-publish is not called by it, and one removed
//      mid-publish is still called by it.
//   3. Callbacks may use the bus: a callback can call subscribe,
//      unsubscribe, or publish on the same bus without deadlocking --
//      like the slides' Welcomer, which subscribes from update().
//   4. Never run a callback while holding your lock. A callback is
//      someone else's code -- it may be slow, or wait for another thread
//      that needs the bus. Other threads must be able to use the bus
//      while a callback runs.
//
// Every method may be called by many threads at the same time.

#include <cstddef>
#include <functional>
#include <string>

class EventBus {
public:
    using Callback = std::function<void(const std::string& event)>;

    long long subscribe(Callback callback) {
        (void)callback;
        // TODO
        return 0;
    }

    bool unsubscribe(long long id) {
        (void)id;
        // TODO
        return false;
    }

    std::size_t publish(const std::string& event) {
        (void)event;
        // TODO
        return 0;
    }

    std::size_t subscriberCount() const {
        // TODO
        return 0;
    }

private:
    // TODO: choose your own representation.
};
