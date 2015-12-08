package log

import (
	"errors"
	"os"

	b "github.com/webdeskltd/log/backends"
	g "github.com/webdeskltd/log/gelf"
	l "github.com/webdeskltd/log/level"
	r "github.com/webdeskltd/log/record"
	u "github.com/webdeskltd/log/uuid"
	w "github.com/webdeskltd/log/writer"

	//"github.com/webdeskltd/debug"
)

// Initialize default log settings
func init() {
	// Карта всех копий logger
	singleton = make(map[string]Log)

	// Устанавливаем в зависимые пакеты функции информирования об ошибках
	b.LogError = Error

	// Defailt backend format
	b.DefaultFormat = default_FORMAT

	// Default public log object
	var uuid, _ = u.ParseUUID(default_LOGUUID)
	singleton[default_LOGUUID] = newLogEssence(uuid)
}

func newLogEssence(uuid u.UUID) *LogEssence {
	var obj = new(LogEssence)
	obj.Id = uuid
	obj.moduleNames = make(map[string]string)
	obj.Initialize()
	var log = []Log{obj}
	obj.Interface = log[0]
	return obj
}

// New log object
func NewLog() Log {
	return newLogEssence(u.TimeUUID()).Interface
}

// Create dafault configuration
func (log *LogEssence) defaultConfiguration() (cnf *Configuration) {
	if testing_mode_one {
		return
	}
	cnf = &Configuration{
		BufferFlushImmediately: true,
		BufferSize:             0,
		Mode:                   make(map[b.BackendName][]l.LevelName),
		Levels:                 make(map[b.BackendName]l.LevelName),
		Formats:                make(map[b.BackendName]string),
		Format:                 default_FORMAT,
		Graylog2: ConfigurationGraylog2{
			Compression: g.COMPRESSION_NONE,
			Source:      log.HostName,
			Protocol:    g.UDP_NETWORK,
			BufferSize:  1000,
		},
		Telegram: ConfigurationTelegram{},
	}
	cnf.Mode[b.NAME_CONSOLE] = nil
	cnf.Levels[b.NAME_CONSOLE] = l.LevelName(l.Map[l.DEFAULT_LEVEL])
	return
}

// Initialize default configuration
func (log *LogEssence) Initialize() *LogEssence {
	var err error
	var cnf *Configuration

	log.SetApplicationName(``)
	log.HostName, err = os.Hostname()
	if testing_mode_one {
		err = errors.New("Hostname not defined")
	}
	if err != nil {
		log.HostName = `undefined`
	}
	// Create default configuration and apply
	cnf = log.defaultConfiguration()
	err = log.Configure(cnf)
	if err != nil {
		Error("Error Configure(): %v\n", err)
	} else {
		log.ready = true
	}

	// Default level writer
	log.defaultLevelLogWriter = w.NewWriter(l.DEFAULT_LEVEL).Resolver(log.ResolveNames).AttachBackends(log.backend)
	if log.interceptStandardLog {
		stdLogConnect(log.defaultLevelLogWriter)
	}

	return log
}

// Resolve resord
func (self *LogEssence) ResolveNames(rec *r.Record) {
	rec.AppName = self.AppName
	rec.HostName = self.HostName
	if _, ok := self.moduleNames[rec.Package]; ok == true {
		rec.Package = self.moduleNames[rec.Package]
	}
	return
}
