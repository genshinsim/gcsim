package mika

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var (
	burstFrames []int
	boxHitmark  = 38
	mineExpired = "secondary-grenades-expired"
)

func init() {
	burstFrames = frames.InitAbilSlice(61) // Q -> N1/Dash/Walk
	burstFrames[action.ActionSkill] = 60
	burstFrames[action.ActionJump] = 60
	burstFrames[action.ActionSwap] = 59
}

func (c *char) Burst(p map[string]int) (action.Info, error) {

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Explosive Grenade",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeSlash,
		Element:    attributes.Pyro,
		Durability: 25,
		FlatDmg:    c.MaxHP() * burst[c.TalentLvlBurst()],
	}

	mineAi := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Secondary Explosive Shell",
		AttackTag:          attacks.AttackTagElementalBurst,
		ICDTag:             attacks.ICDTagElementalBurst,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeDefault,
		Element:            attributes.Pyro,
		Durability:         25,
		Mult:               burstSecondary[c.TalentLvlBurst()],
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}

	player := c.Core.Combat.Player()
	boxPos := geometry.CalcOffsetPoint(player.Pos(), geometry.Point{Y: 3}, player.Direction())
	c.QueueCharTask(func() {
		c.AddStatus(mineExpired, 12*60, true)
		c.Core.QueueAttackWithSnap(ai, c.mineSnap, combat.NewCircleHitOnTarget(boxPos, nil, 6), 0)
	}, boxHitmark)

	c.minePattern = combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 2)
	shellNum := 2
	mineDuration := 4
	c.QueueCharTask(func() {
		for i := 0; i < shellNum; i++ {
			c.Core.QueueAttackWithSnap(mineAi, c.mineSnap, c.minePattern, i*9*2)
		}
		c.secondaryMineNum = 0
	}, 80+mineDuration*60)

	c.ConsumeEnergy(4)
	c.SetCD(action.ActionBurst, 15*60)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
