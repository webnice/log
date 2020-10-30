package fswformatjson // import "github.com/webnice/log/middleware/fswformatjson"

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
