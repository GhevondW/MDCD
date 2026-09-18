package day02

// Day 2 -- API Guarantees: Atomicity.
//
// A record has two fields, A and B. A record is valid only when
// a + b >= 0.
//
// Set(key, a, b) must be all-or-nothing:
//   - If a + b >= 0: store both fields together and return true.
//   - If a + b < 0: the update is invalid. Reject the WHOLE update -- do
//     not change either field -- and return false. This also holds when
//     the key already has a record: the old record must stay exactly as
//     it was.
//
// Get(key) returns the stored record and true, or the zero value and
// false if this key was never successfully set.

type Fields struct {
	A int64
	B int64
}

type RecordStore struct {
	// TODO: choose your own representation.
}

func NewRecordStore() *RecordStore {
	// TODO
	return &RecordStore{}
}

func (r *RecordStore) Set(key string, a, b int64) bool {
	// TODO
	return false
}

func (r *RecordStore) Get(key string) (Fields, bool) {
	// TODO
	return Fields{}, false
}
