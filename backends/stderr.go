package backends

import (
	"os"
)

func NewBackendSTD(f *os.File) (ret *Backend) {
	ret = new(Backend)
	ret.hType = BACKEND_STD
	if f != nil {
		ret.fH = f
	} else {
		ret.fH = os.Stderr
	}
	return
}
