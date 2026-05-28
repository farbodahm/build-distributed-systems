package main

import (
	. "build-distributed-systems/internal/core"
)

func main() {
	var msg Message
	for ScanLine(&msg) {
		bodyType, ok := msg.Body["type"].(string)
		if !ok {
			bodyType = "unknown"
		}
		Log.Print("PARSED: %s|%s|%s", msg.Src, msg.Dest, bodyType)
		Log.Debug("src=%s dest=%s body=%v", msg.Src, msg.Dest, msg.Body)
	}
}
