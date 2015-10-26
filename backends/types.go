package backends

import (
	"container/list"
	"errors"
	"log/syslog"
	"os"

	g "github.com/webdeskltd/log/gelf"
	l "github.com/webdeskltd/log/level"
	m "github.com/webdeskltd/log/message"
	u "github.com/webdeskltd/log/uuid"
)

const (
	BACKEND_CONSOLE    Type = iota // 0 Logging to STDERR or STDOUT
	BACKEND_SYSLOG                 // 1 Logging to Syslog
	BACKEND_FILE                   // 1 Logging to file
	BACKEND_GRAYLOG2               // 2 Logging to graylog2 server
	BACKEND_MEMORYPIPE             // 3 Logging to memory
	BACKEND_TELEGRAM               // 4 Logging to messenger telegram
)

const (
	NAME_CONSOLE    BackendName = "console"
	NAME_SYSLOG     BackendName = "syslog"
	NAME_FILE       BackendName = "file"
	NAME_GRAYLOG2   BackendName = "graylog2"
	NAME_MEMORYPIPE BackendName = "memorypipe"
	NAME_TELEGRAM   BackendName = "telegram"
)

const (
	MODE_NORMAL Mode = iota // 0 Публикуются сообщения начиная от выбранного уровня и ниже. Напримр выбран NOTICE, публковаться будут FATAL, ALERT, CRITICAL, ERROR, WARNING, NOTICE, игнорироваться INFO, DEBUG
	MODE_SELECT             // 1 Публикуются сообщения только выбранных уровней
)

const lengthRecords int = 20000 // The maximum size of the channel buffering incoming messages

var (
	MapTypeName = map[Type]BackendName{
		BACKEND_CONSOLE:    NAME_CONSOLE,
		BACKEND_SYSLOG:     NAME_SYSLOG,
		BACKEND_FILE:       NAME_FILE,
		BACKEND_GRAYLOG2:   NAME_GRAYLOG2,
		BACKEND_MEMORYPIPE: NAME_MEMORYPIPE,
		BACKEND_TELEGRAM:   NAME_TELEGRAM,
	}
	modeName = map[Mode]string{
		MODE_NORMAL: `NORMAL`,
		MODE_SELECT: `SELECT`,
	}
	DefaultFormat string = `` // Set from package log function init()
)

var (
	ErrBackendIdIsEmpty error = errors.New(`Backend id is empty`)
	ErrBackendNotFound  error = errors.New(`Backend not found`)
	ErrBackendIsNull    error = errors.New(`Passed object is null`)
)

type Type int           // Type of backend
type Mode int8          // Mode filtering messages on the level of logging
type BackendName string // Mode name in configuration

type Backends struct {
	RecordsChan chan *m.Message          // Buffered channel incoming messages
	Pool        *list.List               // Pool of backends
	PoolIndex   map[u.UUID]*list.Element // Map of IDs backends
	exitChan    chan bool
	doneChan    chan bool
}

type Backend struct {
	hType        Type           // Type of backend
	hMode        Mode           // Mode filtering messages on the level of logging
	hLevelNormal l.Level        // Logging level (Mode = MODE_NORMAL)
	hLevelSelect []l.Level      // The selected logging levels for a backend (Mode = MODE_SELECT)
	fH           *os.File       // FH for CONSOLE and File backend
	format       string         // String format for backend
	hostname     string         // Hostname
	wSyslog      *syslog.Writer // Writer for syslog backend
	cGraylog2    *g.GelfClient  // Graylog2 client connection
}
