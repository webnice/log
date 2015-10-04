package gelf_test

import (
	"bytes"
	"testing"
	"encoding/json"

	gelf "."
)

type TestingClient struct {
	MessageDataBuffer []gelf.MessageData
}

func NewTestingClient() *TestingClient {
	return &TestingClient {
		MessageDataBuffer: []gelf.MessageData{},
	}
}

func (testingClient *TestingClient) SendMessageData(messageData gelf.MessageData) error {
	testingClient.MessageDataBuffer = append(testingClient.MessageDataBuffer, messageData)
	return nil
}

func TestNewGelfClient(t *testing.T) {
	testingClient := NewTestingClient()
	gelfClient := gelf.NewGelfClient(testingClient, gelf.COMPRESSION_NONE)
	if nil == gelfClient {
		t.Error("NewGelfClient not return gelf clietn")
	}
}

type CustomMessage struct {
	*gelf.Message
	SerialNumber string `json:"_serial_number"`
}

func TestGelfClient_SendMessage(t *testing.T) {
	testingClient := NewTestingClient()
	gelfClient := gelf.NewGelfClient(testingClient, gelf.COMPRESSION_NONE)

	source := "custom_host"
	shortMessage := "custom short"
	level := gelf.LEVEL_WARNING
	serialNumber := "1234-5678-8912"

	customMessage := CustomMessage{
		Message:      gelf.NewMessage(source, level, shortMessage),
		SerialNumber: serialNumber,
	}

	if sendErr := gelfClient.SendMessage(customMessage); nil != sendErr {
		t.Errorf("Send error: %s", sendErr)
		return
	}

	if messageCount := len(testingClient.MessageDataBuffer); messageCount != 1 {
		t.Errorf("Mesasge data should be sended only one time, but sended %d", messageCount)
		return
	}

	messageData := testingClient.MessageDataBuffer[0]
	messageMap := make(map[string]interface{})

	if unmarshalErr := json.Unmarshal(messageData, &messageMap); nil != unmarshalErr {
		t.Errorf("Unmarshal error: %s", unmarshalErr)
		return
	}

	if version := messageMap["version"]; gelf.GELF_VERSION != version {
		t.Errorf(`version is a "%+v", but should be %s`, version, gelf.GELF_VERSION)
	}

	if host := messageMap["host"]; source != host {
		t.Errorf(`host is a "%+v", but should be %s`, host, source)
	}

	if timestamp := messageMap["timestamp"]; nil == timestamp {
		t.Errorf("timestamp is not specified")
	} else if timestampValue := timestamp.(float64); timestampValue == 0 {
		t.Errorf("timestamp is zero")
	}

	if messageLevel := messageMap["level"]; int(level) != int(messageLevel.(float64)) {
		t.Errorf(`Invalid level value, should be "%+v", but have "%+v"`, level, messageLevel)
	}

	if shortMessage != messageMap["short_message"] {
		t.Errorf(`Invalid short message ("%+v" != "%+v")`, shortMessage, messageMap["short_message"])
	}

	if messageSerialNumber := messageMap["_serial_number"]; serialNumber != messageSerialNumber {
		t.Errorf(`Invalid value for custom field (%+v != %+v)`, serialNumber, messageSerialNumber)
	}
}

func TestGelfClient_SendMessage_Gzip(t *testing.T) {
	testingClient := NewTestingClient()
	gelfClient := gelf.NewGelfClient(testingClient, gelf.COMPRESSION_GZIP)

	source := "custom_host"
	shortMessage := "custom short"
	level := gelf.LEVEL_WARNING

	customMessage := CustomMessage{
		Message:      gelf.NewMessage(source, level, shortMessage),
		SerialNumber: "1234-5678-8912",
	}

	if sendErr := gelfClient.SendMessage(customMessage); nil != sendErr {
		t.Errorf("Send error: %s", sendErr)
	}


	if messageCount := len(testingClient.MessageDataBuffer); messageCount != 1 {
		t.Errorf("Mesasge data should be sended only one time, but sended %d", messageCount)
		return
	}

	var messageData = testingClient.MessageDataBuffer[0]

	if !bytes.HasPrefix(messageData, gelf.GZIP_MAGIC_PREFIX) {
		t.Errorf("Gzip magic prefix not found")
	}
}