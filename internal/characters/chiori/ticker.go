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
	c     *core.Core
	alive bool

	cb       func()
	interval int

	onDeath func()
	queuer
}

type queuer func(cb func(), delay int)

// kill stops any existing ticker from ticking
func (g *ticker) kill() {
	g.alive = false
	g.cb = nil
	g.interval = 0
	if g.onDeath != nil {
		g.onDeath()
	}
}

func newTicker(c *core.Core, life int, q queuer) *ticker {
	// note we don't check if life <= 0 here
	// if life is <= 0 then this will cause gadget to kill itself
	// the next time tasks are checked
	g := &ticker{
		alive:  true,
		c:      c,
		queuer: q,
	}
	if g.queuer == nil {
		g.queuer = c.Tasks.Add
	}
	g.queuer(func() {
		if !g.alive {
			return
		}
		g.kill()
	}, life)
	return g
}

func (g *ticker) tick() {
	// do nothing if gadget is dead
	if !g.alive {
		return
	}
	// execute cb
	if g.cb != nil {
		g.cb()
	}
	// queue next action
	if g.interval > 0 {
		g.queuer(g.tick, g.interval)
	}
}
