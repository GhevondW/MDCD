package day02

// Day 2 -- Challenge: Idempotency with Expiry & Conflict.
//
// Like the idempotent payment processor, but keys are only remembered for
// a limited time. There is no real clock: every call receives the current
// time `now` as a plain number, and the numbers never go down from call
// to call. A key recorded at time T is remembered until time T + ttl;
// from now == T + ttl on, it is expired (forgotten).
//
// Charge(key, amount, now) -- three cases:
//   - The key is unknown, or its record has expired: this is a NEW
//     charge (StatusNew). Add amount to the total and remember
//     (key, amount) until now + ttl.
//   - The key is still remembered and amount equals the recorded amount:
//     this is a REPLAY (StatusReplay) -- the caller sent the same request
//     twice. Return the recorded amount and the current total. Change
//     NOTHING: not the total, not the recorded amount, not the expiry
//     time.
//   - The key is still remembered but amount is different: this is a
//     CONFLICT (StatusConflict) -- the same key was used for a different
//     request, a mistake. Return the recorded amount and the current
//     total. Change nothing.
//
// TotalCharged()   sum of all NEW charges.
// ActiveKeys(now)  how many recorded keys are not yet expired at now.

type ChargeStatus int

const (
	StatusNew ChargeStatus = iota
	StatusReplay
	StatusConflict
)

type ChargeOutcome struct {
	Status       ChargeStatus
	Amount       int64 // StatusNew: charged; StatusReplay/StatusConflict: recorded
	RunningTotal int64
}

type ExpiringPaymentProcessor struct {
	// TODO: choose your own representation.
}

func NewExpiringPaymentProcessor(ttl int64) *ExpiringPaymentProcessor {
	// TODO
	return &ExpiringPaymentProcessor{}
}

func (p *ExpiringPaymentProcessor) Charge(key string, amount, now int64) ChargeOutcome {
	// TODO
	return ChargeOutcome{}
}

func (p *ExpiringPaymentProcessor) TotalCharged() int64 {
	// TODO
	return 0
}

func (p *ExpiringPaymentProcessor) ActiveKeys(now int64) int {
	// TODO
	return 0
}
