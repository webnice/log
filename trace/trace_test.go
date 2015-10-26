package trace

import (
	"testing"
	
	r "github.com/webdeskltd/log/record"
)

func TestNewTrace(t *testing.T) {
	var trace = NewTrace()
	if trace == nil {
		t.Errorf("Error in NewTrace()")
		return
	}
}

func TestTraceTrace(t *testing.T) {
	var trace = NewTrace()
	var rec *r.Record
	
	if trace == nil {
		t.Errorf("Error in Trace()")
		return
	}
	trace.Trace(0)
	if trace.GetRecord().FileNameShort != `testing.go` {
		t.Errorf("Error in Trace().FileNameShort")
		return
	}

	trace.Trace(STEP_BACK - 1)
	rec = trace.GetRecord()
	if rec.FileNameShort != `trace_test.go` ||
		rec.FileLine != 31 ||
		rec.Function != `TestTraceTrace` ||
		rec.Module != `trace` {
		t.Errorf("Error in Trace().FileNameShort")
		return
	}
}
