package writer

// import "gopkg.in/webnice/debug.v1"
import (
	"bytes"
	"fmt"

	l "gopkg.in/webnice/log.v2/level"
	m "gopkg.in/webnice/log.v2/message"
)

// New Create new object
func New() Interface {
	var wr = new(impl)
	wr.level = l.Defailt
	return wr
}

// Writer for []byte
func (wr *impl) Write(buf []byte) (ln int, err error) {
	return wr.Writer(bytes.NewBuffer(buf))
}

// Writer for string
func (wr *impl) WriteString(buf string) (ln int, err error) {
	return wr.Writer(bytes.NewBufferString(buf))
}

// Writer for ...any
func (wr *impl) Println(args ...interface{}) {
	wr.Writer(bytes.NewBufferString(fmt.Sprint(args...)))
}

// CleanSpace In the message is deleted whitespace preceding and ends with a message
func (wr *impl) CleanSpace(buf *bytes.Buffer) (ret *bytes.Buffer) {
	ret = bytes.NewBufferString(rexSpaceLast.ReplaceAllString(rexSpaceFirst.ReplaceAllString(buf.String(), ``), ``))
	return
}

// Writer method
func (wr *impl) Writer(buf *bytes.Buffer) (ln int, err error) {
	ln = buf.Len()
	buf = wr.CleanSpace(buf)
	m.New().CallStackCorrect(2).Message(wr.level, buf.String())
	return
}
