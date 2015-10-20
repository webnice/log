package level

import ()

const (
	FATAL    Level = iota // 0 Fatal: system is unusable
	ALERT                 // 1 Alert: action must be taken immediately
	CRITICAL              // 2 Critical: critical conditions
	ERROR                 // 3 Error: error conditions
	WARNING               // 4 Warning: warning conditions
	NOTICE                // 5 Notice: normal but significant condition
	INFO                  // 6 Informational: informational messages
	DEBUG                 // 7 Debug: debug-level messages
)

var Map map[Level]string = map[Level]string{
	FATAL:    `FATAL`,    // Система не стабильна, проолжение работы не возможно
	ALERT:    `ALERT`,    // Система не стабильна но может частично продолжить работу (например запусился один из двух серверов - что-то работает а что-то нет)
	CRITICAL: `CRITICAL`, // Критическая ошибка, часть функционала системы работает не корректно
	ERROR:    `ERROR`,    // Ошибки не прерывающие работу приложения
	WARNING:  `WARNING`,  // Предупреждения
	NOTICE:   `NOTICE`,   // Информационные сообщения
	INFO:     `INFO`,     // Сообщения информационного характера описывающие шаги выполнения алгоритмов приложения
	DEBUG:    `DEBUG`,    // Режим отладки, аналогичен INFO но с подробными данными и дампом переменных
}

var L LevelObject

type Level int8 // Тип уровня журналирования

type LevelObject struct {
	Level
}
