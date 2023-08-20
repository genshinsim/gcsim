package wriothesley

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a1Status = "wriothesley-a1"
	a1ICDKey = "wriothesley-a1-icd"
)

// When Wriothesley's HP is less than 60%, he will obtain a Gracious Rebuke. The next Charged Attack of his
// Normal Attack: Forceful Fists of Frost will be enhanced to become Rebuke: Vaulting Fist. It will not consume
// Stamina, deal 30% increased DMG, and will restore HP for Wriothesley after hitting equal to 30% of his Max HP.
// You can gain a Gracious Rebuke this way once every 5s.
func (c *char) a1Add() {
	if c.StatusIsActive(a1ICDKey) {
		return
	}
	c.AddStatus(a1ICDKey, 5*60, true)

	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.3
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag(a1Status, -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag == attacks.AttackTagExtra {
				return m, true
			}
			return nil, false
		},
	})
}

func (c *char) a1Remove(_ combat.AttackCB) {
	if !c.StatModIsActive(a1Status) {
		return
	}

	c.Core.Player.Heal(player.HealInfo{
		Caller:  c.Index,
		Target:  c.Index,
		Message: "There Shall Be a Plea for Justice",
		Src:     c.MaxHP() * 0.3,
		Bonus:   c.Stat(attributes.Heal),
	})
	c.DeleteStatMod(a1Status)
}
