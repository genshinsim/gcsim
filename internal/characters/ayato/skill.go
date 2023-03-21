package ayato

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillFrames []int

func init() {
	skillFrames = frames.InitAbilSlice(21)
}

const SkillBuffKey = "soukaikanka"

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
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}

	ePos := c.Core.Combat.Player()
	c.a1OnSkill()
	c.Core.Tasks.Add(func() {
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(ePos, nil, 3.5), 0, 0)
		c.a1OnExplosion()
	}, delay)

	// start skill buff on cast
	c.AddStatus(SkillBuffKey, 6*60, true)
	// figure out atk buff
	if c.Base.Cons >= 6 {
		c.c6Ready = true
	}
	c.SetCD(action.ActionSkill, 12*60)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 1.9*60, true)

	count := 1.0
	if c.Core.Rand.Float64() < 0.5 {
		count = 2
	}
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Hydro, c.ParticleDelay) // TODO: this used to be 80 for particle delay
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
		c.DeleteStatus(SkillBuffKey)
		// queue up a4
		c.Core.Tasks.Add(c.a4, 60)
		return false
	}, "ayato-exit")
}
