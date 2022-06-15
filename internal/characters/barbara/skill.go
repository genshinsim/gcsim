package barbara

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

// barbara skill - copied from bennett burst
const skillDuration = 15*60 + 1
const barbSkillKey = "barbskill"

var skillFrames []int

func init() {
	skillFrames = frames.InitAbilSlice(52)
}

func (c *char) Skill(p map[string]int) action.ActionInfo {

	c.Core.Status.Add(barbSkillKey, skillDuration)

	//activate a1
	c.a1()

	//restart a4 counter
	c.a4extendCount = 0

	//hook for buffs; active right away after cast

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Let the Show Begin♪",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Hydro,
		Durability: 25, //TODO: what is 1A GU?
		Mult:       skill[c.TalentLvlSkill()],
	}
	//TODO: review barbara AOE size?
	c.Core.QueueAttack(ai, combat.NewDefCircHit(1, false, combat.TargettableEnemy), 5, 5)
	c.Core.QueueAttack(ai, combat.NewDefCircHit(1, false, combat.TargettableEnemy), 5, 35) // need to confirm timing of this

	stats, _ := c.Stats()
	hpplus := stats[attributes.Heal]
	heal := skillhp[c.TalentLvlSkill()] + skillhpp[c.TalentLvlSkill()]*c.MaxHP()
	//apply right away

	c.skillInitF = c.Core.F
	//add 1 tick each 5s
	//first tick starts at 0
	c.barbaraHealTick(heal, hpplus, c.Core.F)()
	ai.Abil = "Let the Show Begin♪ Wet Tick"
	ai.AttackTag = combat.AttackTagNone
	ai.Mult = 0
	c.barbaraWet(ai, c.Core.F)()

	if c.Base.Cons >= 2 {
		c.c2() //c2 hydro buff
		c.SetCD(action.ActionSkill, 32*60*0.85)
	} else {
		c.SetCD(action.ActionSkill, 32*60)
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.InvalidAction],
		State:           action.SkillState,
	}
}

func (c *char) barbaraHealTick(healAmt float64, hpplus float64, skillInitF int) func() {
	return func() {
		//make sure it's not overwritten
		if c.skillInitF != skillInitF {
			return
		}
		//do nothing if buff expired
		if c.Core.Status.Duration(barbSkillKey) == 0 {
			return
		}
		// c.Core.Log.NewEvent("barbara heal ticking", core.LogCharacterEvent, c.Index)
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  c.Core.Player.Active(),
			Message: "Melody Loop (Tick)",
			Src:     healAmt,
			Bonus:   hpplus,
		})

		// tick per 5 seconds
		c.Core.Tasks.Add(c.barbaraHealTick(healAmt, hpplus, skillInitF), 5*60)
	}
}

func (c *char) barbaraWet(ai combat.AttackInfo, skillInitF int) func() {
	return func() {
		//make sure it's not overwritten
		if c.skillInitF != skillInitF {
			return
		}
		//do nothing if buff expired
		if c.Core.Status.Duration(barbSkillKey) == 0 {
			return
		}
		c.Core.Log.NewEvent("barbara wet ticking", glog.LogCharacterEvent, c.Index)

		c.Core.QueueAttack(ai, combat.NewDefCircHit(1, false, combat.TargettableEnemy), -1, 5)

		// tick per 5 seconds
		c.Core.Tasks.Add(c.barbaraWet(ai, skillInitF), 5*60)
	}
}

func (c *char) a4() {
	//When your active character gains an Elemental Orb/Particle, the duration
	//of the Melody Loop of Let the Show Begin♪ is extended by 1s. The maximum
	//extension is 5s.
	c.Core.Events.Subscribe(event.OnParticleReceived, func(args ...interface{}) bool {
		//TODO: assuming this works no matter who's on field since it just says
		//active char?
		if c.Core.Status.Duration(barbSkillKey) == 0 {
			return false
		}
		if c.a4extendCount == 5 {
			return false
		}

		c.a4extendCount++
		c.Core.Status.Extend(barbSkillKey, 60)

		c.Core.Log.NewEvent("barbara skill extended from a4", glog.LogCharacterEvent, c.Index)

		return false
	}, "barbara-a4")
}
