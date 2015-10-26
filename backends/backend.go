package backends

import (
	//"time"
	"fmt"

	"github.com/webdeskltd/debug"
	l "github.com/webdeskltd/log/level"
)

func init() {
	debug.Nop()
}

// Set mode filtering messages on the level of logging
// 	MODE_NORMAL - Публикуются сообщения начиная от выбранного уровня и ниже. Напримр выбран NOTICE, публковаться будут FATAL, ALERT, CRITICAL, ERROR, WARNING, NOTICE, игнорироваться INFO, DEBUG
//	MODE_SELECT - Публикуются сообщения только выбранных уровней
func (self *Backend) SetMode(mode Mode) *Backend {
	self.hMode = mode
	return self
}

// Set mode = MODE_NORMAL
func (self *Backend) SetModeNormal() *Backend {
	self.hMode = MODE_NORMAL
	return self
}

// Set logging level for Mode = MODE_NORMAL
func (self *Backend) SetLevel(level l.Level) *Backend {
	self.hLevelNormal = level
	return self
}

// Set selected logging levels for Mode = MODE_SELECT
func (self *Backend) SetSelectLevels(levels ...l.Level) *Backend {
	var i int
	self.hLevelSelect = make([]l.Level, len(levels))
	for i = range levels {
		self.hLevelSelect[i] = levels[i]
	}
	return self
}

func (self *Backend) SetFormat(format string) *Backend {
	if format != "" {
		self.format = format
	} else {
		self.format = DefaultFormat
	}
	return self
}

// Блокируемая функция сброса буффера бэкэнда и остановки логирования
func (self *Backend) Stop() (err error) {
	switch self.hType {
	case BACKEND_FILE:
		if self.fH != nil {
			// Close file if open
			err = self.fH.Close()
		}
	}

	fmt.Sprintf("")
	print(fmt.Sprintf("Stop() backend: '%s', mode: '%s'%v|%v\n", MapTypeName[self.hType], modeName[self.hMode], l.Map[self.hLevelNormal], self.hLevelSelect))
	return
}
