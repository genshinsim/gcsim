package finaleofthedeep

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.FinaleOfTheDeep, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// When using an Elemental Skill, ATK will be increased by 12% for 12s,
// and a Bond of Life worth 25% of Max HP will be granted. This effect can be triggered once every 10s.
// When the Bond of Life is cleared, a maximum of 150/187.5/225/262.5/300 ATK will be gained based on 2.4% of the Bond for 12s.
// Bond of Life: Absorbs healing for the character based on its base value, and clears after healing equal to this value is obtained.

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	const icdKey = "finaleofthedeep-icd"
	atk := 0.09 + float64(r)*0.03
	duration := 720 //12s * 60
	icd := 600      //10s * 60
	ATKval := make([]float64, attributes.EndStatType)
	ATKval[attributes.ATKP] = atk

	hp := 0.25
	bondPercentage := 0.018 + float64(r)*0.006
	maxBondAtk := 112.5 + float64(r)*37.5
	fAtkVal := make([]float64, attributes.EndStatType)
	totalHeal := float64(0)

	c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
		if char.StatusIsActive(icdKey) {
			return false
		}
		char.AddStatus(icdKey, icd, true)
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("finaleofthedeep-atk-boost", duration),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				return ATKval, true
			},
		})
		// check for accummulate healing, when enough healing then get the ATK buff
		// absorb healing not implement yet
		// not sure if after (?) seconds the bond of life gonna clear itself, thus not implement yet
		c.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
			healInfo := args[0].(*player.HealInfo)
			healAmt := args[2].(float64)
			maxhp := char.MaxHP()
			if healInfo.Target != -1 && healInfo.Target != char.Index {
				return false
			}
			totalHeal += healAmt
			if totalHeal >= maxhp*hp {
				char.AddStatMod(character.StatMod{
					Base:         modifier.NewBase("finaleofthedeep-flatatk-boost", duration),
					AffectedStat: attributes.ATK,
					Extra:        true,
					Amount: func() ([]float64, bool) {
						bondAtk := maxhp * hp * bondPercentage
						if bondAtk >= maxBondAtk {
							fAtkVal[attributes.ATK] = maxBondAtk
						} else {
							fAtkVal[attributes.ATK] = bondAtk
						}
						return fAtkVal, true
					},
				})

				totalHeal = 0
			}
			return false
		}, fmt.Sprintf("finaleofthedeep-flatatk%v", char.Base.Key.String()))
		return false
	}, fmt.Sprintf("finaleofthedeep-atk%v", char.Base.Key.String()))

	return w, nil
}
