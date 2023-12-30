package wriothesley

import (
	"math"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a1Status = "wriothesley-a1"
	a1ICD    = 5 * 60
	a1ICDKey = "wriothesley-a1-icd"
)

func (c *char) a1Ready() bool {
	return c.CurrentHPRatio() < 0.6 && !c.StatusIsActive(a1ICDKey)
}

// When Wriothesley's HP is less than 60%, he will obtain a Gracious Rebuke. The next Charged Attack of his
// Normal Attack: Forceful Fists of Frost will be enhanced to become Rebuke: Vaulting Fist. It will not consume
// Stamina, deal 50% increased DMG, and will restore HP for Wriothesley after hitting equal to 30% of his Max HP.
// You can gain a Gracious Rebuke this way once every 5s.
func (c *char) a1(ai *combat.AttackInfo, snap *combat.Snapshot) combat.AttackCBFunc {
	if !c.a1Ready() {
		return nil
	}

	// add status that is removed on consumption
	c.AddStatus(a1Status, -1, false)

	// adjust ai
	ai.Abil = "Rebuke: Vaulting Fist"
	ai.HitlagFactor = 0.03
	ai.HitlagHaltFrames = 0.12 * 60

	// 50% increased DMG
	dmg := 0.5
	snap.Stats[attributes.DmgP] += dmg
	c.Core.Log.NewEvent("adding a1", glog.LogCharacterEvent, c.Index).Write("dmg%", dmg)

	// return callback to heal, remove A1 and apply 5s cd
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		// do not proc if a1 not active
		if !c.StatusIsActive(a1Status) {
			return
		}
		// remove A1 and apply CD
		c.DeleteStatus(a1Status)
		c.AddStatus(a1ICDKey, a1ICD, true)

		// heal
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  c.Index,
			Message: "There Shall Be a Plea for Justice",
			Src:     c.caHeal * c.MaxHP(),
			Bonus:   c.Stat(attributes.Heal),
		})
	}
}

// When Wriothesley's current HP increases or decreases, if he is in the Chilling Penalty state conferred by Icefang Rush,
// Chilling Penalty will gain one stack of Prosecution Edict. Max 5 stacks. Each stack will increase Wriothesley's ATK by 6%.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	c.Core.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
		di := args[0].(player.DrainInfo)
		if c.Core.Player.Active() != c.Index {
			return false
		}
		if di.ActorIndex != c.Index {
			return false
		}
		if di.Amount <= 0 {
			return false
		}

		if c.StatusIsActive(skillKey) && c.a4Stack < 5 {
			c.a4Stack++
			c.Core.Log.NewEvent("a4 gained stack", glog.LogCharacterEvent, c.Index).Write("stacks", c.a4Stack)
		}
		return false
	}, "wriothesley-a4-drain")

	c.Core.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		index := args[1].(int)
		amount := args[2].(float64)
		overheal := args[3].(float64)
		if c.Core.Player.Active() != c.Index {
			return false
		}
		if index != c.Index {
			return false
		}
		if amount <= 0 {
			return false
		}
		// do not trigger if at max hp already
		if math.Abs(amount-overheal) <= 1e-9 {
			return false
		}

		if c.StatusIsActive(skillKey) && c.a4Stack < 5 {
			c.a4Stack++
			c.Core.Log.NewEvent("a4 gained stack", glog.LogCharacterEvent, c.Index).Write("stacks", c.a4Stack)
		}
		return false
	}, "wriothesley-a4-heal")
}

func (c *char) applyA4(dur int) {
	if c.Base.Ascension < 4 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("wriothesley-a4", dur),
		AffectedStat: attributes.ATKP,
		Amount: func() ([]float64, bool) {
			m[attributes.ATKP] = float64(c.a4Stack) * 0.06
			return m, true
		},
	})
}

func (c *char) resetA4() {
	if c.Base.Ascension < 4 {
		return
	}
	c.a4Stack = 0
}
