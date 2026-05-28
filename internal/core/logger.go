package core

import (
	"encoding/json"
	"log"
	"os"
)

// Logger routes Print to stdout (protocol output) and level methods to stderr
// for debugging. Safe for concurrent use.
type Logger struct {
	out *log.Logger
	err *log.Logger
}

var Log = &Logger{
	out: log.New(os.Stdout, "", 0),
	err: log.New(os.Stderr, "", 0),
}

// Print writes to stdout. Use for protocol output the grader checks.
func (lg *Logger) Print(format string, args ...interface{}) {
	lg.out.Printf(format, args...)
}

// PrintJSON marshals v and writes it to stdout as a single line. Use for
// protocol replies. On marshal error, logs to stderr and writes nothing.
func (lg *Logger) PrintJSON(v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		lg.Error("marshal: %v", err)
		return
	}
	lg.out.Printf("%s", b)
}

func (lg *Logger) Info(format string, args ...interface{}) {
	lg.err.Printf("[INFO]  "+format, args...)
}

func (lg *Logger) Debug(format string, args ...interface{}) {
	lg.err.Printf("[DEBUG] "+format, args...)
}

func (lg *Logger) Warn(format string, args ...interface{}) {
	lg.err.Printf("[WARN]  "+format, args...)
}

func (lg *Logger) Error(format string, args ...interface{}) {
	lg.err.Printf("[ERROR] "+format, args...)
}
