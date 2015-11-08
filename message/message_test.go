package message

import (
	"errors"
	"testing"
	"time"

	l "github.com/webdeskltd/log/level"
	r "github.com/webdeskltd/log/record"
)

func TestNewMessage(t *testing.T) {
	var rc *r.Record = r.NewRecord()
	var m *Message = NewMessage(rc)
	if m == nil {
		t.Errorf("Error in NewMessage()")
	}
}

func TestDestructor(t *testing.T) {
	var m = NewMessage(r.NewRecord())
	destructor(m)
	if m.Record != nil {
		t.Errorf("Error in destructor()")
	}
}

func TestMessageLevel(t *testing.T) {
	m := NewMessage(r.NewRecord())
	m.Level(l.DEBUG)
	if m.level != l.DEBUG || m.Record.Level != l.DEBUG {
		t.Errorf("Error in Level()")
	}
	m.Level(l.CRITICAL)
	if m.level != l.CRITICAL || m.Record.Level != l.CRITICAL {
		t.Errorf("Error in Level()")
	}
}

func TestMessageWrite(t *testing.T) {
	m := NewMessage(r.NewRecord())
	m.Level(l.DEBUG).Write("|%s", "4CDA289F-4659-4CE7-825E-5BD766F8C808").Prepare()
	if m.Record.Message != "|4CDA289F-4659-4CE7-825E-5BD766F8C808" {
		t.Errorf("Error in Write()")
	}
}

func TestMessageSetResult(t *testing.T) {
	var i int
	var err, resp error

	err = errors.New("Test error")
	m := NewMessage(r.NewRecord())
	go func() {
		t.Logf("Sleep one second")
		time.Sleep(time.Second)
		m.SetResult(101, err)
		t.Logf("Sleep done.")
	}()

	i, resp = m.GetResult()
	if i != 101 || resp != err {
		t.Errorf("Error in SetResult() if GetResult()")
	}

}
