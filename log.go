package log

import (
	"bufio"
	"fmt"
	"os"

	"github.com/webdeskltd/debug"
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
	stdLogConnect()
	
	debug.Nop()
}

// Конструктор сообщений журнала
func newLog() (this *log) {
	this = new(log)
	this.Record = newTrace().Trace(traceStepBack + 1).GetRecord()
	return
}

// Сюда подапаю все сообщения от всех уровней
func (this *log) write(level logLevel, tpl string, args ...interface{}) *log {
	var tmp string

	tmp = fmt.Sprintf(tpl, args...)
	this.Record.SetLevel(int8(level)).SetMessage(tmp).Complete()
	//this.WriteString(tmp)

	//debug.Dumper(this.Record)
	this.Record.Format()


	return this
}
