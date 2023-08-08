package lynette

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) a1Setup() {
	if c.Base.Ascension < 1 {
		return
	}

	// count elemental types
	partyEleTypes := make(map[attributes.Element]bool)
	for _, char := range c.Core.Player.Chars() {
		partyEleTypes[char.Base.Element] = true
	}
	count := len(partyEleTypes)

	// atk% buff setup
	c.a1Buff = make([]float64, attributes.EndStatType)
	c.a1Buff[attributes.ATKP] = 0.08 + float64(count-1)*0.04
}

// Within 10s after using Magic Trick: Astonishing Shift,
// when there are 1/2/3/4 Elemental Types in the party,
// all party members' ATK will be increased by 8%/12%/16%/20% respectively.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	for _, this := range c.Core.Player.Chars() {
		this.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("lynette-a1", 10*60),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return c.a1Buff, true
			},
		})
	}
}

// After the Bogglecat Box summoned by Magic Trick: Astonishing Shift performs Elemental Conversion,
// Lynette's Elemental Burst will deal 15% more DMG.
// This effect will persist until the Bogglecat Box's duration ends.
func (c *char) a4(duration int) {
	if c.Base.Ascension < 4 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.15
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("lynette-a4", duration),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagElementalBurst {
				return nil, false
			}
			return m, true
		},
	})
}
