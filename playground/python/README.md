# MDCD Playground — Python

Plain Python, no third-party dependencies. Each lecture is a
`days/dayXX/` folder: modules with `TODO`s to implement, and `unittest`
cases (stdlib only — no `pip install` needed) that check them.

## Setup

Python 3.9+ and git — that's it.

- **Linux**: `sudo apt install python3 git`
- **macOS**: Python 3 is usually preinstalled; if not, `brew install python`
- **Windows**: install [Python](https://www.python.org/downloads/) (check
  "Add python.exe to PATH" during setup) and [Git](https://git-scm.com/downloads)

## Run the tests

From the `playground/python/` directory, on any OS:

```
python3 -m unittest discover -s . -p "test_*.py" -v
```

(On Windows, use `python` instead of `python3` if that's how it's on PATH.)

A fresh clone runs cleanly but **fails** most tests — expected, the
`TODO`s aren't filled in yet. Edit `days/dayXX/*.py` (the non-`test_`
files); leave `test_*.py` alone.
