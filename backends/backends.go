package backends

import (
	"container/list"
	"runtime"

	m "github.com/webdeskltd/log/message"
	u "github.com/webdeskltd/log/uuid"

	"github.com/webdeskltd/debug"
)

func init() {
	debug.Nop()
}

// Destructor
// Graceful stop all backends in pool
// - Flush all data
// - Close files and network connections
// - Close channel
func destructor(obj *Backends) {
	var err error
	var elm *list.Element
	var item *Backend
	if obj.Pool != nil {
		for elm = obj.Pool.Front(); elm != nil; elm = elm.Next() {
			item = elm.Value.(*Backend)
			if item == nil {
				continue
			}
			err = item.Stop()
			// TODO
			// Обработать ошибку
			// В настоящий момент ошибка игнорируется
			err = err
		}
	}
	if obj.Pool != nil {
		obj.Pool = nil
		close(obj.RecordsChan)
		close(obj.exitChan)
		close(obj.doneChan)
	}
}

func NewBackends() (ret *Backends) {
	ret = new(Backends)
	ret.Pool = list.New()
	ret.PoolIndex = make(map[u.UUID]*list.Element)
	ret.RecordsChan = make(chan *m.Message, lengthRecords)
	ret.exitChan = make(chan bool)
	ret.doneChan = make(chan bool)

	// Log message reader
	go ret.messageWorker()

	// Destructor
	runtime.SetFinalizer(ret, destructor)
	return
}

// Добавление в пул нового backend
func (self *Backends) AddBackend(item *Backend) (ret *u.UUID, err error) {
	ret = new(u.UUID)
	*ret = u.TimeUUID()

	if item == nil {
		err = ErrBackendIsNull
		return
	}
	self.PoolIndex[*ret] = self.Pool.PushBack(item)
	return
}

// Удаление из пула backend по его ID
func (self *Backends) DelBackend(id *u.UUID) (err error) {
	var ok bool
	var elm *list.Element
	var item *Backend

	if id == nil {
		err = ErrBackendIdIsEmpty
		return
	}
	elm, ok = self.PoolIndex[*id]
	if ok == false {
		err = ErrBackendNotFound
		return
	} else {
		item = elm.Value.(*Backend)
		if item != nil {
			err = item.Stop()
			// TODO
			// Обработать ошибку
			// В настоящий момент ошибка игнорируется
			err = err
		}
	}
	self.Pool.Remove(elm)
	delete(self.PoolIndex, *id)
	return
}

// Sending messages to registered backends
// Для того чтобы максимально быстро вернуть управление главной программе используется буфферезированный канал
func (self *Backends) Push(msg *m.Message) {
	if self == nil {
		return
	}
	// Предотвращение блокировки из за переполнения
	if len(self.RecordsChan) == lengthRecords {
		// TODO Надо организовать варианты обработки ситуации переполнения
		// Варианты
		// 1. Выдавать ошибки в STDERR
		// 2. Падать в панику
		panic("Очередь переполнена")
		return
	}
	self.RecordsChan <- msg

	return
}

// Wait for exit backend goroutine
func (self *Backends) Close() {
	self.exitChan <- true
	<-self.doneChan
}
