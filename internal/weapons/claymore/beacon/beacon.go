package beacon

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
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

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	//After the character's Elemental Skill hits an opponent, their ATK will be increased by 20% for 8s.
	//After the character takes DMG, their ATK will be increased by 20% for 8s.
	//The 2 aforementioned effects can be triggered even when the character is not on the field.
	//Additionally, when not protected by a shield, the character's Max HP will be increased by 32%.

	w := &Weapon{}
	r := p.Refine

	stackAtk := .15 + float64(r)*.05

	stackDuration := 480 //8s * 60
	const skillKey = "beacon-of-the-reed-sea-skill"
	const damagedKey = "beacon-of-the-reed-sea-damaged"

	mHP := make([]float64, attributes.EndStatType)
	mHP[attributes.HPP] = 0.24 + float64(r)*.08
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

	mATK := make([]float64, attributes.EndStatType)
	mATK[attributes.ATKP] = stackAtk
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagElementalArt && atk.Info.AttackTag != attacks.AttackTagElementalArtHold {
			return false
		}

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(skillKey, stackDuration),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return mATK, true
			},
		})

		return false
	}, fmt.Sprintf("beacon-of-the-reed-sea-enemy-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
		di := args[0].(player.DrainInfo)
		if di.ActorIndex != char.Index {
			return false
		}
		if di.Amount <= 0 {
			return false
		}
		if !di.External {
			return false
		}

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(damagedKey, stackDuration),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return mATK, true
			},
		})
		return false
	}, fmt.Sprintf("beacon-of-the-reed-sea-player-%v", char.Base.Key.String()))

	return w, nil
}
