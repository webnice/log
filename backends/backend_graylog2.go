package backends // import "github.com/webdeskltd/log/backends"

import (
	"errors"

	g "github.com/webdeskltd/log/gelf"
)

func NewBackendGraylog2(gelf *g.GelfClient) (ret *Backend) {
	ret = new(Backend)
	ret.hType = BACKEND_GRAYLOG2
	if gelf == nil {
		ret = nil
		panic(errors.New("Call NewBackendGraylog2(nil)"))
	}
	ret.cGraylog2 = gelf

	return
}
