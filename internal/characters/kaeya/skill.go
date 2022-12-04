package kaeya

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var skillFrames []int

const skillHitmark = 28

func init() {
	skillFrames = frames.InitAbilSlice(53) // E -> N1
	skillFrames[action.ActionBurst] = 52   // E -> Q
	skillFrames[action.ActionDash] = 25    // E -> D
	skillFrames[action.ActionJump] = 26    // E -> J
	skillFrames[action.ActionSwap] = 49    // E -> Swap
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Frostgnaw",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 50,
		Mult:       skill[c.TalentLvlSkill()],
	}
	a4count := 0
	cb := func(a combat.AttackCB) {
		e, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}

		heal := .15 * (a.AttackEvent.Snapshot.BaseAtk*(1+a.AttackEvent.Snapshot.Stats[attributes.ATKP]) + a.AttackEvent.Snapshot.Stats[attributes.ATK])
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  c.Core.Player.Active(),
			Message: "Cold-Blooded Strike",
			Src:     heal,
			Bonus:   c.Stat(attributes.Heal),
		})
		//if target is frozen after hit then drop additional energy;
		if a4count == 2 {
			return
		}
		if e.AuraContains(attributes.Frozen) {
			a4count++
			c.Core.QueueParticle("kaeya", 1, attributes.Cryo, c.ParticleDelay)
			c.Core.Log.NewEvent("kaeya a4 proc", glog.LogCharacterEvent, c.Index)
		}
	}
	c.Core.QueueAttack(
		ai,
		combat.NewBoxHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: -0.2}, 4, 8),
		0,
		skillHitmark,
		cb,
	)

	// 2 or 3, 1:2 ratio
	var count float64 = 2
	if c.Core.Rand.Float64() < 0.67 {
		count = 3
	}
	c.Core.QueueParticle("kaeya", count, attributes.Cryo, skillHitmark+c.ParticleDelay)

	c.SetCDWithDelay(action.ActionSkill, 360, 25)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel is before skillHitmark
		State:           action.SkillState,
	}
}
