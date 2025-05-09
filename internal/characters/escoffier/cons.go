package escoffier

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c1Dur = 15 * 60
const c1Key = "escoffier-c1"
const c2Key = "escoffier-c2"
const c2Per = 2.4
const c2Dur = 15 * 60
const c4Extra = 6
const c4ExtraScaling = 1
const c4Limit = 7
const c4Key = "escoffier-c4"
const c4Regen = 2.0
const c6Key = "escoffier-c6"
const c6Limit = 6
const c6Scaling = 5
const c6ICD = 0.5 * 60
const c6ICDKey = "escoffier-c6-icd"

func (c *char) c1Init() {
	if c.Base.Cons < 1 {
		return
	}
	if c.Base.Ascension < 4 {
		return
	}
	c.c1Active = false
	for _, char := range c.Core.Player.Chars() {
		switch char.Base.Element {
		case attributes.Cryo, attributes.Hydro:
		default:
			return
		}
	}
	c.c1Active = true
	c.c1Buff = make([]float64, attributes.EndStatType)
	c.c1Buff[attributes.CD] = 0.60
}

func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}
	if c.Base.Ascension < 4 {
		return
	}
	if !c.c1Active {
		return
	}
	for _, char := range c.Core.Player.Chars() {
		// TODO: check if this buff is affected by hitlag on characters or hitlag on escoffier
		// Currently assuming affected by hitlag on characters
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag(c1Key, c1Dur),
			Amount: func(ae *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
				if ae.Info.Element != attributes.Cryo {
					return nil, false
				}
				return c.c1Buff, true
			},
		})
	}
}

func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if c.Index == atk.Info.ActorIndex {
			return false
		}
		if c.Core.Player.Active() != atk.Info.ActorIndex {
			return false
		}
		if !c.StatusIsActive(c2Key) {
			return false
		}
		if c.c2Count <= 0 {
			return false
		}
		amt := c.TotalAtk() * c2Per
		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("Escoffier C2 proc dmg add", glog.LogPreDamageMod, atk.Info.ActorIndex).
				Write("before", atk.Info.FlatDmg).
				Write("addition", amt).
				Write("Cold Dishes left", c.c2Count)
		}
		atk.Info.FlatDmg += amt
		c.c2Count--
		return false
	}, c2Key+"-on-dmg")
}

func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}
	c.AddStatus(c2Key, c2Dur, true)
	c.c2Count = 5
}

func (c *char) c4ExtraCount() int {
	if c.Base.Cons < 4 {
		return 0
	}
	if c.Base.Ascension < 1 {
		return 0
	}
	c.c4Count = c4Limit
	return c4Extra
}

func (c *char) c4ExtraHeal() float64 {
	if c.Base.Cons < 4 {
		return 0
	}
	if c.Base.Ascension < 1 {
		return 0
	}
	if c.c4Count <= 0 {
		return 0
	}

	if c.Core.Rand.Float64() > c.Stat(attributes.CR) {
		return 0
	}
	c.AddEnergy(c4Key, c4Regen)
	c.c4Count--

	return c4ExtraScaling
}

func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if c.Core.Player.Active() != atk.Info.ActorIndex {
			return false
		}
		switch atk.Info.AttackTag {
		case attacks.AttackTagNormal, attacks.AttackTagExtra, attacks.AttackTagPlunge:
		default:
			return false
		}
		if !c.StatusIsActive(skillKey) {
			return false
		}
		if c.c6Count <= 0 {
			return false
		}
		if c.StatusIsActive(c6ICDKey) {
			return false
		}

		c.AddStatus(c6ICDKey, c6ICD, true)
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Special-Grade Frozen Parfait (C6)",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagElementalArt,
			ICDGroup:   attacks.ICDGroupEscoffierSkill,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       c6Scaling,
		}
		// trigger damage
		//TODO: travel time
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 2), 0, 10, c.makeA4CB())

		c.c6Count--

		return false
	}, c6Key+"-on-nacp")
}

func (c *char) c6() {
	if c.Base.Cons < 6 {
		return
	}
	c.c6Count = c6Limit
}
