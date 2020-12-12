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
	ServerAddr *net.UDPAddr
	ChunkSize  uint
}

type MessageId []byte

func NewUdpClient(host string, port uint16, chunkSize uint) (ret *UdpClient, err error) {
	var (
		address string
		addr    *net.UDPAddr
	)

	address = net.JoinHostPort(host, strconv.FormatUint(uint64(port), 10))
	addr, err = net.ResolveUDPAddr(UDP_NETWORK, address)
	if err != nil {
		return
	}
	ret = &UdpClient{
		ServerAddr: addr,
		ChunkSize:  chunkSize,
	}

	return
}

func MustUdpClient(host string, port uint16, chunkSize uint) (ret *UdpClient) {
	var err error

	if ret, err = NewUdpClient(host, port, chunkSize); err != nil {
		panic(err.Error())
	}

	return
}

func (udpClient *UdpClient) Close() error { return nil }

func createMessageId() (ret MessageId, err error) {
	ret = make(MessageId, MESSAGE_ID_SIZE)
	_, err = rand.Read(ret)
	return
}

func (udpClient *UdpClient) SendMessageData(messageData MessageData) (err error) {
	var (
		messageSize, chunkStart, chunkEnd uint
		chunkCount, chunkIndex            byte
		messageId                         MessageId
		messageChunk                      []byte
	)

	messageSize = uint(len(messageData))
	chunkCount = byte(messageSize / udpClient.ChunkSize)
	if messageSize%udpClient.ChunkSize > 0 {
		chunkCount++
	}
	if chunkCount > maxChunkCount {
		return fmt.Errorf("Chunk count is %d, but max possible chunk count is %d", chunkCount, maxChunkCount)
	}
	if messageId, err = createMessageId(); err != nil {
		return
	}
	for chunkIndex = byte(0); chunkIndex < chunkCount; chunkIndex++ {
		chunkStart = uint(chunkIndex) * udpClient.ChunkSize
		chunkEnd = chunkStart + udpClient.ChunkSize
		if chunkEnd >= messageSize {
			chunkEnd = messageSize
		}
		messageChunk = messageData[chunkStart:chunkEnd]
		err = udpClient.sendChunk(messageId, chunkIndex, chunkCount, messageChunk)
		if err != nil {
			return
		}
	}

	return
}

func (udpClient *UdpClient) sendChunk(messageId MessageId, chunkIndex, chunkCount byte, messageChunk []byte) (err error) {
	var (
		udpConn                 *net.UDPConn
		sequenceInfo, chunkData []byte
		start                   int
	)

	if udpConn, err = net.DialUDP(UDP_NETWORK, nil, udpClient.ServerAddr); err != nil {
		return
	}
	defer func() { _ = udpConn.Close() }()
	sequenceInfo = []byte{chunkIndex, chunkCount}
	chunkData = make([]byte, HEADER_SIZE+len(messageChunk))
	start = 0
	start += copy(chunkData[start:], CHUNCK_MAGIC_DATA)
	start += copy(chunkData[start:], messageId)
	start += copy(chunkData[start:], sequenceInfo)
	start += copy(chunkData[start:], messageChunk)
	if _, err = udpConn.Write(chunkData); err != nil {
		return
	}

	return
}
