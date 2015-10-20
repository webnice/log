package log

import (
	"errors"
	"fmt"
	"strings"
	//"os"

	"github.com/webdeskltd/debug"
	"github.com/webdeskltd/log/backends"
	"github.com/webdeskltd/log/record"
	//"github.com/webdeskltd/log/gelf"

	//"github.com/webdeskltd/log/logging"
)

func init() {
	debug.Nop()
}

// Apply new configuration
func (self *Log) Configure(cnf *Configuration) (err error) {
	var i int
	var ok bool
	var mode ModeName

	if cnf == nil {
		err = ERROR_CONFIGURATION_IS_NULL
		return
	}
	self.cnf = cnf

	// Если формат по умолчанию для всех режимов не установлен, используем стандартный формат по умолчанию
	if self.cnf.Format == "" {
		self.cnf.Format = default_FORMAT
	}
	// Проверка формата по умолчанию
	_, err = record.CheckFormat(self.cnf.Format)
	if err != nil {
		return
	}

	for i = range self.cnf.Mode {
		mode = ModeName(strings.ToLower(string(self.cnf.Mode[i])))

		// Если для режима логирования не о определён формат, то присваиваем формат по умолчанию
		if _, ok = self.cnf.Formats[mode]; ok == false {
			self.cnf.Formats[mode] = self.cnf.Format
		}
		// Проверка формата
		_, err = record.CheckFormat(self.cnf.Formats[mode])
		if err != nil {
			return
		}

		switch mode {
		case mode_CONSOLE:
			print(mode_CONSOLE + "\n")
			self.backend = backends.NewBackends()
			self.backend.AddBackend(backends.NewBackendSTD(nil).SetModeNormal())
			//			var lh *logging.LogBackend
			//			lh = logging.NewLogBackend(os.Stderr, ``, 0)
			//			lh.Color = true
			//			self.bStderr = logging.AddModuleLevel(lh)
			//			self.bStderr.SetLevel(level, ``)
			//			self.backends = append(self.backends, self.bStderr)
		case mode_SYSLOG:
			print(mode_SYSLOG + "\n")
			//			var lh *logging.SyslogBackend
			//			lh, err = logging.NewSyslogBackend(self.HostName)
			//			if err != nil {
			//				Fatal("Can't initiate syslog backend %v", err)
			//			}
			//			self.bSyslog = logging.AddModuleLevel(lh)
			//			self.bSyslog.SetLevel(level, ``)
			//			self.backends = append(self.backends, self.bSyslog)
		case mode_FILE:
			print(mode_FILE + "\n")
			//			var lh *fileBackend
			//			if self.cnf.File == "" {
			//				Warning("Not specified log file name but mode 'file' is setted")
			//				continue
			//			}
			//			self.fH, err = os.OpenFile(self.cnf.File, syscall.O_APPEND|syscall.O_CREAT|syscall.O_WRONLY, 0644)
			//			if err != nil {
			//				Fatal("Can't initiate filelog backend %v", err)
			//			}
			//			lh = newFileBackend(self.fH)
			//			self.bFile = logging.AddModuleLevel(lh)
			//			self.bFile.SetLevel(level, ``)
			//			self.backends = append(self.backends, self.bFile)
		case mode_GRAYLOG2:
			print(mode_GRAYLOG2 + "\n")
			//			var hgpc gelf.GelfProtocolClient
			//			var hG *gelf.GelfClient
			//			var lh *gelfBackend
			//			if strings.EqualFold(strings.ToLower(self.cnf.Graylog.Protocol), strings.ToLower(gelf.UDP_NETWORK)) == true {
			//				hgpc = gelf.MustUdpClient(cnf.Graylog.Host, cnf.Graylog.Port, cnf.Graylog.ChunkSize)
			//			}
			//			if strings.EqualFold(strings.ToLower(self.cnf.Graylog.Protocol), strings.ToLower(gelf.TCP_NETWORK)) == true {
			//				hgpc = gelf.MustTcpClient(cnf.Graylog.Host, cnf.Graylog.Port)
			//			}
			//			hG = gelf.NewGelfClient(hgpc, self.cnf.Graylog.Compression)
			//			lh = newGelfBackend(hG, self.HostName, self.HostName)
			//			self.bGraylog = logging.AddModuleLevel(lh)
			//			self.bGraylog.SetLevel(level, ``)
			//			self.backends = append(self.backends, self.bGraylog)
		case mode_MEMPIPE:
			print(mode_MEMPIPE + "\n")
		case mode_TELEGRAM:
			print(mode_TELEGRAM + "\n")
		default:
			err = errors.New(fmt.Sprintf("%s %v", ERROR_UNKNOWN_MODE.Error(), mode))
			return
		}
	}
	//logging.SetBackend(self.backends...)
	//logFormatter = logging.MustStringFormatter(self.cnf.Format)
	//logging.SetFormatter(logFormatter)
	//debug.Dumper(self.cnf)
	return
}
