package core

import (
	"bufio"
	"encoding/json"
	"os"
	"reflect"
)

var stdin = bufio.NewScanner(os.Stdin)

// ScanLine reads the next non-empty line from stdin and unmarshals it into
// target (must be a non-nil pointer). The target is zeroed before each
// unmarshal so omitted fields don't carry over between messages.
// Returns false on EOF.
//
//	var msg Message
//	for ScanLine(&msg) {
//	    // handle msg
//	}
func ScanLine(target interface{}) bool {
	v := reflect.ValueOf(target)
	for stdin.Scan() {
		line := stdin.Bytes()
		if len(line) == 0 {
			continue
		}
		if v.Kind() == reflect.Ptr {
			v.Elem().Set(reflect.Zero(v.Elem().Type()))
		}
		if err := json.Unmarshal(line, target); err != nil {
			Log.Error("parse JSON: %v", err)
			continue
		}
		return true
	}
	if err := stdin.Err(); err != nil {
		Log.Error("scanner: %v", err)
	}
	return false
}
