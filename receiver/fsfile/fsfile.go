package fsfile

import (
	"fmt"
	"os"

	s "github.com/webnice/lv2/sender"

	"github.com/webnice/lv2/middleware"
	"github.com/webnice/lv2/middleware/fswformattext"
)

// const defaultTextFORMAT = `%{color}[%{module:-10s}] %{time:2006-01-02T15:04:05.000Z07:00t} (%{level:-8s}): %{message} (%{package}) (%{function}:%{line}) (%{shortfile}:%{line}) (%{longfile})`
const defaultTextFORMAT = `(%{colorbeg}%{level:1s}:%{level:1d}%{colorend}): %{message} {%{package}/%{shortfile}:%{line}, func: %{function}()}`

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

	// Write Запись среза байт
	Write([]byte) (int, error)
}

// impl is an implementation of package
type impl struct {
	TplText  string              // Шаблон форматирования текста
	FsWriter middleware.FsWriter // Интерфейс записи
}

// New Create new
func New(filename ...string) Interface {
	var (
		rcv *impl
		fnm string
	)

	rcv = new(impl)
	rcv.TplText = defaultTextFORMAT
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
	if _, err = rcv.FsWriter.WriteMessage(msg); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	}
}

// Write Запись среза байт
func (rcv *impl) Write(buf []byte) (n int, err error) {
	if n, err = rcv.FsWriter.Write(buf); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	}

	return
}
