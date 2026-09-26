---
marp: true
theme: gaia
class: invert
paginate: true
footer: 'Races, Critical Sections & Mutexes'
style: |
  @import url('https://fonts.bunny.net/css?family=ibm-plex-sans:400,500,600,700|ibm-plex-sans-condensed:600,700|ibm-plex-mono:400,500,600&display=swap');

  /* Blueprint palette: patterns are, in this deck's own words, "a
     blueprint you adapt" — the UML diagrams already read as technical
     drawings, so the whole theme leans into that instead of a generic
     dark-mode-with-neon-accent look. */
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

  /* The mirror of .good: slides where Day 1–2 code breaks. */
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

  .uml { display: block; margin: 0 auto; }

  /* A small aside -- a caveat that shouldn't compete with the slide. */
  .note { font-size: 18px; color: var(--color-dimmed); margin-top: 1.4em; }
  .note code { font-size: 0.9em; }

  /* Step-by-step reveals: list items written with `*` appear one
     keypress at a time in the HTML export. Items that carry a quote
     panel shouldn't also show a list bullet. */
  li:has(> blockquote), li:has(svg) { list-style: none; }
  li > blockquote { margin: 0.2em 0 0.5em; }

  /* Interleaving tables: one row per keypress. Rows sit flush, with no
     list indent, so they line up under the static header row and the
     two thread columns read as continuous timelines. */
  ul:has(> li > svg.tl) { padding-left: 0; }
  li:has(> svg.tl) { margin: 0; line-height: 0; }
---

<!-- _class: lead invert -->

# Races, Critical Sections & Mutexes

### Modeling, Design & Collaborative Development

---

<!-- _class: lead invert -->

# Races

### Same program, same input, different answers

---

## Last time: two threads, one counter

```cpp
int counter = 0;                 // shared by both threads

void work() {
    for (int i = 0; i < 1000000; ++i) counter++;
}

int main() {
    std::thread t1(work);
    std::thread t2(work);
    t1.join();
    t2.join();
    std::cout << counter;        // 2 × 1,000,000 = 2,000,000?
}
```

---

## Run it

* `$ ./counter` → **1,018,264**
* `$ ./counter` → **1,031,176**
* `$ ./counter` → **1,062,627**
* Same program, same input — a different answer every run, and every
  one of them wrong

---

## `counter++` is not one step

```
load   r ← counter      // read the value into a register
add    r ← r + 1        // add one, inside the register
store  counter ← r      // write it back to memory
```

* Another thread can run between any two of these steps
* On one core, the OS can switch threads right there; on two cores,
  both threads really do run at the same moment

---

<!-- _class: invert bad -->

## One lost update

`counter` starts at 5. Both threads run `counter++` once.

<svg class="uml tl" width="900" height="36" viewBox="0 0 760 30"><g font-family="Helvetica, Arial, sans-serif" font-size="12" font-weight="700" letter-spacing="2" fill="#5cc8ea"><text x="200" y="20" text-anchor="middle">THREAD A</text><text x="500" y="20" text-anchor="middle">THREAD B</text><text x="710" y="20" text-anchor="middle">COUNTER</text></g></svg>

* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">1</text><rect x="60" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="200" y="25" text-anchor="middle" fill="#eaf4fb">load: r = 5</text><text x="710" y="25" text-anchor="middle" fill="#eaf4fb" font-weight="400">5</text></g></svg>
* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">2</text><rect x="360" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="500" y="25" text-anchor="middle" fill="#eaf4fb">load: r = 5</text><text x="710" y="25" text-anchor="middle" fill="#eaf4fb" font-weight="400">5</text></g></svg>
* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">3</text><rect x="60" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="200" y="25" text-anchor="middle" fill="#eaf4fb">add: r = 6</text><text x="710" y="25" text-anchor="middle" fill="#eaf4fb" font-weight="400">5</text></g></svg>
* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">4</text><rect x="360" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="500" y="25" text-anchor="middle" fill="#eaf4fb">add: r = 6</text><text x="710" y="25" text-anchor="middle" fill="#eaf4fb" font-weight="400">5</text></g></svg>
* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">5</text><rect x="60" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="200" y="25" text-anchor="middle" fill="#eaf4fb">store: counter = 6</text><text x="710" y="25" text-anchor="middle" fill="#5cc8ea" font-weight="700">6</text></g></svg>
* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">6</text><rect x="360" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="500" y="25" text-anchor="middle" fill="#eaf4fb">store: counter = 6</text><text x="710" y="25" text-anchor="middle" fill="#f0a06e" font-weight="700">6</text></g></svg>
* > Two increments — the counter went up by **one**. An update was lost.

---

## Sequential, concurrent, parallel

* > **Sequential** — one thing at a time, in order.
* > **Concurrent** — several tasks **in progress** during the same
  > period, taking turns, possibly on one core.
* > **Parallel** — several tasks running **at the same instant**,
  > on different cores.
* The counter was concurrent — and, on a laptop with 10 cores,
  parallel too

---

## Two tasks, three ways to run them

* <svg class="uml" width="860" height="356" viewBox="0 0 760 315"><g font-family="Helvetica, Arial, sans-serif" font-size="15"><text x="10" y="38" fill="#eaf4fb" font-weight="700">Sequential</text><text x="10" y="57" fill="#a9c8de" font-size="12">one core</text><rect x="190" y="22" width="275" height="36" fill="#14385a" stroke="#5cc8ea" stroke-width="1.6"/><text x="327" y="45" text-anchor="middle" fill="#eaf4fb" font-weight="700">A</text><rect x="465" y="22" width="275" height="36" fill="#173f37" stroke="#7ed0a3" stroke-width="1.6"/><text x="602" y="45" text-anchor="middle" fill="#eaf4fb" font-weight="700">B</text><text x="10" y="118" fill="#eaf4fb" font-weight="700">Concurrent</text><text x="10" y="137" fill="#a9c8de" font-size="12">one core, taking turns</text><rect x="190" y="102" width="91" height="36" fill="#14385a" stroke="#5cc8ea" stroke-width="1.6"/><text x="235" y="125" text-anchor="middle" fill="#eaf4fb" font-weight="700">A</text><rect x="281" y="102" width="91" height="36" fill="#173f37" stroke="#7ed0a3" stroke-width="1.6"/><text x="326" y="125" text-anchor="middle" fill="#eaf4fb" font-weight="700">B</text><rect x="372" y="102" width="91" height="36" fill="#14385a" stroke="#5cc8ea" stroke-width="1.6"/><text x="417" y="125" text-anchor="middle" fill="#eaf4fb" font-weight="700">A</text><rect x="463" y="102" width="91" height="36" fill="#173f37" stroke="#7ed0a3" stroke-width="1.6"/><text x="508" y="125" text-anchor="middle" fill="#eaf4fb" font-weight="700">B</text><rect x="554" y="102" width="91" height="36" fill="#14385a" stroke="#5cc8ea" stroke-width="1.6"/><text x="599" y="125" text-anchor="middle" fill="#eaf4fb" font-weight="700">A</text><rect x="645" y="102" width="91" height="36" fill="#173f37" stroke="#7ed0a3" stroke-width="1.6"/><text x="690" y="125" text-anchor="middle" fill="#eaf4fb" font-weight="700">B</text><text x="10" y="198" fill="#eaf4fb" font-weight="700">Parallel</text><text x="10" y="217" fill="#a9c8de" font-size="12">two cores, same instant</text><text x="180" y="200" text-anchor="end" fill="#a9c8de" font-size="12">core 1</text><text x="180" y="244" text-anchor="end" fill="#a9c8de" font-size="12">core 2</text><rect x="190" y="180" width="275" height="32" fill="#14385a" stroke="#5cc8ea" stroke-width="1.6"/><text x="327" y="201" text-anchor="middle" fill="#eaf4fb" font-weight="700">A</text><rect x="190" y="224" width="275" height="32" fill="#173f37" stroke="#7ed0a3" stroke-width="1.6"/><text x="327" y="245" text-anchor="middle" fill="#eaf4fb" font-weight="700">B</text><line x1="190" y1="290" x2="734" y2="290" stroke="#6fa8cc" stroke-width="1.6"/><polygon points="740,290 731,285 731,295" fill="#6fa8cc"/><text x="740" y="310" text-anchor="end" fill="#a9c8de" font-size="12">time</text></g></svg>

* Concurrency is how a program is **structured**. Parallelism is how
  it is **executed**. A one-core machine can be concurrent — never
  parallel.

---

## Sequential vs. concurrent — in code

```cpp
// sequential
download(a);
download(b);

// concurrent — both may be in progress at once
std::thread t1(download, a);
std::thread t2(download, b);
t1.join();
t2.join();
```

---

## Question

* Your phone has 8 cores. Your browser has 100 tabs open.
  Concurrent, parallel — or both?
* > **Both.** All 100 tabs are in progress at once — concurrent. At
  > most 8 of them run at the same instant — parallel.

---

## Why anyone bothers

* **Hardware** — single cores stopped getting faster around 2005;
  machines now grow by adding **more cores**
* **Hiding waits** — while one task waits on disk or network,
  another runs
* **Responsiveness** — a UI must react while work happens behind it
* **Throughput** — a server serves thousands of clients at once

---

## The trade threads make

> Threads **share memory** — so talking to each other is free
> and instant... and that is *exactly* the danger.

* The shared counter was the whole point of the two threads — and the
  poison.

---

<span class="tag">Pitfall</span>

## Two things not to believe

* **Concurrent ≠ parallel**
* **Threads are not free** — each one needs its own stack and its own
  OS records; switching between them costs CPU time too

---

## Race condition vs. data race

* > **Race condition** — the result depends on timing: on how the
  > threads' steps happen to interleave.
* > **Data race** — two threads access the same memory at the same
  > time, at least one of them writes, and nothing coordinates them.
* The data race is the defect in the code; the race condition is what
  you see from outside
* In C++, a data race is **undefined behavior** — the compiler is
  allowed to assume it never happens

---

<!-- _class: lead invert -->

# Critical Sections

### Finding them is most of the work

---

## Critical section

* > A **critical section** is code that touches shared state and must
  > not interleave with another thread touching the same state.
* `counter++` was one: three steps that had to act as one
* The skill is **finding** them — the fix comes after the break

---

## The check-then-act shape

```cpp
if (!cache.contains(key))       // check
    cache.insert(key, value);   // act
```

* The check and the act are two separate steps
* Between them, another thread can make the answer to the check
  out of date
* Two threads both check "not there" — and both insert

---

## One shape, many disguises

* **Lazy init** — `if (instance == nullptr) instance = new T();`
* **Files** — `if (!exists(path)) create(path);`
* **Money** — `if (balance >= 100) balance -= 100;`
* **Queues** — `if (!queue.empty()) job = queue.pop();`

---

<span class="tag">Exercise</span>

## Hunt the critical sections

* Four short snippets follow. For each one, mark every critical
  section — or say there is none.
* Assume many threads call each snippet at the same time.

---

<span class="tag">Snippet 1</span>

## An image cache

```cpp
std::unordered_map<std::string, Image> cache;    // shared

Image load(const std::string& path) {
    auto it = cache.find(path);
    if (it != cache.end()) return it->second;
    Image img = readFromDisk(path);
    cache[path] = img;
    return img;
}
```

* From `find` to `cache[path] = img` — check-then-act. Two threads
  both miss, both read the disk, both write — and two writes at once
  can corrupt the map itself.

---

<span class="tag">Snippet 2</span>

## A sum of squares

```cpp
// `data` is filled once at startup and never changed again
long long sumOfSquares(const std::vector<int>& data) {
    long long total = 0;
    for (int x : data) total += x * x;
    return total;
}
```

* **No race.** `total` lives on each thread's own stack, and `data` is
  only read.

---

<span class="tag">Snippet 3</span>

## Order numbers

```cpp
int lastId = 0;                          // shared

std::string newOrderId() {
    lastId = lastId + 1;
    return "ORD-" + std::to_string(lastId);
}
```

* Both lines. The increment is load-add-store, and reading `lastId`
  again can pick up another thread's increment — two orders, one
  number.

---

<span class="tag">Snippet 4</span>

## A job queue

```cpp
std::queue<Job> jobs;                    // shared, filled by another thread

void worker() {
    for (;;) {
        if (!jobs.empty()) {
            Job j = jobs.front();
            jobs.pop();
            run(j);
        }
    }
}
```

* `empty` → `front` → `pop`: two workers see one job and both take it,
  then one pops an empty queue. `run(j)` is fine — `j` is the thread's
  own copy.

---

## Question

* Is a race possible if both threads only **read**?
* > **No.** A race needs at least one **write**.

---

<!-- _class: lead invert -->

# Everything from Days 1–2, Broken

### The day the guarantees died

---

<span class="tag">Day 1</span>

## Singleton — the promised slide

```cpp
static Singleton* get() {
    if (instance == nullptr)
        instance = new Singleton();   // ← remember this line
    return instance;
}
```

* Day 1: *"That `if` on the previous slide is exactly where. We come
  back to it on Day 3."*
* Two threads call `get()` at the same moment.

---

<!-- _class: invert bad -->

## Singleton — two threads, two instances

<svg class="uml tl" width="900" height="36" viewBox="0 0 760 30"><g font-family="Helvetica, Arial, sans-serif" font-size="12" font-weight="700" letter-spacing="2" fill="#5cc8ea"><text x="200" y="20" text-anchor="middle">THREAD A</text><text x="500" y="20" text-anchor="middle">THREAD B</text><text x="710" y="20" text-anchor="middle">INSTANCE</text></g></svg>

* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">1</text><rect x="60" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="200" y="25" text-anchor="middle" fill="#eaf4fb">instance == nullptr? yes</text><text x="710" y="25" text-anchor="middle" fill="#eaf4fb" font-weight="400">none</text></g></svg>
* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">2</text><rect x="360" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="500" y="25" text-anchor="middle" fill="#eaf4fb">instance == nullptr? yes</text><text x="710" y="25" text-anchor="middle" fill="#eaf4fb" font-weight="400">none</text></g></svg>
* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">3</text><rect x="60" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="200" y="25" text-anchor="middle" fill="#eaf4fb">instance = new Singleton()</text><text x="710" y="25" text-anchor="middle" fill="#5cc8ea" font-weight="700">#1</text></g></svg>
* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">4</text><rect x="360" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="500" y="25" text-anchor="middle" fill="#eaf4fb">instance = new Singleton()</text><text x="710" y="25" text-anchor="middle" fill="#f0a06e" font-weight="700">#2</text></g></svg>
* > Check-then-act again: two "only" instances, and one of them leaks.
* Double-checked locking tries to fix this cheaply — subtle enough
  that it waits for Day 5

---

<span class="tag">Day 1</span>

## Observer — the question you wrote down

> What happens if an observer subscribes **at the exact moment** the
> publisher is notifying?

```cpp
void notify(const Event& e) {
    for (auto* s : subscribers) s->update(e);   // thread A walks the vector
}
void subscribe(Subscriber* s) {
    subscribers.push_back(s);                   // thread B grows it
}
```

---

<!-- _class: invert bad -->

## Observer — the answer

* `push_back` may move the vector to new memory — thread A's loop keeps
  walking the old, freed memory
* Unsubscribe is worse: B removes and deletes its subscriber, A calls
  `update()` on an object that no longer exists
* > The answer is **undefined behavior** — a crash, on a good day.

---

<span class="tag">Day 2</span>

## The bounded stack — the sentence on the board

> An invariant may be broken **inside** a method, as long as it is
> restored before the method returns.

```cpp
void push(int v) {
    data[size] = v;    // step 1: fill the slot
    size++;            // step 2: count it
}
```

* One thread: nobody else runs between step 1 and step 2
* Two threads: the other `push` can run right in between

---

<!-- _class: invert bad -->

## Two pushes, one slot

`size` is 5. Thread A pushes `a`, thread B pushes `b`.

<svg class="uml tl" width="900" height="36" viewBox="0 0 760 30"><g font-family="Helvetica, Arial, sans-serif" font-size="12" font-weight="700" letter-spacing="2" fill="#5cc8ea"><text x="200" y="20" text-anchor="middle">THREAD A · push(a)</text><text x="500" y="20" text-anchor="middle">THREAD B · push(b)</text><text x="710" y="20" text-anchor="middle">STACK</text></g></svg>

* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">1</text><rect x="60" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="200" y="25" text-anchor="middle" fill="#eaf4fb">data[5] = a</text><text x="710" y="25" text-anchor="middle" fill="#5cc8ea" font-weight="700">slot 5: a</text></g></svg>
* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">2</text><rect x="360" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="500" y="25" text-anchor="middle" fill="#eaf4fb">data[5] = b</text><text x="710" y="25" text-anchor="middle" fill="#f0a06e" font-weight="700">slot 5: b</text></g></svg>
* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">3</text><rect x="60" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="200" y="25" text-anchor="middle" fill="#eaf4fb">size++</text><text x="710" y="25" text-anchor="middle" fill="#5cc8ea" font-weight="700">size 6</text></g></svg>
* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">4</text><rect x="360" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="500" y="25" text-anchor="middle" fill="#eaf4fb">size++</text><text x="710" y="25" text-anchor="middle" fill="#f0a06e" font-weight="700">size 7</text></g></svg>
* > Both pushes returned normally — yet `a` is gone, and slot 6 was
  > never written: the next `pop()` returns garbage.
* > The invariant did not get weaker. **"Between method calls" stopped
  > existing.**

---

<span class="tag">Day 2</span>

## Atomicity — the withdrawal

```cpp
bool withdraw(long long amount) {
    if (balance >= amount) {     // check
        balance -= amount;       // act
        return true;
    }
    return false;
}
```

* Each call, on its own, keeps the promise: the balance never goes
  below zero

---

<!-- _class: invert bad -->

## Two withdrawals of 100

`balance` is 100. Both threads call `withdraw(100)`.

<svg class="uml tl" width="900" height="36" viewBox="0 0 760 30"><g font-family="Helvetica, Arial, sans-serif" font-size="12" font-weight="700" letter-spacing="2" fill="#5cc8ea"><text x="200" y="20" text-anchor="middle">THREAD A · withdraw(100)</text><text x="500" y="20" text-anchor="middle">THREAD B · withdraw(100)</text><text x="710" y="20" text-anchor="middle">BALANCE</text></g></svg>

* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">1</text><rect x="60" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="200" y="25" text-anchor="middle" fill="#eaf4fb">balance &gt;= 100? yes</text><text x="710" y="25" text-anchor="middle" fill="#eaf4fb" font-weight="400">100</text></g></svg>
* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">2</text><rect x="360" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="500" y="25" text-anchor="middle" fill="#eaf4fb">balance &gt;= 100? yes</text><text x="710" y="25" text-anchor="middle" fill="#eaf4fb" font-weight="400">100</text></g></svg>
* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">3</text><rect x="60" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="200" y="25" text-anchor="middle" fill="#eaf4fb">balance -= 100</text><text x="710" y="25" text-anchor="middle" fill="#5cc8ea" font-weight="700">0</text></g></svg>
* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">4</text><rect x="360" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="500" y="25" text-anchor="middle" fill="#eaf4fb">balance -= 100</text><text x="710" y="25" text-anchor="middle" fill="#f0a06e" font-weight="700">−100</text></g></svg>
* > Two correct calls, jointly wrong. The promise held **per call** —
  > it said nothing about two calls at once.

---

## The thesis of this course

> Concurrency does not put bugs into bad code. It breaks the
> assumptions that made **good** code good.

* Our job from here on: get the guarantees back — without giving up
  the speed we came for

---

<!-- _class: lead invert -->

# The First Fix: Mutexes

---

## Mutual exclusion

* > **Mutual exclusion** — at most one thread inside a critical section
  > at a time.
* > A **mutex** enforces it. `lock()` waits until the mutex is free,
  > then takes it; `unlock()` gives it back.
* A thread that calls `lock()` while another thread holds it **waits** —
  the OS parks it until the mutex is free

---

## One mutex, two threads

* <svg class="uml" width="880" height="289" viewBox="0 0 760 250"><g font-family="Helvetica, Arial, sans-serif" font-size="15"><text x="10" y="73" fill="#eaf4fb" font-weight="700">Thread A</text><line x1="170" y1="44" x2="170" y2="92" stroke="#5cc8ea" stroke-width="1.6"/><text x="170" y="38" text-anchor="middle" fill="#5cc8ea" font-family="'IBM Plex Mono', Menlo, monospace" font-size="13">lock()</text><rect x="170" y="50" width="260" height="36" fill="#14385a" stroke="#7ed0a3" stroke-width="1.8"/><text x="300" y="73" text-anchor="middle" fill="#eaf4fb">critical section</text><line x1="430" y1="44" x2="430" y2="92" stroke="#5cc8ea" stroke-width="1.6"/><text x="430" y="38" text-anchor="middle" fill="#5cc8ea" font-family="'IBM Plex Mono', Menlo, monospace" font-size="13">unlock()</text><text x="10" y="173" fill="#eaf4fb" font-weight="700">Thread B</text><line x1="240" y1="144" x2="240" y2="192" stroke="#5cc8ea" stroke-width="1.6"/><text x="240" y="138" text-anchor="middle" fill="#5cc8ea" font-family="'IBM Plex Mono', Menlo, monospace" font-size="13">lock()</text><rect x="240" y="150" width="190" height="36" fill="none" stroke="#6fa8cc" stroke-width="1.4" stroke-dasharray="6 5"/><text x="335" y="173" text-anchor="middle" fill="#a9c8de" font-style="italic">waiting</text><rect x="430" y="150" width="260" height="36" fill="#14385a" stroke="#7ed0a3" stroke-width="1.8"/><text x="560" y="173" text-anchor="middle" fill="#eaf4fb">critical section</text><line x1="690" y1="144" x2="690" y2="192" stroke="#5cc8ea" stroke-width="1.6"/><text x="690" y="138" text-anchor="middle" fill="#5cc8ea" font-family="'IBM Plex Mono', Menlo, monospace" font-size="13">unlock()</text><line x1="150" y1="225" x2="744" y2="225" stroke="#6fa8cc" stroke-width="1.6"/><polygon points="750,225 741,220 741,230" fill="#6fa8cc"/><text x="750" y="245" text-anchor="end" fill="#a9c8de" font-size="12">time</text></g></svg>

* B asked for the lock while A held it — so B's critical section runs
  **after** A's, never during

---

## Question

How do we fix the counter?

```cpp
int counter = 0;                 // shared by both threads

void work() {
    for (int i = 0; i < 1000000; ++i) counter++;
}
```

---

<!-- _class: invert good -->

## ✅ The counter, fixed

```cpp
int counter = 0;
std::mutex m;                    // guards counter

void work() {
    for (int i = 0; i < 1000000; ++i) {
        m.lock();                // one thread at a time from here...
        counter++;
        m.unlock();              // ...to here
    }
}
```

* `2000000` — every run

<p class="note">C++ note: <code>lock()</code> / <code>unlock()</code> by hand is here only to show the idea. Real C++ code uses <code>std::lock_guard</code>, which unlocks by itself — even on an early <code>return</code> or an exception.</p>

---

<span class="tag">2 threads × 10,000,000 increments · a 10-core laptop</span>

## Correct — and slower

* **1 thread** — 30 ms, correct
* **2 threads, no mutex** — 23 ms, wrong: about half the increments lost
* **2 threads, one mutex** — 340 ms, correct — and **10× slower** than
  one thread
* Every increment now fights for the same lock: **contention**
* Correct vs. fast — the tension this whole field lives with

---

## A mutex protects data, not code

* Write down which data each mutex guards — right next to the data

```cpp
class Account {
    std::mutex m;         // guards balance
    long long balance;    // touch only while holding m
};
```

* One getter that reads `balance` without `m` — still broken
* Two different mutexes guarding the same data protect nothing

---

<span class="tag">Snippet 1</span>

## Question

How do we fix the image cache?

```cpp
std::unordered_map<std::string, Image> cache;    // shared

Image load(const std::string& path) {
    auto it = cache.find(path);
    if (it != cache.end()) return it->second;
    Image img = readFromDisk(path);
    cache[path] = img;
    return img;
}
```

---

<!-- _class: invert good -->

## ✅ The image cache, fixed

```cpp
std::mutex m;                                  // guards cache

Image load(const std::string& path) {
    m.lock();
    auto it = cache.find(path);
    if (it != cache.end()) {
        Image img = it->second;                // copy it while still locked
        m.unlock();
        return img;
    }
    Image img = readFromDisk(path);            // still locked: correct, but slow
    cache[path] = img;
    m.unlock();
    return img;
}
```

* Find, read and insert are now one step — each file is read once
* The price: while one thread reads the disk, every other thread waits —
  even for images already in the cache. Keeping locks short: Day 4.

---

<span class="tag">Snippet 3</span>

## Question

How do we fix the order numbers?

```cpp
int lastId = 0;                          // shared

std::string newOrderId() {
    lastId = lastId + 1;
    return "ORD-" + std::to_string(lastId);
}
```

---

<!-- _class: invert good -->

## ✅ The order numbers, fixed

```cpp
int lastId = 0;
std::mutex m;                            // guards lastId

std::string newOrderId() {
    m.lock();
    lastId = lastId + 1;
    int id = lastId;                     // read it before unlocking
    m.unlock();
    return "ORD-" + std::to_string(id);  // touches nothing shared
}
```

* Both the increment and the read of the new value happen under the lock
* Building the string needs no lock — keep the critical section small

---

<span class="tag">Snippet 4</span>

## Question

How do we fix the job queue?

```cpp
std::queue<Job> jobs;                    // shared, filled by another thread

void worker() {
    for (;;) {
        if (!jobs.empty()) {
            Job j = jobs.front();
            jobs.pop();
            run(j);
        }
    }
}
```

---

<!-- _class: invert good -->

## ✅ The job queue, fixed

```cpp
std::mutex m;                            // guards jobs — the producer locks it too

void worker() {
    for (;;) {
        m.lock();
        if (jobs.empty()) {
            m.unlock();
            continue;
        }
        Job j = jobs.front();
        jobs.pop();
        m.unlock();
        run(j);                          // outside the lock: j is ours now
    }
}
```

* One lock around `empty` → `front` → `pop`; `run(j)` runs outside it
* An empty queue makes the worker spin and burn a core — waiting properly
  needs a condition variable: Day 4

---

<!-- _class: invert good -->

## ✅ The bounded stack, fixed

```cpp
class BoundedStack {
    std::mutex m;               // guards data and size
    std::vector<int> data;      // `capacity` slots, fixed at construction
    int size = 0;
public:
    void push(int v) {
        m.lock();
        if (size == (int)data.size()) {
            m.unlock();
            throw std::runtime_error("full");
        }
        data[size] = v;
        size++;
        m.unlock();
    }
};
```

* Step 1 and step 2 still happen one after the other — but no other
  thread can get between them. "Between method calls" exists again.

---

<!-- _class: lead invert -->

# What Mutexes Don't Fix

### …and the new problems they bring

---

## Every method locked — still a race

```cpp
class Stack {                   // each method locks m inside
public:
    bool empty();
    int  top();
    void pop();
};

if (!s.empty()) {               // check
    int v = s.top();            // act, part 1
    s.pop();                    // act, part 2
    process(v);
}
```

* Each call is thread-safe. Is this code?

---

<!-- _class: invert bad -->

## Two threads, one thread-safe stack

The stack holds `[3, 7]`. Both threads run the code from the last slide.

<svg class="uml tl" width="900" height="36" viewBox="0 0 760 30"><g font-family="Helvetica, Arial, sans-serif" font-size="12" font-weight="700" letter-spacing="2" fill="#5cc8ea"><text x="200" y="20" text-anchor="middle">THREAD A</text><text x="500" y="20" text-anchor="middle">THREAD B</text><text x="710" y="20" text-anchor="middle">STACK</text></g></svg>

* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">1</text><rect x="60" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="200" y="25" text-anchor="middle" fill="#eaf4fb">empty()? no</text><text x="710" y="25" text-anchor="middle" fill="#eaf4fb" font-weight="400">[3, 7]</text></g></svg>
* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">2</text><rect x="360" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="500" y="25" text-anchor="middle" fill="#eaf4fb">empty()? no</text><text x="710" y="25" text-anchor="middle" fill="#eaf4fb" font-weight="400">[3, 7]</text></g></svg>
* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">3</text><rect x="60" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="200" y="25" text-anchor="middle" fill="#eaf4fb">top() → 7</text><text x="710" y="25" text-anchor="middle" fill="#eaf4fb" font-weight="400">[3, 7]</text></g></svg>
* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">4</text><rect x="360" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="500" y="25" text-anchor="middle" fill="#eaf4fb">top() → 7</text><text x="710" y="25" text-anchor="middle" fill="#eaf4fb" font-weight="400">[3, 7]</text></g></svg>
* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">5</text><rect x="60" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="200" y="25" text-anchor="middle" fill="#eaf4fb">pop()</text><text x="710" y="25" text-anchor="middle" fill="#5cc8ea" font-weight="700">[3]</text></g></svg>
* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">6</text><rect x="360" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="500" y="25" text-anchor="middle" fill="#eaf4fb">pop()</text><text x="710" y="25" text-anchor="middle" fill="#f0a06e" font-weight="700">[ ]</text></g></svg>
* > `7` was processed twice and `3` was lost. Every call was locked — the
  > **sequence** was not. An **API race**: check-then-act, one level up.

---

<!-- _class: invert good -->

## ✅ One call that does the whole job

```cpp
std::optional<int> tryPop() {           // check and act, under one lock
    m.lock();
    if (size == 0) {
        m.unlock();
        return std::nullopt;
    }
    int v = data[--size];
    m.unlock();
    return v;
}

if (auto v = s.tryPop()) process(*v);   // the caller can't split it anymore
```

* A thread-safe class needs an interface where no caller has to combine
  calls: `empty()` + `top()` + `pop()` → one `tryPop()`
* Same idea elsewhere: Java's `ConcurrentHashMap.putIfAbsent`, Go's
  `sync.Map.LoadOrStore`

---

<!-- _class: invert bad -->

## New problem: deadlock

```cpp
void transfer(Account& from, Account& to, long long amount) {
    from.m.lock();
    to.m.lock();
    from.balance -= amount;
    to.balance += amount;
    to.m.unlock();
    from.m.unlock();
}
```

* Thread A calls `transfer(x, y, 10)`; at the same moment, thread B calls
  `transfer(y, x, 20)`.

---

<!-- _class: invert bad -->

## Two transfers, two locks

<svg class="uml tl" width="900" height="36" viewBox="0 0 760 30"><g font-family="Helvetica, Arial, sans-serif" font-size="12" font-weight="700" letter-spacing="2" fill="#5cc8ea"><text x="200" y="20" text-anchor="middle">THREAD A · transfer(x, y)</text><text x="500" y="20" text-anchor="middle">THREAD B · transfer(y, x)</text><text x="710" y="20" text-anchor="middle">LOCKS</text></g></svg>

* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">1</text><rect x="60" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="200" y="25" text-anchor="middle" fill="#eaf4fb">x.m.lock()</text><text x="710" y="25" text-anchor="middle" fill="#5cc8ea" font-weight="700">A has x</text></g></svg>
* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">2</text><rect x="360" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="500" y="25" text-anchor="middle" fill="#eaf4fb">y.m.lock()</text><text x="710" y="25" text-anchor="middle" fill="#5cc8ea" font-weight="700">B has y</text></g></svg>
* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">3</text><rect x="60" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="200" y="25" text-anchor="middle" fill="#eaf4fb">y.m.lock() … waits</text><text x="710" y="25" text-anchor="middle" fill="#eaf4fb" font-weight="400">A waits</text></g></svg>
* <svg class="uml tl" width="900" height="47" viewBox="0 0 760 40"><g font-family="'IBM Plex Mono', Menlo, monospace" font-size="15"><line x1="200" y1="0" x2="200" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><line x1="500" y1="0" x2="500" y2="40" stroke="#6fa8cc" stroke-width="1" stroke-dasharray="3 4" opacity="0.45"/><text x="24" y="25" text-anchor="middle" fill="#a9c8de" font-family="Helvetica, Arial, sans-serif" font-size="13">4</text><rect x="360" y="4" width="280" height="32" fill="#14385a" stroke="#6fa8cc" stroke-width="1.4"/><text x="500" y="25" text-anchor="middle" fill="#eaf4fb">x.m.lock() … waits</text><text x="710" y="25" text-anchor="middle" fill="#f0a06e" font-weight="700">both wait</text></g></svg>
* > **Deadlock**: each thread holds the lock the other one needs. Neither
  > can go on — ever.
* Run both in a loop and the program freezes within seconds. Why it
  happens, and how to prevent it: Day 4.

---

<!-- _class: invert bad -->

## New problem: locking twice

```cpp
long long balance() { m.lock(); long long b = bal; m.unlock(); return b; }

bool withdraw(long long amount) {
    m.lock();
    if (balance() < amount) {    // balance() locks m again: waits for itself
        m.unlock();
        return false;
    }
    bal -= amount;
    m.unlock();
    return true;
}
```

* The thread already holds `m`, so its second `lock()` waits for itself —
  forever
* An easy accident: one locked method calling another. C++ `std::mutex`,
  Go's `sync.Mutex` and Python's `Lock` all wait like this; Java's
  `synchronized` lets the same thread lock twice

---

<!-- _class: invert bad -->

## New problem: someone else's code under your lock

```cpp
void Publisher::notify(const Event& e) {
    m.lock();
    for (auto* s : subscribers)
        s->update(e);                  // someone else's code runs with m held
    m.unlock();
}
void Publisher::subscribe(Subscriber* s) { m.lock(); subscribers.push_back(s); m.unlock(); }

struct Welcomer : Subscriber {         // someone else's code
    Publisher& pub;
    void update(const Event&) override {
        pub.subscribe(new Greeter());  // inside notify(): m.lock() again
    }
};
```

* `update()` calls `subscribe()` → `m.lock()` again: locking twice
* `update()` is slow, or waits for a thread that needs `m` → everyone stops
* Safer: copy the list under the lock, unlock, **then** call them

---

## Mutexes: what they fix, what they bring

* **They fix** — data races, and other threads seeing half-done states
* **They don't fix** — API races: a sequence of locked calls is not one
  locked call
* **They bring** — deadlock, locking twice, waiting on someone else's code,
  and contention: slower code

---

## Practice — `playground/<language>/days/day03`

* **Core** — a ticket dispenser (the counter), a bank account
  (check-then-act, and transfers that must not deadlock), the bounded
  stack, shared (with `tryPush` / `tryPop`)
* **Challenge** — compute once per key (the Singleton), a thread-safe
  event bus (the Observer), the interval booker, shared
* The same tasks in C++, Python, Java and Go

---

<!-- _class: lead invert -->

# The guarantees are restorable

### One mutex · the right data · the right scope
