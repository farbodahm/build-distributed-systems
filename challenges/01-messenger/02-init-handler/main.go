package main

import (
	. "build-distributed-systems/internal/core"
)

type Node struct {
	ID        string
	Peers     []string
	NextMsgID int
}

func main() {
	var msg InitMessage
	var node Node
	for ScanLine(&msg) {
		node.ID = msg.Body.NodeID
		node.Peers = msg.Body.NodeIDs
		Log.Info("initialized node %s with peers %v", node.ID, node.Peers)
		reply := InitOkMessage{
			Src:  msg.Dest,
			Dest: msg.Src,
		}
		reply.Body.Type = "init_ok"
		reply.Body.InReplyTo = msg.Body.MsgID
		reply.Body.MsgID = node.NextMsgID
		node.NextMsgID++

		Log.PrintJSON(reply)
	}
}
