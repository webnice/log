package backends

import (
	"log/syslog"
)

func NewBackendSyslog(hn string) (ret *Backend) {
	var err error
	ret = new(Backend)
	ret.hType = BACKEND_SYSLOG
	ret.wSyslog, err = syslog.New(syslog.LOG_CRIT, hn)
	if err != nil {
		ret = nil
		panic(err)
	}
	return
}
