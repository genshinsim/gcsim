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
		c.Flags.LogDebug = true
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

func NewCore(seed int64, debug bool, cfg core.SimulatorSettings) *core.Core {
	c := newCoreNoQueue(seed, debug)
	switch cfg.QueueMode {
	case core.ActionPriorityList:
		c.Queue = queue.NewQueuer(c)
	case core.SequentialList:
		c.Queue = calcqueue.New(c)
	}

	if cfg.SwapDelay > 0 {
		c.Flags.SwapFrames = cfg.SwapDelay
	}

	return c
}
