package nefer

import (
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (c *char) p1Active() bool {
	return c.Base.Ascension >= 1 && c.ascendantGleam
}

func (c *char) p1Init() {
	if !c.p1Active() {
		return
	}

	c.Core.Events.Subscribe(event.OnDendroCore, func(args ...any) {
		if !c.StatusIsActive(seedWindowKey) {
			return
		}

		core, ok := args[0].(info.Gadget)
		if !ok || core.GadgetTyp() != info.GadgetTypDendroCore {
			return
		}
		if c.isSeedOfDeceit(core) {
			return
		}

		c.replaceDendroCoreWithSeed(core)
	}, "nefer-p1-seed-conversion")
}

func (c *char) startSeedWindow() {
	if !c.p1Active() {
		return
	}

	c.AddStatus(seedWindowKey, 15*60, true)
	c.convertExistingDendroCores()
}

func (c *char) convertExistingDendroCores() {
	for _, gadget := range c.Core.Combat.Gadgets() {
		if gadget == nil || gadget.GadgetTyp() != info.GadgetTypDendroCore {
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
	seeds := c.Core.Combat.GadgetsWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, seedAbsorbRadius), func(g info.Gadget) bool {
		return c.isSeedOfDeceit(g)
	})
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
