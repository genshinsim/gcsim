package flowerwreathedfeathers

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	buffStatus = "flowerwreathedfeathers"
	icdStatus  = "flowerwreathedfeathers-icd"
)

func init() {
	core.RegisterWeaponFunc(keys.FlowerWreathedFeathers, NewWeapon)
}

type Weapon struct {
	Index int

	c        *core.Core
	char     *character.CharWrapper
	stacks   int
	leaveSrc int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Decreases Gliding Stamina consumption by 15%. When using Aimed Shots, the DMG dealt by Charged Attacks
// increases by 6% every 0.5s. This effect can stack up to 6 times and will be removed 10s after leaving
// Aiming Mode.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{
		c:    c,
		char: char,
	}
	r := p.Refine

	// "Gliding Stamina consumption" not implemented

	m := make([]float64, attributes.EndStatType)
	buff := 0.045 + 0.015*float64(r)
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(buffStatus, -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagExtra {
				return nil, false
			}
			m[attributes.DmgP] = buff * float64(w.stacks)
			return m, true
		},
	})

	c.Events.Subscribe(event.OnAimShoot, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}
		if char.StatusIsActive(icdStatus) {
			return false
		}
		char.AddStatus(icdStatus, 0.5*60, true)

		w.leaveSrc = -1
		if w.stacks < 6 {
			w.stacks++
		}
		c.Log.NewEvent("flower-wreathed feathers proc'd", glog.LogWeaponEvent, char.Index).
			Write("stacks", w.stacks)

		return false
	}, fmt.Sprintf("flower-wreathed-aim-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnStateChange, func(args ...interface{}) bool {
		prev := args[0].(action.AnimationState)
		next := args[1].(action.AnimationState)

		if c.Player.Active() != char.Index {
			return false
		}
		if prev != action.AimState || next == action.AimState {
			return false
		}
		if w.leaveSrc != -1 {
			return false
		}
		w.leaveSrc = c.F
		char.QueueCharTask(w.clearBuff(w.leaveSrc), 10*60)

		return false
	}, fmt.Sprintf("flower-wreathed-state-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		prev := args[0].(int)

		if prev != char.Index {
			return false
		}
		if w.leaveSrc != -1 {
			return false
		}
		w.leaveSrc = c.F
		char.QueueCharTask(w.clearBuff(w.leaveSrc), 10*60)

		return false
	}, fmt.Sprintf("flower-wreathed-swap-%v", char.Base.Key.String()))

	return w, nil
}

func (w *Weapon) clearBuff(src int) func() {
	return func() {
		if w.leaveSrc != src {
			return
		}
		if w.c.Player.Active() == w.char.Index && w.c.Player.CurrentState() == action.AimState {
			return
		}

		w.stacks = 0
		w.c.Log.NewEvent("flower-wreathed feathers cleared", glog.LogWeaponEvent, w.char.Index).
			Write("stacks", w.stacks)
	}
}
