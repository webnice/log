// +build windows

package syslog

import (
	"fmt"
	"os"

	s "github.com/webnice/lv2/sender"
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
}

// New Create new
func New() Interface {
	var rcv = new(impl)
	return rcv
}

// SetAddress Назначение адреса syslog сервера
func (rcv *impl) SetAddress(proto string, address string) Interface {
	fmt.Fprintf(os.Stderr, "Syslog not implemented on windows platform\n")
	return rcv
}

// Receiver Message receiver. Output to Syslog
func (rcv *impl) Receiver(msg s.Message) {
	fmt.Fprintf(os.Stderr, "Syslog not implemented on windows platform\n")
}
