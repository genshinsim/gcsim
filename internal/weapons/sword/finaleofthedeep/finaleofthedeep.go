package finaleofthedeep

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
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

// When using an Elemental Skill, ATK will be increased by 12/15/18/21/24% for 15s,
// and a Bond of Life worth 25% of Max HP will be granted. This effect can be triggered once every 10s.
// When the Bond of Life is cleared, a maximum of 150/187.5/225/262.5/300 ATK will be gained
// based on 2.4/3/3.6/4.2/4.8% of the total amount of the Life Bond cleared, lasting for 15s.

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	const icdKey = "finaleofthedeep-icd"
	const bondKey = "finaleofthedeep-bond"
	atk := 0.09 + float64(r)*0.03
	duration := 15 * 60
	icd := 10 * 60
	ATKval := make([]float64, attributes.EndStatType)
	ATKval[attributes.ATKP] = atk

	hp := 0.25
	bondPercentage := 0.018 + float64(r)*0.006
	maxBondAtk := 112.5 + float64(r)*37.5
	fAtkVal := make([]float64, attributes.EndStatType)
	totalHeal := float64(0)

	maxhp := char.MaxHP()
	bondAtk := maxhp * hp * bondPercentage
	if bondAtk >= maxBondAtk {
		bondAtk = maxBondAtk
	}

	c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
		if char.StatusIsActive(icdKey) {
			return false
		}
		char.AddStatus(icdKey, icd, true)
		char.AddStatus(bondKey, -1, true) // not sure if after (?) seconds the bond of life gonna clear itself
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("finaleofthedeep-atk-boost", duration),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return ATKval, true
			},
		})
		return false
	}, fmt.Sprintf("finaleofthedeep-atk%v", char.Base.Key.String()))

	// check for accummulate healing, when enough healing then get the ATK buff
	// not sure if after (?) seconds the bond of life gonna clear itself, thus not implement yet
	c.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		healAmt := args[2].(float64)
		index := args[1].(int)
		if index != char.Index {
			return false
		}
		if !char.StatusIsActive(bondKey) {
			return false
		}
		totalHeal += healAmt
		if totalHeal >= maxhp*hp {
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("finaleofthedeep-bond-flatatk-boost", duration),
				AffectedStat: attributes.ATK,
				Amount: func() ([]float64, bool) {
					fAtkVal[attributes.ATK] = bondAtk
					return fAtkVal, true
				},
			})
			totalHeal = 0
		}
		return false
	}, fmt.Sprintf("finaleofthedeep-flatatk%v", char.Base.Key.String()))
	return w, nil
}
