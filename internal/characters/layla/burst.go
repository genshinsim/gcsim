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

const (
	burstStart = 36

	tickRelease = 56
	tickTravel  = 22
)

func init() {
	burstFrames = frames.InitAbilSlice(79) // Q -> W
	burstFrames[action.ActionAttack] = 65  // Q -> N1
	burstFrames[action.ActionSkill] = 66   // Q -> E
	burstFrames[action.ActionDash] = 67    // Q -> D
	burstFrames[action.ActionJump] = 66    // Q -> J
	burstFrames[action.ActionSwap] = 65    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = tickTravel
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

	c.Core.Status.Add("laylaburst", 12*60+burstStart)

	player := c.Core.Combat.Player()
	burstPos := combat.CalcOffsetPoint(player.Pos(), combat.Point{Y: 1}, player.Direction()) // burst pos
	for delay := burstStart; delay < 12*60+burstStart; delay += 90 {
		c.Core.Tasks.Add(func() {
			trgs := c.Core.Combat.EnemiesWithinRadius(burstPos, 12)
			if len(trgs) == 0 {
				return
			}
			sort.Slice(trgs, func(i, j int) bool { return i < j })
			nearTarget := trgs[0]

			done := false
			cb := func(_ combat.AttackCB) {
				if done {
					return
				}
				done = true

				exist := c.Core.Player.Shields.Get(shield.ShieldLaylaSkill)
				if exist != nil {
					c.addNightStars(1, ICDNightStarBurst)
				}
			}

			c.Core.QueueAttack(
				ai,
				combat.NewCircleHit(
					c.Core.Combat.Player(),
					c.Core.Combat.Enemy(nearTarget),
					nil,
					1.5,
				),
				tickRelease,
				tickRelease+travel,
				cb,
			)
		}, delay)
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
