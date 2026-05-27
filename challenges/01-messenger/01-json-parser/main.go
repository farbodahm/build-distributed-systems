package main

import (
	"bufio"
	"encoding/json"
	"os"

	. "build-distributed-systems/internal/core"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		var msg Message
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			Log.Error("parsing JSON: %v", err)
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
