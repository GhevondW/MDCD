# Team Project — Project Ideas (v2)

Seven projects. Each one is a small version of a real system, built **inside
one process**, using only what the course teaches: mutexes, condition
variables, bounded queues, thread pools, futures, atomics, deadlock
avoidance.

Pick one. Teams of 2–4. Any language (C++, Java, C#, Python, Go, Rust).

## At a glance

| # | Project | Real-life version | Classic problem inside | Level |
|---|---|---|---|---|
| 1 | DAG Task Executor | Make, Airflow, Spark | producer-consumer + dependency counting | ★★ |
| 2 | Mini Transactional Database | SQLite, Postgres | readers-writers + deadlock | ★★★ |
| 3 | Streaming Pipeline | Logstash, Kafka Streams | producer-consumer, backpressure | ★★ |
| 4 | Message Broker | RabbitMQ, SQS | producer-consumer + acknowledgements | ★★ |
| 5 | Cache with TTL | Memcached, Redis | striping + single-flight | ★ |
| 6 | Web Crawler | Googlebot, Scrapy | work queue + "are we done?" | ★★ |
| 7 | Live Search Index | Lucene, Elasticsearch | readers-writers, copy-on-write | ★★ |

★ = good first choice · ★★ = solid challenge · ★★★ = hardest, best for 3–4
people.

---

## 1. DAG Task Executor ★★

**Problem.** A client describes work as a graph. Each node is a function;
an edge means "this node needs the output of that node". The executor runs
the graph on N worker threads: every node starts as soon as all its inputs
are ready, independent nodes run at the same time, and outputs flow along
the edges.

**Real life.**
- `make` / Ninja / Bazel: compile every file in parallel, link at the end.
- Airflow, dbt: "load → clean → join → report" data pipelines.
- Spark: a query becomes a graph of stages.
- GitHub Actions: jobs with `needs:`.

**Example.** Log analysis: `read(file1..file8)` → `parse` each →
`count words` each → `merge all counts` → `top 10` → `write report`.
Eight branches run in parallel, then they join.

**You build.**
- `graph.add(name, function, dependsOn=[...])` — the function receives the
  outputs of its dependencies and returns its own output.
- `graph.validate()` — rejects cycles (and names the cycle) and unknown
  dependencies.
- `executor.run(graph) -> Future<Result>` — many graphs can be submitted
  from many threads at the same time and share one pool of N workers.
- `executor.cancel(run)`, `executor.shutdown()`.
- `Result` — output or error or "skipped" for every node.

**Must work.**
- A node runs only after all its dependencies succeeded, and **exactly
  once**.
- If a node fails: everything downstream is marked *skipped*; independent
  branches keep running. Optional retry count per node.
- Cancel: nodes not yet started never start; running ones finish; the final
  state of every node is consistent.
- An output is released from memory as soon as all its consumers finished.
- Same final results with 1 worker and with 8 workers.
- A graph with 1,000 nodes finishes on a pool of 2 workers (no deadlock from
  workers waiting for each other).
- A huge graph does not starve a tiny graph submitted later.
- Clean shutdown: no thread left running, no node left "running" forever.

**Where it gets hard.**
- Two parents finish at the same instant — the child must be scheduled
  exactly once (the "last finisher wins" counter).
- Cancel racing with a node about to start. Failure racing with completion.
- Fairness between graphs sharing one pool.

**How we check it.**
- Generate random graphs where each node computes a value from its parents.
  Compare with a plain single-threaded evaluation.
- An event log proves: every node started after all its parents finished,
  and ran once.
- Benchmark three shapes: a chain (no speedup possible), a wide fan-out
  (near-linear), a diamond. Explain each result with the *critical path*.
- Naive mode: check "are all parents done?" without a lock → nodes run twice
  or never.

**Stretch.** Cache results by input hash and skip unchanged nodes (like
`make`); run the longest-path nodes first; limit "heavy" nodes to 2 at a
time; print a timeline (who ran what, when).

**Learn on your own.** Topological sort (Kahn's algorithm), critical-path
scheduling.

---

## 2. Mini Transactional Database ★★★

**Problem.** Store key → value data on disk so that a group of changes is
**all-or-nothing**, survives a crash, and many clients can use it at once
without seeing each other's half-finished work. These four promises are
ACID:
- **Atomic** — a transaction fully happens or does not happen at all.
- **Consistent** — the data rules still hold after every transaction.
- **Isolated** — concurrent transactions behave as if they ran one by one.
- **Durable** — once `commit` returns, the data survives a crash.

**Real life.** SQLite, PostgreSQL, the write-ahead log inside LevelDB /
RocksDB. Every bank transfer you ever made.

**Example.** Transfer 100 from account A to B: `A -= 100; B += 100`. A crash
between the two lines must never lose the money; two transfers touching A
at the same time must never lose an update.

**You build.**
- `open(directory)` — loads data and recovers from the log.
- `begin() -> Txn`
- `txn.get(key)`, `txn.put(key, value)`, `txn.delete(key)`
- `txn.commit()`, `txn.abort()`
- Errors: `Deadlock` (transaction was chosen as victim, client retries).

**Must work.**
- **Isolation by strict two-phase locking:** each key has a read lock
  (shared) and a write lock (exclusive). A transaction takes locks as it
  touches keys and releases them all only at commit/abort. Many readers on
  one key at once; a writer is alone.
- **Deadlocks:** transaction 1 locks X then Y, transaction 2 locks Y then X.
  Users choose the order, so you cannot fix it by ordering. Detect it (a
  cycle in the wait-for graph) or use a lock timeout, abort one transaction
  and let the client retry.
- **Durability with a write-ahead log:** a transaction keeps its writes in
  memory until commit. At commit, append **one** log record with all its
  writes and a checksum, `fsync` it, and only then apply to memory and
  return. On startup, replay the valid records; ignore a half-written last
  record. (This design needs no undo.)
- **Group commit:** when 16 transactions commit at the same moment, they
  share one `fsync`. Measure commits/s with and without it.
- At least 16 client threads at once.

**Out of scope.** SQL, indexes, range scans, B-trees, data bigger than RAM.

**Where it gets hard.**
- Lock manager: waiting, waking, upgrading read → write lock (two readers
  both upgrading the same key = deadlock).
- Group commit: one thread does the `fsync`, the others wait for it.
- Commit order in the log must match the order the locks imply.

**How we check it.**
- Bank test: 100 accounts, 16 threads doing random transfers. The total
  money never changes; no negative balance.
- Crash test: run the database in a child process, `kill -9` it at random
  moments while clients run, restart it, check the total and that **every
  commit that returned success is present**.
- Damaged-log test: cut or corrupt the last bytes of the log; recovery must
  still work and ignore the broken record.
- Naive mode: release locks right after each operation → lost updates, and
  the bank total drifts.
- Be honest in the doc: `kill -9` tests your recovery logic, not a power
  cut (the OS still has the data). Say what `fsync` does and does not
  guarantee.

**Stretch.** Checkpoint and shrink the log; snapshot reads so readers never
block (MVCC); range scans.

**Learn on your own.** Write-ahead logging, `fsync` (on macOS plain `fsync`
is not enough — look up `F_FULLFSYNC`), two-phase locking, CRC checksums.

---

## 3. Streaming Pipeline ★★

**Problem.** Process an endless stream of records through stages, for
example *read → parse → filter → enrich → aggregate → write*. Stages have
different speeds. A slow stage must slow the fast ones down (**backpressure**)
instead of letting memory grow without limit.

**Real life.** Logstash and Fluentd (log processing), Kafka Streams / Flink,
ffmpeg filter chains, Unix pipes.

**Example.** A web server writes 50,000 log lines per second. Parse each
line, drop health-check requests, look up the country of the IP (slow),
count requests per country, write the result every second.

**You build.**
- `Pipeline.source(generator).stage(name, function, workers=k)....sink(function)`
- `start()`, `stop(drain=true|false)`
- `stats()` — per stage: records in/out, queue length, % time busy.

**Must work.**
- Bounded queue between every two stages, size configurable.
- A stage with `k` workers processes `k` records in parallel.
- **Ordered mode:** output reaches the sink in input order even when a stage
  has several workers (reorder buffer). **Unordered mode:** fastest first.
- A record that throws goes to a dead-letter list; the pipeline continues.
- A stage can drop a record (filter) or emit several (split).
- `stop(drain=true)`: every record already accepted reaches the sink exactly
  once, then all threads end. `stop(drain=false)`: stop fast, no hang.
- Memory stays flat when the sink is 100× slower than the source.
- `stats()` shows which stage is the bottleneck.

**Where it gets hard.**
- Reorder buffer: a worker holding record #5 is slow, workers behind it
  keep finishing #6, #7, … — the buffer must be bounded without deadlock.
- Shutdown that flows stage by stage in the right order.

**How we check it.**
- Records carry an id. Sink sees every id exactly once (or it is in the
  dead-letter list). Ordered mode: ids strictly increasing.
- Sampler thread: no queue ever longer than its bound.
- Benchmark: add workers to the bottleneck stage and watch the throughput
  rise and the bottleneck move to the next stage.
- Naive mode: unbounded queues → memory grows; no reorder buffer → order
  breaks.

**Stretch.** Split a stream into two branches and join them; windowed
aggregation ("count per 10 seconds"); group records into batches of 100.

**Learn on your own.** Backpressure; sequence numbers for reordering.

---

## 4. Message Broker ★★

**Problem.** Programs send messages to named **topics**; other programs
receive them. If a receiver crashes while handling a message, the message
must not be lost — someone else gets it.

**Real life.** RabbitMQ, Amazon SQS, Kafka, Redis Streams. Example: an
"order placed" message goes to the billing service and the email service;
if one billing worker crashes, another one picks the order up.

**You build.**
- `createTopic(name, maxMessages)`
- `publish(topic, message) -> id` (blocks or fails when the topic is full —
  your choice, documented)
- `subscribe(topic, group) -> Consumer`
- `consumer.poll(timeout) -> Delivery | nothing`
- `ack(delivery)`, `nack(delivery)` (give it back for another try)

**Must work.**
- **Groups:** every group gets every message (billing and email both see
  each order). Inside one group, each message goes to **one** consumer
  (workers share the load).
- **Ack deadline:** a delivered message must be acked in N seconds,
  otherwise it is delivered again, possibly to another consumer
  (at-least-once).
- Consumer disconnects → its unacked messages come back at once.
- After M failed deliveries → message goes to a dead-letter topic.
- A late `ack` for a delivery that already expired is rejected cleanly (the
  message was re-delivered; don't count it twice).
- Order: with one consumer per group, messages arrive in publish order.
  State exactly what you promise with many consumers.
- A slow group never slows another group or the publishers.
- Clean shutdown.

**Where it gets hard.**
- Ack arrives at the same moment the deadline expires.
- Two consumers polling at once must never get the same message.
- New group created while messages are being published.

**How we check it.**
- Load generator publishes 100,000 numbered messages; consumers randomly
  crash, stall, nack. At the end: every message acked exactly once per
  group, none lost.
- Sampler: a message is never in flight to two consumers of the same group
  at the same moment.
- Naive mode: pick-then-mark-in-flight without a lock → duplicate delivery.

**Stretch.** Persist the topic to a file; delayed messages; "same key → same
consumer" so per-key order holds.

**Learn on your own.** At-most-once vs at-least-once vs exactly-once;
visibility timeout (SQS); consumer groups (Kafka).

---

## 5. Cache with TTL ★

**Problem.** Reading from a slow source (database, remote API) is
expensive. Keep recent results in memory. Memory is limited, entries get
old, and when 1,000 threads ask for the same missing key, the slow source
must be asked **once**, not 1,000 times (a *cache stampede*).

**Real life.** Memcached, Redis, Caffeine (Java), CDN caches. Stampedes have
taken down real websites when a popular key expired.

**You build.**
- `get(key)`, `put(key, value, ttl)`, `delete(key)`
- `getOrLoad(key, loader)` — returns the cached value or calls `loader` once
- `stats()` — hits, misses, evictions, loads

**Must work.**
- Capacity N entries; when full, evict the **least recently used**.
- An entry past its TTL is never returned; a background thread removes
  expired entries.
- **Single-flight:** 1,000 threads call `getOrLoad` for the same cold key →
  `loader` runs once, all get the result. If the loader throws, all waiters
  get the error and the next call tries again.
- The loader runs **outside any lock**, and different keys load in parallel.
- Size never exceeds capacity by more than a bound you state.
- `hits + misses` equals the number of `get` calls.
- Compare three designs: one global lock, striped locks, per-shard LRU.

**Where it gets hard.**
- LRU changes its order even on a *read*, so a read-write lock does not help.
- Background expiry racing with `get`.

**How we check it.**
- Loader counts its calls: 1,000 threads, one cold key → exactly 1.
- Values contain their own key, so a mixed-up value is detected.
- Sampler: size ≤ capacity bound at all times.
- Benchmark with a skewed load (90% of requests hit 10% of keys): where does
  striping stop helping?
- Naive mode: check-then-load without single-flight → loader called many
  times.

**Stretch.** Approximate LRU (CLOCK / Redis-style sampling); refresh a hot
entry before it expires; write-behind.

**Learn on your own.** LRU with hash map + linked list; cache stampede;
CLOCK algorithm.

---

## 6. Web Crawler ★★

**Problem.** Visit every page reachable from some start pages, quickly, with
many threads — but never visit a page twice, never hit one website too hard,
and know exactly when the crawl is finished.

**Real life.** Googlebot, Scrapy, the Internet Archive, any "broken link
checker".

**No real internet.** You write (or I give you) a `FakeWeb`: thousands of
generated pages on many hosts, each with links and a random delay of
10–300 ms; some pages fail or time out. Real HTTP is a stretch goal.

**You build.**
- `crawl(seeds, maxPages, maxDepth, threads) -> CrawlResult`
- `pause()`, `resume()`, `stop()`
- `stats()` — pages/s, queue length, errors.

**Must work.**
- Each URL is fetched exactly once (retries of failures aside).
- **Politeness:** at most 2 requests at once to the same host, and at least
  100 ms between requests to the same host. Other hosts keep going; a slow
  host must not freeze the workers.
- Failed fetch → retry up to 3 times with growing delay, then give up.
- **Exact stop:** with `maxPages = 1000` the crawler fetches exactly 1,000.
- **Exact finish:** the crawl ends when the queue is empty **and** no fetch is
  in progress — not before, not never.
- Fair: no host is starved while another has a million pages.

**Where it gets hard.**
- Termination: a worker sees an empty queue while another worker is about to
  add new links. Declaring "done" too early loses pages.
- A visited-set check-then-insert race (two workers fetch the same page).
- Per-host limits without a thread waiting idle.

**How we check it.**
- `FakeWeb` knows the true set of reachable pages: the crawler must find all
  of them (within limits). Same set with 1 thread and 32 threads.
- Fetch log: no URL twice; no moment with 3 requests to one host; gaps ≥ 100
  ms.
- Benchmark: pages/s vs threads. With few hosts, politeness caps the
  speed — show it and explain it.
- Naive mode: plain `if (!visited.contains(url)) visited.add(url)` →
  duplicates.

**Stretch.** Crawl shallow pages first; `robots.txt`-style rules; save the
queue to disk and resume.

**Learn on your own.** Termination detection; URL normalization;
per-host rate limiting.

---

## 7. Live Search Index ★★

**Problem.** Index text documents so that searches are fast. Documents are
added, changed and deleted all the time, while people search. A search must
see a **consistent** state: never half of an update, never a deleted
document mixed with an old index.

**Real life.** Lucene / Elasticsearch / Solr ("near-real-time search"),
database full-text search, the search box on any website.

**Example.** 50,000 articles indexed. Four threads keep adding and editing
articles. Eight threads search "concurrent AND queue". Searches never wait
for the indexers, and new documents show up within a second.

**You build.**
- `addDocument(id, text)`, `updateDocument(id, text)`, `deleteDocument(id)`
- `search(query, k) -> (top k hits, version)` — words joined by AND, ranked
  by how often the words appear (simple TF-IDF is fine)
- Updates become visible within 1 second.

**Must work.**
- An inverted index: word → list of documents.
- Several indexer threads and several search threads at the same time.
- Every search sees the **exact state at one version V**: all updates up to
  V, none after.
- Searches never block on indexing. Choose and justify: immutable segments
  swapped in atomically (copy-on-write), or a read-write lock. Back the
  choice with a benchmark.
- A background thread **merges** small segments into bigger ones without
  stopping searches.
- A deleted document disappears from results at the moment of the version
  that deleted it.

**Where it gets hard.**
- Switching to the new index while searches are running on the old one;
  freeing the old one only when its last search is done.
- Merge running while updates arrive.

**How we check it.**
- Every update has a sequence number. A checker thread re-runs the same
  query by brute force on the document state at version V and compares.
- Marker test: one update changes two words together (`alpha-5` and
  `beta-5`). No search may ever return `alpha-5` with `beta-4`.
- Benchmark: queries/s with 0, 1, 4 indexers; p99 latency while a merge
  runs.
- Naive mode: update the index in place while searching → torn results.

**Stretch.** Phrase queries; prefix search; save segments to disk.

**Learn on your own.** Inverted index; TF-IDF; copy-on-write / immutable
snapshots.

---
