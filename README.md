# Build Distributed Systems

Solutions to https://builddistributedsystem.com/

## A note on AI

I'm not using LLMs to solve these. I enjoy solving them like chess.
Everything is getting super fast right now, so I expect this repo to move slower.
I enjoy the path and the learning.

## Develop

```
make run CHALLENGE=01-messenger/01-json-parser
```

## Submit

```
make verify-submit CHALLENGE=01-messenger/01-json-parser
```

Then paste `challenges/<path>/_submit.go` into the web editor.


### How submit works

1. `bundle` inlines `internal/core` into a single file as `package main`.
2. `awk` strips the package + import lines from your `main.go`.
3. The two bodies are concatenated.
4. `goimports` cleans up the import block.

Result: a self-contained `_submit.go` that compiles in an empty directory.
