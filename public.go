package log

// Level 0
// Fatal: system is unusable
// A "panic" condition - notify all tech staff on call? (earthquake? tornado?) - affects multiple apps/servers/sites...
func Fatal(tpl string, args ...interface{}) {
	newLogMessage().Write(level_FATAL, tpl, args...)
}

// Level 1
// Alert: action must be taken immediately
// Should be corrected immediately - notify staff who can fix the problem - example is loss of backup ISP connection
func Alert(tpl string, args ...interface{}) {
	newLogMessage().Write(level_ALERT, tpl, args...)
}

// Level 2
// Critical: critical conditions
// Should be corrected immediately, but indicates failure in a primary system - fix CRITICAL problems before ALERT - example is loss of primary ISP connection
func Critical(tpl string, args ...interface{}) {
	newLogMessage().Write(level_CRITICAL, tpl, args...)
}

// Level 3
// Error: error conditions
// Non-urgent failures - these should be relayed to developers or admins; each item must be resolved within a given time
func Error(tpl string, args ...interface{}) {
	newLogMessage().Write(level_ERROR, tpl, args...)
}

// Level 4
// Warning: warning conditions
// Warning messages - not an error, but indication that an error will occur if action is not taken, e.g. file system 85% full - each item must be resolved within a given time
func Warning(tpl string, args ...interface{}) {
	newLogMessage().Write(level_WARNING, tpl, args...)
}

// Level 5
// Notice: normal but significant condition
// Events that are unusual but not error conditions - might be summarized in an email to developers or admins to spot potential problems - no immediate action required
func Notice(tpl string, args ...interface{}) {
	newLogMessage().Write(level_NOTICE, tpl, args...)
}

// Level 6
// Informational: informational messages
// Normal operational messages - may be harvested for reporting, measuring throughput, etc - no action required
func Info(tpl string, args ...interface{}) {
	newLogMessage().Write(level_INFO, tpl, args...)
}

// Level 7
// Debug: debug-level messages
// Info useful to developers for debugging the app, not useful during operations
func Debug(tpl string, args ...interface{}) {
	newLogMessage().Write(level_DEBUG, tpl, args...)
}

// Flush log buffer immediately
//func Flush() error {
//	return self.Writer.Flush()
//}

// Close logging
func Close() (err error) {
	err = self.Writer.Flush()
	if err != nil {
		return
	}
	err = self.Close()
	return
}

// Configure log
func Configure(cnf Configuration) error {
	return self.Configure(cnf)
}

// Set application name
// If name is empty, get name from os.Args[0] (string)
func SetApplicationName(name string) {
	self.SetApplicationName(name)

}

// Set module name
// If module name is empty or not set, name equals package name
func SetModuleName(name string) {
	self.SetModuleName(name)
}

// Remover module name
func DelModuleName() {
	self.DelModuleName()
}
