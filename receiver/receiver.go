package receiver

//import "gopkg.in/webnice/debug.v1"
import (
	"gopkg.in/webnice/log.v2/receiver/fsfile"
	"gopkg.in/webnice/log.v2/receiver/fsfilerotation"
	"gopkg.in/webnice/log.v2/receiver/gelf"
	"gopkg.in/webnice/log.v2/receiver/stderr"
	"gopkg.in/webnice/log.v2/receiver/stdout"
	"gopkg.in/webnice/log.v2/receiver/syslog"
)

var (
	// StderrReceiver Read message and output to STDERR
	StderrReceiver = stderr.New()

	// StdoutReceiver Read message and output to STDOUT
	StdoutReceiver = stdout.New()

	// SyslogReceiver Read message and output to SYSLOG
	SyslogReceiver = syslog.New()

	// FsFileReceiver Read message and output to file
	FsFileReceiver = fsfile.New()

	// FsFilerotationReceiver Read message and output to file with time rotation
	FsFilerotationReceiver = fsfilerotation.New()

	// GelfReceiver Read message and output to graylog2 over GELF
	GelfReceiver = gelf.New()
)

// Default receiver
var Default = StderrReceiver
