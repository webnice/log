package log

import (
	"runtime"
	"strings"

	"github.com/webdeskltd/log/record"
)

const (
	traceStepBack    int    = 2
	packageSeparator string = `/`
)

// Trace
type trace struct {
	Record *record.Record
}

func newTrace() *trace {
	return new(trace)
}

func (this *trace) Trace(level int) *trace {
	var ok bool
	var pc uintptr
	var fn *runtime.Func
	var buf []byte
	var tmp []string
	var i int

	this.Record = record.NewRecord()
	if level == 0 {
		level = traceStepBack
	}
	buf = make([]byte, 1<<16)
	pc, this.Record.FileNameLong, this.Record.FileLine, ok = runtime.Caller(level)
	if ok == true {
		fn = runtime.FuncForPC(pc)
		if fn != nil {
			this.Record.Function = fn.Name()
		}
		i = runtime.Stack(buf, true)
		this.Record.CallStack = string(buf[:i])

		tmp = strings.Split(this.Record.Function, packageSeparator)
		if len(tmp) > 1 {
			this.Record.Package += strings.Join(tmp[:len(tmp)-1], packageSeparator)
			this.Record.Function = tmp[len(tmp)-1]
		}
		tmp = strings.SplitN(this.Record.Function, `.`, 2)
		if len(tmp) == 2 {
			if this.Record.Package != "" {
				this.Record.Package += packageSeparator
			}
			this.Record.Package += tmp[0]
			this.Record.Function = tmp[1]
		}

		// Filename short
		tmp = strings.Split(this.Record.FileNameLong, packageSeparator)
		if len(tmp) > 0 {
			this.Record.FileNameShort = tmp[len(tmp)-1]
		}

		// Module name
		tmp = strings.Split(this.Record.Package, packageSeparator)
		if len(tmp) > 0 {
			this.Record.Module = tmp[len(tmp)-1]
		}

		// Custom module name
		if _, ok = self.moduleNames[this.Record.Package]; ok == true {
			this.Record.Module = self.moduleNames[this.Record.Package]
		}
	}
	return this
}

func (this *trace) GetRecord() *record.Record {
	if self.AppName != "" {
		this.Record.AppName = self.AppName
	}
	if self.HostName != "" {
		this.Record.HostName = self.HostName
	}
	return this.Record
}
