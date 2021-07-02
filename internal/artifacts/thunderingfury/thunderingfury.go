package thunderingfury

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterSetFunc("thundering fury", New)
}

func New(c def.Character, s def.Sim, log def.Logger, count int) {
	if count >= 2 {
		m := make([]float64, def.EndStatType)
		m[def.ElectroP] = 0.15
		c.AddMod(def.CharStatMod{
			Key: "tf-2pc",
			Amount: func(a def.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		icd := 0

		s.AddOnTransReaction(func(t def.Target, ds *def.Snapshot) {
			if ds.ActorIndex != c.CharIndex() {
				return
			}
			if icd > s.Frame() {
				return
			}
			switch ds.ReactionType {
			case def.Overload:
			case def.ElectroCharged:
			case def.Superconduct:
			default:
				return
			}
			ds.ReactBonus += 0.4
			icd = s.Frame() + 48
			c.ReduceActionCooldown(def.ActionSkill, 60)
			log.Debugw("thunderfury 4pc proc", "frame", s.Frame(), "event", def.LogArtifactEvent, "reaction", ds.ReactionType, "new cd", c.Cooldown(def.ActionSkill))

		}, fmt.Sprintf("4tf-%v", c.Name()))

	}
	//add flat stat to char
}
