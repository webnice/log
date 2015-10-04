package log

import (
	"regexp"
)

var (
	rexSpaceFirst *regexp.Regexp = regexp.MustCompile(`^[\t\n\f\r ]`)
	rexSpaceLast  *regexp.Regexp = regexp.MustCompile(`[\t\n\f\r ]$`)
)

// Сюда попадают все сообщения от сторонни логеров
// сообщениям присваивается дефолтовый уровень логирования

// Для []byte
func (this *log) Write(buf []byte) (l int, err error) {
	if this.Record == nil {
		this.Record = newTrace().Trace(traceStepBack + 2).GetRecord()
	}
	this.write(defaultLevel, rexSpaceLast.ReplaceAllString(rexSpaceFirst.ReplaceAllString(string(buf), ``), ``))
	l = this.WriteLen
	err = this.WriteErr
	return
}

// Для string
func (this *log) WriteString(buf string) (l int, err error) {
	if this.Record == nil {
		this.Record = newTrace().Trace(traceStepBack + 2).GetRecord()
	}
	this.write(defaultLevel, rexSpaceLast.ReplaceAllString(rexSpaceFirst.ReplaceAllString(buf, ``), ``))
	l = this.WriteLen
	err = this.WriteErr
	return
}
