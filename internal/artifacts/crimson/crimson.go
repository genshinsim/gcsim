package crimson

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterSetFunc("crimson witch of flames", New)
}

func New(c def.Character, s def.Sim, log def.Logger, count int) {
	stacks := 0
	cdTag := fmt.Sprintf("%v-cw-4pc", c.Name())
	if count >= 2 {
		m := make([]float64, def.EndStatType)
		c.AddMod(def.CharStatMod{
			Key: "crimson-2pc",
			Amount: func(a def.AttackTag) ([]float64, bool) {
				if s.Status(cdTag) == 0 {
					stacks = 0
				}
				mult := 0.5*float64(stacks) + 1
				m[def.PyroP] = 0.15 * mult
				if mult > 1 {
					log.Debugw("crimson witch 4pc", "frame", s.Frame(), "event", def.LogArtifactEvent, "char", c.CharIndex(), "mult", mult)
				}
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {

		//post snap shot to increase stacks
		s.AddEventHook(func(s def.Sim) bool {
			if s.ActiveCharIndex() != c.CharIndex() {
				return false
			}
			if s.Status(cdTag) == 0 {
				stacks = 0
			}
			//every exectuion, add 1 stack, to a max of 3, reset cd to 10 seconds
			stacks++
			if stacks > 3 {
				stacks = 3
			}
			log.Debugw("crimson witch 4pc adding stack", "frame", s.Frame(), "event", def.LogArtifactEvent, "current stacks", stacks)
			s.AddStatus(cdTag, 600)
			return false
		}, fmt.Sprintf("cw4s-%v", c.Name()), def.PostSkillHook)

		s.AddOnAmpReaction(func(t def.Target, ds *def.Snapshot) {
			if ds.ActorIndex != c.CharIndex() {
				return
			}
			switch ds.ReactionType {
			case def.Melt:
				ds.ReactBonus += 0.15
			case def.Vaporize:
				ds.ReactBonus += 0.15
			}
		}, cdTag)

		s.AddOnTransReaction(func(t def.Target, ds *def.Snapshot) {
			if ds.ActorIndex != c.CharIndex() {
				return
			}
			switch ds.ReactionType {
			case def.Overload:
				ds.ReactBonus += 0.4
			}
		}, cdTag)

	}
}
