package log

import (
	"bufio"
	"os"
	"time"
	//"fmt"
	
	//"github.com/webdeskltd/debug"
)

// Initialize default log settings
func init() {
	self = new(configuration)
	self.Writer = bufio.NewWriter(os.Stderr)
}

func newLog() (this *log) {
	this = new(log)
	this.Writer = self.Writer
	this.Now = time.Now().In(time.Local)
	this.Trace = newTrace().Trace(traceStepBack + 1)
	return
}

// All level writer
func (this *log) write(level logLevel, tpl string, args ...interface{}) *log {
	self.Writer.WriteString(this.Trace.Package + " " + tpl + "\n")
	
//	debug.Dumper(tpl)
//	fmt.Println("LOG:", tpl, this.Trace.Package, this.Trace.Function)


	if self.FlushImmediately {
		self.Writer.Flush()
	}	
	return this
}
