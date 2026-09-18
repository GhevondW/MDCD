import unittest

from .bounded_stack import BoundedStack


class BoundedStackTest(unittest.TestCase):
    def test_starts_empty(self):
        s = BoundedStack(2)
        self.assertTrue(s.empty())
        self.assertEqual(s.size(), 0)
        self.assertFalse(s.full())

    def test_push_up_to_capacity(self):
        s = BoundedStack(2)
        s.push(5)
        s.push(7)
        self.assertEqual(s.size(), 2)
        self.assertTrue(s.full())

    def test_push_beyond_capacity_raises(self):
        s = BoundedStack(1)
        s.push(1)
        with self.assertRaises(IndexError):
            s.push(2)
        self.assertEqual(s.size(), 1)  # rejected push must not touch state

    def test_pop_returns_lifo_order(self):
        s = BoundedStack(3)
        s.push(1)
        s.push(2)
        s.push(3)
        self.assertEqual(s.pop(), 3)
        self.assertEqual(s.pop(), 2)
        self.assertEqual(s.size(), 1)

    def test_pop_empty_raises(self):
        s = BoundedStack(1)
        with self.assertRaises(IndexError):
            s.pop()

    def test_top_does_not_remove(self):
        s = BoundedStack(2)
        s.push(9)
        self.assertEqual(s.top(), 9)
        self.assertEqual(s.size(), 1)

    def test_top_empty_raises(self):
        s = BoundedStack(1)
        with self.assertRaises(IndexError):
            s.top()

    def test_zero_capacity_is_always_full(self):
        s = BoundedStack(0)
        self.assertTrue(s.full())
        with self.assertRaises(IndexError):
            s.push(1)

    def test_pop_frees_capacity_for_another_push(self):
        s = BoundedStack(1)
        s.push(1)
        self.assertTrue(s.full())
        s.pop()
        self.assertFalse(s.full())
        self.assertTrue(s.empty())
        s.push(2)  # must succeed -- capacity was freed, not stuck full
        self.assertEqual(s.top(), 2)
        self.assertEqual(s.size(), 1)

    def test_top_reflects_most_recent_push(self):
        s = BoundedStack(3)
        s.push(1)
        self.assertEqual(s.top(), 1)
        s.push(2)
        self.assertEqual(s.top(), 2)
        s.push(3)
        self.assertEqual(s.top(), 3)

    def test_empty_after_popping_everything(self):
        s = BoundedStack(2)
        s.push(1)
        s.push(2)
        s.pop()
        s.pop()
        self.assertTrue(s.empty())
        self.assertEqual(s.size(), 0)
        with self.assertRaises(IndexError):
            s.pop()


if __name__ == "__main__":
    unittest.main()
