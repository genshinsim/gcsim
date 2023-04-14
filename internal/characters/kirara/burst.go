package kirara

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

// based on klee frames
// TODO: update frames
var (
	burstFrames []int

	boxHitmark  = 19 * 2
	mineHitmark = 240
	mineExpired = "kirara-cardamoms-expired"
)

func init() {
	burstFrames = frames.InitAbilSlice(33 * 2)
}

// Has one parameter, "hits" determines the number of cardamoms that hit the enemy
func (c *char) Burst(p map[string]int) action.ActionInfo {
	boxAi := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Secret Art: Surprise Dispatch",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Dendro,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}

	c.cardamoms = 6
	if c.Base.Cons >= 1 {
		// Every 8,000 Max HP Kirara possesses will cause her to create 1 extra Cat Grass Cardamom when she uses Secret Art: Surprise Dispatch.
		// A maximum of 4 extra can be created this way.
		bonus := int(c.MaxHP() / 8000)
		if bonus > 4 {
			bonus = 4
		}
		c.cardamoms += bonus
	}

	// TODO: gadgets?
	minehits, ok := p["hits"]
	if !ok {
		minehits = 2
	}
	if minehits > c.cardamoms {
		minehits = c.cardamoms
	}
	mineAi := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Cat Grass Cardamom Explosion",
		AttackTag:          attacks.AttackTagElementalBurst,
		ICDTag:             attacks.ICDTagElementalBurst,
		ICDGroup:           attacks.ICDGroupDefault, // TODO: Mine??
		StrikeType:         attacks.StrikeTypeDefault,
		Element:            attributes.Dendro,
		Durability:         25,
		Mult:               cardamom[c.TalentLvlBurst()],
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}

	c.mineSnap = c.Snapshot(&mineAi)
	c.minePattern = combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 2)

	// box
	c.QueueCharTask(func() {
		c.AddStatus(mineExpired, 12*60, true)
		c.Core.QueueAttackWithSnap(boxAi, c.mineSnap, combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 4), 0)
	}, boxHitmark)

	// mine hits
	c.QueueCharTask(func() {
		c.cardamoms -= minehits
		if c.cardamoms < 0 {
			c.cardamoms = 0
		}
		c.Core.QueueAttackWithSnap(mineAi, c.mineSnap, c.minePattern, 0)
	}, mineHitmark)

	// mine expires
	c.QueueCharTask(func() {
		for i := 0; i < c.cardamoms; i++ {
			c.Core.QueueAttackWithSnap(mineAi, c.mineSnap, c.minePattern, i*9*2)
		}
		c.cardamoms = 0
	}, boxHitmark+12*60)

	if c.Base.Cons >= 6 {
		c.c6()
	}

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(12)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash],
		State:           action.BurstState,
	}
}
