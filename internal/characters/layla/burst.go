package layla

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

var burstFrames []int

const burstStart = 130

func init() {
	burstFrames = frames.InitAbilSlice(78) // Q -> E
	burstFrames[action.ActionAttack] = 77  // Q -> N1
	burstFrames[action.ActionDash] = 62    // Q -> D
	burstFrames[action.ActionJump] = 61    // Q -> J
	burstFrames[action.ActionSwap] = 77    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 26
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Starlight Slug",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		FlatDmg:    burst[c.TalentLvlBurst()] * c.MaxHP(),
	}
	snap := c.Snapshot(&ai)

	c.Core.Status.Add("laylaburst", 12*60+burstStart)

	for delay := burstStart; delay < 12*60+burstStart; delay += 90 {
		c.Core.Tasks.Add(func() {
			// TODO: "When a Starlight Slug hits"
			exist := c.Core.Player.Shields.Get(shield.ShieldLaylaSkill)
			if exist != nil && !c.StatusIsActive(starBurstIcd) && !c.StatusIsActive(shootingStars) {
				c.AddNightStars(1, false)
				c.AddStatus(starBurstIcd, 0.5*60, true)
			}
			c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHit(c.Core.Combat.PrimaryTarget(), 2.5), 0)
		}, delay+travel)
	}

	c.SetCD(action.ActionBurst, 12*60)
	c.ConsumeEnergy(3)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}
