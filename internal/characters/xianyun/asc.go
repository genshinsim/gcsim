package xianyun

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const a4ICDKey = "xianyun-a4-icd"
const a1Key = "xianyun-a1"
const a1Dur = 20 * 60

var a1Crit = []float64{0.0, 0.04, 0.06, 0.08, 0.10}

// a1: For every enemy hit by Driftcloud Wave,
// all party members gain 1 stack of Boost, lasting 20s, max 4 stacks.
// Boost increases Plunge DMG's Crit Rate by 4%/6%/8%/10%,
// each stack's duration is calculated independently.

func (c *char) a1() {
	for idx, char := range c.Core.Player.Chars() {
		mCR := make([]float64, attributes.EndStatType)
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("xianyun-a1-buff", -1),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if atk.Info.AttackTag != attacks.AttackTagPlunge {
					return nil, false
				}
				stackCount := min(c.a1Buffer[idx], 4)
				if stackCount == 0 {
					return nil, false
				}
				mCR[attributes.CR] = a1Crit[stackCount]
				return mCR, true
			},
		})
	}
}

func (c *char) a1cb() combat.AttackCBFunc {
	if c.Base.Ascension < 1 {
		return nil
	}
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}

		for i, char := range c.Core.Player.Chars() {
			idx := i
			c.a1Buffer[idx] += 1
			char.SetTag(a1Key, min(c.a1Buffer[idx], 4))
			char.QueueCharTask(func() {
				char.SetTag(a1Key, min(c.a1Buffer[idx], 4))
			}, a1Dur)
		}
	}
}

// a4: When the Starwicker created by Stars Gather at Dusk has Adeptal Assistance stacks,
// nearby active characters' Plunging Attack shockwave DMG will be increased by 200% of Xianyun's ATK.
// The maximum DMG increase that can be achieved this way is 9000.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		if ae.Info.AttackTag != attacks.AttackTagPlunge {
			return false
		}

		// Collision has 0 durability. Don't buff collision damage
		if ae.Info.Durability == 0 {
			return false
		}

		if !c.StatusIsActive(StarwickerKey) {
			return false
		}

		if c.StatusIsActive(a4ICDKey) {
			return false
		}

		atk := c.Base.Atk*(1+c.Stat(attributes.ATKP)) + c.Stat(attributes.ATK)
		amt := c.a4Ratio * atk

		// A4 cap
		amt = min(c.a4Max, amt)

		c.Core.Log.NewEvent("Xianyun Starwicker proc dmg add", glog.LogPreDamageMod, ae.Info.ActorIndex).
			Write("atk", atk).
			Write("ratio", c.a4Ratio).
			Write("before", ae.Info.FlatDmg).
			Write("addition", amt).
			Write("effect_ends_at", c.StatusExpiry(StarwickerKey)).
			Write("starwicker_left", c.Tags[StarwickerKey])

		ae.Info.FlatDmg += amt
		c.AddStatus(a4ICDKey, 0.4*60, true)

		return false
	}, "xianyun-starwicker-hook")
}
