package log

import (
	"path"
	"regexp"
	"runtime"
	"time"

	"github.com/webdeskltd/log/gelf"
	"github.com/webdeskltd/log/logging"
)

var gelfLevelMap = map[logging.Level]gelf.Level{
	logging.FATAL:    gelf.LEVEL_FATAL,    // Система не стабильна, проолжение работы не возможно
	logging.ALERT:    gelf.LEVEL_ALERT,    // Система не стабильна но может частично продолжить работу (например запусился один из двух серверов - что-то работает а что-то нет)
	logging.CRITICAL: gelf.LEVEL_CRITICAL, // Критическая ошибка, часть функционала системы работает не корректно
	logging.ERROR:    gelf.LEVEL_ERROR,    // Ошибки не прерывающие работу приложения
	logging.WARNING:  gelf.LEVEL_WARNING,  // Предупреждения
	logging.NOTICE:   gelf.LEVEL_NOTICE,   // Информационные сообщения
	logging.INFO:     gelf.LEVEL_INFO,     // Сообщения информационного характера описывающие шаги выполнения алгоритмов приложения
	logging.DEBUG:    gelf.LEVEL_DEBUG,    // Режим отладки, аналогичен INFO но с подробными данными и дампом переменных
}

var fnNameRegexp = regexp.MustCompile("^([[:print:]]+?)\\.((?:\\(\\*?[\\pL_][\\pL_\\pNd]+\\)\\.)?[\\pL_][\\pL_\\pNd]+)(?:[·.]\\d+)?$")

type loggingMessage struct {
	*gelf.Message
	ApplicationName string `json:"_application_name"`
	Function        string `json:"_function"`
}

type gelfBackend struct {
	*gelf.GelfClient
	source      string
	application string
}

func newGelfBackend(gelfClient *gelf.GelfClient, source string, application string) (backend *gelfBackend) {
	backend = new(gelfBackend)
	backend.GelfClient = gelfClient
	backend.source = source
	backend.application = application
	return
}

func (gelfBackend *gelfBackend) createLoggingMessage(record *logging.Record, file, function string, line uint) *loggingMessage {
	message := loggingMessage{
		Message:         gelf.NewMessage(gelfBackend.source, gelfLevelMap[record.Level], record.Message()),
		ApplicationName: gelfBackend.application,
		Function:        function,
	}
	message.Timestamp = float64(record.Time.Unix()) + float64(time.Second)/float64(record.Time.Nanosecond())
	message.Facility = record.Module
	message.File = file
	message.Line = line
	return &message
}

func (gelfBackend *gelfBackend) Log(level logging.Level, calldepth int, record *logging.Record) (err error) {
	file, fnName, line := callerMeta(calldepth + 1)
	message := gelfBackend.createLoggingMessage(record, file, fnName, line)
	err = gelfBackend.SendMessage(message)
	if err != nil {
		Error("Log message sending error: %s", err)
	}
	return
}

func callerMeta(calldepth int) (file, fnName string, line uint) {
	var pc uintptr
	var filePath string
	var lineNum int
	var ok bool

	pc, filePath, lineNum, ok = runtime.Caller(calldepth + 1)
	if ok {
		fullFnName := runtime.FuncForPC(pc).Name()
		fnNameMatches := fnNameRegexp.FindStringSubmatch(fullFnName)
		if len(fnNameMatches) > 1 {
			fileName := path.Base(filePath)
			file = path.Join(fnNameMatches[1], fileName)
			fnName = fnNameMatches[2]
			line = uint(lineNum)
			return
		}
	}
	file = "<unknown>"
	fnName = "<unknown>"
	return
}
