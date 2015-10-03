package log

// Level 0
// Emergency: system is unusable
// A "panic" condition - notify all tech staff on call? (earthquake? tornado?) - affects multiple apps/servers/sites...
func Fatal(tpl string, args ...interface{}) {
	newLog().write(levelFatal, tpl, args...)
}

// Level 1
// Alert: action must be taken immediately
// Should be corrected immediately - notify staff who can fix the problem - example is loss of backup ISP connection
func Alert(tpl string, args ...interface{}) {
	newLog().write(levelAlert, tpl, args...)
}

// Level 2
// Critical: critical conditions
// Should be corrected immediately, but indicates failure in a primary system - fix CRITICAL problems before ALERT - example is loss of primary ISP connection
func Critical(tpl string, args ...interface{}) {
	newLog().write(levelCritical, tpl, args...)
}

// Level 3
// Error: error conditions
// Non-urgent failures - these should be relayed to developers or admins; each item must be resolved within a given time
func Error(tpl string, args ...interface{}) {
	newLog().write(levelError, tpl, args...)
}

// Level 4
// Warning: warning conditions
// Warning messages - not an error, but indication that an error will occur if action is not taken, e.g. file system 85% full - each item must be resolved within a given time
func Warning(tpl string, args ...interface{}) {
	newLog().write(levelWarning, tpl, args...)
}

// Level 5
// Notice: normal but significant condition
// Events that are unusual but not error conditions - might be summarized in an email to developers or admins to spot potential problems - no immediate action required
func Notice(tpl string, args ...interface{}) {
	newLog().write(levelNotice, tpl, args...)
}

// Level 6
// Informational: informational messages
// Normal operational messages - may be harvested for reporting, measuring throughput, etc - no action required
func Info(tpl string, args ...interface{}) {
	newLog().write(levelInfo, tpl, args...)
}

// Level 7
// Debug: debug-level messages
// Info useful to developers for debugging the app, not useful during operations
func Debug(tpl string, args ...interface{}) {
	newLog().write(levelDebug, tpl, args...)
}
