package log

//import "gopkg.in/webnice/debug.v1"
import (
	standardLog "log"
	"os"

	m "gopkg.in/webnice/log.v2/message"
	r "gopkg.in/webnice/log.v2/receiver"
	s "gopkg.in/webnice/log.v2/sender"
	w "gopkg.in/webnice/log.v2/writer"
)

var ess *impl

func init() { ess = newEssence() }

// Create new log object
func newEssence() *impl {
	ess = new(impl)
	ess.writer = w.New()
	ess.sender = s.Gist()
	ess.sender.SetDefaultReceiver(r.Default.Receiver)
	return ess
}

// NewMsg Create new message
func (ess *impl) NewMsg() Log { return m.New() }

// Return writer interface
func (ess *impl) Writer() w.Interface { return ess.writer }

// StandardLogSet Put io writer to log
func (ess *impl) StandardLogSet() Essence {
	standardLog.SetPrefix(``)
	standardLog.SetFlags(0)
	standardLog.SetOutput(ess.writer)
	return ess
}

// StandardLogUnset Reset to defailt
func (ess *impl) StandardLogUnset() Essence {
	standardLog.SetPrefix(``)
	standardLog.SetFlags(standardLog.LstdFlags)
	standardLog.SetOutput(os.Stderr)
	return ess
}
