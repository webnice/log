package fsfilerotation

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/webnice/lv2/middleware"
	"github.com/webnice/lv2/middleware/fswformattext"

	s "github.com/webnice/lv2/sender"
)

// New Create new
func New() Interface {
	var rcv = &impl{
		TplText:      _DefaultTextFORMAT,
		Timezone:     time.Local,
		RotationTime: time.Hour * 24,
		FsWriter:     fswformattext.New(),
	}

	rcv.SetUnlinkFunc(os.Remove)
	rcv.SetFilenamePattern(rcv.defaultFilenamePattern())

	return rcv
}

// Имя файла журнала по умолчанию
func (rcv *impl) defaultFilenamePattern() (ret string) {
	var tmp []string

	if tmp = strings.Split(os.Args[0], string(os.PathSeparator)); len(tmp) > 0 {
		ret = tmp[len(tmp)-1] + `-%Y%m%d.log`
	} else {
		ret = os.Args[0] + `-%Y%m%d.log`
	}

	return
}

// SetPath Установка пути к папке фалов журнала
func (rcv *impl) SetPath(pth string) Interface { rcv.Path = pth; return rcv }

// SetFilenamePattern Установка шаблона имени файла журнала
func (rcv *impl) SetFilenamePattern(fnm string) Interface {
	var (
		tmp string
		rex *regexp.Regexp
	)

	tmp = fnm
	for _, rex = range patternConversion {
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

// SetMaxAge Установка максимального возраста файла журнала до его удаления
// По умолчанию =0 - файлы журналов не удаляются
func (rcv *impl) SetMaxAge(ma time.Duration) Interface { rcv.MaxAge = ma; return rcv }

// SetRotationTime Установка промежутков времени между ротацией файлов
// Значение по умолчанию одни сутки
func (rcv *impl) SetRotationTime(rt time.Duration) Interface { rcv.RotationTime = rt; return rcv }

// SetUnlinkFunc Установка пользовательской функции удаления файлов журнала
// Если приложению требуется не просто удалить файлы, а куда-то их отправить или заархивировать
// Вызывается для каждого файла лога отдельно
func (rcv *impl) SetUnlinkFunc(fn UnlinkFn) Interface { rcv.UnlinkFn = fn; return rcv }

// SetFsWriter Установка функции записи в файл с форматированием
func (rcv *impl) SetFsWriter(fn middleware.FsWriter) Interface { rcv.FsWriter = fn; return rcv }

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

// Rotation Ротация файлов
func (rcv *impl) Rotation(fnm string) (err error) {
	var (
		lockfn, tmpLinkName, pth string
		fh                       *os.File
		fi                       os.FileInfo
		guard                    fnGuard
		matches, toUnlink        []string
		cutoff                   time.Time
	)

	lockfn = fmt.Sprintf("%s_lock", fnm)
	if fh, err = os.OpenFile(lockfn, os.O_CREATE|os.O_EXCL, 0644); err != nil {
		return
	}
	guard.fn = func() {
		_ = fh.Close()
		_ = os.Remove(lockfn)
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
	// MaxAge is not set, not rotating
	if rcv.MaxAge <= 0 {
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

// Ticker Внешний таймер для ротации лог файлов.
// Смена имени файла журнала
func (rcv *impl) Ticker() {
	var (
		err   error
		fnm   string
		isNew bool
	)

	rcv.Lock()
	defer rcv.Unlock()
	if fnm, err = rcv.filename(); err != nil {
		return
	}
	if fnm == rcv.FilenameCurrent {
		return
	}
	isNew = true
	if _, err = os.Stat(fnm); err == nil {
		if rcv.SymbolicLinkName != "" {
			if _, err = os.Lstat(rcv.SymbolicLinkName); err == nil {
				isNew = false
			}
		}
	}
	if isNew {
		if err = rcv.Rotation(fnm); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to rotate: %s\n", err.Error())
		}
	}
	rcv.FilenameCurrent = fnm
}

// Receiver Message receiver
func (rcv *impl) Receiver(msg s.Message) {
	var err error

	rcv.Ticker()
	if _, err = rcv.FsWriter.
		SetFilename(rcv.FilenameCurrent).
		SetFormat(rcv.TplText).
		WriteMessage(msg); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	}
}

// Write Запись среза байт
func (rcv *impl) Write(buf []byte) (n int, err error) {
	rcv.Ticker()
	if n, err = rcv.
		FsWriter.
		SetFilename(rcv.FilenameCurrent).
		Write(buf); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	}

	return
}
