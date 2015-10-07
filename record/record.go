package record

import (
	"syscall"
	"time"

	"github.com/webdeskltd/debug"
	"github.com/webdeskltd/log/uuid"
)

// Игноритуются слеюующие поля структуры
// - Не имеющие тега fmt:""
// - Описанные как fmt:""
// - Описанные как fmt:"-"
// Если формат поля структуры не описан, то подставляется формат %v
type Record struct {
	Id            uuid.UUID `fmt:"id"`                    // %{id}                      - ([16]byte ) Time GUID(UUID) for log message
	Pid           int       `fmt:"pid:d"`                 // %{pid}                     - (int      ) Process id
	AppName       string    `fmt:"application:s"`         // %{application}             - (string   ) Application name basename of os.Args[0]
	HostName      string    `fmt:"hostname:s"`            // %{hostname}                - (string   ) Server host name
	TodayAndNow   time.Time `fmt:"time"`                  // %{time}                    - (time.Time) Time when log occurred
	Level         int8      `fmt:"level:d"`               // %{level}                   - (int8     ) Log level
	Message       string    `fmt:"message:s"`             // %{message}                 - (string   ) Message
	FileNameLong  string    `fmt:"longfile:s"`            // %{longfile}                - (string   ) Full file name and line number: /a/b/c/d.go
	FileNameShort string    `fmt:"shortfile:s"`           // %{shortfile}               - (string   ) Final file name element and line number: d.go
	FileLine      int       `fmt:"line:d"`                // %{line}                    - (int      ) Line number in file
	Package       string    `fmt:"package:s"`             // %{package}                 - (string   ) Full package path, eg. github.com/webdeskltd/log
	Module        string    `fmt:"module:s,shortpkg:s"`   // %{module} or %{shortpkg}   - (string   ) Module name base package path, eg. log
	Function      string    `fmt:"function:s,facility:s"` // %{function} or %{facility} - (string   ) Full function name, eg. PutUint32
	CallStack     string    `fmt:"callstack:s"`           // %{callstack}               - (string   ) Full call stack
	color         bool      `fmt:"color"`                 // %{color}                   - (bool     ) ANSI color for messages in general, based on log level
	colorBeg      bool      `fmt:"colorbeg"`              // %{colorbeg}                - (bool     ) Mark the beginning colored text in message, based on log level
	colorEnd      bool      `fmt:"colorend"`              // %{colorend}                - (bool     ) Mark the ending colored text in message, based on log level
}

func init() {
	debug.Nop()
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

	//debug.Dumper(this)
	return this
}
