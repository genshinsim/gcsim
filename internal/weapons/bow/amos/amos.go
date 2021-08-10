package amos

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("amos' bow", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {

	dmgpers := 0.06 + 0.02*float64(r)

	m := make([]float64, core.EndStatType)
	m[core.DmgP] = 0.09 + 0.03*float64(r)
	c.AddMod(core.CharStatMod{
		Key: "amos",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, a == core.AttackTagNormal || a == core.AttackTagExtra
		},
		Expiry: -1,
	})

	s.AddOnAttackWillLand(func(t core.Target, ds *core.Snapshot) {
		if c.CharIndex() != ds.ActorIndex {
			return
		}
		if ds.AttackTag != core.AttackTagNormal && ds.AttackTag != core.AttackTagExtra {
			return
		}
		//calculate travel time
		travel := float64(s.Frame()-ds.SourceFrame-ds.AnimationFrames) / 60
		stacks := int(travel / 0.1)
		if stacks > 5 {
			stacks = 5
		}
		ds.Stats[core.DmgP] += dmgpers * float64(stacks)
		log.Debugw("amos bow", "frame", s.Frame(), "event", core.LogCalc, "stacks", stacks, "final dmg%", ds.Stats[core.DmgP])

	}, fmt.Sprintf("amos-%v", c.Name()))

}
