package navia

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {

}

// For 4s after using Ceremonial Crystalshot, the DMG dealt by Navia's Normal Attacks,
// Charged Attacks, and Plunging Attacks will be converted into Geo DMG which cannot
// be overridden by other Elemental infusions, and the DMG dealt by Navia's Normal Attacks,
// Charged Attacks, and Plunging Attacks will be increased by 40%.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	c.Core.Log.NewEvent("infusion added", glog.LogCharacterEvent, c.Index)

	// add Damage Bonus
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag("navia-a1-dmg", 60*4), // 4s
		Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			// skip if not normal/charged/plunge
			if atk.Info.AttackTag != attacks.AttackTagNormal &&
				atk.Info.AttackTag != attacks.AttackTagExtra &&
				atk.Info.AttackTag != attacks.AttackTagPlunge {
				return nil, false
			}
			// apply buff
			m[attributes.DmgP] = 0.4
			return m, true
		},
	})
}

// For each Pyro/Electro/Cryo/Hydro party member, Navia gains 20% increased ATK.
// This effect can stack up to 2 times.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	ele := 0
	for _, char := range c.Core.Player.Chars() {
		if char.Base.Element != attributes.Geo && char.Base.
			Element != attributes.Anemo && char.Base.Element != attributes.Dendro {
			ele++
		}
	}
	if ele > 2 {
		ele = 2
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.2 * float64(ele)
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("navia-a4", -1),
		AffectedStat: attributes.ATKP,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}