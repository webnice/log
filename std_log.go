package log

import (
	"io"
	stdLog "log"
	"os"
)

// Put io writer to log
func stdLogConnect(w io.Writer) {
	stdLog.SetPrefix(``)
	stdLog.SetFlags(0)
	stdLog.SetOutput(w)
}

// Reset to defailt
func stdLogClose() {
	stdLog.SetPrefix(``)
	stdLog.SetFlags(stdLog.LstdFlags)
	stdLog.SetOutput(os.Stderr)
}
