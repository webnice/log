package sender // import "github.com/webdeskltd/log/sender"

//import "github.com/webdeskltd/debug"
import (
	"container/list"

	l "github.com/webdeskltd/log/level"
)

func init() {
	singleton = newSender()
	go singleton.Receiver()
}

// newSender Create new message object
func newSender() *impl {
	var snd = new(impl)
	snd.receivers = list.New()
	snd.input = make(chan Message, 100000)
	snd.cancel = make(chan interface{})
	return snd
}

// Gist Sender singleton object interface
func Gist() Interface {
	return singleton
}

// Channel Канал приёма сообщений
func (snd *impl) Channel() chan Message { return snd.input }

// SetDefaultReceiver Определение функции обработки сообщений по умолчанию
func (snd *impl) SetDefaultReceiver(fn Receiver) { snd.defaultReceiver = fn }

// AddSender Добавление нового отправителя сообщений
func (snd *impl) AddSender(fn Receiver) {
	if fn == nil {
		return
	}
	snd.receivers.PushBack(fn)
}

// Удаление всех отправителей сообщений, переключение на дефолтовый отправитель
func (snd *impl) RemoveAllSender() {
	snd.receivers.Init()
}

// Receiver Горутина получающая и обрабатывающая сообщения
func (snd *impl) Receiver() {
	var msg Message
	var rec *list.Element
	for {
		select {
		case msg = <-snd.input:
			if snd.receivers.Len() > 0 {
				for rec = snd.receivers.Front(); rec != nil; rec = rec.Next() {
					rec.Value.(Receiver)(msg)
				}

			} else if snd.defaultReceiver != nil {
				snd.defaultReceiver(msg)
			}
			// Call fatal function
			if msg.Level == l.New().Fatal() {
				fatalFn(msg.Level.Int())
			}
		case <-snd.cancel:
			return
		}
	}
}
