package fsfile // import "github.com/webdeskltd/log/receiver/fsfile"

import (
	"bytes"
	"fmt"
	"os"

	f "github.com/webdeskltd/log/formater"
	s "github.com/webdeskltd/log/sender"

	"github.com/webdeskltd/log/middleware/fswriter"
)

// const _DefaultTextFORMAT = `%{color}[%{module:-10s}] %{time:2006-01-02T15:04:05.000Z07:00t} (%{level:-8s}): %{message} (%{package}) (%{function}:%{line}) (%{shortfile}:%{line}) (%{longfile})`
const _DefaultTextFORMAT = `(%{colorbeg}%{level:1s}:%{level:1d}%{colorend}): %{message} {%{package}/%{shortfile}:%{line}, func: %{function}()}`

// Interface is an interface of package
type Interface interface {
	// SetFilename Set filename
	SetFilename(string) Interface

	// SetFilemode Set filemode
	SetFilemode(os.FileMode) Interface

	// Receiver Message receiver
	Receiver(s.Message)
}

// impl is an implementation of package
type impl struct {
	Formater f.Interface        // Formater interface
	TplText  string             // Шаблон форматирования текста
	FsWriter fswriter.Interface // Интерфейс записи
}

// New Create new
func New(filename ...string) Interface {
	var fnm string
	var rcv = new(impl)
	rcv.TplText = _DefaultTextFORMAT
	rcv.Formater = f.New()
	rcv.FsWriter = fswriter.New()
	for _, fnm = range filename {
		rcv.FsWriter.SetFilename(fnm)
	}
	return rcv
}

// SetFilename Set filename
func (rcv *impl) SetFilename(fnm string) Interface { rcv.FsWriter.SetFilename(fnm); return rcv }

// SetFilemode Set filemode
func (rcv *impl) SetFilemode(fmd os.FileMode) Interface { rcv.FsWriter.SetFilemode(fmd); return rcv }

// Receiver Message receiver. Output to file
func (rcv *impl) Receiver(msg s.Message) {
	var buf *bytes.Buffer
	var err error
	if buf, err = rcv.Formater.Text(msg, rcv.TplText); err != nil {
		buf = bytes.NewBufferString(fmt.Sprintf("Error formatting log message: %s", err.Error()))
	}
	if _, err = rcv.FsWriter.Write(buf.Bytes()); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
