package mualani

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

const burstHitmarks = 108 // adjusted to swap frame

var (
	burstFrames []int
)

func init() {
	burstFrames = frames.InitAbilSlice(180) // charge
	burstFrames[action.ActionAttack] = 167
	burstFrames[action.ActionSkill] = 166
	burstFrames[action.ActionDash] = 167
	burstFrames[action.ActionJump] = 167
	burstFrames[action.ActionWalk] = 167
	burstFrames[action.ActionSwap] = 108
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 70
	}

	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Boomsharka-laka",
		AttackTag:      attacks.AttackTagElementalBurst,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Hydro,
		Durability:     25,
	}
	burstArea := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 5)

	// snapshot at bullet creation
	var snap combat.Snapshot
	stacks := c.a4Stacks
	c.a4Stacks = 0
	c.QueueCharTask(func() {
		snap = c.Snapshot(&ai)
		c.Core.Tasks.Add(func() {
			// TODO: verify if snapshot is used or if maxhp is recalced here
			hp := c.MaxHP()
			ai.FlatDmg = burst[c.TalentLvlBurst()] * hp
			if c.Base.Ascension >= 4 {
				ai.FlatDmg += 0.15 * float64(stacks) * hp
			}

			c.Core.QueueAttackWithSnap(ai, snap, burstArea, 0)
		}, travel)
	}, burstHitmarks)

	c.SetCDWithDelay(action.ActionBurst, 15*60, 0)
	c.ConsumeEnergy(11)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
