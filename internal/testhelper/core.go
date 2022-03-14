package testhelper

import (
	"math/rand"
	"time"

	"github.com/genshinsim/gcsim/internal/tmpl/action"
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

func NewTestCore() *core.Core {
	c := core.New()
	c.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	c.Tasks = task.NewCtrl(&c.Frame)
	c.Events = event.NewCtrl(c)
	c.Status = status.NewCtrl(c)
	c.Energy = energy.NewCtrl(c)
	c.Combat = combat.NewCtrl(c)
	c.Constructs = construct.NewCtrl(c)
	c.Shields = shield.NewCtrl(c)
	c.Health = health.NewCtrl(c)
	c.Action = action.NewCtrl(c)
	c.Queue = queue.NewQueuer(c)
	return c
}
