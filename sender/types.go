package sender // import "github.com/webdeskltd/log/sender"

//import "github.com/webdeskltd/debug"
import (
	"container/list"
	"os"
	"sync"

	l "github.com/webdeskltd/log/level"
	t "github.com/webdeskltd/log/trace"
)

var singleton *impl
var fatalFn func(code int) = os.Exit // Function fatal exit with error code

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
	Level l.Level
	// Trace stack information
	Trace *t.Info
	// Шаблон сообщения
	Pattern string
	// Аргументы шаблона сообщения
	Args []interface{}
	// Ключи сообщения
	Keys map[string]interface{}
}