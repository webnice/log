// Graylog2 GELF client

package gelf

import (
	"fmt"
	"bytes"
	"encoding/json"
	"compress/gzip"
	"compress/flate"
)

type MessageData []byte

type CompressionType string

const (
	COMPRESSION_NONE CompressionType = "none"
	COMPRESSION_GZIP CompressionType = "gzip"
	// COMPRESSION_ZLIB Compression = "zlib"
)

var (
	UNSUPPORTED_COMPRESSION_TYPE = fmt.Errorf("Unsupported compression type")
)

var GZIP_MAGIC_PREFIX = []byte{0x1F, 0x8B}

type GelfProtocolClient interface {
	SendMessageData(messageData MessageData) error
}

type GelfClient struct {
	GelfProtocolClient
	CompressionType
	CompressionLevel int
}

func NewGelfClient(gelfProtocolClient GelfProtocolClient, compressionType CompressionType) *GelfClient {
	return &GelfClient {
		GelfProtocolClient: gelfProtocolClient,
		CompressionType: compressionType,
		CompressionLevel: flate.DefaultCompression,
	}
}

func (gelfClient *GelfClient) SendMessage(message interface{}) error {
	marshaledMessage, marshalErr := json.Marshal(message)
	if nil != marshalErr {
		return marshalErr
	}

	var messageData MessageData

	switch gelfClient.CompressionType {
	case COMPRESSION_NONE:
		messageData = MessageData(marshaledMessage)
	case COMPRESSION_GZIP:
		var gzipErr error
		if messageData, gzipErr = gelfClient.gzipMessageData(marshaledMessage); nil != gzipErr {
			return gzipErr
		}
	default:
		return UNSUPPORTED_COMPRESSION_TYPE
	}

	return gelfClient.SendMessageData(messageData)
}

func (gelfClient *GelfClient) gzipMessageData(messageData []byte) (MessageData, error) {
	buff := new(bytes.Buffer)
	gzipWriter, _ := gzip.NewWriterLevel(buff, gelfClient.CompressionLevel)
	if _, writerErr := gzipWriter.Write(messageData); writerErr != nil {
		return nil, writerErr
	} else if closeErr := gzipWriter.Close(); nil != closeErr {
		return nil, closeErr
	} else {
		return MessageData(buff.Bytes()), nil
	}
}
