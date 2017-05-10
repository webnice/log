package gelf

import (
	"time"

	l "gopkg.in/webnice/log.v2/level"
)

const (
	GELF_VERSION = "1.2"
)

type Message struct {
	Version      string  `json:"version"`
	Host         string  `json:"host"`
	ShortMessage string  `json:"short_message"`
	Timestamp    float64 `json:"timestamp"`
	Level        l.Level `json:"level"`
	FullMessage  string  `json:"full_message,omitempty"`
	Facility     string  `json:"facility,omitempty"`
	Line         uint    `json:"line,omitempty"`
	File         string  `json:"file,omitempty"`
}

func NewMessage(source string, level l.Level, shortMessage string) (message *Message) {
	var tm = time.Now()
	timestamp := float64(tm.Unix()) + float64(time.Second)/float64(tm.Nanosecond())
	message = &Message{
		Version:      GELF_VERSION,
		Host:         source,
		ShortMessage: shortMessage,
		Timestamp:    timestamp,
		Level:        level,
	}
	return
}
