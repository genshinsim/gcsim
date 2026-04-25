package aino

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var burstFrames []int

const (
	burstStart   = 123
	burstKey     = "aino-q"
	burstMarkKey = "aino-burst-mark"
)

func init() {
	burstFrames = frames.InitAbilSlice(60) // Q -> W
	burstFrames[action.ActionAttack] = 55  // Q -> N1
	burstFrames[action.ActionCharge] = 55  // Q -> CA
	burstFrames[action.ActionSkill] = 55   // Q -> E
	burstFrames[action.ActionDash] = 56    // Q -> D
	burstFrames[action.ActionJump] = 56    // Q -> J
	burstFrames[action.ActionSwap] = 54    // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	tickrate, radius, icdGroup, icdTag := c.a1BurstEnhance()
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Cool Your Jets Ducky",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     icdTag,
		ICDGroup:   icdGroup,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	burstArea := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 10)

	for i := 0; i < 14.2*60; i += tickrate {
		c.Core.Tasks.Add(func() {
			// burst tick
			enemy := c.Core.Combat.RandomEnemyWithinArea(
				burstArea,
				func(e info.Enemy) bool {
					return !e.StatusIsActive(burstMarkKey)
				},
			)
			var pos info.Point
			if enemy != nil {
				pos = enemy.Pos()
				enemy.AddStatus(burstMarkKey, 0.8*60, true) // same enemy can't be targeted again for 0.8s
			} else {
				pos = info.CalcRandomPointFromCenter(burstArea.Shape.Pos(), 1.5, 9, c.Core.Rand)
			}
			ai.FlatDmg = c.a4Dmg()
			// TODO: Aino burst travel time
			c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(pos, nil, radius), 0, 10)
		}, i+burstStart)
	}
	c.QueueCharTask(func() { c.AddStatus(burstKey, 14*60, false) }, burstStart)

	c.SetCD(action.ActionBurst, 13.5*60)
	c.ConsumeEnergy(5)

	c.c1OnSkillBurst()
	c.c6OnBurst()
	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
