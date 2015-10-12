package backends

import (
	"container/list"

	"github.com/webdeskltd/log/record"
	"github.com/webdeskltd/log/uuid"

	"github.com/webdeskltd/debug"
)

func init() {
	debug.Nop()
}

func NewBackends() (ret *Backends) {
	ret = new(Backends)
	ret.RecordsChan = make(chan *record.Record, lengthRecords)
	ret.ResetBackends()
	return
}

// Добавление в пул нового backend
func (self *Backends) AddBackend(item *Backend) (ret *uuid.UUID, err error) {
	var id uuid.UUID = uuid.TimeUUID()
	if item == nil {
		err = ErrBackendIsNull
		return
	}
	self.PoolIndex[id] = self.Pool.PushBack(item)
	ret = &id
	return
}

// Удаление из пула backend по его ID
func (self *Backends) DelBackend(id *uuid.UUID) (err error) {
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

// Очистка пула backend
func (self *Backends) ResetBackends() {
	var err error
	var elm *list.Element
	var item *Backend

	// First run
	if self.Pool == nil {
		self.Pool = list.New()
	}
	for elm = self.Pool.Front(); elm != nil; elm = elm.Next() {
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

	self.Pool = list.New()
	self.PoolIndex = make(map[uuid.UUID]*list.Element)
}

// Sending messages to registered backends
// Для того чтобы максимально быстро вернуть управление главной программе используется буфферезированный канал
func (self *Backends) Push(r *record.Record) {
	// Предотвращение блокировки из за переполнения
	if len(self.RecordsChan) == lengthRecords {
		// TODO Надо организовать варианты обработки ситуации переполнения
		// Варианты
		// 1. Выдавать ошибки в STDERR
		// 2. Падать в панику
		// 3. Текущий варинт, все новые сообщения пропускать
		return
	}
	self.RecordsChan <- r
	return
}
