package fswformattext

//import "gopkg.in/webnice/debug.v1"
import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"gopkg.in/webnice/log.v2/middleware"

	f "gopkg.in/webnice/log.v2/formater"
	s "gopkg.in/webnice/log.v2/sender"
)

// New Create new package implementation and return interface
func New(filename ...string) middleware.FsWriter {
	var fnm string
	var fsw = new(impl)
	var tmp = strings.Split(os.Args[0], string(os.PathSeparator))
	if len(tmp) > 0 {
		fsw.SetFilename(tmp[len(tmp)-1] + `.log`)
	} else {
		fsw.SetFilename(os.Args[0] + `.log`)
	}
	for _, fnm = range filename {
		fsw.SetFilename(fnm)
	}
	fsw.Formater = f.New()
	fsw.TplText = _DefaultTextFORMAT
	return fsw
}

// SetFilename Set filename
func (fsw *impl) SetFilename(filename string) middleware.FsWriter {
	fsw.Lock()
	defer fsw.Unlock()
	fsw.Filename = filename
	return fsw
}

// SetFilemode Set filemode
func (fsw *impl) SetFilemode(filemode os.FileMode) middleware.FsWriter {
	fsw.Filemode = filemode
	return fsw
}

// SetFormat Set template line formating
func (fsw *impl) SetFormat(format string) middleware.FsWriter { fsw.TplText = format; return fsw }

// WriteMessage Запись среза байт в файл
func (fsw *impl) WriteMessage(msg s.Message) (n int, err error) {
	var buf *bytes.Buffer
	if buf, err = fsw.Formater.Text(msg, fsw.TplText); err != nil {
		buf = bytes.NewBufferString(fmt.Sprintf("Error formatting log message: %s", err.Error()))
	}
	buf.WriteString("\r\n")
	n, err = fsw.Write(buf.Bytes())
	return
}

// Write Запись среза байт как есть
func (fsw *impl) Write(buf []byte) (n int, err error) {
	var out *os.File
	fsw.Lock()
	defer fsw.Unlock()
	if out, err = os.OpenFile(fsw.Filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.FileMode(0644)); err != nil {
		err = fmt.Errorf("Failed to open file '%s': %s", fsw.Filename, err.Error())
		return
	}
	defer out.Close()
	n, err = out.Write(buf)
	return
}
