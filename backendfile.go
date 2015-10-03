package log

import (
	"os"

	"github.com/webdeskltd/log/logging"
)

type fileBackend struct {
	File *os.File
}

func newFileBackend(file *os.File) (this *fileBackend) {
	this = new(fileBackend)
	this.File = file
	return
}

func (filebackend *fileBackend) Log(level logging.Level, calldepth int, record *logging.Record) (err error) {
	var line string
	line = record.Formatted(calldepth + 1)
	_, err = filebackend.File.WriteString(line + "\n")
	if err != nil {
		return
	}
	err = filebackend.File.Sync()
	return
}
