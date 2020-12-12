package lv2

import (
	l "github.com/webnice/lv2/level"
	m "github.com/webnice/lv2/message"
)

// Fatal Level 0: system is unusable
// A "panic" condition - notify all tech staff on call? (earthquake? tornado?) - affects multiple apps/servers/sites...
func Fatal(args ...interface{}) { ess.NewMsg().Fatal(args...) }

// Fatalf Level 0: system is unusable
// A "panic" condition - notify all tech staff on call? (earthquake? tornado?) - affects multiple apps/servers/sites...
func Fatalf(pattern string, args ...interface{}) { ess.NewMsg().Fatalf(pattern, args...) }

// Alert Level 1: action must be taken immediately
// Should be corrected immediately - notify staff who can fix the problem - example is loss of backup ISP connection
func Alert(args ...interface{}) { ess.NewMsg().Alert(args...) }

// Alertf Level 1: action must be taken immediately
// Should be corrected immediately - notify staff who can fix the problem - example is loss of backup ISP connection
func Alertf(pattern string, args ...interface{}) { ess.NewMsg().Alertf(pattern, args...) }

// Critical Level 2: critical conditions
// Should be corrected immediately, but indicates failure in a primary system - fix CRITICAL problems before ALERT - example is loss of primary ISP connection
func Critical(args ...interface{}) { ess.NewMsg().Critical(args...) }

// Criticalf Level 2: critical conditions
// Should be corrected immediately, but indicates failure in a primary system - fix CRITICAL problems before ALERT - example is loss of primary ISP connection
func Criticalf(pattern string, args ...interface{}) { ess.NewMsg().Criticalf(pattern, args...) }

// Error Level 3: error conditions
// Non-urgent failures - these should be relayed to developers or admins; each item must be resolved within a given time
func Error(args ...interface{}) { ess.NewMsg().Error(args...) }

// Errorf Level 3: error conditions
// Non-urgent failures - these should be relayed to developers or admins; each item must be resolved within a given time
func Errorf(pattern string, args ...interface{}) { ess.NewMsg().Errorf(pattern, args...) }

// Warning Level 4: warning conditions
// Warning messages - not an error, but indication that an error will occur if action is not taken, e.g. file system 85% full - each item must be resolved within a given time
func Warning(args ...interface{}) { ess.NewMsg().Warning(args...) }

// Warningf Level 4: warning conditions
// Warning messages - not an error, but indication that an error will occur if action is not taken, e.g. file system 85% full - each item must be resolved within a given time
func Warningf(pattern string, args ...interface{}) { ess.NewMsg().Warningf(pattern, args...) }

// Notice Level 5: normal but significant condition
// Events that are unusual but not error conditions - might be summarized in an email to developers or admins to spot potential problems - no immediate action required
func Notice(args ...interface{}) { ess.NewMsg().Notice(args...) }

// Noticef Level 5: normal but significant condition
// Events that are unusual but not error conditions - might be summarized in an email to developers or admins to spot potential problems - no immediate action required
func Noticef(pattern string, args ...interface{}) { ess.NewMsg().Noticef(pattern, args...) }

// Info Level 6: informational messages
// Normal operational messages - may be harvested for reporting, measuring throughput, etc - no action required
func Info(args ...interface{}) { ess.NewMsg().Info(args...) }

// Infof Level 6: informational messages
// Normal operational messages - may be harvested for reporting, measuring throughput, etc - no action required
func Infof(pattern string, args ...interface{}) { ess.NewMsg().Infof(pattern, args...) }

// Debug Level 7: debug-level messages
// Info useful to developers for debugging the app, not useful during operations
func Debug(args ...interface{}) { ess.NewMsg().Debug(args...) }

// DebDebugfug Level 7: debug-level messages
// Info useful to developers for debugging the app, not useful during operations
func Debugf(pattern string, args ...interface{}) { ess.NewMsg().Debugf(pattern, args...) }

// Keys Передача в сообщения лога дополнительных информационных полей в виде ключ/значение
func Keys(keys ...Key) Log {
	var tmp []map[string]interface{}
	var i int
	for i = range keys {
		tmp = append(tmp, map[string]interface{}(keys[i]))
	}
	return ess.NewMsg().(m.Interface).Keys(tmp...)
}

// Message send with level and format
func Message(lv l.Level, pat string, args ...interface{}) { ess.NewMsg().Message(lv, pat, args...) }

// Done Flush all buffered messages and exit
func Done() { ess.NewMsg().Done() }

// Gist get interface
func Gist() Essence { return ess }
