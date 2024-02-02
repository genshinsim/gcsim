package lyney

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// If Lyney consumes HP when firing off a Prop Arrow,
// the Grin-Malkin hat summoned by the arrow will, upon hitting an opponent,
// restore 3 Energy to Lyney and increase DMG dealt by 80% of his ATK.
func (c *char) makeA1CB(hpDrained bool) combat.AttackCBFunc {
	if c.Base.Ascension < 1 || !hpDrained {
		return nil
	}
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if done {
			return
		}
		done = true
		c.AddEnergy("lyney-a1", 3)
	}
}

func (c *char) addA1(ai *combat.AttackInfo, hpDrained bool) {
	if c.Base.Ascension < 1 || !hpDrained {
		return
	}
	ai.Mult += 0.8
}

// The DMG Lyney deals to opponents affected by Pyro will receive the following buffs:
// - Increases the DMG dealt by 60%.
// - Each Pyro party member other than Lyney will cause the DMG dealt to increase by an additional 20%.
// Lyney can deal up to 100% increased DMG to opponents affected by Pyro in this way.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	// count up all pyro chars in team
	pyroCount := 0
	for _, char := range c.Core.Player.Chars() {
		if char.Base.Element == attributes.Pyro {
			pyroCount++
		}
	}

	// calc a4 dmg% value
	a4Dmg := 0.6 + float64(pyroCount-1)*0.2
	if a4Dmg > 1 {
		a4Dmg = 1
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = a4Dmg
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("lyney-a4", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			r, ok := t.(core.Reactable)
			if !ok {
				return nil, false
			}
			if !r.AuraContains(attributes.Pyro) {
				return nil, false
			}
			return m, true
		},
	})
}
