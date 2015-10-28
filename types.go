package log

import (
	"bufio"
	"errors"
	"os"

	b "github.com/webdeskltd/log/backends"
	g "github.com/webdeskltd/log/gelf"
	l "github.com/webdeskltd/log/level"
	w "github.com/webdeskltd/log/writer"
)

const (
	default_LOG    string  = `CB7D0E12-C1EC-49CB-A3DD-AD62DE7FB7D8`
	default_FORMAT string  = `%{color}[%{module:-10s}] %{time:2006-01-02T15:04:05.000Z07:00t} (%{level:-8s}): %{message} (%{package}) (%{function}:%{line}) (%{shortfile}:%{line}) (%{longfile})`
	default_LEVEL  l.Level = l.NOTICE
)

var (
	ERROR_LEVEL_UNKNOWN         = errors.New(`Unknown or not supported logging level`)                                     // В конфигурации указан не известный уровень логирования
	ERROR_CONFIGURATION_IS_NULL = errors.New(`The configuration does not initialized. Received nil instead of the object`) // The configuration does not initialized. Received nil instead of the object
	ERROR_UNKNOWN_MODE          = errors.New(`Unknown logging mode`)                                                       // Unknown logging mode
	ERROR_LOG_FILENAME_IS_EMPTY = errors.New(`Log file name not specified`)                                                // Log file name not specified
	ERROR_INIT_FILE_BACKEND     = errors.New(`Can't initiate file backend`)                                                // Can't initiate file backend
)

// Карта всех логгеров
var singleton map[string]*Log
var exit_func func(code int) = os.Exit // Function completed applications for Fatal. In test mode is replaced by Nop()... 100% code test coverage :)
var testing_mode_one bool              // set true from test programm
var testing_mode_two bool              // set true from test programm

type Log struct {
	ready                  bool              // =true - log ready to use
	rescueSTDOUT           *os.File          // Save original STDOUT
	rescueSTDERR           *os.File          // Save original STDERR
	BufferSize             int               // Size memory buffer for log messages
	BufferFlushImmediately bool              // Flush log buffer after call
	Writer                 *bufio.Writer     // Log writer
	AppName                string            // %{program} - Application name
	HostName               string            // %{hostname} - Server host name
	cnf                    *Configuration    // Current configuration
	backend                *b.Backends       // Backend workflow
	defaultLevelLogWriter  *w.Writer         // Writer for standard logging and etc...
	interceptStandardLog   bool              // Current state for standard logging
	moduleNames            map[string]string // Кастомные названия модулей опубликованные через SetModuleName()
}

// Graylog server configuration
type ConfigurationGraylog2 struct {
	Host        string            `yaml:"Host"`        // IP адрес или имя хоста Graylog сервера
	Port        uint16            `yaml:"Port"`        // Порт на котором находится Graylog сервер
	Protocol    string            `yaml:"Protocol"`    // Протокол передачи данных, возможные значения: tcp, udp. По умолчанию: udp
	Source      string            `yaml:"Source"`      // Наименование источника логов
	ChunkSize   uint              `yaml:"ChunkSize"`   // Максимальный размер отправляемого пакета
	Compression g.CompressionType `yaml:"Compression"` // Сжатие передаваемых пакетов данных
	BufferSize  int64             `yaml:"BufferSize"`  // Размер буфера ???
}

// Telegram messenger configuration
type ConfigurationTelegram struct {
}

// Log configuration
type Configuration struct {
	BufferFlushImmediately bool                            `yaml:"BufferFlushImmediately"` // Сбрасывать буффер памяти сразу после записи строки (default: true - unbuffered)
	BufferSize             int                             `yaml:"BufferSize"`             // Размер буфера памяти в байтах (default: 0 - equal unbuffered)
	Mode                   map[b.BackendName][]l.LevelName `yaml:"Mode"`                   // Режим логирования, перечисляются включенные режимы логирования
	Levels                 map[b.BackendName]l.LevelName   `yaml:"Levels"`                 // Уровень логирования для каждого режима логирования
	Formats                map[b.BackendName]string        `yaml:"Formats"`                // Формат строки лога для каждого из режимов. Еслине описан, то берётся Format
	Format                 string                          `yaml:"Format"`                 // Формат по умолчанию для всех режимов
	File                   string                          `yaml:"File"`                   // Режим вывода в файл, путь и имя файла лога
	Graylog2               ConfigurationGraylog2           `yaml:"Graylog2"`               // Настройки подключения к graylog серверу
	Telegram               ConfigurationTelegram           `yaml:"Telegram"`               // Настройка отправки сообщений в telegram
}
