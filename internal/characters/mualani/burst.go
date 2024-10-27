package mualani

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

const burstHitmarks = 185 - 70

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
		FlatDmg:        burst[c.TalentLvlBurst()] * c.MaxHP(),
	}
	burstArea := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 5)

	c.QueueCharTask(func() {
		// the A4 stacks can change during the burst
		ai.FlatDmg += c.a4amount()
		c.Core.QueueAttack(ai, burstArea, 0, 0)
	}, burstHitmarks+travel)

	c.SetCDWithDelay(action.ActionBurst, 15*60, 0)
	c.ConsumeEnergy(11)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
