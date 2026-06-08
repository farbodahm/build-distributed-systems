package main

import (
	"bufio"
	"encoding/json"
	"os"

	. "build-distributed-systems/internal/core"
)

type Message struct {
	Envelope
	Body map[string]interface{} `json:"body"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var msg Message
		if err := json.Unmarshal(line, &msg); err != nil {
			Log.Error("parse JSON: %v", err)
			continue
		}
		bodyType, ok := msg.Body["type"].(string)
		if !ok {
			bodyType = "unknown"
		}
		Log.Print("PARSED: %s|%s|%s", msg.Src, msg.Dest, bodyType)
		Log.Debug("src=%s dest=%s body=%v", msg.Src, msg.Dest, msg.Body)
	}
}
