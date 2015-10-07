package record

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/webdeskltd/log/uuid"
	"github.com/webdeskltd/debug"
)

const (
	TestTpl1 string = `%... %{one} : %{two:7s} ! %{three:8d} ++ %{four}%{five} [%{six:2006-01-02T15:04:05.000Z07:00}] {%{seven} = %{eight:16d}}`
	TestTpl2 string = `%... %{one} : %{two:7s} ! %{int:8d} ++ %{four}%{five} [%{six:2006-01-02T15:04:05.000Z07:00}] {%{seven} = %{eight:16d}}`
)

type TestDictionary1 struct {
	a bool         `fmt:"one"`           // %{one}
	B string       `fmt:"two:s"`         // %{two}
	C int          `fmt:"three:d"`       // %{three}
	D []byte       `fmt:"four:v"`        // %{four}
	E bytes.Buffer `fmt:"five"`          // %{five}
	F time.Time    `fmt:"six:t"`         // %{six}
	G int64        `fmt:"seven,eight:d"` // ${seven} alias %{eight}
	H string       `fmt:"-"`             // skeep
	I string       ``                    // skeep
}

type TestDictionary2 struct {
	Bad string `fmt:"bad:v:"`
}

func init() {
	//debug.Nop()
}

func TestMakeDictionary(t *testing.T) {
	var err error
	var ok bool

	err = makeDictionary(new(TestDictionary1))
	if err != nil {
		t.Errorf("Error in makeDictionary(): %v", err)
		return
	}

	// if exist
	if _, ok = templateNames["one"]; ok == false {
		t.Errorf("Incorrect work makeDictionary(), '%{one}' not found")
		return
	}
	if _, ok = templateNames["two"]; ok == false {
		t.Errorf("Incorrect work makeDictionary(), '%{two}' not found")
		return
	}
	if _, ok = templateNames["three"]; ok == false {
		t.Errorf("Incorrect work makeDictionary(), '%{three}' not found")
		return
	}
	if _, ok = templateNames["four"]; ok == false {
		t.Errorf("Incorrect work makeDictionary(), '%{four}' not found")
		return
	}
	if _, ok = templateNames["five"]; ok == false {
		t.Errorf("Incorrect work makeDictionary(), '%{five}' not found")
		return
	}
	if _, ok = templateNames["six"]; ok == false {
		t.Errorf("Incorrect work makeDictionary(), '%{six}' not found")
		return
	}
	if _, ok = templateNames["seven"]; ok == false {
		t.Errorf("Incorrect work makeDictionary(), '%{seven}' not found")
		return
	}
	if _, ok = templateNames["eight"]; ok == false {
		t.Errorf("Incorrect work makeDictionary(), '%{eight}' not found")
		return
	}

	// Default format
	if templateNames["one"].Format != "%v" {
		t.Errorf("Error in makeDictionary(). %%{one} default format is '%s' expected '%%v'", templateNames["one"].Format)
	}
	if templateNames["five"].Format != "%v" {
		t.Errorf("Error in makeDictionary(). %%{five} default format is '%s' expected '%%v'", templateNames["five"].Format)
	}
	if templateNames["seven"].Format != "%v" {
		t.Errorf("Error in makeDictionary(). %%{seven} default format is '%s' expected '%%v'", templateNames["seven"].Format)
	}

	// Format parse
	if templateNames["two"].Format != "%s" {
		t.Errorf("Error in makeDictionary(). %%{seven} format is '%s' expected '%%s'", templateNames["two"].Format)
	}
	if templateNames["three"].Format != "%d" {
		t.Errorf("Error in makeDictionary(). %%{seven} format is '%s' expected '%%d'", templateNames["three"].Format)
	}
	if templateNames["six"].Format != "%t" {
		t.Errorf("Error in makeDictionary(). %%{seven} format is '%s' expected '%%t'", templateNames["six"].Format)
	}
	if templateNames["seven"].Format != "%v" || templateNames["eight"].Format != "%d" {
		t.Errorf("Error in makeDictionary(). %%{seven} format is '%s' expected '%%v' and %%{eight} format is '%s' expected '%%d'",
			templateNames["seven"].Format, templateNames["eight"].Format)
	}

	// Synonym parse
	if templateNames["seven"].Name != templateNames["eight"].Name || templateNames["eight"].Name != "G" {
		t.Errorf("Error in makeDictionary(). %%{seven} alias %%{eight} is wrong")
	}

	// Bad syntax
	err = makeDictionary(new(TestDictionary2))
	if err == nil {
		t.Errorf("Error in makeDictionary(). Incorrect tag parse")
	}
	if strings.Index(err.Error(), errWrongTag.Error()) != 0 {
		t.Errorf("Error in makeDictionary(). Incorrect tag parse")
	}
}

func TestCheckFormat(t *testing.T) {
	var err error
	var matches [][]int

	err = makeDictionary(new(TestDictionary1))
	if err != nil {
		t.Errorf("Error in makeDictionary(): %v", err)
		return
	}

	matches, err = CheckFormat(TestTpl1)
	if len(matches) != 8 {
		t.Errorf("Error in CheckFormat(). Expected len(matches[][]) = [8][6]")
		return
	}

	matches, err = CheckFormat(TestTpl2)
	if err == nil {
		t.Errorf("Error in CheckFormat(). Expected error.Error('%s:int')", errUnknownVariable)
		return
	}
	if strings.Index(err.Error(), errUnknownVariable.Error()) != 0 {
		t.Errorf("Error in CheckFormat(). Expected error.Error('%s:int')", errUnknownVariable)
		return
	}

	matches, err = CheckFormat(``)
	if err == nil {
		t.Errorf("Error in CheckFormat(). Expected error.Error('%s')", errInvalidFormat)
		return
	}
	if strings.Index(err.Error(), errInvalidFormat.Error()) != 0 {
		t.Errorf("Error in CheckFormat(). Error is error.Error('%s') expected error.Error('%s')", err, errInvalidFormat)
	}
}

func TestRecordGetFormatedElement(t *testing.T) {
	var obj *Record = new(Record)
	var resp string
	var err error
	var matches [][]int

	makeDictionary(obj)
	obj.Id = uuid.TimeUUID()

	matches, err = CheckFormat(TestTpl1)
	resp = obj.getFormatedElement(templateNames[`id`], ``)

	debug.Dumper(matches, err, resp)
	
}
