package message

import (
	"regexp"

	"github.com/webdeskltd/debug"
	l "github.com/webdeskltd/log/level"
	"github.com/webdeskltd/log/record"
	t "github.com/webdeskltd/log/trace"
)

var (
	rexSpaceFirst *regexp.Regexp = regexp.MustCompile(`^[\t\n\f\r ]`)
	rexSpaceLast  *regexp.Regexp = regexp.MustCompile(`[\t\n\f\r ]$`)
)

// Writer object for standard logging and etc...
type Writer struct {
	defaultLevel      l.Level
	resolveModuleName ResolveModuleNameFn
}

// Функция разрешения имени модуля
type ResolveModuleNameFn func(r *record.Record) string

func init() {
	debug.Nop()
}

// Create new writer
func NewWriter(level l.Level) (self *Writer) {
	self = new(Writer)
	self.defaultLevel = level
	return
}

func (self *Writer) Resolver(f ResolveModuleNameFn) *Writer {
	if f != nil {
		self.resolveModuleName = f
	}
	return self
}

// Сюда попадают все сообщения от сторонних логеров
// У сообщений удаляются пробельные символы, предшествующие и заканчивающие сообщение
// Сообщениям присваивается дефолтовый уровень логирования

// Writer for []byte
func (self *Writer) Write(buf []byte) (l int, err error) {
	var msg *Message = NewMessage(t.NewTrace().Trace(t.STEP_BACK + 2).GetRecord())
	if self.resolveModuleName != nil {
		msg.Record.Package = self.resolveModuleName(msg.Record)
	}
	msg.Write(self.defaultLevel, rexSpaceLast.ReplaceAllString(rexSpaceFirst.ReplaceAllString(string(buf), ``), ``))
	l = msg.WriteLen
	err = msg.WriteErr

	debug.Dumper("[]byte: ", buf, msg)

	return
}

// Writer for string
func (self *Writer) WriteString(buf string) (l int, err error) {
	var msg *Message = NewMessage(t.NewTrace().Trace(t.STEP_BACK + 2).GetRecord())
	if self.resolveModuleName != nil {
		msg.Record.Package = self.resolveModuleName(msg.Record)
	}
	msg.Write(self.defaultLevel, rexSpaceLast.ReplaceAllString(rexSpaceFirst.ReplaceAllString(buf, ``), ``))
	l = msg.WriteLen
	err = msg.WriteErr

	debug.Dumper("string: ", buf, msg)

	return
}
