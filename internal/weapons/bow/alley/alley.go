package alley

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("alley hunter", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	//max 10 stacks
	w := weap{}
	w.stacks = param["stack"]
	if w.stacks > 10 {
		w.stacks = 10
	}
	dmg := 0.015 + float64(r)*0.005

	m := make([]float64, def.EndStatType)
	c.AddMod(def.CharStatMod{
		Key: "alley-hunter",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			m[def.DmgP] = dmg * float64(w.stacks)
			return m, true
		},
		Expiry: -1,
	})

	s.AddInitHook(func() {
		w.active = s.ActiveCharIndex() == c.CharIndex()
	})

	s.AddEventHook(func(s def.Sim) bool {
		//if swapped in
		if s.ActiveCharIndex() == c.CharIndex() {
			w.active = true
			c.AddTask(w.decStack(c), "alley-hunter", 240) //start losing every 1 sec at 4 sec
		} else {
			w.active = false
			c.AddTask(w.incStack(c), "alley-hunter", 60)
		}
		return false
	}, fmt.Sprintf("alley-hunter-%v", c.Name()), def.PostSwapHook)

}

type weap struct {
	stacks int
	active bool
}

func (w *weap) decStack(c def.Character) func() {
	return func() {
		if w.active && w.stacks > 0 {
			w.stacks -= 2
			if w.stacks < 0 {
				w.stacks = 0
			}
			c.AddTask(w.decStack(c), "alley-hunter", 60)
		}
	}
}

func (w *weap) incStack(c def.Character) func() {
	return func() {
		if !w.active && w.stacks < 10 {
			w.stacks++
			c.AddTask(w.incStack(c), "alley-hunter", 60)
		}
	}
}
