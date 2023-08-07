package flowingpurity

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
	core.RegisterWeaponFunc(keys.FlowingPurity, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// When using an Elemental Skill, All Elemental DMG Bonus will be increased by 8% for 15s, and a Bond of Life worth 24% of Max HP will be granted.
// This effect can be triggered once every 10s.
// When the Bond of Life is cleared, every 1,000 HP cleared in the process will provide 2% All Elemental DMG Bonus.
// Up to a maximum of 12% All Elemental DMG can be gained this way. This effect lasts 15s.
// Bond of Life: Absorbs healing for the character based on its base value, and clears after healing equal to this value is obtained.

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	const icdKey = "flowingpurity-icd"
	eledmg := 0.06 + float64(r)*0.02
	duration := 900 //15s * 60
	icd := 600      //10s * 60
	m := make([]float64, attributes.EndStatType)
	bond := make([]float64, attributes.EndStatType)
	hp := 0.24 //hpdebt_percentage
	bondPercentage := 0.015 + float64(r)*0.005

	c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
		if char.StatusIsActive(icdKey) {
			return false
		}
		char.AddStatus(icdKey, icd, true)
		for i := attributes.PyroP; i <= attributes.DendroP; i++ {
			m[i] = eledmg
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("flowingpurity-eledmg-boost", duration),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
		char.SetHPDebtByRatio(hp) //set hpdebt before getting healed
		debt := char.CurrentHPDebt()
		if debt >= 6000 {
			debt = 6000 //debt = maxbondp / bondp
		}
		// not sure if after (?) seconds the bond of life gonna clear itself, thus not implement yet
		c.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
			healInfo := args[0].(*player.HealInfo)
			healAmt := args[2].(float64)
			if healInfo.Target != -1 && healInfo.Target != char.Index {
				return false
			}
			if ((healAmt - char.CurrentHPDebt()) <= 0) || char.CurrentHPDebt() > 0 {
				char.ModifyHPDebtByAmount(-healAmt) //reduce heal debt, but there is still heal debt
				return false
			} else {
				char.ModifyHPDebtByAmount(-healAmt)
				char.AddStatMod(character.StatMod{
					Base:         modifier.NewBase("flowingpurity-bondeledmg-boost", duration),
					AffectedStat: attributes.NoStat,
					Amount: func() ([]float64, bool) {
						bondDMGP := (debt / 1000) * bondPercentage //use hp debt since you only get the buff after clearing bond
						for i := attributes.PyroP; i <= attributes.DendroP; i++ {
							bond[i] = bondDMGP
						}
						return bond, true
					},
				})
			}
			return false
		}, fmt.Sprintf("flowingpurity-bondeledmg%v", char.Base.Key.String()))
		return false
	}, fmt.Sprintf("flowingpurity-eledmg%v", char.Base.Key.String()))

	return w, nil
}
