package writer

import (
	"regexp"

	l "github.com/webnice/lv2/level"
)

var (
	rexSpaceFirst *regexp.Regexp = regexp.MustCompile(`^[\t\n\f\r ]+`)
	rexSpaceLast  *regexp.Regexp = regexp.MustCompile(`[\t\n\f\r ]+$`)
)

// Interface is an interface of package
type Interface interface {
	// Writer for []byte
	Write([]byte) (int, error)

	// Writer for string
	WriteString(string) (int, error)

	// Writer for ...any
	Println(...interface{})
}

// impl is an implementation of package
type impl struct {
	level l.Level
}
