package xilonen

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c2key = "xilonen-c2"
const c4key = "xilonen-c4"
const c6key = "xilonen-c6"
const c6IcdKey = "xilonen-c6-icd"
const c6StamKey = "xilonen-c6-stam"

func (c *char) nightsoulDurationMul() float64 {
	if c.Base.Cons < 1 {
		return 1.0
	}
	return 1.0 + 0.45
}

func (c *char) nightsoulConsumptionMul() float64 {
	if c.Base.Cons < 1 {
		return 1.0
	}
	return 1.0 - 0.3
}

func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}

	c.c2Buffs = map[attributes.Element][]float64{
		attributes.Geo:   make([]float64, attributes.EndStatType),
		attributes.Pyro:  make([]float64, attributes.EndStatType),
		attributes.Hydro: make([]float64, attributes.EndStatType),
		attributes.Cryo:  make([]float64, attributes.EndStatType),
	}
	c.c2Buffs[attributes.Geo][attributes.DmgP] = 0.5
	c.c2Buffs[attributes.Pyro][attributes.ATKP] = 0.45
	c.c2Buffs[attributes.Hydro][attributes.HPP] = 0.45
	c.c2Buffs[attributes.Cryo][attributes.CD] = 0.60

	if c.shredElements[attributes.Geo] {
		c.activeGeoSampler(-1)()
		chars := c.Core.Player.Chars()
		for _, ch := range chars {
			if ch.Base.Element != attributes.Geo {
				continue
			}

			ch.AddAttackMod(character.AttackMod{
				Base: modifier.NewBase(c2key, -1),
				Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
					return c.c2Buffs[attributes.Geo], true
				},
			})
		}
	}
}

func (c *char) applyC2Buff(ch *character.CharWrapper) func() {
	return func() {
		if !c.StatusIsActive(activeSamplerKey) {
			return
		}
		ch.AddStatMod(character.StatMod{
			Base: modifier.NewBaseWithHitlag(c2key, 60),
			Amount: func() ([]float64, bool) {
				return c.c2Buffs[ch.Base.Element], true
			},
		})
		c.QueueCharTask(c.applyC2Buff(ch), 0.1*60)
	}
}

func (c *char) c2activate() {
	if c.Base.Cons < 2 {
		return
	}
	chars := c.Core.Player.Chars()
	for _, ch := range chars {
		ele := ch.Base.Element
		switch ele {
		case attributes.Geo:
			// skip, because it's already applied
		case attributes.Pyro, attributes.Hydro, attributes.Cryo:
			c.QueueCharTask(c.applyC2Buff(ch), 0.3*60)
		case attributes.Electro:
			ch.AddEnergy(c2key, 25)
			ch.ReduceActionCooldown(action.ActionBurst, 6*60)
		default:
			continue
		}
	}
}

func (c *char) c4() {
	if c.Base.Cons < 4 {
		return
	}
	for _, char := range c.Core.Player.Chars() {
		char.AddStatus(c4key, 15*60, true) // 15 sec duration
		char.SetTag(c4key, 6)              // 6 c4 stacks
	}
}

func (c *char) c4Init() {
	if c.Base.Cons < 4 {
		return
	}
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)

		switch atk.Info.AttackTag {
		case attacks.AttackTagNormal:
		case attacks.AttackTagExtra:
		case attacks.AttackTagPlunge:
		default:
			return false
		}

		char := c.Core.Player.ByIndex(atk.Info.ActorIndex)
		if !char.StatusIsActive(c4key) || char.Tag(c4key) == 0 {
			return false
		}

		amt := 0.65 * c.TotalDef()
		char.SetTag(c4key, char.Tag(c4key)-1)

		c.Core.Log.NewEvent("xilonen c4 proc dmg add", glog.LogPreDamageMod, atk.Info.ActorIndex).
			Write("before", atk.Info.FlatDmg).
			Write("addition", amt).
			Write("effect_ends_at", c.StatusExpiry(c4key)).
			Write("c4_left", char.Tag(c4key))

		atk.Info.FlatDmg += amt
		return false
	}, fmt.Sprintf("%s-hook", c4key))
}

func (c *char) c6() {
	if c.Base.Cons < 6 {
		return
	}
	if c.StatusIsActive(c6IcdKey) {
		return
	}
	if !c.nightsoulState.HasBlessing() {
		return
	}

	c.AddStatus(c6key, 5*60, true)
	c.AddStatus(c6IcdKey, 15*60, true)

	// "pause" Nightsoul's Blessing time limit countdown
	duration := c.StatusDuration(skillMaxDurKey) + 5*60
	c.setNightsoulExitTimer(duration)

	for i := 1; i <= 4; i++ {
		c.Core.Tasks.Add(func() {
			hpplus := c.Stat(attributes.Heal)
			heal := c.TotalDef() * 1.2
			c.Core.Player.Heal(info.HealInfo{
				Caller:  c.Index,
				Target:  -1,
				Message: "Imperishable Night Carnival",
				Src:     heal,
				Bonus:   hpplus,
			})
		}, i*90)
	}
}

func (c *char) c6Stam() {
	if c.Base.Cons < 6 {
		return
	}
	c.Core.Player.AddStamPercentMod(c6StamKey, -1, func(a action.Action) (float64, bool) {
		if c.StatusIsActive(c6key) {
			return -1, false
		}
		return 0, false
	})
}

func (c *char) c6DmgMult() float64 {
	if c.Base.Cons < 6 {
		return 0.0
	}
	if !c.StatusIsActive(c6key) {
		return 0.0
	}
	return 3.0
}
