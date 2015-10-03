package record

import (
	"syscall"
	"time"

	"github.com/webdeskltd/debug"
	"github.com/webdeskltd/log/uuid"
)

type Record struct {
	Id            uuid.UUID // %{id}                      - ([16]byte ) Time GUID(UUID) for log message
	Pid           int       // %{pid}                     - (int    ) Process id
	AppName       string    // %{application}             - (string   ) Application name basename of os.Args[0]
	HostName      string    // %{hostname}                - (string   ) Server host name
	TodayAndNow   time.Time // %{time}                    - (time.Time) Time when log occurred
	Level         int8      // %{level}                   - (int8     ) Log level
	Message       string    // %{message}                 - (string   ) Message
	Color         bool      // %{color}                   - (bool     ) ANSI color based on log level
	FileNameLong  string    // %{longfile}                - (string   ) Full file name and line number: /a/b/c/d.go
	FileNameShort string    // %{shortfile}               - (string   ) Final file name element and line number: d.go
	FileLine      int       // %{line}                    - (int      ) Line number in file
	Package       string    // %{package}                 - (string   ) Full package path, eg. github.com/webdeskltd/log
	Module        string    // %{module} or %{shortpkg}   - (string   ) Module name base package path, eg. log
	Function      string    // %{function} or %{facility} - (string   ) Full function name, eg. PutUint32
	CallStack     string    // %{callstack}               - (string   ) Full call stack
}

func NewRecord() (self *Record) {
	self = new(Record)
	self.Id = uuid.TimeUUID()
	self.TodayAndNow = time.Now().In(time.Local)
	self.Pid = syscall.Getpid()
	return self
}

// Set log level
func (this *Record) SetLevel(level int8) *Record {
	this.Level = level
	return this
}

// Set message
func (this *Record) SetMessage(msg string) *Record {
	this.Message = msg
	return this
}

// Finishing object
func (this *Record) Complete() *Record {

	debug.Dumper(this)
	return this
}
