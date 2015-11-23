package log

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	b "github.com/webdeskltd/log/backends"
	g "github.com/webdeskltd/log/gelf"
	l "github.com/webdeskltd/log/level"
	r "github.com/webdeskltd/log/record"
	w "github.com/webdeskltd/log/writer"
)

// Prepare configuration
func (self *LogEssence) prepareConfigure(cnf *Configuration) (err error) {
	var i b.BackendName
	var mode b.BackendName
	var n int

	if cnf == nil {
		err = ERROR_CONFIGURATION_IS_NULL
		return
	}

	// Если формат по умолчанию для всех режимов не установлен, используем стандартный формат по умолчанию
	if cnf.Format == "" {
		cnf.Format = default_FORMAT
	}
	b.DefaultFormat = cnf.Format

	// Проверка формата по умолчанию
	_, err = r.CheckFormat(cnf.Format)
	if err != nil {
		return
	}

	for i = range cnf.Mode {
		mode = b.CheckMode(i)
		for n = range cnf.Mode[i] {
			cnf.Mode[i][n] = l.LevelName(strings.ToUpper(string(cnf.Mode[i][n])))
		}
		if mode != i {
			cnf.Mode[mode] = cnf.Mode[i]
			delete(cnf.Mode, i)
		}
	}
	for i = range cnf.Levels {
		cnf.Levels[i] = l.LevelName(strings.ToUpper(string(cnf.Levels[i])))
		mode = b.CheckMode(i)
		if mode != i {
			cnf.Levels[mode] = cnf.Levels[i]
			delete(cnf.Levels, i)
		}
	}
	for i = range cnf.Formats {
		mode = b.CheckMode(i)
		if mode != i {
			cnf.Formats[mode] = cnf.Formats[i]
			delete(cnf.Formats, i)
		}
		// Проверка формата
		_, err = r.CheckFormat(cnf.Formats[mode])
		if err != nil {
			return
		}
	}

	// Graylog2 protocol
	switch strings.ToLower(cnf.Graylog2.Protocol) {
	case "tcp":
		cnf.Graylog2.Protocol = `tcp`
	case "udp":
		cnf.Graylog2.Protocol = `udp`
	default:
		cnf.Graylog2.Protocol = `udp`
	}

	return
}

// Apply new configuration
func (self *LogEssence) Configure(cnf *Configuration) (err error) {
	var bname b.BackendName
	var ok bool

	// Паникаловка :)
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				if err != nil {
					err = errors.New("failed to configure: " + e.Error() + ", " + err.Error())
				} else {
					err = errors.New("failed to configure: " + e.Error())
				}
			}
			self.backend = nil
			self.cnf = nil
		}
	}()

	// Проверка и подготовка конфигурации
	err = self.prepareConfigure(cnf)
	if err != nil {
		return
	}

	// (Ре)Инициализация пула
	self.backend = b.NewBackends()
	if self.interceptStandardLog {
		self.defaultLevelLogWriter = w.NewWriter(default_LEVEL).Resolver(self.ResolveNames).AttachBackends(self.backend)
		stdLogConnect(self.defaultLevelLogWriter)
	}
	for bname = range cnf.Mode {
		var backend *b.Backend
		var bmode b.Mode
		var levels []l.Level
		var level l.Level

		switch bname {
		case b.NAME_CONSOLE:
			backend = b.NewBackendConsole(nil).SetFormat(cnf.Formats[bname])
		case b.NAME_SYSLOG:
			backend = b.NewBackendSyslog(self.HostName).SetFormat(cnf.Formats[bname])
		case b.NAME_FILE:
			var fh *os.File
			if cnf.File == "" {
				panic(ERROR_LOG_FILENAME_IS_EMPTY)
			}
			fh, err = os.OpenFile(cnf.File, syscall.O_APPEND|syscall.O_CREAT|syscall.O_WRONLY, 0644)
			if err != nil {
				panic(ERROR_INIT_FILE_BACKEND)
			}
			backend = b.NewBackendFile(fh).SetFormat(cnf.Formats[bname])
		case b.NAME_GRAYLOG2:
			var hgpc g.GelfProtocolClient
			var hG *g.GelfClient
			if cnf.Graylog2.Protocol == g.UDP_NETWORK {
				hgpc = g.MustUdpClient(cnf.Graylog2.Host, cnf.Graylog2.Port, cnf.Graylog2.ChunkSize)
			}
			if cnf.Graylog2.Protocol == g.TCP_NETWORK {
				hgpc = g.MustTcpClient(cnf.Graylog2.Host, cnf.Graylog2.Port)
			}
			hG = g.NewGelfClient(hgpc, self.cnf.Graylog2.Compression)
			backend = b.NewBackendGraylog2(hG).SetFormat(cnf.Formats[bname])
		case b.NAME_MEMORYPIPE:
			backend = b.NewBackendMemorypipe().SetFormat(cnf.Formats[bname])
		case b.NAME_TELEGRAM:
			backend = b.NewBackendTelegram().SetFormat(cnf.Formats[bname])
		default:
			panic(errors.New(fmt.Sprintf("%s %v", ERROR_UNKNOWN_MODE.Error(), bname)))
		}

		// Устанавливаем уровень или уровни логирования: SetLevel() для NORMAL or SetSelectLevels() для SELECT
		if cnf.Mode[bname] == nil {
			bmode = b.MODE_NORMAL
		} else if len(cnf.Mode[bname]) == 0 {
			bmode = b.MODE_NORMAL
		} else {
			bmode = b.MODE_SELECT
			for n := range cnf.Mode[bname] {
				if _, ok = l.Map2Level[cnf.Mode[bname][n]]; ok {
					levels = append(levels, l.Map2Level[cnf.Mode[bname][n]])
				}
			}
		}
		if bmode == b.MODE_NORMAL {
			if _, ok = cnf.Levels[bname]; ok {
				level = l.Map2Level[cnf.Levels[bname]]
			} else {
				level = default_LEVEL
			}
		}

		// Если backend создан, добавляем его в pool
		if backend != nil {
			self.backend.AddBackend(
				backend.SetLevel(level).
					SetSelectLevels(levels...).
					SetMode(bmode),
			)
		}

	}
	self.cnf = cnf
	return
}
