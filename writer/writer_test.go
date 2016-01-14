package writer

import (
	"testing"

	b "github.com/webdeskltd/log/backends"
	l "github.com/webdeskltd/log/level"
	r "github.com/webdeskltd/log/record"
)

func TestNewWriter(t *testing.T) {
	var obj = NewWriter(l.INFO)
	if obj == nil {
		t.Errorf("Error creating writer object")
		return
	}
	if obj.level != l.INFO {
		t.Errorf("Error initialize logging level in writer object")
		return
	}
	obj = nil
}

func TestWriterResolver(t *testing.T) {
	var obj = NewWriter(l.INFO)

	obj.Resolver(func(rec *r.Record) {}).Write([]byte{})
	obj.Resolver(func(rec *r.Record) {}).WriteString("")
	obj.Resolver(func(rec *r.Record) {}).Println("", "", "")
	if obj.resolver == nil {
		t.Errorf("Error Resolver() not work")
		return
	}
}

func TestWriterResolverAttachBackends(t *testing.T) {
	var bck = b.NewBackends()
	var obj = NewWriter(l.INFO)

	obj.AttachBackends(bck).Write([]byte{})
	obj.AttachBackends(bck).WriteString("")
	if obj.backends == nil {
		t.Errorf("Error AttachBackends() not work")
		return
	}
}

func TestDestructor(t *testing.T) {
	var bck = b.NewBackends()
	var obj = NewWriter(l.INFO)
	obj.Resolver(func(rec *r.Record) {})
	obj.AttachBackends(bck)
	if obj.backends == nil {
		t.Errorf("Error AttachBackends() not work")
		return
	}
	if obj.resolver == nil {
		t.Errorf("Error Resolver() not work")
		return
	}
	destructor(obj)
	if obj.backends != nil || obj.resolver != nil {
		t.Errorf("Error destructor() not work")
		return
	}
}

func TestWriterCleanSpace(t *testing.T) {
	var tmp string = " \t \n \f \rD123\t\n\f\r33333Edn\t\n\f\r \r\t"
	var obj = NewWriter(l.INFO)
	tmp = obj.cleanSpace(tmp)
	if tmp != "D123\t\n\f\r33333Edn" {
		t.Errorf("'%s'", tmp)
		t.Errorf("Error cleanSpace() not work")
		return
	}
}
