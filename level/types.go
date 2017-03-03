package level

const (
	_FATAL    Level = iota // 0 Fatal: system is unusable
	_ALERT                 // 1 Alert: action must be taken immediately
	_CRITICAL              // 2 Critical: critical conditions
	_ERROR                 // 3 Error: error conditions
	_WARNING               // 4 Warning: warning conditions
	_NOTICE                // 5 Notice: normal but significant condition
	_INFO                  // 6 Informational: informational messages
	_DEBUG                 // 7 Debug: debug-level messages
)

// Defailt log level
const Defailt = _NOTICE

var levelMap map[Level]string = map[Level]string{
	_FATAL:    `FATAL`,    // Система не стабильна, проолжение работы не возможно
	_ALERT:    `ALERT`,    // Система не стабильна но может частично продолжить работу (например запусился один из двух серверов - что-то работает а что-то нет)
	_CRITICAL: `CRITICAL`, // Критическая ошибка, часть функционала системы работает не корректно
	_ERROR:    `ERROR`,    // Ошибки не прерывающие работу приложения
	_WARNING:  `WARNING`,  // Предупреждения
	_NOTICE:   `NOTICE`,   // Информационные сообщения
	_INFO:     `INFO`,     // Сообщения информационного характера описывающие шаги выполнения алгоритмов приложения
	_DEBUG:    `DEBUG`,    // Режим отладки, аналогичен INFO но с подробными данными и дампом переменных
}

type Level int8 // Тип уровня журналирования

// Interface is an interface of package
type Interface interface {
	Fatal() Level         // 0 Fatal: system is unusable
	Alert() Level         // 1 Alert: action must be taken immediately
	Critical() Level      // 2 Critical: critical conditions
	Error() Level         // 3 Error: error conditions
	Warning() Level       // 4 Warning: warning conditions
	Notice() Level        // 5 Notice: normal but significant condition
	Informational() Level // 6 Informational: informational messages
	Debug() Level         // 7 Debug: debug-level messages
}

// impl is an implementation of package
type impl struct {
}
