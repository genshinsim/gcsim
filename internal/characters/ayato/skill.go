package ayato

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var skillFrames []int

const skillStart = 21

func init() {
	skillFrames = frames.InitAbilSlice(21)
}

const skillBuffKey = "soukaikanka"

func (c *char) Skill(p map[string]int) action.ActionInfo {
	delay := p["illusion_delay"]
	if delay < 35 {
		delay = 35
	}
	if delay > 6*60 {
		delay = 360
	}

	ai := combat.AttackInfo{
		Abil:       "Kamisato Art: Kyouka",
		ActorIndex: c.Index,
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	c.Core.Tasks.Add(func() {
		c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 3.5), 0, 0)
		// A1:
		// set namisen stacks to max
		c.stacks = c.stacksMax
		c.Core.Log.NewEvent("ayato a1 set namisen stacks to max", glog.LogCharacterEvent, c.Index).
			Write("stacks", c.stacks)
	}, delay)

	//start skill buff on cast
	c.AddStatus(skillBuffKey, 6*60, true)
	//figure out atk buff
	if c.Base.Cons >= 6 {
		c.c6ready = true
	}
	c.SetCD(action.ActionSkill, 12*60)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) generateParticles(ac combat.AttackCB) {
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 114, true)
	var count float64 = 1
	if c.Core.Rand.Float64() < 0.5 {
		count++
	}
	//TODO: this used to be 80 for particle delay
	c.Core.QueueParticle("ayato", count, attributes.Hydro, c.ParticleDelay)
}

func (c *char) skillStacks(ac combat.AttackCB) {
	if c.stacks < c.stacksMax {
		c.stacks++
		c.Core.Log.NewEvent("gained namisen stack", glog.LogCharacterEvent, c.Index).
			Write("stacks", c.stacks)
	}
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		// do nothing if previous char wasn't ayato
		prev := args[0].(int)
		if prev != c.Index {
			return false
		}
		// clear skill status on field exit
		c.stacks = 0
		c.DeleteStatus(skillBuffKey)
		// queue up a4
		c.Core.Tasks.Add(c.a4, 60)
		return false
	}, "ayato-exit")
}
