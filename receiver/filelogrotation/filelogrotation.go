package filelogrotation // import "github.com/webdeskltd/log/receiver/filelogrotation"

//import "github.com/webdeskltd/debug"
import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	f "github.com/webdeskltd/log/formater"
	s "github.com/webdeskltd/log/sender"
)

// const _DefaultTextFORMAT = `%{color}[%{module:-10s}] %{time:2006-01-02T15:04:05.000Z07:00t} (%{level:-8s}): %{message} (%{package}) (%{function}:%{line}) (%{shortfile}:%{line}) (%{longfile})`
const _DefaultTextFORMAT = `%{time:2006-01-02T15:04:05.000Z07:00t} (%{level:-8s}): %{message} {%{package}/%{shortfile}:%{line}, func: %{function}()}`

var patternConversion = []*regexp.Regexp{
	regexp.MustCompile(`%[%+A-Za-z]`),
	regexp.MustCompile(`\*+`),
}

// UnlinkFn Delete old log files function
type UnlinkFn func(string) error

// Interface is an interface of package
type Interface interface {
	// Receiver Message receiver
	Receiver(s.Message)

	// SetPath Установка шаблона файла журнала
	SetPath(string) Interface

	// SetFilenamePattern Установка шаблона файла журнала
	SetFilenamePattern(string) Interface

	// SetTimezone Установка таймзоны для отображения времени в лог файле
	SetTimezone(*time.Location) Interface

	// SetSymlink Установка имени симлинка ведущего на текущий лог файл в ротации (только для *nix OS)
	SetSymlink(string) Interface

	// SetMaxAge Установка максимального возраста файла журнала до его удаления или очистки.
	// По умолчанию =0 - файлы журналов не удаляются и не очищаются
	SetMaxAge(time.Duration) Interface

	// SetRotationTime Установка промежутков времени между ротацией файлов
	// Значение по умолчанию одни сутки
	SetRotationTime(time.Duration) Interface

	// SetUnlinkFunc Установка пользовательской функции удаления файлов журнала
	// Например если приложению требуется не просто удалить файлы а куда-то их отправить или заархивировать
	// Вызывается для каждого файла лога отдельно
	SetUnlinkFunc(UnlinkFn) Interface

	// GetFilename Получение текущего имени файла журнала
	GetFilename() string

	// Write Запись среза байт в файл журнала
	Write([]byte) (int, error)

	// Close Закрытие файлового дескриптора файла журнала
	Close() error
}

// impl is an implementation of package
type impl struct {
	Formater         f.Interface    // Formater interface
	TplText          string         // Шаблон форматирования текста
	Timezone         *time.Location // Таймзона для отображения даты и времени файлов журнала (лога)
	MaxAge           time.Duration  // Максимальный возраст файла журнала до его удаления/очистки
	RotationTime     time.Duration  // Промежутки времени ротации файлов журнала
	Path             string         // Путь к папке размещения файлов журнала
	Filename         string         // Шаблон имени файла журнала
	FilenamePattern  string         // Шаблон файловой системы
	FilenameCurrent  string         // Текущее имя файла журнала в ротации
	SymbolicLinkName string         // Имя симлинка ведущего на текущий лог файл в ротации
	OutFh            *os.File       // Хандлер открытого файла журнала
	UnlinkFn         UnlinkFn       // Функция удаления файлов журнала
	sync.RWMutex
}

// New Create new
func New() Interface {
	var rcv = new(impl)
	rcv.Formater = f.New()
	rcv.TplText = _DefaultTextFORMAT
	rcv.Timezone = time.Local
	rcv.RotationTime = time.Hour * 24
	rcv.SetFilenamePattern(os.Args[0] + `-%Y%m%d.log`)
	rcv.SetUnlinkFunc(os.Remove)
	return rcv
}

// SetPath Установка шаблона файла журнала
func (rcv *impl) SetPath(pth string) Interface { rcv.Path = pth; return rcv }

// SetFilenamePattern Установка шаблона файла журнала
func (rcv *impl) SetFilenamePattern(fnm string) Interface {
	var tmp = fnm
	for _, rex := range patternConversion {
		tmp = rex.ReplaceAllString(tmp, "*")
	}
	rcv.Filename = fnm
	rcv.FilenamePattern = tmp
	return rcv
}

// SetTimezone Установка таймзоны для отображения времени в лог файле
// По умолчанию time.Local
func (rcv *impl) SetTimezone(tz *time.Location) Interface { rcv.Timezone = time.Local; return rcv }

// SetSymlink Установка имени симлинка ведущего на текущий лог файл в ротации (только для *nix OS)
func (rcv *impl) SetSymlink(slnk string) Interface {
	rcv.SymbolicLinkName = rcv.absPath(path.Join(rcv.Path, slnk))
	return rcv
}

// SetMaxAge Установка максимального возраста файла журнала до его удаления или очистки.
// По умолчанию =0 - файлы журналов не удаляются и не очищаются
func (rcv *impl) SetMaxAge(ma time.Duration) Interface { rcv.MaxAge = ma; return rcv }

// SetRotationTime Установка промежутков времени между ротацией файлов
// Значение по умолчанию одни сутки
func (rcv *impl) SetRotationTime(rt time.Duration) Interface { rcv.RotationTime = rt; return rcv }

// SetUnlinkFunc Установка пользовательской функции удаления файлов журнала
// Например если приложению требуется не просто удалить файлы а куда-то их отправить или заархивировать
// Вызывается для каждого файла лога отдельно
func (rcv *impl) SetUnlinkFunc(fn UnlinkFn) Interface { rcv.UnlinkFn = fn; return rcv }

// GetFilename Получение текущего имени файла журнала
func (rcv *impl) GetFilename() string { rcv.Lock(); defer rcv.Unlock(); return rcv.FilenameCurrent }

// Получение абсолютного пути к файлу
func (rcv *impl) absPath(pth string) (ret string) {
	var err error
	if len(pth) > 0 {
		switch pth[0] {
		case '/':
			ret = pth
		default:
			if ret, err = os.Getwd(); err != nil {
				ret = pth
				return
			}
			ret = path.Join(ret, pth)
		}
	}
	return
}

// Генерация имени лог файла
func (rcv *impl) filename() (ret string, err error) {
	var tn, tf time.Time
	var diff time.Duration
	tn = time.Now().In(rcv.Timezone)
	diff = time.Duration(tn.UnixNano()) % rcv.RotationTime
	tf = tn.Add(diff * -1)
	if ret, err = Format(rcv.Filename, tf); err != nil {
		ret = ``
		return
	}
	ret = rcv.absPath(path.Join(rcv.Path, ret))
	return
}

// Write Is implementation io.Writer interface
func (rcv *impl) Write(buf []byte) (n int, err error) {
	var fnm string
	var out *os.File
	var isNew bool
	rcv.Lock()
	defer rcv.Unlock()
	if fnm, err = rcv.filename(); err != nil {
		return
	}
	if fnm == rcv.FilenameCurrent {
		out = rcv.OutFh
	}
	if out == nil {
		isNew = true
		if _, err = os.Stat(fnm); err == nil {
			if rcv.SymbolicLinkName != "" {
				if _, err = os.Lstat(rcv.SymbolicLinkName); err == nil {
					isNew = false
				}
			}
		}
		if out, err = os.OpenFile(fnm, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.FileMode(0644)); err != nil {
			err = fmt.Errorf("Failed to open file '%s': %s", rcv.Filename, err.Error())
			return
		}
		if isNew {
			if err = rcv.Rotation(fnm); err != nil {
				fmt.Fprintf(os.Stderr, "failed to rotate: %s\n", err.Error())
			}
		}
	}
	n, err = out.Write(buf)
	if rcv.OutFh == nil {
		rcv.OutFh = out
	} else if isNew {
		rcv.OutFh.Close()
		rcv.OutFh = out
	}
	rcv.FilenameCurrent = fnm
	return
}

// Rotation Ротация файлов
func (rcv *impl) Rotation(fnm string) (err error) {
	var lockfn, tmpLinkName, pth string
	var fh *os.File
	var fi os.FileInfo
	var guard fnGuard
	var matches, toUnlink []string
	var cutoff time.Time

	lockfn = fmt.Sprintf("%s_lock", fnm)
	if fh, err = os.OpenFile(lockfn, os.O_CREATE|os.O_EXCL, 0644); err != nil {
		return
	}
	guard.fn = func() {
		fh.Close()
		os.Remove(lockfn)
	}
	defer guard.Run()
	if rcv.SymbolicLinkName != "" {
		tmpLinkName = fmt.Sprintf("%s_symlink", fnm)
		if err = os.Symlink(fnm, tmpLinkName); err != nil {
			return
		}
		if err = os.Rename(tmpLinkName, rcv.SymbolicLinkName); err != nil {
			return
		}
	}
	if rcv.MaxAge <= 0 {
		err = fmt.Errorf("SetMaxAge is not set, not rotating")
		return
	}
	if matches, err = filepath.Glob(rcv.absPath(path.Join(rcv.Path, rcv.FilenamePattern))); err != nil {
		return
	}
	cutoff = time.Now().In(rcv.Timezone).Add(rcv.MaxAge * -1)
	for _, pth = range matches {
		// Ignore lock files
		if strings.HasSuffix(pth, "_lock") || strings.HasSuffix(pth, "_symlink") {
			continue
		}
		if fi, err = os.Stat(pth); err != nil {
			continue
		}
		if fi.ModTime().After(cutoff) {
			continue
		}
		toUnlink = append(toUnlink, pth)
	}
	if len(toUnlink) <= 0 {
		return
	}
	guard.Enable()
	go func() {
		for _, pth = range toUnlink {
			_ = rcv.UnlinkFn(pth)
		}
	}()
	return
}

// Close Is implementation io.Closer interface
func (rcv *impl) Close() (err error) {
	rcv.Lock()
	defer rcv.Unlock()
	if rcv.OutFh == nil {
		return
	}
	rcv.OutFh.Close()
	rcv.OutFh = nil
	return
}

// Receiver Message receiver. Output to STDERR
func (rcv *impl) Receiver(msg s.Message) {
	var buf *bytes.Buffer
	var err error
	if buf, err = rcv.Formater.Text(msg, rcv.TplText); err != nil {
		buf = bytes.NewBufferString(fmt.Sprintf("Error formatting log message: %s", err.Error()))
	}
	buf.WriteString("\r\n")
	if _, err = rcv.Write(buf.Bytes()); err != nil {
		fmt.Fprintf(os.Stderr, "Error write log message: %s\n", err.Error())
	}
}
