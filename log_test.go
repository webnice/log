package log

import (
	"os"
	"testing"

	b "github.com/webdeskltd/log/backends"
	l "github.com/webdeskltd/log/level"
	m "github.com/webdeskltd/log/message"
	tr "github.com/webdeskltd/log/trace"
	w "github.com/webdeskltd/log/writer"
)

func TestLogInitialize(t *testing.T) {
	testing_mode_two = true
	singleton[default_LOGUUID].getEssence().InterceptStandardLog(false)
	singleton[default_LOGUUID].getEssence().backend = b.NewBackends()
	testing_mode_two = false

	testing_mode_one = true
	hn, err := os.Hostname()
	if err != nil {
		t.Errorf("Error os.Hostname()")
		return
	}

	lg := NewLog().getEssence()
	if lg.HostName != `undefined` {
		t.Errorf("Error Initialize() os.Hostname() doesn't work")
		return
	}
	if lg.HostName == hn {
		t.Errorf("Error testing_mode_one flag")
		return
	}
	testing_mode_one = false

	lg = NewLog().getEssence()
	lg.defaultLevelLogWriter = w.NewWriter(l.NOTICE)
	lg.InterceptStandardLog(true)
	lg.Initialize()
	if lg.ready == false {
		t.Errorf("Error Initialize()")
		return
	}
	lg.InterceptStandardLog(false)
	err = lg.Close()
	if err != nil {
		t.Errorf("Error Close() logging: %v", err)
		return
	}
}

func TestLogSetModuleName(t *testing.T) {
	testing_mode_two = true
	singleton[default_LOGUUID].getEssence().InterceptStandardLog(false)
	singleton[default_LOGUUID].getEssence().backend = b.NewBackends()

	lg := NewLog().getEssence()
	lg.backend = b.NewBackends()
	if len(lg.moduleNames) != 0 {
		t.Errorf("Error in moduleNames (map[string]string)")
		return
	}
	lg.SetModuleName("TestLogSetModuleName")
	if len(lg.moduleNames) != 1 {
		t.Errorf("Error in SetModuleName()")
		return
	}
	lg.Notice("Test SetModuleName()")
	if nm, ok := lg.moduleNames["testing"]; !ok || nm != "TestLogSetModuleName" {
		t.Errorf("Error in SetModuleName()")
		return
	}
	lg.Close()
	testing_mode_two = false
}

func TestLogDelModuleName(t *testing.T) {
	lg := NewLog().getEssence()
	lg.SetModuleName("TestLogDelModuleName")
	if nm, ok := lg.moduleNames["testing"]; !ok || nm != "TestLogDelModuleName" {
		t.Errorf("Error in SetModuleName()")
		return
	}
	lg.DelModuleName()
	if len(lg.moduleNames) != 0 {
		t.Errorf("Error in DelModuleName()")
		return
	}
}

func TestLogResolveNames(t *testing.T) {
	testing_mode_two = true
	singleton[default_LOGUUID].getEssence().InterceptStandardLog(false)
	singleton[default_LOGUUID].getEssence().backend = b.NewBackends()

	lg := NewLog().getEssence()
	lg.backend = b.NewBackends()
	lg.moduleNames["testing"] = "TestLogResolveNames"
	txt, _ := m.NewMessage(
		tr.NewTrace().
			Trace(tr.STEP_BACK).
			GetRecord().
			Resolver(lg.ResolveNames),
	).Level(l.NOTICE).
		Write("Test SetModuleName()").
		Prepare().
		Record.
		Format(`%{module} | %{package}`)
	if txt != "testing | TestLogResolveNames" {
		t.Errorf("Error in TestLogResolveNames()")
		return
	}
	testing_mode_two = false
}
