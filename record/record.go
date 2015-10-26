package record

import (
	"fmt"
	"syscall"
	"time"

	l "github.com/webdeskltd/log/level"
	u "github.com/webdeskltd/log/uuid"

	"github.com/webdeskltd/debug"
)

// Игноритуются слеюующие поля структуры
// - Не имеющие тега fmt:""
// - Описанные как fmt:""
// - Описанные как fmt:"-"
// Если формат поля структуры не описан, то подставляется формат %v
type Record struct {
	Id            u.UUID    `fmt:"id"`                    // %{id}                      - ([16]byte ) Time GUID(UUID) for log message
	Pid           int       `fmt:"pid:d"`                 // %{pid}                     - (int      ) Process id
	AppName       string    `fmt:"application:s"`         // %{application}             - (string   ) Application name basename of os.Args[0]
	HostName      string    `fmt:"hostname:s"`            // %{hostname}                - (string   ) Server host name
	TodayAndNow   time.Time `fmt:"time"`                  // %{time}                    - (time.Time) Time when log occurred
	Level         l.Level   `fmt:"level:d"`               // %{level}                   - (int8     ) Log level
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

	msgs     []interface{}   `fmt:"-"` // ---------------------------- Original message and params
	resolver RecordResolveFn `fmt:"-"` // ---------------------------- Function name record resolution
}

// Функция разрешения имён записи
type RecordResolveFn func(r *Record)

func init() {
	debug.Nop()
}

func NewRecord() (this *Record) {
	this = new(Record)
	this.Id = u.TimeUUID()
	this.TodayAndNow = time.Now().In(time.Local)
	this.Pid = syscall.Getpid()
	return this
}

// Set log level
func (self *Record) SetLevel(level l.Level) *Record {
	self.Level = level
	return self
}

// Assigning function name resolution
func (self *Record) Resolver(f RecordResolveFn) *Record {
	self.resolver = f
	return self
}

// Set message
func (self *Record) SetMessage(args ...interface{}) *Record {
	self.msgs = args[:]
	return self
}

// Подготовка сообщения
// Выполняется подготовка сообщения перед форматированием и выводом
// На данном этапе из переданных ранее args формируется единое текстовое сообщение
func (self *Record) Prepare() *Record {
	var ok bool
	if len(self.msgs) > 0 {
		switch self.msgs[0].(type) {
		case string:
			ok = true
		}
	}
	if ok && len(self.msgs) > 1 {
		self.Message = fmt.Sprintf(self.msgs[0].(string), self.msgs[1:]...)
	} else {
		self.Message = fmt.Sprint(self.msgs[:]...)
	}
	return self
}

// Finishing object
func (self *Record) End() *Record {
	if self.resolver != nil {
		self.resolver(self)
	}

	//debug.Dumper(this)
	return self
}
