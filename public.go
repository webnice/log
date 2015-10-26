package log

// Configure log
func Configure(cnf *Configuration) error {
	return singleton[default_LOG].Configure(cnf)
}

// Level 0
// Fatal: system is unusable
// A "panic" condition - notify all tech staff on call? (earthquake? tornado?) - affects multiple apps/servers/sites...
func Fatal(args ...interface{}) {
	singleton[default_LOG].Fatal(args...)
}

// Level 1
// Alert: action must be taken immediately
// Should be corrected immediately - notify staff who can fix the problem - example is loss of backup ISP connection
func Alert(args ...interface{}) {
	singleton[default_LOG].Alert(args...)
}

// Level 2
// Critical: critical conditions
// Should be corrected immediately, but indicates failure in a primary system - fix CRITICAL problems before ALERT - example is loss of primary ISP connection
func Critical(args ...interface{}) {
	singleton[default_LOG].Critical(args...)
}

// Level 3
// Error: error conditions
// Non-urgent failures - these should be relayed to developers or admins; each item must be resolved within a given time
func Error(args ...interface{}) {
	singleton[default_LOG].Error(args...)
}

// Level 4
// Warning: warning conditions
// Warning messages - not an error, but indication that an error will occur if action is not taken, e.g. file system 85% full - each item must be resolved within a given time
func Warning(args ...interface{}) {
	singleton[default_LOG].Warning(args...)
}

// Level 5
// Notice: normal but significant condition
// Events that are unusual but not error conditions - might be summarized in an email to developers or admins to spot potential problems - no immediate action required
func Notice(args ...interface{}) {
	singleton[default_LOG].Notice(args...)
}

// Level 6
// Informational: informational messages
// Normal operational messages - may be harvested for reporting, measuring throughput, etc - no action required
func Info(args ...interface{}) {
	singleton[default_LOG].Info(args...)
}

// Level 7
// Debug: debug-level messages
// Info useful to developers for debugging the app, not useful during operations
func Debug(args ...interface{}) {
	singleton[default_LOG].Debug(args...)
}

// Flush log buffer immediately
//func Flush() {
//	singleton[default_LOG].Flush()
//}

// Close logging and reinitialisation defailt log
func Close() (err error) {
	singleton[default_LOG].Close()
	singleton[default_LOG] = NewLog()
	return
}

// Get default log object
func GetDefaultLog() *Log {
	return singleton[default_LOG]
}

// Set application name
// If name is empty, get name from os.Args[0] (string)
func SetApplicationName(name string) {
	singleton[default_LOG].SetApplicationName(name)
}

// Set module name
// If module name is empty or not set, name equals package name
func SetModuleName(name string) {
	singleton[default_LOG].SetModuleName(name)
}

// Remover module name
func DelModuleName() {
	singleton[default_LOG].DelModuleName()
}

// Configuring the interception of communications of a standard log
// flg=true  - intercept is enabled
// flg=false - intercept is desabled
func InterceptStandardLog(flg bool) {
	singleton[default_LOG].InterceptStandardLog(flg)
}

// Configuring the interception of STDOUT
// flg=true  - intercept is enabled
// flg=false - intercept is desabled
//func InterceptSTDOUT(flg bool) {
//	singleton[default_LOG].InterceptSTDOUT(flg)
//}

// Configuring the interception of STDERR
// flg=true  - intercept is enabled
// flg=false - intercept is desabled
//func InterceptSTDERR(flg bool) {
//	singleton[default_LOG].InterceptSTDERR(flg)
//}
