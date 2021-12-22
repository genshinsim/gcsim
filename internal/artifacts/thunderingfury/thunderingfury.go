package thunderingfury

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("thundering fury", New)
	core.RegisterSetFunc("thunderingfury", New)
}

func New(c core.Character, s *core.Core, count int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.ElectroP] = 0.15
		c.AddMod(core.CharStatMod{
			Key: "tf-2pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		icd := 0

		//add +0.4 reaction damage
		c.AddReactBonusMod(core.ReactionBonusMod{
			Key:    "4tf",
			Expiry: -1,
			Amount: func(ai core.AttackInfo) (float64, bool) {
				//overload dmg can't melt or vape so it's fine
				switch ai.AttackTag {
				case core.AttackTagOverloadDamage:
				case core.AttackTagECDamage:
				case core.AttackTagSuperconductDamage:
				default:
					return 0, false
				}
				return 0.4, false
			},
		})

		reduce := func(args ...interface{}) bool {
			atk := args[1].(*core.AttackEvent)
			if atk.Info.ActorIndex != c.CharIndex() {
				return false
			}
			if icd > s.F {
				return false
			}
			icd = s.F + 48
			c.ReduceActionCooldown(core.ActionSkill, 60)
			s.Log.Debugw("thunderfury 4pc proc", "frame", s.F, "event", core.LogArtifactEvent, "reaction", atk.Info.Abil, "new cd", c.Cooldown(core.ActionSkill))
			return false
		}

		s.Events.Subscribe(core.OnOverload, reduce, "4tf"+c.Name())
		s.Events.Subscribe(core.OnElectroCharged, reduce, "4tf"+c.Name())
		s.Events.Subscribe(core.OnSuperconduct, reduce, "4tf"+c.Name())
	}
	//add flat stat to char
}
