package heartstrings

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
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

func init() {
	core.RegisterWeaponFunc(keys.SilvershowerHeartstrings, NewWeapon)
}

type Weapon struct {
	char       *character.CharWrapper
	core       *core.Core
	prevStacks int
	Index      int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

const (
	bondKey          = "heartstrings-bond"
	skillKey         = "heartstrings-skill"
	healingKey       = "heartstrings-healing"
	burstCRKey       = "heartstrings-cr"
	burstCRKeyCancel = "heartstrings-cr-cancel"
)

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	// The equipping character can gain the Remedy effect.
	// When they possess 1/2/3 Remedy stacks, Max HP will increase by 12%/24%/40%.
	// 1 stack may be gained when the following conditions are met:
	// 1 stack for 25s when using an Elemental Skill;
	// 1 stack for 25s when the value of a Bond of Life value increases;
	// 1 stack for 20s for performing healing.
	// Stacks can still be triggered when the equipping character is not on the field.
	// Each stack's duration is counted independently.
	// In addition, when 3 stacks are active, Elemental Burst CRIT Rate will be increased by 28%.
	// This effect will be canceled 4s after falling under 3 stacks.
	w := &Weapon{
		char: char,
		core: c,
	}
	r := p.Refine

	hpStack := 0.09 + float64(r)*0.03
	hpMaxStack := 0.03 + float64(r)*0.01

	// Max HP increase
	mHP := make([]float64, attributes.EndStatType)
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("heartstrings", -1),
		AffectedStat: attributes.HPP,
		Amount: func() ([]float64, bool) {
			stacks := w.Stacks()
			mHP[attributes.HPP] = hpStack * float64(stacks)
			if stacks >= 3 {
				mHP[attributes.HPP] += hpMaxStack
			}
			return mHP, true
		},
	})

	// Using skill
	c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}

		w.AddStack(skillKey, 25*60)
		return false
	}, fmt.Sprintf("heartstrings-%v", char.Base.Key.String()))

	// Gaining Bond
	c.Events.Subscribe(event.OnHPDebt, func(args ...interface{}) bool {
		index := args[0].(int)
		amount := args[1].(float64)

		if char.Index != index || amount <= 0 {
			return false
		}

		w.AddStack(bondKey, 25*60)
		return false
	}, fmt.Sprintf("heartstrings-%v", char.Base.Key.String()))

	// Healing
	c.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		src := args[0].(*info.HealInfo)

		if src.Caller != char.Index {
			return false
		}

		w.AddStack(healingKey, 20*60)
		return false
	}, fmt.Sprintf("heartstrings-%v", char.Base.Key.String()))

	// Burst CR buff if 3 stacks
	mCR := make([]float64, attributes.EndStatType)
	mCR[attributes.CR] = 0.21 + float64(r)*0.07
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(burstCRKey, -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagElementalBurst {
				return nil, false
			}
			if w.Stacks() < 3 && !char.StatusIsActive(burstCRKeyCancel) {
				return nil, false
			}
			return mCR, true
		},
	})

	return w, nil
}

func (w *Weapon) AddStack(name string, duration int) {
	w.char.AddStatus(name, duration, true)
	w.char.QueueCharTask(func() {
		if !w.char.StatusIsActive(name) {
			w.OnUpdateStack()
		}
	}, duration+1)
	w.OnUpdateStack()
}

func (w *Weapon) Stacks() int {
	count := 0
	if w.char.StatusIsActive(skillKey) {
		count++
	}
	if w.char.StatusIsActive(bondKey) {
		count++
	}
	if w.char.StatusIsActive(healingKey) {
		count++
	}
	return count
}

func (w *Weapon) OnUpdateStack() {
	stacks := w.Stacks()
	w.core.Log.NewEvent("heartstrings update stacks", glog.LogWeaponEvent, w.char.Index).
		Write("stacks", stacks).
		Write("bol-stack", w.char.StatusIsActive(bondKey)).
		Write("skill-stack", w.char.StatusIsActive(skillKey)).
		Write("heal-stack", w.char.StatusIsActive(healingKey))

	if w.prevStacks == 3 && stacks < 3 {
		// Elemental Burst CRIT Rate effect will be canceled 4s after falling under 3 stacks.
		w.char.AddStatus(burstCRKeyCancel, 4*60, true)
	}
	w.prevStacks = stacks
}
