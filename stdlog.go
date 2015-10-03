package log

import (
	stdLog "log"
	"os"
)

// Put io writer to log
func stdLogConnect() {
	var w *log = newLog()
	w.Record = nil
	stdLog.SetOutput(w)
}

// Reset to defailt
func stdLogClose() {
	stdLog.SetOutput(os.Stderr)
}
