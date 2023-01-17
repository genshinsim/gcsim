package nahida

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const a1BuffKey = "nahida-a1"
const a4BuffKey = "nahida-a4"

// When unleashing Illusory Heart, the Shrine of Maya will gain the following effects:
// The Elemental Mastery of the active character within the field will be increased by 25% of the Elemental Mastery of the party member with the highest Elemental Mastery.
// You can gain a maximum of 250 Elemental Mastery in this manner.
func (c *char) calcA1Buff() {
	var max float64
	team := c.Core.Player.Chars()
	for _, char := range team {
		em := char.Stat(attributes.EM)
		if em > max {
			max = em
		}
	}
	max = 0.25 * max

	if max > 250 {
		max = 250
	}

	c.a1Buff[attributes.EM] = max
}

func (c *char) applyA1(dur int) {
	for i, char := range c.Core.Player.Chars() {
		idx := i
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase(a1BuffKey, dur),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return c.a1Buff, c.Core.Player.Active() == idx
			},
		})
	}
}

// Each point of Nahida's Elemental Mastery beyond 200 will grant 0.1% Bonus DMG and 0.03% CRIT Rate to Tri-Karma Purification from All Schemes to Know.
// A maximum of 80% Bonus DMG and 24% CRIT Rate can be granted to Tri-Karma Purification in this manner.
func (c *char) a4() {
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(a4BuffKey, -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != combat.AttackTagElementalArt {
				return nil, false
			}
			if !strings.HasPrefix(atk.Info.Abil, "Tri-Karma") {
				return nil, false
			}
			return c.a4Buff, true
		},
	})
}

func (c *char) a4tick() {
	em := c.Stat(attributes.EM)
	var dmgBuff, crBuff float64
	if em > 200 {
		em = em - 200
		dmgBuff = em * 0.001
		if dmgBuff > 0.8 {
			dmgBuff = 0.8
		}
		crBuff = em * 0.0003
		if crBuff > .24 {
			crBuff = .24
		}
	}
	c.a4Buff[attributes.DmgP] = dmgBuff
	c.a4Buff[attributes.CR] = crBuff

	c.Core.Tasks.Add(c.a4tick, 30)
}
