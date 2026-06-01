package gestofthemightywolf

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.GestOfTheMightyWolf, NewWeapon)
}

type Weapon struct {
	Index int
}

const (
	normalAttackStacksKey = "gest-of-the-mighty-wolf-normal-stacks"
	chargeAttackStacksKey = "gest-of-the-mighty-wolf-charge-stacks"
	skillStacksKey        = "gest-of-the-mighty-wolf-skill-stacks"
)

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Increase Plunging Attack CRIT Rate by 16%;
// After a Plunging Attack hits an opponent, Normal, Charged, and Plunging Attack DMG increased by 16% for 10s.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.AtkSpd] = 0.1
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("gest-of-the-mighty-wolf-atkspd", -1),
		AffectedStat: attributes.AtkSpd,
		Amount: func() []float64 {
			return m
		},
	})

	n := make([]float64, attributes.EndStatType)
	char.AddStatMod(character.StatMod{
		Base: modifier.NewBase("gest-of-the-mighty-wolf-stacks", -1),
		Amount: func() []float64 {
			stacks := 0
			if char.StatusIsActive(normalAttackStacksKey) {
				stacks++
			}
			if char.StatusIsActive(chargeAttackStacksKey) {
				stacks += 2
			}
			if char.StatusIsActive(skillStacksKey) {
				stacks += 2
			}
			stacks = min(stacks, 4)

			n[attributes.DmgP] = (0.055 + 0.02*float64(r)) * float64(stacks)

			if c.Player.GetHexereiCount() >= 2 {
				n[attributes.CD] = (0.055 + 0.02*float64(r)) * float64(stacks)
			}

			return n
		},
	})

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}
		if atk.Info.AttackTag != attacks.AttackTagNormal {
			return
		}
		char.AddStatus(normalAttackStacksKey, 4*60, false)
	}, fmt.Sprintf("gest-of-the-mighty-wolf-on-normal-attack-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnChargeAttack, func(args ...any) {
		if c.Player.Active() != char.Index() {
			return
		}
		char.AddStatus(chargeAttackStacksKey, 4*60, false)
	}, fmt.Sprintf("gest-of-the-mighty-wolf-on-charge-attack-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnSkill, func(args ...any) {
		if c.Player.Active() != char.Index() {
			return
		}
		char.AddStatus(skillStacksKey, 4*60, false)
	}, fmt.Sprintf("gest-of-the-mighty-wolf-on-skill-%v", char.Base.Key.String()))

	return w, nil
}
