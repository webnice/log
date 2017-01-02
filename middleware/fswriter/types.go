package fswriter // import "github.com/webdeskltd/log/middleware/fswriter"

//import "github.com/webdeskltd/debug"
import (
	"os"
	"sync"
)

// Interface is an interface of package
type Interface interface {
	// SetFilename Set filename
	SetFilename(string) Interface

	// SetFilemode Set filemode
	SetFilemode(os.FileMode) Interface

	// Write Запись среза байт в файл
	Write([]byte) (int, error)

	// Write Запись строки в файл
	WriteString(string) (int, error)
}

// impl is an implementation of package
type impl struct {
	Filename string      // Имя файла в который записываются данные
	Filemode os.FileMode // filemode файла
	sync.RWMutex
}
