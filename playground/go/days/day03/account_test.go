package day03

import (
	"errors"
	"sync/atomic"
	"testing"
)

func newAccount(t *testing.T, initial int64) *Account {
	t.Helper()
	a, err := NewAccount(initial)
	if err != nil {
		t.Fatalf("NewAccount(%d): unexpected error: %v", initial, err)
	}
	return a
}

func TestAccountStartsWithInitialBalance(t *testing.T) {
	watchdog(t)
	a := newAccount(t, 100)
	if got := a.Balance(); got != 100 {
		t.Errorf("Balance() = %d, want 100", got)
	}
}

func TestAccountNegativeInitialBalanceErrors(t *testing.T) {
	watchdog(t)
	if _, err := NewAccount(-1); !errors.Is(err, ErrInvalidAmount) {
		t.Errorf("NewAccount(-1): got error %v, want ErrInvalidAmount", err)
	}
}

func TestAccountDepositAddsAndWithdrawSubtracts(t *testing.T) {
	watchdog(t)
	a := newAccount(t, 0)
	must(t, a.Deposit(50))
	if ok, err := a.Withdraw(20); err != nil || !ok {
		t.Errorf("Withdraw(20) = (%v, %v), want (true, nil)", ok, err)
	}
	if got := a.Balance(); got != 30 {
		t.Errorf("Balance() = %d, want 30", got)
	}
}

func TestAccountWithdrawExactBalanceSucceeds(t *testing.T) {
	watchdog(t)
	a := newAccount(t, 40)
	if ok, err := a.Withdraw(40); err != nil || !ok {
		t.Errorf("Withdraw(40) = (%v, %v), want (true, nil)", ok, err)
	}
	if got := a.Balance(); got != 0 {
		t.Errorf("Balance() = %d, want 0", got)
	}
}

func TestAccountWithdrawMoreThanBalanceFailsAndChangesNothing(t *testing.T) {
	watchdog(t)
	a := newAccount(t, 30)
	if ok, err := a.Withdraw(31); err != nil || ok {
		t.Errorf("Withdraw(31) = (%v, %v), want (false, nil)", ok, err)
	}
	if got := a.Balance(); got != 30 {
		t.Errorf("Balance() = %d, want 30", got)
	}
}

func TestAccountNonPositiveAmountsErrorAndChangeNothing(t *testing.T) {
	watchdog(t)
	a := newAccount(t, 10)
	for _, amount := range []int64{0, -5} {
		if err := a.Deposit(amount); !errors.Is(err, ErrInvalidAmount) {
			t.Errorf("Deposit(%d): got error %v, want ErrInvalidAmount", amount, err)
		}
		if _, err := a.Withdraw(amount); !errors.Is(err, ErrInvalidAmount) {
			t.Errorf("Withdraw(%d): got error %v, want ErrInvalidAmount", amount, err)
		}
	}
	if got := a.Balance(); got != 10 {
		t.Errorf("Balance() = %d, want 10", got)
	}
}

func TestAccountConcurrentWithdrawalsNeverOverdraw(t *testing.T) {
	watchdog(t)
	// 8 goroutines try to take 40,000 in total from an account holding
	// 10,000. Exactly 10,000 withdrawals of 1 can succeed. An overdraw
	// can only happen in the last moments, when the balance is almost
	// gone -- so one round can miss it, and the test runs 20.
	const rounds = 20
	const goroutines = 8
	const triesPerGoroutine = 5000
	for round := 0; round < rounds; round++ {
		a := newAccount(t, 10000)
		var succeeded atomic.Int64
		runTogether(t, goroutines, func(int) {
			for i := 0; i < triesPerGoroutine; i++ {
				if ok, _ := a.Withdraw(1); ok {
					succeeded.Add(1)
				}
			}
		})

		if got := succeeded.Load(); got != 10000 {
			t.Fatalf("round %d: %d withdrawals of 1 succeeded from a balance of 10000, want exactly 10000", round, got)
		}
		if got := a.Balance(); got != 0 {
			t.Fatalf("round %d: Balance() = %d, want 0", round, got)
		}
	}
}

func TestAccountConcurrentDepositsAreNeverLost(t *testing.T) {
	watchdog(t)
	const goroutines = 8
	const perGoroutine = 20000
	a := newAccount(t, 0)
	runTogether(t, goroutines, func(int) {
		for i := 0; i < perGoroutine; i++ {
			_ = a.Deposit(1)
		}
	})

	if got := a.Balance(); got != goroutines*perGoroutine {
		t.Errorf("Balance() = %d after %d deposits of 1, want %d -- deposits were lost",
			got, goroutines*perGoroutine, goroutines*perGoroutine)
	}
}

func TestAccountConcurrentDepositsAndWithdrawalsBalanceOut(t *testing.T) {
	watchdog(t)
	// 4 goroutines deposit 3 at a time, 4 goroutines withdraw 2 at a
	// time, and one watcher checks that the balance is never seen below
	// zero.
	const perGoroutine = 20000
	a := newAccount(t, 0)
	var withdrawn atomic.Int64
	var sawNegative atomic.Bool
	var workersLeft atomic.Int32
	workersLeft.Store(8)
	runTogether(t, 9, func(g int) {
		if g == 8 { // the watcher
			for workersLeft.Load() > 0 {
				if a.Balance() < 0 {
					sawNegative.Store(true)
				}
			}
			return
		}
		for i := 0; i < perGoroutine; i++ {
			if g < 4 {
				_ = a.Deposit(3)
			} else if ok, _ := a.Withdraw(2); ok {
				withdrawn.Add(2)
			}
		}
		workersLeft.Add(-1)
	})

	if sawNegative.Load() {
		t.Error("the watcher saw a negative balance")
	}
	want := 4*perGoroutine*3 - withdrawn.Load()
	if got := a.Balance(); got != want {
		t.Errorf("Balance() = %d, want %d (all deposits minus the withdrawals that succeeded)", got, want)
	}
}
