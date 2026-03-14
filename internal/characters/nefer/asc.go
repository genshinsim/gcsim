package nefer

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const lunarbloomBonusKey = "nefer-lunarbloom-bonus"

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
