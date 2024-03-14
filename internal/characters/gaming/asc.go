package gaming

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const a1key = "gaming-a1"

/*
For 1s after hitting an opponent with Bestial Ascent's Plunging Attack: Charmed Cloudstrider,

	Gaming will recover 10% of his HP.
*/
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	c.AddStatus(a1key, 0.8*60, true)
	c.QueueCharTask(c.a1Heal, 0.2*60)
}

func (c *char) a1Heal() {
	if !c.StatusIsActive(a1key) {
		return
	}
	c.Core.Player.Heal(player.HealInfo{
		Caller:  c.Index,
		Target:  c.Index,
		Message: "Dance of Amity (A1)",
		Type:    player.HealTypePercent,
		Src:     0.015,
		Bonus:   c.Stat(attributes.Heal),
	})
	c.QueueCharTask(c.a1Heal, 0.2*60)
}

/*
When Gaming has less than 50% HP, he will receive a 20% Incoming Healing Bonus.
When Gaming has 50% HP or more, Plunging Attack: Charmed Cloudstrider will deal 20% more DMG.
TODO: Confirm this is 20% DMG and not a 1.2x multiplier to MVs
*/
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	// Healing part
	mHeal := make([]float64, attributes.EndStatType)
	mHeal[attributes.Heal] = 0.2
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("gaming-a4-heal-bonus", -1),
		AffectedStat: attributes.Heal,
		Amount: func() ([]float64, bool) {
			active := c.Core.Player.ActiveChar()
			if active.CurrentHPRatio() < 0.5 {
				return mHeal, true
			}
			return nil, false
		},
	})

	a4Buff := make([]float64, attributes.EndStatType)
	a4Buff[attributes.PyroP] = 0.2
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("gaming-a4-dmg-bonus", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagPlunge {
				return nil, false
			}
			if c.CurrentHPRatio() < 0.5 {
				return nil, false
			}
			if !strings.Contains(atk.Info.Abil, ePlungeKey) {
				return nil, false
			}
			return a4Buff, true
		},
	})
}
