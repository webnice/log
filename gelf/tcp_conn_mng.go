package gelf

import (
	"net"
	"sync"
	"time"
)

type TcpConnMng struct {
	sync.RWMutex
	net           string
	laddr         *net.TCPAddr
	raddr         *net.TCPAddr
	tcpConnection *net.TCPConn
}

func DialTcpMng(net string, laddr, raddr *net.TCPAddr) (ret *TcpConnMng, err error) {
	ret = &TcpConnMng{
		net:   net,
		laddr: laddr,
		raddr: raddr,
	}
	if err = ret.reconnect(); err != nil {
		ret = nil
		return
	}
	return
}

func (tcpConnMng *TcpConnMng) connection() *net.TCPConn {
	tcpConnMng.RLock()
	defer tcpConnMng.RUnlock()
	return tcpConnMng.tcpConnection
}

func (tcpConnMng *TcpConnMng) reconnect() (err error) {
	var tcpConnection *net.TCPConn
	tcpConnMng.Lock()
	defer tcpConnMng.Unlock()
	if tcpConnMng.tcpConnection != nil {
		_ = tcpConnMng.tcpConnection.Close()
		tcpConnMng.tcpConnection = nil
	}
	tcpConnection, err = net.DialTCP(tcpConnMng.net, tcpConnMng.laddr, tcpConnMng.raddr)
	if err != nil {
		return
	}
	tcpConnMng.tcpConnection = tcpConnection
	return
}

func (tcpConnMng *TcpConnMng) Read(b []byte) (n int, err error) {
	if n, err = tcpConnMng.connection().Read(b); err == nil {
		return
	}
	if err = tcpConnMng.reconnect(); err != nil {
		n = 0
		return
	}
	n, err = tcpConnMng.connection().Read(b)
	return
}

func (tcpConnMng *TcpConnMng) Write(b []byte) (n int, err error) {
	if n, err = tcpConnMng.connection().Write(b); err == nil {
		return
	}
	if err = tcpConnMng.reconnect(); err != nil {
		n = 0
		return
	}
	n, err = tcpConnMng.connection().Write(b)
	return
}

func (tcpConnMng *TcpConnMng) Close() error { return tcpConnMng.connection().Close() }

func (tcpConnMng *TcpConnMng) LocalAddr() net.Addr { return tcpConnMng.connection().LocalAddr() }

func (tcpConnMng *TcpConnMng) RemoteAddr() net.Addr { return tcpConnMng.connection().RemoteAddr() }

func (tcpConnMng *TcpConnMng) SetDeadline(t time.Time) error {
	return tcpConnMng.connection().SetDeadline(t)
}

func (tcpConnMng *TcpConnMng) SetReadDeadline(t time.Time) error {
	return tcpConnMng.connection().SetReadDeadline(t)
}

func (tcpConnMng *TcpConnMng) SetWriteDeadline(t time.Time) error {
	return tcpConnMng.connection().SetWriteDeadline(t)
}
