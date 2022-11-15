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

const (
	skillHitmark     = 11 // Initial Hit
	hpDrainThreshold = 0.2
)

func init() {
	skillFrames = frames.InitAbilSlice(52) // E -> Q
	skillFrames[action.ActionAttack] = 50  // E -> N1
	skillFrames[action.ActionDash] = 12    // E -> D
	skillFrames[action.ActionJump] = 11    // E -> J
	skillFrames[action.ActionSwap] = 50    // E -> Swap
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	// only drain HP when above 20% HP
	if c.HPCurrent/c.MaxHP() > hpDrainThreshold {
		hpdrain := 0.3 * c.HPCurrent
		// The HP consumption from using this skill can only bring her to 20% HP.
		if (c.HPCurrent-hpdrain)/c.MaxHP() <= hpDrainThreshold {
			hpdrain = c.HPCurrent - hpDrainThreshold*c.MaxHP()
		}
		c.Core.Player.Drain(player.DrainInfo{
			ActorIndex: c.Index,
			Abil:       "Sanctifying Ring",
			Amount:     hpdrain,
		})
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Sanctifying Ring",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
		FlatDmg:    c.Stat(attributes.EM) * 0.25,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 4), skillHitmark, skillHitmark)

	// C2: Grass Ring of Sanctification's duration is increased by 3s.
	skilldur := 720
	if c.Base.Cons >= 2 {
		skilldur = 900 //12+3s
	}

	// this gets executed before kuki can expierence hitlag so no need for char queue
	// ring duration starts after hitmark
	c.Core.Tasks.Add(func() {
		// E duration and ticks are not affected by hitlag
		c.Core.Status.Add("kuki-e", skilldur)
		c.Core.Tasks.Add(c.bellTick(), 90) // Assuming this executes every 90 frames = 1.5s
		c.bellActiveUntil = c.Core.F + skilldur
		c.Core.Log.NewEvent("Bell activated", glog.LogCharacterEvent, c.Index).
			Write("expected end", c.bellActiveUntil).
			Write("next expected tick", c.Core.F+90)
	}, 23)

	c.SetCDWithDelay(action.ActionSkill, 15*60, 7)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionJump], // earliest cancel
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
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       skilldot[c.TalentLvlSkill()],
			FlatDmg:    c.Stat(attributes.EM) * 0.25,
		}
		c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 4), 2, 2)

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
			c.Core.QueueParticle("kuki", 1, attributes.Electro, c.ParticleDelay) // TODO: idk the particle timing yet fml (or probability)
		}

		//queue up next hit only if next hit bell is still active
		if c.Core.F+90 <= c.bellActiveUntil {
			c.Core.Tasks.Add(c.bellTick(), 90)
		}
	}
}
