package gelf_test

import (
	"testing"

	gelf "."
)

func TestNewMessage(t *testing.T) {
	source := "camera"
	shortMessage := "Short message"
	message := gelf.NewMessage(source, gelf.LEVEL_INFO, shortMessage)

	if message.Version != gelf.GELF_VERSION {
		t.Errorf("Message have invalid version (have %s, but expected %s)", message.Version, gelf.GELF_VERSION)
	}

	if message.Host != source {
		t.Errorf("Message have invalid source (have %s, but expected %s)", message.Host, source)
	}

	if message.ShortMessage != shortMessage {
		t.Errorf("Message have invalid short message (have %s, but expected %s)", message.ShortMessage, shortMessage)
	}

	if message.Level != gelf.LEVEL_INFO {
		t.Errorf("Message have invalid level (have %d, but expected %d)", message.Level, gelf.LEVEL_INFO)
	}
}
