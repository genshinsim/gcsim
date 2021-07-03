package prayer

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("lost prayer to the sacred winds", weapon)
}

//Increases Movement Speed SPD by 10%. When in battle, earn a 6/8/10/12/14% Elemental DMG Bonus every 4s.
//Max 4 stacks. Lasts until the character falls or leaves combat.
func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	//ignore movement speed
	w := weap{}
	w.stacks = param["stack"]
	if w.stacks > 4 {
		w.stacks = 4
	}
	//check every 4 sec, if active add 1 stack;
	c.AddTask(w.stackCheck(c, s), "prayer-stack", 240)

	//remove stack on swap off
	s.AddEventHook(func(s def.Sim) bool {
		if s.ActiveCharIndex() != c.CharIndex() {
			w.stacks = 0
		}
		return false
	}, fmt.Sprintf("lostprayer-%v", c.Name()), def.PostSwapHook)

	dmg := 0.04 + float64(r)*0.02
	m := make([]float64, def.EndStatType)
	c.AddMod(def.CharStatMod{
		Key:    "lost-prayer",
		Expiry: -1,
		Amount: func(a def.AttackTag) ([]float64, bool) {
			if w.stacks == 0 {
				return nil, false
			}
			p := dmg * float64(w.stacks)
			m[def.PyroP] = p
			m[def.HydroP] = p
			m[def.CryoP] = p
			m[def.ElectroP] = p
			m[def.AnemoP] = p
			m[def.GeoP] = p
			m[def.EleP] = p
			m[def.PhyP] = p
			m[def.DendroP] = p
			return m, true
		},
	})

}

type weap struct {
	stacks int
}

func (w *weap) stackCheck(c def.Character, s def.Sim) func() {
	return func() {
		if s.ActiveCharIndex() == c.CharIndex() {
			w.stacks++
			if w.stacks > 4 {
				w.stacks = 4
			}
		}
		c.AddTask(w.stackCheck(c, s), "prayer-stack", 240)
	}
}
