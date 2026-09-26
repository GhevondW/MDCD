package day03

import (
	"errors"
	"math/rand"
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

// ---- Transfer and Total: two accounts locked at once ----

func TestAccountTransferMovesMoney(t *testing.T) {
	watchdog(t)
	a, b := newAccount(t, 100), newAccount(t, 50)
	if ok, err := a.Transfer(b, 30); !ok || err != nil {
		t.Fatalf("Transfer(b, 30) = (%v, %v), want (true, nil)", ok, err)
	}
	if a.Balance() != 70 || b.Balance() != 80 {
		t.Errorf("balances = %d, %d, want 70, 80", a.Balance(), b.Balance())
	}
}

func TestAccountTransferWithoutEnoughMoneyChangesNothing(t *testing.T) {
	watchdog(t)
	a, b := newAccount(t, 10), newAccount(t, 0)
	if ok, err := a.Transfer(b, 11); ok || err != nil {
		t.Fatalf("Transfer(b, 11) = (%v, %v), want (false, nil)", ok, err)
	}
	if a.Balance() != 10 || b.Balance() != 0 {
		t.Errorf("balances = %d, %d, want 10, 0 -- a failed transfer changed something", a.Balance(), b.Balance())
	}
}

func TestAccountTransferRejectsBadArguments(t *testing.T) {
	watchdog(t)
	a, b := newAccount(t, 10), newAccount(t, 0)
	for _, amount := range []int64{0, -5} {
		if _, err := a.Transfer(b, amount); !errors.Is(err, ErrInvalidAmount) {
			t.Errorf("Transfer(b, %d): want ErrInvalidAmount, got %v", amount, err)
		}
	}
	if _, err := a.Transfer(a, 1); !errors.Is(err, ErrSameAccount) {
		t.Errorf("a transfer to the same account: want ErrSameAccount, got %v", err)
	}
	if a.Balance() != 10 || b.Balance() != 0 {
		t.Errorf("balances = %d, %d, want 10, 0", a.Balance(), b.Balance())
	}
}

func TestAccountTotalAddsBothBalances(t *testing.T) {
	watchdog(t)
	a, b := newAccount(t, 100), newAccount(t, 50)
	for _, pair := range [][2]*Account{{a, b}, {b, a}} {
		if got, err := Total(pair[0], pair[1]); got != 150 || err != nil {
			t.Errorf("Total = (%d, %v), want (150, nil)", got, err)
		}
	}
	if _, err := Total(a, a); !errors.Is(err, ErrSameAccount) {
		t.Errorf("Total(a, a): want ErrSameAccount, got %v", err)
	}
}

func TestAccountOppositeTransfersDoNotDeadlock(t *testing.T) {
	watchdog(t)
	// Half the goroutines move money from x to y, the other half from y
	// to x -- the slides' "Two transfers, two locks", 80,000 times.
	const goroutines = 4
	const perGoroutine = 20000
	x, y := newAccount(t, 1000), newAccount(t, 1000)
	runTogether(t, goroutines, func(g int) {
		for i := 0; i < perGoroutine; i++ {
			if g%2 == 0 {
				_, _ = x.Transfer(y, 1)
			} else {
				_, _ = y.Transfer(x, 1)
			}
		}
	})

	if sum := x.Balance() + y.Balance(); sum != 2000 {
		t.Errorf("balances add up to %d, want 2000 -- money appeared or disappeared", sum)
	}
	if x.Balance() < 0 || y.Balance() < 0 {
		t.Errorf("balances = %d, %d -- one went below zero", x.Balance(), y.Balance())
	}
}

func TestAccountTotalNeverSeesAHalfDoneTransfer(t *testing.T) {
	watchdog(t)
	// Transfers move money back and forth between x and y while a watcher
	// keeps asking for the total. A transfer done as two steps (take from
	// x, then give to y) shows the watcher money in flight.
	const movers = 4
	const perGoroutine = 20000
	x, y := newAccount(t, 1000), newAccount(t, 1000)
	var moversLeft atomic.Int32
	moversLeft.Store(movers)
	var wrongTotal atomic.Int64
	wrongTotal.Store(2000)
	runTogether(t, movers+1, func(g int) {
		if g == movers { // the watcher
			for moversLeft.Load() > 0 {
				if seen, _ := Total(x, y); seen != 2000 {
					wrongTotal.Store(seen)
				}
			}
			return
		}
		for i := 0; i < perGoroutine; i++ {
			if g%2 == 0 {
				_, _ = x.Transfer(y, 7)
			} else {
				_, _ = y.Transfer(x, 7)
			}
		}
		moversLeft.Add(-1)
	})

	if got := wrongTotal.Load(); got != 2000 {
		t.Errorf("Total saw %d, want 2000 every time -- it saw money that was in flight", got)
	}
	if got, _ := Total(x, y); got != 2000 {
		t.Errorf("Total(x, y) = %d at the end, want 2000", got)
	}
}

func TestAccountConcurrentTransfersAmongManyAccountsKeepTheMoney(t *testing.T) {
	watchdog(t)
	// 8 goroutines move random amounts between random pairs of 6
	// accounts, in both directions. No deadlock, no money made or lost,
	// no account below zero.
	const accounts = 6
	const goroutines = 8
	const perGoroutine = 20000
	accs := make([]*Account, accounts)
	for i := range accs {
		accs[i] = newAccount(t, 1000)
	}
	runTogether(t, goroutines, func(g int) {
		rng := rand.New(rand.NewSource(int64(99 + g)))
		for i := 0; i < perGoroutine; i++ {
			from := rng.Intn(accounts)
			to := rng.Intn(accounts - 1)
			if to >= from {
				to++ // any account but from
			}
			_, _ = accs[from].Transfer(accs[to], 1+int64(rng.Intn(50)))
		}
	})

	var sum int64
	for i, a := range accs {
		if a.Balance() < 0 {
			t.Errorf("account %d went below zero: %d", i, a.Balance())
		}
		sum += a.Balance()
	}
	if sum != 1000*accounts {
		t.Errorf("balances add up to %d, want %d -- money appeared or disappeared", sum, 1000*accounts)
	}
}
