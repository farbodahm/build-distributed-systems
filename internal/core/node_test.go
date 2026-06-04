package core

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestNextMsgID(t *testing.T) {
	n := NewNode()
	for i := 0; i < 3; i++ {
		if got := n.NextMsgID(); got != i {
			t.Errorf("call %d: got %d want %d", i, got, i)
		}
	}
}

func TestInit(t *testing.T) {
	n := NewNode()
	n.Init("n1", []string{"n1", "n2"})
	if n.ID != "n1" {
		t.Errorf("ID: got %q want %q", n.ID, "n1")
	}
	if len(n.Peers) != 2 || n.Peers[0] != "n1" || n.Peers[1] != "n2" {
		t.Errorf("Peers: got %v want [n1 n2]", n.Peers)
	}
}

func TestReply(t *testing.T) {
	var buf bytes.Buffer
	orig := Log.out
	Log.out = log.New(&buf, "", 0)
	defer func() { Log.out = orig }()

	n := NewNode()
	var msg InitMessage
	msg.Src = "c0"
	msg.Dest = "n1"
	msg.Body.MsgID = 42

	n.Reply(msg, &InitOkBody{})

	got := strings.TrimSpace(buf.String())
	want := `{"src":"n1","dest":"c0","body":{"type":"init_ok","msg_id":0,"in_reply_to":42}}`
	if got != want {
		t.Errorf("\n got: %s\nwant: %s", got, want)
	}
}

func TestRequestMsgID(t *testing.T) {
	var msg InitMessage
	msg.Body.MsgID = 7

	if got := requestMsgID(msg); got != 7 {
		t.Errorf("value: got %d want 7", got)
	}
	if got := requestMsgID(&msg); got != 7 {
		t.Errorf("pointer: got %d want 7", got)
	}

	noBody := struct{ Src string }{Src: "c0"}
	if got := requestMsgID(noBody); got != 0 {
		t.Errorf("no Body field: got %d want 0", got)
	}

	noMsgID := struct {
		Body struct{ Type string }
	}{}
	if got := requestMsgID(noMsgID); got != 0 {
		t.Errorf("no MsgID field: got %d want 0", got)
	}
}

func TestReplyMsgIDIncrements(t *testing.T) {
	var buf bytes.Buffer
	orig := Log.out
	Log.out = log.New(&buf, "", 0)
	defer func() { Log.out = orig }()

	n := NewNode()
	var msg InitMessage
	msg.Src = "c0"
	msg.Dest = "n1"
	msg.Body.MsgID = 1

	n.Reply(msg, &InitOkBody{})
	n.Reply(msg, &InitOkBody{})

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if !strings.Contains(lines[0], `"msg_id":0`) {
		t.Errorf("line 1 missing msg_id:0: %s", lines[0])
	}
	if !strings.Contains(lines[1], `"msg_id":1`) {
		t.Errorf("line 2 missing msg_id:1: %s", lines[1])
	}
}
