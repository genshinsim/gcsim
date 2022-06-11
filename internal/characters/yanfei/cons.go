package yanfei

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

// Hook for C2:
// Increases Yan Fei's Charged Attack CRIT Rate by 20% against enemies below 50% HP.
func (c *char) c2() {
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod("yanfei-c2", -1, func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
		if atk.Info.AttackTag != combat.AttackTagExtra {
			return nil, false
		}
		if t.HP()/t.MaxHP() >= .5 {
			return nil, false
		}
		m[attributes.CR] = 0.20
		return m, true
	})
}

// Handles C4 shield creation
// When Done Deal is used:
// Creates a shield that absorbs up to 45% of Yan Fei's Max HP for 15s
// This shield absorbs Pyro DMG 250% more effectively
func (c *char) c4() {
	if c.Base.Cons < 4 {
		return
	}
	c.Core.Player.Shields.Add(&shield.Tmpl{
		Src:        c.Core.F,
		ShieldType: shield.ShieldYanfeiC4,
		Name:       "Yanfei C4",
		HP:         c.MaxHP() * .45,
		Ele:        attributes.Pyro,
		Expires:    c.Core.F + 15*60,
	})
}
