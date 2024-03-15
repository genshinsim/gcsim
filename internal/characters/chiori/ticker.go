package chiori

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *char) kill(t *ticker) {
	if t != nil {
		t.kill()
	}
}

// generic ticker for dolls to use
type ticker struct {
	c         *core.Core
	alive     bool
	liveUntil int

	cb       func()
	interval int
}

// kill stops any existing ticker from ticking
func (g *ticker) kill() {
	g.alive = false
	g.cb = nil
	g.interval = 0
}

func newTicker(c *core.Core, life int) *ticker {
	// note we don't check if life <= 0 here
	// if life is <= 0 then this will cause gadget to kill itself
	// the next time tasks are checked
	g := &ticker{
		alive:     true,
		c:         c,
		liveUntil: c.F + life,
	}
	c.Tasks.Add(func() {
		g.kill()
	}, life)
	return g
}

func (g *ticker) tick() {
	// do nothing if gadget is dead
	if !g.alive {
		return
	}
	// sanity check: ideally shouldn't happen
	// TODO: is this < or <=
	if g.liveUntil < g.c.F {
		g.kill()
		return
	}
	// execute cb
	if g.cb != nil {
		g.cb()
	}
	// queue next action
	if g.interval > 0 && g.c.F+g.interval <= g.liveUntil {
		g.c.Tasks.Add(g.tick, g.interval)
	}
}
