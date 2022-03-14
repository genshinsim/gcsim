package thunderingfury

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterSetFunc("thundering fury", New)
	core.RegisterSetFunc("thunderingfury", New)
}

func New(c coretype.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.ElectroP] = 0.15
		c.AddMod(coretype.CharStatMod{
			Key: "tf-2pc",
			Amount: func() ([]float64, bool) {
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
			atk := args[1].(*coretype.AttackEvent)
			if atk.Info.ActorIndex != c.Index() {
				return false
			}
			if s.Player.ActiveChar != c.Index()() {
				return false
			}
			if icd > s.Frame {
				return false
			}
			icd = s.Frame + 48
			c.ReduceActionCooldown(core.ActionSkill, 60)
			s.Log.NewEvent("thunderfury 4pc proc", coretype.LogArtifactEvent, c.Index(), "reaction", atk.Info.Abil, "new cd", c.Cooldown(core.ActionSkill))
			return false
		}

		s.Subscribe(core.OnOverload, reduce, "4tf"+c.Name())
		s.Subscribe(core.OnElectroCharged, reduce, "4tf"+c.Name())
		s.Subscribe(core.OnSuperconduct, reduce, "4tf"+c.Name())
	}
	//add flat stat to char
}
