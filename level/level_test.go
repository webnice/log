package level

import (
	"testing"
)

func TestNew(t *testing.T) {
	var l *LevelObject
	l = New(DEBUG)
	if l == nil {
		t.Errorf("Error in New()")
	}
}

func TestLevelObjectString(t *testing.T) {
	var l *LevelObject

	l = New(FATAL)
	if l.String() != "FATAL" {
		t.Errorf("Error in String(FATAL)")
	}

	l = New(ALERT)
	if l.String() != "ALERT" {
		t.Errorf("Error in String(ALERT)")
	}

	l = New(CRITICAL)
	if l.String() != "CRITICAL" {
		t.Errorf("Error in String(CRITICAL)")
	}

	l = New(ERROR)
	if l.String() != "ERROR" {
		t.Errorf("Error in String(ERROR)")
	}

	l = New(WARNING)
	if l.String() != "WARNING" {
		t.Errorf("Error in String(WARNING)")
	}

	l = New(NOTICE)
	if l.String() != "NOTICE" {
		t.Errorf("Error in String(NOTICE)")
	}

	l = New(INFO)
	if l.String() != "INFO" {
		t.Errorf("Error in String(INFO)")
	}

	l = New(DEBUG)
	if l.String() != "DEBUG" {
		t.Errorf("Error in String(DEBUG)")
	}
}

func TestLevelObjectInt8(t *testing.T) {
	var l *LevelObject
	
	l = New(FATAL)
	if l.Int8() != int8(0) {
		t.Errorf("Error in Int8(FATAL)")
	}

	l = New(ALERT)
	if l.Int8() != int8(1) {
		t.Errorf("Error in Int8(ALERT)")
	}

	l = New(CRITICAL)
	if l.Int8() != int8(2) {
		t.Errorf("Error in Int8(CRITICAL)")
	}

	l = New(ERROR)
	if l.Int8() != int8(3) {
		t.Errorf("Error in Int8(ERROR)")
	}

	l = New(WARNING)
	if l.Int8() != int8(4) {
		t.Errorf("Error in Int8(WARNING)")
	}

	l = New(NOTICE)
	if l.Int8() != int8(5) {
		t.Errorf("Error in Int8(NOTICE)")
	}

	l = New(INFO)
	if l.Int8() != int8(6) {
		t.Errorf("Error in Int8(INFO)")
	}

	l = New(DEBUG)
	if l.Int8() != int8(7) {
		t.Errorf("Error in Int8(DEBUG)")
	}
}
