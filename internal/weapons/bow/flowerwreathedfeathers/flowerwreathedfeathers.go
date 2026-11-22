package flowerwreathedfeathers

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	buffStatus = "flowerwreathedfeathers"
	// one stack is gained every 0.5s spent charging (i.e. in the aiming state)
	stackInterval = 0.5 * 60
	maxStacks     = 6
)

type Weapon struct {
	Index int

	c        *core.Core
	char     *character.CharWrapper
	stacks   int
	leaveSrc int
	stackSrc int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Decreases Gliding Stamina consumption by 15%. When using Aimed Shots, the DMG dealt by Charged
// Attacks increases by 6% every 0.5s. This effect can stack up to 6 times and will be removed 10s
// after leaving Aiming Mode.
//
// Implementation: each aim action starts a charging window; one stack is gained for every 0.5s
// (30 frames) of charge time. Charging stops the moment the aimed shot is released (event.OnAimRelease,
// e.g. Chasca's Shadowhunt volley release / a standard bow shot leaving the bow at its snapshot frame),
// so a whole volley snapshots the stack count it was charged with. The stacks are removed 10s after
// the character leaves Aiming Mode, unless they re-enter it in the meantime.
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
		Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
			if atk.Info.AttackTag != attacks.AttackTagExtra {
				return nil
			}
			m[attributes.DmgP] = buff * float64(w.stacks)
			return m
		},
	})

	// a new aim action starts a new charging window
	c.Events.Subscribe(event.OnAimShoot, func(args ...any) {
		if c.Player.Active() != char.Index() {
			return
		}
		// cancel any pending buff removal and start stacking
		w.leaveSrc = -1
		w.stackSrc = c.F
		char.QueueCharTask(w.stackTick(w.stackSrc), stackInterval)
	}, fmt.Sprintf("flower-wreathed-aim-%v", char.Base.Key.String()))

	// the shot was released: charging is over, stop stacking
	c.Events.Subscribe(event.OnAimRelease, func(args ...any) {
		if len(args) < 1 {
			return
		}
		idx, ok := args[0].(int)
		if !ok || idx != char.Index() {
			return
		}
		w.stackSrc++
	}, fmt.Sprintf("flower-wreathed-release-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnStateChange, func(args ...any) {
		prev := args[0].(action.AnimationState)
		next := args[1].(action.AnimationState)

		if c.Player.Active() != char.Index() {
			return
		}
		if prev != action.AimState || next == action.AimState {
			return
		}
		if w.leaveSrc != -1 {
			return
		}
		// left aiming mode: remove the buff after 10s
		w.leaveSrc = c.F
		char.QueueCharTask(w.clearBuff(w.leaveSrc), 10*60)
	}, fmt.Sprintf("flower-wreathed-state-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnCharacterSwap, func(args ...any) {
		prev := args[0].(int)

		if prev != char.Index() {
			return
		}
		if w.leaveSrc != -1 {
			return
		}
		w.leaveSrc = c.F
		char.QueueCharTask(w.clearBuff(w.leaveSrc), 10*60)
	}, fmt.Sprintf("flower-wreathed-swap-%v", char.Base.Key.String()))

	return w, nil
}

// stackTick gains one stack every 0.5s while the charge is ongoing (character on field, in the
// aiming state, and the shot has not been released yet). The src guard kills stale tick chains
// when a new aim starts charging or when the shot is released.
func (w *Weapon) stackTick(src int) func() {
	return func() {
		if w.stackSrc != src {
			return
		}
		if w.c.Player.Active() != w.char.Index() {
			return
		}
		if w.c.Player.CurrentState() != action.AimState {
			return
		}
		if w.stacks < maxStacks {
			w.stacks++
			w.c.Log.NewEvent("flower-wreathed feathers stack gained", glog.LogWeaponEvent, w.char.Index()).
				Write("stacks", w.stacks)
		}
		w.char.QueueCharTask(w.stackTick(src), stackInterval)
	}
}

func (w *Weapon) clearBuff(src int) func() {
	return func() {
		if w.leaveSrc != src {
			return
		}
		if w.c.Player.Active() == w.char.Index() && w.c.Player.CurrentState() == action.AimState {
			return
		}

		w.stacks = 0
		w.c.Log.NewEvent("flower-wreathed feathers cleared", glog.LogWeaponEvent, w.char.Index()).
			Write("stacks", w.stacks)
	}
}
