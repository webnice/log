package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	b "github.com/webdeskltd/log/backends"
	g "github.com/webdeskltd/log/gelf"
	l "github.com/webdeskltd/log/level"
	r "github.com/webdeskltd/log/record"
	t "github.com/webdeskltd/log/trace"
	w "github.com/webdeskltd/log/writer"

	"github.com/webdeskltd/debug"
)

// Initialize default log settings
func init() {
	// Карта всех копий logger
	singleton = make(map[string]*Log)

	// Устанавливаем в зависимые пакеты функции информирования об ошибках
	b.LogError = Error

	// Defailt backend format
	b.DefaultFormat = default_FORMAT

	// Default public log object
	singleton[default_LOG] = NewLog()

	// Intercept standard logging only first
	singleton[default_LOG].InterceptStandardLog(true)

	debug.Nop()
}

// New log object
func NewLog() (obj *Log) {
	obj = new(Log)
	obj.moduleNames = make(map[string]string)
	obj.Initialize()
	return
}

// Create dafault configuration
func (self *Log) defaultConfiguration() (cnf *Configuration) {
	cnf = &Configuration{
		BufferFlushImmediately: true,
		BufferSize:             0,
		Mode:                   make(map[b.BackendName][]l.LevelName),
		Levels:                 make(map[b.BackendName]l.LevelName),
		Formats:                make(map[b.BackendName]string),
		Format:                 default_FORMAT,
		Graylog2: ConfigurationGraylog2{
			Compression: g.COMPRESSION_NONE,
			Source:      self.HostName,
			Protocol:    g.UDP_NETWORK,
			BufferSize:  1000,
		},
		Telegram: ConfigurationTelegram{},
	}
	cnf.Mode[b.NAME_CONSOLE] = nil
	cnf.Levels[b.NAME_CONSOLE] = l.LevelName(l.Map[default_LEVEL])
	return
}

// Initialize default configuration
func (self *Log) Initialize() *Log {
	var err error
	var cnf *Configuration

	self.SetApplicationName(``)
	self.HostName, err = os.Hostname()
	if err != nil {
		self.HostName = `undefined`
		fmt.Fprintf(os.Stderr, "Error get os.Hostname(): %v\n", err)
	}

	// Create default configuration and apply
	cnf = self.defaultConfiguration()
	err = self.Configure(cnf)
	if err != nil {
		Error("Error Configure(): %v\n", err)
	} else {
		self.ready = true
	}

	// Default level writer
	self.defaultLevelLogWriter = w.NewWriter(default_LEVEL).Resolver(self.ResolveNames).AttachBackends(self.backend)
	if self.interceptStandardLog {
		stdLogConnect(self.defaultLevelLogWriter)
	}

	return self
}

// Set application name
func (self *Log) SetApplicationName(name string) *Log {
	var tmp []string
	self.AppName = name
	if self.AppName == "" {
		tmp = strings.Split(os.Args[0], string(os.PathSeparator))
		if len(tmp) > 0 {
			self.AppName = tmp[len(tmp)-1]
		}
	}
	return self
}

// Set module name
func (self *Log) SetModuleName(name string) *Log {
	var rec *r.Record
	if name != "" {
		rec = t.NewTrace().Trace(t.STEP_BACK + 1).GetRecord()
		self.moduleNames[rec.Package] = name
	}
	return self
}

// Remove module name
func (self *Log) DelModuleName() *Log {
	var rec *r.Record
	rec = t.NewTrace().Trace(t.STEP_BACK + 1).GetRecord()
	delete(self.moduleNames, rec.Package)
	return self
}

// Resolve resord
func (self *Log) ResolveNames(rec *r.Record) {
	rec.AppName = self.AppName
	rec.HostName = self.HostName
	if _, ok := self.moduleNames[rec.Package]; ok == true {
		rec.Package = self.moduleNames[rec.Package]
	}
	return
}

// Configuring the interception of communications of a standard log
// flg=true  - intercept is enabled
// flg=false - intercept is desabled
func (self *Log) InterceptStandardLog(flg bool) *Log {
	self.interceptStandardLog = flg
	if flg {
		stdLogConnect(self.defaultLevelLogWriter)
	} else {
		stdLogClose()
	}
	return self
}

// Configuring the interception of STDOUT
// flg=true  - intercept is enabled
// flg=false - intercept is desabled
//func (self *Log) InterceptSTDOUT(flg bool) *Log {
//	if flg {
//		if self.rescueSTDOUT == nil {
//			self.rescueSTDOUT = os.Stdout
//		}
//	} else {
//		if self.rescueSTDOUT != nil {
//			os.Stdout = self.rescueSTDOUT
//		}
//	}
//	return self
//}

// Configuring the interception of STDERR
// flg=true  - intercept is enabled
// flg=false - intercept is desabled
//func (self *Log) InterceptSTDERR(flg bool) *Log {
//	if flg {
//		if self.rescueSTDERR == nil {
//			self.rescueSTDERR = os.Stderr
//		}
//	} else {
//		if self.rescueSTDERR != nil {
//			os.Stderr = self.rescueSTDERR
//		}
//	}
//	return self
//}

// Close logging
func (self *Log) Close() (err error) {
	// Reset standard logging to default settings
	self.InterceptStandardLog(false)
	self.defaultLevelLogWriter = nil

	// Block programm while goroutine exit
	self.backend.Close()

	// Create new backend object, old object automatic call Stop all backend and destroy
	self.backend = b.NewBackends()

	// Reinitialisation
	singleton[default_LOG] = NewLog()
	self = singleton[default_LOG]
	runtime.GC()
	runtime.Gosched()
	return
}
