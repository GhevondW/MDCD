---
marp: true
theme: gaia
class: invert
paginate: true
footer: 'Primitives, Deadlock & Async Patterns'
style: |
  @import url('https://fonts.bunny.net/css?family=ibm-plex-sans:400,500,600,700|ibm-plex-sans-condensed:600,700|ibm-plex-mono:400,500,600&display=swap');

  /* Same blueprint palette as Days 2–3. */
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

  /* Interactive diagrams — drawn by the script on the last slide.
     Step diagrams move one step per keypress; live ones have buttons. */
  .viz { display: block; position: relative; margin-top: 12px; }
  .viz > svg { display: block; overflow: visible; }
  .viz .noanim * { transition: none !important; }
  .viz .mv { transition: transform .55s cubic-bezier(.45, .05, .25, 1), opacity .35s; }
  .viz .fd { transition: opacity .35s; }
  .viz .col { transition: fill .35s, stroke .35s, opacity .35s; }
  .viz .bar { transform-box: fill-box; transform-origin: left center; transition: transform .9s ease; }
  .viz .flash { animation: viz-flash .9s ease-out; }
  .viz .pulse { animation: viz-pulse 1.1s ease-in-out infinite; }
  @keyframes viz-flash { 0% { opacity: .15; } 35% { opacity: 1; } }
  @keyframes viz-pulse { 50% { opacity: .35; } }
  .viz-step {
    position: absolute; top: 34px; right: 70px;
    font: 500 13px 'IBM Plex Mono', monospace; letter-spacing: .06em;
    color: var(--color-dimmed); opacity: .8;
  }
  .viz-ctl {
    display: flex; flex-wrap: wrap; gap: 8px 16px; align-items: center; margin-top: 10px;
    font: 14px 'IBM Plex Mono', monospace; color: var(--color-dimmed);
  }
  .viz-ctl .grp { display: inline-flex; gap: 0; }
  .viz-ctl .grp button { border-radius: 0; margin-left: -1px; }
  .viz-ctl .grp button:first-child { border-radius: 3px 0 0 3px; margin-left: 0; }
  .viz-ctl .grp button:last-child { border-radius: 0 3px 3px 0; }
  .viz-ctl button {
    font: 500 14px 'IBM Plex Mono', monospace; color: var(--color-foreground);
    background: var(--panel); border: 1px solid var(--line); border-radius: 3px;
    padding: 4px 12px; cursor: pointer;
  }
  .viz-ctl button:hover { border-color: var(--accent); }
  .viz-ctl button.on { background: var(--accent); border-color: var(--accent); color: #0b2136; }
  .viz-ctl label { display: inline-flex; gap: 8px; align-items: center; }
  .viz-ctl input[type=range] { width: 110px; accent-color: #5cc8ea; }
  .viz-ctl .val { color: var(--color-foreground); min-width: 2.6em; }

  /* Step captions: one list item per step; only the newest one shows. */
  .steps ul { display: grid; padding: 0; margin: 12px 0 0; }
  .steps li { grid-area: 1 / 1; list-style: none; font-size: 23px; }
  .steps li[data-bespoke-marp-fragment=active]:has(~ li[data-bespoke-marp-fragment=active]) { visibility: hidden; }
---

<!-- _class: lead invert -->

# Primitives, Deadlock & Async Patterns

### Modeling, Design & Collaborative Development

---

## Producer and consumer — busy waiting

```cpp
std::mutex m;                            // guards jobs

void producer(Job j) {
    m.lock(); jobs.push(j); m.unlock();
}

void consumer() {
    for (;;) {
        m.lock();
        if (jobs.empty()) { m.unlock(); continue; }   // nothing yet: look again… and again
        Job j = jobs.front(); jobs.pop();
        m.unlock();
        run(j);
    }
}
```

* Day 3's job queue. It's correct — but what is wrong with it?

---

## Two kinds of synchronization

<div class="viz" data-viz="kinds"></div>

<div class="steps">

* **Exclusion** — "not at the same time". A is inside; B waits at the door.
* A leaves, B enters — never both at once. That was Day 3: the mutex.
* **Coordination** — "not before it's ready". Nothing to take yet, so the consumer sleeps.
* The producer puts an item and wakes it up. That is today.

</div>

---

## From today: locks that unlock themselves

```cpp
void deposit(long long amount) {
    std::lock_guard lk(m);              // locks here
    if (amount <= 0) return;            // unlocks here...
    balance += amount;
}                                       // ...or here — even on an exception
```

* `lock_guard` locks when it's created and unlocks when the scope ends — no `unlock()` to forget. The code on the next slides uses it.
* Same idea elsewhere: Java `synchronized`, Python `with lock:`, Go `defer mu.Unlock()`, C# `lock (o) { }`

---

<!-- _class: lead invert -->

# Waiting

### Condition variables

---

## Condition variable: a waiting room next to the lock

<div class="viz" data-viz="cvRoom"></div>

---

## The idea: a waiting room next to the lock

* > A **condition variable** is a waiting room attached to a lock. A thread that can't go on sleeps there until another thread says: "something changed — look again".
* **wait** — give the lock back and fall asleep, in one step; on waking up, take the lock again
* **notify** — wake one sleeper · **notify all** — wake every sleeper
* The **condition** itself — "the queue is not empty" — lives in your data, under the lock. The waiting room only carries the wake-up.
* Every language has one: C++ `condition_variable`, Java `wait` / `notify`, Python `Condition`, Go `sync.Cond`, C# `Monitor.Wait`

---

<!-- _class: invert good -->

## ✅ Producer and consumer — waiting on a condition variable

```cpp
std::condition_variable cv;              // the waiting room, next to m

void producer(Job j) {
    { std::lock_guard lk(m); jobs.push(j); }
    cv.notify_one();                     // "something changed — look again"
}

void consumer() {
    for (;;) {
        std::unique_lock lk(m);
        while (jobs.empty()) cv.wait(lk); // nothing yet: sleep, no CPU
        Job j = jobs.front();
        jobs.pop();
        lk.unlock();
        run(j);
    }
}
```

* The consumer sleeps until there is work — and the producer is the one who wakes it

---

## Spin, sleep, or wait — try it

<div class="viz" data-viz="waitModes"></div>

---

## What can we build with this idea?

Every waiting tool is the same recipe: **some state + a lock + a condition variable**

| We build… | A thread waits until… |
|---|---|
| **Bounded blocking queue** | there is room / there is an item |
| **Semaphore** | a permit is free |
| **Latch** | the count reaches 0 |
| **Barrier** | everyone has arrived |
| **Future** | the value is there |
| **Thread pool** | a task arrives |

* > The rest of today: these tools — and what goes wrong when threads wait for each other.

---

<!-- _class: lead invert -->

# The Bounded Blocking Queue

### The structure to know cold

---

## Why bounded? Try it

<div class="viz" data-viz="whyBounded"></div>

- Producer faster than consumer, **unbounded** queue: memory grows until the process dies
- **Bounded**: when it's full, the producer waits — **backpressure**: a full queue slows down whoever fills it

---

<!-- _class: invert good dense -->

## ✅ The bounded blocking queue — the whole thing

```text
BoundedQueue:  items, capacity, mutex m, conditions notFull and notEmpty

push(x):
    lock(m)
    while items is full:
        wait(notFull, m)          // full: unlock m and sleep; when woken, lock m again
    add x to items
    unlock(m)
    notify(notEmpty)              // wake one sleeping consumer

pop():
    lock(m)
    while items is empty:
        wait(notEmpty, m)         // empty: unlock m and sleep; when woken, lock m again
    x = take the first item
    unlock(m)
    notify(notFull)               // wake one sleeping producer
    return x
```

---

<!-- _class: invert good dense -->

## ✅ Shutdown: how does a sleeping consumer ever stop?

```text
close():
    lock(m);  closed = true;  unlock(m)
    notifyAll(notEmpty);  notifyAll(notFull)   // every sleeper must wake up and see "closed"

pop():
    lock(m)
    while items is empty and not closed:
        wait(notEmpty, m)
    if items is empty:                         // closed, and nothing left
        unlock(m)
        return nothing
    x = take the first item
    unlock(m)
    notify(notFull)
    return x
```

* One more piece of state: `closed`. Consumers drain what's left, then get *nothing* and stop · `push` after `close()` returns `false`
* What each call does after `close()` is part of the contract — Day 2's error semantics

---

## Producers, consumers, one bounded queue — try it

<div class="viz" data-viz="queueSim"></div>

---

## Three answers to "the queue is empty"

| | Call | When the queue is empty |
|---|---|---|
| **Balk** | `tryPop()` — Day 3 | returns "nothing" at once |
| **Block** | `pop()` | waits as long as it takes |
| **Time out** | `popFor(100ms)` | waits up to a limit, then gives up |


---

<!-- _class: lead invert -->

# Semaphores

### At most N

---

## At most N

* A mutex means **at most one**. Real limits are often bigger:
* at most **4** downloads at once — or the server blocks you
* at most **10** open database connections
* at most **500** shoppers inside a flash sale — team project #6
* > A **semaphore** is a counter of permits. `acquire()` takes one — and waits while there are none. `release()` gives one back.

---

## A parking lot with N spaces — try it

<div class="viz" data-viz="semaphore"></div>

---

## In code

```cpp
std::counting_semaphore<4> slots(4);      // 4 permits

void download(const Url& u) {
    slots.acquire();                       // waits while 4 downloads run
    fetch(u);                              // throws? release() never runs
    slots.release();
}
```

* An exception between `acquire` and `release` **leaks** a permit. After four failed downloads, nothing downloads ever again — no error, no crash
* Release in a destructor (RAII), like `lock_guard` does — or `try`/`finally` and `with sem:` elsewhere

---

## Semaphore vs. mutex

* A mutex has an **owner**: only the thread that locked it may unlock it
* A semaphore has **no owner**: any thread may `release()` — even one that never acquired
* So a semaphore is also a **signal**: start at 0; one thread waits, another wakes it

```cpp
std::binary_semaphore ready(0);
// thread A:   ready.acquire();               — waits for B
// thread B:   prepare(); ready.release();    — "go"
```

<p class="note">The flip side: as a lock, nothing stops the wrong thread from releasing it — and one extra <code>release()</code> silently raises the limit. Python's <code>BoundedSemaphore</code> raises an error instead.</p>

---

<!-- _class: invert dense -->

## The bounded queue again — with two semaphores

```cpp
std::counting_semaphore<N> freeSlots(N);   // permits = empty slots
std::counting_semaphore<N> usedSlots(0);   // permits = items
std::mutex m;                              // guards items

void push(T v) {
    freeSlots.acquire();                   // wait for an empty slot
    { std::lock_guard lk(m); items.push_back(std::move(v)); }
    usedSlots.release();                   // one more item
}

T pop() {
    usedSlots.acquire();                   // wait for an item
    std::unique_lock lk(m);
    T v = std::move(items.front());
    items.pop_front();
    lk.unlock();
    freeSlots.release();                   // one more empty slot
    return v;
}
```


---

<!-- _class: lead invert -->

# Readers–Writer Locks

### Many readers, or one writer

---

## Readers don't race with readers

* Day 3: *"Is a race possible if both threads only read?"* — **No.**
* Yet with one mutex, every cache lookup waits for every other lookup
* An image cache: 99% lookups, 1% inserts — almost all of that waiting is for nothing
* > A **readers–writer lock**: any number of readers together — **or** one writer, alone.

---

## In code

```cpp
std::shared_mutex m;                                   // guards cache

std::optional<Image> find(const std::string& path) {
    std::shared_lock lk(m);                            // shared: many readers at once
    auto it = cache.find(path);
    if (it == cache.end()) return std::nullopt;
    return it->second;
}

void insert(const std::string& path, Image img) {
    std::unique_lock lk(m);                            // exclusive: one writer, no readers
    cache.insert_or_assign(path, std::move(img));
}
```

* Same idea elsewhere: Java `ReentrantReadWriteLock`, Go `sync.RWMutex`, C# `ReaderWriterLockSlim`, Rust `RwLock` — Python's standard library has none

---

## Many readers, one writer — try it

<div class="viz" data-viz="rwlock"></div>

---

<!-- _class: lead invert -->

# Once, Latch, Barrier

### Three more ways to wait

---

## Once — run it exactly one time

* > **Once**: a piece of code that runs **exactly one time**, no matter how many threads call it at the same moment.
* The first thread runs it; every other thread **waits** until it's done, then goes on
* For one-time setup: a singleton, a config file, a connection pool

---

<!-- _class: invert good -->

<span class="tag">Day 1 → Day 3 → today</span>

## ✅ The Singleton, finally

```cpp
static Singleton& get() {
    static Singleton instance;      // C++11: initialized exactly once;
    return instance;                // other threads wait until it's ready
}

std::once_flag initOnce;            // any other one-time setup
std::call_once(initOnce, [] { loadConfig(); });
```

* The language does the check-then-act for you — correctly
* Same idea elsewhere: Go `sync.Once`, Java holder class, C# `Lazy<T>`, Rust `OnceLock`, Python module import (runs once)
* Double-checked locking — the hand-made version — waits for Day 5

---

## Latch — wait until N things have happened

* > **Latch**: a counter that starts at N. Threads **count down**; a thread that **waits** sleeps until the count reaches 0.
* At 0 it opens and stays open — a latch is used **once**
* For "start when everything is ready": load 3 shards, then start serving

---

## Latch — example

```cpp
std::latch loaded(3);                          // counts down from 3 — once

for (int i = 0; i < 3; ++i)
    workers.emplace_back([&loaded, i] { loadShard(i); loaded.count_down(); });

loaded.wait();                                 // sleeps until the count hits 0
startServing();
```

<div class="viz" data-viz="latch"></div>

<div class="steps">

* `main` waits; three shards are loading — count 3
* Shard 2 is done — 2
* Shard 0 — 1
* Shard 1 — 0: the latch opens, `main` goes on. Elsewhere: Java `CountDownLatch`, Go `sync.WaitGroup`, C# `CountdownEvent`

</div>

---

## Barrier — wait for everyone

* > **Barrier**: N threads work in steps. A thread that finishes a step **waits** at the barrier until all N have arrived — then they all start the next step together.
* Unlike a latch, it opens and then closes again: it is used once per step
* For work that goes step by step: simulations, physics, games — the Falling-Sand project

---

## Barrier — example

```cpp
std::barrier sync(4, []() noexcept { swapBuffers(); });   // runs once per step

void worker(int id) {
    for (int step = 0; step < steps; ++step) {
        updateRegion(id);              // reads the neighbours' cells from the last step
        sync.arrive_and_wait();        // wait for the other three
    }
}
```

* Same idea elsewhere: Java `CyclicBarrier`, Python `threading.Barrier`, C# `Barrier`, Rust `Barrier`

---

## Barrier — try it

<div class="viz" data-viz="barrier"></div>

---

## The toolbox

| You need… | Use |
|---|---|
| no lost updates, no half-done states | **mutex** |
| to wait until something is true | **condition variable** + mutex |
| to hand work from thread to thread | **blocking queue** |
| at most N at once | **semaphore** |
| many readers, rare writers | **readers–writer lock** |
| to do it exactly once | **`call_once`** / static local |
| to wait for N events | **latch** |
| all threads in step | **barrier** |

---

<!-- _class: invert rosetta -->

## The same toolbox, in every language

| | C++ | Java | Go | Python | C# | Rust |
|---|---|---|---|---|---|---|
| mutex | `mutex` | `synchronized` | `sync.Mutex` | `Lock` | `lock` | `Mutex` |
| wait for a condition | `condition_variable` | `wait` / `Condition` | `sync.Cond` | `Condition` | `Monitor.Wait` | `Condvar` |
| blocking queue | build it | `ArrayBlockingQueue` | `chan T` | `queue.Queue` | `Channel<T>` | `sync_channel` |
| semaphore | `counting_semaphore` | `Semaphore` | buffered `chan` | `Semaphore` | `SemaphoreSlim` | tokio `Semaphore` |
| readers–writer | `shared_mutex` | `ReentrantReadWriteLock` | `sync.RWMutex` | — | `ReaderWriterLockSlim` | `RwLock` |
| once | `call_once` | holder class | `sync.Once` | module import | `Lazy<T>` | `OnceLock` |
| latch | `latch` | `CountDownLatch` | `sync.WaitGroup` | — | `CountdownEvent` | — |
| barrier | `barrier` | `CyclicBarrier` | — | `Barrier` | `Barrier` | `Barrier` |

---

<!-- _class: lead invert -->

# Async Patterns

### From threads to tasks

---

## Think in tasks, not threads

* A **task** — a piece of work: "resize this image", "answer this request"
* A **thread** — a worker that runs tasks
* So far: one thread per task. That doesn't scale:
* each thread costs a system call to create, and its own stack — often 0.5–8 MB reserved
* 10,000 requests → 10,000 threads → the OS spends its time switching between them
* > Keep a **few** threads busy with **many** tasks.

---

## Thread pool — try it

<div class="viz" data-viz="pool"></div>

---

<!-- _class: invert good dense -->

## ✅ A thread pool in 20 lines

```cpp
class ThreadPool {
    BlockingQueue<std::function<void()>> tasks{1000};   // the queue from earlier
    std::vector<std::thread> workers;
public:
    explicit ThreadPool(int n) {
        for (int i = 0; i < n; ++i)
            workers.emplace_back([this] {
                while (auto task = tasks.pop())          // nullopt after close(): exit
                    (*task)();
            });
    }
    void submit(std::function<void()> f) { tasks.push(std::move(f)); }
    ~ThreadPool() {
        tasks.close();                                   // workers finish the queue, then stop
        for (auto& w : workers) w.join();
    }
};
```

* The queue's `close()` is the whole shutdown story

---

## How many threads?

* **CPU-bound** tasks — about one thread per core; more only adds switching
* **I/O-bound** tasks — more threads than cores, because most are asleep… or stop blocking altogether: the event loop, soon
* **Queue full?** Make the submitter wait (backpressure), reject the task, or run it on the caller's thread — a policy you choose
* Same idea elsewhere: Java `ExecutorService`, Python `ThreadPoolExecutor`, C# `Task.Run`, Go — the runtime pools goroutines for you

---

## Getting a result back: future and promise

* `submit()` returns nothing — so how does the caller get the answer?
* > A **promise** is the write end and a **future** the read end of a one-time channel for **one** value — or one exception.

```cpp
std::promise<int> p;
std::future<int> f = p.get_future();

std::thread t([&p] { p.set_value(heavyCompute()); });   // writes once
int x = f.get();                                         // waits until it's written
t.join();
```

---

## Errors travel too

```cpp
std::thread t([&p] {
    try { p.set_value(parseConfig("app.conf")); }
    catch (...) { p.set_exception(std::current_exception()); }
});

try { Config c = f.get(); }                       // parseConfig's exception lands here
catch (const ParseError& e) { useDefaults(); }
```

* The exception crosses threads and comes out where the result was expected — Day 2's error semantics, across threads
* `get()` works once; for several waiters there is `std::shared_future`
* Same idea elsewhere: Java `CompletableFuture`, JavaScript `Promise`, Python `concurrent.futures.Future`, C# `Task<T>`, Rust `Future`, Go — a channel of size 1

---

<!-- _class: invert dense -->

## A pool that returns futures

```cpp
template <typename F>
auto submit(F f) -> std::future<decltype(f())> {
    using R = decltype(f());
    auto task = std::make_shared<std::packaged_task<R()>>(std::move(f));
    std::future<R> result = task->get_future();
    tasks.push([task] { (*task)(); });       // a worker runs it; the future gets the result
    return result;
}

auto words = pool.submit([] { return countWords("book.txt"); });
doOtherWork();
std::cout << words.get();                     // waits only if it isn't done yet
```

* `packaged_task` = a function + a promise: running it fills the future — with the value or the exception

---

<!-- _class: invert bad dense -->

## Deadlock with no locks at all

```cpp
ThreadPool pool(2);                              // 2 workers

auto task = [&pool] {
    auto child = pool.submit([] { return 1; });  // the child goes into the queue
    return child.get() + 1;                      // the worker sleeps until the child is done
};
auto a = pool.submit(task);                      // parent A → worker 1
auto b = pool.submit(task);                      // parent B → worker 2: no worker is left
```

<div class="viz" data-viz="poolDeadlock"></div>

<div class="steps">

* A pool with 2 workers. Two parent tasks, A and B, start running
* Each submits a child task — into the queue
* Each waits for its child: both workers are asleep
* No worker is free to run `a` or `b` — **a cycle of waiting**, no lock needed

</div>

---

## How to stay out of it

* Don't wait on a future **inside** a task of the same pool — while it waits, that worker can't run anything
* If you must wait: run that work in a **different** pool, or chain the next step instead of waiting — Java `thenApply`, C# `ContinueWith`, JavaScript `then`
* Or use a pool whose waiting workers **run other queued tasks** in the meantime — work-stealing pools do exactly that

---

## One queue for everyone — and how pools fix it

<div class="viz" data-viz="poolQueues"></div>

* One shared queue is one lock: every submit and every take fights for it — the more workers, the more waiting
* Shard it: each worker gets **its own queue** and pushes and pops there, without a fight
* A worker whose queue runs dry **steals** from another worker's queue — **work stealing**. Java `ForkJoinPool`, Go's scheduler, Rust's rayon and tokio, the .NET thread pool.

---

<!-- _class: lead invert -->

# Wait, Don't Spin

### State · a lock · a waiting room — every tool today is built from these

<script>
/* viz:start */
/* Day 4 interactive diagrams.

   Every <div class="viz" data-viz="NAME"> on a slide gets the diagram
   registered below as def('NAME', ...). Two kinds:

   - step diagrams follow the slide's fragments: each → / Space reveals
     the next caption and moves the diagram one step (← goes back);
   - live diagrams run while their slide is on screen and have their own
     buttons and sliders.

   Everything is drawn as inline SVG with the deck's blueprint palette. */
(function () {
  'use strict';
  if (window.__day4viz) return;
  window.__day4viz = {};

  var NS = 'http://www.w3.org/2000/svg';
  var C = {
    bg: '#0b2136', panel: '#14385a', line: '#6fa8cc', accent: '#5cc8ea',
    good: '#7ed0a3', goodFill: '#173f37', bad: '#f0a06e', badFill: '#43291f',
    warn: '#e8cf74', warnFill: '#3b3520', violet: '#b9a3f0', violetFill: '#2b2950',
    fg: '#eaf4fb', dim: '#a9c8de', faint: '#244a6c', accentFill: '#123d58'
  };
  var MONO = "'IBM Plex Mono', Menlo, monospace";
  var SANS = "'IBM Plex Sans', Helvetica, Arial, sans-serif";
  var FILL = {};
  FILL[C.accent] = C.accentFill; FILL[C.good] = C.goodFill; FILL[C.bad] = C.badFill;
  FILL[C.warn] = C.warnFill; FILL[C.violet] = C.violetFill; FILL[C.line] = C.panel; FILL[C.dim] = C.panel;
  var PALETTE = [C.accent, C.good, C.violet, C.warn, C.bad, '#7fb8ff', '#9be3d0', '#f3b5d8'];
  var uid = 0;

  /* ---------- drawing helpers ---------- */

  function S(tag, attrs, parent) {
    var e = document.createElementNS(NS, tag);
    if (attrs) for (var k in attrs) if (attrs[k] != null) e.setAttribute(k, attrs[k]);
    if (parent) parent.appendChild(e);
    return e;
  }
  function T(parent, x, y, str, o) {
    o = o || {};
    var size = o.size || 15;
    if (size <= 13) size += 1.5;           // small labels must still read on a projector
    var e = S('text', {
      x: x, y: y, 'text-anchor': o.anchor || 'middle', fill: o.fill || C.fg,
      'font-family': o.mono ? MONO : SANS, 'font-size': size,
      'font-weight': o.weight || 400, 'letter-spacing': o.ls,
      'font-style': o.italic ? 'italic' : null, opacity: o.opacity
    }, parent);
    e.textContent = str;
    return e;
  }
  function hdr(parent, x, y, str, anchor, col) {
    return T(parent, x, y, str, { size: 12, weight: 700, ls: 2, fill: col || C.accent, anchor: anchor || 'middle' });
  }
  function R(parent, x, y, w, h, o) {
    o = o || {};
    return S('rect', {
      x: x, y: y, width: w, height: h, rx: o.rx == null ? 3 : o.rx,
      fill: o.fill || 'none', stroke: o.stroke || C.line, 'stroke-width': o.sw || 1.4,
      'stroke-dasharray': o.dash, opacity: o.opacity
    }, parent);
  }
  function L(parent, x1, y1, x2, y2, o) {
    o = o || {};
    return S('line', {
      x1: x1, y1: y1, x2: x2, y2: y2, stroke: o.stroke || C.line, 'stroke-width': o.sw || 1.6,
      'stroke-dasharray': o.dash, opacity: o.opacity, 'marker-end': o.end
    }, parent);
  }
  function Pth(parent, d, o) {
    o = o || {};
    return S('path', {
      d: d, fill: o.fill || 'none', stroke: o.stroke || C.line, 'stroke-width': o.sw || 1.6,
      'stroke-dasharray': o.dash, opacity: o.opacity, 'marker-end': o.end
    }, parent);
  }
  function G(parent, cls) { return S('g', cls ? { 'class': cls } : null, parent); }
  function clear(g) { while (g.firstChild) g.removeChild(g.firstChild); }
  function flash(e) {
    e.classList.remove('flash');
    void e.getBoundingClientRect();
    e.classList.add('flash');
  }
  function fade(e, v) { e.style.opacity = v ? 1 : 0; return e; }
  function rnd(a, b) { return a + Math.random() * (b - a); }
  function bar(parent, x, y, w, h, col) {   // a progress bar; set(p) with p in 0..1
    R(parent, x, y, w, h, { fill: C.bg, stroke: C.faint, sw: 1, rx: 2 });
    var f = R(parent, x, y, w, h, { fill: col || C.good, stroke: 'none', rx: 2 });
    f.classList.add('bar');
    f.style.transform = 'scaleX(0)';
    return { el: f, set: function (p, instant) {
      if (instant) f.style.transition = 'none';
      f.style.transform = 'scaleX(' + Math.max(0, Math.min(1, p)) + ')';
      if (instant) { void f.getBoundingClientRect(); f.style.transition = ''; }
    } };
  }

  function mkSvg(host, w, h) {
    var id = ++uid;
    var svg = S('svg', { viewBox: '0 0 ' + w + ' ' + h, width: w, height: h, 'class': 'noanim' });
    svg.style.width = '100%';
    svg.style.height = 'auto';
    host.appendChild(svg);
    var defs = S('defs', null, svg);
    svg.arrow = {};
    var cols = { accent: C.accent, good: C.good, bad: C.bad, line: C.line, dim: C.dim, warn: C.warn, violet: C.violet };
    Object.keys(cols).forEach(function (n) {
      var m = S('marker', { id: 'ah' + id + n, viewBox: '0 0 10 10', refX: 9, refY: 5,
        markerWidth: 7, markerHeight: 7, orient: 'auto-start-reverse' }, defs);
      S('path', { d: 'M0,0 L10,5 L0,10 z', fill: cols[n] }, m);
      svg.arrow[n] = 'url(#ah' + id + n + ')';
    });
    setTimeout(function () { svg.classList.remove('noanim'); }, 80);
    return svg;
  }

  /* A labelled box that moves (CSS transition on transform) and shows
     run / wait / sleep / bad / good / off states. */
  function Token(parent, label, color, o) {
    o = o || {};
    var w = o.w || 70, h = o.h || 32;
    var g = G(parent, o.still ? null : 'mv');
    var inner = G(g, 'fd');
    var box = R(inner, -w / 2, -h / 2, w, h, { fill: FILL[color] || C.panel, stroke: color, sw: 1.8, rx: 5 });
    box.classList.add('col');
    var lab = T(inner, 0, 5, label, { size: o.size || 14, weight: 700, mono: !!o.mono });
    var sub = T(inner, 0, h / 2 + 17, '', { size: o.subSize || 13, fill: C.dim });
    var zz = T(inner, w / 2 + 6, -h / 2 + 2, 'z z', { size: 12, fill: C.dim, italic: true, anchor: 'start' });
    zz.classList.add('fd');
    zz.style.opacity = 0;
    var tok = {
      g: g, x: 0, y: 0, color: color, w: w, h: h,
      at: function (x, y) {
        tok.x = x; tok.y = y;
        g.style.transform = 'translate(' + x + 'px,' + y + 'px)';
        return tok;
      },
      jump: function (x, y) {      // move without animation
        g.style.transition = 'none';
        tok.at(x, y);
        void g.getBoundingClientRect();
        g.style.transition = '';
        return tok;
      },
      show: function (v) { g.style.opacity = v ? 1 : 0; return tok; },
      note: function (str, col) { sub.textContent = str || ''; sub.setAttribute('fill', col || C.dim); return tok; },
      label: function (str) { lab.textContent = str; return tok; },
      state: function (st) {
        var waiting = st === 'wait' || st === 'sleep';
        box.setAttribute('stroke-dasharray', waiting ? '5 4' : 'none');
        box.setAttribute('stroke', st === 'bad' ? C.bad : waiting ? C.dim : st === 'good' ? C.good : tok.color);
        box.setAttribute('fill', st === 'bad' ? C.badFill : waiting ? C.bg : st === 'good' ? C.goodFill : (FILL[tok.color] || C.panel));
        lab.setAttribute('fill', waiting ? C.dim : C.fg);
        inner.style.opacity = st === 'off' ? 0.22 : st === 'sleep' ? 0.8 : 1;
        zz.style.opacity = st === 'sleep' ? 1 : 0;
        return tok;
      },
      remove: function () { if (g.parentNode) g.parentNode.removeChild(g); }
    };
    return tok;
  }

  /* Key/value panel for "the shared state" next to a diagram. */
  function Panel(parent, x, y, w, h, title, keys, o) {
    o = o || {};
    R(parent, x, y, w, h, { fill: 'rgba(20,56,90,.55)', stroke: C.line, sw: 1.2 });
    hdr(parent, x + 14, y + 22, title, 'start');
    var rows = {};
    keys.forEach(function (k, i) {
      var yy = y + 52 + i * (o.rowH || 30);
      T(parent, x + 14, yy, k, { anchor: 'start', size: 13, fill: C.dim });
      rows[k] = T(parent, x + w - 14, yy, '', { anchor: 'end', size: o.size || 15, mono: true, weight: 600 });
    });
    return {
      set: function (k, v, col) {
        var e = rows[k];
        v = v == null ? '—' : String(v);
        if (e.textContent !== v) { e.textContent = v; flash(e); }
        e.setAttribute('fill', col || C.fg);
      }
    };
  }

  /* Timers for step animations; cleared on every re-render. */
  function Sched() {
    var ids = [];
    return {
      clear: function () { ids.forEach(clearTimeout); ids = []; },
      after: function (ms, fn) { ids.push(setTimeout(fn, ms)); },
      // run fns one after another (numbers are pauses in ms), or all at once
      seq: function (animate, list) {
        var t = 0, self = this;
        list.forEach(function (it) {
          if (typeof it === 'number') { t += it; return; }
          if (!animate) it(); else if (t === 0) it(); else self.after(t, it);
        });
      }
    };
  }

  /* ---------- controls (HTML, under the diagram) ---------- */

  function Controls(host) {
    var d = document.createElement('div');
    d.className = 'viz-ctl';
    host.appendChild(d);
    // arrow keys on a focused slider move the slider, not the slides
    d.addEventListener('keydown', function (e) {
      if (e.target && e.target.tagName === 'INPUT') e.stopPropagation();
    });
    return d;
  }
  function Btn(parent, label, fn) {
    var b = document.createElement('button');
    b.type = 'button';
    b.textContent = label;
    b.addEventListener('click', function () { fn(b); b.blur(); });
    parent.appendChild(b);
    return b;
  }
  function Toggle(parent, opts, value, fn) {
    var wrap = document.createElement('span');
    wrap.className = 'grp';
    parent.appendChild(wrap);
    var btns = opts.map(function (o) {
      var b = Btn(wrap, o[1], function () { set(o[0]); fn(o[0]); });
      b.setAttribute('data-v', String(o[0]));
      return b;
    });
    function set(v) { btns.forEach(function (b) { b.classList.toggle('on', b.getAttribute('data-v') === String(v)); }); }
    set(value);
    return { set: set };
  }
  function Slider(parent, label, min, max, step, value, fn, fmt) {
    fmt = fmt || String;
    var lab = document.createElement('label');
    var name = document.createElement('span');
    name.textContent = label;
    var inp = document.createElement('input');
    inp.type = 'range'; inp.min = min; inp.max = max; inp.step = step; inp.value = value;
    var out = document.createElement('span');
    out.className = 'val';
    out.textContent = fmt(value);
    inp.addEventListener('input', function () { var v = Number(inp.value); out.textContent = fmt(v); fn(v); });
    inp.addEventListener('change', function () { inp.blur(); });
    lab.appendChild(name); lab.appendChild(inp); lab.appendChild(out);
    parent.appendChild(lab);
    return { el: inp, set: function (v) { inp.value = v; out.textContent = fmt(v); } };
  }

  /* ---------- timeline: lanes of steps + a state panel ---------- */

  function Timeline(host, cfg) {
    var W = 1140, rh = cfg.rowH || 38, top = 30, n = cfg.rows.length;
    var pw = cfg.panelW || 320, gap = 26;
    var ph = cfg.panelH || n * rh;
    var H = top + Math.max(n * rh, ph) + 6;
    var svg = mkSvg(host, W, H);
    var x0 = 40, x1 = W - pw - gap;
    var lw = (x1 - x0) / cfg.lanes.length;
    cfg.lanes.forEach(function (lane, i) {
      var cx = x0 + lw * (i + 0.5);
      hdr(svg, cx, 18, lane.name, 'middle', lane.color);
      L(svg, cx, top - 2, cx, top + n * rh, { sw: 1, dash: '3 4', opacity: 0.45 });
    });
    var rows = cfg.rows.map(function (row, i) {
      var g = G(svg, 'fd');
      var y = top + i * rh;
      T(g, 16, y + rh / 2 + 5, String(i + 1), { size: 13, fill: C.dim });
      var span = row.lane < 0;
      var bx = span ? x0 + 8 : x0 + lw * row.lane + 8;
      var bw = span ? (x1 - x0) - 16 : lw - 16;
      var col = row.cls === 'bad' ? C.bad : row.cls === 'good' ? C.good : C.line;
      var box = R(g, bx, y + 3, bw, rh - 6, {
        fill: row.cls === 'bad' ? C.badFill : row.cls === 'good' ? C.goodFill : C.panel, stroke: col, sw: 1.4
      });
      box.classList.add('col');
      T(g, bx + bw / 2, y + rh / 2 + 5, row.text, { size: cfg.fontSize || 14, mono: true });
      g.style.opacity = 0;
      return { g: g, box: box, col: col };
    });
    var panel = Panel(svg, W - pw, top, pw, ph, cfg.stateTitle || 'STATE', cfg.keys);
    return {
      steps: n,
      render: function (k) {
        rows.forEach(function (row, i) {
          row.g.style.opacity = i < k ? 1 : 0;
          row.box.setAttribute('stroke', i === k - 1 && row.col === C.line ? C.accent : row.col);
          row.box.setAttribute('stroke-width', i === k - 1 ? 2.2 : 1.4);
        });
        var st = cfg.state(k);
        cfg.keys.forEach(function (key) {
          var v = st[key];
          if (Array.isArray(v)) panel.set(key, v[0], v[1]); else panel.set(key, v);
        });
      }
    };
  }

  /* ---------- engine ---------- */

  var registry = {}, insts = [];
  function def(name, fn) { registry[name] = fn; }

  // For debugging from the console: show step k of a step diagram, or run
  // a live diagram n frames ahead — __day4viz.show('monitor', 3),
  // __day4viz.run('pool', 1 / 60, 300).
  window.__day4viz.show = function (name, k) {
    insts.forEach(function (inst) { if (inst.name === name && inst.api.render) inst.api.render(k, -1); });
  };
  window.__day4viz.run = function (name, dt, n) {
    insts.forEach(function (inst) {
      if (inst.name !== name || !inst.api.tick) return;
      for (var i = 0; i < (n || 1); i++) inst.api.tick(dt);
    });
  };

  function stepOf(section) {
    var frags = section.querySelectorAll('[data-marpit-fragment]');
    var n = 0, bespoke = false;
    for (var i = 0; i < frags.length; i++) {
      var a = frags[i].getAttribute('data-bespoke-marp-fragment');
      if (a) bespoke = true;
      if (a === 'active') n++;
    }
    return bespoke ? n : frags.length;   // no presenter (e.g. print): last step
  }

  function mount() {
    var hosts = document.querySelectorAll('.viz[data-viz]');
    for (var i = 0; i < hosts.length; i++) {
      var host = hosts[i];
      if (host.__viz) continue;
      var name = host.getAttribute('data-viz'), fn = registry[name];
      var inst = {
        host: host, name: name, section: host.closest('section'),
        slide: host.closest('svg[data-marpit-svg]') || host.closest('section'), step: -1, api: {}
      };
      host.__viz = inst;
      if (!fn) { host.textContent = '[missing diagram: ' + name + ']'; continue; }
      try { inst.api = fn(host, inst) || {}; } catch (err) {
        console.error('[viz] ' + name, err);
        host.textContent = '[diagram error: ' + name + ']';
        continue;
      }
      if (inst.api.steps != null) {
        var frags = inst.section.querySelectorAll('[data-marpit-fragment]').length;
        if (frags !== inst.api.steps) console.warn('[viz] ' + name + ': ' + inst.api.steps + ' steps but ' + frags + ' fragments');
        inst.badge = document.createElement('div');
        inst.badge.className = 'viz-step';
        inst.section.appendChild(inst.badge);      // top-right corner of the slide
      }
      insts.push(inst);
    }
  }

  function sync() {
    for (var i = 0; i < insts.length; i++) {
      var inst = insts[i];
      if (!inst.api.render) continue;
      var n = inst.api.steps, k = stepOf(inst.section);
      if (n != null && k > n) k = n;
      if (k === inst.step) continue;
      var prev = inst.step;
      inst.step = k;
      try { inst.api.render(k, prev); } catch (err) { console.error('[viz] ' + inst.name, err); }
      if (inst.badge) inst.badge.textContent = k === 0 ? 'press → to step' : 'step ' + k + ' / ' + n;
    }
  }

  var lastT = 0;
  function frame(now) {
    var dt = lastT ? Math.min(0.05, (now - lastT) / 1000) : 0;
    lastT = now;
    var presenting = !!document.querySelector('.bespoke-marp-active');
    for (var i = 0; i < insts.length; i++) {
      var inst = insts[i];
      if (!inst.api.tick) continue;
      if (presenting && !inst.slide.classList.contains('bespoke-marp-active')) continue;
      try { inst.api.tick(dt); } catch (err) { console.error('[viz] ' + inst.name, err); inst.api.tick = null; }
    }
    requestAnimationFrame(frame);
  }

  function init() {
    mount();
    sync();
    var pending = false;
    new MutationObserver(function () {
      if (pending) return;
      pending = true;
      requestAnimationFrame(function () { pending = false; sync(); });
    }).observe(document.body, { subtree: true, attributes: true, attributeFilter: ['data-bespoke-marp-fragment'] });
    var st = document.createElement('style');
    st.textContent = 'body[data-bespoke-view=overview] .steps li:not(:last-child){visibility:hidden}';
    document.head.appendChild(st);
    requestAnimationFrame(frame);
  }

  /* ================= before the break: primitives ================= */

  /* Exclusion vs. coordination. */
  def('kinds', function (host) {
    var W = 1140, H = 230, svg = mkSvg(host, W, H), q = Sched();
    hdr(svg, 270, 18, 'EXCLUSION · NOT AT THE SAME TIME');
    R(svg, 200, 50, 180, 150, { fill: 'rgba(20,56,90,.45)' });
    L(svg, 200, 102, 200, 148, { stroke: C.accent, sw: 5 });
    T(svg, 192, 94, 'door', { size: 11, fill: C.dim, anchor: 'end' });
    T(svg, 290, 188, 'critical section', { size: 13, fill: C.dim });
    L(svg, 570, 30, 570, 210, { stroke: C.faint, sw: 1, dash: '4 6' });
    hdr(svg, 865, 18, "COORDINATION · NOT BEFORE IT'S READY");
    R(svg, 830, 102, 70, 46, { fill: C.panel });
    T(svg, 865, 172, 'one slot', { size: 13, fill: C.dim });
    var ping = T(svg, 925, 96, 'notify →', { size: 13, fill: C.accent, weight: 700, anchor: 'start' });
    ping.classList.add('fd');
    var A = Token(svg, 'A', C.accent, { w: 56 }), B = Token(svg, 'B', C.good, { w: 56 });
    var P = Token(svg, 'producer', C.good, { w: 110 }), K = Token(svg, 'consumer', C.accent, { w: 110 });
    var item = Token(svg, 'item', C.warn, { w: 50, h: 26, size: 12 });
    return {
      steps: 4,
      render: function (k, prev) {
        q.clear();
        var fwd = prev === k - 1;
        A.at(80, 100).state('run').note('');
        B.at(80, 160).state('run').note('');
        P.at(690, 125).state('run').note('');
        K.at(1060, 125).state('run').note('');
        item.at(770, 125).show(0);
        fade(ping, 0);
        if (k >= 1) { A.at(290, 125).note('inside'); B.at(140, 125).state('wait').note('waits'); }
        if (k >= 2) { A.at(470, 125).state('off').note(''); B.at(290, 125).state('run').note('inside'); }
        if (k >= 3) K.at(990, 125).state('sleep').note('nothing yet: sleeps');
        if (k >= 4) {
          q.seq(fwd, [
            function () { item.show(1); P.note('puts an item'); },
            60, function () { item.at(865, 125); },
            650, function () { fade(ping, 1); K.state('run').note('woken'); },
            450, function () { item.at(990, 125); },
            550, function () { item.show(0); K.note('takes it'); }
          ]);
        }
      }
    };
  });

  /* Spin vs. sleep-poll vs. condition variable: a small simulation in
     "worker time" (50 simulated ms per real second). */
  def('waitModes', function (host) {
    var W = 1140, H = 226, svg = mkSvg(host, W, H), ctl = Controls(host);
    var MEAS = {
      spin: '5.0 s of CPU per 5 s idle · a job is picked up in 3 µs',
      poll: '0.007 s of CPU per 5 s idle · a job waits 6 ms on average',
      cv: '0.000 s of CPU per 5 s idle · a job is picked up in 17 µs'
    };
    var WIN = 250, SPEED = 50, CHECK = 0.05, POLL = 10, WORK = 15;
    var SX = 260, SW = 850, SY = 46, SH = 44, SLEEP = '#28507a';
    hdr(svg, 120, 18, 'JOB QUEUE');
    R(svg, 30, SY, 180, SH, { fill: C.bg });
    var qG = G(svg);
    var qText = T(svg, 120, SY + SH + 22, '', { size: 13, fill: C.dim });
    hdr(svg, SX, 18, 'THE WORKER · LAST 250 ms OF ITS LIFE', 'start');
    R(svg, SX, SY, SW, SH, { fill: C.bg });
    var segG = G(svg), arrG = G(svg);
    var ly = SY + SH + 22;
    [[C.warn, 'checking “empty?”'], [SLEEP, 'asleep'], [C.good, 'running a job']].forEach(function (it, i) {
      var x = SX + i * 190;
      R(svg, x, ly - 11, 14, 14, { fill: it[0], stroke: 'none', rx: 2 });
      T(svg, x + 22, ly, it[1], { anchor: 'start', size: 13, fill: C.dim });
    });
    T(svg, SX + 3 * 190, ly, '▼ a job arrives', { anchor: 'start', size: 13, fill: C.accent });
    T(svg, 30, 160, 'CPU burned while idle', { anchor: 'start', size: 14, fill: C.dim });
    var cpu = bar(svg, 200, 148, 300, 16, C.warn);
    var cpuText = T(svg, 512, 161, '', { anchor: 'start', size: 16, mono: true, weight: 600 });
    T(svg, 640, 160, 'a job waits before it starts', { anchor: 'start', size: 14, fill: C.dim });
    var latText = T(svg, 1110, 161, '', { anchor: 'end', size: 16, mono: true, weight: 600 });
    var measText = T(svg, 30, 204, '', { anchor: 'start', size: 14, fill: C.accent });

    var mode = 'spin', now = 0, nextArr = 20, queue = [], segs = [], arrivals = [], lat = [];
    var w = { st: 'idle', until: 0, nextCheck: 0 };
    function seg(kind, a, b) {
      if (b <= a) return;
      var last = segs[segs.length - 1];
      if (last && last.k === kind && Math.abs(a - last.b) < 1e-6) last.b = b;
      else segs.push({ k: kind, a: a, b: b });
    }
    function arrive(t) { queue.push(t); arrivals.push(t); nextArr = t + rnd(35, 110); }
    function startJob(t) {
      lat.push(t - queue.shift());
      if (lat.length > 12) lat.shift();
      w.st = 'work';
      w.until = t + WORK;
    }
    function idleAgain(t) { w.st = 'idle'; if (queue.length) startJob(t); else w.nextCheck = t + POLL; }
    function advance(T1) {
      var guard = 0;
      while (now < T1 - 1e-9 && guard++ < 20000) {
        if (w.st === 'work') {
          var end = Math.min(w.until, T1);
          while (nextArr <= end) arrive(nextArr);
          seg('run', now, end);
          now = end;
          if (now >= w.until - 1e-9) idleAgain(now);
        } else if (mode === 'spin') {
          var ta = Math.min(nextArr, T1);
          seg('check', now, ta);
          now = ta;
          if (now >= nextArr - 1e-9) { arrive(now); startJob(now); }
        } else if (mode === 'poll') {
          var tn = Math.min(w.nextCheck, nextArr, T1);
          seg('sleep', now, tn);
          now = tn;
          if (nextArr <= now + 1e-9 && nextArr < w.nextCheck) { arrive(nextArr); continue; }
          if (now >= w.nextCheck - 1e-9) {
            seg('check', now, now + CHECK);
            now += CHECK;
            if (queue.length) startJob(now); else w.nextCheck = now + POLL;
          }
        } else {
          var tc = Math.min(nextArr, T1);
          seg('sleep', now, tc);
          now = tc;
          if (now >= nextArr - 1e-9) { arrive(now); seg('check', now, now + 0.02); now += 0.02; startJob(now); }
        }
      }
    }
    function draw() {
      clear(segG); clear(arrG); clear(qG);
      var t0 = now - WIN, check = 0, sleep = 0;
      segs = segs.filter(function (s) { return s.b > t0 - 5; });
      segs.forEach(function (s) {
        var a = Math.max(s.a, t0), b = s.b;
        if (b <= t0) return;
        if (s.k === 'check') check += b - a; else if (s.k === 'sleep') sleep += b - a;
        var x = SX + (a - t0) / WIN * SW;
        var wd = Math.max(s.k === 'check' ? 3 : 0.5, (b - a) / WIN * SW);
        R(segG, x, SY + 4, Math.min(wd, SX + SW - x), SH - 8, {
          fill: s.k === 'check' ? C.warn : s.k === 'run' ? C.good : SLEEP, stroke: 'none', rx: 0
        });
      });
      arrivals = arrivals.filter(function (t) { return t > t0; });
      arrivals.forEach(function (t) {
        var x = SX + (t - t0) / WIN * SW;
        S('path', { d: 'M' + (x - 6) + ',' + (SY - 13) + ' L' + (x + 6) + ',' + (SY - 13) + ' L' + x + ',' + (SY - 2) + ' z', fill: C.accent }, arrG);
      });
      var idle = check + sleep, burn = idle > 0 ? check / idle : 0;
      cpu.set(burn);
      cpuText.textContent = (burn > 0.995 ? '100' : (burn * 100).toFixed(1)) + '%';
      var avg = lat.length ? lat.reduce(function (s, x) { return s + x; }, 0) / lat.length : 0;
      latText.textContent = lat.length ? (avg < 0.1 ? '≈ 0 ms' : avg.toFixed(1) + ' ms') : '—';
      var n = queue.length;
      qText.textContent = n + (n === 1 ? ' job waiting' : ' jobs waiting');
      for (var i = 0; i < Math.min(n, 6); i++) R(qG, 40 + i * 28, SY + 10, 22, SH - 20, { fill: C.warnFill, stroke: C.warn, sw: 1.2 });
    }
    function reset() {
      segs = []; arrivals = []; queue = []; lat = [];
      w.st = 'idle'; w.nextCheck = now + POLL; nextArr = now + 20;
      measText.textContent = 'Measured on an M4 laptop: ' + MEAS[mode];
    }
    Toggle(ctl, [['spin', 'spin'], ['poll', 'sleep 10 ms'], ['cv', 'condition variable']], mode, function (v) { mode = v; reset(); });
    Btn(ctl, '＋ job now', function () { nextArr = now; });
    reset();
    return { tick: function (dt) { advance(now + dt * SPEED); draw(); } };
  });

  /* Producers → bounded (or unbounded) queue → consumers. */
  function QueueSim(host, opt) {
    var W = 1140, H = opt.H, svg = mkSvg(host, W, H), ctl = Controls(host);
    var cfg = { np: opt.np, nc: opt.nc, cap: opt.cap, pr: opt.pr, cr: opt.cr, unbounded: !!opt.unbounded };
    var OOM = 60, MB = 64, QX = 300, QW = 540, PX = 120, CX = 1020, BW = 170, BH = 38;
    var foot = opt.mem ? 64 : opt.hint ? 52 : 34;
    var midY = 30 + (H - 30 - foot) / 2;
    var layer = G(svg), itemsG = G(svg);
    var banner = T(svg, QX + QW / 2, midY + 62, '', { size: 17, weight: 700, fill: C.bad });
    var rateText = T(svg, QX + QW / 2, H - (opt.mem ? 40 : opt.hint ? 28 : 10), '', { size: 14, mono: true, fill: C.dim });
    var memBar, memText;
    if (opt.mem) {
      T(svg, QX - 10, H - 13, 'memory', { anchor: 'end', size: 13, fill: C.dim });
      memBar = bar(svg, QX, H - 25, QW, 14, C.bad);
      memText = T(svg, QX + QW + 10, H - 13, '', { anchor: 'start', size: 13, mono: true });
    }
    if (opt.hint) T(svg, W / 2, H - 4, opt.hint, { size: 13, fill: C.dim });
    var st, P = [], K = [], slotW = 0, slotStart = 0;

    function place(g, x, y, instant) {
      if (instant) g.style.transition = 'none';
      g.style.transform = 'translate(' + x + 'px,' + y + 'px)';
      if (instant) { void g.getBoundingClientRect(); g.style.transition = ''; }
    }
    function ys(n) {
      var sp = Math.min(54, (H - 40 - foot) / n);
      return function (i) { return midY + (i - (n - 1) / 2) * sp; };
    }
    function agent(x, y, name, color) {
      var g = G(layer);
      var box = R(g, x - BW / 2, y - BH / 2, BW, BH, { fill: FILL[color], stroke: color, sw: 1.6, rx: 4 });
      T(g, x - BW / 2 + 12, y + 5, name, { anchor: 'start', size: 14, weight: 700, mono: true });
      var s = T(g, x + BW / 2 - 10, y + 5, '', { anchor: 'end', size: 13, fill: C.dim });
      var pb = bar(g, x - BW / 2 + 6, y + BH / 2 - 6, BW - 12, 3, color);
      pb.el.style.transition = 'none';
      return { g: g, box: box, s: s, pb: pb, color: color, x: x, y: y, look: '' };
    }
    function look(a, state, text) {
      if (a.s.textContent !== text) a.s.textContent = text;
      if (a.look === state) return;
      a.look = state;
      var waiting = state === 'wait', off = state === 'done';
      a.box.setAttribute('stroke-dasharray', waiting ? '5 4' : 'none');
      a.box.setAttribute('stroke', state === 'bad' ? C.bad : waiting || off ? C.dim : a.color);
      a.box.setAttribute('fill', state === 'bad' ? C.badFill : waiting || off ? C.bg : FILL[a.color]);
      a.g.style.opacity = off ? 0.5 : 1;
    }
    function slotX(i, n) {
      if (!cfg.unbounded) return slotStart + (cfg.cap - 1 - i) * (slotW + 6) + slotW / 2;
      var sp = Math.min(40, (QW - 30) / Math.max(1, n - 1));
      return QX + QW - 16 - i * sp;
    }
    function layoutQueue() {
      var n = st.q.length;
      st.q.forEach(function (it, i) { place(it.g, slotX(i, n), midY); });
    }
    function mkItem(id, x, y) {
      var g = G(itemsG, 'mv');
      R(g, -13, -13, 26, 26, { fill: C.warnFill, stroke: C.warn, sw: 1.3, rx: 3 });
      T(g, 0, 5, String(id % 100), { size: 11, mono: true, weight: 600 });
      place(g, x, y, true);
      return g;
    }
    function dur(rate) { return rnd(0.6, 1.4) / rate; }

    function reset() {
      clear(layer); clear(itemsG);
      banner.textContent = '';
      hdr(layer, PX, 18, 'PRODUCERS');
      hdr(layer, CX, 18, 'CONSUMERS');
      hdr(layer, QX + QW / 2, 18, cfg.unbounded ? 'QUEUE · UNBOUNDED' : 'QUEUE · CAPACITY ' + cfg.cap);
      if (cfg.unbounded) {
        R(layer, QX, midY - 22, QW, 44, { fill: C.bg, stroke: C.line, dash: '6 4' });
      } else {
        slotW = Math.min(46, QW / cfg.cap - 6);
        slotStart = QX + (QW - (cfg.cap * (slotW + 6) - 6)) / 2;
        for (var s = 0; s < cfg.cap; s++) R(layer, slotStart + s * (slotW + 6), midY - 22, slotW, 44, { fill: C.bg, stroke: C.line, dash: '4 3' });
      }
      L(layer, QX - 12, midY, QX - 2, midY, { stroke: C.dim, end: svg.arrow.dim });
      L(layer, QX + QW + 2, midY, QX + QW + 12, midY, { stroke: C.dim, end: svg.arrow.dim });
      var py = ys(cfg.np), cy = ys(cfg.nc);
      P = []; K = [];
      for (var i = 0; i < cfg.np; i++) {
        var t0 = dur(cfg.pr);
        P.push({ a: agent(PX, py(i), 'P' + (i + 1), C.good), s: 'make', left: t0 * rnd(0.3, 1), total: t0, since: 0 });
      }
      for (var j = 0; j < cfg.nc; j++) K.push({ a: agent(CX, cy(j), 'C' + (j + 1), C.accent), s: 'wait', left: 0, total: 1, since: j });
      st = { t: 0, q: [], closed: false, oom: false, seq: 0, made: [], used: [] };
      draw();
    }
    function wake(list, state, fn) {        // notify_one: the longest sleeper wakes
      var best = -1;
      list.forEach(function (x, i) { if (x.s === state && (best < 0 || x.since < list[best].since)) best = i; });
      if (best >= 0) { flash(list[best].a.box); fn(best); }
    }
    function tryPush(i) {
      var p = P[i];
      if (st.closed) { p.s = 'closed'; return; }
      if (cfg.unbounded || st.q.length < cfg.cap) {
        st.seq++;
        st.q.push({ g: mkItem(st.seq, p.a.x + BW / 2 - 14, p.a.y) });
        st.made.push(st.t);
        layoutQueue();
        p.s = 'make'; p.total = p.left = dur(cfg.pr);
        if (cfg.unbounded && opt.mem && st.q.length >= OOM) st.oom = true;
        wake(K, 'wait', tryPop);
      } else { p.s = 'full'; p.since = st.t; }
    }
    function tryPop(i) {
      var k = K[i];
      if (st.q.length) {
        var it = st.q.shift();
        st.used.push(st.t);
        place(it.g, k.a.x - BW / 2 + 18, k.a.y);
        it.g.style.opacity = 0;
        setTimeout(function () { if (it.g.parentNode) it.g.parentNode.removeChild(it.g); }, 600);
        layoutQueue();
        k.s = 'work'; k.total = k.left = dur(cfg.cr);
        wake(P, 'full', tryPush);
      } else if (st.closed) k.s = 'done';
      else { k.s = 'wait'; k.since = st.t; }
    }
    function closeQ() {
      if (st.closed || st.oom) return;
      st.closed = true;
      P.forEach(function (p) { if (p.s === 'full') p.s = 'closed'; });
      K.forEach(function (k, i) { if (k.s === 'wait') tryPop(i); });      // notify_all: everyone re-checks
    }
    function rate(list) {
      while (list.length && list[0] < st.t - 3) list.shift();
      return (list.length / Math.min(3, Math.max(st.t, 0.5))).toFixed(1);
    }
    function draw() {
      P.forEach(function (p) {
        if (st.oom) look(p.a, 'bad', 'killed');
        else if (p.s === 'make') look(p.a, 'run', 'making');
        else if (p.s === 'full') look(p.a, 'wait', 'waits: notFull');
        else look(p.a, 'done', 'closed → false');
        p.a.pb.set(p.s === 'make' ? 1 - p.left / p.total : 0);
      });
      K.forEach(function (k) {
        if (st.oom) look(k.a, 'bad', 'killed');
        else if (k.s === 'work') look(k.a, 'run', 'working');
        else if (k.s === 'wait') look(k.a, 'wait', 'waits: notEmpty');
        else look(k.a, 'done', 'nullopt → exit');
        k.a.pb.set(k.s === 'work' ? 1 - k.left / k.total : 0);
      });
      rateText.textContent = 'made ' + rate(st.made) + '/s · used ' + rate(st.used) + '/s · in the queue: ' + st.q.length +
        (st.closed ? ' · closed' : '');
      if (opt.mem) {
        memBar.set(st.q.length / OOM);
        memText.textContent = (st.q.length * MB / 1024).toFixed(1) + ' GB';
      }
      if (st.oom) banner.textContent = 'OUT OF MEMORY — the process is killed';
    }
    function tick(dt) {
      if (!st.oom) {
        st.t += dt;
        P.forEach(function (p, i) { if (p.s === 'make') { p.left -= dt; if (p.left <= 0) tryPush(i); } });
        K.forEach(function (k, i) { if (k.s === 'work') { k.left -= dt; if (k.left <= 0) tryPop(i); } });
      }
      draw();
    }

    var speed = function (v) { return v + '/s'; };
    if (opt.controls === 'why') {
      Toggle(ctl, [['u', 'unbounded'], ['b', 'bounded · 8']], cfg.unbounded ? 'u' : 'b', function (v) { cfg.unbounded = v === 'u'; reset(); });
      Slider(ctl, 'producer', 0.5, 6, 0.5, cfg.pr, function (v) { cfg.pr = v; }, speed);
      Slider(ctl, 'consumer', 0.5, 6, 0.5, cfg.cr, function (v) { cfg.cr = v; }, speed);
    } else {
      Slider(ctl, 'producers', 1, 4, 1, cfg.np, function (v) { cfg.np = v; reset(); });
      Slider(ctl, 'consumers', 1, 4, 1, cfg.nc, function (v) { cfg.nc = v; reset(); });
      Slider(ctl, 'capacity', 1, 10, 1, cfg.cap, function (v) { cfg.cap = v; reset(); });
      Slider(ctl, 'make', 0.5, 5, 0.5, cfg.pr, function (v) { cfg.pr = v; }, speed);
      Slider(ctl, 'use', 0.5, 5, 0.5, cfg.cr, function (v) { cfg.cr = v; }, speed);
      Btn(ctl, 'close()', closeQ);
    }
    Btn(ctl, 'reset', reset);
    reset();
    return { tick: tick };
  }

  def('whyBounded', function (host) {
    return QueueSim(host, { H: 252, np: 1, nc: 1, cap: 8, pr: 4, cr: 2, unbounded: true, mem: true, controls: 'why' });
  });
  def('queueSim', function (host) {
    return QueueSim(host, {
      H: 300, np: 2, nc: 2, cap: 6, pr: 1.5, cr: 1, controls: 'full',
      hint: 'Try: faster makers → producers wait on notFull · faster users → consumers wait on notEmpty · close() → drain, then nullopt'
    });
  });

  /* A semaphore as a parking lot. */
  def('semaphore', function (host) {
    var W = 1140, H = 250, svg = mkSvg(host, W, H), ctl = Controls(host);
    var N = 3, auto = true, armed = false, seq = 0, arrT = 0.3;
    var GX = 372, BY = 112, BX0 = 430, BWd = 92, BG = 12;
    hdr(svg, 190, 18, 'WAITING · acquire() BLOCKS');
    hdr(svg, 740, 18, 'INSIDE · ONE PERMIT EACH');
    L(svg, GX, 40, GX, 188, { stroke: C.accent, sw: 2, dash: '6 5' });
    T(svg, GX, 208, 'acquire()', { size: 13, mono: true, fill: C.accent });
    T(svg, 1098, 208, 'release()', { size: 13, mono: true, fill: C.good });
    var countText = T(svg, 740, 208, '', { size: 15, mono: true, weight: 600 });
    var banner = T(svg, 570, 240, '', { size: 15, weight: 700, fill: C.bad });
    var moreText = T(svg, 40, BY + 46, '', { anchor: 'start', size: 13, fill: C.dim });
    var baysG = G(svg), toksG = G(svg);
    var bays = [], waiting = [], holders = [], throwBtn;

    function build() {
      clear(baysG); clear(toksG);
      bays = []; waiting = []; holders = []; armed = false;
      if (throwBtn) throwBtn.classList.remove('on');
      for (var i = 0; i < N; i++) {
        var x = BX0 + i * (BWd + BG), g = G(baysG);
        var r = R(g, x, BY - 36, BWd, 72, { fill: 'rgba(20,56,90,.4)', stroke: C.line, dash: '4 4' });
        var lab = T(g, x + BWd / 2, BY + 30, 'permit ' + (i + 1), { size: 11, fill: C.dim });
        bays.push({ x: x + BWd / 2, x0: x, holder: null, leaked: false, g: g, r: r, lab: lab });
      }
    }
    function layoutWaiting() {
      waiting.forEach(function (th, i) {
        th.tok.at(GX - 40 - Math.min(i, 4) * 64, BY).show(i < 5).state('wait').note('');
      });
      moreText.textContent = waiting.length > 5 ? '+' + (waiting.length - 5) + ' more' : '';
    }
    function newThread() {
      seq++;
      var tok = Token(toksG, 'T' + seq, PALETTE[seq % PALETTE.length], { w: 54, h: 30, size: 13 });
      tok.jump(-40, BY);
      waiting.push({ tok: tok, left: 0, bay: null });
      layoutWaiting();
    }
    function admit() {
      bays.forEach(function (b) {
        if (b.holder || b.leaked || !waiting.length) return;
        var th = waiting.shift();
        b.holder = th; th.bay = b; th.left = rnd(1.8, 3.6);
        holders.push(th);
        th.tok.at(b.x, BY - 4).show(1).state('run').note('');
        layoutWaiting();
      });
    }
    function finish(th) {
      holders.splice(holders.indexOf(th), 1);
      var b = th.bay;
      b.holder = null;
      if (armed) {
        armed = false;
        throwBtn.classList.remove('on');
        b.leaked = true;
        b.r.setAttribute('stroke', C.bad); b.r.setAttribute('fill', C.badFill); b.r.setAttribute('stroke-dasharray', 'none');
        L(b.g, b.x0 + 14, BY - 24, b.x0 + BWd - 14, BY + 14, { stroke: C.bad, sw: 2.5 });
        L(b.g, b.x0 + BWd - 14, BY - 24, b.x0 + 14, BY + 14, { stroke: C.bad, sw: 2.5 });
        b.lab.textContent = 'leaked'; b.lab.setAttribute('fill', C.bad);
        th.tok.state('bad').note('');
        setTimeout(function () { th.tok.at(1110, BY - 4); }, 600);
        setTimeout(function () { th.tok.show(0); }, 1000);
        setTimeout(function () { th.tok.remove(); }, 1600);
      } else {
        th.tok.at(1110, BY - 4);
        setTimeout(function () { th.tok.show(0); }, 350);
        setTimeout(function () { th.tok.remove(); }, 1000);
      }
    }
    Slider(ctl, 'permits', 1, 6, 1, N, function (v) { N = v; build(); });
    Btn(ctl, '＋ thread', newThread);
    Toggle(ctl, [['on', 'arrivals on'], ['off', 'off']], 'on', function (v) { auto = v === 'on'; });
    throwBtn = Btn(ctl, '💥 next one throws', function () { armed = true; throwBtn.classList.add('on'); });
    Btn(ctl, 'reset', build);
    build();
    return {
      tick: function (dt) {
        if (auto) { arrT -= dt; if (arrT <= 0) { if (waiting.length < 9) newThread(); arrT = rnd(0.5, 1.1); } }
        holders.slice().forEach(function (th) { th.left -= dt; if (th.left <= 0) finish(th); });
        admit();
        var free = bays.filter(function (b) { return !b.holder && !b.leaked; }).length;
        var leaked = bays.filter(function (b) { return b.leaked; }).length;
        countText.textContent = 'free permits: ' + free + ' / ' + N + (leaked ? ' · leaked: ' + leaked : '');
        banner.textContent = leaked === N ? 'Every permit leaked — nobody gets in, ever. No error, no crash.' : '';
      }
    };
  });

  /* A semaphore keeps an early signal; a condition variable drops it. */
  def('remembers', function (host) {
    var W = 1140, H = 250, svg = mkSvg(host, W, H);
    hdr(svg, 270, 18, 'SEMAPHORE · STARTS AT 0');
    hdr(svg, 870, 18, 'CONDITION VARIABLE');
    L(svg, 570, 30, 570, 230, { stroke: C.faint, sw: 1, dash: '4 6' });
    R(svg, 215, 40, 110, 80, { fill: C.panel });
    T(svg, 270, 60, 'count', { size: 12, fill: C.dim });
    var count = T(svg, 270, 104, '0', { size: 34, mono: true, weight: 700 });
    R(svg, 770, 34, 200, 112, { fill: 'rgba(20,56,90,.3)', dash: '5 5' });
    T(svg, 870, 54, 'waiting room', { size: 12, fill: C.dim });
    var gone = T(svg, 870, 96, '', { size: 14, weight: 700, fill: C.bad });
    var sigL = Pth(svg, 'M432,168 C400,120 360,90 330,82', { stroke: C.good, dash: '5 4', end: svg.arrow.good });
    var sigR = Pth(svg, 'M1032,168 C1010,130 1000,100 974,90', { stroke: C.good, dash: '5 4', end: svg.arrow.good });
    sigL.classList.add('fd'); sigR.classList.add('fd');
    var A1 = Token(svg, 'A', C.accent, { w: 56 }), B1 = Token(svg, 'B', C.good, { w: 56 });
    var A2 = Token(svg, 'A', C.accent, { w: 56 }), B2 = Token(svg, 'B', C.good, { w: 56 });
    var resL = T(svg, 270, 242, '', { size: 15, weight: 700, fill: C.good });
    var resR = T(svg, 870, 242, '', { size: 15, weight: 700, fill: C.bad });
    return {
      steps: 4,
      render: function (k) {
        A1.at(90, 185).state('run').note('');
        B1.at(460, 185).state('run').note('');
        A2.at(690, 185).state('run').note('');
        B2.at(1060, 185).state('run').note('');
        count.textContent = '0'; gone.textContent = ''; resL.textContent = ''; resR.textContent = '';
        fade(sigL, 0); fade(sigR, 0);
        if (k >= 1) { B1.note('release()'); B2.note('notify_one()'); fade(sigL, 1); fade(sigR, 1); }
        if (k >= 2) { count.textContent = '1'; gone.textContent = 'nobody here: gone'; }
        if (k >= 3) {
          fade(sigL, 0); fade(sigR, 0); gone.textContent = '';
          count.textContent = '0';
          A1.at(270, 185).state('good').note('acquire(): goes on');
          A2.at(870, 88).state('sleep').note('wait(): asleep');
        }
        if (k >= 4) { resL.textContent = 'the count remembered'; resR.textContent = 'the wake-up is gone'; }
      }
    };
  });

  /* Readers–writer lock with a policy switch. */
  def('rwlock', function (host) {
    var W = 1140, H = 256, svg = mkSvg(host, W, H), ctl = Controls(host);
    var RX = 450, RW = 340, RY = 34, RH = 170, LY = 120;
    hdr(svg, 220, 18, 'WAITING');
    hdr(svg, RX + RW / 2, 18, 'THE DATA · shared_mutex');
    R(svg, RX, RY, RW, RH, { fill: 'rgba(20,56,90,.5)', stroke: C.accent, sw: 1.8 });
    var inside = T(svg, RX + RW / 2, RY + RH + 22, '', { size: 14, mono: true, weight: 600 });
    var hint = T(svg, 20, H - 6, '', { anchor: 'start', size: 13, fill: C.dim });
    var toks = G(svg);
    var policy = 'readers', storm = false, stormT = 0, seq = 0, waiting = [], inRoom = [];
    function isW(t) { return t.type === 'W'; }
    function slotPos(i) { return [RX + 50 + (i % 4) * 80, RY + 52 + Math.floor(i / 4) * 70]; }
    function layout() {
      waiting.forEach(function (th, i) { th.tok.at(RX - 70 - Math.min(i, 6) * 50, LY).show(i < 7).state('wait'); });
      inRoom.forEach(function (th) {
        if (isW(th)) th.tok.at(RX + RW / 2, RY + RH / 2).state('run').note('');
        else { var p = slotPos(th.slot); th.tok.at(p[0], p[1]).state('run').note(''); }
      });
    }
    function add(type) {
      seq++;
      var tok = Token(toks, type, type === 'W' ? C.bad : C.accent, { w: type === 'W' ? 56 : 40, h: 30, size: 14 });
      tok.jump(-40, LY);
      waiting.push({ type: type, tok: tok, left: 0, since: 0, slot: -1 });
      layout();
    }
    function freeSlot() {
      for (var i = 0; i < 8; i++) if (!inRoom.some(function (t) { return t.slot === i; })) return i;
      return -1;
    }
    function admit() {
      for (var guard = 0; guard < 20; guard++) {
        var wIn = inRoom.some(isW), wWait = waiting.some(isW), firstW = null, done = true;
        for (var j = 0; j < waiting.length; j++) if (isW(waiting[j])) { firstW = waiting[j]; break; }
        for (var i = 0; i < waiting.length; i++) {
          var th = waiting[i];
          var ok = isW(th) ? inRoom.length === 0 && th === firstW
            : !wIn && !(policy === 'writer' && wWait) && freeSlot() >= 0;
          if (ok) {
            waiting.splice(i, 1);
            th.slot = isW(th) ? -1 : freeSlot();
            th.left = isW(th) ? 1.4 : rnd(1.3, 2.3);
            inRoom.push(th);
            done = false;
            break;
          }
        }
        if (done) break;
      }
      layout();
    }
    function reset() {
      clear(toks); waiting = []; inRoom = []; storm = false; stormToggle.set('off');
      add('R'); add('R'); add('W'); add('R');
    }
    Btn(ctl, '＋ reader', function () { add('R'); });
    Btn(ctl, '＋ writer', function () { add('W'); });
    var stormToggle = Toggle(ctl, [['off', 'reader storm off'], ['on', 'on']], 'off', function (v) { storm = v === 'on'; });
    Toggle(ctl, [['readers', 'readers first'], ['writer', 'writer first']], policy, function (v) { policy = v; });
    Btn(ctl, 'reset', reset);
    reset();
    return {
      tick: function (dt) {
        if (storm) { stormT -= dt; if (stormT <= 0) { if (waiting.length < 12) add('R'); stormT = rnd(0.3, 0.55); } }
        inRoom.slice().forEach(function (th) {
          th.left -= dt;
          if (th.left > 0) return;
          inRoom.splice(inRoom.indexOf(th), 1);
          th.tok.at(RX + RW + 70, th.tok.y);
          setTimeout(function () { th.tok.show(0); }, 300);
          setTimeout(function () { th.tok.remove(); }, 900);
        });
        waiting.forEach(function (th) {
          th.since += dt;
          if (isW(th)) th.tok.note('waited ' + th.since.toFixed(1) + ' s', th.since > 3 ? C.bad : C.dim);
        });
        admit();
        var rs = inRoom.filter(function (t) { return !isW(t); }).length;
        inside.textContent = inRoom.some(isW) ? 'inside: one writer, alone' : rs ? 'inside: ' + rs + (rs === 1 ? ' reader' : ' readers') : 'inside: nobody';
        hint.textContent = policy === 'readers'
          ? 'readers first: a waiting writer does not stop new readers from coming in'
          : 'writer first: a waiting writer stops new readers until it has had its turn';
      }
    };
  });

  def('latch', function (host) {
    var W = 1140, H = 150, svg = mkSvg(host, W, H);
    var M = Token(svg, 'main', C.accent, { w: 80 });
    R(svg, 330, 26, 170, 96, { fill: C.panel });
    T(svg, 415, 48, 'latch · count', { size: 12, fill: C.dim });
    var cnt = T(svg, 415, 100, '3', { size: 38, mono: true, weight: 700 });
    L(svg, 505, 74, 610, 74, { stroke: C.dim, dash: '4 4' });
    T(svg, 558, 66, 'count_down()', { size: 12, mono: true, fill: C.dim });
    var bars = [0, 1, 2].map(function (i) {
      var y = 26 + i * 38;
      T(svg, 690, y + 15, 'shard ' + i, { anchor: 'end', size: 13, mono: true, fill: C.dim });
      var b = bar(svg, 705, y + 3, 320, 16, C.good);
      var ok = T(svg, 1040, y + 16, '', { anchor: 'start', size: 14, weight: 700, fill: C.good });
      return { b: b, ok: ok };
    });
    var P = [[0.12, 0.2, 0.08], [0.4, 0.55, 0.7], [0.7, 0.85, 1], [1, 0.92, 1], [1, 1, 1]];
    return {
      steps: 4,
      render: function (k) {
        var p = P[k], done = 0;
        bars.forEach(function (b, i) { b.b.set(p[i]); b.ok.textContent = p[i] >= 1 ? 'done ✓' : ''; if (p[i] >= 1) done++; });
        cnt.textContent = String(3 - done);
        cnt.setAttribute('fill', done === 3 ? C.good : C.fg);
        M.at(120, 70).state('run').note('');
        if (k >= 1) M.state('sleep').note('loaded.wait()');
        if (k >= 4) M.state('good').note('startServing()');
      }
    };
  });

  /* Four threads, one barrier per simulation step. */
  def('barrier', function (host) {
    var W = 1140, H = 252, svg = mkSvg(host, W, H), ctl = Controls(host);
    var X0 = 150, X1 = 940, lanes = [56, 104, 152, 200];
    var on = true, running = true, step = 1, stale = 0;
    hdr(svg, X0, 18, 'ONE STEP OF THE SIMULATION  →', 'start');
    var barLine = L(svg, X1 + 34, 30, X1 + 34, 222, { stroke: C.accent, sw: 3, dash: '7 5' });
    var barText = T(svg, X1 + 34, 22, 'arrive_and_wait()', { size: 13, mono: true, fill: C.accent });
    var stepText = T(svg, 20, H - 6, '', { anchor: 'start', size: 14, mono: true, weight: 600 });
    var doneText = T(svg, 560, H - 6, '', { size: 14, mono: true, fill: C.good });
    var staleText = T(svg, 1120, H - 6, '', { anchor: 'end', size: 14, weight: 700, fill: C.bad });
    var th = lanes.map(function (y, i) {
      L(svg, X0, y, X1, y, { stroke: C.faint, sw: 1 });
      T(svg, 20, y + 5, 'thread ' + (i + 1), { anchor: 'start', size: 13, fill: C.dim, mono: true });
      var tok = Token(svg, 'T' + (i + 1), PALETTE[i], { w: 46, h: 28, size: 13, still: true });
      var lab = T(svg, X1 + 60, y + 5, '', { anchor: 'start', size: 13, mono: true, fill: C.dim });
      return { tok: tok, p: 0, d: rnd(1, 2.6), step: 1, wait: false, lab: lab, y: y, badT: 0 };
    });
    function reset() {
      step = 1; stale = 0; doneText.textContent = '';
      th.forEach(function (t) { t.p = 0; t.d = rnd(1, 2.6); t.step = 1; t.wait = false; t.badT = 0; t.tok.state('run').note(''); });
      fade(barLine, on); fade(barText, on);
    }
    Toggle(ctl, [['on', 'with barrier'], ['off', 'no barrier']], 'on', function (v) { on = v === 'on'; reset(); });
    var runBtn = Btn(ctl, '⏸ pause', function () { running = !running; runBtn.textContent = running ? '⏸ pause' : '▶ run'; });
    Btn(ctl, 'reset', reset);
    reset();
    return {
      tick: function (dt) {
        if (running) {
          th.forEach(function (t, i) {
            if (t.badT > 0) { t.badT -= dt; if (t.badT <= 0) t.tok.state('run').note(''); }
            if (t.wait) return;
            t.p += dt / t.d;
            if (t.p < 1) return;
            if (on) { t.p = 1; t.wait = true; t.tok.state('wait').note('waiting'); return; }
            t.p = 0; t.step++; t.d = rnd(1, 2.6);
            [i - 1, i + 1].forEach(function (j) {        // it now reads its neighbours' cells
              if (j < 0 || j >= th.length || th[j].step === t.step) return;
              stale++; t.badT = 0.7; t.tok.state('bad').note('stale read', C.bad);
            });
          });
          if (on && th.every(function (t) { return t.wait; })) {
            doneText.textContent = 'swapBuffers() · step ' + step + ' done';
            flash(doneText);
            step++;
            th.forEach(function (t) { t.p = 0; t.wait = false; t.d = rnd(1, 2.6); t.step = step; t.tok.state('run').note(''); });
          }
        }
        th.forEach(function (t) {
          t.tok.at(X0 + t.p * (X1 - X0), t.y);
          t.lab.textContent = on ? '' : 'step ' + t.step;
        });
        stepText.textContent = on ? 'every thread is on step ' + step : 'each thread on its own step';
        staleText.textContent = on ? '' : 'reads of half-updated data: ' + stale;
      }
    };
  });

  /* ================= after the break: liveness and async ================= */

  /* Wait-for graph for the two transfers. */
  def('waitFor', function (host) {
    var W = 1140, H = 250, svg = mkSvg(host, W, H);
    function thread(x, y, label, col) {
      S('circle', { cx: x, cy: y, r: 42, fill: FILL[col], stroke: col, 'stroke-width': 2 }, svg);
      T(svg, x, y + 6, label, { size: 17, weight: 700 });
    }
    function lock(x, y, label) {
      R(svg, x - 40, y - 22, 80, 44, { fill: C.panel, stroke: C.line, sw: 1.6 });
      T(svg, x, y + 5, label, { size: 15, mono: true, weight: 600 });
    }
    function edge(d, col, dash, lx, ly, label) {
      var g = G(svg, 'fd');
      var p = Pth(g, d, { stroke: col, sw: 2.4, dash: dash, end: svg.arrow[col === C.accent ? 'accent' : 'good'] });
      p.classList.add('col');
      var t = T(g, lx, ly, label, { size: 13, fill: col });
      g.style.opacity = 0;
      return { g: g, p: p, t: t, col: col, name: col === C.accent ? 'accent' : 'good' };
    }
    var es = [
      edge('M528,44 Q410,30 338,96', C.accent, null, 420, 30, 'x is held by A'),
      edge('M612,206 Q730,220 802,154', C.good, null, 720, 232, 'y is held by B'),
      edge('M336,154 Q410,214 526,206', C.accent, '7 5', 400, 222, 'A waits for y'),
      edge('M804,96 Q730,30 614,44', C.good, '7 5', 740, 30, 'B waits for x')
    ];
    thread(300, 125, 'A', C.accent);
    thread(840, 125, 'B', C.good);
    lock(570, 44, 'lock x');
    lock(570, 206, 'lock y');
    var cyc = T(svg, 570, 132, 'a cycle = deadlock', { size: 20, weight: 700, fill: C.bad });
    cyc.classList.add('fd');
    return {
      steps: 4,
      render: function (k) {
        es.forEach(function (e, i) {
          e.g.style.opacity = i < k ? 1 : 0;
          var hot = k >= 4;
          e.p.setAttribute('stroke', hot ? C.bad : e.col);
          e.p.setAttribute('marker-end', svg.arrow[hot ? 'bad' : e.name]);
          e.t.setAttribute('fill', hot ? C.bad : e.col);
        });
        fade(cyc, k >= 4);
      }
    };
  });

  /* Two threads, two locks, lines for "holds" and "waits for". */
  function TwoLocks(svg) {
    var TX = 190, LXp = 650, ys = { A: 78, B: 188 }, ly = { x: 78, y: 188 };
    var lines = {};
    ['A', 'B'].forEach(function (t) {
      ['x', 'y'].forEach(function (l) {
        lines[t + l] = L(svg, TX + 62, ys[t], LXp - 42, ly[l], { stroke: C.dim, sw: 3 });
        lines[t + l].style.opacity = 0;
      });
    });
    var locks = {};
    ['x', 'y'].forEach(function (l) {
      var r = R(svg, LXp - 40, ly[l] - 24, 80, 48, { fill: C.panel, stroke: C.line, sw: 1.6 });
      T(svg, LXp, ly[l] + 5, 'lock ' + l, { size: 15, mono: true, weight: 600 });
      var who = T(svg, LXp + 52, ly[l] + 5, '', { anchor: 'start', size: 13, fill: C.dim });
      locks[l] = { r: r, who: who };
    });
    var code = {
      A: T(svg, TX, ys.A - 30, '', { size: 13, mono: true, fill: C.dim }),
      B: T(svg, TX, ys.B - 30, '', { size: 13, mono: true, fill: C.dim })
    };
    var toks = {
      A: Token(svg, 'thread A', C.accent, { w: 124, still: true }).at(TX, ys.A),
      B: Token(svg, 'thread B', C.good, { w: 124, still: true }).at(TX, ys.B)
    };
    var cols = { A: C.accent, B: C.good };
    return {
      toks: toks, code: code,
      line: function (t, l, how) {        // how: '' | 'hold' | 'want' | 'fail' | 'dead'
        var e = lines[t + l];
        e.style.opacity = how ? 1 : 0;
        e.setAttribute('stroke', how === 'fail' || how === 'dead' ? C.bad : how === 'want' ? C.dim : cols[t]);
        e.setAttribute('stroke-dasharray', how === 'hold' ? 'none' : '7 5');
      },
      holder: function (l, t) {
        locks[l].r.setAttribute('fill', t ? FILL[cols[t]] : C.panel);
        locks[l].r.setAttribute('stroke', t ? cols[t] : C.line);
        locks[l].who.textContent = t ? 'held by ' + t : 'free';
      }
    };
  }

  /* Opposite lock order freezes; one global order doesn't. */
  def('deadlockLab', function (host) {
    var W = 1140, H = 262, svg = mkSvg(host, W, H), ctl = Controls(host);
    var sc = TwoLocks(svg);
    var banner = T(svg, 570, H - 8, '', { size: 16, weight: 700 });
    var count = { A: T(svg, 1110, 82, '', { anchor: 'end', size: 15, mono: true }), B: T(svg, 1110, 192, '', { anchor: 'end', size: 15, mono: true }) };
    var mode = 'naive', th, locks, dead, still;
    function reset() {
      locks = { x: null, y: null };
      th = {
        A: { first: 'x', second: 'y', ph: 'idle', left: rnd(0.2, 0.6), done: 0 },
        B: { first: mode === 'naive' ? 'y' : 'x', second: mode === 'naive' ? 'x' : 'y', ph: 'idle', left: rnd(0.2, 0.6), done: 0 }
      };
      dead = false; still = 0;
      sc.code.A.textContent = 'transfer(x, y)';
      sc.code.B.textContent = mode === 'naive' ? 'transfer(y, x)' : 'transfer(y, x) · x first';
      draw();
    }
    function stepThread(name, t, dt) {
      switch (t.ph) {
        case 'idle': t.left -= dt; if (t.left <= 0) t.ph = 'want1'; break;
        case 'want1': if (!locks[t.first]) { locks[t.first] = name; t.ph = 'hold1'; t.left = rnd(0.1, 0.25); } break;
        case 'hold1': t.left -= dt; if (t.left <= 0) t.ph = 'want2'; break;
        case 'want2': if (!locks[t.second]) { locks[t.second] = name; t.ph = 'hold2'; t.left = 0.35; } break;
        case 'hold2':
          t.left -= dt;
          if (t.left <= 0) { locks.x = locks.x === name ? null : locks.x; locks.y = locks.y === name ? null : locks.y; t.done++; t.ph = 'idle'; t.left = rnd(0.5, 1.6); }
          break;
      }
    }
    function draw() {
      ['A', 'B'].forEach(function (n) {
        var t = th[n];
        ['x', 'y'].forEach(function (l) {
          var how = locks[l] === n ? 'hold' : (t.ph === 'want1' && t.first === l) || (t.ph === 'want2' && t.second === l) ? (dead ? 'dead' : 'want') : '';
          sc.line(n, l, how);
        });
        var note = t.ph === 'idle' ? 'between transfers' : t.ph === 'want1' ? 'waits for ' + t.first
          : t.ph === 'hold1' ? 'holds ' + t.first : t.ph === 'want2' ? 'holds ' + t.first + ' · waits for ' + t.second : 'holds both · moves money';
        sc.toks[n].state(dead ? 'sleep' : t.ph.indexOf('want') === 0 ? 'wait' : 'run').note(note, dead ? C.bad : C.dim);
        count[n].textContent = t.done + ' transfers';
      });
      sc.holder('x', locks.x); sc.holder('y', locks.y);
      if (dead) { banner.textContent = 'DEADLOCK — both threads asleep, forever · CPU 0%'; banner.setAttribute('fill', C.bad); }
      else { banner.textContent = mode === 'naive' ? 'running… for now' : 'one global order: x before y — no cycle can form'; banner.setAttribute('fill', mode === 'naive' ? C.dim : C.good); }
    }
    Toggle(ctl, [['naive', 'opposite order'], ['ordered', 'always x first']], mode, function (v) { mode = v; reset(); });
    Btn(ctl, 'restart', reset);
    reset();
    return {
      tick: function (dt) {
        if (!dead) {
          stepThread('A', th.A, dt);
          stepThread('B', th.B, dt);
          if (th.A.ph === 'want2' && th.B.ph === 'want2' && locks[th.A.second] === 'B' && locks[th.B.second] === 'A') dead = true;
        }
        draw();
      }
    };
  });

  /* try_lock + back off: in lockstep forever, or unstuck by jitter. */
  def('livelock', function (host) {
    var W = 1140, H = 262, svg = mkSvg(host, W, H), ctl = Controls(host);
    var sc = TwoLocks(svg);
    var banner = T(svg, 570, H - 8, '', { size: 16, weight: 700 });
    var count = { A: T(svg, 1110, 82, '', { anchor: 'end', size: 14, mono: true }), B: T(svg, 1110, 192, '', { anchor: 'end', size: 14, mono: true }) };
    var mode = 'fixed', th, locks;
    sc.code.A.textContent = 'lock x, try y';
    sc.code.B.textContent = 'lock y, try x';
    function reset() {
      locks = { x: null, y: null };
      th = {
        A: { first: 'x', second: 'y', ph: 'take1', left: 0.35, tries: 0, done: 0, got: false },
        B: { first: 'y', second: 'x', ph: 'take1', left: 0.35, tries: 0, done: 0, got: false }
      };
    }
    function rest() { return mode === 'fixed' ? 0.45 : rnd(0.05, 1.0); }
    function release(name) { if (locks.x === name) locks.x = null; if (locks.y === name) locks.y = null; }
    function stepThread(name, t, dt) {
      if (t.ph === 'take1' && !t.got) {               // lock(first): blocks while the other has it
        if (locks[t.first]) return;
        locks[t.first] = name; t.got = true;
      }
      t.left -= dt;
      if (t.left > 0) return;
      switch (t.ph) {
        case 'take1': t.ph = 'try2'; t.left = 0.35; break;
        case 'try2':
          if (!locks[t.second]) { locks[t.second] = name; t.ph = 'both'; t.left = 0.5; }
          else { t.tries++; t.ph = 'back'; t.left = 0.3; }
          break;
        case 'both': release(name); t.done++; t.ph = 'rest'; t.left = rest(); break;
        case 'back': release(name); t.ph = 'rest'; t.left = rest(); break;
        case 'rest': t.ph = 'take1'; t.got = false; t.left = 0.35; break;
      }
    }
    Toggle(ctl, [['fixed', 'back off 450 ms'], ['random', 'back off a random time']], mode, function (v) { mode = v; reset(); });
    Btn(ctl, 'restart', reset);
    reset();
    return {
      tick: function (dt) {
        stepThread('A', th.A, dt);
        stepThread('B', th.B, dt);
        ['A', 'B'].forEach(function (n) {
          var t = th[n];
          sc.line(n, t.first, locks[t.first] === n ? 'hold' : t.ph === 'take1' ? 'want' : '');
          sc.line(n, t.second, locks[t.second] === n ? 'hold' : t.ph === 'try2' ? 'want' : t.ph === 'back' ? 'fail' : '');
          var note = t.ph === 'take1' ? (t.got ? 'locks ' + t.first : 'waits for ' + t.first) : t.ph === 'try2' ? 'tries ' + t.second + '…'
            : t.ph === 'back' ? t.second + ' taken ✗ · lets go' : t.ph === 'both' ? 'has both ✓ · works' : 'backs off';
          sc.toks[n].state(t.ph === 'back' ? 'bad' : t.ph === 'both' ? 'good' : 'run').note(note, t.ph === 'back' ? C.bad : C.dim);
          count[n].textContent = t.tries + ' tries · ' + t.done + ' done';
        });
        sc.holder('x', locks.x); sc.holder('y', locks.y);
        var done = th.A.done + th.B.done;
        banner.textContent = done ? done + ' transfers done — the random waits broke the lockstep'
          : 'busy all the time · ' + (th.A.tries + th.B.tries) + ' tries · 0 done';
        banner.setAttribute('fill', done ? C.good : C.bad);
      }
    };
  });

  /* Lock striping: N locks chosen by hash. */
  def('striping', function (host) {
    var W = 1140, H = 280, svg = mkSvg(host, W, H), ctl = Controls(host);
    var NT = 8, OP = 0.5, N = 1, t = 0;
    var layer = G(svg), toksG = G(svg);
    var info = T(svg, 20, H - 8, '', { anchor: 'start', size: 15, mono: true, weight: 600 });
    var waitInfo = T(svg, 1120, H - 8, '', { anchor: 'end', size: 15, mono: true });
    hdr(svg, 570, 18, 'EACH THREAD: PICK A KEY → LOCK hash(key) % N → WORK 0.5 s → UNLOCK');
    var stripes, threads = [], ops = [];
    for (var i = 0; i < NT; i++) {
      threads.push({ tok: Token(toksG, 'T' + (i + 1), PALETTE[i], { w: 40, h: 26, size: 12 }), st: 'idle', left: rnd(0.05, 0.4), s: -1, home: 90 + i * 135 });
    }
    function build() {
      clear(layer);
      stripes = [];
      var colW = 1080 / N;
      for (var s = 0; s < N; s++) {
        var cx = 30 + colW * (s + 0.5), bw = Math.min(colW - 8, 110);
        R(layer, cx - bw / 2, 52, bw, 46, { fill: C.panel, stroke: C.line, sw: 1.4 });
        T(layer, cx, 47, N > 8 ? String(s) : 'lock ' + s, { size: 12, mono: true, fill: C.dim });
        stripes.push({ cx: cx, holder: null, queue: [] });
      }
      threads.forEach(function (th) { th.st = 'idle'; th.left = rnd(0.05, 0.4); th.s = -1; });
      ops = [];
    }
    function layout() {
      threads.forEach(function (th) {
        if (th.st === 'idle') th.tok.at(th.home, 224).state('run');
      });
      stripes.forEach(function (s) {
        if (s.holder) s.holder.tok.at(s.cx, 75).state('good');
        s.queue.forEach(function (th, i) { th.tok.at(s.cx, 120 + Math.min(i, 3) * 28).state('wait'); });
      });
    }
    function take(s, th) { s.holder = th; th.st = 'hold'; th.left = OP; }
    var opts = [1, 2, 4, 8, 16];
    Slider(ctl, 'locks (stripes)', 0, 4, 1, 0, function (v) { N = opts[v]; build(); }, function (v) { return String(opts[v]); });
    Btn(ctl, 'reset', build);
    build();
    return {
      tick: function (dt) {
        t += dt;
        threads.forEach(function (th) {
          th.left -= dt;
          if (th.left > 0) return;
          if (th.st === 'idle') {
            var s = stripes[Math.floor(Math.random() * 64) % N];
            th.s = stripes.indexOf(s);
            if (!s.holder) take(s, th); else { th.st = 'wait'; th.left = Infinity; s.queue.push(th); }
          } else if (th.st === 'hold') {
            var st = stripes[th.s];
            st.holder = null;
            ops.push(t);
            th.st = 'idle'; th.left = rnd(0.1, 0.35);
            if (st.queue.length) take(st, st.queue.shift());
          }
        });
        layout();
        while (ops.length && ops[0] < t - 3) ops.shift();
        var waiting = threads.filter(function (th) { return th.st === 'wait'; }).length;
        info.textContent = 'throughput: ' + (ops.length / Math.min(3, Math.max(t, 0.5))).toFixed(1) + ' operations/s';
        waitInfo.textContent = 'waiting for a lock right now: ' + waiting + ' of ' + NT;
      }
    };
  });

  /* Amdahl's law as a live chart. */
  def('amdahl', function (host) {
    var W = 1140, H = 250, svg = mkSvg(host, W, H), ctl = Controls(host);
    var s = 0.1, n = 16, X0 = 90, X1 = 1000, Y0 = 220, Y1 = 16;
    var g = G(svg);
    function sp(c) { return 1 / (s + (1 - s) / c); }
    function draw() {
      clear(g);
      var ymax = s > 0 ? Math.min(64, Math.max(8, Math.ceil(1.25 / s / 4) * 4)) : 64;
      var px = function (c) { return X0 + (c - 1) / 63 * (X1 - X0); };
      var py = function (v) { return Y0 - Math.min(v, ymax) / ymax * (Y0 - Y1); };
      for (var i = 0; i <= 4; i++) {
        var v = ymax * i / 4, y = py(v);
        L(g, X0, y, X1, y, { stroke: C.faint, sw: 1 });
        T(g, X0 - 10, y + 5, v + '×', { anchor: 'end', size: 12, mono: true, fill: C.dim });
      }
      [1, 8, 16, 32, 48, 64].forEach(function (c) { T(g, px(c), Y0 + 20, String(c), { size: 12, mono: true, fill: C.dim }); });
      T(g, X1 + 14, Y0 + 20, 'cores', { anchor: 'start', size: 12, fill: C.dim });
      var top = Math.min(64, ymax);
      L(g, px(1), py(1), px(top), py(top), { stroke: C.dim, sw: 1.4, dash: '5 5' });
      T(g, px(top) + 10, py(top) + 4, 'perfect', { anchor: 'start', size: 12, fill: C.dim });
      if (s > 0 && 1 / s <= ymax) {
        L(g, X0, py(1 / s), X1, py(1 / s), { stroke: C.bad, sw: 1.6, dash: '7 5' });
        T(g, X1, py(1 / s) - 8, 'never more than ' + (1 / s).toFixed(1).replace(/\.0$/, '') + '×', { anchor: 'end', size: 14, weight: 700, fill: C.bad });
      }
      var d = '';
      for (var c = 1; c <= 64; c++) d += (c === 1 ? 'M' : 'L') + px(c).toFixed(1) + ',' + py(sp(c)).toFixed(1);
      Pth(g, d, { stroke: C.accent, sw: 3 });
      S('circle', { cx: px(n), cy: py(sp(n)), r: 7, fill: C.accent, stroke: C.bg, 'stroke-width': 2 }, g);
      T(g, px(n) + 12, py(sp(n)) + 24, n + (n === 1 ? ' core' : ' cores') + ' → ' + sp(n).toFixed(1) + '×', { anchor: 'start', size: 15, mono: true, weight: 700 });
    }
    Slider(ctl, 'share of the work under one lock', 0, 50, 1, 10, function (v) { s = v / 100; draw(); }, function (v) { return v + '%'; });
    Slider(ctl, 'cores', 1, 64, 1, 16, function (v) { n = v; draw(); });
    draw();
    return {};
  });

  /* A thread pool: submit → bounded queue → N workers. */
  def('pool', function (host) {
    var W = 1140, H = 272, svg = mkSvg(host, W, H), ctl = Controls(host);
    var NW = 3, rate = 2.5, CAP = 10, midY = 120, QX = 270, SLW = 44;
    var t = 0, nextArr = 0.5, pending = 0, seq = 0, queue = [], workers = [], lat = [], done = 0;
    hdr(svg, 115, 18, 'submit()');
    hdr(svg, QX + (CAP * (SLW + 6)) / 2, 18, 'TASK QUEUE · CAPACITY ' + CAP);
    hdr(svg, 960, 18, 'WORKER THREADS');
    var sub = R(svg, 40, midY - 30, 150, 60, { fill: C.panel, stroke: C.line, sw: 1.6 });
    var subT = T(svg, 115, midY + 5, '', { size: 13 });
    var subN = T(svg, 115, midY + 52, '', { size: 13, fill: C.bad, weight: 600 });
    for (var i = 0; i < CAP; i++) R(svg, QX + i * (SLW + 6), midY - 20, SLW, 40, { fill: C.bg, stroke: C.line, dash: '4 3' });
    L(svg, 194, midY, QX - 6, midY, { stroke: C.dim, end: svg.arrow.dim });
    var wLayer = G(svg), tasksG = G(svg);
    var stats = T(svg, 570, H - 8, '', { size: 14, mono: true, fill: C.dim });
    function slotX(i) { return QX + (CAP - 1 - i) * (SLW + 6) + SLW / 2; }
    function buildWorkers() {
      clear(wLayer);
      var old = workers;
      workers = [];
      var spc = Math.min(50, 214 / NW);
      for (var i = 0; i < NW; i++) {
        var y = midY + (i - (NW - 1) / 2) * spc;
        var g = G(wLayer);
        var box = R(g, 830, y - spc / 2 + 4, 270, spc - 8, { fill: C.bg, stroke: C.dim, dash: '5 4' });
        T(g, 845, y + 5, 'W' + (i + 1), { anchor: 'start', size: 13, mono: true, weight: 700 });
        var b = bar(g, 950, y - 3, 140, 6, C.good);
        b.el.style.transition = 'none';
        workers.push({ y: y, box: box, b: b, task: null, left: 0, total: 1 });
      }
      old.forEach(function (w) { if (w.task) { queue.unshift(w.task); } });
      layout();
    }
    function mkTask() {
      seq++;
      var tok = Token(tasksG, String(seq), PALETTE[seq % PALETTE.length], { w: 38, h: 26, size: 12, mono: true });
      tok.jump(115, midY);
      return { tok: tok, t0: t, dur: rnd(0.8, 1.6) };
    }
    function layout() {
      queue.forEach(function (tk, i) { tk.tok.at(slotX(i), midY).state('wait'); });
      workers.forEach(function (w) {
        if (w.task) w.task.tok.at(900, w.y).state('run');
        w.box.setAttribute('stroke', w.task ? C.good : C.dim);
        w.box.setAttribute('stroke-dasharray', w.task ? 'none' : '5 4');
        w.box.setAttribute('fill', w.task ? C.goodFill : C.bg);
      });
    }
    function reset() {
      clear(tasksG); queue = []; pending = 0; lat = []; done = 0; t = 0;
      workers.forEach(function (w) { w.task = null; });
      buildWorkers();
    }
    Slider(ctl, 'workers', 1, 8, 1, NW, function (v) { NW = v; buildWorkers(); });
    Slider(ctl, 'arriving', 0.5, 8, 0.5, rate, function (v) { rate = v; }, function (v) { return v + '/s'; });
    Btn(ctl, '＋ 20 tasks at once', function () { pending += 20; });
    Btn(ctl, 'reset', reset);
    buildWorkers();
    return {
      tick: function (dt) {
        t += dt;
        nextArr -= dt;
        while (nextArr <= 0) { pending++; nextArr += -Math.log(1 - Math.random()) / rate; }
        while (pending > 0 && queue.length < CAP) { queue.push(mkTask()); pending--; }
        workers.forEach(function (w) {
          if (w.task) {
            w.left -= dt;
            w.b.set(1 - w.left / w.total);
            if (w.left <= 0) {
              var tk = w.task;
              w.task = null; done++;
              tk.tok.at(1130, w.y);
              tk.tok.show(0);
              setTimeout(function () { tk.tok.remove(); }, 700);
            }
          }
          if (!w.task && queue.length) {
            var next = queue.shift();
            lat.push(t - next.t0); if (lat.length > 20) lat.shift();
            w.task = next; w.left = w.total = next.dur;
          }
          if (!w.task) w.b.set(0);
        });
        layout();
        var blocked = pending > 0;
        sub.setAttribute('stroke', blocked ? C.bad : C.line);
        sub.setAttribute('stroke-dasharray', blocked ? '5 4' : 'none');
        subT.textContent = blocked ? 'waits: queue full' : 'submitting';
        subT.setAttribute('fill', blocked ? C.bad : C.fg);
        subN.textContent = blocked ? pending + ' more to submit' : '';
        var avg = lat.length ? lat.reduce(function (a, b) { return a + b; }, 0) / lat.length : 0;
        stats.textContent = 'done: ' + done + ' · a task waits in the queue ' + avg.toFixed(1) + ' s on average';
      }
    };
  });

  /* The inside of a future, line by line. */
  def('futureTrace', function (host) {
    var code = [
      'template <typename T> class Shared {  // both ends point here',
      '    std::mutex m;',
      '    std::condition_variable cv;',
      '    std::optional<T> value;',
      'public:',
      '    void set(T v) {                    // promise::set_value',
      '        { std::lock_guard lk(m); value = std::move(v); }',
      '        cv.notify_all();',
      '    }',
      '    T get() {                          // future::get',
      '        std::unique_lock lk(m);',
      '        cv.wait(lk, [&] { return value.has_value(); });',
      '        return *value;',
      '    }',
      '};'
    ];
    var W = 1140, LH = 22, top = 10, CX = 92, H = top + code.length * LH + 12, svg = mkSvg(host, W, H);
    R(svg, 0, 0, 750, H, { fill: 'rgba(8,24,40,.6)', stroke: C.faint, sw: 1 });
    var hlC = R(svg, 4, 0, 742, LH, { fill: 'rgba(92,200,234,.18)', stroke: C.accent, sw: 1, rx: 2 });
    var hlW = R(svg, 4, 0, 742, LH, { fill: 'rgba(126,208,163,.18)', stroke: C.good, sw: 1, rx: 2 });
    var tagC = T(svg, 12, 0, 'caller ▸', { anchor: 'start', size: 12, weight: 700, fill: C.accent });
    var tagW = T(svg, 12, 0, 'worker ▸', { anchor: 'start', size: 12, weight: 700, fill: C.good });
    [hlC, hlW, tagC, tagW].forEach(function (e) { e.classList.add('col'); e.style.transition = 'opacity .3s, transform .35s ease'; });
    var KW = /\b(template|typename|class|public|void|return)\b/;
    code.forEach(function (line, i) {
      var y = top + i * LH + 16;
      var ci = line.indexOf('//'), src = ci >= 0 ? line.slice(0, ci) : line, cm = ci >= 0 ? line.slice(ci) : '';
      var tx = S('text', { x: CX, y: y, 'font-family': MONO, 'font-size': 14, fill: C.fg, 'xml:space': 'preserve' }, svg);
      tx.style.whiteSpace = 'pre';
      src.split(/(\b(?:template|typename|class|public|void|return)\b)/).forEach(function (part) {
        if (!part) return;
        var sp = S('tspan', { fill: KW.test(part) ? '#8fd3f0' : C.fg }, tx);
        sp.textContent = part;
      });
      if (cm) { var c2 = S('tspan', { fill: C.dim }, tx); c2.textContent = cm; }
    });
    var panel = Panel(svg, 780, 0, 360, 170, 'STATE', ['lock m', 'value', 'waiting room', 'caller']);
    var STEPS = [
      [-1, -1, ['free', 'empty', 'empty', '—']],
      [10, -1, ['caller', 'empty', 'empty', 'in get()']],
      [11, -1, ['free', 'empty', 'caller', 'asleep']],
      [11, 6, ['worker', '42', 'caller', 'asleep']],
      [11, 7, ['free', '42', 'empty', 'woken · wants m']],
      [11, -1, ['caller', '42', 'empty', 'checks again: yes']],
      [12, -1, ['caller', '42', 'empty', ['returns 42', C.good]]]
    ];
    function put(hl, tag, line) {
      hl.style.opacity = line < 0 ? 0 : 1;
      tag.style.opacity = line < 0 ? 0 : 1;
      if (line < 0) return;
      hl.style.transform = 'translateY(' + (top + line * LH + 1) + 'px)';
      tag.style.transform = 'translateY(' + (top + line * LH + 16) + 'px)';
    }
    return {
      steps: 6,
      render: function (k) {
        var s = STEPS[k];
        put(hlC, tagC, s[0]);
        put(hlW, tagW, s[1]);
        var keys = ['lock m', 'value', 'waiting room', 'caller'];
        s[2].forEach(function (v, i) { if (Array.isArray(v)) panel.set(keys[i], v[0], v[1]); else panel.set(keys[i], v); });
      }
    };
  });

  /* Single-flight: one disk read, everyone else waits on the same future. */
  def('singleFlight', function (host) {
    var W = 1140, H = 270, svg = mkSvg(host, W, H), q = Sched();
    hdr(svg, 100, 18, 'THREADS');
    hdr(svg, 520, 18, 'files · map<path, shared_future>');
    R(svg, 340, 74, 360, 120, { fill: 'rgba(20,56,90,.5)', stroke: C.line });
    var entry = T(svg, 520, 140, '', { size: 16, mono: true, weight: 600 });
    entry.classList.add('fd');
    hdr(svg, 985, 18, 'DISK');
    var disk = R(svg, 880, 34, 210, 96, { fill: C.panel, stroke: C.line, sw: 1.6 });
    var busy = T(svg, 985, 76, '', { size: 14, fill: C.warn, weight: 600 });
    var reads = T(svg, 985, 112, '', { size: 15, mono: true });
    var summary = T(svg, 985, 200, '', { size: 18, weight: 700, fill: C.good });
    var T1 = Token(svg, 'T1', C.accent, { w: 60 }), T2 = Token(svg, 'T2', C.violet, { w: 60 }), T3 = Token(svg, 'T3', C.warn, { w: 60 });
    return {
      steps: 6,
      render: function (k, prev) {
        q.clear();
        var fwd = prev === k - 1;
        T1.at(100, 60).state('run').note('');
        T2.at(100, 140).state('run').note('');
        T3.at(100, 220).state('run').note('');
        fade(entry, 0); entry.textContent = '“cat.png” → ⏳ future';
        entry.setAttribute('fill', C.fg);
        busy.textContent = ''; busy.classList.remove('pulse');
        disk.setAttribute('stroke', C.line);
        reads.textContent = 'reads: 0'; summary.textContent = '';
        if (k >= 1) { T1.at(280, 60).note("miss: I'll read it"); fade(entry, 1); }
        if (k >= 2) {
          T1.at(780, 70).note('reads · no lock held');
          busy.textContent = 'reading cat.png…'; busy.classList.add('pulse');
          disk.setAttribute('stroke', C.warn); reads.textContent = 'reads: 1';
        }
        if (k >= 3) T2.at(280, 140).state('sleep').note('f.get(): waits');
        if (k >= 4) T3.at(280, 220).state('sleep').note('f.get(): waits');
        if (k >= 5) {
          busy.textContent = ''; busy.classList.remove('pulse'); disk.setAttribute('stroke', C.line);
          entry.textContent = '“cat.png” → ✓ image'; entry.setAttribute('fill', C.good);
          T1.at(280, 60).state('good').note('set_value(image)');
          q.seq(fwd, [500, function () {
            T2.state('good').note('same image');
            T3.state('good').note('same image');
          }]);
        }
        if (k >= 6) summary.textContent = '3 requests · 1 read';
      }
    };
  });

  /* Fan-out / fan-in over 4 chunks. */
  def('fanOut', function (host) {
    var W = 1140, H = 290, svg = mkSvg(host, W, H);
    var CW = 14.5, X0 = 106, cols = [C.accent, C.good, C.violet, C.warn];
    var data = [], partial = [0, 0, 0, 0];
    for (var i = 0; i < 64; i++) { data.push((i * 37) % 19 + 1); partial[Math.floor(i / 16)] += data[i] * data[i]; }
    var total = partial.reduce(function (a, b) { return a + b; }, 0);
    hdr(svg, X0, 18, 'THE DATA · 64 NUMBERS', 'start');
    var cells = data.map(function (v, i) {
      var r = R(svg, X0 + i * CW, 28, CW - 2, 26, { fill: C.panel, stroke: 'none', rx: 1 });
      r.classList.add('col');
      return r;
    });
    var cx = [0, 1, 2, 3].map(function (j) { return X0 + (16 * j + 8) * CW - 1; });
    var g2 = G(svg, 'fd'), g3 = G(svg, 'fd'), g4 = G(svg, 'fd');
    var bars = [];
    cx.forEach(function (x, j) {
      L(g2, x, 58, x, 104, { stroke: cols[j], end: svg.arrow[['accent', 'good', 'violet', 'warn'][j]] });
      R(g2, x - 92, 108, 184, 56, { fill: FILL[cols[j]], stroke: cols[j], sw: 1.6 });
      T(g2, x, 130, 'task ' + (j + 1) + ' · sums its chunk', { size: 13 });
      bars.push(bar(g2, x - 80, 144, 160, 8, cols[j]));
      T(g3, x, 194, 'future ✓ ' + partial[j].toLocaleString('en-US'), { size: 15, mono: true, weight: 600, fill: cols[j] });
      L(g4, x, 202, 570 + (j - 1.5) * 40, 236, { stroke: C.dim, end: svg.arrow.dim });
    });
    R(g4, 420, 240, 300, 40, { fill: C.goodFill, stroke: C.good, sw: 1.8 });
    T(g4, 570, 266, 'total = ' + total.toLocaleString('en-US'), { size: 17, mono: true, weight: 700 });
    return {
      steps: 5,
      render: function (k) {
        cells.forEach(function (c, i) { c.setAttribute('fill', k >= 1 ? FILL[cols[Math.floor(i / 16)]] : C.panel); c.setAttribute('stroke', k >= 1 ? cols[Math.floor(i / 16)] : 'none'); });
        fade(g2, k >= 2); fade(g3, k >= 3); fade(g4, k >= 4);
        bars.forEach(function (b) { b.set(k >= 2 ? 1 : 0); });
      }
    };
  });

  /* A task waiting for a subtask in the same full pool. */
  def('poolDeadlock', function (host) {
    var W = 1140, H = 250, svg = mkSvg(host, W, H);
    hdr(svg, 260, 18, 'TASK QUEUE');
    R(svg, 80, 76, 360, 100, { fill: C.bg, stroke: C.line, dash: '5 4' });
    hdr(svg, 850, 18, 'POOL · 2 WORKER THREADS');
    R(svg, 640, 34, 420, 76, { fill: C.bg, stroke: C.line });
    R(svg, 640, 142, 420, 76, { fill: C.bg, stroke: C.line });
    T(svg, 656, 76, 'worker 1', { anchor: 'start', size: 13, mono: true, fill: C.dim });
    T(svg, 656, 184, 'worker 2', { anchor: 'start', size: 13, mono: true, fill: C.dim });
    var hot = G(svg, 'fd');
    Pth(hot, 'M756,52 C640,36 300,40 214,108', { stroke: C.bad, sw: 2, dash: '6 5', end: svg.arrow.bad });
    T(hot, 520, 34, 'A waits for a', { size: 13, fill: C.bad });
    Pth(hot, 'M756,196 C640,238 380,232 316,146', { stroke: C.bad, sw: 2, dash: '6 5', end: svg.arrow.bad });
    T(hot, 560, 244, 'B waits for b', { size: 13, fill: C.bad });
    L(hot, 444, 126, 632, 126, { stroke: C.bad, sw: 2, dash: '6 5', end: svg.arrow.bad });
    T(hot, 538, 116, 'need a free worker', { size: 13, fill: C.bad });
    var A = Token(svg, 'A', C.accent, { w: 60 }), B = Token(svg, 'B', C.good, { w: 60 });
    var a = Token(svg, 'a', C.accent, { w: 44, h: 28 }), b = Token(svg, 'b', C.good, { w: 44, h: 28 });
    return {
      steps: 4,
      render: function (k) {
        A.at(800, 62).state('run').note('').show(k >= 1);
        B.at(800, 170).state('run').note('').show(k >= 1);
        a.at(800, 62).state('wait').show(0);
        b.at(800, 170).state('wait').show(0);
        fade(hot, 0);
        if (k >= 1) { A.note('running'); B.note('running'); }
        if (k >= 2) { a.at(200, 126).show(1); b.at(310, 126).show(1); A.note('submit(a)'); B.note('submit(b)'); }
        if (k >= 3) { A.state('sleep').note('a.get()'); B.state('sleep').note('b.get()'); }
        if (k >= 4) { fade(hot, 1); a.state('bad'); b.state('bad'); }
      }
    };
  });

  /* A three-stage pipeline joined by bounded queues. */
  def('pipeline', function (host) {
    var W = 1140, H = 196, svg = mkSvg(host, W, H), ctl = Controls(host);
    var CAP = 4, midY = 104, SW = 50;
    var stages = [
      { name: 'read', x: 140, rate: 3 }, { name: 'resize', x: 570, rate: 1.5 }, { name: 'write', x: 1000, rate: 2.5 }
    ];
    var qx = [252, 682];
    var t = 0, seq = 0, done = [], itemsG;
    stages.forEach(function (s) {
      s.box = R(svg, s.x - 80, midY - 40, 160, 80, { fill: C.panel, stroke: C.line, sw: 1.6 });
      T(svg, s.x - 18, midY - 10, s.name, { size: 16, weight: 700 });
      s.st = T(svg, s.x, midY + 20, '', { size: 12, fill: C.dim });
      s.b = bar(svg, s.x - 66, midY + 30, 132, 4, C.good);
      s.b.el.style.transition = 'none';
      s.tag = T(svg, s.x, midY - 52, '', { size: 13, weight: 700, fill: C.warn });
    });
    qx.forEach(function (x) {
      for (var i = 0; i < CAP; i++) R(svg, x + i * (SW + 6), midY - 20, SW, 40, { fill: C.bg, stroke: C.line, dash: '4 3' });
    });
    hdr(svg, qx[0] + 109, midY - 32, 'QUEUE · 4');
    hdr(svg, qx[1] + 109, midY - 32, 'QUEUE · 4');
    itemsG = G(svg);
    var info = T(svg, 570, H - 10, '', { size: 15, mono: true, weight: 600 });
    var queues = [[], []];
    function qX(qi, i) { return qx[qi] + (CAP - 1 - i) * (SW + 6) + SW / 2; }
    function mkItem() {
      seq++;
      var g = G(itemsG, 'mv');
      R(g, -13, -13, 26, 26, { fill: C.warnFill, stroke: C.warn, sw: 1.3, rx: 3 });
      T(g, 0, 5, String(seq % 100), { size: 11, mono: true, weight: 600 });
      return g;
    }
    function place(g, x, y, instant) {
      if (instant) g.style.transition = 'none';
      g.style.transform = 'translate(' + x + 'px,' + y + 'px)';
      if (instant) { void g.getBoundingClientRect(); g.style.transition = ''; }
    }
    function reset() {
      clear(itemsG); queues = [[], []]; done = []; t = 0;
      stages.forEach(function (s) { s.item = null; s.state = 'idle'; s.left = 0; s.total = 1; });
    }
    function dur(s) { return rnd(0.8, 1.2) / s.rate; }
    stages.forEach(function (s) {
      Slider(ctl, s.name, 0.5, 4, 0.5, s.rate, function (v) { s.rate = v; }, function (v) { return v + '/s'; });
    });
    Btn(ctl, 'reset', reset);
    reset();
    return {
      tick: function (dt) {
        t += dt;
        for (var i = stages.length - 1; i >= 0; i--) {
          var s = stages[i];
          if (s.state === 'work') { s.left -= dt; if (s.left <= 0) s.state = 'out'; }
          if (s.state === 'out') {
            if (i === stages.length - 1) {
              var g = s.item;
              place(g, 1125, midY); g.style.opacity = 0;
              setTimeout(function () { if (g.parentNode) g.parentNode.removeChild(g); }, 600);
              done.push(t); s.item = null; s.state = 'idle';
            } else if (queues[i].length < CAP) { queues[i].push(s.item); s.item = null; s.state = 'idle'; }
            else s.state = 'blocked';
          } else if (s.state === 'blocked' && queues[i].length < CAP) { queues[i].push(s.item); s.item = null; s.state = 'idle'; }
          if (s.state === 'idle') {
            if (i === 0) { s.item = mkItem(); place(s.item, s.x - 110, midY, true); }
            else if (queues[i - 1].length) s.item = queues[i - 1].shift();
            if (s.item) { s.state = 'work'; s.left = s.total = dur(s); }
          }
        }
        queues.forEach(function (qq, qi) { qq.forEach(function (g, j) { place(g, qX(qi, j), midY); }); });
        var slow = stages.reduce(function (m, s) { return s.rate < m.rate ? s : m; }, stages[0]);
        stages.forEach(function (s, i) {
          if (s.item && s.state !== 'idle') place(s.item, s.x + 48, midY - 12);
          s.b.set(s.state === 'work' ? 1 - s.left / s.total : s.state === 'blocked' ? 1 : 0);
          var label = s.state === 'work' ? 'working' : s.state === 'blocked' ? 'output full: waits' : 'input empty: waits';
          s.st.textContent = label;
          s.box.setAttribute('stroke', s.state === 'work' ? C.good : s.state === 'blocked' ? C.bad : C.dim);
          s.box.setAttribute('stroke-dasharray', s.state === 'work' ? 'none' : '5 4');
          s.tag.textContent = s === slow ? 'slowest stage' : '';
        });
        while (done.length && done[0] < t - 4) done.shift();
        info.textContent = 'out of the pipeline: ' + (done.length / Math.min(4, Math.max(t, 0.5))).toFixed(1) + ' items/s';
      }
    };
  });

  /* An event loop: one thread, a ready queue, tasks waiting on I/O. */
  def('eventLoop', function (host) {
    var W = 1140, H = 300, svg = mkSvg(host, W, H), ctl = Controls(host);
    var IO = ['socket read', 'HTTP call', 'DB query', 'timer', 'file read'];
    hdr(svg, 190, 18, 'READY · waiting for the thread');
    R(svg, 40, 64, 310, 76, { fill: C.bg, stroke: C.line, dash: '5 4' });
    hdr(svg, 535, 18, 'THE ONE THREAD');
    var box = R(svg, 385, 40, 300, 124, { fill: C.panel, stroke: C.accent, sw: 2 });
    var doing = T(svg, 535, 150, '', { size: 13, fill: C.dim });
    hdr(svg, 925, 18, 'WAITING ON I/O · not using the thread');
    R(svg, 730, 30, 380, 210, { fill: 'rgba(20,56,90,.25)', stroke: C.line, dash: '5 4' });
    var os = R(svg, 730, 250, 380, 42, { fill: C.panel, stroke: C.line });
    var osT = T(svg, 920, 276, 'OS · epoll / kqueue / IOCP: "this I/O is ready"', { size: 13, fill: C.dim });
    Pth(svg, 'M726,271 L18,271 L18,102 L34,102', { stroke: C.dim, dash: '4 4', end: svg.arrow.dim });
    var stats = T(svg, 44, 178, '', { anchor: 'start', size: 14, mono: true });
    var stats2 = T(svg, 44, 202, '', { anchor: 'start', size: 14, mono: true });
    var tasksG = G(svg);
    var tasks, ready, running, blockNext = false, blockLeft = 0, paused = false, slices = 0;
    function reset() {
      clear(tasksG);
      tasks = ['A', 'B', 'C', 'D', 'E'].map(function (n, i) {
        var tok = Token(tasksG, n, PALETTE[i], { w: 44, h: 30 });
        tok.jump(70 + i * 54, 102);
        return { name: n, tok: tok, st: 'ready', left: 0, io: '', since: 0, slot: -1 };
      });
      ready = tasks.slice();
      running = null; blockLeft = 0; blockNext = false; slices = 0;
    }
    function ioSlot() {
      for (var i = 0; i < 6; i++) if (!tasks.some(function (t) { return t.st === 'io' && t.slot === i; })) return i;
      return 0;
    }
    var runBtn = Btn(ctl, '⏸ pause', function () { paused = !paused; runBtn.textContent = paused ? '▶ run' : '⏸ pause'; });
    var blockBtn = Btn(ctl, '⏳ block the loop for 3 s', function () { blockNext = true; blockBtn.classList.add('on'); });
    Btn(ctl, 'reset', reset);
    reset();
    return {
      tick: function (dt) {
        if (paused) return;
        var osHit = false;
        tasks.forEach(function (t) {
          if (t.st === 'ready') t.since += dt;
          if (t.st !== 'io') return;
          t.left -= dt;
          if (t.left <= 0) { t.st = 'ready'; t.since = 0; ready.push(t); osHit = true; }
        });
        if (osHit) { flash(os); flash(osT); }
        if (running) {
          if (blockLeft > 0) {
            blockLeft -= dt;
            if (blockLeft <= 0) { running.left = 0; blockBtn.classList.remove('on'); }
          } else running.left -= dt;
          if (blockLeft <= 0 && running.left <= 0) {
            running.st = 'io'; running.left = rnd(1.5, 3.5);
            running.io = IO[Math.floor(Math.random() * IO.length)];
            running.slot = ioSlot();
            running = null; slices++;
          }
        }
        if (!running && ready.length) {
          running = ready.shift();
          running.st = 'run'; running.left = rnd(0.5, 0.9);
          if (blockNext) { blockNext = false; blockLeft = 3; }
        }
        ready.forEach(function (t, i) { t.tok.at(318 - Math.min(i, 5) * 54, 102).state('wait').note(i < 6 ? t.since.toFixed(1) + ' s' : '', t.since > 1.5 ? C.bad : C.dim); });
        tasks.forEach(function (t) {
          if (t.st === 'io') {
            t.tok.at(830 + (t.slot % 2) * 180, 64 + Math.floor(t.slot / 2) * 62).state('sleep').note('await ' + t.io);
          }
        });
        if (running) running.tok.at(535, 92).state(blockLeft > 0 ? 'bad' : 'run').note('');
        box.setAttribute('stroke', blockLeft > 0 ? C.bad : C.accent);
        doing.textContent = !running ? 'idle: nothing is ready' : blockLeft > 0
          ? 'time.sleep(3): nobody else can run' : 'runs ' + running.name + ' until its next await';
        doing.setAttribute('fill', blockLeft > 0 ? C.bad : C.dim);
        var worst = ready.reduce(function (m, t) { return Math.max(m, t.since); }, 0);
        stats.textContent = 'slices run: ' + slices;
        stats2.textContent = 'longest wait for the thread: ' + worst.toFixed(1) + ' s';
        stats2.setAttribute('fill', worst > 1.5 ? C.bad : C.fg);
      }
    };
  });

  def('awaitRace', function (host) {
    return Timeline(host, {
      lanes: [{ name: 'TASK A', color: C.accent }, { name: 'TASK B', color: C.violet }],
      rowH: 34,
      rows: [
        { lane: 0, text: 'balance >= 100? yes' },
        { lane: 0, text: 'await audit_log() · suspended' },
        { lane: 1, text: 'balance >= 100? yes' },
        { lane: 1, text: 'await audit_log() · suspended' },
        { lane: 0, text: 'resumes · balance -= 100' },
        { lane: 1, text: 'resumes · balance -= 100', cls: 'bad' }
      ],
      stateTitle: 'ONE THREAD',
      keys: ['balance', 'running now'],
      state: function (k) {
        return {
          balance: k >= 6 ? ['−100', C.bad] : k >= 5 ? '0' : '100',
          'running now': ['—', 'task A', 'task A', 'task B', 'task B', 'task A', 'task B'][k]
        };
      }
    });
  });


  /* The condition-variable picture: the room guarded by the mutex, the
     buffer inside, and the waiting room on its side. */
  def('cvRoom', function (host) {
    var W = 1140, H = 462, svg = mkSvg(host, W, H);
    var A = { stroke: C.accent, sw: 2.2, end: svg.arrow.accent };
    var WALL = { stroke: C.fg, sw: 2.6 };
    function tick(x, y, vertical) {
      if (vertical) L(svg, x - 6, y, x + 6, y, { stroke: C.accent, sw: 2.2 });
      else L(svg, x, y - 6, x, y + 6, { stroke: C.accent, sw: 2.2 });
    }
    function bigT(x, y) { T(svg, x, y, 'T', { size: 32, weight: 600 }); }
    function label(x, y, str, anchor, col) { T(svg, x, y, str, { size: 17, mono: true, anchor: anchor || 'start', fill: col || C.fg }); }

    // threads arriving
    bigT(578, 52); bigT(612, 30); bigT(640, 62);
    label(680, 46, 'threads', 'start', C.dim);
    // the mutex: everything inside the dotted line
    R(svg, 410, 100, 690, 300, { stroke: C.bad, sw: 1.6, dash: '2 4', rx: 0 });
    label(1094, 90, 'mutex', 'end', C.dim);
    // mutex.lock(): in through the top door
    tick(610, 82, true); L(svg, 610, 82, 610, 152, A);
    label(596, 117, 'mutex.lock()', 'end', C.accent);
    // the room
    L(svg, 450, 130, 560, 130, WALL); L(svg, 660, 130, 1060, 130, WALL);
    L(svg, 1060, 130, 1060, 380, WALL);
    L(svg, 450, 380, 600, 380, WALL); L(svg, 700, 380, 1060, 380, WALL);
    L(svg, 450, 130, 450, 182, WALL); L(svg, 450, 330, 450, 380, WALL);
    // the waiting room on the side
    L(svg, 450, 182, 140, 182, WALL); L(svg, 140, 182, 140, 330, WALL); L(svg, 140, 330, 450, 330, WALL);
    label(295, 170, 'waiting room', 'middle', C.dim);
    L(svg, 458, 228, 458, 284, WALL);                         // the sleeping thread can't get through
    bigT(282, 226); bigT(228, 268); bigT(282, 310);
    L(svg, 250, 257, 458, 257, { stroke: C.fg, sw: 2 });
    // notEmpty.wait(mutex): from the room into the waiting room
    tick(580, 215); L(svg, 580, 215, 312, 215, A);
    label(596, 221, 'notEmpty.wait(mutex)', 'start', C.accent);
    // a woken thread comes back in
    tick(312, 299); L(svg, 312, 299, 580, 299, A);
    // notEmpty.notify_one(): wakes one sleeper
    label(476, 352, 'notEmpty.notify_one()', 'start', C.accent);
    S('circle', { cx: 560, cy: 360, r: 3.5, fill: C.accent }, svg);
    Pth(svg, 'M560,360 L560,370 L282,370 L282,320', A);
    // the buffer
    R(svg, 720, 236, 240, 62, { stroke: C.fg, sw: 2, rx: 2 });
    [[760, C.bad], [815, '#8fb2e8'], [870, C.good]].forEach(function (b) {
      S('circle', { cx: b[0], cy: 267, r: 21, fill: b[1], stroke: C.bg, 'stroke-width': 2 }, svg);
    });
    label(905, 326, 'buffer', 'middle', C.fg);
    // mutex.unlock(): out through the bottom door
    tick(650, 362, true); L(svg, 650, 362, 650, 452, A);
    label(666, 432, 'mutex.unlock()', 'start', C.accent);
    return {};
  });


  /* One shared task queue vs. a queue per worker (work stealing). */
  def('poolQueues', function (host) {
    var W = 1140, H = 250, svg = mkSvg(host, W, H);
    function task(x, y, col) { R(svg, x - 12, y - 12, 24, 24, { fill: FILL[col], stroke: col, sw: 1.3, rx: 3 }); }
    function worker(x, y, label) {
      R(svg, x - 38, y - 18, 76, 36, { fill: C.panel, stroke: C.line, sw: 1.6, rx: 5 });
      T(svg, x, y + 5, label, { size: 14, weight: 700, mono: true });
    }
    // left: one queue, one lock, everyone fights for it
    hdr(svg, 270, 18, 'ONE SHARED QUEUE');
    R(svg, 130, 48, 280, 50, { fill: C.bg, stroke: C.bad, sw: 1.8 });
    [170, 210, 250, 290, 330, 370].forEach(function (x) { task(x, 73, C.warn); });
    T(svg, 420, 66, '1 lock', { anchor: 'start', size: 13, mono: true, fill: C.bad });
    [70, 200, 340, 470].forEach(function (x, i) {
      worker(x, 196, 'W' + (i + 1));
      L(svg, x, 176, 270 + (x - 270) * 0.25, 104, { stroke: C.bad, sw: 1.6, dash: '5 4', end: svg.arrow.bad });
    });
    T(svg, 270, 244, 'every take waits for the same lock', { size: 13, fill: C.bad });
    L(svg, 570, 30, 570, 222, { stroke: C.faint, sw: 1, dash: '4 6' });
    // right: a queue per worker; an idle worker steals
    hdr(svg, 860, 18, 'ONE QUEUE PER WORKER · WORK STEALING');
    var cols = [660, 780, 900, 1020], counts = [3, 2, 0, 4];
    cols.forEach(function (x, i) {
      R(svg, x - 50, 48, 100, 50, { fill: C.bg, stroke: C.good, sw: 1.6 });
      for (var k = 0; k < counts[i]; k++) task(x - 30 + k * 20, 73, C.warn);
      L(svg, x, 102, x, 172, { stroke: C.good, sw: 1.8, end: svg.arrow.good });
      worker(x, 196, 'W' + (i + 1));
    });
    T(svg, 900, 78, 'empty', { size: 12, fill: C.dim });
    Pth(svg, 'M930,178 C960,140 1010,128 1040,104', { stroke: C.accent, sw: 2.2, dash: '6 4', end: svg.arrow.accent });
    T(svg, 1000, 150, 'steal', { anchor: 'start', size: 14, weight: 700, fill: C.accent });
    return {};
  });

  if (document.readyState === 'loading') document.addEventListener('DOMContentLoaded', init);
  else init();
})();
/* viz:end */
</script>
