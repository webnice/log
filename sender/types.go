package sender // import "github.com/webnice/log/v2/sender"

import (
	"container/list"
	"os"
	"sync"

	l "github.com/webnice/log/v2/level"
	t "github.com/webnice/log/v2/trace"
)

var (
	singleton *impl
	fatalFn   func(code int) = os.Exit // Function fatal exit with error code
)

// Interface is an interface of package
type Interface interface {
	// Channel Канал приёма сообщений
	Channel() chan Message

	// SetDefaultReceiver Определение функции обработки сообщений по умолчанию
	SetDefaultReceiver(Receiver)

	// AddSender Добавление нового отправителя сообщений
	AddSender(Receiver)

	// Удаление всех отправителей сообщений, переключение на дефолтовый отправитель
	RemoveAllSender()

	// Flush Ожидание обработки всех буфферизированных сообщений
	// Перезапуск Receiver и выход
	Flush()
}

// Receiver Функция приёма и обработки сообщений
type Receiver func(Message)

// impl is an implementation of package
type impl struct {
	input           chan Message     // Буферизированный лог сообщений
	cancel          chan interface{} // Канал завершения работы Receiver()
	doCancelDone    sync.WaitGroup   // Ожидание завершения Receiver()
	defaultReceiver Receiver         // Функция обработки сообщений по умолчанию
	receivers       *list.List       // Список отправителей сообщений
}

// Message data structure
type Message struct {
	// Level Уровень сообщения
	Level l.Level `json:"level"`

	// Trace stack information
	Trace *t.Info `json:"trace"`

	// Шаблон сообщения
	Pattern string `json:"pattern"`

	// Аргументы шаблона сообщения
	Args []interface{} `json:"args"`

	// Ключи сообщения
	Keys map[string]interface{} `json:"keys"`
}
