# MDCD Playground — Go

Plain Go, stdlib `testing` only — no third-party dependencies. Each
lecture is a `days/dayXX/` package: files with `TODO`s to implement, and
`_test.go` files that check them.

## Setup

Go 1.21+ and git.

- **Linux**: `sudo apt install golang-go git` (or the [official installer](https://go.dev/doc/install)
  for a newer version than your distro ships)
- **macOS**: `brew install go`
- **Windows**: install [Go](https://go.dev/dl/) and [Git](https://git-scm.com/downloads)

## Run the tests

From the `playground/go/` directory, on any OS:

```
go test ./...
```

For just one day: `go test ./days/day02/...`.

A fresh clone builds cleanly but **fails** most tests — expected, the
`TODO`s aren't filled in yet. Edit `days/dayXX/*.go` (the non-`_test.go`
files); leave `*_test.go` alone.
