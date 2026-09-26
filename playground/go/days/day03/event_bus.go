package day03

// Day 3 -- Challenge: A Thread-Safe Event Bus (Observer, again).
//
// Day 1 left a question to write down: "What happens if an observer
// subscribes at the exact moment the publisher is notifying?" This is
// Day 1's publisher, made safe for many goroutines -- and for callbacks
// that call back into the bus.
//
//	Subscribe(callback)  register a callback. Returns a new id (> 0);
//	                     no two calls ever return the same id.
//	Unsubscribe(id)      remove that subscriber. Returns false if the id
//	                     is unknown or was already removed.
//	Publish(event)       call every subscriber with event, on the
//	                     calling goroutine, before returning. Returns how
//	                     many callbacks were called.
//	SubscriberCount()    how many subscribers are registered right now.
//
// Rules:
//
//  1. Order: a publish calls subscribers in the order they subscribed.
//  2. Snapshot: a publish calls exactly the subscribers registered at
//     the moment it starts. Changes made while it runs -- by another
//     goroutine, or by a callback -- only affect later publishes: a
//     subscriber added mid-publish is not called by it, and one removed
//     mid-publish is still called by it.
//  3. Callbacks may use the bus: a callback can call Subscribe,
//     Unsubscribe, or Publish on the same bus without deadlocking.
//  4. Never run a callback while holding your lock. A callback is
//     someone else's code -- it may be slow, or wait for another
//     goroutine that needs the bus. Other goroutines must be able to use
//     the bus while a callback runs.
//
// Every method may be called by many goroutines at the same time.

type EventBus struct {
	// TODO: choose your own representation.
}

func NewEventBus() *EventBus {
	// TODO
	return &EventBus{}
}

func (b *EventBus) Subscribe(callback func(event string)) int64 {
	// TODO
	return 0
}

func (b *EventBus) Unsubscribe(id int64) bool {
	// TODO
	return false
}

func (b *EventBus) Publish(event string) int {
	// TODO
	return 0
}

func (b *EventBus) SubscriberCount() int {
	// TODO
	return 0
}
