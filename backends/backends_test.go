package backends

import (
	"strings"
	"testing"

	m "github.com/webdeskltd/log/message"
	u "github.com/webdeskltd/log/uuid"
)

func NewBackendTesting() (ret *Backend) {
	ret = new(Backend)
	ret.hType = BACKEND_CONSOLE
	ret.reader = ret.readerTesting
	return
}

func (self *Backend) readerTesting(msg *m.Message) {
}

func TestNewBackends(t *testing.T) {
	var b *Backends = NewBackends()
	if b == nil {
		t.Errorf("Error in Backends(): nil")
	}
	b.Close()
}

func TestDestructor(t *testing.T) {
	b := NewBackends()
	b.Close()
	destructor(b)
	if b.Pool != nil {
		t.Errorf("Error in destructor()")
	}
}

func TestBackendsAddBackend(t *testing.T) {
	var bck *Backend = NewBackendTesting()
	var obj *Backends = NewBackends()
	var uid1, uid2 *u.UUID
	var err error

	uid1, err = obj.AddBackend(bck)
	if err != nil {
		t.Errorf("Error in AddBackend()")
	}
	uid2, err = obj.AddBackend(nil)
	if err == nil {
		t.Errorf("Error in AddBackend(nil)")
	} else if strings.Index(err.Error(), ErrBackendIsNull.Error()) != 0 {
		t.Errorf("Error in AddBackend(): %v", err)
	}
	if uid1 == nil || uid2 == nil {
		t.Errorf("Error in AddBackend(), uuid wrong")
		return
	}
	if uid1.String() == uid2.String() {
		t.Errorf("Error in AddBackend(), uuid wrong")
	}
	obj.Close()
}
