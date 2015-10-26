package backends

import ()

func NewBackendTelegram() (ret *Backend) {
	ret = new(Backend)
	ret.hType = BACKEND_TELEGRAM

	return
}
