package message

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	l "github.com/webnice/lv2/level"
	s "github.com/webnice/lv2/sender"
	t "github.com/webnice/lv2/trace"
	u "github.com/webnice/lv2/uuid"
)

// New Create new message object
func New() Interface {
	var msg = new(impl)
	msg.keys = make(map[string]interface{})
	msg.callStack = 4
	return msg
}

// Fatal Level 0: system is unusable
// A "panic" condition - notify all tech staff on call? (earthquake? tornado?) - affects multiple apps/servers/sites...
func (msg *impl) Fatal(args ...interface{}) {
	msg.Message(l.New().Fatal(), "", args...)
}

// Fatalf Level 0: system is unusable
// A "panic" condition - notify all tech staff on call? (earthquake? tornado?) - affects multiple apps/servers/sites...
func (msg *impl) Fatalf(pattern string, args ...interface{}) {
	msg.Message(l.New().Fatal(), pattern, args...)
}

// Alert Level 1: action must be taken immediately
// Should be corrected immediately - notify staff who can fix the problem - example is loss of backup ISP connection
func (msg *impl) Alert(args ...interface{}) {
	msg.Message(l.New().Alert(), "", args...)
}

// Alertf Level 1: action must be taken immediately
// Should be corrected immediately - notify staff who can fix the problem - example is loss of backup ISP connection
func (msg *impl) Alertf(pattern string, args ...interface{}) {
	msg.Message(l.New().Alert(), pattern, args...)
}

// Critical Level 2: critical conditions
// Should be corrected immediately, but indicates failure in a primary system - fix CRITICAL problems before ALERT - example is loss of primary ISP connection
func (msg *impl) Critical(args ...interface{}) {
	msg.Message(l.New().Critical(), "", args...)
}

// Criticalf Level 2: critical conditions
// Should be corrected immediately, but indicates failure in a primary system - fix CRITICAL problems before ALERT - example is loss of primary ISP connection
func (msg *impl) Criticalf(pattern string, args ...interface{}) {
	msg.Message(l.New().Critical(), pattern, args...)
}

// Error Level 3: error conditions
// Non-urgent failures - these should be relayed to developers or admins; each item must be resolved within a given time
func (msg *impl) Error(args ...interface{}) {
	msg.Message(l.New().Error(), "", args...)
}

// Errorf Level 3: error conditions
// Non-urgent failures - these should be relayed to developers or admins; each item must be resolved within a given time
func (msg *impl) Errorf(pattern string, args ...interface{}) {
	msg.Message(l.New().Error(), pattern, args...)
}

// Warning Level 4: warning conditions
// Warning messages - not an error, but indication that an error will occur if action is not taken, e.g. file system 85% full - each item must be resolved within a given time
func (msg *impl) Warning(args ...interface{}) {
	msg.Message(l.New().Warning(), "", args...)
}

// Warningf Level 4: warning conditions
// Warning messages - not an error, but indication that an error will occur if action is not taken, e.g. file system 85% full - each item must be resolved within a given time
func (msg *impl) Warningf(pattern string, args ...interface{}) {
	msg.Message(l.New().Warning(), pattern, args...)
}

// Notice Level 5: normal but significant condition
// Events that are unusual but not error conditions - might be summarized in an email to developers or admins to spot potential problems - no immediate action required
func (msg *impl) Notice(args ...interface{}) {
	msg.Message(l.New().Notice(), "", args...)
}

// Noticef Level 5: normal but significant condition
// Events that are unusual but not error conditions - might be summarized in an email to developers or admins to spot potential problems - no immediate action required
func (msg *impl) Noticef(pattern string, args ...interface{}) {
	msg.Message(l.New().Notice(), pattern, args...)
}

// Info Level 6: informational messages
// Normal operational messages - may be harvested for reporting, measuring throughput, etc - no action required
func (msg *impl) Info(args ...interface{}) {
	msg.Message(l.New().Informational(), "", args...)
}

// Infof Level 6: informational messages
// Normal operational messages - may be harvested for reporting, measuring throughput, etc - no action required
func (msg *impl) Infof(pattern string, args ...interface{}) {
	msg.Message(l.New().Informational(), pattern, args...)
}

// Debug Level 7: debug-level messages
// Info useful to developers for debugging the app, not useful during operations
func (msg *impl) Debug(args ...interface{}) {
	msg.Message(l.New().Debug(), "", args...)
}

// DebDebugfug Level 7: debug-level messages
// Info useful to developers for debugging the app, not useful during operations
func (msg *impl) Debugf(pattern string, args ...interface{}) {
	msg.Message(l.New().Debug(), pattern, args...)
}

// Keys add to message
func (msg *impl) Keys(keys ...map[string]interface{}) Interface {
	var (
		n int
		k string
	)

	for n = range keys {
		for k = range keys[n] {
			msg.keys[k] = keys[n][k]
		}
	}
	msg.CallStackCorrect(-1)

	return msg
}

// CallStackCorrect Correction detect original call function
func (msg *impl) CallStackCorrect(delta int) Interface {
	msg.callStack += delta
	return msg
}

// Message send with level and format
func (msg *impl) Message(level l.Level, pattern string, args ...interface{}) {
	var (
		rec s.Message
		err error
		tmp []string
	)

	rec.Level = level
	rec.Trace = t.New().Trace(msg.callStack).Info()
	rec.Pattern = pattern
	rec.Args = args
	rec.Keys = msg.keys

	rec.Trace.Id = u.TimeUUID()
	rec.Trace.Pid = syscall.Getpid()
	if tmp = strings.Split(os.Args[0], string(os.PathSeparator)); len(tmp) > 0 {
		rec.Trace.AppName = tmp[len(tmp)-1]
	}
	if rec.Trace.HostName, err = os.Hostname(); err != nil {
		rec.Trace.HostName = "Hostname not defined"
	}
	rec.Trace.TodayAndNow = time.Now().In(time.Local)
	rec.Trace.Level = level
	if pattern == "" {
		rec.Trace.Message = fmt.Sprint(args...)
	} else {
		rec.Trace.Message = fmt.Sprintf(pattern, args...)
	}

	s.Gist().Channel() <- rec
}

// Done Flush all buffered messages and exit
func (msg *impl) Done() { s.Gist().Flush() }
