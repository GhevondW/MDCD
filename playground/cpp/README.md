# MDCD Playground — C++

CMake + GoogleTest C++ project. Each lecture is a `days/dayXX/` folder:
headers with `TODO`s to implement, and GoogleTest cases that check them.
Needs only a C++17 compiler and CMake — GoogleTest is fetched by
[CPM](https://github.com/cpm-cmake/CPM.cmake) (vendored in `cmake/CPM.cmake`)
on first configure.

## Setup

Compiler + CMake (>= 3.16) + git.

- **Linux**: `sudo apt install build-essential cmake git`
- **macOS**: `xcode-select --install` (clang) + `brew install cmake`
- **Windows**: [Visual Studio Community](https://visualstudio.microsoft.com/)
  with "Desktop development with C++" (includes MSVC + CMake) — open
  `playground/cpp/` via `File > Open > Folder...`. Or CLI: CMake +
  [Build Tools for Visual Studio](https://visualstudio.microsoft.com/downloads/#build-tools-for-visual-studio-2022)
  (run from a "Developer Command Prompt for VS") or MSYS2's `mingw-w64-gcc`.

## Build and test

```
cmake -S . -B build
cmake --build build --config Release
cd build && ctest --output-on-failure -C Release
```

`--config`/`-C Release` only matters on Windows; harmless elsewhere. First
configure needs internet (fetches GoogleTest once, then it's cached).

A fresh clone builds cleanly but **fails** most tests — expected, the
`TODO`s aren't filled in yet. Edit `days/dayXX/include/**/*.h`; leave
`tests/` alone.
