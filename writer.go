package log

import (
	//"fmt"

	"github.com/webdeskltd/log/logging"
)

func (this *log) Write(buf []byte) (l int, err error) {
	var tmp []byte
	//	var i int

	if this.Record == nil {
		this.Record = newTrace().Trace(traceStepBack + 2).GetRecord().SetLevel(int8(defaultLevel)).SetMessage(string(buf)).Complete()
	}

	//	tmp = []byte(fmt.Sprintf("- [%s] [%s] [%s:%d] ", this.Trace.Package, this.Trace.Function, this.Trace.File, this.Trace.Line))
	//	i, err = self.Writer.Write(tmp)
	//	l += i
	//	if err != nil {
	//		return
	//	}
	//	i, err = self.Writer.Write(buf)
	//	l += i

	lll := logging.MustGetLogger("123")
	lll.Critical("%s", tmp)

	// Flush
	if self.BufferFlushImmediately {
		this.Flush()
	}
	this.Flush()

	return
}

func (this *log) WriteString(buf string) (l int, err error) {
	l, err = this.Write([]byte(buf + "\n"))
	return
}

func (this *log) Flush() (err error) {
	err = self.Writer.Flush()
	return
}
