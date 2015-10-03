package log

import (
	//"bufio"
	"time"
)

func newLog() (self *log) {
	self = new(log)
	self.Now = time.Now().In(time.Local)
	self.Trace = newTrace().Trace(traceStepBack + 1)
	return
}


// All level writer
func (self *log) write(level logLevel, tpl string, args ...interface{}) {
	//debug.Dumper(self)
}
