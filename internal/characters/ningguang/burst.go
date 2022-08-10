package ningguang

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var burstFrames []int

var (
	burstHitmarks       = []int{62, 97, 106, 110, 116, 124}
	burstScreenHitmarks []int
)

func init() {
	burstFrames = frames.InitAbilSlice(127)
	burstFrames[action.ActionDash] = 99
	burstFrames[action.ActionJump] = 100
	burstFrames[action.ActionWalk] = 99
	burstFrames[action.ActionSwap] = 98
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 0
	}

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Starshatter",
		AttackTag:          combat.AttackTagElementalBurst,
		ICDTag:             combat.ICDTagElementalBurst,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeBlunt,
		Element:            attributes.Geo,
		Durability:         50,
		Mult:               burst[c.TalentLvlBurst()],
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}

	// fires 6 normally
	for _, hitmark := range burstHitmarks {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), 0.1, false, combat.TargettableEnemy),
			hitmark,
			hitmark+travel,
		)
	}
	// if jade screen is active add 6 jades
	if c.Core.Constructs.Destroy(c.lastScreen) {
		ai.Abil = "Starshatter (Jade Screen Gems)"
		for i := 6; i < 12; i++ {
			c.Core.QueueAttackWithSnap(
				ai,
				c.skillSnapshot,
				combat.NewCircleHit(c.Core.Combat.Player(), 0.1, false, combat.TargettableEnemy),
				burstHitmarks[len(burstHitmarks)-1]+30+travel,
			) // TODO: figure out jade screen hitmarks
		}
	}

	if c.Base.Cons >= 6 {
		c.jadeCount = 7
		c.Core.Log.NewEvent("c6 - adding star jade", glog.LogCharacterEvent, c.Index).
			Write("count", c.jadeCount)
	}

	c.ConsumeEnergy(3)
	c.SetCD(action.ActionBurst, 720)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap],
		State:           action.BurstState,
	}
}
