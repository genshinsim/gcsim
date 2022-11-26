package layla

import (
	"sort"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

var burstFrames []int

const burstStart = 36

func init() {
	burstFrames = frames.InitAbilSlice(79) // Q -> W
	burstFrames[action.ActionAttack] = 66  // Q -> N1
	burstFrames[action.ActionSkill] = 66   // Q -> E
	burstFrames[action.ActionDash] = 66    // Q -> D
	burstFrames[action.ActionJump] = 65    // Q -> J
	burstFrames[action.ActionSwap] = 64    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 22
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Starlight Slug",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		FlatDmg:    burst[c.TalentLvlBurst()] * c.MaxHP(),
	}
	// TODO: snapshot?
	snap := c.Snapshot(&ai)

	c.Core.Status.Add("laylaburst", 12*60+burstStart)

	x, y := c.Core.Combat.Player().Pos() // burst pos
	for delay := burstStart; delay < 12*60+burstStart; delay += 90 {
		c.Core.Tasks.Add(func() {
			nearTarget := -1
			trgs := c.Core.Combat.EnemiesWithinRadius(x, y, 12)
			if len(trgs) > 0 {
				sort.Slice(trgs, func(i, j int) bool { return i < j })
				nearTarget = trgs[0]
			}

			done := false
			cb := func(_ combat.AttackCB) {
				if done {
					return
				}
				done = true

				exist := c.Core.Player.Shields.Get(shield.ShieldLaylaSkill)
				if exist != nil {
					c.AddNightStars(1, ICDNightStarBurst)
				}
			}

			c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHit(c.Core.Combat.Enemy(nearTarget), 1.5), 0, cb)
		}, delay+travel)
	}

	c.SetCD(action.ActionBurst, 12*60)
	c.ConsumeEnergy(6)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}
