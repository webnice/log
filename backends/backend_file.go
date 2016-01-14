package backends

//import "github.com/webdeskltd/debug"
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

	txt, err = msg.Record.Format(self.format)
	if err != nil {
		if LogError != nil {
			LogError("Error Record.Format(): %v", err)
		}
	}
	self.fH.WriteString(txt + "\n")
}
