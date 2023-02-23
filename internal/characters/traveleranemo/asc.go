package traveleranemo

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

// The last hit of a Normal Attack combo unleashes a wind blade, dealing 60% of ATK as Anemo DMG to all opponents in its path.
func (c *char) a1() {
	if c.Base.Ascension < 1 || c.NormalCounter != c.NormalHitNum-1 {
		return
	}
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Slitting Wind (A1)",
		AttackTag:  attacks.AttackTagNormal,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupPoleExtraAttack,
		StrikeType: combat.StrikeTypeSlash,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       0.6,
	}
	c.QueueCharTask(func() {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				nil,
				1,
			),
			0,
			0,
		)
	}, a1Hitmark[c.gender])
}

const a4ICDKey = "traveleranemo-a4-icd"

// Palm Vortex kills regenerate 2% HP for 5s.
// This effect can only occur once every 5s.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	c.Core.Events.Subscribe(event.OnTargetDied, func(args ...interface{}) bool {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagElementalArt {
			return false
		}
		if c.StatusIsActive(a4ICDKey) {
			return false
		}

		c.AddStatus(a4ICDKey, 300, true)

		for i := 0; i < 5; i = i + 1 {
			c.QueueCharTask(func() {
				c.Core.Player.Heal(player.HealInfo{
					Caller:  c.Index,
					Target:  c.Index,
					Message: "Second Wind",
					Type:    player.HealTypePercent,
					Src:     0.02,
				})
			}, (i+1)*60) // healing starts 1s after death
		}

		return false
	}, "traveleranemo-a4")
}
