package klee

import (
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var (
	skillFrames         []int
	skillCanceledFrames []int
)

var bounceHitmarks = []int{71, 111, 140}

const mineHitmark = 240

func init() {
	skillFrames = frames.InitAbilSlice(75)
	skillFrames[action.ActionAttack] = 66
	skillFrames[action.ActionCharge] = 69
	skillFrames[action.ActionSkill] = 68
	skillFrames[action.ActionBurst] = 34
	skillFrames[action.ActionDash] = 37
	skillFrames[action.ActionJump] = 35
	skillFrames[action.ActionSwap] = 74
	skillCanceledFrames = make([]int, len(skillFrames))
	copy(skillCanceledFrames, skillFrames)
	skillCanceledFrames[action.ActionBurst] = 5
}

// Has two parameters, "bounce" determines the number of bounces that hit
// "mine" determines the number of mines that hit the enemy
func (c *char) Skill(p map[string]int) action.ActionInfo {
	type attackData struct {
		ai   combat.AttackInfo
		snap combat.Snapshot
	}
	bounce, ok := p["bounce"]
	if !ok {
		bounce = 1
	}
	bounceAttacks := make([]attackData, bounce)
	for i := range bounceAttacks {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Jumpy Dumpty",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     combat.ICDTagKleeFireDamage,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeBlunt,
			Element:    attributes.Pyro,
			Durability: 25,
			Mult:       jumpy[c.TalentLvlSkill()],
		}
		// 3rd bounce is 2B
		if i == 2 {
			ai.Durability = 50
		}
		bounceAttacks[i] = attackData{
			ai:   ai,
			snap: c.Snapshot(&ai),
		}
	}
	minehits, ok := p["mine"]
	if !ok {
		minehits = 2
	}
	mineAttacks := make([]attackData, minehits)
	mineAi := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Jumpy Dumpty Mine Hit",
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             combat.ICDTagKleeFireDamage,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeDefault,
		Element:            attributes.Pyro,
		Durability:         25,
		Mult:               mine[c.TalentLvlSkill()],
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}
	for i := range mineAttacks {
		mineAttacks[i] = attackData{
			ai:   mineAi,
			snap: c.Snapshot(&mineAi),
		}
	}

	release, ok := p["release"]
	if !ok {
		release = 1
	}

	adjustedFrames := skillFrames
	if release == 0 {
		adjustedFrames = skillCanceledFrames
	}

	canQueueAfter := math.MaxInt32
	for _, f := range adjustedFrames {
		if f < canQueueAfter {
			canQueueAfter = f
		}
	}
	actionInfo := action.ActionInfo{
		Frames:          frames.NewAbilFunc(adjustedFrames),
		AnimationLength: adjustedFrames[action.InvalidAction],
		CanQueueAfter:   canQueueAfter,
		State:           action.SkillState,
	}
	cooldownDelay := 33
	actionInfo.QueueAction(func() {
		if release == 0 {
			c.Core.Log.NewEvent("attempted klee skill cancel without burst", glog.LogWarnings, -1)
		}
		particleCB := c.makeParticleCB()
		for i, data := range bounceAttacks {
			c.Core.QueueAttackWithSnap(
				data.ai,
				data.snap,
				combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 4),
				bounceHitmarks[i]-cooldownDelay,
				c.makeA1CB(),
				particleCB,
			)
		}
		for _, data := range mineAttacks {
			c.Core.QueueAttackWithSnap(
				data.ai,
				data.snap,
				combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 2),
				mineHitmark-cooldownDelay,
				c.c2,
			)
		}
		c.c1(bounceHitmarks[0] - cooldownDelay)
		c.SetCD(action.ActionSkill, 1200)
	}, cooldownDelay)
	return actionInfo
}

func (c *char) makeParticleCB() combat.AttackCBFunc {
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != combat.TargettableEnemy {
			return
		}
		if done {
			return
		}
		done = true
		c.Core.QueueParticle(c.Base.Key.String(), 4, attributes.Pyro, c.ParticleDelay)
	}
}
