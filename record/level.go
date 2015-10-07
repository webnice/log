package record

import (

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

type (
	logLevel  int8
)
