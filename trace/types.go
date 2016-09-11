package trace // import "github.com/webdeskltd/log/trace"

import "time"
import u "github.com/webdeskltd/log/uuid"
import l "github.com/webdeskltd/log/level"

const (
	_STACKBACK        int    = 2
	_PACKAGESEPARATOR string = `/`
)

// Interface is an interface of package
type Interface interface {
	Trace(int) Interface
	Info() *Info
}

// impl is an implementation of package
type impl struct {
	info Info
}

// Info Trace information
// Игноритуются слеюующие поля структуры
// - Не имеющие тега fmt:""
// - Описанные как fmt:""
// - Описанные как fmt:"-"
// Если формат поля структуры не описан, то подставляется формат %v
type Info struct {
	Id            u.UUID    `fmt:"id"`                    // %{id}                      - ([16]byte ) Time GUID(UUID) for log message
	Pid           int       `fmt:"pid:d"`                 // %{pid}                     - (int      ) Process id
	AppName       string    `fmt:"application:s"`         // %{application}             - (string   ) Application name basename of os.Args[0]
	HostName      string    `fmt:"hostname:s"`            // %{hostname}                - (string   ) Server host name
	TodayAndNow   time.Time `fmt:"time:t"`                // %{time}                    - (time.Time) Time when log occurred
	Level         l.Level   `fmt:"level:d"`               // %{level}                   - (int8     ) Log level
	Message       string    `fmt:"message:s"`             // %{message}                 - (string   ) Message
	Color         bool      `fmt:"color"`                 // %{color}                   - (bool     ) ANSI color for messages in general, based on log level
	ColorBeg      bool      `fmt:"colorbeg"`              // %{colorbeg}                - (bool     ) Mark the beginning colored text in message, based on log level
	ColorEnd      bool      `fmt:"colorend"`              // %{colorend}                - (bool     ) Mark the ending colored text in message, based on log level
	FileNameLong  string    `fmt:"longfile:s"`            // %{longfile}                - (string   ) Full file name and line number: /a/b/c/d.go
	FileNameShort string    `fmt:"shortfile:s"`           // %{shortfile}               - (string   ) Final file name element and line number: d.go
	FileLine      int       `fmt:"line:d"`                // %{line}                    - (int      ) Line number in file
	Package       string    `fmt:"package:s"`             // %{package}                 - (string   ) Full package path, eg. github.com/webdeskltd/log
	Module        string    `fmt:"module:s,shortpkg:s"`   // %{module} or %{shortpkg}   - (string   ) Module name base package path, eg. log
	Function      string    `fmt:"function:s,facility:s"` // %{function} or %{facility} - (string   ) Full function name, eg. PutUint32
	CallStack     string    `fmt:"callstack:s"`           // %{callstack}               - (string   ) Full call stack
}
