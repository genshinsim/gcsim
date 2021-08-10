package thunderingfury

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterSetFunc("thundering fury", New)
}

func New(c core.Character, s core.Sim, log core.Logger, count int) {
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

		s.AddOnTransReaction(func(t core.Target, ds *core.Snapshot) {
			if ds.ActorIndex != c.CharIndex() {
				return
			}
			if icd > s.Frame() {
				return
			}
			switch ds.ReactionType {
			case core.Overload:
			case core.ElectroCharged:
			case core.Superconduct:
			default:
				return
			}
			ds.ReactBonus += 0.4
			icd = s.Frame() + 48
			c.ReduceActionCooldown(core.ActionSkill, 60)
			log.Debugw("thunderfury 4pc proc", "frame", s.Frame(), "event", core.LogArtifactEvent, "reaction", ds.ReactionType, "new cd", c.Cooldown(core.ActionSkill))

		}, fmt.Sprintf("4tf-%v", c.Name()))

	}
	//add flat stat to char
}
