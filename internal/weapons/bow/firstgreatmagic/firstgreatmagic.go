package firstgreatmagic

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.TheFirstGreatMagic, NewWeapon)
}

type Weapon struct {
	Index            int
	core             *core.Core
	char             *character.CharWrapper
	atkStackVal      float64
	sameElement      int
	differentElement int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error {
	// calc Gimmick and Theatrics stacks
	for _, x := range w.core.Player.Chars() {
		if x.Base.Element == w.char.Base.Element { // includes wielder
			w.sameElement++
			continue
		}
		w.differentElement++
	}
	// cap element counts for calcing the buff values
	if w.sameElement > 3 {
		w.sameElement = 3
	}
	if w.differentElement > 3 {
		w.differentElement = 3
	}

	// Gimmick buff
	mAtk := make([]float64, attributes.EndStatType)
	mAtk[attributes.ATKP] = w.atkStackVal * float64(w.sameElement)
	w.char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("thefirstgreatmagic-atk", -1),
		AffectedStat: attributes.ATKP,
		Amount: func() ([]float64, bool) {
			return mAtk, true
		},
	})

	// Theatrics buff
	// TODO: movement speed is not implemented

	return nil
}

// DMG dealt by Charged Attacks increased by 16/20/24/28/32%.
// For every party member with the same Elemental Type as the wielder (including the wielder themselves), gain 1 Gimmick stack.
// For every party member with a different Elemental Type from the wielder, gain 1 Theatrics stack.
// When the wielder has 1/2/3 or more Gimmick stacks, ATK will be increased by 16%/32%/48% / 20%/40%/60% / 24%/48%/72% / 28%/56%/84% / 32%/64%/96%.
// When the wielder has 1/2/3 or more Theatrics stacks, Movement SPD will be increased by 4%/7%/10% / 6%/9%/12% / 8%/11%/14% / 10%/13%/16% / 12%/15%/18%.
func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{
		core: c,
		char: char,
	}
	r := p.Refine

	// CA buff
	mDmg := make([]float64, attributes.EndStatType)
	mDmg[attributes.DmgP] = (0.12 + float64(r)*0.04)
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("thefirstgreatmagic-dmg%", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagExtra {
				return nil, false
			}
			return mDmg, true
		},
	})

	w.atkStackVal = (0.12 + float64(r)*0.04)

	return w, nil
}
