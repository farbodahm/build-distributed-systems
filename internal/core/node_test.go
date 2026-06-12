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
	var init InitMessage
	init.Body.MsgID = 7
	if got := init.RequestMsgID(); got != 7 {
		t.Errorf("InitMessage: got %d want 7", got)
	}

	var echo EchoMessage
	echo.Body.MsgID = 9
	if got := echo.RequestMsgID(); got != 9 {
		t.Errorf("EchoMessage: got %d want 9", got)
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

func TestRunDispatchesByType(t *testing.T) {
	var buf bytes.Buffer
	orig := Log.out
	Log.out = log.New(&buf, "", 0)
	defer func() { Log.out = orig }()

	feedStdin(t, `{"src":"c0","dest":"n1","body":{"type":"init","msg_id":1,"node_id":"n1","node_ids":["n1"]}}`+"\n"+
		`{"src":"c1","dest":"n1","body":{"type":"echo","msg_id":2,"echo":"hi"}}`+"\n")

	n := NewNode()
	n.OnInit(func(m InitMessage) error {
		n.Init(m.Body.NodeID, m.Body.NodeIDs)
		n.Reply(m, &InitOkBody{})
		return nil
	})
	n.OnEcho(func(m EchoMessage) error {
		n.Reply(m, &EchoOkBody{Echo: m.Body.Echo})
		return nil
	})

	n.Run()

	got := strings.Split(strings.TrimSpace(buf.String()), "\n")
	want := []string{
		`{"src":"n1","dest":"c0","body":{"type":"init_ok","msg_id":0,"in_reply_to":1}}`,
		`{"src":"n1","dest":"c1","body":{"type":"echo_ok","msg_id":1,"in_reply_to":2,"echo":"hi"}}`,
	}
	if len(got) != len(want) {
		t.Fatalf("got %d replies, want %d:\n%s", len(got), len(want), buf.String())
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("reply %d:\n got: %s\nwant: %s", i, got[i], want[i])
		}
	}

	if n.ID != "n1" {
		t.Errorf("init handler did not set node ID: got %q", n.ID)
	}
}

func TestRegisterHandlerDuplicatePanics(t *testing.T) {
	n := NewNode()
	n.OnInit(func(InitMessage) error { return nil })

	defer func() {
		if recover() == nil {
			t.Error("expected panic on duplicate handler registration")
		}
	}()
	n.OnInit(func(InitMessage) error { return nil })
}
