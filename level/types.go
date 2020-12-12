package level

const (
	levelFatal    Level = iota // 0 Fatal: system is unusable
	levelAlert                 // 1 Alert: action must be taken immediately
	levelCritical              // 2 Critical: critical conditions
	levelError                 // 3 Error: error conditions
	levelWarning               // 4 Warning: warning conditions
	levelNotice                // 5 Notice: normal but significant condition
	levelInfo                  // 6 Informational: informational messages
	levelDebug                 // 7 Debug: debug-level messages
)

// Defailt log level
const Defailt = levelNotice

var levelMap map[Level]string = map[Level]string{
	levelFatal:    `FATAL`,    // Система не стабильна, проолжение работы не возможно
	levelAlert:    `ALERT`,    // Система не стабильна но может частично продолжить работу (например запусился один из двух серверов - что-то работает а что-то нет)
	levelCritical: `CRITICAL`, // Критическая ошибка, часть функционала системы работает не корректно
	levelError:    `ERROR`,    // Ошибки не прерывающие работу приложения
	levelWarning:  `WARNING`,  // Предупреждения
	levelNotice:   `NOTICE`,   // Информационные сообщения
	levelInfo:     `INFO`,     // Сообщения информационного характера описывающие шаги выполнения алгоритмов приложения
	levelDebug:    `DEBUG`,    // Режим отладки, аналогичен INFO но с подробными данными и дампом переменных
}

type Level int8 // Тип уровня журналирования

// Interface is an interface of package
type Interface interface {
	// Fatal - 0 Fatal: system is unusable
	Fatal() Level

	// Alert - 1 Alert: action must be taken immediately
	Alert() Level

	// Critical - 2 Critical: critical conditions
	Critical() Level

	// Error - 3 Error: error conditions
	Error() Level

	// Warning - 4 Warning: warning conditions
	Warning() Level

	// Notice - 5 Notice: normal but significant condition
	Notice() Level

	// Informational - 6 Informational: informational messages
	Informational() Level

	// Debug - 7 Debug: debug-level messages
	Debug() Level
}

// impl is an implementation of package
type impl struct {
}
