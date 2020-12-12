package stdout

import (
	"bytes"
	"fmt"
	"os"

	f "github.com/webnice/lv2/formater"
	s "github.com/webnice/lv2/sender"
)

// const defaultTextFORMAT = `%{color}[%{module:-10s}] %{time:2006-01-02T15:04:05.000Z07:00t} (%{level:-8s}): %{message} (%{package}) (%{function}:%{line}) (%{shortfile}:%{line}) (%{longfile})`
const defaultTextFORMAT = `(%{colorbeg}%{level:1s}:%{level:1d}%{colorend}): %{message} {%{package}/%{shortfile}:%{line}, func: %{function}()}`

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
	var rcv = &impl{
		Formater: f.New(),
		TplText:  defaultTextFORMAT,
	}

	return rcv
}

// Receiver Message receiver. Output to STDOUT
func (rcv *impl) Receiver(msg s.Message) {
	var (
		err error
		buf *bytes.Buffer
	)

	if buf, err = rcv.Formater.Text(msg, rcv.TplText); err != nil {
		buf = bytes.NewBufferString(fmt.Sprintf("formatting log message error: %s", err))
	}
	_, _ = fmt.Fprintln(os.Stdout, buf.String())
}
