package gelf

import (
	"net"
	"strconv"
	"testing"
	"time"
)

const (
	tcpHost = "127.0.0.1"
	tcpPort = uint16(12201)
)

func TestNewTcpClient(t *testing.T) {
	listenTcpMessage(tcpHost, tcpPort, 10*time.Millisecond, func(messageDataChan <-chan []byte, errorChan <-chan error) {
		validTcpClient, validTcpClientErr := NewTcpClient(tcpHost, tcpPort)
		if nil == validTcpClient {
			t.Error("Valid TcpClient is not created")
		}
		if nil != validTcpClientErr {
			t.Errorf("Valid TcpClient has error: %s", validTcpClientErr)
		}

		invalidTcpClient, invalidTcpClientErr := NewTcpClient("300.300.300.300", tcpPort)
		if nil != invalidTcpClient {
			t.Error("Invalid TcpClient is created")
		}
		if nil == invalidTcpClientErr {
			t.Error("Invalid TcpClient not has error")
		}
	})
}

func acceptTcpMessage(tcpListener *net.TCPListener, deadlineTime time.Time, messageDataChan chan<- []byte, errorChan chan<- error) {
	tcpConn, acceptErr := tcpListener.AcceptTCP()

	if nil != acceptErr {
		errorChan <- acceptErr
	} else {
		defer func() { _ = tcpConn.Close() }()
	}

	if setDeadlineErr := tcpConn.SetDeadline(deadlineTime); nil != setDeadlineErr {
		errorChan <- setDeadlineErr
	}

	//	messageData := make([]byte, 2 << 16)
	//	if readSize, readErr := tcpConn.Read(messageData); nil != readErr {
	//		errorChan <- readErr
	//	} else {
	//		messageDataChan <- messageData[0:readSize]
	//	}

}

func listenTcpMessage(host string, port uint16, timeout time.Duration, callback func(<-chan []byte, <-chan error)) {
	hostWithPort := net.JoinHostPort(host, strconv.FormatUint(uint64(port), 10))
	tcpAdd, resolveErr := net.ResolveTCPAddr(TCP_NETWORK, hostWithPort)
	if nil != resolveErr {
		panic(resolveErr)
	}
	tcpListener, listenErr := net.ListenTCP(TCP_NETWORK, tcpAdd)
	if nil != listenErr {
		panic(listenErr)
	} else {
		defer func() { _ = tcpListener.Close() }()
	}
	deadlineTime := time.Now().Add(timeout)
	if setDeadlineErr := tcpListener.SetDeadline(deadlineTime); nil != setDeadlineErr {
		panic(setDeadlineErr)
	}
	messageDataChan := make(chan []byte)
	errorChan := make(chan error)
	defer close(messageDataChan)
	defer close(errorChan)
	go acceptTcpMessage(tcpListener, deadlineTime, messageDataChan, errorChan)
	callback(messageDataChan, errorChan)
}
