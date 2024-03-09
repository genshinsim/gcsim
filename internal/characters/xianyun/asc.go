package xianyun

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
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
		i := idx
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("xianyun-a1-buff", -1),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if atk.Info.AttackTag != attacks.AttackTagPlunge {
					return nil, false
				}
				stackCount := min(c.a1Buffer[i], 4)
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
			char.AddStatus(a1Key, a1Dur, true)
			char.SetTag(a1Key, min(c.a1Buffer[idx], 4))
			char.QueueCharTask(func() {
				// tags currently aren't visible in the results UI
				// the user can still access it using .char.tags.xianyun-a1
				c.a1Buffer[idx] -= 1
				char.SetTag(a1Key, min(c.a1Buffer[idx], 4))
			}, a1Dur)
		}
	}
}

func (c *char) a4StartUpdate() {
	if c.Base.Ascension < 4 {
		return
	}

	c.a4src = c.Core.F
	c.Core.Tasks.Add(c.a4AtkUpdate(c.Core.F), 0)
}

func (c *char) a4AtkUpdate(src int) func() {
	return func() {
		if c.a4src != src {
			return
		}
		c.a4Atk = c.getTotalAtk()
		c.Core.Tasks.Add(c.a4AtkUpdate(src), 0.5*60)
	}
}

// a4: When the Starwicker created by Stars Gather at Dusk has Adeptal Assistance stacks,
// nearby active characters' Plunging Attack shockwave DMG will be increased by 200% of Xianyun's ATK.
// The maximum DMG increase that can be achieved this way is 9000.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	c.a4Max = 9000
	c.a4Ratio = 2.0

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		if ae.Info.AttackTag != attacks.AttackTagPlunge {
			return false
		}

		// Collision has 0 durability. Don't buff collision damage
		if ae.Info.Durability == 0 {
			return false
		}

		if !c.StatusIsActive(player.XianyunAirborneBuff) {
			return false
		}

		if c.StatusIsActive(a4ICDKey) {
			return false
		}

		// A4 cap
		amt := min(c.a4Max, c.a4Ratio*c.a4Atk)

		c.Core.Log.NewEvent("Xianyun A4 proc dmg add", glog.LogPreDamageMod, ae.Info.ActorIndex).
			Write("atk", c.a4Atk).
			Write("ratio", c.a4Ratio).
			Write("addition", amt)

		ae.Info.FlatDmg += amt
		c.AddStatus(a4ICDKey, 0.4*60, true)

		return false
	}, "xianyun-starwicker-hook")
}
