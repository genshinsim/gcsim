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

const starwickerKey = "xianyun-starwicker-count"
const starwickerICDKey = "xianyun-starwicker-icd"

// a1: For every enemy hit by Driftcloud Wave,
// all party members gain 1 stack of Boost, lasting 20s, max 4 stacks.
// Boost increases Plunge DMG's Crit Rate by 4%/6%/8%/10%,
// each stack's duration is calculated independently.
func (c *char) a1() combat.AttackCBFunc {
	c.a1Count = 0
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		c.a1Count += 1
		if c.a1Count > 4 {
			c.a1Count = 4
		}
		boost := 0
		stackDuration := 20 * 60
		var crBuff float64 = 0.0
		if c.a1Count != 0 {
			crBuff = 0.02 + 0.02*float64(c.a1Count)
		}
		mCR := make([]float64, attributes.EndStatType)
		mCR[attributes.CR] = crBuff
		addStackMod := func(idx int, duration int) {
			for _, char := range c.Core.Player.Chars() {
				char.AddAttackMod(character.AttackMod{
					Base: modifier.NewBaseWithHitlag("Xianyun-a1-buff", stackDuration),
					Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
						if atk.Info.AttackTag != attacks.AttackTagPlunge {
							return nil, false
						}
						return mCR, true
					},
				})
			}

		}
		addStackMod(boost, 20)
		boost = (boost + 1) % 4
	}
}

// TODO: ?? plunge
// a4: When the Starwicker created by Stars Gather at Dusk has Adeptal Assistance stacks,
// nearby active characters' Plunging Attack shockwave DMG will be increased by 170% of Xianyun's ATK.
// The maximum DMG increase that can be achieved this way is 8,500.
// TODO: Each Plunging Attack shockwave DMG instance can only apply this increased DMG effect to a single opponent.
// Each character can trigger this effect once every 0.4s.
func (c *char) a4() {
	c.SetTag(starwickerKey, 8)
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if c.StatusIsActive(starwickerICDKey) {
			return false
		}
		if !c.StatusIsActive(burstKey) {
			return false
		}

		var a4ScalingATKp float64 = 170 / 100

		if c.Base.Cons >= 2 {
			a4ScalingATKp = 340 / 100
		}
		switch atk.Info.AttackTag {
		case attacks.AttackTagPlunge:
		default:
			return false
		}

		// char := c.Core.Player.ByIndex(atk.Info.ActorIndex)
		if !c.StatusIsActive(starwickerKey) {
			return false
		}

		if c.Tags[starwickerKey] > 0 {
			stats, _ := c.Stats()
			amt := a4ScalingATKp * ((c.Base.Atk+c.Weapon.BaseAtk)*(1+stats[attributes.ATKP]) + stats[attributes.ATK])

			//A4 cap
			if c.Base.Cons < 2 && amt >= 8500 {
				amt = 8500
			} else if amt >= 13500 {
				amt = 13500
			}

			c.Tags[starwickerKey]--

			if c.Core.Flags.LogDebug {
				c.Core.Log.NewEvent("Xianyun Starwicker proc dmg add", glog.LogPreDamageMod, atk.Info.ActorIndex).
					Write("before", atk.Info.FlatDmg).
					Write("addition", amt).
					Write("effect_ends_at", c.StatusExpiry(starwickerKey)).
					Write("starwicker_left", c.Tags[starwickerKey])
			}

			atk.Info.FlatDmg += amt
			c.AddStatus(starwickerICDKey, 0.4*60, true)
		}

		return false
	}, "xianyun-starwicker-hook")
}
