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

func TestScanTyped(t *testing.T) {
	feedStdin(t, `{"src":"c0","dest":"n1","body":{"type":"init","msg_id":1,"node_id":"n1","node_ids":["n1"]}}`+"\n"+
		`{"src":"c0","dest":"n1","body":{"type":"echo","msg_id":2,"echo":"hi"}}`+"\n")

	var msg Incoming
	if !ScanTyped(&msg) {
		t.Fatal("expected a message")
	}
	init, isInit := msg.(InitMessage)
	if !isInit {
		t.Fatalf("expected InitMessage, got %T", msg)
	}
	if init.Body.NodeID != "n1" || init.Body.MsgID != 1 {
		t.Errorf("got %+v", init)
	}

	if !ScanTyped(&msg) {
		t.Fatal("expected a second message")
	}
	echo, isEcho := msg.(EchoMessage)
	if !isEcho {
		t.Fatalf("expected EchoMessage, got %T", msg)
	}
	if echo.Body.Echo != "hi" {
		t.Errorf("echo: got %q", echo.Body.Echo)
	}

	if ScanTyped(&msg) {
		t.Error("expected false at EOF")
	}
}

func TestScanTypedSkipsUnknownType(t *testing.T) {
	feedStdin(t, `{"src":"c0","dest":"n1","body":{"type":"mystery","msg_id":1}}`+"\n"+
		`{"src":"c0","dest":"n1","body":{"type":"echo","msg_id":2,"echo":"hi"}}`+"\n")

	var msg Incoming
	if !ScanTyped(&msg) {
		t.Fatal("expected to skip unknown type and return echo")
	}
	if _, isEcho := msg.(EchoMessage); !isEcho {
		t.Fatalf("expected EchoMessage, got %T", msg)
	}
}

func TestScanTypedSkipsBadJSON(t *testing.T) {
	feedStdin(t, "not json\n"+`{"src":"c0","dest":"n1","body":{"type":"init","msg_id":2,"node_id":"n1"}}`+"\n")

	var msg Incoming
	if !ScanTyped(&msg) {
		t.Fatal("expected true after skipping bad line")
	}
	if _, isInit := msg.(InitMessage); !isInit {
		t.Fatalf("expected InitMessage, got %T", msg)
	}
}
