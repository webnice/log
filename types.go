package log // import "github.com/webdeskltd/log"

import l "github.com/webdeskltd/log/level"
import w "github.com/webdeskltd/log/writer"
import s "github.com/webdeskltd/log/sender"
import f "github.com/webdeskltd/log/formater"

// const defaultTextFORMAT = `%{color}[%{module:-10s}] %{time:2006-01-02T15:04:05.000Z07:00t} (%{level:-8s}): %{message} (%{package}) (%{function}:%{line}) (%{shortfile}:%{line}) (%{longfile})`
const defaultTextFORMAT = `(%{colorbeg}%{level:1s}:%{level:1d}%{colorend}): %{message} {%{package}/%{shortfile}:%{line}, func: %{function}()}`

// Key type, used in Keys(Key{})
type Key map[string]interface{}

// Log interface
type Log interface {
	// Fatal Level 0: system is unusable
	Fatal(...interface{})
	// Fatalf Level 0: system is unusable
	Fatalf(string, ...interface{})
	// Alert Level 1: action must be taken immediately
	Alert(...interface{})
	// Alertf Level 1: action must be taken immediately
	Alertf(string, ...interface{})
	// Critical Level 2: critical conditions
	Critical(...interface{})
	// Criticalf Level 2: critical conditions
	Criticalf(string, ...interface{})
	// Error Level 3: error conditions
	Error(...interface{})
	// Errorf Level 3: error conditions
	Errorf(string, ...interface{})
	// Warning Level 4: warning conditions
	Warning(...interface{})
	// Warningf Level 4: warning conditions
	Warningf(string, ...interface{})
	// Notice Level 5: normal but significant condition
	Notice(...interface{})
	// Noticef Level 5: normal but significant condition
	Noticef(string, ...interface{})
	// Info Level 6: informational messages
	Info(...interface{})
	// Infof Level 6: informational messages
	Infof(string, ...interface{})
	// Debug Level 7: debug-level messages
	Debug(...interface{})
	// Debugf Level 7: debug-level messages
	Debugf(string, ...interface{})

	// Message send with level and format
	Message(l.Level, string, ...interface{})
	// Done Flush all buffered messages and exit
	Done()
}

// Объект логера
type impl struct {
	level    l.Level     // Log level severity or below
	writer   w.Interface // Writer interface
	sender   s.Interface // Sender interface
	formater f.Interface // Formater interface
	tplText  string      // Шаблон форматирования текста
}

// Essence interface
type Essence interface {
	// Return writer interface
	Writer() w.Interface
	// StandardLogSet Put io writer to log
	StandardLogSet() Essence
	// StandardLogUnset Reset to defailt
	StandardLogUnset() Essence
}
