// +build !windows

package syslog

import (
	"bytes"
	"fmt"
	"log/syslog"
	"os"
	"strings"

	f "gopkg.in/webnice/log.v2/formater"
	l "gopkg.in/webnice/log.v2/level"
	s "gopkg.in/webnice/log.v2/sender"
)

// const _DefaultTextFORMAT = `%{color}[%{module:-10s}] %{time:2006-01-02T15:04:05.000Z07:00t} (%{level:-8s}): %{message} (%{package}) (%{function}:%{line}) (%{shortfile}:%{line}) (%{longfile})`
const _DefaultTextFORMAT = `%{message} {%{package}/%{shortfile}:%{line}, func: %{function}()}`
const (
	_DefaultNetwork  = `udp`
	_DefaultAddress  = `localhost:514`
	_DefaultPriority = syslog.LOG_INFO
	_DefaultPreffix  = ``
)

// Interface is an interface of package
type Interface interface {
	// Receiver Message receiver
	Receiver(s.Message)

	// SetAddress Назначение адреса syslog сервера
	SetAddress(string, string) Interface
}

// impl is an implementation of package
type impl struct {
	Formater f.Interface // Formater interface
	TplText  string      // Шаблон форматирования текста
	Network  string
	Address  string
	Priority syslog.Priority
	Preffix  string
}

// New Create new
func New() Interface {
	var rcv = new(impl)
	rcv.TplText = _DefaultTextFORMAT
	rcv.Formater = f.New()
	rcv.Network = _DefaultNetwork
	rcv.Address = _DefaultAddress
	rcv.Priority = _DefaultPriority
	rcv.Preffix = _DefaultPreffix
	return rcv
}

// SetAddress Назначение адреса syslog сервера
func (rcv *impl) SetAddress(proto string, address string) Interface {
	switch strings.ToLower(proto) {
	case "udp", "tcp":
		rcv.Network = strings.ToLower(proto)
	default:
		rcv.Network = _DefaultNetwork
	}
	rcv.Address = address
	return rcv
}

// Receiver Message receiver. Output to Syslog
func (rcv *impl) Receiver(msg s.Message) {
	var err error
	var buf *bytes.Buffer
	var level = l.New()
	var wr *syslog.Writer

	if rcv.Network == _DefaultNetwork && rcv.Address == _DefaultAddress {
		wr, err = syslog.New(rcv.Priority, rcv.Preffix)
	} else {
		wr, err = syslog.Dial(rcv.Network, rcv.Address, rcv.Priority, rcv.Preffix)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error dial to syslog: %s\n", err.Error())
		return
	}
	defer wr.Close()
	if buf, err = rcv.Formater.Text(msg, rcv.TplText); err != nil {
		buf = bytes.NewBufferString(fmt.Sprintf("Error formatting log message: %s", err.Error()))
	}
	switch msg.Level {
	case level.Fatal():
		err = wr.Emerg(buf.String())
	case level.Alert():
		err = wr.Alert(buf.String())
	case level.Critical():
		err = wr.Crit(buf.String())
	case level.Error():
		err = wr.Err(buf.String())
	case level.Warning():
		err = wr.Warning(buf.String())
	case level.Notice():
		err = wr.Notice(buf.String())
	case level.Informational():
		err = wr.Info(buf.String())
	case level.Debug():
		err = wr.Debug(buf.String())
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error write to syslog: %s\n", err.Error())
	}
}
