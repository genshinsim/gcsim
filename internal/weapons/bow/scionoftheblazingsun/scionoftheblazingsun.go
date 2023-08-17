package scionoftheblazingsun

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.ScionOfTheBlazingSun, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// After a Charged Attack hits an opponent, a Sunfire Arrow will descend upon the opponent hit, dealing 60% ATK as DMG,
// and applying the Heartsearer effect to the opponent damaged by said Arrow for 10s. Opponents affected by Heartsearer
// take 28% more Charged Attack DMG from the wielder. A Sunfire Arrow can be triggered once every 10s.
func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	sunfireMult := 0.45 + float64(r)*0.15

	const icdKey = "scion-icd"
	const debuffKey = "scion-heartsearer"

	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.21 + 0.07*float64(r)
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("scion", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			e, ok := t.(*enemy.Enemy)
			if !ok {
				return nil, false
			}
			if !e.StatusIsActive(debuffKey) {
				return nil, false
			}
			if atk.Info.AttackTag != attacks.AttackTagExtra {
				return nil, false
			}

			return m, true
		},
	})

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagExtra {
			return false
		}

		ai := combat.AttackInfo{
			ActorIndex: char.Index,
			Abil:       "Sunfire Arrow",
			AttackTag:  attacks.AttackTagWeaponSkill,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Physical,
			Durability: 100,
			Mult:       sunfireMult,
		}
		c.QueueAttack(ai, combat.NewCircleHitOnTarget(t, nil, 3.5), 0, 1)

		char.AddStatus(icdKey, 10*60, true)
		t.AddStatus(debuffKey, 10*60, true)

		return false
	}, fmt.Sprintf("scion-%v", char.Base.Key.String()))
	return w, nil
}
