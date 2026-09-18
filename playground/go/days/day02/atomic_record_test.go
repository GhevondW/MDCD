package day02

import "testing"

func TestRecordStoreGetOnUnknownKeyIsFalse(t *testing.T) {
	store := NewRecordStore()
	if _, ok := store.Get("ghost"); ok {
		t.Error("expected ok=false for unknown key")
	}
}

func TestRecordStoreValidSetCommitsBothFields(t *testing.T) {
	store := NewRecordStore()
	if !store.Set("acct1", 10, 5) {
		t.Fatal("expected Set to succeed")
	}
	r, ok := store.Get("acct1")
	if !ok || r.A != 10 || r.B != 5 {
		t.Fatalf("unexpected record: %+v ok=%v", r, ok)
	}
}

func TestRecordStoreInvalidSetIsRejected(t *testing.T) {
	store := NewRecordStore()
	if store.Set("acct2", -5, -5) {
		t.Error("expected Set to fail")
	}
	if _, ok := store.Get("acct2"); ok {
		t.Error("expected key to be absent after rejected set")
	}
}

func TestRecordStoreRejectedSetLeavesExistingRecordUntouched(t *testing.T) {
	// The core atomicity check: a rejected update on an EXISTING key must
	// not partially apply -- field A must not change even though it would
	// have been "valid" applied in isolation.
	store := NewRecordStore()
	store.Set("acct1", 10, 5)
	if store.Set("acct1", -20, 3) { // -20 + 3 < 0 -> reject
		t.Fatal("expected Set to fail")
	}
	r, ok := store.Get("acct1")
	if !ok || r.A != 10 || r.B != 5 {
		t.Fatalf("record must be unchanged, got %+v ok=%v", r, ok)
	}
}

func TestRecordStoreZeroSumIsValid(t *testing.T) {
	store := NewRecordStore()
	if !store.Set("a", 0, 0) {
		t.Fatal("expected Set to succeed")
	}
	r, ok := store.Get("a")
	if !ok || r.A != 0 || r.B != 0 {
		t.Fatalf("unexpected record: %+v ok=%v", r, ok)
	}
}

func TestRecordStoreNegativeFieldWithCompensatingPositiveIsValid(t *testing.T) {
	// The rule is a + b >= 0, NOT "both fields individually
	// non-negative" -- a single negative field with a large enough
	// positive one is fine.
	store := NewRecordStore()
	if !store.Set("acct3", -5, 10) {
		t.Fatal("expected Set to succeed")
	}
	r, ok := store.Get("acct3")
	if !ok || r.A != -5 || r.B != 10 {
		t.Fatalf("unexpected record: %+v ok=%v", r, ok)
	}
}

func TestRecordStoreValidSetOverwritesPreviousValidRecord(t *testing.T) {
	store := NewRecordStore()
	store.Set("acct1", 10, 5)
	if !store.Set("acct1", 1, 2) {
		t.Fatal("expected Set to succeed")
	}
	r, _ := store.Get("acct1")
	if r.A != 1 || r.B != 2 {
		t.Fatalf("unexpected record: %+v", r)
	}
}

func TestRecordStoreKeysAreIndependent(t *testing.T) {
	store := NewRecordStore()
	store.Set("acct1", 10, 5)
	store.Set("acct2", 1, 1)
	if store.Set("acct2", -100, 0) { // reject on acct2 only
		t.Fatal("expected Set to fail")
	}
	r1, _ := store.Get("acct1")
	r2, _ := store.Get("acct2")
	if r1.A != 10 || r1.B != 5 {
		t.Errorf("acct1 unexpectedly changed: %+v", r1)
	}
	if r2.A != 1 || r2.B != 1 {
		t.Errorf("acct2 unexpectedly changed: %+v", r2)
	}
}
