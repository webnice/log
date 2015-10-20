package log

import (
	//"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/webdeskltd/debug"
	"github.com/webdeskltd/log/backends"
	"github.com/webdeskltd/log/gelf"
	l "github.com/webdeskltd/log/level"
	m "github.com/webdeskltd/log/message"
	"github.com/webdeskltd/log/record"
	t "github.com/webdeskltd/log/trace"

	//"github.com/webdeskltd/log/gelf"

	//"github.com/webdeskltd/log/logging"
)

// Initialize default log settings
func init() {
	singleton = make(map[string]*Log)

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
		Mode:                   []ModeName{mode_CONSOLE},
		Levels:                 make(map[ModeName]LevelName),
		Formats:                make(map[ModeName]string),
		Format:                 default_FORMAT,
		Graylog2: ConfigurationGraylog2{
			Compression: gelf.COMPRESSION_NONE,
			Source:      self.HostName,
			Protocol:    gelf.UDP_NETWORK,
			BufferSize:  1000,
		},
		Telegram: ConfigurationTelegram{},
	}
	cnf.Levels[mode_CONSOLE] = LevelName(l.Map[default_LEVEL])
	return
}

// Initialize default configuration
func (self *Log) Initialize() *Log {
	var err error
	var cnf *Configuration

	// Default level writer
	self.defaultLevelLogWriter = m.NewWriter(default_LEVEL).Resolver(self.ResolveModuleName)

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
		fmt.Fprintf(os.Stderr, "Error Configure(): %v\n", err)
	} else {
		self.ready = true
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
	var r *record.Record
	if name != "" {
		r = t.NewTrace().Trace(t.STEP_BACK + 1).GetRecord()
		self.moduleNames[r.Package] = name
	}
	return self
}

// Remove module name
func (self *Log) DelModuleName() *Log {
	var r *record.Record
	r = t.NewTrace().Trace(t.STEP_BACK + 1).GetRecord()
	delete(self.moduleNames, r.Package)
	return self
}

// Resolve application name
func (self *Log) ResolveModuleName(r *record.Record) (name string) {
	var ok bool
	_, ok = self.moduleNames[r.Package]
	if ok {
		name = self.moduleNames[r.Package]
	} else {
		name = r.Package
	}
	return
}

// Configuring the interception of communications of a standard log
// flg=true  - intercept is enabled
// flg=false - intercept is desabled
func (self *Log) InterceptStandardLog(flg bool) *Log {
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
func (self *Log) InterceptSTDOUT(flg bool) *Log {
	if flg {
		if self.rescueSTDOUT == nil {
			self.rescueSTDOUT = os.Stdout
		}
	} else {
		if self.rescueSTDOUT != nil {
			os.Stdout = self.rescueSTDOUT
		}
	}
	return self
}

// Configuring the interception of STDERR
// flg=true  - intercept is enabled
// flg=false - intercept is desabled
func (self *Log) InterceptSTDERR(flg bool) *Log {
	if flg {
		if self.rescueSTDERR == nil {
			self.rescueSTDERR = os.Stderr
		}
	} else {
		if self.rescueSTDERR != nil {
			os.Stderr = self.rescueSTDERR
		}
	}
	return self
}

// Close logging
func (self *Log) Close() (err error) {
	// Flush buffer
	// err = self.Writer.Flush()

	// Reset standard logging to default settings
	self.InterceptStandardLog(false)

	// Create new backend object, old object automatic call Stop all backend and destroy
	self.backend = backends.NewBackends()

	// Reinitialisation
	//	selfMap["main"] = new(Log)
	//	selfMap["main"].Initialize()
	return
}
