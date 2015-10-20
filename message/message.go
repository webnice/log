package message

import (
	"fmt"

	"github.com/webdeskltd/debug"
	l "github.com/webdeskltd/log/level"
	"github.com/webdeskltd/log/record"
	// t "github.com/webdeskltd/log/trace"
)

type Message struct {
	Record   *record.Record
	level    l.Level
	WriteLen int
	WriteErr error
}

func init() {
	debug.Nop()
}

// Конструктор сообщений журнала
func NewMessage(record *record.Record) (this *Message) {
	this = new(Message)
	this.Record = record
	// this.Record = t.NewTrace().Trace(t.STEP_BACK + 1).GetRecord()
	return
}

// Устанавливаем уровень логирования для сообщения
func (self *Message) Level(level l.Level) *Message {
	self.level = level
	self.Record.Level = level
	return self
}

// Сюда попадают все сообщения от всех уровней логирования
func (this *Message) Write(args ...interface{}) *Message {
	var tmp string

	tmp = fmt.Sprintf("%v", args...)
	this.Record.SetMessage(tmp).End()
	//	this.WriteString(tmp)

	//	debug.Dumper(this.Record)

	var test string = ` 1: %{id}
 2: %{pid:8d}                  - (int      ) Process id
 3: %{application}             - (string   ) Application name basename of os.Args[0]
 4: %{hostname}                - (string   ) Server host name
 5: %{time}                    - (time.Time) Time when log occurred
 6: %{level:-8d}               - (int8     ) Log level
 7: %{message}                 - (string   ) Message
 8: %{color}                   - %{begcolor}(bool     ) ANSI color based on log level%{endcolor}
 9: %{longfile}                - (string   ) Full file name and line number: /a/b/c/d.go
10: %{shortfile}               - (string   ) Final file name element and line number: d.go
11: %{line}                    - (int      ) Line number in file
12: %{package}                 - (string   ) Full package path, eg. github.com/webdeskltd/log
13: %{module} or %{shortpkg}   - (string   ) Module name base package path, eg. log
14: %{function} or %{facility} - (string   ) Full function name, eg. PutUint32
15: %{callstack}               - (string   ) Full call stack

"%{color}[%{module:-10s}] %{time:2006-01-02T15:04:05.000Z07:00} (%{level:7s}): %{message} (%{package}) (%{function}:%{line}) (%{shortfile}:%{line}) (%{longfile})"
%{level:.1s}`
	this.Record.Format(test)

	//self.backend.Push(this.Record)
	return this
}
