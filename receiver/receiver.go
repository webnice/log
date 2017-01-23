package receiver // import "github.com/webdeskltd/log/receiver"

//import "github.com/webdeskltd/debug"
import (
	"github.com/webdeskltd/log/receiver/fsfile"
	"github.com/webdeskltd/log/receiver/fsfilerotation"
	"github.com/webdeskltd/log/receiver/stderr"
	"github.com/webdeskltd/log/receiver/stdout"
	"github.com/webdeskltd/log/receiver/syslog"
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
)

// Default receiver
var Default = StderrReceiver