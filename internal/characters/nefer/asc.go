package nefer

import (
	"math"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const (
	lunarbloomBonusKey = "nefer-lunarbloom-bonus"
	p2DewWindowKey     = "nefer-p2-dew-window"
	p2DewRateKey       = "nefer-p2-dew-rate"
	p2DewWindowDur     = 5 * 60
)

func (c *char) p2Active() bool {
	return c.Base.Ascension >= 4
}

func (c *char) p2DewRateBonus() float64 {
	bonusSteps := math.Floor(max(c.Stat(attributes.EM)-500, 0) / 100)
	return min(bonusSteps*0.1, 0.5)
}

func (c *char) p2Init() {
	if !c.p2Active() {
		return
	}

	c.Core.Player.AddVerdantDewRateMod(p2DewRateKey, -1, func() (float64, bool) {
		if !c.StatusIsActive(shadowDanceKey) || !c.slitherActive() || !c.StatusIsActive(p2DewWindowKey) {
			return 0, false
		}
		return c.p2DewRateBonus(), false
	})

	c.Core.Events.Subscribe(event.OnLunarBloom, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex < 0 || atk.Info.ActorIndex >= len(c.Core.Player.Chars()) {
			return
		}
		c.AddStatus(p2DewWindowKey, p2DewWindowDur, true)
	}, p2DewWindowKey)
}

func (c *char) lunarbloomInit() {
	c.Core.Flags.Custom[reactable.LunarBloomEnableKey] = 1

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.AttackTag != attacks.AttackTagDirectLunarBloom {
			return
		}
		atk.Info.BaseDmgBonus += min(c.Stat(attributes.EM)*0.000175, 0.14)
	}, lunarbloomBonusKey)
}
