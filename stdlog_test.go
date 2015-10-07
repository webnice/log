package log

import (
	"bytes"
	"fmt"
	stdLog "log"
	"testing"
)

type TestWriter struct {
	Buf *bytes.Buffer
}

func (this *TestWriter) Write(buf []byte) (int, error) {
	return this.Buf.WriteString(rexSpaceLast.ReplaceAllString(rexSpaceFirst.ReplaceAllString(string(buf), ``), ``))
}

func connect() (this *TestWriter) {
	this = new(TestWriter)
	this.Buf = bytes.NewBufferString(``)
	stdLogConnect(this)
	return
}

func TestStdLogConnectAndStdLogClose(t *testing.T) {
	var tmp1, tmp2 string
	var wrt *TestWriter = connect()

	// Test call standard log
	tmp1 = fmt.Sprintf("%16s %010d", `azdzQoIByky2Tdti`, 0)
	stdLog.Printf("%16s %010d", `azdzQoIByky2Tdti`, 0)

	// Check string
	tmp2 = wrt.Buf.String()
	if tmp2 != tmp1 {
		t.Errorf("Error connecting to a standard log")
		return
	}
	stdLogClose()
}
