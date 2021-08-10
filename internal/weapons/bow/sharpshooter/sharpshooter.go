package sharpshooter

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("sharpshooter's oath", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
	dmg := 0.18 + float64(r)*0.06
	s.AddOnAttackWillLand(func(t core.Target, ds *core.Snapshot) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.HitWeakPoint {
			ds.Stats[core.DmgP] += dmg
			log.Debugw("sharpshooter", "frame", s.Frame(), "event", core.LogWeaponEvent, "final dmg%", ds.Stats[core.DmgP])
		}
	}, fmt.Sprintf("sharpshooter-%v", c.Name()))
}
