package day02

// Day 2 -- API Guarantees: Idempotency.
//
// Every charge request comes with a key. The same key must never be
// charged twice.
//
// Charge(key, amount):
//   - If this key was never seen before: charge it. Add amount to the
//     total and return ChargeResult{amount, runningTotal, WasNew: true}.
//   - If this key was already charged: do NOT charge again. Return the
//     amount of the FIRST charge and the current running total, with
//     WasNew=false. TotalCharged() must stay the same.
//
// TotalCharged() returns the sum of all real (first-time) charges.

type ChargeResult struct {
	Amount       int64
	RunningTotal int64
	WasNew       bool
}

type PaymentProcessor struct {
	// TODO: choose your own representation.
}

func NewPaymentProcessor() *PaymentProcessor {
	// TODO
	return &PaymentProcessor{}
}

func (p *PaymentProcessor) Charge(key string, amount int64) ChargeResult {
	// TODO
	return ChargeResult{}
}

func (p *PaymentProcessor) TotalCharged() int64 {
	// TODO
	return 0
}
