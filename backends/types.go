package backends

import (
	"container/list"
	"errors"
	"os"

	l "github.com/webdeskltd/log/level"
	"github.com/webdeskltd/log/record"
	"github.com/webdeskltd/log/uuid"
)

const (
	BACKEND_STD      Type = iota // 0 Logging to STDERR or STDOUT
	BACKEND_FILE                 // 1 Logging to file
	BACKEND_GRAYLOG2             // 2 Logging to graylog2 server
	BACKEND_MEMORY               // 3 Logging to memory pool
	BACKEND_TELEGRAM             // 4 Logging to messenger telegram
)

const (
	MODE_NORMAL Mode = iota // 0 Публикуются сообщения начиная от выбранного уровня и ниже. Напримр выбран NOTICE, публковаться будут FATAL, ALERT, CRITICAL, ERROR, WARNING, NOTICE, игнорироваться INFO, DEBUG
	MODE_SELECT             // 1 Публикуются сообщения только выбранных уровней
)

const lengthRecords int = 20000 // The maximum size of the channel buffering incoming messages

type Type int  // Type of backend
type Mode int8 // Mode filtering messages on the level of logging

type Backends struct {
	RecordsChan chan *record.Record         // Buffered channel incoming messages
	Pool        *list.List                  // Pool of backends
	PoolIndex   map[uuid.UUID]*list.Element // Map of IDs backends
}

type Backend struct {
	hType        Type      // Type of backend
	hMode        Mode      // Mode filtering messages on the level of logging
	hLevelNormal l.Level   // Logging level (Mode = MODE_NORMAL)
	hLevelSelect []l.Level // The selected logging levels for a backend (Mode = MODE_SELECT)
	fH           *os.File  // FH for STD and File backend
}

var (
	typeName = map[Type]string{
		BACKEND_STD:      `STD`,
		BACKEND_FILE:     `FILE`,
		BACKEND_GRAYLOG2: `GRAYLOG2`,
		BACKEND_MEMORY:   `MEMORY`,
		BACKEND_TELEGRAM: `TELEGRAM`,
	}
	modeName = map[Mode]string{
		MODE_NORMAL: `NORMAL`,
		MODE_SELECT: `SELECT`,
	}
)

var (
	ErrBackendIdIsEmpty error = errors.New(`Backend id is empty`)
	ErrBackendNotFound  error = errors.New(`Backend not found`)
	ErrBackendIsNull    error = errors.New(`Passed object is null`)
)
