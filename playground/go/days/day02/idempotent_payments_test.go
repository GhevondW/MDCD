package day02

import "testing"

func TestPaymentProcessorFirstChargeIsNew(t *testing.T) {
	p := NewPaymentProcessor()
	r := p.Charge("orderA", 100)
	if !r.WasNew || r.Amount != 100 || r.RunningTotal != 100 {
		t.Fatalf("unexpected result: %+v", r)
	}
	if p.TotalCharged() != 100 {
		t.Errorf("expected TotalCharged() 100, got %d", p.TotalCharged())
	}
}

func TestPaymentProcessorRepeatKeyIsNotCharged(t *testing.T) {
	p := NewPaymentProcessor()
	p.Charge("orderA", 100)
	r := p.Charge("orderA", 100)
	if r.WasNew {
		t.Error("expected WasNew false on repeat")
	}
	if r.Amount != 100 || r.RunningTotal != 100 {
		t.Fatalf("unexpected result: %+v", r)
	}
	if p.TotalCharged() != 100 {
		t.Errorf("expected TotalCharged() 100, got %d", p.TotalCharged())
	}
}

func TestPaymentProcessorDifferentKeysAccumulate(t *testing.T) {
	p := NewPaymentProcessor()
	p.Charge("orderA", 100)
	p.Charge("orderB", 250)
	if p.TotalCharged() != 350 {
		t.Errorf("expected TotalCharged() 350, got %d", p.TotalCharged())
	}
}

func TestPaymentProcessorRepeatIgnoresNewAmountArgument(t *testing.T) {
	// A caller retrying a request might resend the same amount; a repeat
	// must still report the ORIGINAL amount, never a re-derived one.
	p := NewPaymentProcessor()
	p.Charge("x", 5)
	r := p.Charge("x", 5)
	if r.Amount != 5 {
		t.Errorf("expected Amount 5, got %d", r.Amount)
	}
	if p.TotalCharged() != 5 {
		t.Errorf("expected TotalCharged() 5, got %d", p.TotalCharged())
	}
}

func TestPaymentProcessorManyRepeatsStillChargeOnce(t *testing.T) {
	p := NewPaymentProcessor()
	p.Charge("x", 5)
	p.Charge("x", 5)
	p.Charge("x", 5)
	p.Charge("y", 10)
	if p.TotalCharged() != 15 {
		t.Errorf("expected TotalCharged() 15, got %d", p.TotalCharged())
	}
}

func TestPaymentProcessorRepeatIsIdempotentEvenWithOtherChargesBetween(t *testing.T) {
	p := NewPaymentProcessor()
	p.Charge("A", 100)
	p.Charge("B", 50)
	r := p.Charge("A", 100) // repeat of A, not adjacent to the first
	if r.WasNew || r.Amount != 100 || r.RunningTotal != 150 {
		t.Fatalf("unexpected result: %+v", r)
	}
	p.Charge("C", 25)
	if p.TotalCharged() != 175 { // A never double-counted
		t.Errorf("expected TotalCharged() 175, got %d", p.TotalCharged())
	}
}

func TestPaymentProcessorZeroAmountChargeIsStillTrackedByKey(t *testing.T) {
	p := NewPaymentProcessor()
	first := p.Charge("free-trial", 0)
	if !first.WasNew || first.Amount != 0 {
		t.Fatalf("unexpected result: %+v", first)
	}
	repeat := p.Charge("free-trial", 0)
	if repeat.WasNew {
		t.Error("expected WasNew false on repeat")
	}
	if p.TotalCharged() != 0 {
		t.Errorf("expected TotalCharged() 0, got %d", p.TotalCharged())
	}
}
