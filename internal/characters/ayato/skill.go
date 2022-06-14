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

func (c *char) Skill(p map[string]int) action.ActionInfo {
	delay := p["illusion_delay"]
	if delay < 35 {
		delay = 35
	}
	if delay > 6*60 {
		delay = 360
	}
	hitlag := p["hitlag_extend"]

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

	c.Core.Status.Add("soukaikanka", 6*60+skillStart+hitlag) //add animation to the duration
	c.Core.Log.NewEvent("Soukai Kanka acivated", glog.LogCharacterEvent, c.Index, "expiry", c.Core.F+6*60+skillStart+hitlag)
	//figure out atk buff
	if c.Base.Cons >= 6 {
		c.c6ready = true

	}
	c.SetCD(action.ActionSkill, 12*60)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		Post:            skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

// clear skill status on field exit
func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		c.stacks = 0
		c.Core.Status.Delete("soukaikanka")
		c.a4()
		return false
	}, "ayato-exit")
}
