---
marp: true
theme: gaia
class: invert
paginate: true
footer: 'Atomics & Concurrency Patterns'
style: |
  @import url('https://fonts.bunny.net/css?family=ibm-plex-sans:400,500,600,700|ibm-plex-sans-condensed:600,700|ibm-plex-mono:400,500,600&display=swap');

  /* Same blueprint palette as Days 2–4. */
  section {
    --line: #6fa8cc;
    --panel: #14385a;
    --accent: #5cc8ea;
    --good: #7ed0a3;
    --bad: #f0a06e;

    --color-background: #0b2136;
    --color-foreground: #eef4f9;
    --color-dimmed: #a9c8de;
    --color-highlight: #5cc8ea;
    --color-background-stripe: rgba(92, 200, 234, .07);

    font-family: 'IBM Plex Sans', Helvetica, Arial, sans-serif;
    font-size: 25px;
    background-color: var(--color-background);
    background-image:
      repeating-linear-gradient(90deg, rgba(111, 168, 204, .06) 0 1px, transparent 1px 48px),
      repeating-linear-gradient(0deg, rgba(111, 168, 204, .06) 0 1px, transparent 1px 48px);
  }

  section.lead h1 {
    font-family: 'IBM Plex Sans Condensed', 'IBM Plex Sans', sans-serif;
    font-size: 56px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.02em;
  }
  section.lead h1::after {
    content: "";
    display: block;
    width: 64px;
    height: 3px;
    margin: 20px auto 0;
    background: var(--accent);
  }
  section.lead h3 { font-family: 'IBM Plex Sans', sans-serif; font-size: 30px; font-weight: 400; opacity: 0.85; }

  h2 {
    font-family: 'IBM Plex Sans Condensed', 'IBM Plex Sans', sans-serif;
    font-size: 34px;
    font-weight: 700;
    position: relative;
  }
  h2::after {
    content: "";
    position: absolute;
    left: 1px;
    bottom: -8px;
    width: 40px;
    height: 3px;
    background: var(--accent);
  }

  code { font-size: 0.8em; font-family: 'IBM Plex Mono', Menlo, monospace; }
  pre, pre code { font-family: 'IBM Plex Mono', Menlo, monospace; }

  /* Long code listings: a little smaller so they fit with their notes. */
  section.dense pre code { font-size: 0.7em; }

  .tag {
    display: inline-block;
    font-family: 'IBM Plex Mono', monospace;
    font-size: 14px;
    font-weight: 500;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: var(--accent);
    background: var(--panel);
    border: 1px solid var(--line);
    border-radius: 3px;
    padding: 3px 10px;
    opacity: 1;
  }

  blockquote {
    margin: 0.6em 0 0;
    padding: 10px 20px;
    background: var(--panel);
    border-left: 4px solid var(--accent);
    border-radius: 0 4px 4px 0;
    opacity: 1;
  }
  blockquote::before, blockquote::after { content: none; }

  section.good { border-left: 6px solid var(--good); }
  section.good h2 { color: var(--good); }
  section.good h2::after { background: var(--good); }

  section.bad { border-left: 6px solid var(--bad); }
  section.bad h2 { color: var(--bad); }
  section.bad h2::after { background: var(--bad); }

  section footer {
    font-family: 'IBM Plex Mono', monospace;
    font-size: 13px;
    letter-spacing: 0.06em;
    text-transform: uppercase;
  }
  section::after {
    content: attr(data-marpit-pagination) "\2009/\2009" attr(data-marpit-pagination-total);
    font-family: 'IBM Plex Mono', monospace;
    font-size: 13px;
  }

  .note { font-size: 18px; color: var(--color-dimmed); margin-top: 1.4em; }
  .note code { font-size: 0.9em; }

  li:has(> blockquote) { list-style: none; }
  li > blockquote { margin: 0.2em 0 0.5em; }

  table { font-size: 20px; border-collapse: collapse; margin-top: 14px; }
  table th, table td { background: transparent !important; border: 0; border-bottom: 1px solid rgba(111, 168, 204, .35); padding: 7px 14px; text-align: left; }
  table th { color: var(--accent) !important; font-weight: 600; }
  section.rosetta table { font-size: 15px; }
  section.rosetta table td, section.rosetta table th { padding: 5px 8px; }
  section.rosetta table code { font-size: 0.95em; }

---

<!-- _class: lead invert -->

# Atomics & Concurrency Patterns

### Modeling, Design & Collaborative Development

---

<span class="tag">Day 3</span>

## The counter — correct, and slow

```cpp
int counter = 0;
std::mutex m;                       // guards counter

void work() {
    for (int i = 0; i < 1000000; ++i) {
        std::lock_guard lk(m);      // lock — maybe sleep, wake up — unlock
        counter++;                  // ...for one tiny ++
    }
}
```

* Day 3: correct — and 10× slower than one thread doing all the work
* A lock, two million times, to guard **one** `int`. Is there anything smaller than a lock?

---

<!-- _class: lead invert -->

# Atomic Operations

### One step that no thread can split

---

## `counter++` in one step

* Day 3: `counter++` is three steps — load, add, store — and another thread can run between them
* > An **atomic operation** happens in **one** step: every other thread sees it not started yet, or finished — never halfway.
* The CPU has instructions for this. For that one instruction, its core keeps the memory to itself — no lock object, no sleeping, no OS
* Same idea elsewhere: Java `AtomicInteger`, Go `atomic.Int64`, C# `Interlocked`, Rust `AtomicI64`, JavaScript `Atomics` — Python has none: use a `Lock`

---

<!-- _class: invert good -->

## ✅ The counter, fixed without a lock

```cpp
std::atomic<int> counter{0};

void work() {
    for (int i = 0; i < 1000000; ++i)
        counter++;                  // one atomic step: fetch_add(1)
}
```

* `2000000` — every run, and no mutex
* `std::atomic<int>` looks like an `int`, but every operation on it is one indivisible step

---

## Two threads, one atomic add

`counter` is 5. Both threads run `counter++` at the same moment.

| | Core 1 · thread A | Core 2 · thread B | `counter` |
|---|---|---|---|
| 1 | holds the memory: read 5, write 6 — one step | waits — a few nanoseconds, in hardware | 6 |
| 2 | | holds the memory: read 6, write 7 | 7 |

* Two adds, +2 — nothing lost
* Each add also gives back the value it replaced: A gets 5, B gets 6 — never the same one twice

---

## The menu: one variable, one step

| Operation | What happens, in one step | Used for |
|---|---|---|
| **load / store** | read or write the whole value — never half of it | flags, ready values |
| **exchange** | write a new value, get the old one back | "take it, leave X" |
| **fetch_add / fetch_sub** | add, get the old value back | counters, tickets, IDs |
| **compare-and-swap** | write a new value **only if** it is still the one I expect | everything else |

* Every one of them touches **one** variable

---

<!-- _class: invert good -->

<span class="tag">Day 3</span>

## ✅ The order numbers, fixed without a lock

```cpp
std::atomic<int> lastId{0};

std::string newOrderId() {
    int id = lastId.fetch_add(1) + 1;     // add, and get the value — one step
    return "ORD-" + std::to_string(id);
}
```

* Day 3's bug: the add and the read of the new value were two steps
* `fetch_add` returns the value it replaced — the add and the read are now **one** step. Every caller gets a different number

---

<!-- _class: invert bad -->

## ❌ An atomic variable — still a race

```cpp
std::string newOrderId() {
    lastId++;                                   // atomic
    return "ORD-" + std::to_string(lastId);     // atomic too
}
```

| | Thread A | Thread B | `lastId` |
|---|---|---|---|
| 1 | `lastId++` | | 1 |
| 2 | | `lastId++` | 2 |
| 3 | reads 2 → `ORD-2` | | 2 |
| 4 | | reads 2 → **`ORD-2` again** | 2 |

* Each line is atomic — the pair is not

---

<!-- _class: lead invert -->

# Compare-and-Swap

### Check-then-act in one step

---

<!-- _class: invert bad -->

## ❌ The withdrawal, with an atomic balance

```cpp
std::atomic<long long> balance{100};

bool withdraw(long long amount) {
    if (balance >= amount) {        // check — atomic
        balance -= amount;          // act — atomic
        return true;
    }
    return false;
}
```

* Day 3's interleaving still works: both threads check 100, both subtract — **−100**
* > An atomic **variable** doesn't make your **code** atomic. Only each single operation is.
* `fetch_sub` can't help — it subtracts without checking. We need the check **and** the act in one step

---

## Compare-and-swap

* > **compare-and-swap(x, expected, new)**: "if `x` is still `expected`, set it to `new` — and tell me if it worked." One step.
* It's check-then-act — done by the hardware, so nothing can get in between the check and the act
* It **fails** if someone changed `x` since you read it. Then: read again, decide again, try again
* It's the universal one: every other atomic operation can be built from it

---

<!-- _class: invert bad dense -->

## ❌ Why not build CAS from `load` and `store`?

```cpp
bool myCas(std::atomic<long long>& x, long long expected, long long desired) {
    if (x.load() != expected) return false;    // check — atomic
    x.store(desired);                          // act — atomic
    return true;
}
```

`balance` is 100; both threads call `myCas(balance, 100, 0)` to withdraw 100.

| | Thread A | Thread B | `balance` |
|---|---|---|---|
| 1 | `load` → 100 — matches | | 100 |
| 2 | | `load` → 100 — matches | 100 |
| 3 | `store(0)` → `true` | | 0 |
| 4 | | `store(0)` → `true` | 0 |

* Both won: 200 paid out of 100. Two atomic steps leave a gap — CAS must be **one** instruction

---

<!-- _class: invert good -->

## ✅ The withdrawal with compare-and-swap

```cpp
bool withdraw(long long amount) {
    long long cur = balance.load();
    while (cur >= amount) {
        if (balance.compare_exchange_weak(cur, cur - amount))
            return true;            // still cur? then cur - amount — done
        // failed: someone changed it; cur now holds the new value
    }
    return false;                   // not enough money
}
```

* The **CAS loop**: read → decide → CAS → if it failed, decide again with the fresh value
* Never negative: a withdrawal only lands if the balance is still the one it checked

<p class="note">C++: when it fails, <code>compare_exchange</code> writes the current value into <code>cur</code>. The <code>_weak</code> version may fail even when the values match — harmless inside a loop.</p>

---

## Two withdrawals of 100 — with CAS

`balance` is 100. Both threads call `withdraw(100)`.

| | Thread A | Thread B | `balance` |
|---|---|---|---|
| 1 | `cur` = 100 | | 100 |
| 2 | | `cur` = 100 | 100 |
| 3 | CAS: still 100? → 0 ✓ | | 0 |
| 4 | | CAS: still 100? **no** — fails, `cur` = 0 | 0 |
| 5 | | 0 ≥ 100? no → `false` | 0 |

* The same start as Day 3 — but the failed CAS hands B the fresh value, and B decides again
* One withdrawal wins; the balance is 0 — never −100

---

<!-- _class: invert dense -->

## The CAS loop: any update to one variable

```cpp
T cur = x.load();
T next;
do {
    next = f(cur);                                 // compute from what you saw
} while (!x.compare_exchange_weak(cur, next));     // still what I saw? store it. No? again
```

Example — the highest latency seen so far:

```cpp
std::atomic<long> peak{0};

void recordLatency(long ms) {
    long cur = peak.load();
    while (ms > cur && !peak.compare_exchange_weak(cur, ms)) {}
}
```

* Works for max, min, "add but never above 100", "set if not set yet" — anything computed from **one** variable

---

## Optimistic vs. pessimistic

* A **mutex is pessimistic**: "someone may get in my way — lock first, then work"
* A **CAS loop is optimistic**: "nobody will — work first, then check; if I was wrong, redo it"
* Few threads in a fight → optimistic wins: nobody waits, nobody sleeps
* Many threads in a fight → optimistic loses: they redo the same work again and again. A CAS loop never sleeps — it spins, like Day 4's busy waiting

---

## The same idea, far from C++

* **Databases** — optimistic locking:
  `UPDATE accounts SET balance = 0, version = 8 WHERE id = 1 AND version = 7`
  0 rows changed? Someone was faster: read again, try again
* **HTTP** — `PUT` with `If-Match: "v7"` → `412 Precondition Failed` if it changed in between
* **Git** — `git push` is rejected if the remote moved since you fetched: pull, try again
* > Compare-and-swap isn't a C++ feature. It's how to change shared state **without holding a lock**.

---

<!-- _class: lead invert -->

# What an Atomic Costs

### Cache lines

---

## Why one shared atomic gets slow

* Memory travels between RAM and the cores' caches in blocks: **cache lines** — 64 bytes on most CPUs
* To write, a core must hold the line **alone** — every other core's copy is thrown away
* One counter, many cores adding: the line travels from core to core, on every add
* > An atomic doesn't remove the fight — it makes it cheaper. Many threads writing one variable is still a line, in hardware.

---

<!-- _class: invert good -->

## ✅ Don't share it: a counter per thread

```cpp
struct alignas(64) Slot { std::atomic<long> n{0}; };   // one cache line each
Slot perThread[8];

void hit(int t) { perThread[t].n++; }         // stays in this core's cache

long total() {                                // rare: add them up
    long sum = 0;
    for (auto& s : perThread) sum += s.n;
    return sum;
}
```

* Writes stay on their own core; only reading the total touches every line — fine for counters and statistics
* Without `alignas(64)`: separate counters, **one** line — it travels anyway. That's **false sharing**
* Same idea elsewhere: Java `LongAdder` — and Day 4's queue per worker

---

<!-- _class: lead invert -->

# Building With Atomics

### Flags, locks, reference counts

---

## What can we build with atomics?

| We build… | From |
|---|---|
| counters, statistics | `fetch_add` |
| unique IDs, ticket numbers | `fetch_add` — the old value it returns |
| a stop flag | an atomic `bool` |
| a spinlock — and the fast path of every mutex | `exchange` |
| reference counts — `shared_ptr`, Rust `Arc` | `fetch_add` / `fetch_sub` |
| max, min, any "update if" | a CAS loop |

---

<!-- _class: invert bad -->

## ❌ A stop flag — a plain `bool`

```cpp
bool running = true;               // shared, no lock

void worker() {
    while (running) step();
}

// main thread:
running = false;                   // "please stop"
t.join();
```

* One thread writes, another reads, no lock: a **data race**. So the compiler may assume nobody else changes `running` — read it **once**, and loop forever. Optimized builds really do this
* ✅ `std::atomic<bool> running{true};` — every check really reads it
* Same idea elsewhere: Java `AtomicBoolean` or `volatile`, Go `atomic.Bool`, Rust `AtomicBool`, C# `CancellationToken`

<p class="note">C++ <code>volatile</code> is not the fix — in C++ it makes nothing atomic. Java's <code>volatile</code> is a different thing, and does work here.</p>

---

<span class="tag">Your turn</span>

## Challenge: build a mutex

Only today's tools — atomic `load`, `store`, `exchange`, compare-and-swap. No `std::mutex`.

```cpp
class Mutex {
    // your state here
public:
    void lock();      // returns only when this thread owns the lock
    void unlock();    // lets another thread in
};
```

* Two threads must never be between `lock()` and `unlock()` at the same time
* What is the state? What does `lock()` do when the lock is already taken?
* "If it's free, take it" — what does that look like as **one** step?

---

<!-- _class: invert good dense -->

## ✅ A simple mutex

```cpp
class Mutex {
    std::atomic<bool> locked{false};
public:
    void lock() {
        bool expected = false;
        while (!locked.compare_exchange_weak(expected, true))   // free? then take it — one step
            expected = false;                                   // it was taken: try again
    }
    void unlock() { locked.store(false); }                      // free again
};
```

* "If free, take it" is check-then-act — so it must be **one** CAS. A `load` + `store` lets two threads in
* Also works: `while (locked.exchange(true)) {}` — "set it taken; was it taken before?"
* While it's taken, `lock()` spins — a **spinlock**: fine for a few instructions; Day 4: don't spin for long
* A real mutex = this CAS for the fast path + Day 4's waiting room when it's taken

---

## Reference counts — atomics you already use

```cpp
void retain()  { refs.fetch_add(1); }

void release() {
    if (refs.fetch_sub(1) == 1)       // the old value was 1: I was the last one
        delete this;
}
```

* Every copy of a `std::shared_ptr` does `fetch_add`; every destructor does `fetch_sub`
* "Was I the last?" must be decided in the **same** step as the decrement — as two steps, two threads could both see 0 and delete twice
* Same idea elsewhere: Rust `Arc`, Swift ARC, Objective-C, COM `AddRef` / `Release`

---

## Atomic or mutex?

| Situation | Use |
|---|---|
| one variable, one simple update — count, flag, ID, max | **atomic** |
| a rule that ties two or more variables together | **mutex** |
| a whole structure — map, queue, list | **mutex** — or a tested library |
| this lock is the bottleneck, and you measured it | **shard it** — a counter or a queue per thread |

* > When in doubt, a mutex: clear and correct first, fast second.

---

<!-- _class: invert rosetta -->

## The same atomics, in every language

| | C++ | Java | Go | Python | C# | Rust |
|---|---|---|---|---|---|---|
| atomic integer | `atomic<long>` | `AtomicLong` | `atomic.Int64` | — use a `Lock` | `Interlocked` on a `long` | `AtomicI64` |
| add | `fetch_add` | `getAndAdd` | `Add` | — | `Interlocked.Add` | `fetch_add` |
| compare-and-swap | `compare_exchange_weak` | `compareAndSet` | `CompareAndSwap` | — | `Interlocked.CompareExchange` | `compare_exchange` |
| flag | `atomic<bool>` | `AtomicBoolean` | `atomic.Bool` | `threading.Event` | `CancellationToken` | `AtomicBool` |
| pointer / reference | `atomic<T*>` | `AtomicReference` | `atomic.Pointer[T]` | — | `Interlocked.CompareExchange` | `AtomicPtr` |
| counter per thread | build it | `LongAdder` | build it | — | build it | build it |

---

<!-- _class: lead invert -->

# Concurrency Design Patterns

### Primitives are the parts — these are the shapes

---

## Same problem, different shapes

* So far: **threads + locks** — share the memory, guard it (Day 3) · **pool + futures** — share the work, get results back (Day 4)
* Every shape answers one question: *who may touch the state?*
* Three more shapes today: **actors**, **the event loop**, **channels**
* The trick in all three: stop sharing mutable state — then there is nothing to lock

---

## Actors: state with one owner

* > An **actor** owns its state — nobody else touches it. To change it, you send it a **message**; it handles its messages one at a time.
* The mailbox is Day 4's blocking queue — one per actor, one thread reading it
* No locks inside: only one thread ever touches the state
* Same idea elsewhere: Erlang / Elixir processes, Akka (Java, Scala), Orleans (C#), Swift `actor`

---

<!-- _class: invert good -->

## ✅ A logger as an actor

```cpp
class Logger {                                    // an actor
    BlockingQueue<std::string> mailbox{1000};     // Day 4's queue
    std::thread worker{[this] {
        while (auto line = mailbox.pop())         // one message at a time
            std::cout << *line << '\n';           // only this thread touches cout
    }};
public:
    void log(std::string line) { mailbox.push(std::move(line)); }   // any thread: just send
    ~Logger() { mailbox.close(); worker.join(); }                   // finish the mailbox, stop
};
```

* Any thread calls `logger.log("...")` — lines never mix, and there is no mutex around `cout`

---

## Actors: what can still go wrong

* Two actors each wait for the other's reply → **deadlock** — with no lock at all, like Day 4's pool
* A full mailbox: bounded → the sender waits (backpressure); unbounded → memory grows (Day 4: why bounded?)
* One slow message holds up every message behind it — keep handlers short
* Asking an actor for a value is a message too: the answer is already old when you read it

---

## Event loop: one thread, many waits

* Most server work is **waiting**: for the network, the disk, the database
* A thread per request: 10,000 requests → 10,000 threads, almost all asleep (Day 4)
* > **Event loop**: **one** thread and a queue of ready tasks. A task runs until it would wait — then it gives the thread back, and returns to the queue when its data arrives.
* Many tasks in flight, one running at a time: concurrency **without** parallelism
* Same idea elsewhere: JavaScript / Node.js, Python `asyncio`, C# `async` / `await`, Rust tokio, nginx

---

## The loop itself

```text
ready = queue of tasks that can run now

loop forever:
    task = ready.pop()
    run task until it reaches an await          // it would have to wait
    when its data arrives — network, disk, timer:
        ready.push(task)                        // it can go on
```

* Day 4's queue again — one consumer, and the producers are "your data is here"
* `async` / `await`: code that reads top to bottom; every `await` is a place where the task gives the thread back

---

## Channels: share memory by communicating

* > Go: "Don't communicate by sharing memory; share memory by communicating."
* A **channel** is Day 4's blocking queue, built into the language: send waits while it's full, receive waits while it's empty
* Goroutines don't share the data — they pass it along. Whoever holds a value owns it
* Same idea elsewhere: Rust `mpsc`, Kotlin `Channel`, C# `System.Threading.Channels`, Clojure `core.async`

---

<!-- _class: invert good -->

## ✅ A channel (Go)

```go
ch := make(chan int, 10)          // Day 4's bounded queue, built in

go func() {                       // producer
    for i := 1; i <= 3; i++ {
        ch <- i                   // send: waits while the channel is full
    }
    close(ch)                     // "no more values"
}()

for v := range ch {               // receive: waits while it's empty; ends after close
    fmt.Println(v)
}
```

* No mutex, no condition variable — the channel does Day 4's whole queue: waiting, backpressure, `close()`

---

## `select` — wait on several channels at once

```go
select {
case msg := <-requests:
    handle(msg)                          // a request came in
case <-time.After(time.Second):
    fmt.Println("nothing for 1 s")       // time out
case <-done:
    return                               // shutdown signal
}
```

* Sleeps until **one** of them is ready — no busy waiting
* Day 4's three answers to "empty": block (receive), time out (`time.After`), balk (`default:`)

---

<!-- _class: invert rosetta -->

## Same problems, different shapes

| | The state is… | Threads wait on… | Typical bug | Good for |
|---|---|---|---|---|
| **threads + locks** | shared, behind locks | locks | races, deadlock | small shared structures |
| **pool + futures** | shared tasks; results come back | `future.get()` | waiting inside the pool | CPU work, many small jobs |
| **event loop** | owned by one thread | nothing — tasks `await` | blocking the loop; races across `await` | many network waits |
| **actors** | owned by each actor | its mailbox | actors waiting for each other | many independent things — accounts, players |
| **channels** | passed along, one owner at a time | send / receive | nobody closes; nobody receives | pipelines: stage → stage |

* > Your team project: which shape, and why? You will defend it.

---

<!-- _class: lead invert -->

# Pick the Shape

### One variable: an atomic · a few: a mutex · the whole program: a pattern
