package backends

import (
	"time"

	"github.com/webdeskltd/debug"
)

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
func (self *Backend) SetLevel(level Level) *Backend {
	self.hLevelNormal = level
	return self
}

// Set selected logging levels for Mode = MODE_SELECT
func (self *Backend) SetSelectLevels(levels ...Level) *Backend {
	var i int
	self.hLevelSelect = make([]Level, len(levels))
	for i = range levels {
		self.hLevelSelect[i] = levels[i]
	}
	return self
}

// Блокируемая функция сброса буффера бэкэнда и остановки логирования
func (self *Backend) Stop() (err error) {
	debug.Dumper("Stop() start", self.hLevelNormal, self.hLevelSelect, self.hMode, self.hType)
	time.Sleep(time.Second * 2)
	debug.Dumper("Stop() end")
	return
}
