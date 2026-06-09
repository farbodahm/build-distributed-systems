package main

import (
	. "build-distributed-systems/internal/core"
)

func main() {
	node := NewNode()

	var msg Incoming
	for ScanTyped(&msg) {
		switch m := msg.(type) {
		case InitMessage:
			node.Init(m.Body.NodeID, m.Body.NodeIDs)
			Log.Info("initialized node %s with peers %v", node.ID, node.Peers)
			node.Reply(m, &InitOkBody{})
		case EchoMessage:
			node.Reply(m, &EchoOkBody{Echo: m.Body.Echo})
		}
	}
}
