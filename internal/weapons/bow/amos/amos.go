package amos

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("amos' bow", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {

	dmgpers := 0.06 + 0.02*float64(r)

	m := make([]float64, def.EndStatType)
	m[def.DmgP] = 0.09 + 0.03*float64(r)
	c.AddMod(def.CharStatMod{
		Key: "amos",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			return m, a == def.AttackTagNormal || a == def.AttackTagExtra
		},
		Expiry: -1,
	})

	s.AddOnAttackWillLand(func(t def.Target, ds *def.Snapshot) {
		if c.CharIndex() != ds.ActorIndex {
			return
		}
		if ds.AttackTag != def.AttackTagNormal && ds.AttackTag != def.AttackTagExtra {
			return
		}
		//calculate travel time
		travel := float64(s.Frame()-ds.SourceFrame-ds.AnimationFrames) / 60
		stacks := int(travel / 0.1)
		if stacks > 5 {
			stacks = 5
		}
		ds.Stats[def.DmgP] += dmgpers * float64(stacks)
		log.Debugw("amos bow", "frame", s.Frame(), "event", def.LogCalc, "stacks", stacks, "final dmg%", ds.Stats[def.DmgP])

	}, fmt.Sprintf("amos-%v", c.Name()))

}
