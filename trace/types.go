package trace // import "github.com/webnice/log/v2/trace"

import (
	"time"

	l "github.com/webnice/log/v2/level"
	u "github.com/webnice/log/v2/uuid"
)

const (
	stackBack        int    = 2
	packageSeparator string = `/`
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
	Id            u.UUID    `fmt:"id"                     json:"id"`                  // %{id}                      - ([16]byte ) Time GUID(UUID) for log message
	Pid           int       `fmt:"pid:d"                  json:"pid"`                 // %{pid}                     - (int      ) Process id
	AppName       string    `fmt:"application:s"          json:"appName"`             // %{application}             - (string   ) Application name basename of os.Args[0]
	HostName      string    `fmt:"hostname:s"             json:"hostName"`            // %{hostname}                - (string   ) Server host name
	TodayAndNow   time.Time `fmt:"time:t"                 json:"todayAndNow"`         // %{time}                    - (time.Time) Time when log occurred
	Level         l.Level   `fmt:"level:d"                json:"level"`               // %{level}                   - (int8     ) Log level
	Message       string    `fmt:"message:s"              json:"message,omitempty"`   // %{message}                 - (string   ) Message
	Color         bool      `fmt:"color"                  json:"-"`                   // %{color}                   - (bool     ) ANSI color for messages in general, based on log level
	ColorBeg      bool      `fmt:"colorbeg"               json:"-"`                   // %{colorbeg}                - (bool     ) Mark the beginning colored text in message, based on log level
	ColorEnd      bool      `fmt:"colorend"               json:"-"`                   // %{colorend}                - (bool     ) Mark the ending colored text in message, based on log level
	FileNameLong  string    `fmt:"longfile:s"             json:"filenameLong"`        // %{longfile}                - (string   ) Full file name and line number: /a/b/c/d.go
	FileNameShort string    `fmt:"shortfile:s"            json:"filenameShort"`       // %{shortfile}               - (string   ) Final file name element and line number: d.go
	FileLine      int       `fmt:"line:d"                 json:"fileLine"`            // %{line}                    - (int      ) Line number in file
	Package       string    `fmt:"package:s"              json:"package"`             // %{package}                 - (string   ) Full package path
	Module        string    `fmt:"module:s,shortpkg:s"    json:"module"`              // %{module} or %{shortpkg}   - (string   ) Module name base package path, eg. log
	Function      string    `fmt:"function:s,facility:s"  json:"function"`            // %{function} or %{facility} - (string   ) Full function name, eg. PutUint32
	CallStack     string    `fmt:"callstack:s"            json:"callStack,omitempty"` // %{callstack}               - (string   ) Full call stack
}
