package backends

import (
	"container/list"
	"time"

	//l "github.com/webdeskltd/log/level"
	m "github.com/webdeskltd/log/message"
	//r "github.com/webdeskltd/log/record"

	"github.com/webdeskltd/debug"
)

func init() {
	debug.Nop()
}

// Процесс читает из канала сообщения
func (self *Backends) messageReader() {
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
	var txt string
	var i int
	var ok bool
	var err error

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
		// Форматируем сообщение
		txt, err = msg.Record.Format(pool[i].format)
		// Ошибка не должна никогда возникать так как формат проверяется при конфигуриговании
		// Но лучше перебдеть и проинформировать, чем недобдеть
		if err != nil {
			if LogError != nil {
				LogError("Error Record.Format(): %v", err)
			}
		}
		txt = txt
		print(pool[i].format); print("\n")
		print(txt); print("\n")

		//		print(self.Pool.Len())
		//		print(" ")
		//		print(MapTypeName[pool[i].hType])
		//		print(" ")
		//		print(l.Map[msg.Record.Level])
		//		print(" '")
		//		print(msg.Record.Message)
		//		print("'\n")
		//debug.Dumper(MapTypeName[pool[i].hType], modeName[pool[i].hMode], l.Map[pool[i].hLevelNormal], pool[i].hLevelSelect)

	}

	// Устанавливаем длинну записанного сообщения
	msg.SetResult(len(msg.Record.Message), nil)

	return msg
}
