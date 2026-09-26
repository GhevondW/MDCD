package day03;

/**
 * Day 3 -- Races: the broken counter.
 *
 * A ticket machine at a bank: every customer who walks in takes the next
 * number. Many threads take tickets at the same time.
 *
 *   next()    hand out the next ticket: 1 for the first call, then 2,
 *             3, ... No two calls may ever get the same number, and no
 *             number may be skipped.
 *   issued()  how many tickets have been handed out so far.
 *
 * Both may be called by many threads at the same time.
 *
 * Remember the lecture: counter++ is load -> add -> store, three steps
 * another thread can slip between. next() has exactly that shape.
 */
public class TicketDispenser {
    // TODO: choose your own representation.

    public long next() {
        // TODO
        return 0;
    }

    public long issued() {
        // TODO
        return 0;
    }
}
