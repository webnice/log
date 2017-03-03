package fsfilerotation

//import "gopkg.in/webnice/debug.v1"
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
