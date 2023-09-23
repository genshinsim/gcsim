package ningguang

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var (
	burstFrames   []int
	burstHitmarks = []int{62, 97, 106, 110, 116, 124}
)

func init() {
	burstFrames = frames.InitAbilSlice(127)
	burstFrames[action.ActionDash] = 99
	burstFrames[action.ActionJump] = 100
	burstFrames[action.ActionWalk] = 99
	burstFrames[action.ActionSwap] = 98
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 0
	}

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Starshatter",
		AttackTag:          attacks.AttackTagElementalBurst,
		ICDTag:             attacks.ICDTagElementalBurst,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		Element:            attributes.Geo,
		Durability:         50,
		Mult:               burst[c.TalentLvlBurst()],
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}
	ap := combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 0.5)

	// fires 6 normally
	jade, ok := p["jade"]
	if !ok {
		jade = 6
	}
	for i := 0; i < jade; i++ {
		c.Core.QueueAttack(ai, ap, 0, burstHitmarks[i]+travel)
	}
	// if jade screen is active add 6 jades
	if c.Core.Constructs.Destroy(c.lastScreen) {
		screen, ok := p["screen"]
		if !ok {
			screen = 6
		}
		ai.Abil = "Starshatter (Jade Screen Gems)"
		for i := 0; i < screen; i++ {
			c.Core.QueueAttackWithSnap(ai, c.skillSnapshot, ap, burstHitmarks[len(burstHitmarks)-1]+30+travel) // TODO: figure out jade screen hitmarks
		}
	}

	if c.Base.Cons >= 6 {
		c.jadeCount = 7
		c.Core.Log.NewEvent("c6 - adding star jade", glog.LogCharacterEvent, c.Index).
			Write("count", c.jadeCount)
	}

	c.ConsumeEnergy(3)
	c.SetCD(action.ActionBurst, 720)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap],
		State:           action.BurstState,
	}, nil
}
