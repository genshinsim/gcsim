package razor

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const hexereiICDKey = "razor-hexerei-icd"

// Decreases Claw and Thunder's CD by 18%.
func (c *char) a1CDReduction(cd int) int {
	if c.Base.Ascension < 1 {
		return cd
	}
	return int(float64(cd) * 0.82)
}

// Using Lightning Fang resets the CD of Claw and Thunder.
func (c *char) a1CDReset() {
	if c.Base.Ascension < 1 {
		return
	}
	c.ResetActionCooldown(action.ActionSkill)
}

// When Razor's Energy is below 50%, increases Energy Recharge by 30%.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	c.a4Bonus = make([]float64, attributes.EndStatType)
	c.a4Bonus[attributes.ER] = 0.3
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("razor-a4", -1),
		AffectedStat: attributes.ER,
		Amount: func() []float64 {
			if c.Energy/c.EnergyMax >= 0.5 {
				return nil
			}
			return c.a4Bonus
		},
	})
}

func (c *char) thunderFallCB() {
	if c.StatusIsActive(hexereiICDKey) {
		return
	}

	c.c6HexereiMod()

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Surge of Lightning",
		AttackTag:  attacks.AttackTagNone,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       1.5,
	}

	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		info.Point{Y: 2},
		5,
	)
	c.AddStatus(hexereiICDKey, 60+74, false)

	c.Core.QueueAttack(ai, ap, 0, 37)
	c.Core.Tasks.Add(
		func() {
			c.AddEnergy("razor-surge-of-lightning", 7)
		}, 74,
	)
}
