package fswformatjson // import "github.com/webdeskltd/log/middleware/fswformatjson"

//import "github.com/webdeskltd/debug"
import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/webdeskltd/log/middleware"

	s "github.com/webdeskltd/log/sender"
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
func (fsw *impl) SetFormat(f string) middleware.FsWriter { return fsw }

// Write Запись среза байт в файл
func (fsw *impl) Write(msg s.Message) (n int, err error) {
	var out *os.File
	var buf *bytes.Buffer
	fsw.Lock()
	defer fsw.Unlock()
	buf = bytes.NewBufferString(``)
	if err = json.NewEncoder(buf).Encode(msg); err != nil {
		fmt.Fprintf(os.Stderr, "Error encode json: %s", err.Error())
		return
	}
	if out, err = os.OpenFile(fsw.Filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.FileMode(0644)); err != nil {
		err = fmt.Errorf("Failed to open file '%s': %s", fsw.Filename, err.Error())
		return
	}
	defer out.Close()
	n, err = out.Write(buf.Bytes())
	return
}
