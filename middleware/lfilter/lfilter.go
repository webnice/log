package lfilter // import "github.com/webnice/log/middleware/lfilter"

import (
	"sync"

	l "github.com/webnice/log/level"
	s "github.com/webnice/log/sender"
)

//import "github.com/webnice/debug/v1"

// Interface is an interface of package
type Interface interface {
	// Receiver Message receiver
	Receiver(s.Message)
}

// ReceiverFn Receiver function
type ReceiverFn func(s.Message)

// Filter settings
type Filter map[l.Level]ReceiverFn

// impl is an implementation of package
type impl struct {
	Filter Filter
	sync.Mutex
}

// New Create new
func New(f ...Filter) Interface {
	var i int
	var j l.Level
	var lft = new(impl)
	lft.Filter = make(Filter)
	for i = range f {
		for j = range f[i] {
			lft.Filter[j] = f[i][j]
		}
	}
	return lft
}

// Receiver Message receiver and send selected level to other receiver
func (lft *impl) Receiver(msg s.Message) {
	var ok bool
	var fn ReceiverFn
	lft.Lock()
	defer lft.Unlock()
	if fn, ok = lft.Filter[msg.Level]; ok {
		if fn != nil {
			fn(msg)
		}
	}
}
