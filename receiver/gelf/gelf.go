package gelf // import "github.com/webnice/log/v2/receiver/gelf"

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	f "github.com/webnice/log/v2/formater"
	g "github.com/webnice/log/v2/gelf"
	s "github.com/webnice/log/v2/sender"
	t "github.com/webnice/log/v2/trace"
)

const (
	// defaultTextFORMAT = `%{color}[%{module:-10s}] %{time:2006-01-02T15:04:05.000Z07:00t} (%{level:-8s}): %{message} (%{package}) (%{function}:%{line}) (%{shortfile}:%{line}) (%{longfile})`
	defaultTextFORMAT = `%{message}`

	keyTCP             = `tcp`
	keyUDP             = `udp`
	gelfVersion        = "1.2"
	defaultNetwork     = keyUDP
	defaultCompression = g.COMPRESSION_NONE
	defaultChunkSize   = uint(1400)
)

// Interface is an interface of package
type Interface interface {
	// Receiver Message receiver
	Receiver(s.Message)

	// SetAddress Назначение адреса graylog2 сервера
	SetAddress(proto string, host string, port uint16) Interface

	// SetCompression Назначение сжатия для отправляемых данных
	SetCompression(compress string) Interface
}

// impl is an implementation of package
type impl struct {
	Formater    f.Interface   // Formater interface
	TplText     string        // Шаблон форматирования текста
	Client      *g.GelfClient // GELF client interface
	Network     string
	Host        string
	Port        uint16
	ChunkSize   uint
	Compression g.CompressionType
}

type msg struct {
	Version      string                 `json:"version"`                // [required] GELF spec version
	Host         string                 `json:"host"`                   // [required] he name of the host, source or application that sent this message
	ShortMessage string                 `json:"short_message"`          // [required] A short descriptive message
	FullMessage  string                 `json:"full_message,omitempty"` // A long message that can i.e. contain a backtrace
	Timestamp    float64                `json:"timestamp"`              // [required] Seconds since UNIX epoch with optional decimal places for milliseconds; SHOULD be set by client library
	Level        int8                   `json:"level"`                  // [required] The level equal to the standard syslog levels
	LevelString  string                 `json:"levelString"`            // Уровень лога в строковом эквиваленте
	Facility     string                 `json:"facility,omitempty"`     // [deprecated] Объект или пакет отправляющий сообщение
	Keys         map[string]interface{} `json:"keys,omitempty"`         // Дополнительные ключи сообщения и их значения
	t.Info
}

// New Create new
func New() Interface {
	var rcv = &impl{
		Formater:    f.New(),
		TplText:     defaultTextFORMAT,
		Client:      nil,
		Network:     defaultNetwork,
		ChunkSize:   defaultChunkSize,
		Compression: defaultCompression,
	}
	runtime.SetFinalizer(rcv, destructor)

	return rcv
}

func destructor(obj *impl) {
	if obj.Client == nil {
		return
	}
	defer func() { _ = obj.Client.Close() }()
}

// SetAddress Назначение адреса syslog сервера
func (rcv *impl) SetAddress(proto string, host string, port uint16) Interface {
	switch strings.ToLower(proto) {
	case keyUDP, keyTCP:
		rcv.Network = strings.ToLower(proto)
	default:
		rcv.Network = defaultNetwork
	}
	rcv.Host = host
	rcv.Port = port

	return rcv
}

// SetCompression Назначение сжатия для отправляемых данных
func (rcv *impl) SetCompression(compress string) Interface {
	rcv.Compression = g.CompressionType(compress)
	return rcv
}

func (rcv *impl) client() (ret *g.GelfClient, err error) {
	var pc g.GelfProtocolClient

	switch rcv.Network {
	case keyUDP:
		pc, err = g.NewUdpClient(rcv.Host, rcv.Port, defaultChunkSize)
	case keyTCP:
		pc, err = g.NewTcpClient(rcv.Host, rcv.Port)
	}
	if err != nil {
		return
	}
	ret = g.NewGelfClient(pc, rcv.Compression)

	return
}

// Receiver Message receiver. Output to STDERR
func (rcv *impl) Receiver(inp s.Message) {
	var (
		err        error
		msb        *msg
		buf, short *bytes.Buffer
		i          int
		rn         rune
	)

	if rcv.Client == nil {
		rcv.Client, err = rcv.client()
	}
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "create graylog GELF client error: %s", err)
		return
	}
	// Create object to send
	msb = new(msg)
	msb.Info = *inp.Trace
	msb.Version = gelfVersion
	msb.Host = inp.Trace.HostName
	msb.Timestamp = float64(msb.TodayAndNow.Unix()) + float64(time.Second)/float64(msb.TodayAndNow.Nanosecond())
	msb.Level = inp.Level.Int8()
	msb.LevelString = inp.Level.String()
	msb.Facility = inp.Trace.Package
	// Copy keys
	msb.Keys = make(map[string]interface{})
	for key := range inp.Keys {
		msb.Keys[key] = inp.Keys[key]
	}
	// Create full message by template
	buf, err = rcv.Formater.Text(inp, rcv.TplText)
	if err != nil {
		buf = bytes.NewBufferString(fmt.Sprintf("formatting log message error: %s", err))
		_, _ = fmt.Fprintln(os.Stderr, buf.String())
		return
	}
	msb.FullMessage = buf.String()
	// Create short message
	short = bytes.NewBuffer(bytes.Replace(bytes.Replace(buf.Bytes(), []byte("\r"), []byte{}, -1), []byte("\n"), []byte(" "), -1))
	for i = 0; i < 76 && short.Len() > 0; i++ {
		if rn, _, err = short.ReadRune(); err != nil {
			break
		}
		msb.ShortMessage += string(rn)
	}
	if len(msb.ShortMessage) < len(msb.FullMessage) {
		msb.ShortMessage += "..."
	}
	// Send object
	if err = rcv.Client.SendMessage(msb); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "send message to graylog error: %s", err)
		return
	}
}
