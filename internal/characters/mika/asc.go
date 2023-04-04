package mika

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a1Stacks = "detector-stacks"
	a1Buff   = "detector-buff"
)

func (c *char) addDetectorStack() {
	stacks := c.Tag(a1Stacks)

	if stacks < c.maxDeterctorStacks {
		stacks++
		c.Core.Log.NewEvent("add detector stack", glog.LogCharacterEvent, c.Index).
			Write("stacks", stacks).
			Write("maxstacks", c.maxDeterctorStacks)
	}
	c.SetTag(a1Stacks, stacks)
}

// Per the following circumstances, the Soulwind state caused by Starfrost Swirl will grant characters the Detector effect,
// increasing their Physical DMG by 10% when they are on the field.
// - If the Flowfrost Arrow hits more than one opponent, each additional opponent hit will generate 1 Detector stack.
// - When a Rimestar Shard hits an opponent, it will generate 1 Detector stack. Each Rimestar Shard can trigger the effect 1 time.
//
// The Soulwind state can have a maximum of 3 Detector stacks, and if Starfrost Swirl is cast again during this duration, the pre-existing
// Soulwind state and all its Detector stacks will be cleared.
func (c *char) a1() {
	m := make([]float64, attributes.EndStatType)
	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(a1Buff, -1),
			AffectedStat: attributes.PhyP,
			Amount: func() ([]float64, bool) {
				if !char.StatusIsActive(skillBuffKey) {
					c.SetTag(a1Stacks, 0)
					return nil, false
				}
				m[attributes.PhyP] = 0.1 * float64(c.Tag(a1Stacks))
				return m, true
			},
		})

	}
}

// When an active character affected by both Skyfeather Song's Eagleplume and Starfrost Swirl's Soulwind at once scores a CRIT Hit with their
// attacks, Soulwind will grant them 1 stack of Detector from Suppressive Barrage. During a single instance of Soulwind, 1 Detector stack
// can be gained in this manner.
// Additionally, the maximum number of stacks that can be gained through Soulwind alone is increased by 1.
// Requires Suppressive Barrage to be unlocked first.
func (c *char) a4() {
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		if c.a4Stack {
			return false
		}

		atk := args[1].(*combat.AttackEvent)
		char := c.Core.Player.ByIndex(atk.Info.ActorIndex)
		if char.Index != c.Core.Player.Active() {
			return false
		}

		if !char.StatModIsActive(skillBuffKey) || !c.StatusIsActive(healKey) {
			return false
		}

		crit := args[3].(bool)
		if !crit {
			return false
		}

		c.addDetectorStack()
		c.a4Stack = true
		return false
	}, "mika-a4")
}
