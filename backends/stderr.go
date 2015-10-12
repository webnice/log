package backends

import (
	"os"
)

func NewBackendSTD(f *os.File) (ret *Backend) {
	ret = new(Backend)
	ret.hType = BACKEND_STD
	if f == nil {
		f = os.Stderr
	} else if f != os.Stderr && f != os.Stdout {
		f = os.Stderr
	}
	return
}
