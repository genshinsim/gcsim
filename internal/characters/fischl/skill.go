package fischl

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var skillFrames []int

const skillOzSpawn = 32

func init() {
	skillFrames = frames.InitAbilSlice(43)
	skillFrames[action.ActionDash] = 14
	skillFrames[action.ActionJump] = 16
	skillFrames[action.ActionSwap] = 42
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	//always trigger electro no ICD on initial summon
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Oz (Summon)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupFischl,
		StrikeType: combat.StrikeTypePierce,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       birdSum[c.TalentLvlSkill()],
	}

	if c.Base.Cons >= 2 {
		ai.Mult += 2
	}
	//hitmark is 5 frames after oz spawns
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 1, false, combat.TargettableEnemy), skillOzSpawn, skillOzSpawn+5)

	c.Core.Tasks.Add(func() {
		// activate status at CanQueueAfter because adding the status later than that breaks things
		// if you do E -> D then everything's fine, but any other cancel breaks the conditional
		// example:
		// "fischl skill" followed by a ".status.fischloz == 0" check
		dur := 600 + skillOzSpawn - skillFrames[action.ActionDash]
		if c.Base.Cons == 6 {
			dur += 120
		}
		c.Core.Status.Add("fischloz", dur)
	}, skillFrames[action.ActionDash])

	//set on field oz to be this one
	c.Core.Tasks.Add(func() {
		c.queueOz("Skill")
	}, skillOzSpawn)

	// CD Delay is 18 frames, but things break if Delay > CanQueueAfter
	// so we add 18 to the duration instead. this probably mess up CDR stuff
	c.SetCD(action.ActionSkill, 25*60+18) //18 frames until CD starts

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) queueOz(src string) {
	dur := 600
	if c.Base.Cons == 6 {
		dur += 120
	}
	c.ozActiveUntil = c.Core.F + dur
	c.ozSource = c.Core.F

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Oz (%v)", src),
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupFischl,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       birdAtk[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)
	c.ozSnapshot = combat.AttackEvent{
		Info:        ai,
		Snapshot:    snap,
		Pattern:     combat.NewDefSingleTarget(c.Core.Combat.DefaultTarget, combat.TargettableEnemy),
		SourceFrame: c.Core.F,
	}
	c.Core.Tasks.Add(c.ozTick(c.Core.F), 60)
	c.Core.Log.NewEvent("Oz activated", glog.LogCharacterEvent, c.Index).
		Write("source", src).
		Write("expected end", c.ozActiveUntil).
		Write("next expected tick", c.Core.F+60)
}

func (c *char) ozTick(src int) func() {
	return func() {
		c.Core.Log.NewEvent("Oz checking for tick", glog.LogCharacterEvent, c.Index).
			Write("src", src)
		//if src != ozSource then this is no longer the same oz, do nothing
		if src != c.ozSource {
			return
		}
		c.Core.Log.NewEvent("Oz ticked", glog.LogCharacterEvent, c.Index).
			Write("next expected tick", c.Core.F+60).
			Write("active", c.ozActiveUntil).
			Write("src", src)
		//trigger damage
		ae := c.ozSnapshot
		c.Core.QueueAttackEvent(&ae, 0)
		//check for orb
		//Particle check is 67% for particle, from datamine
		//TODO: this delay used to be 120
		if c.Core.Rand.Float64() < .67 {
			c.Core.QueueParticle("fischl", 1, attributes.Electro, c.ParticleDelay)
		}

		//queue up next hit only if next hit oz is still active
		if c.Core.F+60 <= c.ozActiveUntil {
			c.Core.Tasks.Add(c.ozTick(src), 60)
		}
	}
}
