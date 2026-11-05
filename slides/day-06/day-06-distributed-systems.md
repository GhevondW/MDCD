---
marp: true
theme: gaia
class: invert
paginate: true
footer: 'Distributed Systems'
style: |
  @import url('https://fonts.bunny.net/css?family=ibm-plex-sans:400,500,600,700|ibm-plex-sans-condensed:600,700|ibm-plex-mono:400,500,600&display=swap');

  /* Same blueprint palette as Days 2–5. */
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

  .uml { display: block; margin: 0 auto; overflow: visible; }

  /* Step-by-step reveals: list items written with `*` appear one
     keypress at a time. Items that carry a quote panel or a diagram
     shouldn't also show a list bullet. */
  li:has(> blockquote), li:has(svg) { list-style: none; }

  /* Sequence diagrams: one message per keypress. Rows sit flush under
     the static header row, so the lanes read as continuous lines. */
  ul:has(> li > svg.tl) { padding-left: 0; }
  li:has(> svg.tl) { margin: 0; line-height: 0; }
  .sp { height: 6px; }

  .cols { display: flex; gap: 28px; align-items: center; }
  .cols > div { flex: 1; min-width: 0; }

  /* ---------- Diagrams: inline SVG, animated with CSS only ---------- */
  .d text { font-family: 'IBM Plex Sans', Helvetica, Arial, sans-serif; font-size: 15px; fill: #eaf4fb; text-anchor: middle; }
  .d .t { font-weight: 700; }
  .d .s { font-size: 12px; fill: #a9c8de; }
  .d .m { font-family: 'IBM Plex Mono', Menlo, monospace; font-size: 14px; }
  .d .cap { font-size: 12px; font-weight: 700; letter-spacing: 2px; fill: #5cc8ea; }
  .d .st { text-anchor: start; }
  .d .gd { fill: #7ed0a3; }
  .d .bd { fill: #f0a06e; }
  .d .ac { fill: #5cc8ea; }
  .d .bx { fill: #14385a; stroke: #6fa8cc; stroke-width: 1.6; }
  .d .bx.g { stroke: #7ed0a3; }
  .d .bx.b { stroke: #f0a06e; }
  .d .bx.a { stroke: #5cc8ea; stroke-width: 2; }
  .d .bx.dim { opacity: .45; }
  .d .grp { fill: none; stroke: #6fa8cc; stroke-width: 1.2; stroke-dasharray: 6 5; opacity: .6; }
  .d .ln { stroke: #6fa8cc; stroke-width: 1.6; fill: none; }
  .d .ln.b { stroke: #f0a06e; }
  .d .ln.g { stroke: #7ed0a3; }
  .d .ln.a { stroke: #5cc8ea; }
  .d .ln.x { stroke: #f0a06e; stroke-width: 3; }
  .d .hd { fill: #6fa8cc; }
  .d .hd.b { fill: #f0a06e; }
  .d .hd.g { fill: #7ed0a3; }
  .d .hd.a { fill: #5cc8ea; }
  .d .lane { stroke: #6fa8cc; stroke-width: 1; stroke-dasharray: 3 4; opacity: .45; }
  .d .pkt { fill: #5cc8ea; }
  .d .pkt.r { fill: #7ed0a3; }
  .d .pkt.b { fill: #f0a06e; }
  .d .fly, .d .hop, .d .seg, .d .lost, .d .once { opacity: 0; }

  /* Looping messages (whole slide): a packet moves by (--dx, --dy) at
     the start of every --t cycle, after an initial --d delay. */
  @media screen {
    .d .fly { animation: fly var(--t, 3s) linear var(--d, 0s) infinite both; }
    .d .hop { animation: hop var(--t, 6s) linear var(--d, 0s) infinite both; }
    .d .seg { animation: seg var(--t, 12s) linear var(--d, 0s) infinite both; }
    .d .lost { animation: lost var(--t, 3s) linear var(--d, 0s) infinite both; }
    .d .boom { animation: boom var(--t, 3s) linear var(--d, 0s) infinite both; }
    .d .pulse { animation: pulse 1.4s ease-in-out infinite; }
    .d .pop { animation: late .4s ease-out var(--d, 0s) both; }
    .d .spin { transform-box: view-box; transform-origin: var(--ox) var(--oy); animation: spin var(--t, 6s) linear infinite; }

    /* On keypress (sequence rows): the arrow draws itself, a packet
       runs along it, then the parts marked .late appear. */
    [data-bespoke-marp-fragment=active] .d .draw { animation: draw .6s ease-out var(--d, 0s) both; }
    [data-bespoke-marp-fragment=active] .d .once { animation: once .6s ease-out var(--d, 0s) both; }
    [data-bespoke-marp-fragment=active] .d .late { animation: late .3s ease-out calc(var(--d, 0s) + .55s) both; }
  }
  .d .draw { stroke-dasharray: 1; }
  @keyframes fly {
    0% { transform: translate(0, 0); opacity: 0; }
    4% { opacity: 1; }
    40% { transform: translate(var(--dx), var(--dy)); opacity: 1; }
    46%, 100% { transform: translate(var(--dx), var(--dy)); opacity: 0; }
  }
  @keyframes hop {
    0% { transform: translate(0, 0); opacity: 0; }
    2% { opacity: 1; }
    14% { transform: translate(var(--dx), var(--dy)); opacity: 1; }
    16%, 100% { transform: translate(var(--dx), var(--dy)); opacity: 0; }
  }
  @keyframes seg {
    0% { transform: translate(0, 0); opacity: 0; }
    1% { opacity: 1; }
    6% { transform: translate(var(--dx), var(--dy)); opacity: 1; }
    7%, 100% { transform: translate(var(--dx), var(--dy)); opacity: 0; }
  }
  @keyframes lost {
    0% { transform: translate(0, 0); opacity: 0; }
    4% { opacity: 1; }
    20% { transform: translate(var(--dx), var(--dy)); opacity: 1; }
    24%, 100% { transform: translate(var(--dx), var(--dy)); opacity: 0; }
  }
  @keyframes boom { 0%, 19% { opacity: 0; } 22%, 55% { opacity: 1; } 62%, 100% { opacity: 0; } }
  @keyframes pulse { 50% { opacity: .3; } }
  @keyframes spin { to { transform: rotate(360deg); } }
  @keyframes draw { from { stroke-dashoffset: 1; } to { stroke-dashoffset: 0; } }
  @keyframes once {
    0% { transform: translate(0, 0); opacity: 0; }
    8%, 85% { opacity: 1; }
    100% { transform: translate(var(--dx), var(--dy)); opacity: 0; }
  }
  @keyframes late { from { opacity: 0; } to { opacity: 1; } }
  .d .en { text-anchor: end; }

  /* Where we are in the stack: a small map in the top-right corner. */
  .mini { position: absolute; top: 44px; right: 70px; }

---

<!-- _class: lead invert -->

# Distributed Systems

### Modeling, Design & Collaborative Development

---

<span class="tag">Days 3–5</span>

## So far: one machine

<svg class="uml d " width="1000" height="246" viewBox="0 0 1000 246"><rect class="grp" x="30" y="8" width="420" height="200" rx="4"/><text class="cap" x="240" y="30">ONE MACHINE</text><path class="ln " d="M111.1 91.2L159.7 132.8"/><polygon class="hd " points="165.8,138 155.7,135.9 162.1,128.4"/><polygon class="hd " points="105,86 115.1,88.1 108.6,95.6"/><rect class="bx " x="55" y="48" width="100" height="38" rx="3"/><text class="t " x="105" y="72">thread</text><path class="ln " d="M240 94L240 130"/><polygon class="hd " points="240,138 235.1,129 244.9,129"/><polygon class="hd " points="240,86 244.9,95 235.1,95"/><rect class="bx " x="190" y="48" width="100" height="38" rx="3"/><text class="t " x="240" y="72">thread</text><path class="ln " d="M368.9 91.2L320.3 132.8"/><polygon class="hd " points="314.2,138 317.9,128.4 324.3,135.9"/><polygon class="hd " points="375,86 371.4,95.6 364.9,88.1"/><rect class="bx " x="325" y="48" width="100" height="38" rx="3"/><text class="t " x="375" y="72">thread</text><rect class="bx a" x="110" y="138" width="260" height="46" rx="3"/><text class="t " x="240" y="166">shared memory</text><text class="s" x="240" y="236" style="font-size:14px">threads share memory → locks, atomics</text><path class="ln a" d="M468 110L524 110"/><polygon class="hd a" points="532,110 523,115 523,105"/><text class="cap" x="500" y="98">TODAY</text><path class="ln " d="M630 47L910 47"/><path class="ln " d="M910 47L770 177"/><path class="ln " d="M770 177L630 47"/><circle class="pkt  fly" cx="630" cy="47" r="6" style="--dx:280px;--dy:0px;--t:3s;--d:0s"/><circle class="pkt  fly" cx="910" cy="47" r="6" style="--dx:-140px;--dy:130px;--t:3s;--d:1s"/><circle class="pkt  fly" cx="770" cy="177" r="6" style="--dx:-140px;--dy:-130px;--t:3s;--d:2s"/><rect class="bx " x="560" y="22" width="140" height="50" rx="3"/><text class="t " x="630" y="45">machine</text><text class="s" x="630" y="62">own memory</text><rect class="bx " x="840" y="22" width="140" height="50" rx="3"/><text class="t " x="910" y="45">machine</text><text class="s" x="910" y="62">own memory</text><rect class="bx " x="700" y="152" width="140" height="50" rx="3"/><text class="t " x="770" y="175">machine</text><text class="s" x="770" y="192">own memory</text><text class="s" x="770" y="236" style="font-size:14px">machines share nothing → only messages</text></svg>

- Days 3–5: many threads, **one** memory — guarded by locks and atomics
- Day 5's actors and channels already shared nothing: they sent messages
- Today the threads are on **different computers**. Only the messages are left

---

## What is a distributed system?

- > A **distributed system** is a group of computers that share **no memory**, talk only by **sending messages** over a network, and work together so that their users see **one** system.

<svg class="uml d " width="1000" height="204" viewBox="0 0 1000 204"><rect class="grp" x="190" y="4" width="800" height="196" rx="6"/><text class="cap" x="590" y="24">ONE SYSTEM — TO ITS USERS</text><path class="ln " d="M315 62L545 112"/><path class="ln " d="M315 162L545 112"/><path class="ln " d="M315 62L775 62"/><path class="ln " d="M315 162L775 162"/><path class="ln " d="M775 62L775 162"/><path class="ln " d="M120 50L232 61.2"/><polygon class="hd " points="240,62 230.6,66 231.5,56.2"/><circle class="pkt  fly" cx="120" cy="50" r="6" style="--dx:120px;--dy:12px;--t:3s;--d:0.0s"/><path class="ln " d="M120 112L232.6 158.9"/><polygon class="hd " points="240,162 229.8,163.1 233.6,154"/><circle class="pkt  fly" cx="120" cy="112" r="6" style="--dx:120px;--dy:50px;--t:3s;--d:1.0s"/><path class="ln " d="M120 174L232 162.8"/><polygon class="hd " points="240,162 231.5,167.8 230.6,158"/><circle class="pkt  fly" cx="120" cy="174" r="6" style="--dx:120px;--dy:-12px;--t:3s;--d:2.0s"/><circle class="pkt  fly" cx="315" cy="62" r="6" style="--dx:230px;--dy:50px;--t:3s;--d:0.5s"/><circle class="pkt  fly" cx="315" cy="162" r="6" style="--dx:460px;--dy:0px;--t:3s;--d:1.5s"/><circle class="pkt r fly" cx="775" cy="62" r="6" style="--dx:0px;--dy:100px;--t:3s;--d:2.2s"/><circle class="pkt  fly" cx="315" cy="62" r="6" style="--dx:460px;--dy:0px;--t:3s;--d:1.1s"/><rect class="bx " x="20" y="30" width="100" height="40" rx="3"/><text class="t " x="70" y="55">user</text><rect class="bx " x="20" y="92" width="100" height="40" rx="3"/><text class="t " x="70" y="117">user</text><rect class="bx " x="20" y="154" width="100" height="40" rx="3"/><text class="t " x="70" y="179">user</text><rect class="bx " x="240" y="40" width="150" height="44" rx="3"/><text class="t " x="315" y="67">web server</text><rect class="bx " x="240" y="140" width="150" height="44" rx="3"/><text class="t " x="315" y="167">web server</text><rect class="bx " x="470" y="90" width="150" height="44" rx="3"/><text class="t " x="545" y="117">cache</text><rect class="bx " x="700" y="40" width="150" height="44" rx="3"/><text class="t " x="775" y="67">database</text><rect class="bx " x="700" y="140" width="150" height="44" rx="3"/><text class="t " x="775" y="167">database copy</text></svg>

- You use one every day: web search, your bank, WhatsApp, an online game


---

## Why not one big computer?

- **Too much work** — Day 2's link shortener: 1 billion links, 100 million users a day. More than one machine can store and serve
- **Machines die** — disks fail, power goes out, a bad deploy crashes it. One machine: everything stops. Copies: the others take over
- **Distance** — light needs time. From Yerevan, a round trip to Frankfurt takes at least 30 ms; to California, at least 110 ms
- > If one machine is enough, use one machine. Distribution is not free — today is about its price.

<!-- 30 ms / 110 ms: great-circle distance (2,984 km / 11,246 km), there and back, at ~200 km per ms in fibre. Real routes are longer. -->

---

## Threads vs. machines: what changes

| | Threads (Days 3–5) | Machines (today) |
|---|---|---|
| They talk by | reading and writing shared memory | sending messages |
| One "talk" takes | ~100 ns | ~0.5 ms nearby · ~150 ms across the world |
| A crash | kills the process — every thread together | kills **one** machine — the rest keep running |
| Lost on the way? | never | messages get lost |
| The clock | one | one per machine — they drift apart |

---

## Why can't machines agree on the time?

<svg class="uml d " width="1000" height="172" viewBox="0 0 1000 172"><rect class="bx " x="20" y="46" width="170" height="64" rx="3"/><text class="t " x="105" y="76">time server</text><text class="s" x="105" y="93">12:00:00.000</text><path class="ln a" d="M190 78L310 78"/><polygon class="hd a" points="318,78 309,83 309,73"/><circle class="pkt  fly" cx="190" cy="78" r="6" style="--dx:122px;--dy:0px;--t:3s;--d:0s"/><text class="m ac" x="254" y="66" style="font-size:12px">it's 12:00:00</text><text class="s bd" x="254" y="100">arrives ? ms later</text><rect class="grp" x="330" y="4" width="660" height="164" rx="6"/><text class="cap" x="660" y="22">EACH MACHINE COUNTS TIME WITH ITS OWN QUARTZ CRYSTAL</text><circle class="bx" cx="440" cy="74" r="38"/><path class="ln " d="M440 45L440 38"/><path class="ln " d="M456 46.3L458 42.8"/><path class="ln " d="M467.7 58L471.2 56"/><path class="ln " d="M469 74L476 74"/><path class="ln " d="M467.7 90L471.2 92"/><path class="ln " d="M456 101.7L458 105.2"/><path class="ln " d="M440 103L440 110"/><path class="ln " d="M424 101.7L422 105.2"/><path class="ln " d="M412.3 90L408.8 92"/><path class="ln " d="M411 74L404 74"/><path class="ln " d="M412.3 58L408.8 56"/><path class="ln " d="M424 46.3L422 42.8"/><path class="ln" style="stroke-width:4;stroke-linecap:round" d="M440 74L424.4 63.1"/><path class="ln" style="stroke-width:3;stroke-linecap:round" d="M440 74L464.4 59.9"/><path class="ln a spin" style="stroke-width:2;--ox:440px;--oy:74px;--t:5.6s" d="M440 82L440 43"/><circle cx="440" cy="74" r="3.5" style="fill:#5cc8ea"/><text class="t" x="440" y="134">machine A</text><text class="s" x="440" y="151">runs fast</text><circle class="bx" cx="660" cy="74" r="38"/><path class="ln " d="M660 45L660 38"/><path class="ln " d="M676 46.3L678 42.8"/><path class="ln " d="M687.7 58L691.2 56"/><path class="ln " d="M689 74L696 74"/><path class="ln " d="M687.7 90L691.2 92"/><path class="ln " d="M676 101.7L678 105.2"/><path class="ln " d="M660 103L660 110"/><path class="ln " d="M644 101.7L642 105.2"/><path class="ln " d="M632.3 90L628.8 92"/><path class="ln " d="M631 74L624 74"/><path class="ln " d="M632.3 58L628.8 56"/><path class="ln " d="M644 46.3L642 42.8"/><path class="ln" style="stroke-width:4;stroke-linecap:round" d="M660 74L644.4 63.1"/><path class="ln" style="stroke-width:3;stroke-linecap:round" d="M660 74L684.4 59.9"/><path class="ln a spin" style="stroke-width:2;--ox:660px;--oy:74px;--t:6.0s" d="M660 82L660 43"/><circle cx="660" cy="74" r="3.5" style="fill:#5cc8ea"/><text class="t" x="660" y="134">machine B</text><text class="s" x="660" y="151">on time — for now</text><circle class="bx" cx="880" cy="74" r="38"/><path class="ln " d="M880 45L880 38"/><path class="ln " d="M896 46.3L898 42.8"/><path class="ln " d="M907.7 58L911.2 56"/><path class="ln " d="M909 74L916 74"/><path class="ln " d="M907.7 90L911.2 92"/><path class="ln " d="M896 101.7L898 105.2"/><path class="ln " d="M880 103L880 110"/><path class="ln " d="M864 101.7L862 105.2"/><path class="ln " d="M852.3 90L848.8 92"/><path class="ln " d="M851 74L844 74"/><path class="ln " d="M852.3 58L848.8 56"/><path class="ln " d="M864 46.3L862 42.8"/><path class="ln" style="stroke-width:4;stroke-linecap:round" d="M880 74L864.4 63.1"/><path class="ln" style="stroke-width:3;stroke-linecap:round" d="M880 74L904.4 59.9"/><path class="ln a spin" style="stroke-width:2;--ox:880px;--oy:74px;--t:6.4s" d="M880 82L880 43"/><circle cx="880" cy="74" r="3.5" style="fill:#5cc8ea"/><text class="t" x="880" y="134">machine C</text><text class="s" x="880" y="151">runs slow</text></svg>

- **Drift**: a quartz crystal counts the time; no two tick at exactly the same speed, and heat and age change them. Google's Spanner plans for up to 200 µs a second — 17 s a day, uncorrected
- **Syncing is a message**: "it's 12:00:00" arrives late, by an amount you can't know — you only see the round trip. NTP: ~1 ms off on a LAN, tens of ms over the internet
- > Clocks can be **close**, never **equal** — and a corrected clock can jump back. Don't order events on different machines by timestamps.

<!-- Spanner (OSDI 2012): worst-case drift 200 µs/s, polled every 30 s → 0–6 ms; 200 µs/s × 86,400 s = 17.3 s. NTP numbers: Wikipedia "Network Time Protocol". Order events with sequence numbers or logical clocks (Lamport, 1978); measure durations with a monotonic clock. -->

---

<!-- _class: lead invert -->

# The Network

### How a message gets from here to there

---

## Sending one message: four problems

Your laptop in Yerevan sends `GET /abc123` to a server in Frankfurt. Someone has to solve:

1. **The next hop** — turn bits into signals on one wire or radio, and reach the next device
2. **Across networks** — find one machine among billions, through thousands of networks
3. **The right program** — hand the data to one program on that machine, and cope with lost or mixed-up packets
4. **The meaning** — both sides must agree what the bytes mean

- > One program that solves all four would be enormous — and every application would write it again.

---

## The solution: a stack of layers

<svg class="uml d " width="920" height="253.9" viewBox="0 0 1000 276"><rect x="8" y="4" width="300" height="268" rx="8" style="fill:rgba(92,200,234,.12);stroke:rgba(92,200,234,.9);stroke-width:1.5"/><rect x="692" y="4" width="300" height="268" rx="8" style="fill:rgba(92,200,234,.12);stroke:rgba(92,200,234,.9);stroke-width:1.5"/><rect x="318" y="112" width="364" height="160" rx="8" style="fill:rgba(126,208,163,.12);stroke:rgba(126,208,163,.9);stroke-width:1.5"/><text class="cap" x="158" y="22">IN YOUR COMPUTER</text><text class="cap" x="842" y="22">IN THE SERVER'S COMPUTER</text><text class="cap" x="500" y="130" style="fill:#7ed0a3">IN THE NETWORK</text><path class="ln " d="M290 226L400 226"/><path class="ln " d="M600 226L710 226"/><text class="s" x="345" y="250">Wi-Fi</text><text class="s" x="655" y="250">fibre</text><rect class="bx " x="70" y="34" width="220" height="48" rx="3"/><text class="t " x="180" y="56">4 · Application</text><text class="s" x="180" y="73">your program · HTTP, DNS</text><rect class="bx " x="70" y="90" width="220" height="48" rx="3"/><text class="t " x="180" y="112">3 · Transport</text><text class="s" x="180" y="129">the OS · TCP, UDP</text><rect class="bx " x="70" y="146" width="220" height="48" rx="3"/><text class="t " x="180" y="168">2 · Internet</text><text class="s" x="180" y="185">the OS · IP</text><rect class="bx " x="70" y="202" width="220" height="48" rx="3"/><text class="t " x="180" y="224">1 · Link</text><text class="s" x="180" y="241">network card · Wi-Fi, Ethernet</text><rect class="bx g" x="400" y="146" width="200" height="48" rx="3"/><text class="t " x="500" y="168">2 · Internet</text><text class="s" x="500" y="185">router: reads the IP address</text><rect class="bx g" x="400" y="202" width="200" height="48" rx="3"/><text class="t " x="500" y="224">1 · Link</text><text class="s" x="500" y="241">router: Wi-Fi, fibre</text><rect class="bx " x="710" y="34" width="220" height="48" rx="3"/><text class="t " x="820" y="56">4 · Application</text><text class="s" x="820" y="73">your program · HTTP, DNS</text><rect class="bx " x="710" y="90" width="220" height="48" rx="3"/><text class="t " x="820" y="112">3 · Transport</text><text class="s" x="820" y="129">the OS · TCP, UDP</text><rect class="bx " x="710" y="146" width="220" height="48" rx="3"/><text class="t " x="820" y="168">2 · Internet</text><text class="s" x="820" y="185">the OS · IP</text><rect class="bx " x="710" y="202" width="220" height="48" rx="3"/><text class="t " x="820" y="224">1 · Link</text><text class="s" x="820" y="241">network card · Wi-Fi, Ethernet</text><path class="ln a" d="M40 40L40 238"/><polygon class="hd a" points="40,246 35,237 45,237"/><text class="cap" x="40" y="264">send</text><path class="ln a" d="M960 246L960 48"/><polygon class="hd a" points="960,40 965,49 955,49"/><text class="cap" x="960" y="264">receive</text><circle class="pkt  seg" cx="180" cy="58" r="7" style="--dx:0px;--dy:56px;--t:12s;--d:0.0s"/><circle class="pkt  seg" cx="180" cy="114" r="7" style="--dx:0px;--dy:56px;--t:12s;--d:0.72s"/><circle class="pkt  seg" cx="180" cy="170" r="7" style="--dx:0px;--dy:56px;--t:12s;--d:1.44s"/><circle class="pkt  seg" cx="180" cy="226" r="7" style="--dx:320px;--dy:0px;--t:12s;--d:2.16s"/><circle class="pkt  seg" cx="500" cy="226" r="7" style="--dx:0px;--dy:-56px;--t:12s;--d:2.88s"/><circle class="pkt  seg" cx="500" cy="170" r="7" style="--dx:0px;--dy:56px;--t:12s;--d:3.6s"/><circle class="pkt  seg" cx="500" cy="226" r="7" style="--dx:320px;--dy:0px;--t:12s;--d:4.32s"/><circle class="pkt  seg" cx="820" cy="226" r="7" style="--dx:0px;--dy:-56px;--t:12s;--d:5.04s"/><circle class="pkt  seg" cx="820" cy="170" r="7" style="--dx:0px;--dy:-56px;--t:12s;--d:5.76s"/><circle class="pkt  seg" cx="820" cy="114" r="7" style="--dx:0px;--dy:-56px;--t:12s;--d:6.48s"/></svg>

- Each layer solves **one** problem, using only the layer below. Send: **down**; receive: **up**
- **In the computer**: 4 is your program · 3 and 2 the operating system · 1 the network card
- **In the network**: wires, radio, routers — layers 1 and 2 only. A router never looks at TCP or HTTP

<p class="note">The OSI model has 7 layers; the internet is built on these 4 (RFC 1122).</p>

---

## Down the stack: an envelope inside an envelope

<svg class="uml d " width="1000" height="128" viewBox="0 0 1000 128"><g class="pop" style="--d:0.2s"><rect class="bx a" x="360" y="82" width="580" height="40" rx="3"/><text class="t m " x="650" y="107">GET /abc123 HTTP/1.1 …</text></g><g class="pop" style="--d:2.3s"><rect class="bx " x="60" y="82" width="100" height="40" rx="3"/><text class="t " x="110" y="107">Ethernet</text><path class="ln " d="M60 22L940 22"/><path class="ln " d="M60 18L60 26"/><path class="ln " d="M940 18L940 26"/><text class="s st" x="66" y="17">Ethernet frame</text></g><g class="pop" style="--d:1.6s"><rect class="bx " x="160" y="82" width="100" height="40" rx="3"/><text class="t " x="210" y="107">IP</text><path class="ln " d="M160 44L940 44"/><path class="ln " d="M160 40L160 48"/><path class="ln " d="M940 40L940 48"/><text class="s st" x="166" y="39">IP packet</text></g><g class="pop" style="--d:0.9s"><rect class="bx " x="260" y="82" width="100" height="40" rx="3"/><text class="t " x="310" y="107">TCP</text><path class="ln " d="M260 66L940 66"/><path class="ln " d="M260 62L260 70"/><path class="ln " d="M940 62L940 70"/><text class="s st" x="266" y="61">TCP segment</text></g></svg>

| Layer | Puts in its header | For example |
|---|---|---|
| **4 · Application** | the message itself | `GET /abc123 HTTP/1.1` |
| **3 · Transport** | ports: which program, on each side | `51234 → 443` |
| **2 · Internet** | IP addresses: which machine, on each side | `192.168.1.10 → 172.66.147.243` |
| **1 · Link** | MAC addresses: this hop only | laptop → Wi-Fi router |

- The receiver opens the envelopes in reverse: each layer reads its own header and passes the rest up

---

<span class="tag">Layer 1 · Link</span>

<div class="mini"><svg class="uml d " width="152" height="110" viewBox="0 0 152 110"><rect x="1" y="2" width="150" height="23" rx="2" style="fill:#14385a;stroke:#6fa8cc;stroke-width:1;opacity:0.7"/><text class="s" x="76" y="18" style="font-size:12px">4 · Application</text><rect x="1" y="29" width="150" height="23" rx="2" style="fill:#14385a;stroke:#6fa8cc;stroke-width:1;opacity:0.7"/><text class="s" x="76" y="45" style="font-size:12px">3 · Transport</text><rect x="1" y="56" width="150" height="23" rx="2" style="fill:#14385a;stroke:#6fa8cc;stroke-width:1;opacity:0.7"/><text class="s" x="76" y="72" style="font-size:12px">2 · Internet</text><rect x="1" y="83" width="150" height="23" rx="2" style="fill:#5cc8ea;stroke:#6fa8cc;stroke-width:1;opacity:1"/><text class="t" x="76" y="99" style="font-size:12px;fill:#0b2136">1 · Link</text></svg></div>

## The next hop

> **Problem:** the laptop, the phone and the router share one Wi-Fi radio. How does a frame reach the right device?

<svg class="uml d " width="900" height="154.8" viewBox="0 0 1000 172"><rect class="grp" x="10" y="6" width="650" height="160" rx="6"/><text class="cap" x="335" y="24">YOUR LOCAL NETWORK</text><path class="ln " d="M180 58L265 90"/><path class="ln " d="M180 128L265 98"/><path class="ln " d="M415 94L480 94"/><path class="ln " d="M630 94L760 94"/><circle class="pkt  hop" cx="180" cy="58" r="6" style="--dx:85px;--dy:32px;--t:6s;--d:0s"/><circle class="pkt  hop" cx="415" cy="94" r="6" style="--dx:65px;--dy:0px;--t:6s;--d:0.84s"/><circle class="pkt r hop" cx="630" cy="94" r="6" style="--dx:130px;--dy:0px;--t:6s;--d:1.68s"/><text class="s ac" x="448" y="76">to: router</text><text class="s gd" x="695" y="76">a new frame</text><rect class="bx " x="30" y="34" width="150" height="48" rx="3"/><text class="t " x="105" y="56">laptop</text><text class="s" x="105" y="73">a4:83:e7:2b:19:0c</text><rect class="bx " x="30" y="104" width="150" height="48" rx="3"/><text class="t " x="105" y="126">phone</text><text class="s" x="105" y="143">6e:1d:90:c4:52:a8</text><rect class="bx " x="265" y="70" width="150" height="48" rx="3"/><text class="t " x="340" y="99">Wi-Fi / switch</text><rect class="bx a" x="480" y="70" width="150" height="48" rx="3"/><text class="t " x="555" y="92">router</text><text class="s" x="555" y="109">the way out</text><rect class="bx dash" x="760" y="70" width="220" height="48" rx="3" style="stroke-dasharray:6 5"/><text class="s" x="870" y="99" style="font-size:14px">the rest of the internet</text></svg>

- Every network card has a **MAC address** — 48 bits, like `a4:83:e7:2b:19:0c`
- Data travels in **frames** "from MAC A to MAC B" — at most ~1500 bytes of data on Ethernet
- A frame stays in the **local** network. To go further it goes to the **router**; the next link gets a new frame

---

<span class="tag">Layer 2 · Internet</span>

<div class="mini"><svg class="uml d " width="152" height="110" viewBox="0 0 152 110"><rect x="1" y="2" width="150" height="23" rx="2" style="fill:#14385a;stroke:#6fa8cc;stroke-width:1;opacity:0.7"/><text class="s" x="76" y="18" style="font-size:12px">4 · Application</text><rect x="1" y="29" width="150" height="23" rx="2" style="fill:#14385a;stroke:#6fa8cc;stroke-width:1;opacity:0.7"/><text class="s" x="76" y="45" style="font-size:12px">3 · Transport</text><rect x="1" y="56" width="150" height="23" rx="2" style="fill:#5cc8ea;stroke:#6fa8cc;stroke-width:1;opacity:1"/><text class="t" x="76" y="72" style="font-size:12px;fill:#0b2136">2 · Internet</text><rect x="1" y="83" width="150" height="23" rx="2" style="fill:#14385a;stroke:#6fa8cc;stroke-width:1;opacity:0.7"/><text class="s" x="76" y="99" style="font-size:12px">1 · Link</text></svg></div>

## From network to network

> **Problem:** the server is not on my network. Between us: thousands of networks and billions of machines.

<svg class="uml d " width="800" height="177.6" viewBox="0 0 1000 222"><path class="ln " d="M75 102L280 102"/><path class="ln " d="M280 102L490 40"/><path class="ln " d="M280 102L490 164"/><path class="ln " d="M490 40L700 102"/><path class="ln " d="M490 164L700 102"/><path class="ln " d="M700 102L925 102"/><circle class="pkt  hop" cx="75" cy="102" r="6" style="--dx:205px;--dy:0px;--t:6s;--d:0s"/><circle class="pkt  hop" cx="280" cy="102" r="6" style="--dx:210px;--dy:-62px;--t:6s;--d:0.84s"/><circle class="pkt  hop" cx="490" cy="40" r="6" style="--dx:210px;--dy:62px;--t:6s;--d:3.0s"/><circle class="pkt  hop" cx="700" cy="102" r="6" style="--dx:225px;--dy:0px;--t:6s;--d:3.84s"/><circle class="pkt r hop" cx="75" cy="102" r="6" style="--dx:205px;--dy:0px;--t:6s;--d:0.5s"/><circle class="pkt r hop" cx="280" cy="102" r="6" style="--dx:210px;--dy:62px;--t:6s;--d:1.34s"/><circle class="pkt r hop" cx="490" cy="164" r="6" style="--dx:210px;--dy:-62px;--t:6s;--d:2.18s"/><circle class="pkt r hop" cx="700" cy="102" r="6" style="--dx:225px;--dy:0px;--t:6s;--d:3.02s"/><circle class="pkt b hop" cx="75" cy="102" r="6" style="--dx:205px;--dy:0px;--t:6s;--d:1.0s"/><circle class="pkt b lost" cx="280" cy="102" r="6" style="--dx:105px;--dy:-31px;--t:6s;--d:1.84s"/><g class="boom" style="--t:6s;--d:1.84s"><path class="ln x" d="M378 64L392 78M392 64L378 78"/></g><rect class="bx " x="20" y="80" width="110" height="44" rx="3"/><text class="t " x="75" y="107">client</text><rect class="bx " x="240" y="80" width="80" height="44" rx="3"/><text class="t " x="280" y="107">router</text><rect class="bx " x="450" y="18" width="80" height="44" rx="3"/><text class="t " x="490" y="45">router</text><rect class="bx " x="450" y="142" width="80" height="44" rx="3"/><text class="t " x="490" y="169">router</text><rect class="bx " x="660" y="80" width="80" height="44" rx="3"/><text class="t " x="700" y="107">router</text><rect class="bx " x="870" y="80" width="110" height="44" rx="3"/><text class="t " x="925" y="107">server</text><circle class="pkt " cx="170" cy="210" r="6"/><text class="s st" x="184" y="215" style="font-size:13px">packet 1 waits in a busy router</text><circle class="pkt r" cx="450" cy="210" r="6"/><text class="s st" x="464" y="215" style="font-size:13px">packet 2 overtakes it</text><circle class="pkt b" cx="680" cy="210" r="6"/><text class="s st" x="694" y="215" style="font-size:13px">packet 3 is lost</text></svg>

- **IP address** — which machine: `192.168.1.20` (IPv4), `2001:db8::20` (IPv6)
- **Routers** pass each packet one hop closer, each on its own path — `traceroute google.com` shows them
- Only **best effort**: packets can be **lost**, **delayed**, **duplicated**, or **reordered**

---

<span class="tag">Layer 3 · Transport</span>

<div class="mini"><svg class="uml d " width="152" height="110" viewBox="0 0 152 110"><rect x="1" y="2" width="150" height="23" rx="2" style="fill:#14385a;stroke:#6fa8cc;stroke-width:1;opacity:0.7"/><text class="s" x="76" y="18" style="font-size:12px">4 · Application</text><rect x="1" y="29" width="150" height="23" rx="2" style="fill:#5cc8ea;stroke:#6fa8cc;stroke-width:1;opacity:1"/><text class="t" x="76" y="45" style="font-size:12px;fill:#0b2136">3 · Transport</text><rect x="1" y="56" width="150" height="23" rx="2" style="fill:#14385a;stroke:#6fa8cc;stroke-width:1;opacity:0.7"/><text class="s" x="76" y="72" style="font-size:12px">2 · Internet</text><rect x="1" y="83" width="150" height="23" rx="2" style="fill:#14385a;stroke:#6fa8cc;stroke-width:1;opacity:0.7"/><text class="s" x="76" y="99" style="font-size:12px">1 · Link</text></svg></div>

## To the right program

> **Problem:** the packet reached the machine. Which of its programs gets it — and who fixes what IP lost or mixed up?

<svg class="uml d " width="800" height="172.8" viewBox="0 0 1000 216"><rect class="bx " x="20" y="80" width="180" height="56" rx="3"/><text class="t " x="110" y="106">your laptop</text><text class="s" x="110" y="123">192.168.1.10</text><rect class="grp" x="470" y="6" width="510" height="206" rx="4"/><text class="t" x="725" y="30">machine 192.168.1.20</text><path class="ln a" d="M200 108L492.1 69.1"/><polygon class="hd a" points="500,68 491.7,74.1 490.4,64.3"/><text class="m ac" x="352" y="76">to 192.168.1.20 : 443</text><circle class="pkt  fly" cx="200" cy="108" r="6" style="--dx:294px;--dy:-39px;--t:3.2s;--d:0s"/><rect class="bx a" x="500" y="46" width="450" height="44" rx="3"/><text class="t st" x="520" y="73">web server</text><text class="m en ac" x="930" y="73">:443</text><rect class="bx " x="500" y="102" width="450" height="44" rx="3"/><text class="t st" x="520" y="129">SSH server</text><text class="m en ac" x="930" y="129">:22</text><rect class="bx " x="500" y="158" width="450" height="44" rx="3"/><text class="t st" x="520" y="185">PostgreSQL</text><text class="m en ac" x="930" y="185">:5432</text></svg>

- **Port** — which program: a number 0–65535. 80 HTTP · 443 HTTPS · 22 SSH · 53 DNS · 5432 PostgreSQL
- Two services: **TCP** — a reliable, ordered byte stream; **UDP** — single datagrams, nothing promised

---

<span class="tag">Layer 3 · Transport</span>

<div class="mini"><svg class="uml d " width="152" height="110" viewBox="0 0 152 110"><rect x="1" y="2" width="150" height="23" rx="2" style="fill:#14385a;stroke:#6fa8cc;stroke-width:1;opacity:0.7"/><text class="s" x="76" y="18" style="font-size:12px">4 · Application</text><rect x="1" y="29" width="150" height="23" rx="2" style="fill:#5cc8ea;stroke:#6fa8cc;stroke-width:1;opacity:1"/><text class="t" x="76" y="45" style="font-size:12px;fill:#0b2136">3 · Transport</text><rect x="1" y="56" width="150" height="23" rx="2" style="fill:#14385a;stroke:#6fa8cc;stroke-width:1;opacity:0.7"/><text class="s" x="76" y="72" style="font-size:12px">2 · Internet</text><rect x="1" y="83" width="150" height="23" rx="2" style="fill:#14385a;stroke:#6fa8cc;stroke-width:1;opacity:0.7"/><text class="s" x="76" y="99" style="font-size:12px">1 · Link</text></svg></div>

## TCP or UDP?

| | **TCP** | **UDP** |
|---|---|---|
| Connection | yes — a handshake first | none — just send |
| A lost packet | sent again, automatically | gone |
| Order | bytes arrive in order | any order |
| You send | a **stream** of bytes | separate **datagrams** |
| Used by | the web (HTTP/1.1, HTTP/2), SSH, databases | DNS, video calls, games, HTTP/3 (QUIC) |

- TCP also slows down when the receiver or the network can't keep up
- > TCP's "reliable" means: while the connection lives, bytes arrive in order — or not at all. When it breaks, you don't know what the other side got.

---

<span class="tag">Layer 3 · Transport</span>

<div class="mini"><svg class="uml d " width="152" height="110" viewBox="0 0 152 110"><rect x="1" y="2" width="150" height="23" rx="2" style="fill:#14385a;stroke:#6fa8cc;stroke-width:1;opacity:0.7"/><text class="s" x="76" y="18" style="font-size:12px">4 · Application</text><rect x="1" y="29" width="150" height="23" rx="2" style="fill:#5cc8ea;stroke:#6fa8cc;stroke-width:1;opacity:1"/><text class="t" x="76" y="45" style="font-size:12px;fill:#0b2136">3 · Transport</text><rect x="1" y="56" width="150" height="23" rx="2" style="fill:#14385a;stroke:#6fa8cc;stroke-width:1;opacity:0.7"/><text class="s" x="76" y="72" style="font-size:12px">2 · Internet</text><rect x="1" y="83" width="150" height="23" rx="2" style="fill:#14385a;stroke:#6fa8cc;stroke-width:1;opacity:0.7"/><text class="s" x="76" y="99" style="font-size:12px">1 · Link</text></svg></div>

## TCP opens a connection: the handshake

<svg class="uml d tl" width="1000" height="30" viewBox="0 0 1000 30"><text class="cap" x="220" y="20">CLIENT</text><text class="cap" x="640" y="20">SERVER</text><text class="cap" x="870" y="20"></text></svg>

* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="220" y1="0" x2="220" y2="44"/><line class="lane" x1="640" y1="0" x2="640" y2="44"/><text class="s" x="14" y="33">1</text><path class="ln  draw" pathLength="1" d="M225 29L627 29" style="--d:0s"/><g class="late" style="--d:0s"><polygon class="hd " points="635,29 626,34 626,24.1"/></g><circle class="pkt  once" cx="225" cy="29" r="5" style="--dx:402px;--dy:0px;--t:3s;--d:0s"/><g class="late" style="--d:-0.55s"><text class="m " x="430" y="20">SYN</text></g><g class="late" style="--d:0.3s"><text x="870" y="33">client: let's talk</text></g></svg>
* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="220" y1="0" x2="220" y2="44"/><line class="lane" x1="640" y1="0" x2="640" y2="44"/><text class="s" x="14" y="33">2</text><path class="ln  draw" pathLength="1" d="M635 29L233 29" style="--d:0s"/><g class="late" style="--d:0s"><polygon class="hd " points="225,29 234,24.1 234,34"/></g><circle class="pkt r once" cx="635" cy="29" r="5" style="--dx:-402px;--dy:0px;--t:3s;--d:0s"/><g class="late" style="--d:-0.55s"><text class="m " x="430" y="20">SYN-ACK</text></g><g class="late" style="--d:0.3s"><text x="870" y="33">server: OK, let's talk</text></g></svg>
* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="220" y1="0" x2="220" y2="44"/><line class="lane" x1="640" y1="0" x2="640" y2="44"/><text class="s" x="14" y="33">3</text><path class="ln  draw" pathLength="1" d="M225 29L627 29" style="--d:0s"/><g class="late" style="--d:0s"><polygon class="hd " points="635,29 626,34 626,24.1"/></g><circle class="pkt  once" cx="225" cy="29" r="5" style="--dx:402px;--dy:0px;--t:3s;--d:0s"/><g class="late" style="--d:-0.55s"><text class="m " x="430" y="20">ACK</text></g><g class="late" style="--d:0.3s"><text class="gd" x="870" y="33">connected ✓</text></g></svg>
* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="220" y1="0" x2="220" y2="44"/><line class="lane" x1="640" y1="0" x2="640" y2="44"/><text class="s" x="14" y="33">4</text><path class="ln  draw" pathLength="1" d="M225 29L627 29" style="--d:0s"/><g class="late" style="--d:0s"><polygon class="hd " points="635,29 626,34 626,24.1"/></g><circle class="pkt  once" cx="225" cy="29" r="5" style="--dx:402px;--dy:0px;--t:3s;--d:0s"/><g class="late" style="--d:-0.55s"><text class="m " x="430" y="20">GET /abc123 …</text></g><g class="late" style="--d:0.3s"><text x="870" y="33">now the request</text></g></svg>
* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="220" y1="0" x2="220" y2="44"/><line class="lane" x1="640" y1="0" x2="640" y2="44"/><text class="s" x="14" y="33">5</text><path class="ln  draw" pathLength="1" d="M635 29L233 29" style="--d:0s"/><g class="late" style="--d:0s"><polygon class="hd " points="225,29 234,24.1 234,34"/></g><circle class="pkt r once" cx="635" cy="29" r="5" style="--dx:-402px;--dy:0px;--t:3s;--d:0s"/><g class="late" style="--d:-0.55s"><text class="m " x="430" y="20">302 Found …</text></g><g class="late" style="--d:0.3s"><text x="870" y="33">and the response</text></g></svg>

<div class="sp"></div>

* One full round trip **before** the request can even leave. California ↔ Europe: ~150 ms just to say hello. HTTPS adds one more round trip (TLS 1.3)

---

<span class="tag">Layer 4 · Application</span>

<div class="mini"><svg class="uml d " width="152" height="110" viewBox="0 0 152 110"><rect x="1" y="2" width="150" height="23" rx="2" style="fill:#5cc8ea;stroke:#6fa8cc;stroke-width:1;opacity:1"/><text class="t" x="76" y="18" style="font-size:12px;fill:#0b2136">4 · Application</text><rect x="1" y="29" width="150" height="23" rx="2" style="fill:#14385a;stroke:#6fa8cc;stroke-width:1;opacity:0.7"/><text class="s" x="76" y="45" style="font-size:12px">3 · Transport</text><rect x="1" y="56" width="150" height="23" rx="2" style="fill:#14385a;stroke:#6fa8cc;stroke-width:1;opacity:0.7"/><text class="s" x="76" y="72" style="font-size:12px">2 · Internet</text><rect x="1" y="83" width="150" height="23" rx="2" style="fill:#14385a;stroke:#6fa8cc;stroke-width:1;opacity:0.7"/><text class="s" x="76" y="99" style="font-size:12px">1 · Link</text></svg></div>

## What the bytes mean

> **Problem:** the bytes arrived, in order. What do they mean? And how did we find the server — we only knew its **name**?

<svg class="uml d " width="1000" height="156" viewBox="0 0 1000 156"><path class="ln " d="M197.9 70.7L412.1 35.3"/><polygon class="hd " points="420,34 411.9,40.4 410.3,30.6"/><polygon class="hd " points="190,72 198.1,65.6 199.7,75.4"/><path class="ln " d="M198 98.4L712 125.6"/><polygon class="hd " points="720,126 710.8,130.5 711.3,120.6"/><polygon class="hd " points="190,98 199.2,93.5 198.7,103.4"/><text class="m ac" x="300" y="38" style="font-size:13px">example.com?</text><text class="m gd" x="330" y="76" style="font-size:13px">172.66.147.243</text><text class="s" x="470" y="128">then: connect · TCP · HTTP</text><circle class="pkt  hop" cx="190" cy="72" r="6" style="--dx:230px;--dy:-38px;--t:6s;--d:0s"/><circle class="pkt r hop" cx="420" cy="34" r="6" style="--dx:-230px;--dy:38px;--t:6s;--d:0.9s"/><circle class="pkt  hop" cx="190" cy="98" r="6" style="--dx:530px;--dy:28px;--t:6s;--d:1.8s"/><circle class="pkt r hop" cx="720" cy="126" r="6" style="--dx:-530px;--dy:-28px;--t:6s;--d:2.7s"/><rect class="bx " x="20" y="60" width="170" height="50" rx="3"/><text class="t " x="105" y="90">your program</text><rect class="bx " x="420" y="8" width="220" height="50" rx="3"/><text class="t " x="530" y="31">DNS server</text><text class="s" x="530" y="48">port 53, usually UDP</text><rect class="bx " x="720" y="100" width="260" height="50" rx="3"/><text class="t " x="850" y="123">example.com</text><text class="s" x="850" y="140">172.66.147.243</text></svg>

- **Protocols** give the bytes meaning: HTTP (the web), SMTP (email), SSH (a remote shell)
- **DNS** turns a name into an IP address — before almost every connection
- This is the layer your program lives in. Everything below just moves bytes

---

## What the network does not guarantee

| What can happen | What your program must do |
|---|---|
| a message is **lost** — or its reply is | time out, then retry |
| a message arrives **twice** | survive the repeat: idempotency |
| messages arrive **out of order** | number them; don't trust arrival order |
| a message is **late** — 1 ms today, 2 s tomorrow | never wait forever: always a timeout |
| the other machine **crashed** — or is only slow | plan for "unknown": you can't tell which |
| the connection **breaks** halfway | assume the other side got part of it |
| someone on the way **reads or changes** the bytes | encrypt: TLS (HTTPS) |
| the network **splits** — both halves alive | decide: refuse, or answer and repair later |

- > Your program must stay correct when any of these happens. The rest of today is how.

---

<!-- _class: lead invert -->

# Client & Server

### The smallest distributed system — and how it serves many

---

## Client and server

<svg class="uml d " width="1000" height="204" viewBox="0 0 1000 204"><path class="ln " d="M197.9 34L712.1 99"/><polygon class="hd " points="720,100 710.5,103.8 711.7,94"/><polygon class="hd " points="190,33 199.5,29.2 198.3,39"/><circle class="pkt  fly" cx="190" cy="33" r="6" style="--dx:530px;--dy:67px;--t:3s;--d:0.0s"/><circle class="pkt r fly" cx="720" cy="100" r="6" style="--dx:-530px;--dy:-67px;--t:3s;--d:1.4s"/><path class="ln " d="M198 100L712 100"/><polygon class="hd " points="720,100 711,105 711,95"/><polygon class="hd " points="190,100 199,95 199,105"/><circle class="pkt  fly" cx="190" cy="100" r="6" style="--dx:530px;--dy:0px;--t:3s;--d:1.0s"/><circle class="pkt r fly" cx="720" cy="100" r="6" style="--dx:-530px;--dy:0px;--t:3s;--d:2.4s"/><path class="ln " d="M197.9 166L712.1 101"/><polygon class="hd " points="720,100 711.7,106 710.5,96.2"/><polygon class="hd " points="190,167 198.3,161 199.5,170.8"/><circle class="pkt  fly" cx="190" cy="167" r="6" style="--dx:530px;--dy:-67px;--t:3s;--d:2.0s"/><circle class="pkt r fly" cx="720" cy="100" r="6" style="--dx:-530px;--dy:67px;--t:3s;--d:3.4s"/><rect class="bx " x="20" y="10" width="170" height="46" rx="3"/><text class="t " x="105" y="38">browser</text><rect class="bx " x="20" y="77" width="170" height="46" rx="3"/><text class="t " x="105" y="105">phone app</text><rect class="bx " x="20" y="144" width="170" height="46" rx="3"/><text class="t " x="105" y="172">another service</text><rect class="bx a" x="720" y="55" width="260" height="90" rx="3"/><text class="t " x="850" y="98">server</text><text class="s" x="850" y="115">listens on port 443</text><text class="s" x="455" y="197" style="font-size:13px">request →     ← response</text></svg>

- **Server**: starts first, **listens** on a known port, waits — and serves many clients
- **Client**: connects when it needs something, sends a **request**, waits for the **response**
- Browser → web server; app → API; your code → database. A database is a server too
- A server is often a client of another server — that's how big systems are built

---

## Sockets: the OS's network API

<svg class="uml d tl" width="1000" height="30" viewBox="0 0 1000 30"><text class="cap" x="250" y="20">CLIENT</text><text class="cap" x="690" y="20">SERVER</text><text class="cap" x="900" y="20"></text></svg>

* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="250" y1="0" x2="250" y2="44"/><line class="lane" x1="690" y1="0" x2="690" y2="44"/><text class="s" x="14" y="33">1</text><g class="late" style="--d:-0.45s"><rect class="bx " x="539.4" y="6" width="301.2" height="32" rx="3"/><text class="m " x="690" y="27" style="font-size:13px">socket() · bind(9000) · listen()</text></g></svg>
* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="250" y1="0" x2="250" y2="44"/><line class="lane" x1="690" y1="0" x2="690" y2="44"/><text class="s" x="14" y="33">2</text><g class="late" style="--d:-0.45s"><rect class="bx " x="603.9" y="6" width="172.2" height="32" rx="3"/><text class="m " x="690" y="27" style="font-size:13px">accept()  … waits</text></g></svg>
* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="250" y1="0" x2="250" y2="44"/><line class="lane" x1="690" y1="0" x2="690" y2="44"/><text class="s" x="14" y="33">3</text><path class="ln  draw" pathLength="1" d="M255 29L677 29" style="--d:0s"/><g class="late" style="--d:0s"><polygon class="hd " points="685,29 676,34 676,24.1"/></g><circle class="pkt  once" cx="255" cy="29" r="5" style="--dx:422px;--dy:0px;--t:3s;--d:0s"/><g class="late" style="--d:-0.55s"><text class="m " x="470" y="20">connect() — the TCP handshake</text></g><g class="late" style="--d:0.4s"><text class="gd" x="900" y="33">accept() returns</text></g></svg>
* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="250" y1="0" x2="250" y2="44"/><line class="lane" x1="690" y1="0" x2="690" y2="44"/><text class="s" x="14" y="33">4</text><path class="ln  draw" pathLength="1" d="M255 29L677 29" style="--d:0s"/><g class="late" style="--d:0s"><polygon class="hd " points="685,29 676,34 676,24.1"/></g><circle class="pkt  once" cx="255" cy="29" r="5" style="--dx:422px;--dy:0px;--t:3s;--d:0s"/><g class="late" style="--d:-0.55s"><text class="m " x="470" y="20">send(request)</text></g><g class="late" style="--d:0.4s"><text x="900" y="33">recv()</text></g></svg>
* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="250" y1="0" x2="250" y2="44"/><line class="lane" x1="690" y1="0" x2="690" y2="44"/><text class="s" x="14" y="33">5</text><path class="ln  draw" pathLength="1" d="M685 29L263 29" style="--d:0s"/><g class="late" style="--d:0s"><polygon class="hd " points="255,29 264,24.1 264,34"/></g><circle class="pkt r once" cx="685" cy="29" r="5" style="--dx:-422px;--dy:0px;--t:3s;--d:0s"/><g class="late" style="--d:-0.55s"><text class="m " x="470" y="20">send(response)</text></g></svg>
* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="250" y1="0" x2="250" y2="44"/><line class="lane" x1="690" y1="0" x2="690" y2="44"/><text class="s" x="14" y="33">6</text><path class="ln  draw" pathLength="1" d="M255 29L677 29" style="--d:0s"/><g class="late" style="--d:0s"><polygon class="hd " points="685,29 676,34 676,24.1"/></g><circle class="pkt  once" cx="255" cy="29" r="5" style="--dx:422px;--dy:0px;--t:3s;--d:0s"/><g class="late" style="--d:-0.55s"><text class="m " x="470" y="20">close()</text></g><g class="late" style="--d:0.4s"><text x="900" y="33">recv() returns 0</text></g></svg>

<div class="sp"></div>

* Day 2: these are **system calls**. A socket is a handle the OS gives you — like an open file

---

<!-- _class: invert good dense -->

## ✅ An echo server

```cpp
int srv = socket(AF_INET, SOCK_STREAM, 0);            // IPv4 + TCP
sockaddr_in addr{};                                   // address 0.0.0.0: every network card
addr.sin_family = AF_INET;
addr.sin_port = htons(9000);                          // port 9000
bind(srv, (sockaddr*)&addr, sizeof addr);
listen(srv, SOMAXCONN);                               // ready: clients may connect now

for (;;) {
    int conn = accept(srv, nullptr, nullptr);         // waits for the next client
    char buf[1024];
    ssize_t n;
    while ((n = recv(conn, buf, sizeof buf, 0)) > 0)  // 0 means: the client closed
        send(conn, buf, n, 0);                        // echo it back
    close(conn);
}
```

- The same system calls under every language: Python `socket`, Java `ServerSocket`, Go `net.Listen`, C++ Boost.Asio

<p class="note">Linux: &lt;sys/socket.h&gt;, &lt;netinet/in.h&gt;, &lt;unistd.h&gt;. Error checks left out to fit the slide.</p>

---

<!-- _class: dense -->

## The client

```cpp
int fd = socket(AF_INET, SOCK_STREAM, 0);
sockaddr_in addr{};
addr.sin_family = AF_INET;
addr.sin_port = htons(9000);
inet_pton(AF_INET, "127.0.0.1", &addr.sin_addr);
connect(fd, (sockaddr*)&addr, sizeof addr);           // the TCP handshake happens here

send(fd, "hello", 5, 0);
char buf[1024];
ssize_t n = recv(fd, buf, sizeof buf, 0);
printf("%.*s\n", (int)n, buf);                        // hello
close(fd);
```


---

<!-- _class: invert bad -->

## TCP is a stream, not messages

<svg class="uml d " width="1000" height="148" viewBox="0 0 1000 148"><rect class="bx " x="20" y="22" width="220" height="40" rx="3"/><text class="t m " x="130" y="47">send("hello")</text><rect class="bx " x="20" y="98" width="220" height="40" rx="3"/><text class="t m " x="130" y="123">send("world")</text><rect class="grp" x="290" y="40" width="440" height="80" rx="6"/><text class="cap" x="510" y="32">TCP: A STREAM OF BYTES</text><path class="ln " d="M240 42L281.1 66"/><polygon class="hd " points="288,70 277.7,69.7 282.7,61.2"/><path class="ln " d="M240 118L281 95.8"/><polygon class="hd " points="288,92 282.4,100.6 277.7,91.9"/><path class="ln " d="M732 80L770 80"/><polygon class="hd " points="778,80 769,85 769,75"/><text class="m ac t fly" x="310" y="87" style="--dx:405px;--dy:0px;--t:5s;--d:0.0s">h</text><text class="m ac t fly" x="310" y="87" style="--dx:405px;--dy:0px;--t:5s;--d:0.28s">e</text><text class="m ac t fly" x="310" y="87" style="--dx:405px;--dy:0px;--t:5s;--d:0.56s">l</text><text class="m ac t fly" x="310" y="87" style="--dx:405px;--dy:0px;--t:5s;--d:0.8400000000000001s">l</text><text class="m ac t fly" x="310" y="87" style="--dx:405px;--dy:0px;--t:5s;--d:1.12s">o</text><text class="m gd t fly" x="310" y="87" style="--dx:405px;--dy:0px;--t:5s;--d:1.4000000000000001s">w</text><text class="m gd t fly" x="310" y="87" style="--dx:405px;--dy:0px;--t:5s;--d:1.6800000000000002s">o</text><text class="m gd t fly" x="310" y="87" style="--dx:405px;--dy:0px;--t:5s;--d:1.9600000000000002s">r</text><text class="m gd t fly" x="310" y="87" style="--dx:405px;--dy:0px;--t:5s;--d:2.24s">l</text><text class="m gd t fly" x="310" y="87" style="--dx:405px;--dy:0px;--t:5s;--d:2.5200000000000005s">d</text><rect class="bx " x="780" y="52" width="200" height="56" rx="3"/><text class="t m " x="880" y="78">recv(…, 1024)</text><text class="s" x="880" y="95">→ 10 bytes: helloworld</text><text class="s" x="510" y="140" style="font-size:13px">no boundary between “hello” and “world” — just bytes</text></svg>

```cpp
send(fd, "hello", 5, 0);
send(fd, "world", 5, 0);
recv(conn, buf, 1024, 0);    // on the server, a moment later: 10 bytes, "helloworld"
```

- TCP keeps the **bytes** in order. It does **not** keep your message boundaries
- `recv` returns **whatever has arrived**, up to the size you ask for: half a message, or two. One 4 MB send arrived in pieces — the first `recv` got 128,000 bytes

---

<!-- _class: invert good -->

## ✅ Framing: where does a message end?

<svg class="uml d " width="1000" height="126" viewBox="0 0 1000 126"><text class="cap st" x="20" y="45">DELIMITER</text><text class="cap st" x="20" y="105">LENGTH FIRST</text><rect class="bx " x="150" y="22" width="150" height="36" rx="3"/><text class="t m " x="225" y="45">hello</text><rect class="bx a" x="300" y="22" width="50" height="36" rx="3"/><text class="t m " x="325" y="45">\n</text><rect class="bx " x="350" y="22" width="150" height="36" rx="3"/><text class="t m " x="425" y="45">world</text><rect class="bx a" x="500" y="22" width="50" height="36" rx="3"/><text class="t m " x="525" y="45">\n</text><rect class="bx a" x="150" y="82" width="150" height="36" rx="3"/><text class="t m " x="225" y="105">00 00 00 05</text><rect class="bx " x="300" y="82" width="150" height="36" rx="3"/><text class="t m " x="375" y="105">hello</text><rect class="bx a" x="450" y="82" width="150" height="36" rx="3"/><text class="t m " x="525" y="105">00 00 00 05</text><rect class="bx " x="600" y="82" width="150" height="36" rx="3"/><text class="t m " x="675" y="105">world</text><text class="s st" x="570" y="45">read until you see \n</text><text class="s st" x="770" y="98">read 4 bytes: 5 —</text><text class="s st" x="770" y="114">then exactly 5 more</text></svg>

- **Delimiter**: end every message with a special byte, like a newline. HTTP ends each header line with `\r\n`
- **Length first**: send the size in 4 bytes, then exactly that many bytes — most binary protocols
- Either way, the receiver calls `recv` in a loop until it has a whole message

---

## HTTP: the protocol everyone speaks

```text
→ GET /abc123 HTTP/1.1                 method · path · version
→ Host: short.example                  headers, one per line
→                                      an empty line: the request ends

← HTTP/1.1 302 Found                   status
← Location: https://example.com/a/very/long/link
← Content-Length: 0                    body length: no body
```

- Text over TCP, framed by lines and `Content-Length`. HTTP/2 and HTTP/3 are binary, with the same meaning
- Status codes: **2xx** done · **3xx** look elsewhere · **4xx** the client's mistake · **5xx** the server's
- Day 2's link shortener, `GET /{code}` → `302` — this is what it looks like on the wire

<p class="note">Captured from Python's http.server; its Date and Server headers are trimmed.</p>

---

## HTTP methods and REST

| Request | Means | Safe to repeat? |
|---|---|---|
| `GET /urls/abc123` | read it | ✅ changes nothing |
| `PUT /urls/abc123` | create or replace it | ✅ idempotent |
| `DELETE /urls/abc123` | remove it | ✅ idempotent |
| `POST /urls` | create a new one | ❌ twice = two links |
| `PATCH /urls/abc123` | change part of it | ❌ not promised |

- **REST**: URLs name the things (resources), methods are the verbs, bodies are usually JSON
- `POST /urls {"long_url": "https://…"}` → `201 Created {"code": "abc123"}`
- The last column decides what a client may **retry** — later today

---

<!-- _class: invert good -->

## ✅ A thread per connection

```cpp
void handle(int conn) {                               // runs in its own thread
    char buf[1024];
    ssize_t n;
    while ((n = recv(conn, buf, sizeof buf, 0)) > 0)
        send(conn, buf, n, 0);
    close(conn);
}

for (;;) {                                            // srv: socket + bind + listen, as before
    int conn = accept(srv, nullptr, nullptr);
    std::thread(handle, conn).detach();               // one thread per client
}
```

- A slow client blocks only its own thread. 200 clients at once: all answered
- Handlers share data? Then it's Day 3 again: locks
- The price: a thread per client — a stack each, and the OS switching between thousands

---

## A thread pool, and the slow client

<svg class="uml d " width="1000" height="214" viewBox="0 0 1000 214"><rect class="bx " x="20" y="88" width="140" height="44" rx="3"/><text class="t m " x="90" y="115">accept()</text><text class="cap" x="300" y="82">QUEUE</text><rect class="bx" x="200" y="92" width="34" height="36" rx="2"/><rect class="bx" x="240" y="92" width="34" height="36" rx="2"/><rect class="bx" x="280" y="92" width="34" height="36" rx="2"/><rect class="bx" x="320" y="92" width="34" height="36" rx="2"/><rect class="bx" x="360" y="92" width="34" height="36" rx="2"/><path class="ln " d="M160 110L190 110"/><polygon class="hd " points="198,110 189,115 189,105"/><path class="ln " d="M398 110L470 31"/><path class="ln b" d="M660 31L760 31"/><path class="ln " d="M398 110L470 83"/><path class="ln b" d="M660 83L760 83"/><path class="ln " d="M398 110L470 135"/><path class="ln " d="M660 135L760 135"/><path class="ln " d="M398 110L470 187"/><path class="ln " d="M660 187L760 187"/><circle class="pkt  fly" cx="160" cy="110" r="6" style="--dx:40px;--dy:0px;--t:2.5s;--d:0s"/><circle class="pkt  hop" cx="398" cy="110" r="6" style="--dx:72px;--dy:25px;--t:5s;--d:1.0s"/><circle class="pkt  hop" cx="398" cy="110" r="6" style="--dx:72px;--dy:77px;--t:5s;--d:3.5s"/><circle class="pkt  fly" cx="660" cy="135" r="6" style="--dx:100px;--dy:0px;--t:2.5s;--d:0.5s"/><circle class="pkt r fly" cx="760" cy="187" r="6" style="--dx:-100px;--dy:0px;--t:2.5s;--d:1.5s"/><rect class="bx b" x="470" y="10" width="190" height="42" rx="3"/><text class="t " x="565" y="29">worker 1</text><text class="s" x="565" y="46">slow client: 30 s</text><rect class="bx b" x="760" y="10" width="130" height="42" rx="3"/><text class="t " x="825" y="36">client</text><rect class="bx b" x="470" y="62" width="190" height="42" rx="3"/><text class="t " x="565" y="81">worker 2</text><text class="s" x="565" y="98">slow client: 30 s</text><rect class="bx b" x="760" y="62" width="130" height="42" rx="3"/><text class="t " x="825" y="88">client</text><rect class="bx " x="470" y="114" width="190" height="42" rx="3"/><text class="t " x="565" y="133">worker 3</text><text class="s" x="565" y="150">short request</text><rect class="bx " x="760" y="114" width="130" height="42" rx="3"/><text class="t " x="825" y="140">client</text><rect class="bx " x="470" y="166" width="190" height="42" rx="3"/><text class="t " x="565" y="185">worker 4</text><text class="s" x="565" y="202">short request</text><rect class="bx " x="760" y="166" width="130" height="42" rx="3"/><text class="t " x="825" y="192">client</text></svg>

- Day 4's pool: `accept()` hands each connection to one of N workers
- Threads and memory stay bounded — but N slow clients take all N workers, and everyone else waits
- Good when requests are short. A long connection — chat, streaming — holds its worker the whole time

---

<!-- _class: invert good dense -->

## ✅ An event loop: one thread, many connections

```cpp
int ep = epoll_create1(0);
epoll_event ev{EPOLLIN, {.fd = srv}};
epoll_ctl(ep, EPOLL_CTL_ADD, srv, &ev);               // watch the listening socket
epoll_event ready[64];
char buf[1024];
for (;;) {
    int n = epoll_wait(ep, ready, 64, -1);            // sleep until some socket is ready
    for (int i = 0; i < n; ++i) {
        int fd = ready[i].data.fd;
        if (fd == srv) {                              // a new client: watch it too
            epoll_event e{EPOLLIN, {.fd = accept(srv, nullptr, nullptr)}};
            epoll_ctl(ep, EPOLL_CTL_ADD, e.data.fd, &e);
        } else if (ssize_t k = recv(fd, buf, sizeof buf, 0); k > 0)
            send(fd, buf, k, 0);                      // a client sent data: echo it
        else
            close(fd);                                // it closed: epoll forgets it
    }
}
```

- `epoll` (Linux) reports which sockets are **ready** — one thread served 1000 clients

---

<!-- _class: lead invert -->

# When the Network Fails

### Timeouts, retries, idempotency

---

<!-- _class: invert bad -->

## You sent a request. No reply.

<svg class="uml d tl" width="1000" height="30" viewBox="0 0 1000 30"><text class="cap" x="150" y="20">CLIENT</text><text class="cap" x="590" y="20">SERVER</text><text class="cap" x="870" y="20">WHAT HAPPENED</text></svg>

* <svg class="uml d tl" width="1000" height="66" viewBox="0 0 1000 66"><line class="lane" x1="150" y1="0" x2="150" y2="66"/><line class="lane" x1="590" y1="0" x2="590" y2="66"/><text class="s" x="14" y="47.6">1</text><path class="ln  draw" pathLength="1" d="M155 43.6L364 43.6" style="--d:0s"/><g class="late" style="--d:0s"><path class="ln x" d="M364 37.6L376 49.6M376 37.6L364 49.6"/></g><circle class="pkt b once" cx="155" cy="43.6" r="5" style="--dx:209px;--dy:0px;--t:3s;--d:0s"/><g class="late" style="--d:-0.55s"><text class="m " x="260" y="34.6">withdraw 100</text></g><g class="late" style="--d:0.3s"><text class="bd" x="870" y="47.6">request lost</text></g></svg>
* <svg class="uml d tl" width="1000" height="66" viewBox="0 0 1000 66"><line class="lane" x1="150" y1="0" x2="150" y2="66"/><line class="lane" x1="590" y1="0" x2="590" y2="66"/><text class="s" x="14" y="47.6">2</text><path class="ln  draw" pathLength="1" d="M155 43.6L577 43.6" style="--d:0s"/><g class="late" style="--d:0s"><polygon class="hd " points="585,43.6 576,48.5 576,38.6"/></g><circle class="pkt  once" cx="155" cy="43.6" r="5" style="--dx:422px;--dy:0px;--t:3s;--d:0s"/><g class="late" style="--d:-0.55s"><text class="m " x="370" y="34.6">withdraw 100</text></g><g class="late" style="--d:0.3s"><rect class="bx b" x="622.5" y="6" width="155" height="54" rx="3"/><text class="m " x="700" y="38" style="font-size:13px">busy … or dead?</text></g><g class="late" style="--d:0.6s"><text class="bd" x="870" y="47.6">server too slow</text></g></svg>
* <svg class="uml d tl" width="1000" height="66" viewBox="0 0 1000 66"><line class="lane" x1="150" y1="0" x2="150" y2="66"/><line class="lane" x1="590" y1="0" x2="590" y2="66"/><text class="s" x="14" y="47.6">3</text><path class="ln  draw" pathLength="1" d="M155 27.7L577 27.7" style="--d:0s"/><g class="late" style="--d:0s"><polygon class="hd " points="585,27.7 576,32.7 576,22.8"/></g><circle class="pkt  once" cx="155" cy="27.7" r="5" style="--dx:422px;--dy:0px;--t:3s;--d:0s"/><g class="late" style="--d:-0.55s"><text class="m " x="370" y="18.7">withdraw 100</text></g><g class="late" style="--d:0.2s"><rect class="bx g" x="644" y="6" width="112" height="54" rx="3"/><text class="m " x="700" y="38" style="font-size:13px">done: −100</text></g><path class="ln  draw" pathLength="1" d="M585 56.8L376 56.8" style="--d:1.0s"/><g class="late" style="--d:1.0s"><path class="ln x" d="M364 50.8L376 62.8M376 50.8L364 62.8"/></g><circle class="pkt r once" cx="585" cy="56.8" r="5" style="--dx:-209px;--dy:0px;--t:3s;--d:1.0s"/><g class="late" style="--d:0.45s"><text class="m " x="480" y="47.8">OK</text></g><g class="late" style="--d:1.2s"><text class="bd" x="870" y="47.6">done — reply lost</text></g></svg>

<div class="sp"></div>

* > From the client, all three look the **same**: silence. You **cannot tell** a slow machine from a dead one — or "not done" from "done".

---

## Timeouts: the only detector you have

<svg class="uml d " width="1000" height="146" viewBox="0 0 1000 146"><path class="ln " d="M70 96L950 96"/><polygon class="hd " points="958,96 949,101 949,91"/><path class="ln " d="M70 92L70 100"/><text class="s" x="70" y="118">0 s</text><path class="ln " d="M230 92L230 100"/><text class="s" x="230" y="118">0.5 s</text><path class="ln " d="M390 92L390 100"/><text class="s" x="390" y="118">1 s</text><path class="ln " d="M550 92L550 100"/><text class="s" x="550" y="118">1.5 s</text><path class="ln " d="M710 92L710 100"/><text class="s" x="710" y="118">2 s</text><path class="ln " d="M870 92L870 100"/><text class="s" x="870" y="118">2.5 s</text><rect x="76.4" y="62" width="121.6" height="34" rx="3" style="fill:rgba(126,208,163,.18);stroke:#7ed0a3;stroke-width:1"/><text class="s gd" x="137.2" y="54">normal replies: 20–400 ms</text><line x1="710" y1="34" x2="710" y2="104" style="stroke:#f0a06e;stroke-width:2;stroke-dasharray:5 4"/><text class="s bd" x="710" y="26">timeout: stop waiting</text><text class="s ac" x="70" y="140">request sent</text><circle class="pkt  fly" cx="70" cy="80" r="7" style="--dx:636px;--dy:0px;--t:4s;--d:0s"/><g class="boom" style="--t:4s;--d:0s"><path class="ln x" d="M702 72L718 88M718 72L702 88"/></g></svg>

- No reply yet? Wait at most **T**, then stop waiting. Many APIs wait **forever** by default — a blocking `recv`, many HTTP clients
- Too **short**: you give up on servers that were only slow, and your retries add load. Too **long**: users wait, threads pile up
- Choose from data: a bit above the slowest normal response

---

<!-- _class: invert bad -->

## A retry can do it twice

<svg class="uml d tl" width="1000" height="30" viewBox="0 0 1000 30"><text class="cap" x="200" y="20">CLIENT</text><text class="cap" x="600" y="20">SERVER</text><text class="cap" x="880" y="20">CHARGED</text></svg>

* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="200" y1="0" x2="200" y2="44"/><line class="lane" x1="600" y1="0" x2="600" y2="44"/><text class="s" x="14" y="33">1</text><path class="ln  draw" pathLength="1" d="M205 29L587 29" style="--d:0s"/><g class="late" style="--d:0s"><polygon class="hd " points="595,29 586,34 586,24.1"/></g><circle class="pkt  once" cx="205" cy="29" r="5" style="--dx:382px;--dy:0px;--t:3s;--d:0s"/><g class="late" style="--d:-0.55s"><text class="m " x="400" y="20">charge 100</text></g></svg>
* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="200" y1="0" x2="200" y2="44"/><line class="lane" x1="600" y1="0" x2="600" y2="44"/><text class="s" x="14" y="33">2</text><g class="late" style="--d:-0.45s"><rect class="bx " x="518.2" y="6" width="163.6" height="32" rx="3"/><text class="m " x="600" y="27" style="font-size:13px">charges the card</text></g><g class="late" style="--d:0.2s"><text x="880" y="33">100</text></g></svg>
* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="200" y1="0" x2="200" y2="44"/><line class="lane" x1="600" y1="0" x2="600" y2="44"/><text class="s" x="14" y="33">3</text><path class="ln  draw" pathLength="1" d="M595 29L406 29" style="--d:0s"/><g class="late" style="--d:0s"><path class="ln x" d="M394 23L406 35M406 23L394 35"/></g><circle class="pkt r once" cx="595" cy="29" r="5" style="--dx:-189px;--dy:0px;--t:3s;--d:0s"/><g class="late" style="--d:-0.55s"><text class="m " x="500" y="20">OK</text></g></svg>
* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="200" y1="0" x2="200" y2="44"/><line class="lane" x1="600" y1="0" x2="600" y2="44"/><text class="s" x="14" y="33">4</text><g class="late" style="--d:-0.45s"><rect class="bx b" x="122.5" y="6" width="155" height="32" rx="3"/><text class="m " x="200" y="27" style="font-size:13px">timeout → retry</text></g></svg>
* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="200" y1="0" x2="200" y2="44"/><line class="lane" x1="600" y1="0" x2="600" y2="44"/><text class="s" x="14" y="33">5</text><path class="ln  draw" pathLength="1" d="M205 29L587 29" style="--d:0s"/><g class="late" style="--d:0s"><polygon class="hd " points="595,29 586,34 586,24.1"/></g><circle class="pkt  once" cx="205" cy="29" r="5" style="--dx:382px;--dy:0px;--t:3s;--d:0s"/><g class="late" style="--d:-0.55s"><text class="m " x="400" y="20">charge 100</text></g></svg>
* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="200" y1="0" x2="200" y2="44"/><line class="lane" x1="600" y1="0" x2="600" y2="44"/><text class="s" x="14" y="33">6</text><g class="late" style="--d:-0.45s"><rect class="bx " x="518.2" y="6" width="163.6" height="32" rx="3"/><text class="m " x="600" y="27" style="font-size:13px">charges the card</text></g><g class="late" style="--d:0.2s"><text class="bd t" x="880" y="33">200 ✗</text></g></svg>

<div class="sp"></div>

* The server did nothing wrong. The client did nothing wrong. The customer paid twice

---

<!-- _class: invert good -->

## ✅ Idempotency keys: a retry that's safe

<svg class="uml d tl" width="1000" height="30" viewBox="0 0 1000 30"><text class="cap" x="200" y="20">CLIENT</text><text class="cap" x="600" y="20">SERVER</text><text class="cap" x="880" y="20">CHARGED</text></svg>

* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="200" y1="0" x2="200" y2="44"/><line class="lane" x1="600" y1="0" x2="600" y2="44"/><text class="s" x="14" y="33">1</text><path class="ln  draw" pathLength="1" d="M205 29L587 29" style="--d:0s"/><g class="late" style="--d:0s"><polygon class="hd " points="595,29 586,34 586,24.1"/></g><circle class="pkt  once" cx="205" cy="29" r="5" style="--dx:382px;--dy:0px;--t:3s;--d:0s"/><g class="late" style="--d:-0.55s"><text class="m " x="400" y="20">charge 100 · key K</text></g></svg>
* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="200" y1="0" x2="200" y2="44"/><line class="lane" x1="600" y1="0" x2="600" y2="44"/><text class="s" x="14" y="33">2</text><g class="late" style="--d:-0.45s"><rect class="bx " x="475.2" y="6" width="249.6" height="32" rx="3"/><text class="m " x="600" y="27" style="font-size:13px">new K: charge, save K → OK</text></g><g class="late" style="--d:0.2s"><text x="880" y="33">100</text></g></svg>
* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="200" y1="0" x2="200" y2="44"/><line class="lane" x1="600" y1="0" x2="600" y2="44"/><text class="s" x="14" y="33">3</text><path class="ln  draw" pathLength="1" d="M595 29L406 29" style="--d:0s"/><g class="late" style="--d:0s"><path class="ln x" d="M394 23L406 35M406 23L394 35"/></g><circle class="pkt r once" cx="595" cy="29" r="5" style="--dx:-189px;--dy:0px;--t:3s;--d:0s"/><g class="late" style="--d:-0.55s"><text class="m " x="500" y="20">OK</text></g></svg>
* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="200" y1="0" x2="200" y2="44"/><line class="lane" x1="600" y1="0" x2="600" y2="44"/><text class="s" x="14" y="33">4</text><g class="late" style="--d:-0.45s"><rect class="bx b" x="88.1" y="6" width="223.8" height="32" rx="3"/><text class="m " x="200" y="27" style="font-size:13px">timeout → retry, same K</text></g></svg>
* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="200" y1="0" x2="200" y2="44"/><line class="lane" x1="600" y1="0" x2="600" y2="44"/><text class="s" x="14" y="33">5</text><path class="ln  draw" pathLength="1" d="M205 29L587 29" style="--d:0s"/><g class="late" style="--d:0s"><polygon class="hd " points="595,29 586,34 586,24.1"/></g><circle class="pkt  once" cx="205" cy="29" r="5" style="--dx:382px;--dy:0px;--t:3s;--d:0s"/><g class="late" style="--d:-0.55s"><text class="m " x="400" y="20">charge 100 · key K</text></g></svg>
* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="200" y1="0" x2="200" y2="44"/><line class="lane" x1="600" y1="0" x2="600" y2="44"/><text class="s" x="14" y="33">6</text><g class="late" style="--d:-0.45s"><rect class="bx g" x="643.9" y="6" width="172.2" height="32" rx="3"/><text class="m " x="730" y="27" style="font-size:13px">seen K: no charge</text></g><path class="ln g draw" pathLength="1" d="M595 29L213 29" style="--d:0.4s"/><g class="late" style="--d:0.4s"><polygon class="hd g" points="205,29 214,24.1 214,34"/></g><circle class="pkt r once" cx="595" cy="29" r="5" style="--dx:-382px;--dy:0px;--t:3s;--d:0.4s"/><g class="late" style="--d:-0.15s"><text class="m " x="400" y="20">the saved OK</text></g><g class="late" style="--d:0.2s"><text class="gd t" x="880" y="33">100 ✓</text></g></svg>

<div class="sp"></div>

* The client makes a unique **key** per operation (a UUID) and sends the **same** key with every retry. The server remembers key → result
* Two copies at the same moment? Check-then-act — Day 3: store the key atomically first. Stripe's API: the `Idempotency-Key` header

---

## Delivery guarantees

| | How | Lost? | Twice? | Good for |
|---|---|---|---|---|
| **at-most-once** | send once, never retry | can be | never | metrics, video frames |
| **at-least-once** | retry until acknowledged | never | can be | most things |
| **"exactly-once"** | at-least-once + idempotent handling | never | no **effect** twice | payments, orders |

- > Exactly-once **delivery** over a network is impossible. An exactly-once **effect** is what you build.

---

## Retry politely: backoff and jitter

> **Problem:** a server is slow for a second. Thousands of clients time out — and all retry at the same moment.

<svg class="uml d " width="1000" height="172" viewBox="0 0 1000 172"><text class="cap st" x="20" y="54" style="fill:#f0a06e">NO JITTER</text><path class="ln " d="M190 58L820 58"/><path class="ln " d="M340 55L340 61"/><path class="ln " d="M490 55L490 61"/><path class="ln " d="M790 55L790 61"/><text class="cap st" x="20" y="136" style="fill:#7ed0a3">BACKOFF + JITTER</text><path class="ln " d="M190 140L820 140"/><path class="ln " d="M340 137L340 143"/><path class="ln " d="M490 137L490 143"/><path class="ln " d="M790 137L790 143"/><text class="s" x="340" y="166">1 s</text><text class="s" x="490" y="166">2 s</text><text class="s" x="790" y="166">4 s</text><g class="pop" style="--d:0.3s"><circle class="pkt b" cx="340" cy="50" r="4"/></g><g class="pop" style="--d:0.3s"><circle class="pkt b" cx="340" cy="42" r="4"/></g><g class="pop" style="--d:0.3s"><circle class="pkt b" cx="340" cy="34" r="4"/></g><g class="pop" style="--d:0.3s"><circle class="pkt b" cx="340" cy="26" r="4"/></g><g class="pop" style="--d:0.3s"><circle class="pkt b" cx="340" cy="18" r="4"/></g><g class="pop" style="--d:0.3s"><circle class="pkt b" cx="340" cy="10" r="4"/></g><g class="pop" style="--d:0.8s"><circle class="pkt b" cx="490" cy="50" r="4"/></g><g class="pop" style="--d:0.8s"><circle class="pkt b" cx="490" cy="42" r="4"/></g><g class="pop" style="--d:0.8s"><circle class="pkt b" cx="490" cy="34" r="4"/></g><g class="pop" style="--d:0.8s"><circle class="pkt b" cx="490" cy="26" r="4"/></g><g class="pop" style="--d:0.8s"><circle class="pkt b" cx="490" cy="18" r="4"/></g><g class="pop" style="--d:0.8s"><circle class="pkt b" cx="490" cy="10" r="4"/></g><g class="pop" style="--d:1.3s"><circle class="pkt b" cx="790" cy="50" r="4"/></g><g class="pop" style="--d:1.3s"><circle class="pkt b" cx="790" cy="42" r="4"/></g><g class="pop" style="--d:1.3s"><circle class="pkt b" cx="790" cy="34" r="4"/></g><g class="pop" style="--d:1.3s"><circle class="pkt b" cx="790" cy="26" r="4"/></g><g class="pop" style="--d:1.3s"><circle class="pkt b" cx="790" cy="18" r="4"/></g><g class="pop" style="--d:1.3s"><circle class="pkt b" cx="790" cy="10" r="4"/></g><g class="pop" style="--d:2.0s"><circle class="pkt r" cx="235" cy="132" r="4"/></g><g class="pop" style="--d:2.08s"><circle class="pkt r" cx="310" cy="124" r="4"/></g><g class="pop" style="--d:2.16s"><circle class="pkt r" cx="272.5" cy="132" r="4"/></g><g class="pop" style="--d:2.24s"><circle class="pkt r" cx="205" cy="124" r="4"/></g><g class="pop" style="--d:2.32s"><circle class="pkt r" cx="325" cy="132" r="4"/></g><g class="pop" style="--d:2.4s"><circle class="pkt r" cx="287.5" cy="124" r="4"/></g><g class="pop" style="--d:2.5s"><circle class="pkt r" cx="250" cy="132" r="4"/></g><g class="pop" style="--d:2.58s"><circle class="pkt r" cx="445" cy="124" r="4"/></g><g class="pop" style="--d:2.66s"><circle class="pkt r" cx="355" cy="132" r="4"/></g><g class="pop" style="--d:2.74s"><circle class="pkt r" cx="227.5" cy="124" r="4"/></g><g class="pop" style="--d:2.82s"><circle class="pkt r" cx="400" cy="132" r="4"/></g><g class="pop" style="--d:2.9s"><circle class="pkt r" cx="482.5" cy="124" r="4"/></g><g class="pop" style="--d:3.0s"><circle class="pkt r" cx="580" cy="132" r="4"/></g><g class="pop" style="--d:3.08s"><circle class="pkt r" cx="325" cy="124" r="4"/></g><g class="pop" style="--d:3.16s"><circle class="pkt r" cx="715" cy="132" r="4"/></g><g class="pop" style="--d:3.24s"><circle class="pkt r" cx="460" cy="124" r="4"/></g><g class="pop" style="--d:3.32s"><circle class="pkt r" cx="775" cy="132" r="4"/></g><g class="pop" style="--d:3.4s"><circle class="pkt r" cx="655" cy="124" r="4"/></g><text class="s bd st" x="828" y="40" style="font-size:12px">everyone at the same moments</text><text class="s gd st" x="828" y="136" style="font-size:12px">spread out</text></svg>

- **Exponential backoff**: wait longer after each failure — 0.1, 0.2, 0.4, 0.8 s… **Jitter**: wait a random part of it, so clients spread out ("full jitter", AWS)
- Retry only what is **idempotent** — and give up after a few tries

---

<!-- _class: lead invert -->

# A Web Application, Step by Step

### A database, a backend — then more of each

---

## Where does the data live? A database

> **Problem:** the link shortener keeps its links in a `std::unordered_map`. Restart the server: every link is gone. Start a second server: it has its own, different map.

<svg class="uml d " width="1000" height="150" viewBox="0 0 1000 150"><path class="ln " d="M248 44.6L592 69.4"/><polygon class="hd " points="600,70 590.7,74.3 591.4,64.4"/><polygon class="hd " points="240,44 249.3,39.7 248.6,49.6"/><path class="ln " d="M248 109.6L592 90.4"/><polygon class="hd " points="600,90 591.3,95.4 590.7,85.6"/><polygon class="hd " points="240,110 248.7,104.6 249.3,114.4"/><circle class="pkt  fly" cx="240" cy="44" r="6" style="--dx:360px;--dy:26px;--t:3s;--d:0s"/><circle class="pkt r fly" cx="600" cy="90" r="6" style="--dx:-360px;--dy:20px;--t:3s;--d:1.5s"/><text class="m ac" x="420" y="46" style="font-size:12px">write: abc123 → https://…</text><text class="m ac" x="420" y="124" style="font-size:12px">read: abc123?</text><rect class="bx " x="40" y="20" width="200" height="48" rx="3"/><text class="t " x="140" y="49">backend server 1</text><rect class="bx " x="40" y="86" width="200" height="48" rx="3"/><text class="t " x="140" y="115">backend server 2</text><path class="bx a" d="M600 23L600 129A110 9 0 0 0 820 129L820 23"/><ellipse class="bx a" cx="710" cy="23" rx="110" ry="9"/><text class="t" x="710" y="78.5">database</text><text class="s" x="710" y="95.5">on disk · shared by both</text><text class="s" x="900" y="70">survives</text><text class="s" x="900" y="88">restarts</text></svg>

- A **database** is a separate server: it keeps the data **on disk**, lets many clients read and write **at the same time**, and answers questions about the data
- Examples: PostgreSQL and MySQL (tables and SQL), MongoDB (documents)

---

## What a database promises: ACID

> **Problem:** a transfer is two writes: A −100, B +100. A crash between them — or another client in the middle — and money appears or disappears.

<svg class="uml d " width="900" height="104.4" viewBox="0 0 1000 116"><rect class="grp" x="40" y="22" width="760" height="62" rx="6"/><text class="cap" x="420" y="16">ONE TRANSACTION</text><path class="ln " d="M180 53L222 53"/><polygon class="hd " points="230,53 221,58 221,48"/><path class="ln " d="M390 53L432 53"/><polygon class="hd " points="440,53 431,58 431,48"/><path class="ln " d="M600 53L642 53"/><polygon class="hd " points="650,53 641,58 641,48"/><rect class="bx " x="60" y="31" width="120" height="44" rx="3"/><text class="t m " x="120" y="58">BEGIN</text><rect class="bx " x="230" y="31" width="160" height="44" rx="3"/><text class="t m " x="310" y="58">A: −100</text><rect class="bx " x="440" y="31" width="160" height="44" rx="3"/><text class="t m " x="520" y="58">B: +100</text><rect class="bx g" x="650" y="31" width="130" height="44" rx="3"/><text class="t m " x="715" y="58">COMMIT</text><text class="t gd" x="870" y="49">both</text><text class="s gd" x="870" y="67">happened</text><g class="pop" style="--d:1.2s"><g class=""><path class="ln x" d="M408 46L422 60M422 46L408 60"/></g><text class="s bd" x="415" y="106" style="font-size:13px">crash here → the database undoes A: as if nothing happened</text></g></svg>

| | The promise | You know it from |
|---|---|---|
| **A**tomicity | the whole transaction, or none of it | Day 2: all-or-nothing |
| **C**onsistency | your rules always hold: unique codes, balance ≥ 0 | Day 2: invariants |
| **I**solation | others never see a half-done transaction | Day 3: races and locks |
| **D**urability | after "commit", it survives a crash — it's on disk | — |

<p class="note">ACID: Härder and Reuter, 1983. Isolation has levels — PostgreSQL's default, Read Committed, is weaker than "one at a time".</p>

---

## A web application: frontend, backend, database

<svg class="uml d " width="1000" height="218" viewBox="0 0 1000 218"><path class="ln " d="M197.5 97.3L352.5 42.7"/><polygon class="hd " points="360,40 353.2,47.7 349.9,38.3"/><polygon class="hd " points="190,100 196.8,92.3 200.1,101.7"/><path class="ln " d="M197.7 124.1L352.3 165.9"/><polygon class="hd " points="360,168 350,170.4 352.6,160.9"/><polygon class="hd " points="190,122 200,119.6 197.4,129.1"/><path class="ln " d="M588 168L732 168"/><polygon class="hd " points="740,168 731,172.9 731,163.1"/><polygon class="hd " points="580,168 589,163.1 589,172.9"/><text class="s" x="262" y="56">the page</text><text class="s" x="262" y="166">the API: JSON</text><text class="s" x="660" y="158">SQL</text><circle class="pkt  hop" cx="190" cy="100" r="6" style="--dx:170px;--dy:-60px;--t:6s;--d:0s"/><circle class="pkt r hop" cx="360" cy="40" r="6" style="--dx:-170px;--dy:60px;--t:6s;--d:0.9s"/><circle class="pkt  hop" cx="190" cy="122" r="6" style="--dx:170px;--dy:46px;--t:6s;--d:1.8s"/><circle class="pkt  hop" cx="580" cy="168" r="6" style="--dx:160px;--dy:0px;--t:6s;--d:2.7s"/><circle class="pkt r hop" cx="740" cy="168" r="6" style="--dx:-160px;--dy:0px;--t:6s;--d:3.6s"/><circle class="pkt r hop" cx="360" cy="168" r="6" style="--dx:-170px;--dy:-46px;--t:6s;--d:4.5s"/><rect class="bx " x="20" y="83" width="170" height="56" rx="3"/><text class="t " x="105" y="109">browser / app</text><text class="s" x="105" y="126">on the user's device</text><rect class="bx " x="360" y="12" width="220" height="56" rx="3"/><text class="t " x="470" y="38">frontend server</text><text class="s" x="470" y="55">pages, JavaScript, images</text><rect class="bx a" x="360" y="140" width="220" height="56" rx="3"/><text class="t " x="470" y="166">backend server</text><text class="s" x="470" y="183">your code: the API</text><path class="bx " d="M740 135L740 203A110 9 0 0 0 960 203L960 135"/><ellipse class="bx " cx="850" cy="135" rx="110" ry="9"/><text class="t" x="850" y="171.5">database</text><text class="s" x="850" y="188.5">the data · ACID</text></svg>

- **Frontend**: what the user sees — a page or an app. A frontend server only sends its files (HTML, JavaScript, images), often from a CDN
- **Backend**: your code — the API. It checks the rules and talks to the database: Day 2's `POST /urls` and `GET /{code}`
- **Database**: the data, and its ACID promises

---

## The app gets popular: two ways to grow

> **Problem:** 100 million users a day (Day 2). One backend and one database can't keep up.

| | **Scale up**: a bigger machine | **Scale out**: more machines |
|---|---|---|
| Your code | stays the same | must work on many machines |
| The limit | the biggest machine you can buy | hardly any |
| The price | grows faster than the power | grows with the number of machines |
| One machine dies | everything stops | the others carry on |
| New problems | none | everything in today's lecture |

- > Usual path: scale up first. Scale out when one machine is too small — or one failure would cost too much.

---

<!-- _class: lead invert -->

# Caching

### The cheapest way to serve more reads

---

## A cache: keep hot answers close

> **Problem:** every request asks the one database — about 1000 reads for every write (Day 2), mostly for the same few popular links.

<svg class="uml d " width="1000" height="148" viewBox="0 0 1000 148"><path class="ln " d="M197.9 63L422.1 35"/><polygon class="hd " points="430,34 421.7,40 420.5,30.2"/><polygon class="hd " points="190,64 198.3,58 199.5,67.8"/><path class="ln " d="M197.9 87L422.1 115"/><polygon class="hd " points="430,116 420.5,119.8 421.7,110"/><polygon class="hd " points="190,86 199.5,82.2 198.3,92"/><circle class="pkt  hop" cx="190" cy="64" r="6" style="--dx:240px;--dy:-30px;--t:6s;--d:0s"/><circle class="pkt b hop" cx="430" cy="34" r="6" style="--dx:-240px;--dy:30px;--t:6s;--d:0.9s"/><circle class="pkt  hop" cx="190" cy="86" r="6" style="--dx:240px;--dy:30px;--t:6s;--d:1.8s"/><circle class="pkt r hop" cx="430" cy="116" r="6" style="--dx:-240px;--dy:-30px;--t:6s;--d:2.7s"/><circle class="pkt r hop" cx="190" cy="64" r="6" style="--dx:240px;--dy:-30px;--t:6s;--d:3.6s"/><circle class="pkt  hop" cx="190" cy="64" r="6" style="--dx:240px;--dy:-30px;--t:6s;--d:4.5s"/><circle class="pkt r hop" cx="430" cy="34" r="6" style="--dx:-240px;--dy:30px;--t:6s;--d:5.3s"/><rect class="bx " x="20" y="50" width="170" height="50" rx="3"/><text class="t " x="105" y="80">app server</text><rect class="bx a" x="430" y="8" width="200" height="50" rx="3"/><text class="t " x="530" y="31">cache</text><text class="s" x="530" y="48">in memory · fast</text><rect class="bx " x="430" y="92" width="200" height="50" rx="3"/><text class="t " x="530" y="115">database</text><text class="s" x="530" y="132">on disk · slower</text><text class="cap st" x="670" y="30">1</text><text class="st bd" x="692" y="30" style="font-size:14px">ask the cache — a miss</text><text class="cap st" x="670" y="60">2</text><text class="st " x="692" y="60" style="font-size:14px">ask the database</text><text class="cap st" x="670" y="90">3</text><text class="st gd" x="692" y="90" style="font-size:14px">put the answer in the cache</text><text class="cap st" x="670" y="120">4</text><text class="st gd" x="692" y="120" style="font-size:14px">next time: a hit</text></svg>

- A **cache** keeps a copy of hot data in memory, close to whoever reads it. A **hit** is fast; a **miss** goes to the database, then keeps the answer
- **Hit ratio**: the share of reads the cache answers. 1000 reads per write → most reads can be hits
- The price: **stale data** — the copy can be older than the database

---

## Where caches live

<svg class="uml d " width="1000" height="152" viewBox="0 0 1000 152"><path class="ln " d="M178 64L207 64"/><polygon class="hd " points="215,64 206,69 206,59"/><polygon class="hd " points="170,64 179,59 179,69"/><path class="ln " d="M373 64L412 64"/><polygon class="hd " points="420,64 411,69 411,59"/><polygon class="hd " points="365,64 374,59 374,69"/><path class="ln " d="M598 64L637 64"/><polygon class="hd " points="645,64 636,69 636,59"/><polygon class="hd " points="590,64 599,59 599,69"/><path class="ln " d="M823 64L842 64"/><polygon class="hd " points="850,64 841,69 841,59"/><polygon class="hd " points="815,64 824,59 824,69"/><circle class="pkt  seg" cx="90" cy="64" r="6" style="--dx:200px;--dy:0px;--t:12s;--d:0s"/><circle class="pkt r seg" cx="290" cy="64" r="6" style="--dx:-200px;--dy:0px;--t:12s;--d:0.72s"/><circle class="pkt  seg" cx="90" cy="64" r="6" style="--dx:200px;--dy:0px;--t:12s;--d:1.8s"/><circle class="pkt  seg" cx="290" cy="64" r="6" style="--dx:215px;--dy:0px;--t:12s;--d:2.52s"/><circle class="pkt  seg" cx="505" cy="64" r="6" style="--dx:225px;--dy:0px;--t:12s;--d:3.24s"/><circle class="pkt  seg" cx="730" cy="64" r="6" style="--dx:185px;--dy:0px;--t:12s;--d:3.96s"/><circle class="pkt r seg" cx="915" cy="64" r="6" style="--dx:-185px;--dy:0px;--t:12s;--d:4.68s"/><circle class="pkt r seg" cx="730" cy="64" r="6" style="--dx:-225px;--dy:0px;--t:12s;--d:5.4s"/><circle class="pkt r seg" cx="505" cy="64" r="6" style="--dx:-215px;--dy:0px;--t:12s;--d:6.12s"/><circle class="pkt r seg" cx="290" cy="64" r="6" style="--dx:-200px;--dy:0px;--t:12s;--d:6.84s"/><rect class="bx " x="10" y="36" width="160" height="56" rx="3"/><text class="t " x="90" y="62">browser cache</text><text class="s" x="90" y="79">on the device</text><rect class="bx " x="215" y="36" width="150" height="56" rx="3"/><text class="t " x="290" y="62">CDN</text><text class="s" x="290" y="79">near the user</text><rect class="bx " x="420" y="36" width="170" height="56" rx="3"/><text class="t " x="505" y="62">backend memory</text><text class="s" x="505" y="79">each copy its own</text><rect class="bx " x="645" y="36" width="170" height="56" rx="3"/><text class="t " x="730" y="62">shared cache</text><text class="s" x="730" y="79">Redis, Memcached</text><path class="bx " d="M850 33L850 97A70 9 0 0 0 990 97L990 33"/><ellipse class="bx " cx="920" cy="33" rx="70" ry="9"/><text class="t" x="920" y="67.5">database</text><text class="s" x="920" y="84.5">the truth</text><path class="ln a" d="M980 126L28 126"/><polygon class="hd a" points="20,126 29,121 29,130.9"/><text class="s ac" x="500" y="146">closer to the user: faster — and harder to keep fresh</text></svg>

- **Browser and CDN** keep files and whole pages. The server says for how long: `Cache-Control: max-age=3600`
- **Backend memory**: the fastest — but every backend copy has its own, and they disagree
- **Shared cache** (Redis, Memcached): one cache for all backends, one network hop away

---

## Caching strategies

| Strategy | How it works | The price |
|---|---|---|
| **Cache-aside** | the app reads the cache; on a miss it reads the database and fills the cache. On a write it updates the database and **deletes** the cache entry | the first read after a change is slow |
| **Read-through** | the app asks only the cache; the cache itself loads misses from the database | the cache must know the database |
| **Write-through** | every write goes to the cache **and** the database, together | each write waits for both |
| **Write-behind** | write to the cache, "OK" at once; the database gets it later | a cache crash loses those writes |

- Most apps: **cache-aside**. On a write, **delete** the entry instead of updating it — a delete is idempotent, and two racing updates can't leave the wrong value behind (Day 3)

---

## Keeping a cache fresh

> **Problem:** someone changes a link. The cache keeps returning the old one — for how long?

- **TTL** (time to live): every entry expires after, say, an hour — stale for at most an hour
- **Invalidate** on write: delete the entry when the data changes; the next read refills it
- **Eviction**: memory is full → throw out the **least recently used** entry (LRU)

<svg class="uml d " width="1000" height="104" viewBox="0 0 1000 104"><text class="cap st" x="40" y="18">MOST RECENTLY USED</text><text class="cap en" x="960" y="18">LEAST RECENTLY USED</text><g class="pop" style="--d:0.4s"><rect class="bx a" x="20" y="30" width="130" height="44" rx="3"/><text class="t m " x="85" y="50">x9k</text><text class="s" x="85" y="67">just read</text></g><rect class="bx " x="190" y="30" width="130" height="44" rx="3"/><text class="t m " x="255" y="57">abc123</text><rect class="bx " x="340" y="30" width="130" height="44" rx="3"/><text class="t m " x="405" y="57">zz7</text><rect class="bx " x="490" y="30" width="130" height="44" rx="3"/><text class="t m " x="555" y="57">m0n</text><rect class="bx " x="640" y="30" width="130" height="44" rx="3"/><text class="t m " x="705" y="57">tt5</text><rect class="bx " x="790" y="30" width="130" height="44" rx="3"/><text class="t m " x="855" y="57">q42</text><g class="pop" style="--d:1.2s"><g class=""><path class="ln x" d="M829 26L881 78M881 26L829 78"/></g><text class="s bd" x="855" y="96">memory full: thrown out</text></g><path class="ln a" d="M160 52L178 52"/><polygon class="hd a" points="186,52 177,57 177,47"/></svg>

- **Stampede**: a hot entry expires and 1000 requests miss at once — all hit the database. Let one refill it while the others wait (Day 4: call once)

---

<!-- _class: lead invert -->

# Load Balancing

### More copies of the backend

---

## Load balancer: one address, many backends

> **Problem:** scale out: three copies of the backend. Users know one address — who picks the copy, and stops sending to a dead one?

<svg class="uml d " width="880" height="190.1" viewBox="0 0 1000 216"><path class="ln " d="M130 40L213.7 105.1"/><polygon class="hd " points="220,110 209.9,108.4 215.9,100.6"/><circle class="pkt  fly" cx="130" cy="40" r="6" style="--dx:90px;--dy:70px;--t:3s;--d:0.0s"/><path class="ln " d="M130 110L212 110"/><polygon class="hd " points="220,110 211,115 211,105"/><circle class="pkt  fly" cx="130" cy="110" r="6" style="--dx:90px;--dy:0px;--t:3s;--d:1.0s"/><path class="ln " d="M130 180L213.7 114.9"/><polygon class="hd " points="220,110 215.9,119.4 209.9,111.6"/><circle class="pkt  fly" cx="130" cy="180" r="6" style="--dx:90px;--dy:-70px;--t:3s;--d:2.0s"/><path class="ln " d="M370 110L493.1 38"/><polygon class="hd " points="500,34 494.7,42.8 489.7,34.3"/><path class="ln " d="M670 34L800 110"/><path class="ln " d="M370 110L492 110"/><polygon class="hd " points="500,110 491,115 491,105"/><path class="ln " d="M670 110L800 110"/><path class="ln b" d="M370 110L493.1 182"/><polygon class="hd b" points="500,186 489.7,185.7 494.7,177.2"/><path class="ln b" d="M670 186L800 110"/><g class=""><path class="ln x" d="M428 141L442 155M442 141L428 155"/></g><circle class="pkt  fly" cx="370" cy="110" r="6" style="--dx:130px;--dy:-76px;--t:2s;--d:0.4s"/><circle class="pkt  fly" cx="370" cy="110" r="6" style="--dx:130px;--dy:0px;--t:2s;--d:1.4s"/><circle class="pkt  fly" cx="670" cy="34" r="6" style="--dx:130px;--dy:76px;--t:2s;--d:1.2s"/><circle class="pkt  fly" cx="670" cy="110" r="6" style="--dx:130px;--dy:0px;--t:2s;--d:0.2s"/><rect class="bx " x="20" y="20" width="110" height="40" rx="3"/><text class="t " x="75" y="45">client</text><rect class="bx " x="20" y="90" width="110" height="40" rx="3"/><text class="t " x="75" y="115">client</text><rect class="bx " x="20" y="160" width="110" height="40" rx="3"/><text class="t " x="75" y="185">client</text><rect class="bx a" x="220" y="75" width="150" height="70" rx="3"/><text class="t " x="295" y="108">load balancer</text><text class="s" x="295" y="125">health checks</text><rect class="bx " x="500" y="10" width="170" height="48" rx="3"/><text class="t " x="585" y="39">server 1</text><rect class="bx " x="500" y="86" width="170" height="48" rx="3"/><text class="t " x="585" y="115">server 2</text><rect class="bx b" x="500" y="162" width="170" height="48" rx="3"/><text class="t " x="585" y="184">server 3</text><text class="s" x="585" y="201">health check failed</text><rect class="bx " x="800" y="86" width="170" height="48" rx="3"/><text class="t " x="885" y="115">database</text></svg>

- The **load balancer** owns the one public address and passes each request to one of the copies
- It **health-checks** every backend and skips the dead ones — Day 4's thread pool, with machines as workers
- Run **two** load balancers — otherwise it is the new single point of failure

---

## How does the balancer pick a server?

| Algorithm | Picks | Good when |
|---|---|---|
| **Round-robin** | 1, 2, 3, 1, 2, 3, … | requests cost about the same |
| **Weighted** | a bigger machine gets more | the machines differ |
| **Least connections** | the copy with the fewest open requests | some requests are slow |
| **Hash of the user** (IP, user ID) | the same copy for the same user | a copy keeps something per user ("sticky") |

- A **transport-level** balancer (called L4) sees only IP addresses and ports and passes whole TCP connections — fast
- An **application-level** balancer (L7) reads the HTTP request, so `/api` and `/images` can go to different servers. nginx and HAProxy do both

<p class="note">L4 and L7 are the OSI layer numbers; in our 4-layer stack they are the transport and application layers.</p>

---

## Copies need a stateless backend

> **Problem:** the balancer sends Ann's first request to copy 1, the next one to copy 3. If copy 1 kept her session in memory, copy 3 has never heard of her.

<svg class="uml d " width="1000" height="212" viewBox="0 0 1000 212"><line class="lane" x1="500" y1="10" x2="500" y2="205"/><text class="cap" x="250" y="16">STATELESS</text><path class="ln " d="M250 66L100 108"/><path class="ln " d="M100 146L250 172"/><path class="ln " d="M250 66L250 108"/><path class="ln " d="M250 146L250 172"/><path class="ln " d="M250 66L400 108"/><path class="ln " d="M400 146L250 172"/><circle class="pkt  fly" cx="250" cy="66" r="6" style="--dx:-150px;--dy:42px;--t:3s;--d:0.0s"/><circle class="pkt  fly" cx="250" cy="66" r="6" style="--dx:0px;--dy:42px;--t:3s;--d:1.0s"/><circle class="pkt  fly" cx="250" cy="66" r="6" style="--dx:150px;--dy:42px;--t:3s;--d:2.0s"/><rect class="bx " x="180" y="28" width="140" height="38" rx="3"/><text class="t " x="250" y="52">load balancer</text><rect class="bx " x="40" y="108" width="120" height="38" rx="3"/><text class="t " x="100" y="132">server 1</text><rect class="bx " x="190" y="108" width="120" height="38" rx="3"/><text class="t " x="250" y="132">server 2</text><rect class="bx " x="340" y="108" width="120" height="38" rx="3"/><text class="t " x="400" y="132">server 3</text><rect class="bx " x="160" y="172" width="180" height="36" rx="3"/><text class="t " x="250" y="195">database: sessions</text><text class="cap" x="750" y="16">STATEFUL</text><path class="ln " d="M750 66L600 108"/><path class="ln " d="M750 66L750 108"/><path class="ln " d="M750 66L900 108"/><circle class="pkt  fly" cx="750" cy="66" r="6" style="--dx:0px;--dy:42px;--t:3s;--d:0.0s"/><circle class="pkt  fly" cx="750" cy="66" r="6" style="--dx:0px;--dy:42px;--t:3s;--d:1.0s"/><circle class="pkt  fly" cx="750" cy="66" r="6" style="--dx:0px;--dy:42px;--t:3s;--d:2.0s"/><rect class="bx " x="680" y="28" width="140" height="38" rx="3"/><text class="t " x="750" y="52">load balancer</text><rect class="bx " x="540" y="108" width="120" height="38" rx="3"/><text class="t " x="600" y="132">server 1</text><rect class="bx a" x="690" y="108" width="120" height="38" rx="3"/><text class="t " x="750" y="132">server 2</text><rect class="bx " x="840" y="108" width="120" height="38" rx="3"/><text class="t " x="900" y="132">server 3</text><text class="s ac" x="750" y="166" style="font-size:13px">Ann's session lives in server 2</text><text class="s bd" x="750" y="192" style="font-size:13px">server 2 dies → Ann's session is gone</text></svg>

- **Stateless**: nothing kept between requests — the state lives in the database or the shared cache. Any copy can answer; a crash loses nothing
- **Stateful**: sessions in memory, so the balancer must keep the user on one copy (sticky). When that copy crashes, the session is gone

---

<!-- _class: lead invert -->

# Database Replication

### More copies of the database

---

## Copies of the database: master and slaves

> **Problem:** one database machine. If it dies, every link is gone — and every read still waits for it.

<svg class="uml d " width="1000" height="192" viewBox="0 0 1000 192"><path class="ln " d="M160 102L272 102"/><polygon class="hd " points="280,102 271,107 271,97"/><path class="ln " d="M460 92L592.5 40.9"/><polygon class="hd " points="600,38 593.4,45.9 589.8,36.6"/><path class="ln " d="M460 112L592.5 159.3"/><polygon class="hd " points="600,162 589.9,163.6 593.2,154.3"/><path class="ln " d="M853.3 87.6L786.7 44.4"/><polygon class="hd " points="780,40 790.2,40.8 784.8,49.1"/><polygon class="hd " points="860,92 849.8,91.2 855.2,82.9"/><path class="ln " d="M853.1 116.1L786.9 155.9"/><polygon class="hd " points="780,160 785.2,151.1 790.3,159.6"/><polygon class="hd " points="860,112 854.8,120.9 849.7,112.4"/><text class="s" x="220" y="92">writes</text><text class="s" x="540" y="106" style="font-size:12px">a moment later</text><text class="s" x="812" y="106">reads</text><circle class="pkt  fly" cx="160" cy="102" r="6" style="--dx:120px;--dy:0px;--t:4s;--d:0s"/><circle class="pkt  fly" cx="460" cy="92" r="6" style="--dx:140px;--dy:-54px;--t:4s;--d:1.7s"/><circle class="pkt  fly" cx="460" cy="112" r="6" style="--dx:140px;--dy:50px;--t:4s;--d:2.4s"/><circle class="pkt r fly" cx="780" cy="40" r="6" style="--dx:80px;--dy:52px;--t:4s;--d:0.5s"/><circle class="pkt r fly" cx="780" cy="160" r="6" style="--dx:80px;--dy:-48px;--t:4s;--d:2.5s"/><rect class="bx " x="20" y="80" width="140" height="44" rx="3"/><text class="t " x="90" y="107">writer</text><rect class="bx a" x="280" y="72" width="180" height="60" rx="3"/><text class="t " x="370" y="100">master</text><text class="s" x="370" y="117">all writes</text><rect class="bx " x="600" y="14" width="180" height="48" rx="3"/><text class="t " x="690" y="43">slave 1</text><rect class="bx " x="600" y="138" width="180" height="48" rx="3"/><text class="t " x="690" y="167">slave 2</text><rect class="bx " x="860" y="78" width="120" height="48" rx="3"/><text class="t " x="920" y="107">readers</text></svg>

- **Master–slave**: all writes go to the **master**; it sends every change to its **slaves**; reads can go to any copy
- **Synchronous**: the master waits for a slave before "OK" — safe, slower. **Asynchronous**: "OK" at once — fast, but a master that dies may lose its last writes

<p class="note">Newer names for the same idea: primary and replica, leader and follower.</p>

---

## Reads to the slaves, writes to the master

> **Problem:** 1000 reads for every write — and the master still answers all of them.

<svg class="uml d " width="1000" height="206" viewBox="0 0 1000 206"><rect class="grp" x="600" y="96" width="390" height="106" rx="6"/><text class="cap" x="795" y="92">SLAVES</text><path class="ln a" d="M190 88L352.3 44.1"/><polygon class="hd a" points="360,42 352.6,49.1 350,39.6"/><path class="ln g" d="M190 112L590 149.3"/><polygon class="hd g" points="598,150 588.6,154.1 589.5,144.2"/><path class="ln " d="M540 56L632.9 104.3"/><polygon class="hd " points="640,108 629.7,108.2 634.3,99.5"/><text class="s ac" x="262" y="54">writes</text><text class="s gd" x="400" y="148">reads: any slave</text><text class="s" x="640" y="66">every change</text><circle class="pkt  fly" cx="190" cy="88" r="6" style="--dx:170px;--dy:-46px;--t:3s;--d:0s"/><circle class="pkt r fly" cx="190" cy="112" r="6" style="--dx:408px;--dy:38px;--t:3s;--d:0.9s"/><circle class="pkt  fly" cx="540" cy="56" r="6" style="--dx:100px;--dy:52px;--t:3s;--d:1.6s"/><rect class="bx " x="20" y="72" width="170" height="56" rx="3"/><text class="t " x="105" y="98">backend</text><text class="s" x="105" y="115">reads ≫ writes</text><path class="bx a" d="M360 15L360 71A90 9 0 0 0 540 71L540 15"/><ellipse class="bx a" cx="450" cy="15" rx="90" ry="9"/><text class="t" x="450" y="45.5">master</text><text class="s" x="450" y="62.5">all writes</text><path class="bx " d="M615 125L615 179A85 9 0 0 0 785 179L785 125"/><ellipse class="bx " cx="700" cy="125" rx="85" ry="9"/><text class="t" x="700" y="154.5">slave 1</text><text class="s" x="700" y="171.5">reads</text><path class="bx " d="M805 125L805 179A85 9 0 0 0 975 179L975 125"/><ellipse class="bx " cx="890" cy="125" rx="85" ry="9"/><text class="t" x="890" y="154.5">slave 2</text><text class="s" x="890" y="171.5">reads</text></svg>

- The backend sends **writes to the master** and **reads to a slave** — the database driver or a proxy can route them
- More slaves → more reads. Writes don't scale this way: every write still goes to the one master — and to every slave
- Slaves also take backups and heavy reports, away from the master

---

<!-- _class: invert bad -->

## Replication lag: where did my comment go?

<svg class="uml d tl" width="1000" height="30" viewBox="0 0 1000 30"><text class="cap" x="150" y="20">USER</text><text class="cap" x="500" y="20">MASTER</text><text class="cap" x="850" y="20">SLAVE</text></svg>

* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="150" y1="0" x2="150" y2="44"/><line class="lane" x1="500" y1="0" x2="500" y2="44"/><line class="lane" x1="850" y1="0" x2="850" y2="44"/><text class="s" x="14" y="33">1</text><path class="ln  draw" pathLength="1" d="M155 29L487 29" style="--d:0s"/><g class="late" style="--d:0s"><polygon class="hd " points="495,29 486,34 486,24.1"/></g><circle class="pkt  once" cx="155" cy="29" r="5" style="--dx:332px;--dy:0px;--t:3s;--d:0s"/><g class="late" style="--d:-0.55s"><text class="m " x="325" y="20">post comment</text></g></svg>
* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="150" y1="0" x2="150" y2="44"/><line class="lane" x1="500" y1="0" x2="500" y2="44"/><line class="lane" x1="850" y1="0" x2="850" y2="44"/><text class="s" x="14" y="33">2</text><path class="ln  draw" pathLength="1" d="M495 29L163 29" style="--d:0s"/><g class="late" style="--d:0s"><polygon class="hd " points="155,29 164,24.1 164,34"/></g><circle class="pkt r once" cx="495" cy="29" r="5" style="--dx:-332px;--dy:0px;--t:3s;--d:0s"/><g class="late" style="--d:-0.55s"><text class="m " x="325" y="20">OK</text></g></svg>
* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="150" y1="0" x2="150" y2="44"/><line class="lane" x1="500" y1="0" x2="500" y2="44"/><line class="lane" x1="850" y1="0" x2="850" y2="44"/><text class="s" x="14" y="33">3</text><path class="ln  draw" pathLength="1" d="M155 29L837 29" style="--d:0s"/><g class="late" style="--d:0s"><polygon class="hd " points="845,29 836,34 836,24.1"/></g><circle class="pkt  once" cx="155" cy="29" r="5" style="--dx:682px;--dy:0px;--t:3s;--d:0s"/><g class="late" style="--d:-0.55s"><text class="m " x="500" y="20">get comments</text></g></svg>
* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="150" y1="0" x2="150" y2="44"/><line class="lane" x1="500" y1="0" x2="500" y2="44"/><line class="lane" x1="850" y1="0" x2="850" y2="44"/><text class="s" x="14" y="33">4</text><path class="ln b draw" pathLength="1" d="M845 29L163 29" style="--d:0s"/><g class="late" style="--d:0s"><polygon class="hd b" points="155,29 164,24.1 164,34"/></g><circle class="pkt b once" cx="845" cy="29" r="5" style="--dx:-682px;--dy:0px;--t:3s;--d:0s"/><g class="late" style="--d:-0.55s"><text class="m bd" x="500" y="20">[ ]  — my comment is gone!</text></g></svg>
* <svg class="uml d tl" width="1000" height="44" viewBox="0 0 1000 44"><line class="lane" x1="150" y1="0" x2="150" y2="44"/><line class="lane" x1="500" y1="0" x2="500" y2="44"/><line class="lane" x1="850" y1="0" x2="850" y2="44"/><text class="s" x="14" y="33">5</text><path class="ln g draw" pathLength="1" d="M505 29L837 29" style="--d:0s"/><g class="late" style="--d:0s"><polygon class="hd g" points="845,29 836,34 836,24.1"/></g><circle class="pkt r once" cx="505" cy="29" r="5" style="--dx:332px;--dy:0px;--t:3s;--d:0s"/><g class="late" style="--d:-0.55s"><text class="m " x="675" y="20">+ comment  (1 s late)</text></g></svg>

<div class="sp"></div>

* > **Eventual consistency**: stop writing, and all copies **eventually** agree — with no promise about **when**.
* A fix for this case — **read your own writes**: read what you just changed from the master

---

## When the master dies: failover

> **Problem:** the master stops answering. Writes have nowhere to go.

<svg class="uml d " width="820" height="173.8" viewBox="0 0 1000 212"><path class="ln b" d="M180 82L392.1 47.3"/><polygon class="hd b" points="400,46 391.9,52.3 390.3,42.6"/><text class="s bd" x="290" y="52">no answer</text><rect class="bx " x="20" y="60" width="160" height="50" rx="3"/><text class="t " x="100" y="90">backend</text><path class="bx b" d="M400 15L400 71A90 9 0 0 0 580 71L580 15"/><ellipse class="bx b" cx="490" cy="15" rx="90" ry="9"/><text class="t" x="490" y="45.5">master</text><text class="s" x="490" y="62.5">stopped answering</text><g class="pop" style="--d:0.3s"><g class=""><path class="ln x" d="M588 10L612 34M612 10L588 34"/></g></g><path class="bx " d="M400 125L400 177A90 9 0 0 0 580 177L580 125"/><ellipse class="bx " cx="490" cy="125" rx="90" ry="9"/><text class="t" x="490" y="153.5">slave 1</text><text class="s" x="490" y="170.5">the newest data</text><path class="bx " d="M720 125L720 177A90 9 0 0 0 900 177L900 125"/><ellipse class="bx " cx="810" cy="125" rx="90" ry="9"/><text class="t" x="810" y="160.5">slave 2</text><g class="pop" style="--d:2.0s"><path class="ln g" d="M180 96L390.2 144.2"/><polygon class="hd g" points="398,146 388.1,148.8 390.3,139.2"/><text class="s gd" x="258" y="140">writes go here now</text></g><g class="pop" style="--d:1.2s"><rect x="396" y="112" width="188" height="78" rx="10" style="fill:none;stroke:#7ed0a3;stroke-width:3"/><text class="s gd" x="490" y="205">promoted: the new master</text></g><g class="pop" style="--d:2.6s"><path class="ln g" d="M582 151L710 151"/><polygon class="hd g" points="718,151 709,155.9 709,146.1"/><text class="s gd" x="650" y="141">follows it</text></g></svg>

* **Detect**: the master misses its heartbeats — a **timeout** says "dead" (maybe it is only slow)
* **Promote** the slave with the newest data; point the backends and the other slaves at it
* Asynchronous replication: writes the old master never sent are **lost**
* The old master returns: it must become a slave. Two masters at once (**split brain**) is the worst case

---

## The whole picture

<svg class="uml d " width="1100" height="253" viewBox="0 0 1000 230"><rect class="grp" x="290" y="10" width="170" height="216" rx="6"/><text class="cap" x="375" y="26">BACKENDS</text><path class="ln " d="M100 108L126.5 53.2"/><polygon class="hd " points="130,46 130.5,56.3 121.6,51.9"/><path class="ln " d="M100 122L122 122"/><polygon class="hd " points="130,122 121,127 121,117"/><path class="ln " d="M250 122L282 122"/><polygon class="hd " points="290,122 281,127 281,117"/><path class="ln " d="M460 56L512.1 47.3"/><polygon class="hd " points="520,46 511.9,52.4 510.3,42.6"/><path class="ln " d="M460 118L710.2 57.9"/><polygon class="hd " points="718,56 710.4,62.9 708.1,53.3"/><path class="ln g" d="M460 184L680 185.9"/><polygon class="hd g" points="688,186 679,190.9 679,181"/><path class="ln " d="M800 88L760 150"/><path class="ln " d="M800 88L920 150"/><circle class="pkt  fly" cx="100" cy="122" r="6" style="--dx:30px;--dy:0px;--t:3s;--d:0s"/><circle class="pkt  fly" cx="250" cy="122" r="6" style="--dx:40px;--dy:0px;--t:3s;--d:0.5s"/><circle class="pkt  fly" cx="460" cy="56" r="6" style="--dx:60px;--dy:-10px;--t:3s;--d:1.0s"/><circle class="pkt  fly" cx="460" cy="118" r="6" style="--dx:258px;--dy:-62px;--t:3s;--d:1.5s"/><circle class="pkt r fly" cx="460" cy="184" r="6" style="--dx:228px;--dy:2px;--t:3s;--d:2.0s"/><rect class="bx " x="10" y="92" width="90" height="44" rx="3"/><text class="t " x="55" y="119">users</text><rect class="bx " x="130" y="20" width="120" height="50" rx="3"/><text class="t " x="190" y="43">CDN</text><text class="s" x="190" y="60">files</text><rect class="bx a" x="130" y="96" width="120" height="52" rx="3"/><text class="t " x="190" y="120">load</text><text class="s" x="190" y="137">balancer ×2</text><rect class="bx " x="305" y="40" width="140" height="44" rx="3"/><text class="t " x="375" y="60">backend</text><text class="s" x="375" y="77">stateless</text><rect class="bx " x="305" y="104" width="140" height="44" rx="3"/><text class="t " x="375" y="124">backend</text><text class="s" x="375" y="141">stateless</text><rect class="bx " x="305" y="168" width="140" height="44" rx="3"/><text class="t " x="375" y="188">backend</text><text class="s" x="375" y="205">stateless</text><rect class="bx " x="520" y="22" width="150" height="48" rx="3"/><text class="t " x="595" y="44">cache</text><text class="s" x="595" y="61">Redis</text><path class="bx a" d="M720 21L720 79A80 9 0 0 0 880 79L880 21"/><ellipse class="bx a" cx="800" cy="21" rx="80" ry="9"/><text class="t" x="800" y="52.5">master</text><text class="s" x="800" y="69.5">writes</text><path class="bx " d="M690 159L690 205A70 9 0 0 0 830 205L830 159"/><ellipse class="bx " cx="760" cy="159" rx="70" ry="9"/><text class="t" x="760" y="184.5">slave</text><text class="s" x="760" y="201.5">reads</text><path class="bx " d="M850 159L850 205A70 9 0 0 0 990 205L990 159"/><ellipse class="bx " cx="920" cy="159" rx="70" ry="9"/><text class="t" x="920" y="184.5">slave</text><text class="s" x="920" y="201.5">reads</text></svg>


---

<!-- _class: lead invert -->

# Start Simple, Grow Step by Step

### A database · a cache · a load balancer · replicas — each one when you need it
