package ineffa

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const a4Key = "ineffa-a4"

func (c *char) a1OnDischarge() {
	if c.Base.Ascension < 1 {
		return
	}

	// check if LC cloud is active

	a1Atk := func() {
		ai := info.AttackInfo{
			ActorIndex:       c.Index(),
			Abil:             "Birgitta (A1)",
			AttackTag:        attacks.AttackTagLCDamage,
			ICDTag:           attacks.ICDTagNone,
			ICDGroup:         attacks.ICDGroupDirectLunarCharged,
			StrikeType:       attacks.StrikeTypeDefault,
			Element:          attributes.Electro,
			Mult:             0.65,
			IgnoreDefPercent: 1,
		}
		c.Core.QueueAttack(ai, combat.NewBoxHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 4, 4.1), 0, 0)
	}

	c.Core.Tasks.Add(a1Atk, 30)
}

func (c *char) a4Init() {
	m := make([]float64, attributes.EndStatType)
	// TODO: Is this buff hitlag affected on Ineffa only? Or is it hitlag affected per character?
	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(a4Key+"-buff", -1),
			Extra:        true,
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				if c.Core.Player.Active() != char.Index() && c.Index() != char.Index() {
					return nil, false
				}
				if !c.StatusIsActive(a4Key) {
					return nil, false
				}

				stats := c.SelectStat(true, attributes.BaseATK, attributes.ATKP, attributes.ATK)
				m[attributes.EM] = stats.TotalATK() * 0.06
				return m, true
			},
		})
	}
}

func (c *char) a4OnBurst() {
	if c.Base.Ascension < 4 {
		return
	}

	// TODO when does this start and end?
	c.AddStatus(a4Key, 20*60, true)
}
