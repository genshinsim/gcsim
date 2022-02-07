package simulation

import (
	"math/rand"

	"github.com/genshinsim/gcsim/internal/evtlog"
	"github.com/genshinsim/gcsim/internal/tmpl/action"
	"github.com/genshinsim/gcsim/internal/tmpl/calcqueue"
	"github.com/genshinsim/gcsim/internal/tmpl/combat"
	"github.com/genshinsim/gcsim/internal/tmpl/construct"
	"github.com/genshinsim/gcsim/internal/tmpl/energy"
	"github.com/genshinsim/gcsim/internal/tmpl/event"
	"github.com/genshinsim/gcsim/internal/tmpl/health"
	"github.com/genshinsim/gcsim/internal/tmpl/queue"
	"github.com/genshinsim/gcsim/internal/tmpl/shield"
	"github.com/genshinsim/gcsim/internal/tmpl/status"
	"github.com/genshinsim/gcsim/internal/tmpl/task"
	"github.com/genshinsim/gcsim/pkg/core"
)

func newCoreNoQueue(seed int64, debug bool) *core.Core {
	c := core.New()
	if debug {
		c.Log = evtlog.NewCtrl(c, 500)
	}
	c.Rand = rand.New(rand.NewSource(seed))
	c.Tasks = task.NewCtrl(&c.F)
	c.Events = event.NewCtrl(c)
	c.Status = status.NewCtrl(c)
	c.Energy = energy.NewCtrl(c)
	c.Combat = combat.NewCtrl(c)
	c.Constructs = construct.NewCtrl(c)
	c.Shields = shield.NewCtrl(c)
	c.Health = health.NewCtrl(c)
	c.Action = action.NewCtrl(c)
	return c
}

func NewDefaultCoreWithCalcQueue(seed int64) *core.Core {
	c := newCoreNoQueue(seed, false)
	c.Queue = calcqueue.New(c)
	return c
}

func NewDefaultCore(seed int64) *core.Core {
	c := newCoreNoQueue(seed, false)
	c.Queue = queue.NewQueuer(c)
	return c
}

func NewDefaultCoreWithDebug(seed int64) *core.Core {
	c := newCoreNoQueue(seed, true)
	c.Queue = queue.NewQueuer(c)
	return c
}

func NewDefaultCoreWithDebugCalcQueue(seed int64) *core.Core {
	c := newCoreNoQueue(seed, true)
	c.Queue = calcqueue.New(c)
	return c
}
