package beacon

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.BeaconOfTheReedSea, NewWeapon)
}

type Weapon struct {
	Index  int
	stacks int
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
	w.stacks = p.Params["stacks"]

	//not modeling damage yet, so no duration
	//stackDuration := 720 // 12s * 60

	mATK := make([]float64, attributes.EndStatType)
	mHP := make([]float64, attributes.EndStatType)
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("beacon-atk", -1),
		AffectedStat: attributes.ATKP,
		Amount: func() ([]float64, bool) {
			if w.stacks >= 2 {
				w.stacks = 2
			}

			atkbonus := stackAtk * float64(w.stacks)

			mATK[attributes.ATKP] = atkbonus

			return mATK, true
		},
	})
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("beacon-hp", -1),
		AffectedStat: attributes.HPP,
		Amount: func() ([]float64, bool) {
			if !c.Player.Shields.PlayerIsShielded() {
				mHP[attributes.HPP] = 0.24 + float64(r)*.08
			}
			return mHP, true
		},
	})

	return w, nil
}
