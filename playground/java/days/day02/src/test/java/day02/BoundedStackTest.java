package day02;

import org.junit.jupiter.api.Test;

import java.util.NoSuchElementException;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class BoundedStackTest {

    @Test
    void startsEmpty() {
        BoundedStack s = new BoundedStack(2);
        assertTrue(s.isEmpty());
        assertEquals(0, s.size());
        assertFalse(s.isFull());
    }

    @Test
    void pushUpToCapacity() {
        BoundedStack s = new BoundedStack(2);
        s.push(5);
        s.push(7);
        assertEquals(2, s.size());
        assertTrue(s.isFull());
    }

    @Test
    void pushBeyondCapacityThrows() {
        BoundedStack s = new BoundedStack(1);
        s.push(1);
        assertThrows(IllegalStateException.class, () -> s.push(2));
        assertEquals(1, s.size()); // rejected push must not touch state
    }

    @Test
    void popReturnsLifoOrder() {
        BoundedStack s = new BoundedStack(3);
        s.push(1);
        s.push(2);
        s.push(3);
        assertEquals(3, s.pop());
        assertEquals(2, s.pop());
        assertEquals(1, s.size());
    }

    @Test
    void popEmptyThrows() {
        BoundedStack s = new BoundedStack(1);
        assertThrows(NoSuchElementException.class, s::pop);
    }

    @Test
    void topDoesNotRemove() {
        BoundedStack s = new BoundedStack(2);
        s.push(9);
        assertEquals(9, s.top());
        assertEquals(1, s.size());
    }

    @Test
    void topEmptyThrows() {
        BoundedStack s = new BoundedStack(1);
        assertThrows(NoSuchElementException.class, s::top);
    }

    @Test
    void zeroCapacityIsAlwaysFull() {
        BoundedStack s = new BoundedStack(0);
        assertTrue(s.isFull());
        assertThrows(IllegalStateException.class, () -> s.push(1));
    }

    @Test
    void popFreesCapacityForAnotherPush() {
        BoundedStack s = new BoundedStack(1);
        s.push(1);
        assertTrue(s.isFull());
        s.pop();
        assertFalse(s.isFull());
        assertTrue(s.isEmpty());
        s.push(2); // must succeed -- capacity was freed, not stuck full
        assertEquals(2, s.top());
        assertEquals(1, s.size());
    }

    @Test
    void topReflectsMostRecentPush() {
        BoundedStack s = new BoundedStack(3);
        s.push(1);
        assertEquals(1, s.top());
        s.push(2);
        assertEquals(2, s.top());
        s.push(3);
        assertEquals(3, s.top());
    }

    @Test
    void emptyAfterPoppingEverything() {
        BoundedStack s = new BoundedStack(2);
        s.push(1);
        s.push(2);
        s.pop();
        s.pop();
        assertTrue(s.isEmpty());
        assertEquals(0, s.size());
        assertThrows(NoSuchElementException.class, s::pop);
    }
}
