package core

import (
	"bufio"
	"strings"
	"testing"
)

func feedStdin(t *testing.T, input string) {
	t.Helper()
	orig := stdin
	stdin = bufio.NewScanner(strings.NewReader(input))
	t.Cleanup(func() { stdin = orig })
}

func TestScanLine(t *testing.T) {
	feedStdin(t, `{"src":"c0","dest":"n1","body":{"type":"init","msg_id":1,"node_id":"n1","node_ids":["n1"]}}`+"\n")

	var msg InitMessage
	if !ScanLine(&msg) {
		t.Fatal("expected true")
	}
	if msg.Src != "c0" || msg.Body.NodeID != "n1" || msg.Body.MsgID != 1 {
		t.Errorf("got %+v", msg)
	}
	if ScanLine(&msg) {
		t.Error("expected false at EOF")
	}
}

func TestScanLineSkipsBadJSON(t *testing.T) {
	feedStdin(t, "not json\n"+`{"src":"c0","dest":"n1","body":{"type":"init","msg_id":2,"node_id":"n1"}}`+"\n")

	var msg InitMessage
	if !ScanLine(&msg) {
		t.Fatal("expected true after skipping bad line")
	}
	if msg.Body.MsgID != 2 {
		t.Errorf("expected msg_id 2, got %d", msg.Body.MsgID)
	}
}
