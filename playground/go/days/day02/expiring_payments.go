package day02

// Day 2 -- Challenge: Idempotency with Expiry & Conflict.
//
// A payment processor whose idempotency keys EXPIRE. Time is a logical
// clock: every call takes `now` (non-decreasing across calls); a key
// recorded at time T with time-to-live `ttl` is expired once now >= T + ttl.
//
// Charge(key, amount, now):
//   - key unknown, or its record expired  -> charge: add amount to the
//     total, record (key, amount) with expiry now + ttl, return StatusNew.
//   - key active, same amount             -> return StatusReplay with the
//     originally recorded amount and the current total. Nothing changes --
//     not the total, not the recorded amount, not the expiry.
//   - key active, different amount        -> return StatusConflict with
//     the recorded amount and the current total. Nothing changes.
//
// TotalCharged()   sum of all charges that returned StatusNew.
// ActiveKeys(now)  number of recorded keys not yet expired at now.

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
	ttl int64
	// TODO: choose your own representation.
}

func NewExpiringPaymentProcessor(ttl int64) *ExpiringPaymentProcessor {
	return &ExpiringPaymentProcessor{ttl: ttl}
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
