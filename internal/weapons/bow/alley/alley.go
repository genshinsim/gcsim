package alley

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.AlleyHunter, NewWeapon)
}

type Weapon struct {
	stacks           int
	active           bool
	lastActiveChange int
	Index            int
	core             *core.Core
	char             *character.CharWrapper
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }

// Initiate off-field stacking if off-field at start of the sim
func (w *Weapon) Init() error {
	w.active = w.core.Player.Active() == w.char.Index
	if !w.active {
		w.core.Tasks.Add(w.incStack(w.char, w.core.F), 1)
	}
	w.lastActiveChange = w.core.F
	return nil
}

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	//While the character equipped with this weapon is in the party but not on the field, their DMG
	//increases by 2% every second up to a max of 20%. When the character is on the field for more than 4s,
	//the aforementioned DMG buff decreases by 4% per second until it reaches 0%.
	r := p.Refine

	//max 10 stacks
	w := Weapon{
		core: c,
		char: char,
	}
	w.stacks = p.Params["stacks"]
	if w.stacks > 10 {
		w.stacks = 10
	}
	dmg := 0.015 + float64(r)*0.005

	m := make([]float64, attributes.EndStatType)
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("alley-hunter", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			m[attributes.DmgP] = dmg * float64(w.stacks)
			return m, true
		},
	})

	key := fmt.Sprintf("alley-hunter-%v", char.Base.Key.String())

	c.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		prev := args[0].(int)
		next := args[1].(int)
		w.lastActiveChange = c.F
		if next == char.Index {
			w.active = true
			w.char.QueueCharTask(w.decStack(char, c.F), 240) // on field for more than 4s, start decreasing stacks
		} else if prev == char.Index {
			w.active = false
			c.Tasks.Add(w.incStack(char, c.F), 60)
		}
		return false
	}, key)

	return &w, nil
}

func (w *Weapon) decStack(c *character.CharWrapper, src int) func() {
	return func() {
		if w.active && w.stacks > 0 && src == w.lastActiveChange {
			w.stacks -= 2
			if w.stacks < 0 {
				w.stacks = 0
			}
			w.core.Log.NewEvent("Alley lost stack", glog.LogWeaponEvent, w.char.Index).
				Write("stacks:", w.stacks).
				Write("last_swap", w.lastActiveChange).
				Write("source", src)
			w.char.QueueCharTask(w.decStack(c, src), 60)
		}
	}
}

func (w *Weapon) incStack(c *character.CharWrapper, src int) func() {
	return func() {
		if !w.active && w.stacks < 10 && src == w.lastActiveChange {
			w.stacks++
			w.core.Log.NewEvent("Alley gained stack", glog.LogWeaponEvent, w.char.Index).
				Write("stacks:", w.stacks).
				Write("last_swap", w.lastActiveChange).
				Write("source", src)
			w.core.Tasks.Add(w.incStack(c, src), 60)
		}
	}
}
