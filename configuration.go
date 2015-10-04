package log

import (
	"os"
	"strings"
	"syscall"

	"github.com/webdeskltd/debug"
	"github.com/webdeskltd/log/gelf"
	"github.com/webdeskltd/log/record"

	"github.com/webdeskltd/log/logging"
)

func init() {
	debug.Nop()
}

// Apply new configuration
func (self *configuration) Configure(cnf Configuration) (err error) {
	var i int
	var ok bool
	var level logging.Level
	var logFormatter *logging.StringFormatterStruct

	self.cnf = new(Configuration)
	*self.cnf = cnf

	// Если формат по умолчанию для всех режимов не установлен, используем стандартный формат по умолчанию
	if self.cnf.Format == "" {
		self.cnf.Format = defaultFormat
	}
	// Проверка формата
	err = record.CheckFormat(self.cnf.Format)
	if err != nil {
		return
	}

	for i = range self.cnf.Mode {
		level, err = logging.LogLevel(string(self.cnf.Levels[ConfigurationModeName(self.cnf.Mode[i])]))
		if err != nil {
			Error("Error log level. Mode: '%v', Level: '%v'", self.cnf.Mode[i], string(self.cnf.Levels[ConfigurationModeName(self.cnf.Mode[i])]))
			level = logging.NOTICE
			err = nil
		}
		// Если для режима логирования не о определён формат, то присваиваем формат по умолчанию
		if _, ok = self.cnf.Formats[self.cnf.Mode[i]]; ok == false {
			self.cnf.Formats[self.cnf.Mode[i]] = self.cnf.Format
		}
		// Проверка формата
		err = record.CheckFormat(self.cnf.Formats[self.cnf.Mode[i]])
		if err != nil {
			return
		}

		switch self.cnf.Mode[i] {
		case mode_CONSOLE:
			var lh *logging.LogBackend
			lh = logging.NewLogBackend(os.Stderr, ``, 0)
			lh.Color = true
			self.bStderr = logging.AddModuleLevel(lh)
			self.bStderr.SetLevel(level, ``)
			self.backends = append(self.backends, self.bStderr)
		case mode_SYSLOG:
			var lh *logging.SyslogBackend
			lh, err = logging.NewSyslogBackend(self.HostName)
			if err != nil {
				Fatal("Can't initiate syslog backend %v", err)
			}
			self.bSyslog = logging.AddModuleLevel(lh)
			self.bSyslog.SetLevel(level, ``)
			self.backends = append(self.backends, self.bSyslog)
		case mode_FILE:
			var lh *fileBackend
			if self.cnf.File == "" {
				Warning("Not specified log file name but mode 'file' is setted")
				continue
			}
			self.fH, err = os.OpenFile(self.cnf.File, syscall.O_APPEND|syscall.O_CREAT|syscall.O_WRONLY, 0644)
			if err != nil {
				Fatal("Can't initiate filelog backend %v", err)
			}
			lh = newFileBackend(self.fH)
			self.bFile = logging.AddModuleLevel(lh)
			self.bFile.SetLevel(level, ``)
			self.backends = append(self.backends, self.bFile)
		case mode_GRAYLOG:
			var hgpc gelf.GelfProtocolClient
			var hG *gelf.GelfClient
			var lh *gelfBackend
			if strings.EqualFold(strings.ToLower(self.cnf.Graylog.Protocol), strings.ToLower(gelf.UDP_NETWORK)) == true {
				hgpc = gelf.MustUdpClient(cnf.Graylog.Host, cnf.Graylog.Port, cnf.Graylog.ChunkSize)
			}
			if strings.EqualFold(strings.ToLower(self.cnf.Graylog.Protocol), strings.ToLower(gelf.TCP_NETWORK)) == true {
				hgpc = gelf.MustTcpClient(cnf.Graylog.Host, cnf.Graylog.Port)
			}
			hG = gelf.NewGelfClient(hgpc, self.cnf.Graylog.Compression)
			lh = newGelfBackend(hG, self.HostName, self.HostName)
			self.bGraylog = logging.AddModuleLevel(lh)
			self.bGraylog.SetLevel(level, ``)
			self.backends = append(self.backends, self.bGraylog)
		default:
			Fatal("Unknown logging mode %v", self.cnf.Mode[i])
		}
	}
	logging.SetBackend(self.backends...)
	logFormatter = logging.MustStringFormatter(self.cnf.Format)
	logging.SetFormatter(logFormatter)
	//debug.Dumper(self.cnf)
	return
}

func (self *configuration) SetApplicationName(name string) {
	var tmp []string
	self.AppName = name
	if self.AppName == "" {
		tmp = strings.Split(os.Args[0], string(os.PathSeparator))
		if len(tmp) > 0 {
			self.AppName = tmp[len(tmp)-1]
		}
	}
	return
}

// Set module name
func (self *configuration) SetModuleName(name string) {
	var r *record.Record
	if name != "" {
		r = newTrace().Trace(traceStepBack + 1).GetRecord()
		self.moduleNames[r.Package] = name
	}
	return
}

// Close logging
func (self *configuration) Close() (err error) {
	logging.Reset()
	stdLogClose()
	if self.fH != nil {
		err = self.fH.Sync()
		if err != nil {
			return
		}
		err = self.fH.Close()
	}
	return
}
