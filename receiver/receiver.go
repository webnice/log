package receiver // import "github.com/webnice/log/v2/receiver"

import (
	"github.com/webnice/log/v2/receiver/fsfile"
	"github.com/webnice/log/v2/receiver/fsfilerotation"
	"github.com/webnice/log/v2/receiver/gelf"
	"github.com/webnice/log/v2/receiver/stderr"
	"github.com/webnice/log/v2/receiver/stdout"
	"github.com/webnice/log/v2/receiver/syslog"
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
