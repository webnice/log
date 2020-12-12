package message

import (
	l "github.com/webnice/lv2/level"
)

// Interface is an interface of package
type Interface interface {
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

	// Keys add to message
	Keys(...map[string]interface{}) Interface

	// CallStackCorrect Correction detect original call function
	CallStackCorrect(int) Interface

	// Message send with level and format
	Message(l.Level, string, ...interface{})

	// Done Flush all buffered messages and exit
	Done()
}

// impl is an implementation of package
type impl struct {
	callStack int                    // Глубина стэка
	keys      map[string]interface{} // Ключи сообщения
}
