package message // import "github.com/webdeskltd/log/message"

import (
	"runtime"

	l "github.com/webdeskltd/log/level"
	r "github.com/webdeskltd/log/record"
)

type Message struct {
	Record   *r.Record
	WriteLen int
	WriteErr error
	level    l.Level
	written  chan bool
}

// Message object destructor
func destructor(obj *Message) {
	obj.Record = nil
	close(obj.written)
}

// Конструктор сообщений журнала
func NewMessage(record *r.Record) (this *Message) {
	this = new(Message)
	this.Record = record
	this.written = make(chan bool, 1000)
	runtime.SetFinalizer(this, destructor)
	return
}

// Устанавливаем уровень логирования для сообщения
func (self *Message) Level(level l.Level) *Message {
	self.level = level
	self.Record.Level = level
	return self
}

// Записываем сообщение в объект resord
// Пока без какой-либо обработки чтобы ускорить возврат к основной программе
// Сюда попадают все сообщения от всех уровней логирования
func (self *Message) Write(args ...interface{}) *Message {
	self.Record.SetMessage(args...)
	return self
}

// Подготовка сообщения
// Выполняется подготовка сообщения перед форматированием и выводом
// На данном этапе из переданных ранее args формируется единое текстовое сообщение
func (self *Message) Prepare() *Message {
	self.Record.Prepare()
	return self
}

// Вызывается после окончания обработки сообщения
func (self *Message) SetResult(i int, err error) *Message {
	self.WriteLen, self.WriteErr = i, err
	self.written <- true
	return self
}

// Return result after write message
func (self *Message) GetResult() (int, error) {
	<-self.written
	return self.WriteLen, self.WriteErr
}
