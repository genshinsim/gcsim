package sealord

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
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.LuxuriousSeaLord, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	//Increases Elemental Burst DMG by 12%. When Elemental Burst hits opponents,
	//there is a 100% chance of summoning a huge onrush of tuna that deals 100%
	//ATK as AoE DMG. This effect can occur once every 15s.
	w := &Weapon{}
	r := p.Refine

	//perm burst dmg increase
	burstDmgIncrease := .09 + float64(r)*0.03
	val := make([]float64, attributes.EndStatType)
	val[attributes.DmgP] = burstDmgIncrease
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("luxurious-sea-lord", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag == attacks.AttackTagElementalBurst {
				return val, true
			}
			return nil, false
		},
	})

	tunaDmg := .75 + float64(r)*0.25
	const icdKey = "sealord-icd"
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
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
		if atk.Info.AttackTag != attacks.AttackTagElementalBurst {
			return false
		}
		char.AddStatus(icdKey, 900, true)
		ai := combat.AttackInfo{
			ActorIndex: char.Index,
			Abil:       "Luxurious Sea-Lord Proc",
			AttackTag:  attacks.AttackTagWeaponSkill,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Physical,
			Durability: 100,
			Mult:       tunaDmg,
		}
		trg := args[0].(combat.Target)
		c.QueueAttack(ai, combat.NewCircleHitOnTarget(trg, nil, 3), 0, 1)

		return false
	}, fmt.Sprintf("sealord-%v", char.Base.Key.String()))
	return w, nil
}
