package ororon

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var burstFrames []int

var burstTicks = []int{69, 63, 60, 60, 60, 60, 60, 60, 60}

const burstHitmark = 36 // Initial Hit

func init() {
	burstFrames = frames.InitAbilSlice(62) // Q -> N1
	burstFrames[action.ActionDash] = 59
	burstFrames[action.ActionJump] = 60
	burstFrames[action.ActionSwap] = 61
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	// first zap has no icd and hits everyone
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Ritual DMG",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 2.3}, 4, 4), // TODO: confirm box size
		burstHitmark,
		burstHitmark,
		c.makeC2cb(),
	)

	ai = combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Soundwave Collision DMG",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagOroronElmentalBurst,
		ICDGroup:   attacks.ICDGroupOroronElementalBurst,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       soundwave[c.TalentLvlBurst()],
	}

	progress := 0
	for i := 0; i < 9; i++ {
		progress += int(float64(burstTicks[i]) * c.c4BurstInterval())
		c.QueueCharTask(func() {
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHitFanAngle(
					c.Core.Combat.Player(),
					c.Core.Combat.PrimaryTarget(),
					geometry.Point{Y: 1},
					10, 20,
				), // TODO: confirm size
				0,
				0,
				c.makeC2cb(),
			)
		}, progress)
	}
	c.c2OnBurst()
	c.c6OnBurst()
	c.SetCDWithDelay(action.ActionBurst, 15*60, 0)
	c.ConsumeEnergy(21)
	c.c4EnergyRestore()
	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}, nil
}
