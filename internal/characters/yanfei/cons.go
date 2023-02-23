package yanfei

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// Hook for C2:
// Increases Yan Fei's Charged Attack CRIT Rate by 20% against enemies below 50% HP.
func (c *char) c2() {
	if c.Core.Combat.DamageMode {
		m := make([]float64, attributes.EndStatType)
		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("yanfei-c2", -1),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if atk.Info.AttackTag != attacks.AttackTagExtra {
					return nil, false
				}
				x, ok := t.(*enemy.Enemy)
				if !ok {
					return nil, false
				}
				if x.HP()/x.MaxHP() >= .5 {
					return nil, false
				}
				m[attributes.CR] = 0.20
				return m, true
			},
		})
	}
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
