package log

import (
	"bufio"
	"errors"
	"os"

	"github.com/webdeskltd/log/gelf"
	"github.com/webdeskltd/log/logging"
	"github.com/webdeskltd/log/record"
)

const (
	level_FATAL    logLevel = iota // 0 Fatal: system is unusable
	level_ALERT                    // 1 Alert: action must be taken immediately
	level_CRITICAL                 // 2 Critical: critical conditions
	level_ERROR                    // 3 Error: error conditions
	level_WARNING                  // 4 Warning: warning conditions
	level_NOTICE                   // 5 Notice: normal but significant condition
	level_INFO                     // 6 Informational: informational messages
	level_DEBUG                    // 7 Debug: debug-level messages
)

const (
	mode_CONSOLE ConfigurationModeName = "console"
	mode_SYSLOG  ConfigurationModeName = "syslog"
	mode_FILE    ConfigurationModeName = "file"
	mode_GRAYLOG ConfigurationModeName = "graylog"
)

const (
	defaultFormat string   = `"%{color}[%{module:-10s}] %{time:2006-01-02T15:04:05.000Z07:00} (%{level:7s}): %{message} (%{package}) (%{function}:%{line}) (%{shortfile}:%{line}) (%{longfile})"`
	defaultLevel  logLevel = level_NOTICE
)

var (
	self     *configuration      // Singleton
	levelMap map[logLevel]string = map[logLevel]string{
		level_FATAL:    `FATAL`,    // Система не стабильна, проолжение работы не возможно
		level_ALERT:    `ALERT`,    // Система не стабильна но может частично продолжить работу (например запусился один из двух серверов - что-то работает а что-то нет)
		level_CRITICAL: `CRITICAL`, // Критическая ошибка, часть функционала системы работает не корректно
		level_ERROR:    `ERROR`,    // Ошибки не прерывающие работу приложения
		level_WARNING:  `WARNING`,  // Предупреждения
		level_NOTICE:   `NOTICE`,   // Информационные сообщения
		level_INFO:     `INFO`,     // Сообщения информационного характера описывающие шаги выполнения алгоритмов приложения
		level_DEBUG:    `DEBUG`,    // Режим отладки, аналогичен INFO но с подробными данными и дампом переменных
	}
	errLevelUnknown = errors.New(`Unknown or not supported logging level`) // В конфигурации указан не известный уровень логирования
)

type logLevel int8

type configuration struct {
	BufferSize             int                    // Size memory buffer for log messages
	BufferFlushImmediately bool                   // Flush log buffer after call
	Writer                 *bufio.Writer          // Log writer
	AppName                string                 // %{program} - Application name
	HostName               string                 // %{hostname} - Server host name
	cnf                    *Configuration         // Current configuration
	fH                     *os.File               // File handle
	backends               []logging.Backend      // All mode logging
	bStderr                logging.LeveledBackend // Mode: "console" STDERR
	bSyslog                logging.LeveledBackend // Mode: "syslog" SYSTEM SYSLOG
	bFile                  logging.LeveledBackend // Mode: "file"
	bGraylog               logging.LeveledBackend // Mode: "graylog"
	moduleNames            map[string]string      // Кастомные названия модулей опубликованные через SetModuleName()
}

type log struct {
	Record   *record.Record
	WriteLen int
	WriteErr error
}

type writer struct {
}

type ConfigurationModeName string
type ConfigurationLevelName string

// Graylog server configuration
type ConfigurationGraylog struct {
	Host        string               `yaml:"Host"`        // IP адрес или имя хоста Graylog сервера
	Port        uint16               `yaml:"Port"`        // Порт на котором находится Graylog сервер
	Protocol    string               `yaml:"Protocol"`    // Протокол передачи данных, возможные значения: tcp, udp. По умолчанию: udp
	Source      string               `yaml:"Source"`      // Наименование источника логов
	ChunkSize   uint                 `yaml:"ChunkSize"`   // Максимальный размер отправляемого пакета
	Compression gelf.CompressionType `yaml:"Compression"` // Сжатие передаваемых пакетов данных
	BufferSize  int64                `yaml:"BufferSize"`  // Размер буфера ???
}

// Log configuration
type Configuration struct {
	BufferFlushImmediately bool                                             `yaml:"BufferFlushImmediately"` // Сбрасывать буффер памяти сразу после записи строки (default: true - unbuffered)
	BufferSize             int                                              `yaml:"BufferSize"`             // Размер буфера памяти в байтах (default: 0 - equal unbuffered)
	Mode                   []ConfigurationModeName                          `yaml:"Mode"`                   // Режим логирования, перечисляются включенные режимы логирования
	Levels                 map[ConfigurationModeName]ConfigurationLevelName `yaml:"Levels"`                 // Уровень логирования для каждого режима логирования
	Formats                map[ConfigurationModeName]string                 `yaml:"Formats"`                // Формат строки лога для каждого из режимов. Еслине описан, то берётся Format
	Format                 string                                           `yaml:"Format"`                 // Формат по умолчанию для всех режимов
	File                   string                                           `yaml:"File"`                   // Режим вывода в файл, путь и имя файла лога
	Graylog                ConfigurationGraylog                             `yaml:"Graylog"`                // Настройки подключения к graylog серверу
}
