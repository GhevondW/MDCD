package day03

// Day 3 -- Critical sections: check-then-act, and two locks at once.
//
// A bank account shared by many goroutines -- think of one card used by
// several people at the same moment. The balance must never go below
// zero.
//
//	NewAccount(initial)  start with initial. initial must be >= 0; if
//	                     not, return ErrInvalidAmount.
//	Deposit(amount)      add amount. amount must be > 0; if not, return
//	                     ErrInvalidAmount and change nothing.
//	Withdraw(amount)     amount must be > 0; if not, return
//	                     ErrInvalidAmount and change nothing.
//	                     If balance >= amount: subtract it, return true.
//	                     Otherwise return false and change nothing.
//	Balance()            the current balance.
//	a.Transfer(to, amount)
//	                     move amount from a to to, as ONE step: no
//	                     goroutine may ever see the money gone from one
//	                     account but not yet in the other. amount must
//	                     be > 0 (else ErrInvalidAmount) and to must be a
//	                     different account (else ErrSameAccount). If a's
//	                     balance < amount, return false and change
//	                     nothing.
//	Total(a, b)          a's balance plus b's balance, read as ONE step:
//	                     every transfer between a and b happens either
//	                     completely before it or completely after it. a
//	                     and b must be different accounts (else
//	                     ErrSameAccount).
//
// Everything except NewAccount may be called by many goroutines at the
// same time.
//
//	if balance >= amount { balance -= amount }
//
// is a check and an act -- another goroutine can slip in between them.
//
// Transfer and Total need two accounts locked at once. Two goroutines
// that lock the same two accounts in opposite orders wait for each other
// forever -- the slides' "Two transfers, two locks". Hint: if every
// goroutine locks any two accounts in the same order (say, by a number
// each account gets when it is created), that cannot happen. And a
// method that already holds a lock must not call another method that
// takes the same lock: a sync.Mutex locked twice by the same goroutine
// waits for itself forever.

import "errors"

var ErrInvalidAmount = errors.New("invalid amount")
var ErrSameAccount = errors.New("same account")

type Account struct {
	// TODO: choose your own representation.
}

func NewAccount(initial int64) (*Account, error) {
	// TODO
	return &Account{}, nil
}

func (a *Account) Deposit(amount int64) error {
	// TODO
	return nil
}

func (a *Account) Withdraw(amount int64) (bool, error) {
	// TODO
	return false, nil
}

func (a *Account) Balance() int64 {
	// TODO
	return 0
}

func (a *Account) Transfer(to *Account, amount int64) (bool, error) {
	// TODO
	return false, nil
}

func Total(a, b *Account) (int64, error) {
	// TODO
	return 0, nil
}
