package beacon

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
	core.RegisterWeaponFunc(keys.BeaconOfTheReedSea, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	//After the character's Elemental Skill hits an opponent, their ATK will be increased by 20% for 8s.
	//After the character takes DMG, their ATK will be increased by 20% for 8s.
	//The 2 aforementioned effects can be triggered even when the character is not on the field.
	//Additionally, when not protected by a shield, the character's Max HP will be increased by 32%.

	w := &Weapon{}
	r := p.Refine

	stackAtk := .15 + float64(r)*.05
	damaged := p.Params["damaged"]

	stackDuration := 480 //8s * 60
	const skillKey = "beacon-of-the-reed-sea-skill"
	const damagedKey = "beacon-of-the-reed-sea-damaged"

	mATK := make([]float64, attributes.EndStatType)
	mHP := make([]float64, attributes.EndStatType)
	mHP[attributes.HPP] = 0.24 + float64(r)*.08
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("beacon-of-the-reed-sea-atk", -1),
		AffectedStat: attributes.ATKP,
		Amount: func() ([]float64, bool) {
			count := 0

			if char.StatusIsActive(skillKey) {
				count++
			}
			if char.StatusIsActive(damagedKey) {
				count++
			}

			atkbonus := stackAtk * float64(count)

			mATK[attributes.ATKP] = atkbonus

			return mATK, true
		},
	})
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("beacon-of-the-reed-sea-hp", -1),
		AffectedStat: attributes.HPP,
		Amount: func() ([]float64, bool) {
			if c.Player.Shields.PlayerIsShielded() {
				return nil, false
			}
			return mHP, true
		},
	})

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}

		if atk.Info.AttackTag == attacks.AttackTagElementalArt || atk.Info.AttackTag == attacks.AttackTagElementalArtHold {
			char.AddStatus(skillKey, stackDuration, true)
			if damaged > 0 {
				char.AddStatus(damagedKey, stackDuration, true)
			}

		}
		return false
	}, fmt.Sprintf("beacon-of-the-reed-sea-%v", char.Base.Key.String()))

	return w, nil
}
