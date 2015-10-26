package log

import (
	"os"

	l "github.com/webdeskltd/log/level"
	m "github.com/webdeskltd/log/message"
	t "github.com/webdeskltd/log/trace"
)

// Level 0
// Fatal: system is unusable
// A "panic" condition - notify all tech staff on call? (earthquake? tornado?) - affects multiple apps/servers/sites...
func (self *Log) Fatal(args ...interface{}) *Log {
	self.backend.Push(
		m.NewMessage(
			t.NewTrace().
				Trace(t.STEP_BACK + 1).
				GetRecord().
				Resolver(self.ResolveNames),
		).Level(l.FATAL).
			Write(args...),
	)
	self.Close()
	os.Exit(1)
	return self
}

// Level 1
// Alert: action must be taken immediately
// Should be corrected immediately - notify staff who can fix the problem - example is loss of backup ISP connection
func (self *Log) Alert(args ...interface{}) *Log {
	self.backend.Push(
		m.NewMessage(
			t.NewTrace().
				Trace(t.STEP_BACK + 1).
				GetRecord().
				Resolver(self.ResolveNames),
		).Level(l.ALERT).
			Write(args...),
	)
	return self
}

// Level 2
// Critical: critical conditions
// Should be corrected immediately, but indicates failure in a primary system - fix CRITICAL problems before ALERT - example is loss of primary ISP connection
func (self *Log) Critical(args ...interface{}) *Log {
	self.backend.Push(
		m.NewMessage(
			t.NewTrace().
				Trace(t.STEP_BACK + 1).
				GetRecord().
				Resolver(self.ResolveNames),
		).Level(l.CRITICAL).
			Write(args...),
	)
	return self
}

// Level 3
// Error: error conditions
// Non-urgent failures - these should be relayed to developers or admins; each item must be resolved within a given time
func (self *Log) Error(args ...interface{}) *Log {
	self.backend.Push(
		m.NewMessage(
			t.NewTrace().
				Trace(t.STEP_BACK + 1).
				GetRecord().
				Resolver(self.ResolveNames),
		).Level(l.ERROR).
			Write(args...),
	)
	return self
}

// Level 4
// Warning: warning conditions
// Warning messages - not an error, but indication that an error will occur if action is not taken, e.g. file system 85% full - each item must be resolved within a given time
func (self *Log) Warning(args ...interface{}) *Log {
	self.backend.Push(
		m.NewMessage(
			t.NewTrace().
				Trace(t.STEP_BACK + 1).
				GetRecord().
				Resolver(self.ResolveNames),
		).Level(l.WARNING).
			Write(args...),
	)
	return self
}

// Level 5
// Notice: normal but significant condition
// Events that are unusual but not error conditions - might be summarized in an email to developers or admins to spot potential problems - no immediate action required
func (self *Log) Notice(args ...interface{}) *Log {
	self.backend.Push(
		m.NewMessage(
			t.NewTrace().
				Trace(t.STEP_BACK + 1).
				GetRecord().
				Resolver(self.ResolveNames),
		).Level(l.NOTICE).
			Write(args...),
	)
	return self
}

// Level 6
// Informational: informational messages
// Normal operational messages - may be harvested for reporting, measuring throughput, etc - no action required
func (self *Log) Info(args ...interface{}) *Log {
	self.backend.Push(
		m.NewMessage(
			t.NewTrace().
				Trace(t.STEP_BACK + 1).
				GetRecord().
				Resolver(self.ResolveNames),
		).Level(l.INFO).
			Write(args...),
	)
	return self
}

// Level 7
// Debug: debug-level messages
// Info useful to developers for debugging the app, not useful during operations
func (self *Log) Debug(args ...interface{}) *Log {
	self.backend.Push(
		m.NewMessage(
			t.NewTrace().
				Trace(t.STEP_BACK + 1).
				GetRecord().
				Resolver(self.ResolveNames),
		).Level(l.DEBUG).
			Write(args...),
	)
	return self
}

// Flush log buffer immediately
//func (self *Log) Flush() *Log {
//	return self
//}
