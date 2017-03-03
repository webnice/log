package fswformatjson

//import "gopkg.in/webnice/debug.v1"
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
