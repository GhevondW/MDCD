package day02

import "testing"

func TestTxnKVSetThenGet(t *testing.T) {
	kv := NewTxnKV()
	kv.Set("a", "1")
	v, ok := kv.Get("a")
	if !ok {
		t.Fatal("expected key to be present")
	}
	if v != "1" {
		t.Errorf("expected %q, got %q", "1", v)
	}
}

func TestTxnKVGetMissingKeyReturnsNothing(t *testing.T) {
	kv := NewTxnKV()
	if _, ok := kv.Get("nope"); ok {
		t.Error("expected ok=false for a missing key")
	}
}

func TestTxnKVOverwriteReplacesValue(t *testing.T) {
	kv := NewTxnKV()
	kv.Set("a", "1")
	kv.Set("a", "2")
	if v, _ := kv.Get("a"); v != "2" {
		t.Errorf("expected %q, got %q", "2", v)
	}
}

func TestTxnKVRemoveReturnsTrueOnlyIfKeyWasVisible(t *testing.T) {
	kv := NewTxnKV()
	kv.Set("a", "1")
	if !kv.Remove("a") {
		t.Error("expected Remove to return true for a visible key")
	}
	if _, ok := kv.Get("a"); ok {
		t.Error("expected key to be gone after Remove")
	}
	if kv.Remove("a") {
		t.Error("expected Remove to return false for an already-removed key")
	}
	if kv.Remove("never-set") {
		t.Error("expected Remove to return false for a key never set")
	}
}

func TestTxnKVCommitWithoutOpenTransactionFails(t *testing.T) {
	kv := NewTxnKV()
	if kv.Commit() {
		t.Error("expected Commit to return false with no open transaction")
	}
}

func TestTxnKVRollbackWithoutOpenTransactionFails(t *testing.T) {
	kv := NewTxnKV()
	if kv.Rollback() {
		t.Error("expected Rollback to return false with no open transaction")
	}
}

func TestTxnKVSetInsideTransactionIsVisibleBeforeCommit(t *testing.T) {
	kv := NewTxnKV()
	kv.Begin()
	kv.Set("a", "1")
	v, ok := kv.Get("a")
	if !ok {
		t.Fatal("expected key to be present")
	}
	if v != "1" {
		t.Errorf("expected %q, got %q", "1", v)
	}
}

func TestTxnKVRollbackDiscardsEveryChangeOfTheTransaction(t *testing.T) {
	kv := NewTxnKV()
	kv.Set("a", "old")
	kv.Begin()
	kv.Set("a", "new")
	kv.Set("b", "1")
	if !kv.Rollback() {
		t.Error("expected Rollback to return true")
	}
	if v, _ := kv.Get("a"); v != "old" {
		t.Errorf("expected %q, got %q", "old", v)
	}
	if _, ok := kv.Get("b"); ok {
		t.Error("expected b to be gone after Rollback")
	}
}

func TestTxnKVCommitMakesChangesPermanent(t *testing.T) {
	kv := NewTxnKV()
	kv.Begin()
	kv.Set("a", "1")
	if !kv.Commit() {
		t.Error("expected Commit to return true")
	}
	v, ok := kv.Get("a")
	if !ok {
		t.Fatal("expected key to be present")
	}
	if v != "1" {
		t.Errorf("expected %q, got %q", "1", v)
	}
	if kv.Rollback() { // nothing left open to roll back
		t.Error("expected Rollback to return false")
	}
	if v, _ := kv.Get("a"); v != "1" {
		t.Errorf("expected %q, got %q", "1", v)
	}
}

func TestTxnKVRemoveInsideTransactionHidesOuterValue(t *testing.T) {
	kv := NewTxnKV()
	kv.Set("a", "base")
	kv.Begin()
	if !kv.Remove("a") {
		t.Error("expected Remove to return true")
	}
	if _, ok := kv.Get("a"); ok { // gone as seen from inside
		t.Error("expected key to be hidden inside the transaction")
	}
}

func TestTxnKVRollbackRestoresValueRemovedInsideTransaction(t *testing.T) {
	kv := NewTxnKV()
	kv.Set("a", "base")
	kv.Begin()
	kv.Remove("a")
	kv.Rollback()
	v, ok := kv.Get("a")
	if !ok {
		t.Fatal("expected key to be restored after Rollback")
	}
	if v != "base" {
		t.Errorf("expected %q, got %q", "base", v)
	}
}

func TestTxnKVCommittedRemoveDeletesFromTheStore(t *testing.T) {
	kv := NewTxnKV()
	kv.Set("a", "base")
	kv.Begin()
	kv.Remove("a")
	if !kv.Commit() {
		t.Error("expected Commit to return true")
	}
	if _, ok := kv.Get("a"); ok {
		t.Error("expected key to be deleted from the store")
	}
}

func TestTxnKVNestedInnerRollbackKeepsOuterChanges(t *testing.T) {
	kv := NewTxnKV()
	kv.Begin()
	kv.Set("outer", "1")
	kv.Begin()
	kv.Set("inner", "2")
	if !kv.Rollback() { // discards only the inner transaction
		t.Error("expected Rollback to return true")
	}
	if _, ok := kv.Get("inner"); ok {
		t.Error("expected inner to be gone")
	}
	v, ok := kv.Get("outer")
	if !ok {
		t.Fatal("expected outer to be present")
	}
	if v != "1" {
		t.Errorf("expected %q, got %q", "1", v)
	}
}

func TestTxnKVInnerCommitMergesIntoParentNotIntoStore(t *testing.T) {
	kv := NewTxnKV()
	kv.Begin()
	kv.Begin()
	kv.Set("a", "1")
	if !kv.Commit() { // inner commit: `a` now belongs to the OUTER txn
		t.Error("expected Commit to return true")
	}
	if v, _ := kv.Get("a"); v != "1" {
		t.Errorf("expected %q, got %q", "1", v)
	}
	if !kv.Rollback() { // outer rollback must take `a` down with it
		t.Error("expected Rollback to return true")
	}
	if _, ok := kv.Get("a"); ok {
		t.Error("expected a to be gone after the outer rollback")
	}
}

func TestTxnKVInnerCommittedRemoveStaysInsideParent(t *testing.T) {
	kv := NewTxnKV()
	kv.Set("a", "base")
	kv.Begin()
	kv.Begin()
	kv.Remove("a")
	if !kv.Commit() { // remove merged into the outer txn
		t.Error("expected Commit to return true")
	}
	if _, ok := kv.Get("a"); ok {
		t.Error("expected a to be hidden after the inner commit")
	}
	if !kv.Rollback() { // outer rollback cancels the remove too
		t.Error("expected Rollback to return true")
	}
	v, ok := kv.Get("a")
	if !ok {
		t.Fatal("expected a to be restored after the outer rollback")
	}
	if v != "base" {
		t.Errorf("expected %q, got %q", "base", v)
	}
}

func TestTxnKVSetAfterRemoveInSameTransactionResurrectsKey(t *testing.T) {
	kv := NewTxnKV()
	kv.Set("a", "base")
	kv.Begin()
	kv.Remove("a")
	kv.Set("a", "reborn")
	if v, _ := kv.Get("a"); v != "reborn" {
		t.Errorf("expected %q, got %q", "reborn", v)
	}
	kv.Commit()
	if v, _ := kv.Get("a"); v != "reborn" {
		t.Errorf("expected %q, got %q", "reborn", v)
	}
}

func TestTxnKVValueIsVisibleThroughMultipleUntouchedLayers(t *testing.T) {
	kv := NewTxnKV()
	kv.Set("a", "base")
	kv.Begin()
	kv.Begin()
	kv.Begin()
	v, ok := kv.Get("a")
	if !ok {
		t.Fatal("expected key to be present")
	}
	if v != "base" {
		t.Errorf("expected %q, got %q", "base", v)
	}
}
