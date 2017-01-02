package fswriter // import "github.com/webdeskltd/log/middleware/fswriter"

//import "github.com/webdeskltd/debug"
import (
	"fmt"
	"os"
	"strings"
)

// New Create new package implementation and return interface
func New(filename ...string) Interface {
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
func (fsw *impl) SetFilename(filename string) Interface {
	fsw.Lock()
	defer fsw.Unlock()
	fsw.Filename = filename
	return fsw
}

// SetFilemode Set filemode
func (fsw *impl) SetFilemode(filemode os.FileMode) Interface { fsw.Filemode = filemode; return fsw }

// Write Запись среза байт в файл
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

// Write Запись строки в файл
func (fsw *impl) WriteString(buf string) (n int, err error) { return fsw.Write([]byte(buf)) }
