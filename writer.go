package log

import (
	"regexp"
)

var (
	rexSpaceFirst *regexp.Regexp = regexp.MustCompile(`^[\t\n\f\r ]`)
	rexSpaceLast  *regexp.Regexp = regexp.MustCompile(`[\t\n\f\r ]$`)
)

// Сюда попадают все сообщения от сторонних логеров
// У сообщений удаляются пробельные символы, предшествующие и заканчивающие сообщение
// Сообщениям присваивается дефолтовый уровень логирования

// Для []byte
func (this *Writer) Write(buf []byte) (l int, err error) {
	var msg *logMessage = newLogMessage()
	msg.Record = newTrace().Trace(traceStepBack + 2).GetRecord()
	msg.Write(defaultLevel, rexSpaceLast.ReplaceAllString(rexSpaceFirst.ReplaceAllString(string(buf), ``), ``))
	l = msg.WriteLen
	err = msg.WriteErr
	return
}

// Для string
func (this *Writer) WriteString(buf string) (l int, err error) {
	var msg *logMessage = newLogMessage()
	msg.Record = newTrace().Trace(traceStepBack + 2).GetRecord()
	msg.Write(defaultLevel, rexSpaceLast.ReplaceAllString(rexSpaceFirst.ReplaceAllString(buf, ``), ``))
	l = msg.WriteLen
	err = msg.WriteErr
	return
}
