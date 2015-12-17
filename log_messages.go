package log

import (
	"os"
	"runtime"
	"strings"

	b "github.com/webdeskltd/log/backends"
	l "github.com/webdeskltd/log/level"
	m "github.com/webdeskltd/log/message"
	r "github.com/webdeskltd/log/record"
	t "github.com/webdeskltd/log/trace"
	u "github.com/webdeskltd/log/uuid"
	w "github.com/webdeskltd/log/writer"

	//"github.com/webdeskltd/debug"
)

func (log *LogEssence) leveled(level l.Level, args ...interface{}) *LogEssence {
	log.backend.Push(
		m.NewMessage(
			t.NewTrace().
				Trace(t.STEP_BACK + 2).
				GetRecord().
				Resolver(log.ResolveNames),
		).Level(level).
			Write(args...),
	)
	return log
}

// Level 0
// Fatal: system is unusable
// A "panic" condition - notify all tech staff on call? (earthquake? tornado?) - affects multiple apps/servers/sites...
func (log *LogEssence) Fatal(args ...interface{}) {
	log.leveled(l.FATAL, args...).Close()
	exit_func(1)
	return
}

// Level 1
// Alert: action must be taken immediately
// Should be corrected immediately - notify staff who can fix the problem - example is loss of backup ISP connection
func (log *LogEssence) Alert(args ...interface{}) Log {
	return log.leveled(l.ALERT, args...).Interface
}

// Level 2
// Critical: critical conditions
// Should be corrected immediately, but indicates failure in a primary system - fix CRITICAL problems before ALERT - example is loss of primary ISP connection
func (log *LogEssence) Critical(args ...interface{}) Log {
	return log.leveled(l.CRITICAL, args...).Interface
}

// Level 3
// Error: error conditions
// Non-urgent failures - these should be relayed to developers or admins; each item must be resolved within a given time
func (log *LogEssence) Error(args ...interface{}) Log {
	return log.leveled(l.ERROR, args...).Interface
}

// Level 4
// Warning: warning conditions
// Warning messages - not an error, but indication that an error will occur if action is not taken, e.g. file system 85% full - each item must be resolved within a given time
func (log *LogEssence) Warning(args ...interface{}) Log {
	return log.leveled(l.WARNING, args...).Interface
}

// Level 5
// Notice: normal but significant condition
// Events that are unusual but not error conditions - might be summarized in an email to developers or admins to spot potential problems - no immediate action required
func (log *LogEssence) Notice(args ...interface{}) Log {
	return log.leveled(l.NOTICE, args...).Interface
}

// Level 6
// Informational: informational messages
// Normal operational messages - may be harvested for reporting, measuring throughput, etc - no action required
func (log *LogEssence) Info(args ...interface{}) Log {
	return log.leveled(l.INFO, args...).Interface
}

// Level 7
// Debug: debug-level messages
// Info useful to developers for debugging the app, not useful during operations
func (log *LogEssence) Debug(args ...interface{}) Log {
	return log.leveled(l.DEBUG, args...).Interface
}

// Message To send a message to the log with the level of logging
func (log *LogEssence) Message(level l.Level, args ...interface{}) Log {
	ll := log.leveled(level, args...)
	if level == l.FATAL {
		ll.Close()
		exit_func(1)
	}
	return ll.Interface
}

// Close logging
func (log *LogEssence) Close() (err error) {
	// Reset standard logging to default settings
	log.InterceptStandardLog(false)
	log.defaultLevelLogWriter = nil

	// Block programm while goroutine exit
	log.backend.Close()

	// Create new backend object, old object automatic call Stop all backend and destroy
	log.backend = b.NewBackends()

	// Reinitialisation
	var uuid, _ = u.ParseUUID(default_LOGUUID)
	singleton[default_LOGUUID] = newLogEssence(uuid)

	runtime.GC()
	runtime.Gosched()
	return
}

// Set application name
func (log *LogEssence) SetApplicationName(name string) Log {
	var tmp []string
	log.AppName = name
	if log.AppName == "" {
		tmp = strings.Split(os.Args[0], string(os.PathSeparator))
		if len(tmp) > 0 {
			log.AppName = tmp[len(tmp)-1]
		}
	}
	return log.Interface
}

// Set module name
func (log *LogEssence) SetModuleName(name string) Log {
	var rec *r.Record
	if name != "" {
		rec = t.NewTrace().Trace(t.STEP_BACK + 1).GetRecord()
		log.moduleNames[rec.Package] = name
	}
	return log.Interface
}

// Remove module name
func (self *LogEssence) DelModuleName() Log {
	var rec *r.Record
	rec = t.NewTrace().Trace(t.STEP_BACK + 1).GetRecord()
	delete(self.moduleNames, rec.Package)
	return self.Interface
}

// Configuring the interception of communications of a standard log
// flg=true  - intercept is enabled
// flg=false - intercept is desabled
func (self *LogEssence) InterceptStandardLog(flg bool) Log {
	self.interceptStandardLog = flg
	if flg {
		stdLogConnect(self.defaultLevelLogWriter)
	} else {
		stdLogClose()
	}
	return self.Interface
}

// GetWriter Returns the standard writer to logging
func (self *LogEssence) GetWriter() *w.Writer {
	return self.defaultLevelLogWriter
}
