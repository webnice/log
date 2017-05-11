package gelf

//import "gopkg.in/webnice/debug.v1"
import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	f "gopkg.in/webnice/log.v2/formater"
	g "gopkg.in/webnice/log.v2/gelf"
	s "gopkg.in/webnice/log.v2/sender"
	t "gopkg.in/webnice/log.v2/trace"
)

// const _DefaultTextFORMAT = `%{color}[%{module:-10s}] %{time:2006-01-02T15:04:05.000Z07:00t} (%{level:-8s}): %{message} (%{package}) (%{function}:%{line}) (%{shortfile}:%{line}) (%{longfile})`
const _DefaultTextFORMAT = `%{message}`
const _GelfVersion = "1.2"
const _DefaultNetwork = `udp`
const _DefaultCompression = g.COMPRESSION_NONE
const _DefaultChunkSize = uint(1400)

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
	Version     string                 `json:"version"`
	Timestamp   float64                `json:"timestamp"`
	Host        string                 `json:"host"`
	LevelString string                 `json:"levelString"`
	Facility    string                 `json:"facility,omitempty"`
	Message     string                 `json:"message,omitempty"`
	Keys        map[string]interface{} `json:"keys"`
	t.Info
}

// New Create new
func New() Interface {
	var rcv = new(impl)
	rcv.Formater = f.New()
	rcv.TplText = _DefaultTextFORMAT
	rcv.Network = _DefaultNetwork
	rcv.ChunkSize = _DefaultChunkSize
	rcv.Compression = _DefaultCompression
	runtime.SetFinalizer(rcv, destructor)
	return rcv
}

func destructor(obj *impl) {
	if obj.Client == nil {
		return
	}
	defer obj.Client.Close()
}

// SetAddress Назначение адреса syslog сервера
func (rcv *impl) SetAddress(proto string, host string, port uint16) Interface {
	switch strings.ToLower(proto) {
	case "udp", "tcp":
		rcv.Network = strings.ToLower(proto)
	default:
		rcv.Network = _DefaultNetwork
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
	case "udp":
		pc, err = g.NewUdpClient(rcv.Host, rcv.Port, _DefaultChunkSize)
	case "tcp":
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
	var err error
	var msb *msg
	var buf *bytes.Buffer

	if rcv.Client == nil {
		rcv.Client, err = rcv.client()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error create graylog GELF client: %s", err.Error())
		return
	}

	msb = new(msg)
	msb.Info = *inp.Trace
	msb.Version = _GelfVersion
	msb.Timestamp = float64(msb.TodayAndNow.Unix()) + float64(time.Second)/float64(msb.TodayAndNow.Nanosecond())
	msb.Host = inp.Trace.HostName
	msb.LevelString = inp.Level.String()
	msb.Facility = inp.Trace.Package
	msb.Keys = make(map[string]interface{})
	for key := range inp.Keys {
		msb.Keys[key] = inp.Keys[key]
	}

	buf, err = rcv.Formater.Text(inp, rcv.TplText)
	if err != nil {
		buf = bytes.NewBufferString(fmt.Sprintf("Error formatting log message: %s", err.Error()))
		fmt.Fprintln(os.Stderr, buf.String())
		return
	}
	msb.Message = buf.String()

	if err = rcv.Client.SendMessage(msb); err != nil {
		fmt.Fprintf(os.Stderr, "Error send message to graylog: %s", err.Error())
		return
	}
}
