package trace // import "github.com/webnice/log/v2/trace"

import (
	"runtime"
	"strings"
)

// New Create new object
func New() Interface {
	var trc = new(impl)
	return trc
}

// Trace Get call information with stack back level
func (trc *impl) Trace(stackLevel int) Interface {
	var (
		ok  bool
		pc  uintptr
		fn  *runtime.Func
		buf []byte
		tmp []string
		i   int
	)

	if stackLevel == 0 {
		stackLevel = stackBack
	}
	buf = make([]byte, 1<<16)
	pc, trc.info.FileNameLong, trc.info.FileLine, ok = runtime.Caller(stackLevel)
	if ok == true {
		fn = runtime.FuncForPC(pc)
		if fn != nil {
			trc.info.Function = fn.Name()
		}
		i = runtime.Stack(buf, true)
		trc.info.CallStack = string(buf[:i])

		tmp = strings.Split(trc.info.Function, packageSeparator)
		if len(tmp) > 1 {
			trc.info.Package += strings.Join(tmp[:len(tmp)-1], packageSeparator)
			trc.info.Function = tmp[len(tmp)-1]
		}
		tmp = strings.SplitN(trc.info.Function, `.`, 2)
		if len(tmp) == 2 {
			if trc.info.Package != "" {
				trc.info.Package += packageSeparator
			}
			trc.info.Package += tmp[0]
			trc.info.Function = tmp[1]
		}

		// Filename short
		tmp = strings.Split(trc.info.FileNameLong, packageSeparator)
		if len(tmp) > 0 {
			trc.info.FileNameShort = tmp[len(tmp)-1]
		}

		// Module name
		tmp = strings.Split(trc.info.Package, packageSeparator)
		if len(tmp) > 0 {
			trc.info.Module = tmp[len(tmp)-1]
		}
	}

	return trc
}

// Info Return trace information
func (trc *impl) Info() *Info { return &trc.info }
