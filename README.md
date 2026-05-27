# Build Distributed Systems

Solutions to https://builddistributedsystem.com/

## Setup

```
make install-maelstrom    # optional, for later Maelstrom-backed tracks
```

(The `submit` pipeline pulls `bundle` and `goimports` on demand via `go run`.)

## Develop

```
make run CHALLENGE=01-messenger/01-json-parser \
  < challenges/01-messenger/01-json-parser/testdata/input.txt
```

## Submit

```
make verify-submit CHALLENGE=01-messenger/01-json-parser
```

Then paste `challenges/<path>/submit.go` into the web editor.

## Layout

```
challenges/<track>/<task>/main.go   ← what you edit (dot-imports internal/core)
internal/core/                      ← one package, bundled into each submit.go
tools/                              ← helper scripts
Makefile                            ← orchestration
```

## How submit works

1. `bundle` inlines `internal/core` into a single file as `package main`.
2. `awk` strips the package + import lines from your `main.go`.
3. The two bodies are concatenated.
4. `goimports` cleans up the import block.

Result: a self-contained `submit.go` that compiles in an empty directory.
