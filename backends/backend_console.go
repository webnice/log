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
	txt, err = msg.Record.Format(self.format)
	if err != nil {
		if LogError != nil {
			LogError("Error Record.Format(): %v", err)
		}
	}
	self.fH.WriteString(txt + "\n")
}
