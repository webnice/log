package backends

import (
	"os"

	m "github.com/webdeskltd/log/message"
)

func NewBackendFile(f *os.File) (ret *Backend) {
	ret = new(Backend)
	ret.reader = ret.readerFile
	ret.hType = BACKEND_FILE
	ret.fH = f
	ret.fH.Seek(0, 2)
	return
}

func (self *Backend) readerFile(msg *m.Message) {
	var txt string
	var err error

	// Форматируем сообщение
	txt, err = msg.Record.Format(self.format)
	// Ошибка не должна никогда возникать так как формат проверяется при конфигуриговании
	// Но лучше перебдеть и проинформировать, чем недобдеть
	if err != nil {
		if LogError != nil {
			LogError("Error Record.Format(): %v", err)
		}
	}

	print(txt)
	print("\n")
}
