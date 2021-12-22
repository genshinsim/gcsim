package prayer

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("lost prayer to the sacred winds", weapon)
	core.RegisterWeaponFunc("lostprayertothesacredwinds", weapon)
}

//Increases Movement Speed SPD by 10%. When in battle, earn a 6/8/10/12/14% Elemental DMG Bonus every 4s.
//Max 4 stacks. Lasts until the character falls or leaves combat.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	//ignore movement speed
	w := weap{}
	w.stacks = param["stack"]
	if w.stacks > 4 {
		w.stacks = 4
	}
	//check every 4 sec, if active add 1 stack;
	char.AddTask(w.stackCheck(char, c), "prayer-stack", 240)

	//remove stack on swap off
	c.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		if c.ActiveChar != char.CharIndex() {
			w.stacks = 0
		}
		return false
	}, fmt.Sprintf("lostprayer-%v", char.Name()))

	dmg := 0.04 + float64(r)*0.02
	char.AddMod(core.CharStatMod{
		Key:    "lost-prayer",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			m := make([]float64, core.EndStatType)
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
			m[core.DendroP] = p
			return m, true
		},
	})

}

type weap struct {
	stacks int
}

func (w *weap) stackCheck(char core.Character, c *core.Core) func() {
	return func() {
		if c.ActiveChar == char.CharIndex() {
			w.stacks++
			if w.stacks > 4 {
				w.stacks = 4
			}
		}
		char.AddTask(w.stackCheck(char, c), "prayer-stack", 240)
	}
}
