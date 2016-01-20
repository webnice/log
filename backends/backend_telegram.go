package backends // import "github.com/webdeskltd/log/backends"

import ()

func NewBackendTelegram() (ret *Backend) {
	ret = new(Backend)
	ret.hType = BACKEND_TELEGRAM

	return
}
