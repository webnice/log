package fsfile // import "github.com/webdeskltd/log/receiver/fsfile"

import (
	"fmt"
	"os"

	s "github.com/webdeskltd/log/sender"

	"github.com/webdeskltd/log/middleware"
	"github.com/webdeskltd/log/middleware/fswformattext"
)

// const _DefaultTextFORMAT = `%{color}[%{module:-10s}] %{time:2006-01-02T15:04:05.000Z07:00t} (%{level:-8s}): %{message} (%{package}) (%{function}:%{line}) (%{shortfile}:%{line}) (%{longfile})`
const _DefaultTextFORMAT = `(%{colorbeg}%{level:1s}:%{level:1d}%{colorend}): %{message} {%{package}/%{shortfile}:%{line}, func: %{function}()}`

// Interface is an interface of package
type Interface interface {
	// SetFilename Set filename
	SetFilename(string) Interface

	// SetFilemode Set filemode
	SetFilemode(os.FileMode) Interface

	// SetFormat Set template line formating
	SetFormat(string) Interface

	// Receiver Message receiver
	Receiver(s.Message)
}

// impl is an implementation of package
type impl struct {
	TplText  string              // Шаблон форматирования текста
	FsWriter middleware.FsWriter // Интерфейс записи
}

// New Create new
func New(filename ...string) Interface {
	var fnm string
	var rcv = new(impl)
	rcv.TplText = _DefaultTextFORMAT
	rcv.FsWriter = fswformattext.New()
	for _, fnm = range filename {
		rcv.FsWriter.SetFilename(fnm)
	}
	return rcv
}

// SetFilename Set filename
func (rcv *impl) SetFilename(fnm string) Interface { rcv.FsWriter.SetFilename(fnm); return rcv }

// SetFilemode Set filemode
func (rcv *impl) SetFilemode(fmd os.FileMode) Interface { rcv.FsWriter.SetFilemode(fmd); return rcv }

// SetFormat Set template line formating
func (rcv *impl) SetFormat(format string) Interface { rcv.TplText = format; return rcv }

// SetFsWriter Установка функции записи в файл с форматированием
func (rcv *impl) SetFsWriter(fn middleware.FsWriter) Interface { rcv.FsWriter = fn; return rcv }

// Receiver Message receiver. Output to file
func (rcv *impl) Receiver(msg s.Message) {
	var err error
	rcv.FsWriter.SetFormat(rcv.TplText)
	if _, err = rcv.FsWriter.Write(msg); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
