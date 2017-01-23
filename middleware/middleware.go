package middleware // import "github.com/webdeskltd/log/middleware"

import (
	"os"

	s "github.com/webdeskltd/log/sender"
)

// Interface of filesystem writer
type FsWriter interface {
	// SetFilename Set filename
	SetFilename(string) FsWriter

	// SetFilemode Set filemode
	SetFilemode(os.FileMode) FsWriter

	// SetFormat Set template line formating
	SetFormat(string) FsWriter

	// WriteMessage Запись структуры данных с форматированием
	WriteMessage(s.Message) (int, error)

	// Write Запись среза байт
	Write([]byte) (int, error)
}
