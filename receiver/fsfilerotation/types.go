package fsfilerotation

import (
	"regexp"
	"sync"
	"time"

	"github.com/webnice/lv2/middleware"

	s "github.com/webnice/lv2/sender"
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

	// Write Запись среза байт
	Write([]byte) (int, error)

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

	// SetFsWriter Установка функции записи в файл с форматированием
	SetFsWriter(middleware.FsWriter) Interface

	// GetFilename Получение текущего имени файла журнала
	GetFilename() string

	// Ticker Внешний таймер для ротации лог файлов.
	// Если в лог файл пишется мало информации, но необходимо своевременно производить ротацию, вызывается Ticker() по таймеру
	Ticker()
}

// impl is an implementation of package
type impl struct {
	TplText          string              // Шаблон форматирования текста
	Timezone         *time.Location      // Таймзона для отображения даты и времени файлов журнала (лога)
	MaxAge           time.Duration       // Максимальный возраст файла журнала до его удаления/очистки
	RotationTime     time.Duration       // Промежутки времени ротации файлов журнала
	Path             string              // Путь к папке размещения файлов журнала
	Filename         string              // Шаблон имени файла журнала
	FilenamePattern  string              // Шаблон файловой системы
	FilenameCurrent  string              // Текущее имя файла журнала в ротации
	SymbolicLinkName string              // Имя симлинка ведущего на текущий лог файл в ротации
	UnlinkFn         UnlinkFn            // Функция удаления файлов журнала
	FsWriter         middleware.FsWriter // Интерфейс записи
	sync.RWMutex
}
