package travelerdendro

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const a1Key = "dmc-a1"

// Gets removed on swap - from Kolibri
func (c *char) a1Init() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		prev := args[0].(int)
		prevChar := c.Core.Player.ByIndex(prev)
		prevChar.DeleteStatMod(a1Key)
		return false
	}, "dmc-a1-remove")
}

func (c *char) a1Buff(delay int) {
	m := make([]float64, attributes.EndStatType)
	// A1/C6 buff ticks every 0.3s and applies for 1s. probably counting from gadget spawn - from Kolibri
	c.Core.Tasks.Add(func() {
		if c.Core.Status.Duration(burstKey) <= 0 {
			return
		}
		if !combat.TargetIsWithinArea(c.Core.Combat.Player(), combat.NewCircleHitOnTarget(c.burstPos, nil, c.burstRadius)) {
			return
		}
		m[attributes.EM] = float64(6 * c.burstOverflowingLotuslight)
		active := c.Core.Player.ActiveChar()
		active.AddStatMod(character.StatMod{
			Base:         modifier.NewBase(a1Key, 60),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}, delay)
}

func (c *char) a1Stack(delay int) {
	c.Core.Tasks.Add(func() {
		if c.Core.Status.Duration(burstKey) > 0 && c.burstOverflowingLotuslight < 10 { //burst isn't expired, and stacks aren't capped
			c.burstOverflowingLotuslight += 1
		}
	}, delay)
}

// Every point of Elemental Mastery the Traveler possesses increases the DMG dealt
// by Razorgrass Blade by 0.15% and the DMG dealt by Surgent Manifestation by 0.1%.
func (c *char) a4Init() {
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag("dmc-a4", -1),
		Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			switch atk.Info.AttackTag {
			case combat.AttackTagElementalArt:
				m[attributes.DmgP] = c.Stat(attributes.EM) * 0.0015
				return m, true
			case combat.AttackTagElementalBurst:
				m[attributes.DmgP] = c.Stat(attributes.EM) * 0.001
				return m, true
			default:
				return nil, false
			}
		},
	})
}
