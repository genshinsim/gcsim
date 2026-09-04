package nefer

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/gadget"
)

const seedDuration = 12 * 60

type seedGadget struct {
	*gadget.Gadget
}

func newSeedGadget(c *core.Core, pos info.Point) *seedGadget {
	g := &seedGadget{}
	g.Gadget = gadget.New(c, pos, 2, info.GadgetTypDendroCore) // ~~~
	g.Duration = seedDuration
	return g
}

func (g *seedGadget) HandleAttack(atk *info.AttackEvent) float64 {
	g.Core.Events.Emit(event.OnGadgetHit, g, atk)
	return 0
}

func (c *char) startSeedWindow() {
	if c.Core.Player.GetMoonsignLevel() < 2 {
		return
	}

	c.AddStatus(seedWindowKey, 15*60, true)
	c.convertExistingDendroCores()
}

func (c *char) convertExistingDendroCores() {
	for _, gadget := range c.Core.Combat.GadgetsWithinArea(combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.Player(), nil, 20), nil) {
		if gadget == nil {
			continue
		}
		if gadget.GadgetTyp() != info.GadgetTypDendroCore {
			continue
		}
		if c.isSeedOfDeceit(gadget) {
			continue
		}
		c.replaceDendroCoreWithSeed(gadget)
	}
}

func (c *char) replaceDendroCoreWithSeed(core info.Gadget) {
	pos := core.Pos()
	c.Core.Combat.RemoveGadget(core.Key())
	seed := newSeedGadget(c.Core, pos)
	c.Core.Combat.AddGadget(seed)
	if c.Core.Flags.LogDebug {
		c.Core.Log.NewEvent("nefer seed of deceit created", glog.LogCharacterEvent, c.Index()).Write("expiry", c.Core.F+seedDuration)
	}
}

func (c *char) isSeedOfDeceit(g info.Gadget) bool {
	if g == nil {
		return false
	}
	_, ok := g.(*seedGadget)
	return ok
}

func (c *char) absorbSeedsOfDeceit() {
	seeds := c.Core.Combat.GadgetsWithinArea(
		combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			nil,
			seedAbsorbRadius,
		),
		func(g info.Gadget) bool {
			return c.isSeedOfDeceit(g)
		},
	)
	if len(seeds) == 0 {
		return
	}
	count := len(seeds)
	for _, seed := range seeds {
		c.Core.Combat.RemoveGadget(seed.Key())
	}
	c.addVeilStacks(count)
	if c.Core.Flags.LogDebug {
		c.Core.Log.NewEvent("nefer seeds of deceit absorbed", glog.LogCharacterEvent, c.Index()).Write("absorbed", count).Write("radius_assumption", seedAbsorbRadius).Write("veil", c.currentVeilStacks())
	}
}
