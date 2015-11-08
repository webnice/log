package log

import (
	"strings"
	"testing"

	b "github.com/webdeskltd/log/backends"
	g "github.com/webdeskltd/log/gelf"
	l "github.com/webdeskltd/log/level"
)

func TestLogPrepareConfigure(t *testing.T) {
	var cnf *Configuration
	var err error
	var log *Log = NewLog()

	cnf = &Configuration{
		BufferFlushImmediately: true,
		BufferSize:             0,
		Mode:                   make(map[b.BackendName][]l.LevelName),
		Levels:                 make(map[b.BackendName]l.LevelName),
		Formats:                make(map[b.BackendName]string),
		Format:                 "",
		Graylog2: ConfigurationGraylog2{
			Compression: g.COMPRESSION_NONE,
			Source:      "test.local",
			Protocol:    g.UDP_NETWORK,
			BufferSize:  1000,
		},
		Telegram: ConfigurationTelegram{},
	}

	err = log.prepareConfigure(cnf)
	if err != nil {
		t.Errorf("Error in prepareConfigure()")
	} else if cnf.Format != default_FORMAT {
		t.Errorf("Error Default_FORMAT in prepareConfigure()")
	}

	cnf.Format = "%{ERROR_FOR_TEST}"
	err = log.prepareConfigure(cnf)
	if err == nil {
		t.Errorf("Error in prepareConfigure()")
	}

	cnf.Format = ""
	cnf.Mode = map[b.BackendName][]l.LevelName{
		b.NAME_CONSOLE:    []l.LevelName{"fatal"},
		b.NAME_FILE:       []l.LevelName{"fatal"},
		b.NAME_GRAYLOG2:   []l.LevelName{"fatal"},
		b.NAME_MEMORYPIPE: []l.LevelName{"fatal"},
		b.NAME_SYSLOG:     []l.LevelName{"fatal"},
		b.NAME_TELEGRAM:   []l.LevelName{"fatal"},
	}
	err = log.prepareConfigure(cnf)
	if err != nil {
		t.Errorf("Error in prepareConfigure()")
	}

	cnf.Mode = map[b.BackendName][]l.LevelName{
		b.BackendName("FOR_TEST"): []l.LevelName{"FOR_TEST"},
	}
	cnf.Levels = map[b.BackendName]l.LevelName{
		b.BackendName("FOR_TEST"): l.LevelName("FOR_TEST"),
	}
	cnf.Formats = map[b.BackendName]string{
		b.BackendName("coNsSOLE"): default_FORMAT,
	}
	err = log.prepareConfigure(cnf)
	if err != nil {
		t.Errorf("Error in prepareConfigure(): %v", err)
	}

	cnf.Formats["file"] = ""
	err = log.prepareConfigure(cnf)
	if err == nil {
		t.Errorf("Error in prepareConfigure()")
	}
	delete(cnf.Formats, "file")

	cnf.Graylog2.Protocol = "TCP"
	err = log.prepareConfigure(cnf)
	if err != nil {
		t.Errorf("Error in prepareConfigure(): %v", err)
	} else if cnf.Graylog2.Protocol != "tcp" {
		t.Errorf("Error in prepareConfigure()")
	}

	cnf.Graylog2.Protocol = ""
	err = log.prepareConfigure(cnf)
	if err != nil {
		t.Errorf("Error in prepareConfigure(): %v", err)
	} else if cnf.Graylog2.Protocol != "udp" {
		t.Errorf("Error in prepareConfigure()")
	}
}

func TestLogConfigure(t *testing.T) {
	var err error
	var log *Log = NewLog()
	var cnf *Configuration

	cnf = &Configuration{
		BufferFlushImmediately: true,
		BufferSize:             0,
		Levels:                 make(map[b.BackendName]l.LevelName),
		Formats:                make(map[b.BackendName]string),
		Format:                 "",
		Graylog2: ConfigurationGraylog2{
			Compression: g.COMPRESSION_NONE,
			Source:      "test.local",
			Protocol:    g.UDP_NETWORK,
			BufferSize:  1000,
		},
		Telegram: ConfigurationTelegram{},
		Mode: map[b.BackendName][]l.LevelName{
			b.NAME_SYSLOG:             []l.LevelName{"fatal"},
			b.BackendName("FOR_TEST"): []l.LevelName{"fatal"},
		},
	}
	err = log.Configure(cnf)
	if err == nil {
		t.Errorf("Error in Configure()")
	} else if strings.Index(err.Error(), ERROR_UNKNOWN_MODE.Error()) < 0 {
		t.Errorf("Error in Configure()")
	}

	cnf.Mode = map[b.BackendName][]l.LevelName{
		b.NAME_CONSOLE: []l.LevelName{},
		b.NAME_SYSLOG:  []l.LevelName{"fatal"},
	}
	err = log.Configure(cnf)
	if err != nil {
		t.Errorf("Error in Configure(): %v", err)
	}

	cnf.Mode = map[b.BackendName][]l.LevelName{
		b.NAME_FILE: []l.LevelName{"fatal"},
	}
	cnf.File = ""
	err = log.Configure(cnf)
	if err == nil {
		t.Errorf("Error in Configure()")
	} else if strings.Index(err.Error(), ERROR_LOG_FILENAME_IS_EMPTY.Error()) < 0 {
		t.Errorf("Error in Configure()")
	}

	cnf.File = "."
	err = log.Configure(cnf)
	if err == nil {
		t.Errorf("Error in Configure()")
	} else if strings.Index(err.Error(), ERROR_INIT_FILE_BACKEND.Error()) < 0 {
		t.Errorf("Error in Configure(): %v", err)
	}

	cnf.File = "test.log"
	err = log.Configure(cnf)
	if err != nil {
		t.Errorf("Error in Configure(): %v", err)
	}

	cnf.File = ""
	cnf.Mode = map[b.BackendName][]l.LevelName{
		b.NAME_MEMORYPIPE: []l.LevelName{"fatal"},
	}
	err = log.Configure(cnf)
	if err != nil {
		t.Errorf("Error in Configure(): %v", err)
	}

	cnf.Mode = map[b.BackendName][]l.LevelName{
		b.NAME_TELEGRAM: []l.LevelName{"fatal"},
	}
	err = log.Configure(cnf)
	if err != nil {
		t.Errorf("Error in Configure(): %v", err)
	}

	cnf.Mode = map[b.BackendName][]l.LevelName{
		b.NAME_GRAYLOG2: []l.LevelName{"fatal"},
	}
	err = log.Configure(cnf)
	if err != nil {
		t.Errorf("Error in Configure(): %v", err)
	}

	cnf.Mode = map[b.BackendName][]l.LevelName{
		b.NAME_GRAYLOG2: []l.LevelName{"fatal"},
	}
	cnf.Graylog2.Protocol = "tcp"
	err = log.Configure(cnf)
	if err != nil {
		t.Errorf("Error in Configure(): %v", err)
	}

}
