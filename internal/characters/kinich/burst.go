package kinich

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var burstFrames []int
var burstHitMarks = []int{30, 20, 40, 40, 40, 40, 40}

const (
	consumeEnergyDelay = 5
)

func init() {
	burstFrames = frames.InitAbilSlice(58) // Q
	burstFrames[action.ActionAttack] = 57  // Q -> N1
	burstFrames[action.ActionSkill] = 57   // Q -> E
	burstFrames[action.ActionJump] = 57    // Q -> J
	burstFrames[action.ActionSwap] = 56    // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	c.ajawSrc = c.Core.F
	if c.nightsoulState.HasBlessing() {
		c.durationExtended = true
		c.resetNightsoulExitTimer((10 + 1.7) * 60)
	}
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Hail to the Almighty Dragonlord (Skill DMG)",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), geometry.Point{Y: 0.2}, 4), burstHitMarks[0], burstHitMarks[0])
	c.Core.Tasks.Add(func() { c.QueueLaser(2, c.Core.F) }, burstHitMarks[1])
	c.ConsumeEnergy(consumeEnergyDelay)
	c.SetCDWithDelay(action.ActionBurst, 18*60, 1)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) QueueLaser(step, src int) {
	if src != c.ajawSrc {
		return
	}
	// duration expired
	if c.Core.F-c.ajawSrc > 15*60 {
		return
	}
	// condition to track number of hits just in case
	if step == 7 {
		c.ajawSrc = -1
		return
	}
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Hail to the Almighty Dragonlord (Dragon Breath DMG)",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	c.Core.QueueAttack(ai, combat.NewSingleTargetHit(c.Core.Combat.PrimaryTarget().Key()), 0, 0)
	c.Core.Tasks.Add(func() { c.QueueLaser(step+1, src) }, burstHitMarks[step])
}

func (c *char) resetNightsoulExitTimer(duration int) {
	src := c.Core.F
	timePassed := src - c.nightsoulSrc
	duration -= timePassed
	c.QueueCharTask(func() {
		c.durationExtended = false
		c.cancelNightsoul()
	}, duration)
}
