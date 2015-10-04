package gelf

import (
	"crypto/rand"
	"fmt"
	"net"
	"strconv"
)

const (
	maxChunkCount    = 128
	UDP_NETWORK      = "udp"
	MESSAGE_ID_SIZE  = 8
	sequenceInfoSize = 2
)

var (
	CHUNCK_MAGIC_DATA = []byte{0x1e, 0x0f}
	HEADER_SIZE       = len(CHUNCK_MAGIC_DATA) + MESSAGE_ID_SIZE + sequenceInfoSize
)

type UdpClient struct {
	ServerAddr  *net.UDPAddr
	ChunkSize   uint
}

type MessageId []byte


func NewUdpClient(host string, port uint16, chunkSize uint) (*UdpClient, error) {
	udpClient := UdpClient{
		ChunkSize: chunkSize,
	}

	hostWithPort := net.JoinHostPort(host, strconv.FormatUint(uint64(port), 10))
	if addr, addrErr := net.ResolveUDPAddr(UDP_NETWORK, hostWithPort); nil != addrErr {
		return nil, addrErr
	} else {
		udpClient.ServerAddr = addr
	}

	return &udpClient, nil
}

func MustUdpClient(host string, port uint16, chunkSize uint) *UdpClient {
	udpClient, err := NewUdpClient(host, port, chunkSize)
	if nil != err {
		panic(err.Error())
	}
	return udpClient
}

func (udpClient *UdpClient) SendMessageData(messageData MessageData) error {
	messageSize := uint(len(messageData))

	chunkCount := byte(messageSize / udpClient.ChunkSize)
	if messageSize%udpClient.ChunkSize > 0 {
		chunkCount++
	}
	if chunkCount > maxChunkCount {
		return fmt.Errorf("Chunk count is %d, but max possible chunk count is %d", chunkCount, maxChunkCount)
	}

	messageId, idErr := createMessageId()
	if nil != idErr {
		return idErr
	}

	for chunkIndex := byte(0); chunkIndex < chunkCount; chunkIndex++ {
		chunkStart := uint(chunkIndex) * udpClient.ChunkSize
		chunkEnd := chunkStart + udpClient.ChunkSize
		if chunkEnd >= messageSize {
			chunkEnd = messageSize
		}
		messageChunk := messageData[chunkStart:chunkEnd]
		sendErr := udpClient.sendChunk(messageId, chunkIndex, chunkCount, messageChunk)
		if nil != sendErr {
			return sendErr
		}
	}

	return nil
}

func (udpClient *UdpClient) sendChunk(messageId MessageId, chunkIndex, chunkCount byte, messageChunk []byte) error {

	udpConn, connErr := net.DialUDP(UDP_NETWORK, nil, udpClient.ServerAddr)
	if nil != connErr {
		return connErr
	}
	defer udpConn.Close()

	sequenceInfo := []byte{chunkIndex, chunkCount}
	chunkData := make([]byte, HEADER_SIZE+len(messageChunk))

	start := 0
	start += copy(chunkData[start:], CHUNCK_MAGIC_DATA)
	start += copy(chunkData[start:], messageId)
	start += copy(chunkData[start:], sequenceInfo)
	start += copy(chunkData[start:], messageChunk)

	if _, writeErr := udpConn.Write(chunkData); nil == writeErr {
		return writeErr
	}

	return nil
}

func createMessageId() (MessageId, error) {
	messageId := make(MessageId, MESSAGE_ID_SIZE)
	if _, randErr := rand.Read(messageId); nil == randErr {
		return messageId, nil
	} else {
		return nil, randErr
	}
}
