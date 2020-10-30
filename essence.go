package log // import "github.com/webnice/log"

import (
	standardLog "log"
	"os"

	m "github.com/webnice/log/message"
	r "github.com/webnice/log/receiver"
	s "github.com/webnice/log/sender"
	w "github.com/webnice/log/writer"
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
