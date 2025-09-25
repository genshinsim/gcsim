package kirara

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	burstFrames []int

	mineSnapshot = 34
	boxHitmark   = 38
	mineExpired  = "kirara-cardamoms-expired"
)

func init() {
	burstFrames = frames.InitAbilSlice(58) // Q -> Walk
	burstFrames[action.ActionAttack] = 57
	burstFrames[action.ActionSkill] = 57
	burstFrames[action.ActionDash] = 52
	burstFrames[action.ActionJump] = 52
	burstFrames[action.ActionSwap] = 57
}

// Has one parameter, "hits" determines the number of cardamoms that hit the enemy
func (c *char) Burst(p map[string]int) (action.Info, error) {
	boxAi := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Secret Art: Surprise Dispatch",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}

	c.cardamoms = 6
	if c.Base.Cons >= 1 {
		// Every 8,000 Max HP Kirara possesses will cause her to create 1 extra Cat Grass Cardamom when she uses Secret Art: Surprise Dispatch.
		// A maximum of 4 extra can be created this way.
		bonus := min(int(c.MaxHP()/8000), 4)
		c.cardamoms += bonus
	}

	minehits, ok := p["hits"]
	if !ok {
		minehits = 2
	}
	mineHitmark, ok := p["mine_delay"]
	if !ok {
		mineHitmark = 180
	}
	if minehits > c.cardamoms {
		minehits = c.cardamoms
	}
	mineAi := info.AttackInfo{
		ActorIndex:         c.Index(),
		Abil:               "Cat Grass Cardamom Explosion",
		AttackTag:          attacks.AttackTagElementalBurst,
		ICDTag:             attacks.ICDTagElementalBurst,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeDefault,
		Element:            attributes.Dendro,
		Durability:         25,
		Mult:               cardamom[c.TalentLvlBurst()],
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}

	c.QueueCharTask(func() {
		c.mineSnap = c.Snapshot(&mineAi)
	}, mineSnapshot)

	c.minePattern = combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 2)

	// box
	player := c.Core.Combat.Player()
	boxPos := info.CalcOffsetPoint(player.Pos(), info.Point{Y: 3}, player.Direction())
	c.QueueCharTask(func() {
		c.AddStatus(mineExpired, 12*60, true)
		c.Core.QueueAttack(boxAi, combat.NewCircleHitOnTarget(boxPos, nil, 6), 0, 0)
	}, boxHitmark)

	// mine hits
	c.QueueCharTask(func() {
		for i := 0; i < minehits; i++ {
			c.Core.QueueAttackWithSnap(mineAi, c.mineSnap, c.minePattern, 0)
		}
		c.cardamoms -= minehits
		if c.cardamoms < 0 {
			c.cardamoms = 0
		}
	}, mineHitmark)

	// mine expires
	c.QueueCharTask(func() {
		for i := 0; i < c.cardamoms; i++ {
			c.Core.QueueAttackWithSnap(mineAi, c.mineSnap, c.minePattern, i*9*2)
		}
		c.cardamoms = 0
	}, 80+12*60)

	if c.Base.Cons >= 6 {
		c.c6()
	}

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(7)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash],
		State:           action.BurstState,
	}, nil
}
