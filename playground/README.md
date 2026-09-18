# MDCD Playground

Hands-on tasks for each lecture, implemented once per language so you can
work in whichever you're using this semester. Same problems, same rules,
same tests conceptually — just idiomatic to each language.

| Language | Setup + commands |
|---|---|
| C++ (CMake + GoogleTest) | [`cpp/README.md`](cpp/README.md) |
| Python (unittest, stdlib only) | [`python/README.md`](python/README.md) |
| Java (Maven + JUnit 5) | [`java/README.md`](java/README.md) |
| Go (stdlib `testing`) | [`go/README.md`](go/README.md) |

Each is self-contained under its own folder — pick one, follow its README,
ignore the rest. CI builds and tests all four on every push.

Each day has two tiers: a **core** set covering the lecture's ideas, and
**challenge** problems — harder, with an algorithmic component. Do the core
set first; the challenge set is where the real fun is.

A fresh clone builds/compiles cleanly in every language but **fails** most
tests — expected, the `TODO`s haven't been filled in yet.
