package backends

import ()

func NewBackendMemorypipe() (ret *Backend) {
	ret = new(Backend)
	ret.hType = BACKEND_MEMORYPIPE

	return
}
