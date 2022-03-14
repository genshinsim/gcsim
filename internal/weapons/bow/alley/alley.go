package alley

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("alley hunter", weapon)
	core.RegisterWeaponFunc("alleyhunter", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	//max 10 stacks
	w := weap{}
	w.stacks = param["stacks"]
	if w.stacks > 10 {
		w.stacks = 10
	}
	dmg := 0.015 + float64(r)*0.005

	m := make([]float64, core.EndStatType)
	char.AddPreDamageMod(coretype.PreDamageMod{
		Key:    "alley-hunter",
		Expiry: -1,
		Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
			m[core.DmgP] = dmg * float64(w.stacks)
			return m, true
		},
	})

	key := fmt.Sprintf("alley-hunter-%v", char.Name())

	c.Subscribe(core.OnInitialize, func(args ...interface{}) bool {
		w.active = c.ActiveChar == char.Index()
		w.lastActiveChange = c.Frame
		return true
	}, key)

	c.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		prev := args[0].(int)
		next := args[1].(int)
		w.lastActiveChange = c.Frame
		if next == char.Index() {
			w.active = true
			c.Tasks.Add(w.decStack(char, c.Frame), 240)
		} else if prev == char.Index() {
			w.active = false
			c.Tasks.Add(w.incStack(char, c.Frame), 60)
		}
		return false
	}, key)

	return "alleyhunter"
}

type weap struct {
	stacks           int
	active           bool
	lastActiveChange int
}

func (w *weap) decStack(c coretype.Character, src int) func() {
	return func() {
		if w.active && w.stacks > 0 && src == w.lastActiveChange {
			w.stacks -= 2
			if w.stacks < 0 {
				w.stacks = 0
			}
			c.AddTask(w.decStack(c, src), "alley-hunter", 60)
		}
	}
}

func (w *weap) incStack(c coretype.Character, src int) func() {
	return func() {
		if !w.active && w.stacks < 10 && src == w.lastActiveChange {
			w.stacks++
			c.AddTask(w.incStack(c, src), "alley-hunter", 60)
		}
	}
}
