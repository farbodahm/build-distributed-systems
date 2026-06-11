package main

import (
	. "build-distributed-systems/internal/core"
)

func main() {
	node := NewNode()

	node.RegisterHandler(MsgTypeInit, func(msg Incoming) error {
		m := msg.(InitMessage)
		node.Init(m.Body.NodeID, m.Body.NodeIDs)
		Log.Info("initialized node %s with peers %v", node.ID, node.Peers)
		node.Reply(m, &InitOkBody{})
		return nil
	})

	node.RegisterHandler(MsgTypeEcho, func(msg Incoming) error {
		m := msg.(EchoMessage)
		Log.Info("received echo request: %s", m.Body.Echo)
		node.Reply(m, &EchoOkBody{Echo: m.Body.Echo})
		return nil
	})

	if err := node.Run(); err != nil {
		Log.Error("run: %v", err)
	}
}
