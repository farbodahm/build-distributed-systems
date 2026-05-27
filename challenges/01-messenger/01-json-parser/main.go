package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	. "build-distributed-systems/internal/core"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		var msg Message
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing JSON:", err)
			continue
		}
		bodyType, ok := msg.Body["type"].(string)
		if !ok {
			bodyType = "unknown"
		}
		fmt.Printf("PARSED: %s|%s|%s\n", msg.Src, msg.Dest, bodyType)
		fmt.Fprintf(os.Stderr, "DEBUG: src=%s dest=%s body=%v\n", msg.Src, msg.Dest, msg.Body)
	}
}
