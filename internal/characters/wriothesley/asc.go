package wriothesley

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
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

// When Wriothesley's current HP increases or decreases, if he is in the Chilling Penalty state conferred by Icefang Rush,
// Chilling Penalty will gain one stack of Prosecution Edict. Max 5 stacks. Each stack will increase Wriothesley's ATK by 6%.
func (c *char) a4() {
	m := make([]float64, attributes.EndStatType)
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("wriothesley-a4", -1),
		AffectedStat: attributes.ATKP,
		Amount: func() ([]float64, bool) {
			m[attributes.ATKP] = float64(c.a4Stack) * 0.06
			return m, true
		},
	})

	c.Core.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
		di := args[0].(player.DrainInfo)
		if c.Core.Player.Active() != c.Index { // TODO: works off-field?
			return false
		}
		if di.ActorIndex != c.Index {
			return false
		}
		if di.Amount <= 0 {
			return false
		}

		if c.a4Stack < 5 {
			c.a4Stack++
		}
		return false
	}, "wriothesley-a4-drain")

	c.Core.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		index := args[1].(int)
		amount := args[2].(float64)
		if c.Core.Player.Active() != c.Index { // TODO: works off-field?
			return false
		}
		if index != c.Index {
			return false
		}
		if amount <= 0 {
			return false
		}

		if c.a4Stack < 5 {
			c.a4Stack++
		}
		return false
	}, "wriothesley-a4-heal")
}
