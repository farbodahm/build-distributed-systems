package main

import (
	. "build-distributed-systems/internal/core"
)

func main() {
	var msg InitMessage
	for ScanLine(&msg) {
		Log.Debug("got: %+v", msg)
	}
}
