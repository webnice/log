package backends

import (
	"container/list"
	"time"

	l "github.com/webdeskltd/log/level"
	m "github.com/webdeskltd/log/message"

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
	var i int
	var ok bool

	//debug.Dumper(l.Map[msg.Record.Level], msg.Record.FileNameShort, msg.Record.FileLine)

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

	l.New(l.ALERT)
	for i = range pool {
		print(self.Pool.Len())
		print(" ")
		print(MapTypeName[pool[i].hType])
		print(" ")
		print(l.Map[msg.Record.Level])
		print("\n")
		//debug.Dumper(MapTypeName[pool[i].hType], modeName[pool[i].hMode], l.Map[pool[i].hLevelNormal], pool[i].hLevelSelect)

	}

	//debug.Dumper(msg)
	msg.SetResult(100, nil)

	return msg
}
