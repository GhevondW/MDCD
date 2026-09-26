package day03

// Day 3 -- Critical sections: check-then-act.
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
//
// Deposit, Withdraw and Balance may be called by many goroutines at the
// same time.
//
//	if balance >= amount { balance -= amount }
//
// is a check and an act -- another goroutine can slip in between them.

import "errors"

var ErrInvalidAmount = errors.New("invalid amount")

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
