package stdout // import "github.com/webdeskltd/log/receiver/stdout"

import (
	"bytes"
	"fmt"
	"os"

	f "github.com/webdeskltd/log/formater"
	s "github.com/webdeskltd/log/sender"
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

// Receiver Message receiver. Output to STDOUT
func (rcv *impl) Receiver(msg s.Message) {
	var buf *bytes.Buffer
	var err error
	if buf, err = rcv.Formater.Text(msg, rcv.TplText); err != nil {
		fmt.Fprintf(os.Stdout, "Error formationg log message: %s", err.Error())
		return
	}
	fmt.Fprintln(os.Stdout, buf.String())
}
