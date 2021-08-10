package crimson

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterSetFunc("crimson witch of flames", New)
}

func New(c core.Character, s core.Sim, log core.Logger, count int) {
	stacks := 0
	cdTag := fmt.Sprintf("%v-cw-4pc", c.Name())
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		c.AddMod(core.CharStatMod{
			Key: "crimson-2pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				if s.Status(cdTag) == 0 {
					stacks = 0
				}
				mult := 0.5*float64(stacks) + 1
				m[core.PyroP] = 0.15 * mult
				if mult > 1 {
					log.Debugw("crimson witch 4pc", "frame", s.Frame(), "event", core.LogArtifactEvent, "char", c.CharIndex(), "mult", mult)
				}
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {

		//post snap shot to increase stacks
		s.AddEventHook(func(s core.Sim) bool {
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
			log.Debugw("crimson witch 4pc adding stack", "frame", s.Frame(), "event", core.LogArtifactEvent, "current stacks", stacks)
			s.AddStatus(cdTag, 600)
			return false
		}, fmt.Sprintf("cw4s-%v", c.Name()), core.PostSkillHook)

		s.AddOnAmpReaction(func(t core.Target, ds *core.Snapshot) {
			if ds.ActorIndex != c.CharIndex() {
				return
			}
			switch ds.ReactionType {
			case core.Melt:
				ds.ReactBonus += 0.15
			case core.Vaporize:
				ds.ReactBonus += 0.15
			}
		}, cdTag)

		s.AddOnTransReaction(func(t core.Target, ds *core.Snapshot) {
			if ds.ActorIndex != c.CharIndex() {
				return
			}
			switch ds.ReactionType {
			case core.Overload:
				ds.ReactBonus += 0.4
			}
		}, cdTag)

	}
}
