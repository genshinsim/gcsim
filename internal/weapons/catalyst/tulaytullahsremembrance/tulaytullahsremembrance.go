package tulaytullahsremembrance

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
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
	core.RegisterWeaponFunc(keys.TulaytullahsRemembrance, NewWeapon)
}

type Weapon struct {
	Index  int
	stacks int
	src    int
	core   *core.Core
}

const (
	icdKey    = "tulaytullahsremembrance-icd"
	atkSpdKey = "tulaytullahsremembrance-atkspd"
	buffKey   = "tulaytullahsremembrance-na-dmg"
)

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Normal Attack SPD is increased by 10/12.5/15/17.5/20%.
// After the wielder unleashes an Elemental Skill, Normal Attack DMG will increase by 4.8/6/7.2/8.4/9.6% every second for 14s.
// After this character hits an opponent with a Normal Attack during this duration, Normal Attack DMG will be increased by 9.6/12/14.4/16.8/19.2%.
// This increase can be triggered once every 0.3s. The maximum Normal Attack DMG increase per single duration of the overall effect is 48/60/72/84/96%.
// The effect will be removed when the wielder leaves the field, and using the Elemental Skill again will reset all DMG buffs.
func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{
		core: c,
	}
	r := p.Refine

	// attack speed part
	mAtkSpd := make([]float64, attributes.EndStatType)
	mAtkSpd[attributes.AtkSpd] = 0.075 + float64(r)*0.025
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(atkSpdKey, -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			if c.Player.CurrentState() != action.NormalAttackState {
				return nil, false
			}
			return mAtkSpd, true
		},
	})

	// normal attack dmg part
	incDmg := 0.036 + float64(r)*0.012
	mDmg := make([]float64, attributes.EndStatType)
	c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}

		// remove stacks on skill in any case
		w.stacks = 0

		// gain 1 stack every 1s for 14s after using skill
		// no need to check for the 14s part, because it will stop when it reaches max stacks anyways
		w.src = c.F
		char.QueueCharTask(w.incStack(char, c.F), 60)

		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag(buffKey, 14*60),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if atk.Info.AttackTag != combat.AttackTagNormal {
					return nil, false
				}
				mDmg[attributes.DmgP] = incDmg * float64(w.stacks)
				return mDmg, true
			},
		})
		return false
	}, fmt.Sprintf("tulaytullahsremembrance-%v", char.Base.Key.String()))

	// gain 2 stacks on normal attack dmg, 0.3s icd
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if atk.Info.AttackTag != combat.AttackTagNormal {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		char.AddStatus(icdKey, 0.3*60, true)

		previous := w.stacks
		w.stacks += 2
		if w.stacks > 10 {
			w.stacks = 10
		}
		gain := w.stacks - previous
		if gain == 0 {
			return false
		}
		gainMsg := "2 stacks"
		if gain == 1 {
			gainMsg = "1 stack"
		}
		w.core.Log.NewEvent(fmt.Sprintf("Tulaytullah's Remembrance gained %v via normal attack", gainMsg), glog.LogWeaponEvent, char.Index).
			Write("stacks", w.stacks)
		return false
	}, fmt.Sprintf("tulaytullahsremembrance-ondmg-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		prev := args[0].(int)
		if prev != char.Index {
			return false
		}
		if !char.StatusIsActive(buffKey) {
			return false
		}
		// remove stacks, invalidate incStack task and remove buff on swap
		w.stacks = 0
		w.src = -1
		char.DeleteStatus(buffKey)
		return false
	}, fmt.Sprintf("tulaytullahsremembrance-exit-%v", char.Base.Key.String()))

	return w, nil
}

func (w *Weapon) incStack(char *character.CharWrapper, src int) func() {
	return func() {
		if w.stacks > 9 {
			return
		}
		if src != w.src {
			return
		}
		w.stacks++
		w.core.Log.NewEvent("Tulaytullah's Remembrance gained stack via timer", glog.LogWeaponEvent, char.Index).
			Write("stacks", w.stacks).
			Write("weapon effect start", w.src).
			Write("source", src)
		char.QueueCharTask(w.incStack(char, src), 60)
	}
}
