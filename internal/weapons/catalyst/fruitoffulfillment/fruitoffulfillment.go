package fruitoffulfillment

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.FruitOfFulfillment, NewWeapon)
}

type Weapon struct {
	Index  int
	core   *core.Core
	char   *character.CharWrapper
	stacks int
	// Required to check for stack loss
	stackLossTimer int
	lastStackGain  int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Obtain the "Wax and Wane" effect after an Elemental Reaction is triggered, gaining 24/27/30/33/36 Elemental Mastery while losing 5% ATK.
// For every 0.3s, 1 stack of Wax and Wane can be gained. Max 5 stacks.
// For every 6s that go by without an Elemental Reaction being triggered, 1 stack will be lost.
// This effect can be triggered even when the character is off-field.
func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{
		core: c,
		char: char,
	}
	r := p.Refine

	em := 21 + float64(r)*3
	atkLoss := -0.05

	w.stackLossTimer = 360 // 6s * 60

	const buffKey = "fruitoffulfillment"
	const icdKey = "fruitoffulfillment-icd"

	m := make([]float64, attributes.EndStatType)
	w.char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(buffKey, -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			m[attributes.EM] = em * float64(w.stacks)
			m[attributes.ATKP] = atkLoss * float64(w.stacks)
			return m, true
		},
	})

	f := func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != w.char.Index {
			return false
		}
		if w.char.StatusIsActive(icdKey) {
			return false
		}
		w.char.AddStatus(icdKey, 18, true)

		w.stacks++
		if w.stacks > 5 {
			w.stacks = 5
		}

		w.lastStackGain = c.F
		w.char.QueueCharTask(w.checkStackLoss(c.F), w.stackLossTimer)

		w.core.Log.NewEvent("fruitoffulfillment gained stack", glog.LogWeaponEvent, w.char.Index).
			Write("stacks", w.stacks)

		return false
	}

	for i := event.Event(event.ReactionEventStartDelim + 1); i < event.ReactionEventEndDelim; i++ {
		w.core.Events.Subscribe(i, f, fmt.Sprintf("fruitoffulfillment-%v", w.char.Base.Key.String()))
	}

	return w, nil
}

// Helper function to check for stack loss
// called after every stack gain
func (w *Weapon) checkStackLoss(src int) func() {
	return func() {
		if w.lastStackGain != src {
			w.core.Log.NewEvent("fruitoffulfillment stack loss check ignored, src diff", glog.LogWeaponEvent, w.char.Index).
				Write("src", src).
				Write("new src", w.lastStackGain)
			return
		}
		w.stacks--
		w.core.Log.NewEvent("fruitoffulfillment lost stack", glog.LogWeaponEvent, w.char.Index).
			Write("stacks", w.stacks).
			Write("last_stack_change", w.lastStackGain)

		// queue up again if we still have stacks
		if w.stacks > 0 {
			w.char.QueueCharTask(w.checkStackLoss(src), w.stackLossTimer)
		}
	}
}
