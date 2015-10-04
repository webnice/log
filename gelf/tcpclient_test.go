package gelf_test

import (
	"net"
	"testing"
	"strconv"
	"time"
	"bytes"

	gelf "."
)

const (
	tcpHost = "127.0.0.1"
	tcpPort = uint16(12201)
)

func TestNewTcpClient(t *testing.T) {
	listenTcpMessage(tcpHost, tcpPort, 10 * time.Millisecond, func(messageDataChan <- chan []byte, errorChan <- chan error) {
		validTcpClient, validTcpClientErr := gelf.NewTcpClient(tcpHost, tcpPort)
		if nil == validTcpClient {
			t.Error("Valid TcpClient is not created")
		}
		if nil != validTcpClientErr {
			t.Errorf("Valid TcpClient has error: %s", validTcpClientErr)
		}

		invalidTcpClient, invalidTcpClientErr := gelf.NewTcpClient("300.300.300.300", tcpPort)
		if nil != invalidTcpClient {
			t.Error("Invalid TcpClient is created")
		}
		if nil == invalidTcpClientErr {
			t.Error("Invalid TcpClient not has error")
		}
	})
}

func TestTcpClient_SendMessage(t *testing.T) {

	listenTcpMessage(tcpHost, tcpPort, 3 * time.Second, func(messageDataChan <- chan []byte, errorChan <- chan error) {
		tcpClient, tcpClientErr := gelf.NewTcpClient(tcpHost, tcpPort)
		if nil != tcpClientErr {
			t.Errorf("TcpClient error: %s", tcpClientErr)
			return
		}

		messageData := createLongMessageData(4.5)

		if sendErr := tcpClient.SendMessageData(messageData); nil != sendErr {
			t.Errorf("Send error: %s", sendErr)
			return
		}

		select {
		case sendedMessageData := <- messageDataChan:
			if false == bytes.Equal(sendedMessageData[0:len(sendedMessageData)-1], messageData) {
				t.Errorf("Recived data(len: %d) not equal to sended data(len: %d)", len(sendedMessageData) - 1, len(messageData))
			}
		case senderError := <- errorChan:
			t.Errorf("Sender error: %s", senderError)
		}
	})
}

func acceptTcpMessage(tcpListener *net.TCPListener, deadlineTime time.Time, messageDataChan chan <- []byte,  errorChan chan <- error) {
	tcpConn, acceptErr := tcpListener.AcceptTCP()

	if nil != acceptErr {
		errorChan <- acceptErr
	} else {
		defer tcpConn.Close()
	}

	if setDeadlineErr := tcpConn.SetDeadline(deadlineTime); nil != setDeadlineErr {
		errorChan <- setDeadlineErr
	}

	messageData := make([]byte, 2 << 16)

	if readSize, readErr := tcpConn.Read(messageData); nil != readErr {
		errorChan <- readErr
	} else {
		messageDataChan <- messageData[0:readSize]
	}

}

func listenTcpMessage(host string, port uint16, timeout time.Duration, callback func(<- chan []byte, <- chan error)) {
	hostWithPort := net.JoinHostPort(host, strconv.FormatUint(uint64(port), 10))
	tcpAdd, resolveErr := net.ResolveTCPAddr(gelf.TCP_NETWORK, hostWithPort)
	if nil != resolveErr {
		panic(resolveErr)
	}

	tcpListener, listenErr := net.ListenTCP(gelf.TCP_NETWORK, tcpAdd)
	if nil != listenErr {
		panic(listenErr)
	} else {
		defer tcpListener.Close()
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
