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
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	c.Core.Tasks.Add(func() {
		c.Core.QueueAttack(ai, combat.NewDefCircHit(3.5, false, combat.TargettableEnemy), 0, 0)
		//add a namisen stack
		if c.stacks < c.stacksMax {
			c.stacks++
		}
	}, delay)

	//start skill buff after animation
	//TODO: make sure this isn't causing bugs?
	c.QueueCharTask(func() {
		c.AddStatus(skillBuffKey, 6*60, true)
	}, skillStart)
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
	c.Core.QueueParticle("ayato", count, attributes.Hydro, 80)
}

func (c *char) skillStacks(ac combat.AttackCB) {
	if c.stacks < c.stacksMax {
		c.stacks++
		c.Core.Log.NewEvent("gained namisen stack", glog.LogCharacterEvent, c.Index).
			Write("stacks", c.stacks)
	}
}

// clear skill status on field exit
func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		c.stacks = 0
		c.DeleteStatus(skillBuffKey)
		c.a4()
		return false
	}, "ayato-exit")
}
