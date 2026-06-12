package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

var stdin = bufio.NewScanner(os.Stdin)

// ScanTyped reads the next non-empty line, looks at its body "type" field, and
// stores the matching concrete message type in *target. Type-switch on it to
// handle it. Returns false on EOF.
//
//	var msg Incoming
//	for ScanTyped(&msg) {
//	    switch m := msg.(type) {
//	    case InitMessage:
//	        // m.Body.NodeID ...
//	    case EchoMessage:
//	        // m.Body.Echo ...
//	    }
//	}
func ScanTyped(target *Incoming) bool {
	for stdin.Scan() {
		line := stdin.Bytes()
		if len(line) == 0 {
			continue
		}
		var probe struct {
			Body struct {
				Type IncomingMessageType `json:"type"`
			} `json:"body"`
		}
		if err := json.Unmarshal(line, &probe); err != nil {
			Log.Error("parse JSON: %v", err)
			continue
		}
		msg, err := decodeByType(probe.Body.Type, line)
		if err != nil {
			Log.Error("%v", err)
			continue
		}
		*target = msg
		return true
	}
	if err := stdin.Err(); err != nil {
		Log.Error("scanner: %v", err)
	}
	return false
}

// decodeByType unmarshals line into the concrete message type named by t.
func decodeByType(t IncomingMessageType, line []byte) (Incoming, error) {
	switch t {
	case MsgTypeInit:
		var m InitMessage
		err := json.Unmarshal(line, &m)
		return m, err
	case MsgTypeEcho:
		var m EchoMessage
		err := json.Unmarshal(line, &m)
		return m, err
	default:
		return nil, fmt.Errorf("unknown message type: %q", t)
	}
}
