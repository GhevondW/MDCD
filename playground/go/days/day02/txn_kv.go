package day02

// Day 2 -- Challenge: Transactional Key-Value Store.
//
// In-memory string->string store with NESTED transactions.
//
//	Set(key, value)   store a value (in the innermost open transaction,
//	                  or directly if none is open)
//	Get(key)          value currently visible for key, with ok=false
//	                  when nothing is visible
//	Remove(key)       delete key; returns whether it was visible
//	Begin()           open a new transaction (transactions nest)
//	Commit()          merge the innermost transaction's changes -- writes
//	                  AND deletes -- into its parent transaction (or into
//	                  the store if it has no parent). false if no
//	                  transaction is open.
//	Rollback()        discard the innermost transaction's changes.
//	                  false if no transaction is open.
//
// Reads always see the innermost state: a value set (or a key removed)
// inside an open transaction is visible to Get() immediately.

type TxnKV struct {
	// TODO: choose your own representation.
}

func NewTxnKV() *TxnKV {
	return &TxnKV{}
}

func (kv *TxnKV) Set(key, value string) {
	// TODO
}

func (kv *TxnKV) Get(key string) (string, bool) {
	// TODO
	return "", false
}

func (kv *TxnKV) Remove(key string) bool {
	// TODO
	return false
}

func (kv *TxnKV) Begin() {
	// TODO
}

func (kv *TxnKV) Commit() bool {
	// TODO
	return false
}

func (kv *TxnKV) Rollback() bool {
	// TODO
	return false
}
