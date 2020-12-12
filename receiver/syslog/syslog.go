// +build !windows

package syslog

import (
	"bytes"
	"fmt"
	"log/syslog"
	"os"
	"strings"

	f "github.com/webnice/lv2/formater"
	l "github.com/webnice/lv2/level"
	s "github.com/webnice/lv2/sender"
)

const (
	// defaultTextFORMAT = `%{color}[%{module:-10s}] %{time:2006-01-02T15:04:05.000Z07:00t} (%{level:-8s}): %{message} (%{package}) (%{function}:%{line}) (%{shortfile}:%{line}) (%{longfile})`
	defaultTextFORMAT = `%{message} {%{package}/%{shortfile}:%{line}, func: %{function}()}`

	keyTCP          = `tcp`
	keyUDP          = `udp`
	defaultNetwork  = keyUDP
	defaultAddress  = `localhost:514`
	defaultPriority = syslog.LOG_INFO
	defaultPrefix   = ``
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
	Prefix   string
}

// New Create new
func New() Interface {
	var rcv = &impl{
		Formater: f.New(),
		TplText:  defaultTextFORMAT,
		Network:  defaultNetwork,
		Address:  defaultAddress,
		Priority: defaultPriority,
		Prefix:   defaultPrefix,
	}
	return rcv
}

// SetAddress Назначение адреса syslog сервера
func (rcv *impl) SetAddress(proto string, address string) Interface {
	switch strings.ToLower(proto) {
	case keyUDP, keyTCP:
		rcv.Network = strings.ToLower(proto)
	default:
		rcv.Network = defaultNetwork
	}
	rcv.Address = address

	return rcv
}

// Receiver Message receiver. Output to Syslog
func (rcv *impl) Receiver(msg s.Message) {
	var (
		err   error
		buf   *bytes.Buffer
		wr    *syslog.Writer
		level l.Interface
	)

	level = l.New()
	if rcv.Network == defaultNetwork && rcv.Address == defaultAddress {
		wr, err = syslog.New(rcv.Priority, rcv.Prefix)
	} else {
		wr, err = syslog.Dial(rcv.Network, rcv.Address, rcv.Priority, rcv.Prefix)
	}
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "syslog dial error: %s\n", err)
		return
	}
	defer func() { _ = wr.Close() }()
	if buf, err = rcv.Formater.Text(msg, rcv.TplText); err != nil {
		buf = bytes.NewBufferString(fmt.Sprintf("formatting log message error: %s", err))
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
		_, _ = fmt.Fprintf(os.Stderr, "Error write to syslog: %s\n", err.Error())
	}
}
