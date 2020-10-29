package sender // import "github.com/webnice/log/v2/sender"

import (
	"container/list"

	l "github.com/webnice/log/v2/level"
)

func init() {
	singleton = newSender()
	newReceiver()
}

// newSender Create new message object
func newSender() *impl {
	var snd = &impl{
		input:     make(chan Message, 100000),
		cancel:    make(chan interface{}),
		receivers: list.New(),
	}
	return snd
}

func newReceiver() {
	singleton.doCancelDone.Add(1)
	go func() {
		defer singleton.doCancelDone.Done()
		singleton.Receiver()
	}()
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
func (snd *impl) RemoveAllSender() { snd.receivers.Init() }

// Receiver Горутина получающая и обрабатывающая сообщения
func (snd *impl) Receiver() {
	var (
		msg  Message
		rec  *list.Element
		exit bool
	)

	for {
		if len(snd.input) == 0 && exit {
			return
		}
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
			exit = true
		}
	}
}

// Ожидание обработки всех буфферизированных сообщений
// Перезапуск Receiver и выход
func (snd *impl) Flush() {
	snd.cancel <- true
	snd.doCancelDone.Wait()
	newReceiver()
}
