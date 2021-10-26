package thunderingfury

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
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

		s.Events.Subscribe(core.OnTransReaction, func(args ...interface{}) bool {
			ds := args[1].(*core.Snapshot)
			if ds.ActorIndex != c.CharIndex() {
				return false
			}
			switch ds.ReactionType {
			case core.Overload:
			case core.ElectroCharged:
			case core.Superconduct:
			default:
				return false
			}
			//react bonus should always apply
			ds.ReactBonus += 0.4
			//icd only applies to cd reduction
			if icd > s.F {
				return false
			}
			icd = s.F + 48
			c.ReduceActionCooldown(core.ActionSkill, 60)
			s.Log.Debugw("thunderfury 4pc proc", "frame", s.F, "event", core.LogArtifactEvent, "reaction", ds.ReactionType, "new cd", c.Cooldown(core.ActionSkill))
			return false
		}, fmt.Sprintf("4tf-%v", c.Name()))
	}
	//add flat stat to char
}
