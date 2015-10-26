package backends

import (
	"os"
)

func NewBackendFile(f *os.File) (ret *Backend) {
	ret = new(Backend)
	ret.hType = BACKEND_FILE
	ret.fH = f
	ret.fH.Seek(0, 2)
	return
}
