package main

import (
	. "build-distributed-systems/internal/core"
)

func main() {
	node := NewNode()

	var msg InitMessage
	for ScanLine(&msg) {
		node.Init(msg.Body.NodeID, msg.Body.NodeIDs)
		Log.Info("initialized node %s with peers %v", node.ID, node.Peers)
		node.Reply(msg, &InitOkBody{})
	}

}
