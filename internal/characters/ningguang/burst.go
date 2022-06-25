package ningguang

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var burstFrames []int

const burstHitmark = 97

func init() {
	burstFrames = frames.InitAbilSlice(97)
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Starshatter",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}

	// TODO: hitmark timing
	// fires 6 normally
	// geo applied 1 4 7 10, +3 pattern; or 0 3 6 9
	for i := 0; i < 6; i++ {
		c.Core.QueueAttack(ai, combat.NewDefCircHit(0.1, false, combat.TargettableEnemy), burstHitmark, burstHitmark+travel)
	}
	// if jade screen is active add 6 jades
	if c.Core.Constructs.Destroy(c.lastScreen) {
		ai.Abil = "Starshatter (Jade Screen Gems)"
		for i := 6; i < 12; i++ {
			c.Core.QueueAttackWithSnap(ai, c.skillSnapshot, combat.NewDefCircHit(0.1, false, combat.TargettableEnemy), burstHitmark+travel)
		}
		// do we need to log this?
		c.Core.Log.NewEvent("extra 6 gems from jade screen", glog.LogCharacterEvent, c.Index)
	}

	if c.Base.Cons >= 6 {
		c.Tags["jade"] = 7
		c.Core.Log.NewEvent("c6 - adding star jade", glog.LogCharacterEvent, c.Index, "count", c.Tags["jade"])
	}

	c.ConsumeEnergy(8)
	c.SetCDWithDelay(action.ActionBurst, 720, 8)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstHitmark,
		State:           action.BurstState,
	}
}
