package columbina

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

func init() {
	burstFrames = frames.InitAbilSlice(120)
	burstFrames[action.ActionSwap] = 120
}

const (
	burstKey     = "columbina-q"
	burstBuffKey = "columbina-q-buff"
	withinTimer  = 1.2 * 60
	burstDur     = 20 * 60
	burstHitmark = 105
)

func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Burst",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		UseHP:      true,
		Mult:       burst[c.TalentLvlBurst()],
	}
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6.5)
	// TODO: is the field the same size as the hitbox?
	c.Core.QueueAttack(ai, ap, burstHitmark, burstHitmark)

	// how burst works:
	// every 0.8s, it checks if the characters are in the field
	// if the character is in the field, they get a "in_field" status
	// while they have the "in_field" status, they also get the lunar reaction buff

	c.burstArea = combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: 1}, 20)
	c.Core.Tasks.Add(func() {
		c.burstSrc = c.Core.F
		src := c.Core.F
		c.Core.Status.Add(burstKey, burstDur)
		for i := 0; i <= burstDur; i += 0.8 * 60 {
			c.Core.Tasks.Add(func() {
				// don't tick if another burst has already started
				if c.burstSrc != src {
					return
				}

				// don't apply anything if outside of burst area
				if !c.Core.Combat.Player().IsWithinArea(c.burstArea) {
					return
				}

				c.applyBurstBuff()
			}, i)
		}
	}, burstHitmark)

	c.ConsumeEnergy(5)
	c.SetCD(action.ActionBurst, 15*60)
	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) applyBurstBuff() {
	for _, char := range c.Core.Player.Chars() {
		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBaseWithHitlag(burstBuffKey, withinTimer),
			Amount: func(ai info.AttackInfo) float64 {
				if !attacks.AttackTagIsLunar(ai.AttackTag) {
					return 0
				}

				if c.Core.Combat.Debug {
					c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, char.Index(), "Adding columbina burst react bonus")
				}
				return burstBuff[c.TalentLvlBurst()]
			},
		})
	}
}
