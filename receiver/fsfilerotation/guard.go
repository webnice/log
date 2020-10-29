package fsfilerotation // import "github.com/webnice/log/v2/receiver/fsfilerotation"

import "sync"

type fnGuard struct {
	enable bool
	fn     func()
	mutex  sync.Mutex
}

func (g *fnGuard) Enable() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.enable = true
}

func (g *fnGuard) Run() { g.fn() }
