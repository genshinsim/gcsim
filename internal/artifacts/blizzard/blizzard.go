package blizzard

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterSetFunc("blizzard strayer", New)
}

func New(c def.Character, s def.Sim, log def.Logger, count int) {
	if count >= 2 {
		m := make([]float64, def.EndStatType)
		m[def.CryoP] = 0.15
		c.AddMod(def.CharStatMod{
			Key: "bs-2pc",
			Amount: func(a def.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		s.AddOnAttackWillLand(func(t def.Target, ds *def.Snapshot) {
			if ds.ActorIndex != c.CharIndex() {
				return
			}
			switch t.AuraType() {
			case def.Cryo:
				ds.Stats[def.CR] += .2
				log.Debugw("blizzard strayer 4pc on cryo", "frame", s.Frame(), "event", def.LogCalc, "char", c.CharIndex(), "new crit", ds.Stats[def.CR])
			case def.Frozen:
				ds.Stats[def.CR] += .4
				log.Debugw("blizzard strayer 4pc on frozen", "frame", s.Frame(), "event", def.LogCalc, "char", c.CharIndex(), "new crit", ds.Stats[def.CR])
			}
		}, fmt.Sprintf("bs4-%v", c.Name()))
	}
	//add flat stat to char
}
