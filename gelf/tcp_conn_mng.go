package gelf

import (
	"net"
	"time"
	"sync"
)

type TcpConnMng struct {
	sync.RWMutex
	net             string
	laddr           *net.TCPAddr
	raddr           *net.TCPAddr
	tcpConnection   *net.TCPConn
}

func DialTcpMng(net string, laddr, raddr *net.TCPAddr) (*TcpConnMng, error) {
	tcpConnMng := &TcpConnMng {
		net:             net,
		laddr:           laddr,
		raddr:           raddr,
	}

	if connErr := tcpConnMng.reconnect(); nil != connErr {
		return nil, connErr
	} else {
		return tcpConnMng, nil
	}
}

func (tcpConnMng *TcpConnMng) connection() *net.TCPConn {
	tcpConnMng.RLock()
	defer tcpConnMng.RUnlock()
	return tcpConnMng.tcpConnection
}

func (tcpConnMng *TcpConnMng) reconnect() error {
	tcpConnMng.Lock()
	defer tcpConnMng.Unlock()

	if tcpConnMng.tcpConnection != nil {
		tcpConnMng.tcpConnection.Close()
		tcpConnMng.tcpConnection = nil
	}

	tcpConnection, tcpErr := net.DialTCP(
		tcpConnMng.net,
		tcpConnMng.laddr,
		tcpConnMng.raddr,
	)

	if nil != tcpErr {
		return tcpErr
	}

	tcpConnMng.tcpConnection = tcpConnection

	return nil
}

func (tcpConnMng *TcpConnMng) Read(b []byte) (n int, err error) {
	if n, err := tcpConnMng.connection().Read(b); nil == err {
		return n, err
	}

	if connErr := tcpConnMng.reconnect(); nil != connErr {
		return 0, connErr
	}

	return tcpConnMng.connection().Read(b)
}

func (tcpConnMng *TcpConnMng) Write(b []byte) (n int, err error) {
	if n, err := tcpConnMng.connection().Write(b); nil == err {
		return n, err
	}

	if connErr := tcpConnMng.reconnect(); nil != connErr {
		return 0, connErr
	}

	return tcpConnMng.connection().Write(b)
}

func (tcpConnMng *TcpConnMng) Close() error {
	return tcpConnMng.connection().Close()
}

func (tcpConnMng *TcpConnMng) LocalAddr() net.Addr {
	return tcpConnMng.connection().LocalAddr()
}

func (tcpConnMng *TcpConnMng) RemoteAddr() net.Addr {
	return tcpConnMng.connection().RemoteAddr()
}

func (tcpConnMng *TcpConnMng) SetDeadline(t time.Time) error {
	return tcpConnMng.connection().SetDeadline(t)
}

func (tcpConnMng *TcpConnMng) SetReadDeadline(t time.Time) error {
	return tcpConnMng.connection().SetReadDeadline(t)
}

func (tcpConnMng *TcpConnMng) SetWriteDeadline(t time.Time) error {
	return tcpConnMng.connection().SetWriteDeadline(t)
}
