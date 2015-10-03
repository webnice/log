package log

import (
	"time"
)

type logLevel int8

const (
	levelFatal    logLevel = 0
	levelAlert    logLevel = 1
	levelCritical logLevel = 2
	levelError    logLevel = 3
	levelWarning  logLevel = 4
	levelNotice   logLevel = 5
	levelInfo     logLevel = 6
	levelDebug    logLevel = 7
)

type log struct {
	Now   time.Time // Call time
	Trace *trace    // Caller information
}
