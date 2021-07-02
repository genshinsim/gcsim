package sharpshooter

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("sharpshooter's oath", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	dmg := 0.18 + float64(r)*0.06
	s.AddOnAttackWillLand(func(t def.Target, ds *def.Snapshot) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.HitWeakPoint {
			ds.Stats[def.DmgP] += dmg
			log.Debugw("sharpshooter", "frame", s.Frame(), "event", def.LogWeaponEvent, "final dmg%", ds.Stats[def.DmgP])
		}
	}, fmt.Sprintf("sharpshooter-%v", c.Name()))
}
