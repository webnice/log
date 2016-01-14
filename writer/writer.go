package writer

import (
	"fmt"
	"regexp"
	"runtime"

	b "github.com/webdeskltd/log/backends"
	l "github.com/webdeskltd/log/level"
	m "github.com/webdeskltd/log/message"
	r "github.com/webdeskltd/log/record"
	t "github.com/webdeskltd/log/trace"
)

var (
	rexSpaceFirst *regexp.Regexp = regexp.MustCompile(`^[\t\n\f\r ]+`)
	rexSpaceLast  *regexp.Regexp = regexp.MustCompile(`[\t\n\f\r ]+$`)
)

// Writer object for standard logging and etc...
type Writer struct {
	level    l.Level
	resolver r.RecordResolveFn
	backends *b.Backends
}

// Деструктор
func destructor(obj *Writer) {
	obj.backends = nil
	obj.resolver = nil
}

// Create new writer
func NewWriter(level l.Level) (self *Writer) {
	self = new(Writer)
	self.level = level
	runtime.SetFinalizer(self, destructor)
	return
}

func (self *Writer) Resolver(f r.RecordResolveFn) *Writer {
	self.resolver = f
	return self
}

func (self *Writer) AttachBackends(bck *b.Backends) *Writer {
	self.backends = bck
	return self
}

// In the message is deleted whitespace preceding and ends with a message
func (self *Writer) cleanSpace(buf string) string {
	return rexSpaceLast.ReplaceAllString(rexSpaceFirst.ReplaceAllString(buf, ``), ``)
}

// Writer for third party loggers, messages from third-party data loggers are intercepted and sent to this writer
// All messages assigned to the default logging level

func (wr *Writer) write(buf string) (ln int, err error) {
	var msg *m.Message
	msg = m.NewMessage(
		t.NewTrace().
			Trace(t.STEP_BACK + 3).
			GetRecord().
			Resolver(wr.resolver),
	).
		Level(
		l.NewFromMesssage(buf, wr.level).Level,
	).
		Write(wr.cleanSpace(buf))
	if wr.backends != nil {
		wr.backends.Push(msg)
		ln, err = msg.GetResult()
	} else {
		// backend is not initialized, no place to send messages
	}
	return
}

// Writer for []byte
func (wr *Writer) Write(buf []byte) (ln int, err error) {
	return wr.write(string(buf))
}

// Writer for string
func (wr *Writer) WriteString(buf string) (ln int, err error) {
	return wr.write(buf)
}

// Writer for ...any
func (wr *Writer) Println(v ...interface{}) {
	wr.write(fmt.Sprint(v...))
}
