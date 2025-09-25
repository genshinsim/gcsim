package skirk

import (
	"errors"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var (
	skillFrames     []int
	skillHoldFrames []int
)

const (
	maxSerpentsSubtlety    = 100
	skillGainSS            = 25
	skillKey               = "seven-phase-flash"
	skillDelay             = 19
	skillDur               = 754
	particleICD            = 15 * 60
	particleICDKey         = "skirk-particle-icd"
	skillHoldGainSS        = 18
	skillAbsorbRiftAnimKey = "skirk-hold-e-anim"
)

func init() {
	skillFrames = frames.InitAbilSlice(43) // E -> W
	skillFrames[action.ActionAttack] = 21
	skillFrames[action.ActionCharge] = 21
	skillFrames[action.ActionBurst] = 31
	skillFrames[action.ActionDash] = 29
	skillFrames[action.ActionJump] = 30
	skillFrames[action.ActionSwap] = 29

	skillHoldFrames = frames.InitAbilSlice(18)
	skillHoldFrames[action.ActionSwap] = 15
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillKey) {
		return action.Info{}, errors.New("skill cannot be used while in seven-phase flash")
	}

	h := p["hold"]
	if h > 0 {
		return c.skillHold(p)
	}
	return c.skillTap()
}

func (c *char) skillTap() (action.Info, error) {
	c.QueueCharTask(func() { c.enterSkillState() }, skillDelay)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionAttack],
		State:           action.SkillState,
	}, nil
}

func (c *char) enterSkillState() {
	c.skillSrc = c.Core.F
	c.AddStatus(skillKey, skillDur, false)
	c.AddSerpentsSubtlety(c.Base.Key.String()+"-skill", 45.0)
	c.c2OnSkill()
	c.serpentsReduceTask(c.skillSrc)
	src := c.skillSrc
	c.Core.Tasks.Add(func() { c.exitSkillState(src) }, skillDur)
}

func (c *char) exitSkillState(src int) {
	if c.skillSrc != src {
		return
	}
	c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index(), "exit skirk skill").Write("src", src)
	c.skillSrc = -1
	c.DeleteAttackMod(c2Key)
	c.DeleteStatus(skillKey)
	c.DeleteStatus(burstExtinctKey)
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
	// TODO: max duration of hold E?
	extraDuration := min(duration, 184) - 1 // subtract 1 because frames are listed as the minimum already
	c.QueueCharTask(func() {
		c.AddSerpentsSubtlety(c.Base.Key.String()+"-skill-hold", 45.0)
		c.c2OnSkill()
		c.absorbVoidRifts()
	}, skillHoldGainSS)

	// status used to absorb void rifts constantly during the hold E animation
	c.AddStatus(skillAbsorbRiftAnimKey, extraDuration, true)

	c.SetCDWithDelay(action.ActionSkill, 8*60, extraDuration+skillHoldGainSS)

	return action.Info{
		Frames: func(next action.Action) int {
			return skillHoldFrames[next] + extraDuration
		},
		AnimationLength: skillHoldFrames[action.InvalidAction] + extraDuration,
		CanQueueAfter:   skillHoldFrames[action.ActionSwap] + extraDuration,
		State:           action.SkillState,
	}, nil
}

func (c *char) particleInit() {
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...any) bool {
		atk := args[1].(*info.AttackEvent)
		_, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		if atk.Info.ActorIndex != c.Index() {
			return false
		}

		if atk.Info.Element != attributes.Cryo {
			return false
		}

		if c.StatusIsActive(particleICDKey) {
			return false
		}
		c.AddStatus(particleICDKey, particleICD, false)

		count := 4.0
		c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Cryo, c.ParticleDelay)

		return false
	}, c.Base.Key.String()+"-particles")
}
