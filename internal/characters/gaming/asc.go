package gaming

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const a1Key = "gaming-a1"

// After Bestial Ascent's Plunging Attack: Charmed Cloudstrider hits an opponent,
// Gaming will regain 1.5% of his Max HP once every 0.2s for 0.8s.
func (c *char) makeA1CB() combat.AttackCBFunc {
	if c.Base.Ascension < 1 {
		return nil
	}

	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if c.StatusIsActive(a1Key) {
			return
		}
		c.AddStatus(a1Key, 0.8*60, true)
		c.QueueCharTask(c.a1Heal, 0.2*60)
	}
}

func (c *char) a1Heal() {
	if !c.StatusIsActive(a1Key) {
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

// When Gaming has less than 50% HP, he will receive a 20% Incoming Healing Bonus.
// When Gaming has 50% HP or more, Plunging Attack: Charmed Cloudstrider will deal 20% more DMG.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	mHeal := make([]float64, attributes.EndStatType)
	mHeal[attributes.Heal] = 0.2
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("gaming-a4-heal-bonus", -1),
		AffectedStat: attributes.Heal,
		Amount: func() ([]float64, bool) {
			if c.CurrentHPRatio() >= 0.5 {
				return nil, false
			}
			return mHeal, true
		},
	})

	mDmg := make([]float64, attributes.EndStatType)
	mDmg[attributes.DmgP] = 0.2
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("gaming-a4-dmg-bonus", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if c.CurrentHPRatio() < 0.5 {
				return nil, false
			}
			if atk.Info.Abil != ePlungeKey {
				return nil, false
			}
			return mDmg, true
		},
	})
}
