package day02

// Day 2 -- API Guarantees: Atomicity.
//
// A record has two fields, A and B. Set() must be all-or-nothing: if
// a + b < 0 the update is invalid and must be REJECTED IN FULL -- neither
// field may change, whether the key already existed or not. Applying one
// field before validating the other is exactly the bug this catches.

type Fields struct {
	A int64
	B int64
}

type RecordStore struct {
	store map[string]Fields
}

func NewRecordStore() *RecordStore {
	return &RecordStore{store: make(map[string]Fields)}
}

// TODO: if a + b >= 0, commit BOTH fields together and return true.
// Otherwise leave the store completely unchanged and return false.
func (r *RecordStore) Set(key string, a, b int64) bool {
	return false
}

// TODO: return the record and true if key has been successfully set, or
// the zero value and false otherwise.
func (r *RecordStore) Get(key string) (Fields, bool) {
	return Fields{}, false
}
