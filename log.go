package log

import (
	"bufio"
	"fmt"
	"os"

	"github.com/webdeskltd/log/gelf"
)

// Initialize default log settings
func init() {
	var err error
	var cnf Configuration

	self = new(configuration)
	self.moduleNames = make(map[string]string)
	self.Writer = bufio.NewWriterSize(os.Stderr, self.BufferSize)
	self.AppName = os.Args[0]
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
			Proto:       gelf.UDP_NETWORK,
		},
	}
	cnf.Levels[mode_CONSOLE] = ConfigurationLevelName(levelMap[defaultLevel])
	self.Configure(cnf)
	stdLogConnect()
}

func newLog() (this *log) {
	this = new(log)
	this.Record = newTrace().Trace(traceStepBack + 1).GetRecord()
	return
}

// All level writer
func (this *log) write(level logLevel, tpl string, args ...interface{}) *log {
	var tmp string

	tmp = fmt.Sprintf(tpl, args...)
	this.Record.SetLevel(int8(level)).SetMessage(tmp).Complete()
	this.WriteString(tmp)
	return this
}
