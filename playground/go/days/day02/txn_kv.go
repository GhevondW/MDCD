package day02

// Day 2 -- Challenge: Transactional Key-Value Store.
//
// A key-value store (string -> string) with transactions. A transaction
// can be opened inside another transaction (they nest).
//
//	Set(key, value)   store a value. If a transaction is open, the change
//	                  belongs to that transaction; otherwise it is applied
//	                  directly to the store.
//	Get(key)          the value visible right now, with ok=false if the
//	                  key does not exist. Changes made inside an open
//	                  transaction are visible immediately.
//	Remove(key)       delete the key. Returns true if the key was visible
//	                  before the call. A delete made inside a transaction
//	                  also belongs to that transaction.
//	Begin()           open a new transaction.
//	Commit()          take everything the innermost transaction changed --
//	                  writes AND deletes -- and move it into the
//	                  transaction one level out (or into the store itself
//	                  if there is no outer transaction). Returns false if
//	                  no transaction is open.
//	Rollback()        throw away everything the innermost transaction
//	                  changed. Returns false if no transaction is open.

type TxnKV struct {
	// TODO: choose your own representation.
}

func NewTxnKV() *TxnKV {
	// TODO
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
