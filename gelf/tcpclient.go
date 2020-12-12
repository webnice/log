package gelf

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
)

const TCP_NETWORK = "tcp"

var (
	MESSAGE_SEPARATOR             = []byte{0}
	_ErrorCompressionNotSupported = fmt.Errorf("compression not supported")
)

type TcpClient struct {
	ServerAddr *net.TCPAddr
	connection net.Conn
}

func NewTcpClient(host string, port uint16) (ret *TcpClient, err error) {
	var address string
	var ipAddr *net.TCPAddr
	var connection *TcpConnMng
	address = net.JoinHostPort(host, strconv.FormatUint(uint64(port), 10))
	if ipAddr, err = net.ResolveTCPAddr(TCP_NETWORK, address); err != nil {
		return
	}
	if connection, err = DialTcpMng(TCP_NETWORK, nil, ipAddr); err != nil {
		return
	}
	ret = &TcpClient{
		ServerAddr: ipAddr,
		connection: connection,
	}
	return
}

func MustTcpClient(host string, port uint16) (ret *TcpClient) {
	var err error

	if ret, err = NewTcpClient(host, port); err != nil {
		panic(err.Error())
	}

	return
}

func (tcpClient *TcpClient) SendMessageData(message MessageData) (err error) {
	var messageWithSeparator MessageData

	if bytes.HasPrefix(message, GZIP_MAGIC_PREFIX) {
		return _ErrorCompressionNotSupported
	}
	messageWithSeparator = append(message, MESSAGE_SEPARATOR...)
	if _, err = tcpClient.connection.Write(messageWithSeparator); err != nil {
		return
	}

	return
}

func (tcpClient *TcpClient) Close() error { return tcpClient.Close() }
