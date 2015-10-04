package gelf_test

import (
	"bytes"
	"fmt"
	"math"
	"net"
	"testing"

	gelf "."
)

const (
	udpHost = "127.0.0.1"
	udpPort = uint16(12201)
	chunkSize  = 1500
)

var (
	indexOffset     = len(gelf.CHUNCK_MAGIC_DATA) + gelf.MESSAGE_ID_SIZE
	countOffset     = indexOffset + 1
)

func TestNewUdpClient(t *testing.T) {
	gelfClient, gelfErr := gelf.NewUdpClient(udpHost, udpPort, chunkSize)
	if nil != gelfErr {
		t.Error(gelfErr)
	}

	if nil == gelfClient {
		t.Error("gelfClient is nil, but shuld be not")
	}
}

func TestUdpClient_SendMessageData_Simple(t *testing.T) {
	gelfClient, gelfErr := gelf.NewUdpClient(udpHost, udpPort, chunkSize)
	if nil != gelfErr {
		t.Fatal(gelfErr)
	}

	shortMessageData := createLongMessageData(0.5)

	listenChunks(1, func(serverBuff chan []byte, serverErr chan error) {
		if sendErr := gelfClient.SendMessageData(shortMessageData); nil != sendErr {
			t.Fatal(sendErr)
		}
		select {
		case shortBuff := <-serverBuff:
			if chunkErr := validateChunk(shortBuff, 1); nil != chunkErr {
				t.Error(chunkErr)
			}
		case shortErr := <-serverErr:
			t.Error(shortErr)
		}
	})
}

func TestUdpClient_SendMessageData_Chunks(t *testing.T) {
	gelfClient, gelfErr := gelf.NewUdpClient(udpHost, udpPort, chunkSize)
	if nil != gelfErr {
		t.Fatal(gelfErr)
	}

	longMessageMultiple := 2.5
	chunkCount := byte(math.Ceil(longMessageMultiple))
	longMessageData := createLongMessageData(longMessageMultiple)

	listenChunks(chunkCount, func(serverBuff chan []byte, serverErr chan error) {
		if sendErr := gelfClient.SendMessageData(longMessageData); nil != sendErr {
			t.Fatal(sendErr)
		}
		for i := byte(0); i < chunkCount; i++ {
			select {
			case longBuff := <-serverBuff:
				if chunkErr := validateChunk(longBuff, chunkCount); nil != chunkErr {
					t.Error(chunkErr)
				}
			case longErr := <-serverErr:
				t.Error(longErr)
			}
		}
	})

}

func handlerChunks(serverConn *net.UDPConn, chunkCount byte, serverBuff chan []byte, serverErr chan error) {
	for ; chunkCount > 0; chunkCount++ {
		buff := make([]byte, chunkSize)
		if _, _, readErr := serverConn.ReadFromUDP(buff); readErr == nil {
			serverBuff <- buff
		} else {
			serverErr <- readErr
		}
	}
}

func createLongMessageData(longMessageMultiple float64) gelf.MessageData {

	longMessageSize := int(longMessageMultiple * chunkSize)

	longMessageData := make(gelf.MessageData, longMessageSize)
	for i, _ := range longMessageData {
		longMessageData[i] = byte(i % 100)
	}

	return longMessageData

}

func listenChunks(chunkCount byte, callback func(chan []byte, chan error)) {
	buff := make(chan []byte, chunkCount)
	err := make(chan error, chunkCount)

	hostWithPort := fmt.Sprintf("%s:%d", udpHost, udpPort)
	serverAddr, addrErr := net.ResolveUDPAddr(gelf.UDP_NETWORK, hostWithPort)
	if nil != addrErr {
		panic(addrErr)
	}
	serverConn, connErr := net.ListenUDP(gelf.UDP_NETWORK, serverAddr)
	if nil != connErr {
		panic(connErr.Error())
	}
	defer serverConn.Close()
	go handlerChunks(serverConn, chunkCount, buff, err)
	callback(buff, err)
}

func validateChunk(chunkData []byte, chunkCount byte) error {

	if !bytes.HasPrefix(chunkData, gelf.CHUNCK_MAGIC_DATA) {
		return fmt.Errorf("Chunk not have magic prefix")
	}

	if countValue := chunkData[countOffset]; countValue != chunkCount {
		return fmt.Errorf("Invalid chunk count. Should be %d, but have %d",
			chunkCount, countValue)
	}

	return nil
}
