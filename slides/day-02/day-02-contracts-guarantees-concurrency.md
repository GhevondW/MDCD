---
marp: true
theme: gaia
class: invert
paginate: true
footer: 'Contracts, Guarantees & Concurrency'
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

  /* Step-by-step reveals: list items written with `*` appear one
     keypress at a time in the HTML export. Items that carry a quote
     panel shouldn't also show a list bullet. */
  li:has(> blockquote), li:has(svg) { list-style: none; }
  li > blockquote { margin: 0.2em 0 0.5em; }
---

<!-- _class: lead invert -->

# Contracts, Guarantees & Concurrency

### Modeling, Design & Collaborative Development

---

<!-- _class: lead invert -->

# Contracts & Invariants

---

## Every function makes a promise

* > **Precondition** — what the caller must guarantee before calling.
* > **Postcondition** — what the function guarantees after it returns.
* > **Invariant** — what stays true about the object *between* calls.

---

## A bounded stack

* <svg class="uml" width="640" height="180" viewBox="0 0 700 200"><g font-family="Helvetica, Arial, sans-serif" font-size="15"><rect x="30" y="50" width="60" height="60" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="60" y="85" text-anchor="middle" fill="#eaf4fb">v0</text><rect x="100" y="50" width="60" height="60" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="130" y="85" text-anchor="middle" fill="#eaf4fb">v1</text><rect x="170" y="50" width="60" height="60" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="200" y="85" text-anchor="middle" fill="#eaf4fb">v2</text><rect x="240" y="50" width="60" height="60" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="270" y="85" text-anchor="middle" fill="#eaf4fb">v3</text><rect x="310" y="50" width="60" height="60" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="340" y="85" text-anchor="middle" fill="#eaf4fb">v4</text><rect x="380" y="50" width="60" height="60" fill="none" stroke="#6fa8cc" stroke-width="1.4" stroke-dasharray="5 4" opacity="0.5"/><rect x="450" y="50" width="60" height="60" fill="none" stroke="#6fa8cc" stroke-width="1.4" stroke-dasharray="5 4" opacity="0.5"/><rect x="520" y="50" width="60" height="60" fill="none" stroke="#6fa8cc" stroke-width="1.4" stroke-dasharray="5 4" opacity="0.5"/><line x1="30" y1="35" x2="370" y2="35" stroke="#7ed0a3" stroke-width="1.6"/><line x1="30" y1="30" x2="30" y2="40" stroke="#7ed0a3" stroke-width="1.6"/><line x1="370" y1="30" x2="370" y2="40" stroke="#7ed0a3" stroke-width="1.6"/><text x="200" y="22" text-anchor="middle" fill="#7ed0a3" font-weight="700">size = 5</text><line x1="30" y1="130" x2="580" y2="130" stroke="#6fa8cc" stroke-width="1.4"/><text x="305" y="152" text-anchor="middle" fill="#eaf4fb" opacity="0.8">capacity = 8</text></g></svg>

* Invariant: `0 <= size <= capacity`; elements `0 .. size-1` are valid.

---

## push() — precondition and postcondition

```cpp
class BoundedStack {
    vector<int> data;      // fixed capacity, set at construction
    int size = 0;
public:
    // precondition:  size < capacity  (not full)
    void push(int v) {
        data[size] = v;     // postcondition: size == old(size) + 1
        size++;              //                data[size-1] == v
    }
    bool full() const { return size == (int)data.capacity(); }
};
```

---

## Who checks what

* **Defensive checks** — validate at the boundary, fail loud on bad input
* **Narrow, documented contracts** — trust the caller, check nothing
  (faster, riskier)
* **Assert** — catches a programmer bug; typically compiled out in release
* **Exception** — a recoverable condition; part of the documented contract

---

## Who checks what — in code

```cpp
assert(!full());                                // programmer bug — stripped in release
if (full()) throw std::runtime_error("full");   // recoverable — part of the contract
```

---

<!-- _class: lead invert -->

# What a Well-Designed API Guarantees

---

## An API is a promise

> An API is a **promise to strangers** — callers who will never
> read its source code.

* Three promises worth naming: **idempotency**, **atomicity**,
  **clear error semantics**.

---

## Idempotency

* > **Idempotent**: calling twice has the same effect as calling once.
* Callers **retry** on timeout — did the first call even happen?
* HTTP: `PUT` is idempotent by contract; `POST` is not
* Stripe-style **idempotency keys**: "make sure this order is paid,"
  not "charge the card again"

---

## Idempotency — in code

```cpp
if (seen(idempotencyKey))
    return cachedResult(idempotencyKey);

auto result = chargeCard(amount);
remember(idempotencyKey, result);
return result;
```

---

## Atomicity

* > **Atomic**: all-or-nothing. No caller ever sees a half-done state.
* A game is saving and the power goes out. What do you want to find
  when you restart?
* Your **old save**, complete and working — not a broken file that is
  half old, half new
* Money transfer: never take money out without putting it in

---

## Atomicity — in code

```cpp
writeFile("save.tmp", gameState);   // 1. write the new save into a NEW file
rename("save.tmp", "save.dat");     // 2. swap it in — one instant step
```

* A crash before the swap: the old save is untouched.
* A crash after the swap: the new save is complete.
* At no moment does a half-written save exist under the real name.

---

## Clear error semantics

* > After an error, the caller must know: **did it happen? Can I retry?**
* Three levels of promise an operation can make:
* > **Never fails** — always succeeds, no error possible. The strongest.
* > **Strong** — on failure, **nothing changed**: as if it was never
  > called. Always safe to retry.
* > **Basic** — on failure, some work may already be done, but
  > everything is still **valid and usable**: no broken invariants.
* Anything less — a failure halfway through, invariants broken — is
  not a level, it is a bug

---

## Clear error semantics — in code

```cpp
void badSave() {
    file.write(header);   // succeeds
    file.write(body);     // throws — header is already written,
}                          // the caller can't tell what stuck
```

---

<!-- _class: lead invert -->

# Requirements & Communicating Designs

### From a vague ask to a design others can follow

---

## The path from idea to design

* <svg class="uml" width="700" height="170" viewBox="0 0 760 185"><g font-family="Helvetica, Arial, sans-serif" font-size="15"><text x="76" y="30" text-anchor="middle" fill="#5cc8ea" font-size="12" font-weight="700" letter-spacing="2">STEP 1</text><rect x="10" y="45" width="132" height="60" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="76" y="81" text-anchor="middle" fill="#eaf4fb" font-weight="700">Requirements</text><text x="76" y="130" text-anchor="middle" fill="#a9c8de" font-size="12">what &amp; how well</text><line x1="142" y1="75" x2="158" y2="75" stroke="#6fa8cc" stroke-width="1.6"/><polygon points="160,75 151,70 151,80" fill="#6fa8cc"/><text x="226" y="30" text-anchor="middle" fill="#5cc8ea" font-size="12" font-weight="700" letter-spacing="2">STEP 2</text><rect x="160" y="45" width="132" height="60" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="226" y="81" text-anchor="middle" fill="#eaf4fb" font-weight="700">Core entities</text><text x="226" y="130" text-anchor="middle" fill="#a9c8de" font-size="12">name the things</text><line x1="292" y1="75" x2="308" y2="75" stroke="#6fa8cc" stroke-width="1.6"/><polygon points="310,75 301,70 301,80" fill="#6fa8cc"/><text x="376" y="30" text-anchor="middle" fill="#5cc8ea" font-size="12" font-weight="700" letter-spacing="2">STEP 3</text><rect x="310" y="45" width="132" height="60" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="376" y="81" text-anchor="middle" fill="#eaf4fb" font-weight="700">Interface</text><text x="376" y="130" text-anchor="middle" fill="#a9c8de" font-size="12">the operations</text><line x1="442" y1="75" x2="458" y2="75" stroke="#6fa8cc" stroke-width="1.6"/><polygon points="460,75 451,70 451,80" fill="#6fa8cc"/><text x="526" y="30" text-anchor="middle" fill="#5cc8ea" font-size="12" font-weight="700" letter-spacing="2">STEP 4</text><rect x="460" y="45" width="132" height="60" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="526" y="74" text-anchor="middle" fill="#eaf4fb" font-weight="700">High-level</text><text x="526" y="93" text-anchor="middle" fill="#eaf4fb" font-weight="700">design</text><text x="526" y="130" text-anchor="middle" fill="#a9c8de" font-size="12">boxes &amp; arrows</text><line x1="592" y1="75" x2="608" y2="75" stroke="#6fa8cc" stroke-width="1.6"/><polygon points="610,75 601,70 601,80" fill="#6fa8cc"/><text x="676" y="30" text-anchor="middle" fill="#5cc8ea" font-size="12" font-weight="700" letter-spacing="2">STEP 5</text><rect x="610" y="45" width="132" height="60" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="676" y="81" text-anchor="middle" fill="#eaf4fb" font-weight="700">Deep dives</text><text x="676" y="130" text-anchor="middle" fill="#a9c8de" font-size="12">the hard parts</text></g></svg>

* Five steps, always in this order — each step's output feeds the next.

---

## Step 1 — Requirements

* > **Functional** — what the system does.
* > **Non-functional** — how well it must do it: how fast it answers,
  > how many users it can serve, whether it may ever lose data.

---

## Step 1 — the questions that do the work

* Who uses it, and how often?
* How big does the data get?
* What happens when something fails?
* What must **never** happen?
* *Your team project will start from exactly this list.*

---

<span class="tag">Example — link shortener</span>

## Step 1 — answered for the link shortener

* **Does**: shorten a URL (custom name, expiry — optional); redirect
  to the original
* **Not now**: accounts, click statistics
* **Numbers**: 1B links · 100M users a day · ~1000 reads per write
* **Must**: one code → one URL · redirect under 100 ms · almost
  never down

---

## Step 2 — Core entities

The main "things" the system keeps and moves.

* Name the nouns first — names only: no tables, no fields, no database yet
* If you can't name the things, you don't understand the system yet

<span class="tag">Example — link shortener</span>

* <svg class="uml" width="700" height="150" viewBox="0 0 760 165"><g font-family="Helvetica, Arial, sans-serif" font-size="15"><rect x="20" y="35" width="190" height="70" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="115" y="65" text-anchor="middle" fill="#eaf4fb" font-weight="700">User</text><text x="115" y="88" text-anchor="middle" fill="#eaf4fb" font-size="12" opacity="0.75">who creates links</text><line x1="210" y1="70" x2="276" y2="70" stroke="#6fa8cc" stroke-width="1.6"/><polygon points="285,70 276,65 276,75" fill="#6fa8cc"/><text x="247" y="58" text-anchor="middle" fill="#a9c8de" font-size="12">creates</text><rect x="285" y="35" width="190" height="70" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="380" y="65" text-anchor="middle" fill="#eaf4fb" font-weight="700">Short URL</text><text x="380" y="88" text-anchor="middle" fill="#eaf4fb" font-size="12" opacity="0.75">the code · expiry</text><line x1="475" y1="70" x2="541" y2="70" stroke="#6fa8cc" stroke-width="1.6"/><polygon points="550,70 541,65 541,75" fill="#6fa8cc"/><text x="508" y="56" text-anchor="middle" fill="#a9c8de" font-size="11">points to</text><rect x="550" y="35" width="190" height="70" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="645" y="65" text-anchor="middle" fill="#eaf4fb" font-weight="700">Original URL</text><text x="645" y="88" text-anchor="middle" fill="#eaf4fb" font-size="12" opacity="0.75">the long link</text></g></svg>

---

## Step 3 — Interface

How the outside world talks to the system.

* Write each operation: its name, what goes in, what comes out
* Every operation is a **contract** — the promises from the start of
  today apply here, word for word

<span class="tag">Example — link shortener</span>

* `POST /urls { long_url, alias?, expires? } → short_url`
* `GET /{code}` → `302`, redirect to the original

---

## Step 4 — High-level design

The first full picture of the system.

* Boxes and arrows: one box per job, one arrow per conversation
* Follow a single request through the whole picture, end to end
* A sequence diagram shows who calls whom, in what order; use UML
  when the drawing must be precise

<span class="tag">Example — link shortener</span>

* <svg class="uml" width="640" height="110" viewBox="0 0 700 120"><g font-family="Helvetica, Arial, sans-serif" font-size="15"><rect x="20" y="25" width="170" height="60" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="105" y="60" text-anchor="middle" fill="#eaf4fb" font-weight="700">Browser</text><line x1="199" y1="55" x2="246" y2="55" stroke="#6fa8cc" stroke-width="1.6"/><polygon points="255,55 246,50 246,60" fill="#6fa8cc"/><polygon points="190,55 199,50 199,60" fill="#6fa8cc"/><rect x="255" y="25" width="170" height="60" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="340" y="60" text-anchor="middle" fill="#eaf4fb" font-weight="700">Web service</text><line x1="434" y1="55" x2="481" y2="55" stroke="#6fa8cc" stroke-width="1.6"/><polygon points="490,55 481,50 481,60" fill="#6fa8cc"/><polygon points="425,55 434,50 434,60" fill="#6fa8cc"/><rect x="490" y="25" width="170" height="60" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="575" y="60" text-anchor="middle" fill="#eaf4fb" font-weight="700">Database</text></g></svg>

* Three boxes are enough — one box per job

---

## Step 5 — Deep dives

Zoom into the hard parts.

* The hard parts are usually the non-functional requirements
* What breaks first when traffic grows? What happens on a crash,
  halfway through a write?
* This is where today's guarantees — idempotency, atomicity — get
  designed in

---

<span class="tag">Example — link shortener</span>

## The finished design

* <svg class="uml" width="700" height="300" viewBox="0 0 760 325"><g font-family="Helvetica, Arial, sans-serif" font-size="14"><rect x="10" y="130" width="120" height="60" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="70" y="165" text-anchor="middle" fill="#eaf4fb" font-weight="700">Browser</text><line x1="139" y1="160" x2="171" y2="160" stroke="#6fa8cc" stroke-width="1.6"/><polygon points="180,160 171,155 171,165" fill="#6fa8cc"/><polygon points="130,160 139,155 139,165" fill="#6fa8cc"/><rect x="180" y="130" width="130" height="60" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="245" y="155" text-anchor="middle" fill="#eaf4fb" font-weight="700">Load balancer</text><text x="245" y="175" text-anchor="middle" fill="#eaf4fb" font-size="11" opacity="0.75">splits the traffic</text><line x1="310" y1="145" x2="380" y2="85" stroke="#6fa8cc" stroke-width="1.4"/><line x1="310" y1="175" x2="380" y2="235" stroke="#6fa8cc" stroke-width="1.4"/><rect x="380" y="55" width="140" height="60" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="450" y="80" text-anchor="middle" fill="#eaf4fb" font-weight="700">Read service</text><text x="450" y="100" text-anchor="middle" fill="#eaf4fb" font-size="11" opacity="0.75">GET /{code}</text><rect x="380" y="205" width="140" height="60" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="450" y="230" text-anchor="middle" fill="#eaf4fb" font-weight="700">Write service</text><text x="450" y="250" text-anchor="middle" fill="#eaf4fb" font-size="11" opacity="0.75">POST /urls — creates</text><line x1="520" y1="70" x2="600" y2="45" stroke="#6fa8cc" stroke-width="1.4"/><line x1="520" y1="100" x2="600" y2="140" stroke="#6fa8cc" stroke-width="1.4"/><line x1="520" y1="220" x2="600" y2="160" stroke="#6fa8cc" stroke-width="1.4"/><line x1="520" y1="250" x2="600" y2="270" stroke="#6fa8cc" stroke-width="1.4"/><rect x="600" y="15" width="145" height="55" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="672" y="38" text-anchor="middle" fill="#eaf4fb" font-weight="700">Cache</text><text x="672" y="57" text-anchor="middle" fill="#eaf4fb" font-size="11" opacity="0.75">hot codes, in memory</text><rect x="600" y="120" width="145" height="65" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="672" y="147" text-anchor="middle" fill="#eaf4fb" font-weight="700">Database</text><text x="672" y="167" text-anchor="middle" fill="#eaf4fb" font-size="11" opacity="0.75">all links · plus copies</text><rect x="600" y="245" width="145" height="55" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="672" y="268" text-anchor="middle" fill="#eaf4fb" font-weight="700">Counter</text><text x="672" y="287" text-anchor="middle" fill="#eaf4fb" font-size="11" opacity="0.75">one global number</text></g></svg>
---

<!-- _class: lead invert -->

# Concurrency

### What is actually running, right now, on your machine?

---

## What is an operating system?

* > An **operating system (OS)** is the software that manages hardware
  > and gives every program a controlled, shared way to use it.
* Manages the **CPU** (who runs, when), **memory** (who owns what),
  storage, devices, network
* Every program talks to hardware **through** the OS — never directly
* Linux, Windows, macOS — same job, different implementation

---

## A program vs. a process

* > A **program** is a passive file on disk: instructions + data.
  > A **process** is that program, **loaded and running**.
* The same program can run as many independent processes at once
  (two terminals, same shell binary — two separate processes)
* The OS gives each process its own memory, its own file handles,
  its own identity (a PID)

---

## A process's memory image

* <svg class="uml" width="640" height="360" viewBox="0 0 700 380"><g font-family="Helvetica, Arial, sans-serif" font-size="17"><rect x="220" y="20" width="260" height="80" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="350" y="55" text-anchor="middle" fill="#eaf4fb" font-weight="700">Stack</text><text x="350" y="78" text-anchor="middle" fill="#eaf4fb" font-size="14" opacity="0.75">locals, call frames — grows ↓</text><rect x="220" y="100" width="260" height="45" fill="none" stroke="#6fa8cc" stroke-width="1.2" stroke-dasharray="6 5" opacity="0.6"/><text x="350" y="127" text-anchor="middle" fill="#eaf4fb" font-size="13" opacity="0.6">unused (free space)</text><rect x="220" y="145" width="260" height="85" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="350" y="182" text-anchor="middle" fill="#eaf4fb" font-weight="700">Heap</text><text x="350" y="205" text-anchor="middle" fill="#eaf4fb" font-size="14" opacity="0.75">malloc / new — grows ↑</text><rect x="220" y="230" width="260" height="60" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="350" y="265" text-anchor="middle" fill="#eaf4fb" font-weight="700">Data / BSS</text><rect x="220" y="290" width="260" height="60" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="350" y="325" text-anchor="middle" fill="#eaf4fb" font-weight="700">Code / Text</text><text x="350" y="365" text-anchor="middle" fill="#eaf4fb" font-size="13" font-style="italic" opacity="0.7">one process — one address space</text></g></svg>

---

## What's inside the image

* **Code / Text** — the instructions, read-only
* **Data / BSS** — global and static variables
* **Heap** — memory asked for while running, grows on demand
* **Stack** — function calls and locals, grows and shrinks automatically

---

## Same layout, in code

```c
int global_counter;              /* Data / BSS */

void handle_request(void) {
    int local_id = 42;           /* Stack */
    char *buf = malloc(1024);    /* Heap */
}
```

---

## How a program actually runs

* > The OS **loads** the program's instructions into memory, points
  > the CPU at the entry point, and gets out of the way — until it
  > needs the CPU back.
* A CPU core is a fast **fetch → decode → execute** loop
* One core runs **one instruction stream at a time**
* The OS's **scheduler** decides which process runs next, and for
  how long — a **context switch** saves one process's registers and
  loads another's
* On one core, things only *look* simultaneous — the switching is
  just that fast; with several cores, some of it is real

---

## Collaboration with the OS — system calls

> A process runs in restricted **user mode** — it cannot touch
> hardware, or another process's memory, directly.

* To do anything that matters — read a file, send a packet, get more
  memory, start another process — a program asks the OS with a
  **system call**
* The CPU switches to privileged **kernel mode**, the OS does the
  work, control returns to user mode
* Familiar syscalls: `read` / `write` / `open` (files), `mmap`
  (memory), `fork` / `exec` (processes)
* Nearly everything a program does beyond pure computation crosses
  this boundary — including, later, thread creation and locks

---

## System calls — files

```c
int fd = open("data.txt", O_RDONLY);   // syscall: ask the OS to open
read(fd, buf, sizeof(buf));            // syscall: ask the OS to read
close(fd);                             // syscall: ask the OS to close
```

---

## System calls — processes

```c
pid_t pid = fork();        // syscall: ask the OS for a new process
if (pid == 0) {
    exec("./worker", ...); // syscall: replace this process's image
} else {
    wait(NULL);             // syscall: block until the child exits
}
```

---

## System calls — memory

```c
void *buf = mmap(NULL, 4096, PROT_READ | PROT_WRITE,
                  MAP_PRIVATE | MAP_ANONYMOUS, -1, 0);
// syscall: ask the OS for more address space
munmap(buf, 4096);
```

---

## Recap: a process, revisited

A **process** = a running program: its own address space, its own
file handles, at least one flow of execution through the memory
image above.

* What if one process needs to do more than one thing **at once**,
  without paying for a whole second process?

---

## A thread

> A **thread** is an independent flow of execution **inside** a
> process.

* <svg class="uml" width="700" height="380" viewBox="0 0 760 400"><g font-family="Helvetica, Arial, sans-serif" font-size="16"><rect x="140" y="20" width="200" height="70" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="240" y="50" text-anchor="middle" fill="#eaf4fb" font-weight="700">Stack — thread A</text><text x="240" y="72" text-anchor="middle" fill="#eaf4fb" font-size="13" opacity="0.75">own locals, own IP</text><rect x="420" y="20" width="200" height="70" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="520" y="50" text-anchor="middle" fill="#eaf4fb" font-weight="700">Stack — thread B</text><text x="520" y="72" text-anchor="middle" fill="#eaf4fb" font-size="13" opacity="0.75">own locals, own IP</text><line x1="240" y1="90" x2="240" y2="150" stroke="#6fa8cc" stroke-width="1.6"/><line x1="520" y1="90" x2="520" y2="150" stroke="#6fa8cc" stroke-width="1.6"/><rect x="120" y="150" width="520" height="90" fill="#14385a" stroke="#7ed0a3" stroke-width="1.8"/><text x="380" y="188" text-anchor="middle" fill="#7ed0a3" font-weight="700">Heap — shared</text><text x="380" y="211" text-anchor="middle" fill="#eaf4fb" font-size="14" opacity="0.85">both threads read &amp; write the same objects</text><rect x="120" y="250" width="520" height="50" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="380" y="281" text-anchor="middle" fill="#eaf4fb" font-weight="700">Data / BSS — shared</text><rect x="120" y="310" width="520" height="50" fill="#14385a" stroke="#6fa8cc" stroke-width="1.6"/><text x="380" y="341" text-anchor="middle" fill="#eaf4fb" font-weight="700">Code / Text — shared</text></g></svg>

---

## What that picture means

* Multiple threads = multiple stacks, multiple instruction pointers
* **One shared heap, one shared set of globals**

---

## A thread — in code

```cpp
int counter = 0;                 // shared between threads

void increment() { counter++; }

std::thread t1(increment);
std::thread t2(increment);
t1.join();
t2.join();
```
