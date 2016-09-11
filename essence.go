package log // import "github.com/webdeskltd/log"

//import "github.com/webdeskltd/debug"
import (
	"bytes"
	"fmt"
	standardLog "log"
	"os"

	f "github.com/webdeskltd/log/formater"
	m "github.com/webdeskltd/log/message"
	s "github.com/webdeskltd/log/sender"
	w "github.com/webdeskltd/log/writer"
)

var ess *impl

func init() {
	ess = newEssence()
}

// Create new log object
func newEssence() *impl {
	ess = new(impl)
	ess.formater = f.New()
	ess.writer = w.New()
	ess.sender = s.Gist()
	ess.sender.SetDefaultReceiver(ess.DefaultReceiver)
	ess.tplText = defaultTextFORMAT
	return ess
}

// NewMsg Create new message
func (ess *impl) NewMsg() Log {
	return m.New()
}

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

// DefaultReceiver Default message receiver
func (ess *impl) DefaultReceiver(msg s.Message) {
	var buf *bytes.Buffer
	var err error
	if buf, err = ess.formater.Text(msg, ess.tplText); err != nil {
		fmt.Fprintf(os.Stderr, "Error formationg log message: %s", err.Error())
		return
	}
	fmt.Fprintln(os.Stderr, buf.String())
}
