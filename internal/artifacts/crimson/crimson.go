package crimson

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("crimson witch of flames", New)
	core.RegisterSetFunc("crimsonwitchofflames", New)
}

func New(c core.Character, s *core.Core, count int) {
	stacks := 0
	key := fmt.Sprintf("%v-cw-4pc", c.Name())
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		c.AddMod(core.CharStatMod{
			Key: "crimson-2pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				if s.Status.Duration(key) == 0 {
					stacks = 0
				}
				mult := 0.5*float64(stacks) + 1
				m[core.PyroP] = 0.15 * mult
				if mult > 1 {
					s.Log.Debugw("crimson witch 4pc", "frame", s.F, "event", core.LogArtifactEvent, "char", c.CharIndex(), "mult", mult)
				}
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {

		//post snap shot to increase stacks
		s.Events.Subscribe(core.PostSkill, func(args ...interface{}) bool {
			if s.ActiveChar != c.CharIndex() {
				return false
			}
			if s.Status.Duration(key) == 0 {
				stacks = 0
			}
			//every exectuion, add 1 stack, to a max of 3, reset cd to 10 seconds
			stacks++
			if stacks > 3 {
				stacks = 3
			}
			s.Log.Debugw("crimson witch 4pc adding stack", "frame", s.F, "event", core.LogArtifactEvent, "current stacks", stacks)
			s.Status.AddStatus(key, 600)
			return false
		}, fmt.Sprintf("cw4s-%v", c.Name()))

		c.AddReactBonusMod(core.ReactionBonusMod{
			Key:    "4cw",
			Expiry: -1,
			Amount: func(ai core.AttackInfo) (float64, bool) {
				//overload dmg can't melt or vape so it's fine
				if ai.AttackTag == core.AttackTagOverloadDamage {
					return 0.4, false
				}
				if ai.Amped {
					return 0.15, false
				}
				return 0, false
			},
		})

	}
}
