package heartstrings

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.SilvershowerHeartstrings, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

const (
	bondKey    = "heartstrings-bond"
	skillKey   = "heartstrings-skill"
	healingKey = "heartstrings-healing"
)

// TODO: clear burst cr buff 4s after losing 3 stacks <-- I have no idea how to do this
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
	w := &Weapon{}
	r := p.Refine

	stack := 0.09 + float64(r)*0.03
	max := 0.04 + float64(r) - 1

	mHP := make([]float64, attributes.EndStatType)

	// Using skill
	c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}

		char.AddStatus(skillKey, 25*60, true)
		return false
	}, fmt.Sprintf("heartstrings-%v", char.Base.Key.String()))

	// Gaining Bond
	c.Events.Subscribe(event.OnHPDebt, func(args ...interface{}) bool {
		index := args[0].(int)
		amount := args[1].(float64)

		if char.Index != index || amount <= 0 {
			return false
		}

		char.AddStatus(bondKey, 25*60, true)
		return false
	}, fmt.Sprintf("heartstrings-%v", char.Base.Key.String()))

	// Healing
	c.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		src := args[0].(*info.HealInfo)

		if src.Caller != char.Index {
			return false
		}

		char.AddStatus(healingKey, 20*60, true)
		return false
	}, fmt.Sprintf("heartstrings-%v", char.Base.Key.String()))

	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("heartstrings", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			count := 0
			if char.StatusIsActive(skillKey) {
				count++
			}
			if char.StatusIsActive(bondKey) {
				count++
			}
			if char.StatusIsActive(healingKey) {
				count++
			}
			maxhpbonus := stack * float64(count)
			if count >= 3 {
				maxhpbonus += max

				// should this be here
				// Burst CR buff if 3 stacks
				mCR := make([]float64, attributes.EndStatType)
				mCR[attributes.CR] = 0.21 + float64(r)*0.07
				char.AddAttackMod(character.AttackMod{
					Base: modifier.NewBase("heartstrings-cr", -1),
					Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
						if atk.Info.AttackTag != attacks.AttackTagElementalBurst {
							return nil, false
						}
						return mCR, true
					},
				})
			}
			mHP[attributes.HPP] = maxhpbonus
			return mHP, true
		},
	})
	return w, nil
}
