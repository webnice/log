package log

import (
	"bufio"
	"fmt"
	"os"

	"github.com/webdeskltd/debug"
	"github.com/webdeskltd/log/backends"
	"github.com/webdeskltd/log/gelf"
)

// Initialize default log settings
func init() {
	var err error
	var cnf Configuration

	self = new(configuration)
	self.moduleNames = make(map[string]string)
	self.Writer = bufio.NewWriterSize(os.Stderr, self.BufferSize)
	self.HostName, err = os.Hostname()
	if err != nil {
		self.HostName = `undefined`
		fmt.Fprintf(os.Stderr, "Error get os.Hostname(): %v\n", err)
	}
	cnf = Configuration{
		BufferFlushImmediately: true,
		BufferSize:             0,
		Mode:                   []ConfigurationModeName{mode_CONSOLE},
		Levels:                 make(map[ConfigurationModeName]ConfigurationLevelName),
		Formats:                make(map[ConfigurationModeName]string),
		Format:                 defaultFormat,
		Graylog: ConfigurationGraylog{
			Compression: gelf.COMPRESSION_NONE,
			Source:      self.HostName,
			Protocol:    gelf.UDP_NETWORK,
			BufferSize:  1000,
		},
	}
	cnf.Levels[mode_CONSOLE] = ConfigurationLevelName(levelMap[defaultLevel])
	self.Configure(cnf)
	self.SetApplicationName(``)

	// Create default dackend - STDERR
	self.backend = backends.NewBackends()
	self.backend.AddBackend(backends.NewBackendSTD(nil).SetModeNormal())

	// Setup standard logging
	self.stdLogWriter = new(Writer)
	stdLogConnect(self.stdLogWriter)

	debug.Nop()
}

// Конструктор сообщений журнала
func newLogMessage() (this *logMessage) {
	this = new(logMessage)
	this.Record = newTrace().Trace(traceStepBack + 1).GetRecord()
	return
}

// Сюда попадают все сообщения от всех уровней логирования
func (this *logMessage) Write(level logLevel, tpl string, args ...interface{}) *logMessage {
	var tmp string

	tmp = fmt.Sprintf(tpl, args...)
	this.Record.SetLevel(int8(level)).SetMessage(tmp).Complete()
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

	self.backend.Push(this.Record)
	return this
}
