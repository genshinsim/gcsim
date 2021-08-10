package prayer

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("lost prayer to the sacred winds", weapon)
}

//Increases Movement Speed SPD by 10%. When in battle, earn a 6/8/10/12/14% Elemental DMG Bonus every 4s.
//Max 4 stacks. Lasts until the character falls or leaves combat.
func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
	//ignore movement speed
	w := weap{}
	w.stacks = param["stack"]
	if w.stacks > 4 {
		w.stacks = 4
	}
	//check every 4 sec, if active add 1 stack;
	c.AddTask(w.stackCheck(c, s), "prayer-stack", 240)

	//remove stack on swap off
	s.AddEventHook(func(s core.Sim) bool {
		if s.ActiveCharIndex() != c.CharIndex() {
			w.stacks = 0
		}
		return false
	}, fmt.Sprintf("lostprayer-%v", c.Name()), core.PostSwapHook)

	dmg := 0.04 + float64(r)*0.02
	m := make([]float64, core.EndStatType)
	c.AddMod(core.CharStatMod{
		Key:    "lost-prayer",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if w.stacks == 0 {
				return nil, false
			}
			p := dmg * float64(w.stacks)
			m[core.PyroP] = p
			m[core.HydroP] = p
			m[core.CryoP] = p
			m[core.ElectroP] = p
			m[core.AnemoP] = p
			m[core.GeoP] = p
			m[core.EleP] = p
			m[core.PhyP] = p
			m[core.DendroP] = p
			return m, true
		},
	})

}

type weap struct {
	stacks int
}

func (w *weap) stackCheck(c core.Character, s core.Sim) func() {
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
