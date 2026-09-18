package day02

import "testing"

func TestExpiringPaymentsFirstChargeIsNew(t *testing.T) {
	p := NewExpiringPaymentProcessor(10)
	r := p.Charge("k", 100, 0)
	if r.Status != StatusNew || r.Amount != 100 || r.RunningTotal != 100 {
		t.Fatalf("unexpected outcome: %+v", r)
	}
	if p.TotalCharged() != 100 {
		t.Errorf("expected TotalCharged() 100, got %d", p.TotalCharged())
	}
}

func TestExpiringPaymentsSameAmountWithinTtlIsReplay(t *testing.T) {
	p := NewExpiringPaymentProcessor(10)
	p.Charge("k", 100, 0)
	r := p.Charge("k", 100, 5)
	if r.Status != StatusReplay || r.Amount != 100 || r.RunningTotal != 100 {
		t.Fatalf("unexpected outcome: %+v", r)
	}
	if p.TotalCharged() != 100 {
		t.Errorf("expected TotalCharged() 100, got %d", p.TotalCharged())
	}
}

func TestExpiringPaymentsDifferentAmountWithinTtlIsConflict(t *testing.T) {
	p := NewExpiringPaymentProcessor(10)
	p.Charge("k", 100, 0)
	r := p.Charge("k", 250, 5)
	if r.Status != StatusConflict {
		t.Fatalf("expected StatusConflict, got %+v", r)
	}
	if r.Amount != 100 { // the RECORDED amount, not 250
		t.Errorf("expected Amount 100, got %d", r.Amount)
	}
	if r.RunningTotal != 100 {
		t.Errorf("expected RunningTotal 100, got %d", r.RunningTotal)
	}
	if p.TotalCharged() != 100 { // nothing was charged
		t.Errorf("expected TotalCharged() 100, got %d", p.TotalCharged())
	}
}

func TestExpiringPaymentsExpiredKeyIsChargedAgain(t *testing.T) {
	p := NewExpiringPaymentProcessor(10)
	p.Charge("k", 100, 0)
	r := p.Charge("k", 100, 20)
	if r.Status != StatusNew {
		t.Errorf("expected StatusNew, got %+v", r)
	}
	if p.TotalCharged() != 200 {
		t.Errorf("expected TotalCharged() 200, got %d", p.TotalCharged())
	}
}

func TestExpiringPaymentsExpiryBoundaryIsInclusive(t *testing.T) {
	// Recorded at 0 with ttl 10 -> expired exactly at now == 10.
	p := NewExpiringPaymentProcessor(10)
	p.Charge("k", 100, 0)
	if s := p.Charge("k", 100, 9).Status; s != StatusReplay {
		t.Errorf("expected StatusReplay at now=9, got %v", s)
	}
	if s := p.Charge("k", 100, 10).Status; s != StatusNew {
		t.Errorf("expected StatusNew at now=10, got %v", s)
	}
	if p.TotalCharged() != 200 {
		t.Errorf("expected TotalCharged() 200, got %d", p.TotalCharged())
	}
}

func TestExpiringPaymentsReplayDoesNotExtendExpiry(t *testing.T) {
	p := NewExpiringPaymentProcessor(10)
	p.Charge("k", 100, 0)
	if s := p.Charge("k", 100, 9).Status; s != StatusReplay {
		t.Errorf("expected StatusReplay at now=9, got %v", s)
	}
	// If the replay at 9 had refreshed the record, the key would still be
	// active at 10. It must not be.
	if s := p.Charge("k", 100, 10).Status; s != StatusNew {
		t.Errorf("expected StatusNew at now=10, got %v", s)
	}
}

func TestExpiringPaymentsConflictChangesNothing(t *testing.T) {
	p := NewExpiringPaymentProcessor(10)
	p.Charge("k", 100, 0)
	p.Charge("k", 250, 5)         // conflict
	r := p.Charge("k", 100, 9)    // original amount again
	if r.Status != StatusReplay { // record still 100
		t.Errorf("expected StatusReplay, got %+v", r)
	}
	if r.Amount != 100 {
		t.Errorf("expected Amount 100, got %d", r.Amount)
	}
	if p.TotalCharged() != 100 {
		t.Errorf("expected TotalCharged() 100, got %d", p.TotalCharged())
	}
	if s := p.Charge("k", 100, 10).Status; s != StatusNew { // expiry not shifted
		t.Errorf("expected StatusNew at now=10, got %v", s)
	}
}

func TestExpiringPaymentsAfterExpiryDifferentAmountIsNewNotConflict(t *testing.T) {
	p := NewExpiringPaymentProcessor(10)
	p.Charge("k", 100, 0)
	r := p.Charge("k", 250, 15)
	if r.Status != StatusNew || r.Amount != 250 {
		t.Fatalf("unexpected outcome: %+v", r)
	}
	if p.TotalCharged() != 350 {
		t.Errorf("expected TotalCharged() 350, got %d", p.TotalCharged())
	}
	// The record was replaced: 250 is now the recorded amount.
	if s := p.Charge("k", 250, 16).Status; s != StatusReplay {
		t.Errorf("expected StatusReplay, got %v", s)
	}
	if s := p.Charge("k", 100, 16).Status; s != StatusConflict {
		t.Errorf("expected StatusConflict, got %v", s)
	}
}

func TestExpiringPaymentsKeysAreIndependent(t *testing.T) {
	p := NewExpiringPaymentProcessor(10)
	p.Charge("a", 100, 0)
	p.Charge("b", 50, 8)
	if p.TotalCharged() != 150 {
		t.Errorf("expected TotalCharged() 150, got %d", p.TotalCharged())
	}
	if s := p.Charge("a", 100, 12).Status; s != StatusNew { // a expired
		t.Errorf("expected StatusNew for a, got %v", s)
	}
	if s := p.Charge("b", 50, 12).Status; s != StatusReplay { // b still active
		t.Errorf("expected StatusReplay for b, got %v", s)
	}
}

func TestExpiringPaymentsActiveKeysCountsOnlyUnexpired(t *testing.T) {
	p := NewExpiringPaymentProcessor(10)
	p.Charge("a", 1, 0) // expires at 10
	p.Charge("b", 2, 5) // expires at 15
	p.Charge("c", 3, 8) // expires at 18
	if n := p.ActiveKeys(0); n != 3 {
		t.Errorf("expected ActiveKeys(0) 3, got %d", n)
	}
	if n := p.ActiveKeys(12); n != 2 {
		t.Errorf("expected ActiveKeys(12) 2, got %d", n)
	}
	if n := p.ActiveKeys(15); n != 1 {
		t.Errorf("expected ActiveKeys(15) 1, got %d", n)
	}
	if n := p.ActiveKeys(18); n != 0 {
		t.Errorf("expected ActiveKeys(18) 0, got %d", n)
	}
}

func TestExpiringPaymentsZeroTtlMeansEveryChargeIsNew(t *testing.T) {
	p := NewExpiringPaymentProcessor(0)
	if s := p.Charge("k", 100, 5).Status; s != StatusNew {
		t.Errorf("expected StatusNew, got %v", s)
	}
	if s := p.Charge("k", 100, 5).Status; s != StatusNew {
		t.Errorf("expected StatusNew, got %v", s)
	}
	if p.TotalCharged() != 200 {
		t.Errorf("expected TotalCharged() 200, got %d", p.TotalCharged())
	}
	if n := p.ActiveKeys(5); n != 0 {
		t.Errorf("expected ActiveKeys(5) 0, got %d", n)
	}
}

func TestExpiringPaymentsReplayReportsCurrentTotalNotTotalAtFirstCharge(t *testing.T) {
	p := NewExpiringPaymentProcessor(100)
	p.Charge("a", 100, 0)
	p.Charge("b", 50, 1)
	r := p.Charge("a", 100, 2)
	if r.Status != StatusReplay {
		t.Errorf("expected StatusReplay, got %+v", r)
	}
	if r.RunningTotal != 150 {
		t.Errorf("expected RunningTotal 150, got %d", r.RunningTotal)
	}
}

func TestExpiringPaymentsTotalIsSumOfNewChargesOnly(t *testing.T) {
	p := NewExpiringPaymentProcessor(10)
	p.Charge("a", 100, 0) // new
	p.Charge("a", 100, 1) // replay
	p.Charge("a", 999, 2) // conflict
	p.Charge("b", 50, 3)  // new
	p.Charge("a", 25, 12) // expired -> new
	if p.TotalCharged() != 175 {
		t.Errorf("expected TotalCharged() 175, got %d", p.TotalCharged())
	}
}
