package fswformattext // import "github.com/webnice/log/v2/middleware/fswformattext"

import (
	"os"
	"sync"

	f "github.com/webnice/log/v2/formater"
)

// const _DefaultTextFORMAT = `%{color}[%{module:-10s}] %{time:2006-01-02T15:04:05.000Z07:00t} (%{level:-8s}): %{message} (%{package}) (%{function}:%{line}) (%{shortfile}:%{line}) (%{longfile})`
const _DefaultTextFORMAT = `(%{colorbeg}%{level:1s}:%{level:1d}%{colorend}): %{message} {%{package}/%{shortfile}:%{line}, func: %{function}()}`

// impl is an implementation of package
type impl struct {
	Filename string      // Имя файла в который записываются данные
	Filemode os.FileMode // filemode файла
	Formater f.Interface // Formater interface
	TplText  string      // Шаблон форматирования текста
	sync.RWMutex
}
