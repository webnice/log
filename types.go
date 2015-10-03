package log

import (
	"bufio"
	"time"
)

const (
	level_FATAL    logLevel = iota // 0 Fatal: system is unusable
	level_ALERT                    // 1 Alert: action must be taken immediately
	level_CRITICAL                 // 2 Critical: critical conditions
	level_ERROR                    // 3 Error: error conditions
	level_WARNING                  // 4 Warning: warning conditions
	level_NOTICE                   // 5 Notice: normal but significant condition
	level_INFO                     // 6 Informational: informational messages
	level_DEBUG                    // 7 Debug: debug-level messages
)

const (
	defaultFormat string = `"%{color}[%{module:-10s}] %{time:2006-01-02T15:04:05.000Z07:00} (%{level:7s}): %{message} (%{package}) (%{function}:%{line}) (%{shortfile}:%{line}) (%{file})"`
)

var (
	self *configuration // Singleton
)

type logLevel int8

type configuration struct {
	FlushImmediately bool          // Flush log buffer after call
	Writer           *bufio.Writer // Log writer
}

type log struct {
	Now    time.Time     // Call time
	Trace  *trace        // Caller information
	Writer *bufio.Writer // Log writer
}
