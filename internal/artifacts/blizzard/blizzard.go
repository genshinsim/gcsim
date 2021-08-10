package blizzard

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterSetFunc("blizzard strayer", New)
}

func New(c core.Character, s core.Sim, log core.Logger, count int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.CryoP] = 0.15
		c.AddMod(core.CharStatMod{
			Key: "bs-2pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		s.AddOnAttackWillLand(func(t core.Target, ds *core.Snapshot) {
			if ds.ActorIndex != c.CharIndex() {
				return
			}
			switch t.AuraType() {
			case core.Cryo:
				ds.Stats[core.CR] += .2
				log.Debugw("blizzard strayer 4pc on cryo", "frame", s.Frame(), "event", core.LogCalc, "char", c.CharIndex(), "new crit", ds.Stats[core.CR])
			case core.Frozen:
				ds.Stats[core.CR] += .4
				log.Debugw("blizzard strayer 4pc on frozen", "frame", s.Frame(), "event", core.LogCalc, "char", c.CharIndex(), "new crit", ds.Stats[core.CR])
			}
		}, fmt.Sprintf("bs4-%v", c.Name()))
	}
	//add flat stat to char
}
