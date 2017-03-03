package stderr

import (
	"bytes"
	"fmt"
	"os"

	f "gopkg.in/webnice/log.v2/formater"
	s "gopkg.in/webnice/log.v2/sender"
)

// const _DefaultTextFORMAT = `%{color}[%{module:-10s}] %{time:2006-01-02T15:04:05.000Z07:00t} (%{level:-8s}): %{message} (%{package}) (%{function}:%{line}) (%{shortfile}:%{line}) (%{longfile})`
const _DefaultTextFORMAT = `(%{colorbeg}%{level:1s}:%{level:1d}%{colorend}): %{message} {%{package}/%{shortfile}:%{line}, func: %{function}()}`

// Interface is an interface of package
type Interface interface {
	// Receiver Message receiver
	Receiver(s.Message)
}

// impl is an implementation of package
type impl struct {
	Formater f.Interface // Formater interface
	TplText  string      // Шаблон форматирования текста
}

// New Create new
func New() Interface {
	var rcv = new(impl)
	rcv.Formater = f.New()
	rcv.TplText = _DefaultTextFORMAT
	return rcv
}

// Receiver Message receiver. Output to STDERR
func (rcv *impl) Receiver(msg s.Message) {
	var buf *bytes.Buffer
	var err error
	if buf, err = rcv.Formater.Text(msg, rcv.TplText); err != nil {
		buf = bytes.NewBufferString(fmt.Sprintf("Error formatting log message: %s", err.Error()))
	}
	fmt.Fprintln(os.Stderr, buf.String())
}
