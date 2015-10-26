package backends

import (
	"os"
)

func NewBackendConsole(f *os.File) (ret *Backend) {
	ret = new(Backend)
	ret.hType = BACKEND_CONSOLE
	if f != nil {
		ret.fH = f
	} else {
		ret.fH = os.Stderr
	}
	return
}
