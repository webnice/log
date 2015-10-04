package gelf

import (
	"time"
)

const (
	GELF_VERSION = "1.2"
)

type Level int8

const (
	LEVEL_FATAL    Level = iota // 0 Система не стабильна, проолжение работы не возможно
	LEVEL_ALERT                 // 1 Система не стабильна но может частично продолжить работу (например запусился один из двух серверов - что-то работает а что-то нет)
	LEVEL_CRITICAL              // 2 Критическая ошибка, часть функционала системы работает не корректно
	LEVEL_ERROR                 // 3 Ошибки не прерывающие работу приложения
	LEVEL_WARNING               // 4 Предупреждения
	LEVEL_NOTICE                // 5 Информационные сообщения
	LEVEL_INFO                  // 6 Сообщения информационного характера описывающие шаги выполнения алгоритмов приложения
	LEVEL_DEBUG                 // 7 Режим отладки, аналогичен INFO но с подробными данными и дампом переменных
)

type Message struct {
	Version      string  `json:"version"`
	Host         string  `json:"host"`
	ShortMessage string  `json:"short_message"`
	Timestamp    float64 `json:"timestamp"`
	Level        Level   `json:"level"`
	FullMessage  string  `json:"full_message,omitempty"`
	Facility     string  `json:"facility,omitempty"`
	Line         uint    `json:"line,omitempty"`
	File         string  `json:"file,omitempty"`
}

func NewMessage(source string, level Level, shortMessage string) *Message {
	t := time.Now()
	timestamp := float64(t.Unix()) + float64(time.Second)/float64(t.Nanosecond())
	message := Message{
		Version:      GELF_VERSION,
		Host:         source,
		ShortMessage: shortMessage,
		Timestamp:    timestamp,
		Level:        level,
	}
	return &message
}
