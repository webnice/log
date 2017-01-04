package fswformatjson // import "github.com/webdeskltd/log/middleware/fswformatjson"

//import "github.com/webdeskltd/debug"
import (
	"os"
	"sync"
)

// impl is an implementation of package
type impl struct {
	Filename string      // Имя файла в который записываются данные
	Filemode os.FileMode // filemode файла
	sync.RWMutex
}
