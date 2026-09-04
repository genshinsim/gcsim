package nefer

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const (
	lunarbloomBonusKey = "nefer-lunarbloom-bonus"
	a4DewWindowKey     = "nefer-p2-dew-window"
	a4DewRateKey       = "nefer-p2-dew-rate"
	a4DewWindowDur     = 5 * 60
)

func (c *char) a1Init() {
	if c.Base.Ascension < 1 {
		return
	}

	if c.Core.Player.GetMoonsignLevel() < 2 {
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
	}, "nefer-a1-seed-conversion")
}

func (c *char) a4Init() {
	if c.Base.Ascension < 4 {
		return
	}

	c.Core.Player.AddVerdantDewRateMod(a4DewRateKey, -1, func() (float64, bool) {
		if !c.StatusIsActive(shadowDanceKey) || !c.StatusIsActive(a4DewWindowKey) {
			return 0, false
		}
		return c.a4BuffAmt, false
	})

	c.Core.Events.Subscribe(event.OnLunarBloom, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex < 0 || atk.Info.ActorIndex >= len(c.Core.Player.Chars()) {
			return
		}
		c.AddStatus(a4DewWindowKey, a4DewWindowDur, true)
		c.updateDewRateBonusStart(c.Core.F)
	}, a4DewWindowKey)
}

func (c *char) updateDewRateBonusStart(src int) {
	c.a4BuffSrc = src
	c.updateDewRateBonus(src)
}

func (c *char) updateDewRateBonus(src int) {
	if src != c.a4BuffSrc {
		return
	}

	if !c.StatusIsActive(a4DewWindowKey) {
		return
	}

	c.a4BuffAmt = min(max(c.Stat(attributes.EM)-500, 0)/1000, 0.5)

	c.QueueCharTask(func() {
		c.updateDewRateBonus(src)
	}, 0.5*60)
}

func (c *char) lunarbloomInit() {
	c.Core.Flags.Custom[reactable.LunarBloomEnableKey] = 1

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.AttackTag != attacks.AttackTagDirectLunarBloom {
			return
		}
		bonus := min(c.Stat(attributes.EM)*0.000175, 0.14)

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("nefer adding lunarbloom base damage", glog.LogCharacterEvent, c.Index()).Write("bonus", bonus)
		}

		atk.Info.BaseDmgBonus += bonus
	}, lunarbloomBonusKey)
}
