package writer

import (
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

// Writer for []byte
func (self *Writer) Write(buf []byte) (ln int, err error) {
	var msg *m.Message
	msg = m.NewMessage(
		t.NewTrace().
			Trace(t.STEP_BACK + 2).
			GetRecord().
			Resolver(self.resolver),
	).Level(self.level).
		Write(self.cleanSpace(string(buf)))
	if self.backends != nil {
		self.backends.Push(msg)
		ln, err = msg.GetResult()
	} else {
		// backend не инициализирован, отправлять сообщения некуда
	}
	return
}

// Writer for string
func (self *Writer) WriteString(buf string) (ln int, err error) {
	var msg *m.Message
	msg = m.NewMessage(
		t.NewTrace().
			Trace(t.STEP_BACK + 2).
			GetRecord().
			Resolver(self.resolver),
	).Level(self.level).
		Write(self.cleanSpace(buf))
	if self.backends != nil {
		self.backends.Push(msg)
		ln, err = msg.GetResult()
	} else {
		// backend не инициализирован, отправлять сообщения некуда
	}
	return
}
