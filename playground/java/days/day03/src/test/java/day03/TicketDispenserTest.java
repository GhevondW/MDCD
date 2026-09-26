package day03;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Timeout;

import java.util.Arrays;
import java.util.concurrent.atomic.AtomicBoolean;

import static day03.RunTogether.runTogether;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Timeout(value = 60, threadMode = Timeout.ThreadMode.SEPARATE_THREAD)
class TicketDispenserTest {

    @Test
    void firstTicketIsOne() {
        TicketDispenser d = new TicketDispenser();
        assertEquals(1, d.next());
    }

    @Test
    void ticketsCountUpByOne() {
        TicketDispenser d = new TicketDispenser();
        assertEquals(1, d.next());
        assertEquals(2, d.next());
        assertEquals(3, d.next());
    }

    @Test
    void issuedCountsHandedOutTickets() {
        TicketDispenser d = new TicketDispenser();
        assertEquals(0, d.issued());
        d.next();
        d.next();
        assertEquals(2, d.issued());
    }

    @Test
    void concurrentTicketsAreUniqueAndGapless() throws Exception {
        int threads = 8;
        int perThread = 20_000;
        TicketDispenser d = new TicketDispenser();
        long[][] got = new long[threads][perThread];
        runTogether(threads, t -> {
            for (int i = 0; i < perThread; i++) got[t][i] = d.next();
        });

        long[] all = Arrays.stream(got).flatMapToLong(Arrays::stream).sorted().toArray();
        for (int i = 0; i < all.length; i++) {
            assertEquals(i + 1, all[i], "a ticket was handed out twice, or a number was skipped");
        }
        assertEquals(threads * perThread, d.issued());
    }

    @Test
    void eachThreadSeesItsTicketsIncrease() throws Exception {
        int threads = 8;
        int perThread = 20_000;
        TicketDispenser d = new TicketDispenser();
        long[][] got = new long[threads][perThread];
        runTogether(threads, t -> {
            for (int i = 0; i < perThread; i++) got[t][i] = d.next();
        });

        for (long[] g : got) {
            for (int i = 1; i < g.length; i++) {
                assertTrue(g[i - 1] < g[i], "a later ticket had a smaller number");
            }
        }
    }

    @Test
    void issuedNeverGoesBackwardsWhileTicketsAreTaken() throws Exception {
        int takers = 4;
        int perThread = 20_000;
        TicketDispenser d = new TicketDispenser();
        AtomicBoolean wentBackwards = new AtomicBoolean(false);
        runTogether(takers + 1, t -> {
            if (t == takers) { // the watcher
                long last = 0;
                for (int i = 0; i < 20_000; i++) {
                    long now = d.issued();
                    if (now < last) wentBackwards.set(true);
                    last = now;
                }
                return;
            }
            for (int i = 0; i < perThread; i++) d.next();
        });

        assertFalse(wentBackwards.get(), "issued() went backwards while tickets were taken");
        assertEquals(takers * perThread, d.issued());
    }
}
