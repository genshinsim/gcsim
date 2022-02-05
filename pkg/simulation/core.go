package simulation

import (
	"log"
	"math/rand"

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
	"go.uber.org/zap"
)

func newCore(seed int64, logger *zap.SugaredLogger) *core.Core {
	c := core.New()
	c.Log = logger
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
	c.Queue = queue.NewQueuer(c)
	return c
}

func NewDefaultCoreWithDefaultLogger(seed int64) *core.Core {
	logger, err := core.NewDefaultLogger(false, false, nil)
	if err != nil {
		log.Panicf("error building default logger, shouldn't happen: %v\n", err)
	}
	c := newCore(seed, logger)

	return c
}

func NewDefaultCoreWithCustomLogger(seed int64, logger *zap.SugaredLogger) *core.Core {
	return newCore(seed, logger)
}
