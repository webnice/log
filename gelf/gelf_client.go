package gelf

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"fmt"
)

const (
	COMPRESSION_NONE CompressionType = "none"
	COMPRESSION_GZIP CompressionType = "gzip"
	// COMPRESSION_ZLIB CompressionType = "zlib" // Не реализовано
)

var (
	UNSUPPORTED_COMPRESSION_TYPE = fmt.Errorf("unsupported compression type")
	GZIP_MAGIC_PREFIX            = []byte{0x1F, 0x8B}
)

type MessageData []byte

type CompressionType string

type GelfProtocolClient interface {
	SendMessageData(messageData MessageData) error
	Close() error
}

type GelfClient struct {
	GelfProtocolClient
	CompressionType
	CompressionLevel int
}

func NewGelfClient(gelfProtocolClient GelfProtocolClient, compressionType CompressionType) *GelfClient {
	return &GelfClient{
		GelfProtocolClient: gelfProtocolClient,
		CompressionType:    compressionType,
		CompressionLevel:   flate.DefaultCompression,
	}
}

func (gelfClient *GelfClient) SendMessage(message interface{}) (err error) {
	var (
		buf         []byte
		messageData MessageData
	)

	if buf, err = json.Marshal(message); err != nil {
		return
	}
	switch gelfClient.CompressionType {
	case COMPRESSION_NONE:
		messageData = MessageData(buf)
	case COMPRESSION_GZIP:
		messageData, err = gelfClient.gzipMessageData(buf)
		if err != nil {
			return
		}
	default:
		return UNSUPPORTED_COMPRESSION_TYPE
	}

	return gelfClient.SendMessageData(messageData)
}

func (gelfClient *GelfClient) gzipMessageData(messageData []byte) (ret MessageData, err error) {
	var (
		buff       *bytes.Buffer
		gzipWriter *gzip.Writer
	)

	buff = bytes.NewBufferString(``)
	gzipWriter, _ = gzip.NewWriterLevel(buff, gelfClient.CompressionLevel)
	if _, err = gzipWriter.Write(messageData); err != nil {
		return
	} else if err = gzipWriter.Close(); err != nil {
		return
	}
	ret = MessageData(buff.Bytes())

	return
}
