package gelf

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
)

const TCP_NETWORK = "tcp"

var MESSAGE_SEPARATOR = []byte{0}

type TcpClient struct {
	ServerAddr *net.TCPAddr
	connection net.Conn
}

func NewTcpClient(host string, port uint16) (*TcpClient, error) {
	hostWithPort := net.JoinHostPort(host, strconv.FormatUint(uint64(port), 10))

	ipAddr, resolveErr := net.ResolveTCPAddr(TCP_NETWORK, hostWithPort)
	if nil != resolveErr {
		return nil, resolveErr
	}

	connection, dialErr := DialTcpMng(TCP_NETWORK, nil, ipAddr)
	if nil != dialErr {
		return nil, dialErr
	}

	return &TcpClient{
		ServerAddr: ipAddr,
		connection: connection,
	}, nil

}

func MustTcpClient(host string, port uint16) *TcpClient {
	tcpClient, err := NewTcpClient(host, port)
	if nil != err {
		panic(err.Error())
	}
	return tcpClient
}

var compressionNotSupported = fmt.Errorf("Compression not supported")

func (tcpClient *TcpClient) SendMessageData(message MessageData) (err error) {
	if bytes.HasPrefix(message, GZIP_MAGIC_PREFIX) {
		return compressionNotSupported
	}
	messageWithSeparator := append(message, MESSAGE_SEPARATOR...)

	if _, err = tcpClient.connection.Write(messageWithSeparator); err != nil {
		return
	}
	return
}
