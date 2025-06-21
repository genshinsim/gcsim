package skirk

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var skillFrames []int
var skillHoldFrames []int

const (
	maxSerpentsSubtlety = 100
	skillGainSS         = 25
	skillKey            = "seven-phase-slash"
	particleICDKey      = "skirk-particle-icd"
	skillHoldGainSS     = 16
)

func init() {
	skillFrames = frames.InitAbilSlice(34)
	skillHoldFrames = frames.InitAbilSlice(16)
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	h := p["hold"]
	if h > 0 {
		return c.skillHold(p)
	}
	return c.skillTap()
}

func (c *char) skillTap() (action.Info, error) {
	if c.StatusIsActive(skillKey) {
		c.exitSkillState(c.skillSrc)
	} else {
		c.QueueCharTask(func() { c.enterSkillState() }, skillGainSS)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel is before skillHitmark
		State:           action.SkillState,
	}, nil
}

func (c *char) enterSkillState() {
	c.skillSrc = c.Core.F
	c.AddStatus(skillKey, 12.5*60, false)
	c.AddSerpentsSubtlety(c.Base.Key.String()+"-skill", 45.0)
	c.c2OnSkill()
	c.serpentsReduceTask(c.skillSrc)
	src := c.skillSrc
	c.Core.Tasks.Add(func() { c.exitSkillState(src) }, 12.5*60)
}

func (c *char) exitSkillState(src int) {
	if c.skillSrc != src {
		return
	}
	c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index, "exit skirk skill").Write("src", src)
	c.skillSrc = -1
	c.DeleteAttackMod(c2Key)
	c.DeleteStatus(skillKey)
	c.SetCD(action.ActionSkill, 8*60)
	c.ConsumeSerpentsSubtlety(0, c.Base.Key.String()+"-skill-exit")
}

func (c *char) serpentsReduceTask(src int) {
	const tickInterval = .2
	c.Core.Tasks.Add(func() {
		if c.skillSrc != src {
			return
		}
		// reduce 1.4 point every 12f, which is 7 per second
		c.ReduceSerpentsSubtlety(c.Base.Key.String()+"skill", 1.4)
		if c.serpentsSubtlety == 0 && c.StatusIsActive(skillKey) {
			c.exitSkillState(src)
		}
		c.serpentsReduceTask(src)
	}, 60*tickInterval)
}

func (c *char) skillHold(p map[string]int) (action.Info, error) {
	duration := p["hold"]
	if duration < 10 {
		duration = 10
	}

	c.QueueCharTask(func() {
		c.AddSerpentsSubtlety(c.Base.Key.String()+"-skill-hold", 45.0)
		c.c2OnSkill()

		c.absorbVoidRift()
	}, skillHoldGainSS)

	c.SetCDWithDelay(action.ActionSkill, 8*60, duration)

	return action.Info{
		Frames: func(next action.Action) int {
			return skillHoldFrames[next] + duration
		},
		AnimationLength: skillHoldFrames[action.InvalidAction] + duration,
		CanQueueAfter:   skillHoldFrames[action.ActionDash] + duration, // earliest cancel is before skillHitmark
		State:           action.SkillState,
	}, nil
}

func (c *char) particleInit() {
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		_, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		if atk.Info.ActorIndex != c.Index {
			return false
		}

		if atk.Info.Element != attributes.Cryo {
			return false
		}

		if c.StatusIsActive(particleICDKey) {
			return false
		}
		c.AddStatus(particleICDKey, 15*60, false)

		count := 4.0
		c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Cryo, c.ParticleDelay)

		return false
	}, c.Base.Key.String()+"-particles")
}
