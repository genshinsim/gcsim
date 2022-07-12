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
	skillFrames = frames.InitAbilSlice(58)
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
			c.Core.QueueParticle("kaeya", 1, attributes.Cryo, c.Core.Flags.ParticleDelay)
			c.Core.Log.NewEvent("kaeya a4 proc", glog.LogCharacterEvent, c.Index)
		}
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(1, false, combat.TargettableEnemy), 0, skillHitmark, cb)

	//2 or 3 1:1 ratio
	var count float64 = 2
	if c.Core.Rand.Float64() < 0.67 {
		count = 3
	}
	c.Core.QueueParticle("kaeya", count, attributes.Cryo, skillHitmark+c.Core.Flags.ParticleDelay)

	c.SetCD(action.ActionSkill, 360+28) //+28 since cd starts 28 frames in

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillHitmark,
		State:           action.SkillState,
	}
}
