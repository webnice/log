package backends // import "github.com/webdeskltd/log/backends"

import (
	"container/list"
	"time"

	m "github.com/webdeskltd/log/message"

	"github.com/webdeskltd/debug"
)

func init() {
	debug.Nop()
}

// Процесс читает из канала сообщения
func (self *Backends) messageWorker() {
	var msg *m.Message
	var exit bool
	// Дочитать сообщения до конца только потом завершиться
	for exit == false || len(self.RecordsChan) != 0 {
		select {
		case msg = <-self.RecordsChan:
			self.shuffle(msg)
		case <-self.exitChan:
			exit = true
		default:
			time.Sleep(time.Second / 10)
		}
	}
	self.doneChan <- true
}

// Выбор бэкэндов из пула
// Отправка сообщений в зависимости от режима и уровня логирования
func (self *Backends) shuffle(msg *m.Message) *m.Message {
	var item *list.Element
	var bck *Backend
	var pool []*Backend
	var i int
	var ok bool

	// Отбираем backend логгеры подходящие для уровня сообщения
	for item = self.Pool.Front(); item != nil; item = item.Next() {
		bck = item.Value.(*Backend)
		// Выбор по режиму
		ok = false
		switch bck.hMode {
		case MODE_NORMAL:
			if bck.hLevelNormal >= msg.Record.Level {
				ok = true
			}
		case MODE_SELECT:
			for i = range bck.hLevelSelect {
				if bck.hLevelSelect[i] == msg.Record.Level {
					ok = true
				}
			}
		}
		if ok {
			pool = append(pool, bck)
		}
	}

	// Если есть backend обработчики, то готовим сообщение
	if len(pool) > 0 {
		msg.Prepare()
	}

	for i = range pool {
		if pool[i].reader != nil {
			pool[i].reader(msg)
		}
	}

	// Устанавливаем длинну записанного сообщения
	msg.SetResult(len(msg.Record.Message), nil)

	return msg
}
