package day03

// Day 3 -- Races: the broken counter.
//
// A ticket machine at a bank: every customer who walks in takes the next
// number. Many goroutines take tickets at the same time.
//
//	Next()    hand out the next ticket: 1 for the first call, then 2,
//	          3, ... No two calls may ever get the same number, and no
//	          number may be skipped.
//	Issued()  how many tickets have been handed out so far.
//
// Both may be called by many goroutines at the same time.
//
// Remember the lecture: counter++ is load -> add -> store, three steps
// another goroutine can slip between. Next() has exactly that shape.

type TicketDispenser struct {
	// TODO: choose your own representation.
}

func NewTicketDispenser() *TicketDispenser {
	// TODO
	return &TicketDispenser{}
}

func (d *TicketDispenser) Next() int64 {
	// TODO
	return 0
}

func (d *TicketDispenser) Issued() int64 {
	// TODO
	return 0
}
