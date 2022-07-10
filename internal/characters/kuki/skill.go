package kuki

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var skillFrames []int

const skillHitmark = 58

func init() {
	skillFrames = frames.InitAbilSlice(58)
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	//remove some hp
	if 0.7*(c.HPCurrent/c.MaxHP()) > 0.2 {
		c.HPCurrent = 0.7 * c.HPCurrent
	} else if (c.HPCurrent / c.MaxHP()) > 0.2 { //check if below 20%
		c.HPCurrent = 0.2 * c.MaxHP()
	}
	//TODO: damage frame

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Sanctifying Ring",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypePierce,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
		FlatDmg:    c.Stat(attributes.EM) * 0.25,
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(1, false, combat.TargettableEnemy), skillHitmark, skillHitmark)

	// C2: Grass Ring of Sanctification's duration is increased by 3s.
	skilldur := 720
	if c.Base.Cons >= 2 {
		skilldur = 900 //12+3s
	}

	c.SetCD(action.ActionSkill, skillHitmark+15*60) // what's the diff between f and a again? Nice question Yakult
	c.Core.Tasks.Add(c.bellTick(), 90)              //Assuming this executes every 90 frames-1.5s
	c.bellActiveUntil = c.Core.F + skilldur
	c.Core.Log.NewEvent("Bell activated", glog.LogCharacterEvent, c.Index).
		Write("expected end", c.bellActiveUntil).
		Write("next expected tick", c.Core.F+90)

	c.Core.Status.Add("kukibell", skilldur)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.InvalidAction],
		State:           action.SkillState,
	}
}

func (c *char) bellTick() func() {
	return func() {
		c.Core.Log.NewEvent("Bell ticking", glog.LogCharacterEvent, c.Index)

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Grass Ring of Sanctification",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagElementalArt,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypePierce,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       skilldot[c.TalentLvlSkill()],
			FlatDmg:    c.Stat(attributes.EM) * 0.25,
		}
		c.Core.QueueAttack(ai, combat.NewDefCircHit(1, false, combat.TargettableEnemy), 2, 2)

		//A4 is considered here
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  c.Core.Player.Active(),
			Message: "Grass Ring of Sanctification Healing",
			Src:     (skillhealpp[c.TalentLvlSkill()]*c.MaxHP() + skillhealflat[c.TalentLvlSkill()] + c.Stat(attributes.EM)*0.75),
			Bonus:   c.Stat(attributes.Heal),
		})

		c.Core.Log.NewEvent("Bell ticked", glog.LogCharacterEvent, c.Index).
			Write("next expected tick", c.Core.F+90).
			Write("active", c.bellActiveUntil)
		//trigger damage
		//TODO: Check for snapshots

		//c.Core.QueueAttackEvent(&ae, 0)
		//check for orb
		//Particle check is 45% for particle
		if c.Core.Rand.Float64() < .45 {
			c.Core.QueueParticle("kuki", 1, attributes.Electro, 100) // TODO: idk the particle timing yet fml (or probability)
		}

		//queue up next hit only if next hit bell is still active
		if c.Core.F+90 <= c.bellActiveUntil {
			c.Core.Tasks.Add(c.bellTick(), 90)
		}
	}
}
