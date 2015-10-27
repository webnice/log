package backends

import (
	"os"

	m "github.com/webdeskltd/log/message"

	"github.com/webdeskltd/debug"
)

func init() {
	debug.Nop()
}

func NewBackendConsole(f *os.File) (ret *Backend) {
	ret = new(Backend)
	ret.reader = ret.readerConsole
	ret.hType = BACKEND_CONSOLE
	if f != nil {
		ret.fH = f
	} else {
		ret.fH = os.Stderr
	}
	return
}

func (self *Backend) readerConsole(msg *m.Message) {
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
	txt = txt
	//print(self.format); print("\n")
	print(txt); print("\n")

}
