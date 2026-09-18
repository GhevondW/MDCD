package day02

// Day 2 -- API Guarantees: Idempotency.
//
// The SAME idempotency key must never be charged twice: a repeat of a key
// you've already seen must return the result of the FIRST charge -- same
// amount, unchanged running total -- not charge again.

type ChargeResult struct {
	Amount       int64
	RunningTotal int64
	WasNew       bool
}

type PaymentProcessor struct {
	seen  map[string]int64
	total int64
}

func NewPaymentProcessor() *PaymentProcessor {
	return &PaymentProcessor{seen: make(map[string]int64)}
}

// TODO: if key was already charged, return the ORIGINAL amount and the
// CURRENT running total, with WasNew=false -- and don't change
// TotalCharged(). Otherwise remember it, add it to the total, and return
// it with WasNew=true.
func (p *PaymentProcessor) Charge(key string, amount int64) ChargeResult {
	return ChargeResult{}
}

func (p *PaymentProcessor) TotalCharged() int64 {
	// TODO
	return 0
}
