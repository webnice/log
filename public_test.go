package log

import (
	"testing"

	b "github.com/webdeskltd/log/backends"
	l "github.com/webdeskltd/log/level"
)

func TestConfigure(t *testing.T) {
	singleton[default_LOGUUID].getEssence().InterceptStandardLog(false)
	singleton[default_LOGUUID].getEssence().backend = b.NewBackends()
	err := Configure(nil)
	if err != ERROR_CONFIGURATION_IS_NULL {
		t.Errorf("Error Configure()")
		return
	}
}

func TestFatal(t *testing.T) {
	singleton[default_LOGUUID].getEssence().InterceptStandardLog(false)
	singleton[default_LOGUUID].getEssence().backend = b.NewBackends()
	var code int = 0
	exit_func = func(c int) {
		code = c
	}
	Fatal()
	if code == 0 {
		t.Errorf("Error Fatal()")
		return
	}
}

func TestAlert(t *testing.T) {
	singleton[default_LOGUUID].getEssence().InterceptStandardLog(false)
	singleton[default_LOGUUID].getEssence().backend = b.NewBackends()
	Alert()
}

func TestCritical(t *testing.T) {
	singleton[default_LOGUUID].getEssence().InterceptStandardLog(false)
	singleton[default_LOGUUID].getEssence().backend = b.NewBackends()
	Critical()
}

func TestError(t *testing.T) {
	singleton[default_LOGUUID].getEssence().InterceptStandardLog(false)
	singleton[default_LOGUUID].getEssence().backend = b.NewBackends()
	Error()
}

func TestWarning(t *testing.T) {
	singleton[default_LOGUUID].getEssence().InterceptStandardLog(false)
	singleton[default_LOGUUID].getEssence().backend = b.NewBackends()
	Warning()
}

func TestNotice(t *testing.T) {
	singleton[default_LOGUUID].getEssence().InterceptStandardLog(false)
	singleton[default_LOGUUID].getEssence().backend = b.NewBackends()
	Notice()
}

func TestInfo(t *testing.T) {
	singleton[default_LOGUUID].getEssence().InterceptStandardLog(false)
	singleton[default_LOGUUID].getEssence().backend = b.NewBackends()
	Info()
}

func TestDebug(t *testing.T) {
	singleton[default_LOGUUID].getEssence().InterceptStandardLog(false)
	singleton[default_LOGUUID].getEssence().backend = b.NewBackends()
	Debug()
}

func TestMessage(t *testing.T) {
	singleton[default_LOGUUID].getEssence().InterceptStandardLog(false)
	singleton[default_LOGUUID].getEssence().backend = b.NewBackends()
	Message(l.NOTICE)
	var code int = 0
	exit_func = func(c int) {
		code = c
	}
	Message(l.FATAL)
	if code == 0 {
		t.Errorf("Error Fatal()")
		return
	}
}

func TestClose(t *testing.T) {
	singleton[default_LOGUUID].getEssence().InterceptStandardLog(false)
	singleton[default_LOGUUID].getEssence().backend = b.NewBackends()
	Close()
}

func TestGetDefaultLog(t *testing.T) {
	singleton[default_LOGUUID].getEssence().InterceptStandardLog(false)
	singleton[default_LOGUUID].getEssence().backend = b.NewBackends()
	lg := GetDefaultLog()
	if lg == nil {
		t.Errorf("Error GetDefaultLog()")
	}
}

func TestSetApplicationName(t *testing.T) {
	SetApplicationName(`21D1478D-8460-4BC6-A55E-769F2CD653C1`)
	if singleton[default_LOGUUID].getEssence().AppName != `21D1478D-8460-4BC6-A55E-769F2CD653C1` {
		t.Errorf("Error SetApplicationName()")
	}
}

func TestSetModuleName(t *testing.T) {
	SetModuleName(`242A9053-CA60-4733-A0BB-F446FF5AD124`)
	if singleton[default_LOGUUID].getEssence().moduleNames["github.com/webdeskltd/log"] != `242A9053-CA60-4733-A0BB-F446FF5AD124` {
		t.Errorf("Error SetModuleName()")
	}
}

func TestDelModuleName(t *testing.T) {
	SetModuleName(`242A9053-CA60-4733-A0BB-F446FF5AD124`)
	if singleton[default_LOGUUID].getEssence().moduleNames["github.com/webdeskltd/log"] != `242A9053-CA60-4733-A0BB-F446FF5AD124` {
		t.Errorf("Error SetModuleName()")
	} else {
		DelModuleName()
		if _, ok := singleton[default_LOGUUID].getEssence().moduleNames["github.com/webdeskltd/log"]; ok {
			t.Errorf("Error DelModuleName()")
		}
	}
}

func TestInterceptStandardLog(t *testing.T) {
	if singleton[default_LOGUUID].getEssence().interceptStandardLog == false {
		InterceptStandardLog(true)
		if singleton[default_LOGUUID].getEssence().interceptStandardLog != true {
			t.Errorf("Error SetApplicationName(true)")
			return
		}
	}
	if singleton[default_LOGUUID].getEssence().interceptStandardLog == true {
		InterceptStandardLog(false)
		if singleton[default_LOGUUID].getEssence().interceptStandardLog != false {
			t.Errorf("Error SetApplicationName(false)")
			return
		}
	}
}

func TestLogEssenceGetWriter(t *testing.T) {
	a := singleton[default_LOGUUID].getEssence().defaultLevelLogWriter
	b := singleton[default_LOGUUID].GetWriter()
	if a != b {
		t.Errorf("Error GetWriter()")
		return
	}
	b = GetWriter()
	if a != b {
		t.Errorf("Error GetWriter()")
		return
	}
}
