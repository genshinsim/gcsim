package flowingpurity

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
	core.RegisterWeaponFunc(keys.FlowingPurity, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// When using an Elemental Skill, All Elemental DMG Bonus will be increased by 8/10/12/14/16% for 15s,
// and a Bond of Life worth 24% of Max HP will be granted. This effect can be triggered once every 10s.
// When the Bond of Life is cleared, every 1,000 HP cleared in the process will provide 2/2.5/3/3.5/4% All Elemental DMG Bonus,
// up to a maximum of 12/15/18/21/24%. This effect lasts 15s.

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	const icdKey = "flowingpurity-icd"
	const bondKey = "flowingpurity-bond"
	eledmg := 0.06 + float64(r)*0.02
	duration := 15 * 60
	icd := 10 * 60

	m := make([]float64, attributes.EndStatType)
	for i := attributes.PyroP; i <= attributes.DendroP; i++ {
		m[i] = eledmg
	}
	bond := make([]float64, attributes.EndStatType)
	hp := 0.24
	bondPercentage := 0.015 + float64(r)*0.005
	bondDMGPCap := 0.09 + float64(r)*0.03
	debt := 0.0

	c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		char.AddStatus(icdKey, icd, true)

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("flowingpurity-eledmg-boost", duration),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})

		if !char.StatusIsActive(bondKey) {
			debt = 0
		}
		char.AddStatus(bondKey, -1, true)

		char.SetHPDebtByRatio(hp)
		debt += char.CurrentHPDebt()
		bondDMGP := (debt / 1000) * bondPercentage // use hp debt since you only get the buff after clearing bond anyway
		if bondDMGP > bondDMGPCap {
			bondDMGP = bondDMGPCap
		}
		for i := attributes.PyroP; i <= attributes.DendroP; i++ {
			bond[i] = bondDMGP
		}

		return false
	}, fmt.Sprintf("flowingpurity-eledmg%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		index := args[1].(int)
		if index != char.Index {
			return false
		}
		if char.CurrentHPDebt() > 0 {
			return false
		}
		if !char.StatusIsActive(bondKey) {
			return false
		}
		char.DeleteStatus(bondKey)

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("flowingpurity-bond-eledmg-boost", duration),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				return bond, true
			},
		})
		return false
	}, fmt.Sprintf("flowingpurity-bondeledmg%v", char.Base.Key.String()))
	return w, nil
}
