package finaleofthedeep

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.FinaleOfTheDeep, NewWeapon)
}

type Weapon struct {
	Index         int
	collectedDebt float64
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error {
	w.collectedDebt = 0.
	return nil
}

// When using an Elemental Skill, ATK will be increased by 12/15/18/21/24% for 15s,
// and a Bond of Life worth 25% of Max HP will be granted. This effect can be triggered once every 10s.
// When the Bond of Life is cleared, a maximum of 150/187.5/225/262.5/300 ATK will be gained
// based on 2.4/3/3.6/4.2/4.8% of the total amount of the Life Bond cleared, lasting for 15s.

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	const icdKey = "finaleofthedeep-icd"
	atk := 0.09 + float64(r)*0.03
	duration := 15 * 60
	icd := 10 * 60

	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = atk
	bond := make([]float64, attributes.EndStatType)
	hp := 0.25
	bondPercentage := 0.018 + float64(r)*0.006
	bondAtkCap := 112.5 + float64(r)*37.5

	c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		char.AddStatus(icdKey, icd, true)

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("finaleofthedeep-atk-boost", duration),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})

		char.ModifyHPDebtByRatio(hp)

		return false
	}, fmt.Sprintf("finaleofthedeep-atk%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnHPDebt, func(args ...interface{}) bool {
		index := args[0].(int)
		debtChange := args[1].(float64)
		if index != char.Index {
			return false
		}

		if debtChange < 0 {
			w.collectedDebt += -float64(debtChange)
		}
		if char.CurrentHPDebt() > 0 {
			return false
		}
		if w.collectedDebt < 0.0001 {
			return false
		}

		bondAtk := min(bondAtkCap, w.collectedDebt*bondPercentage)
		bond[attributes.ATK] = bondAtk
		w.collectedDebt = 0

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("finaleofthedeep-bond-flatatk-boost", duration),
			AffectedStat: attributes.ATK,
			Amount: func() ([]float64, bool) {
				return bond, true
			},
		})
		return false
	}, fmt.Sprintf("finaleofthedeep-flatatk%v", char.Base.Key.String()))
	return w, nil
}
