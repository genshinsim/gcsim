package ifa

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	skillFrames       []int
	skillCancelFrames []int
)

const (
	// skillHitmarks      = 3
	plungeAvailableKey = "ifa-plunge-available"
)

func init() {
	skillFrames = frames.InitAbilSlice(31) // E -> N
	skillFrames[action.ActionCharge] = 28
	skillFrames[action.ActionSkill] = 17
	skillFrames[action.ActionBurst] = 5
	skillFrames[action.ActionDash] = 28
	skillFrames[action.ActionSwap] = 589 + 42 // wait for nightsoul to run out and fall onto the ground

	skillCancelFrames = frames.InitAbilSlice(69) // E -> Walk
	skillCancelFrames[action.ActionAttack] = 50
	skillCancelFrames[action.ActionCharge] = 49
	skillCancelFrames[action.ActionLowPlunge] = 6
	skillCancelFrames[action.ActionBurst] = 49
	skillCancelFrames[action.ActionDash] = 44
	skillCancelFrames[action.ActionJump] = 58
	skillCancelFrames[action.ActionSwap] = 55
}

func (c *char) reduceNightsoulPoints(val float64) {
	c.nightsoulState.ConsumePoints(val)
	c.checkNS()
}

// Checks the current number of nightsoul points and exits nightsoul if there aren't enough. Returns the status of NS after the check
func (c *char) checkNS() {
	if c.nightsoulState.Points() < 0.001 {
		c.exitNightsoul()
	}
}

func (c *char) enterNightsoul() {
	c.nightsoulState.EnterBlessing(80)
	c.nightsoulSrc = c.Core.F
	c.Core.Tasks.Add(c.nightsoulPointReduceFunc(c.nightsoulSrc), 4)
}

func (c *char) nigthsoulFallingMsg() {
	c.Core.Log.NewEvent("nightsoul ended, falling", glog.LogCharacterEvent, c.Index())
}

func (c *char) exitNightsoul() {
	if !c.nightsoulState.HasBlessing() {
		return
	}
	c.Core.Player.SwapCD = 37
	c.nigthsoulFallingMsg()

	c.nightsoulState.ExitBlessing()
	c.nightsoulState.ClearPoints()
	c.nightsoulSrc = -1
	c.SetCD(action.ActionSkill, 7.5*60)
	c.NormalHitNum = normalHitNum
	c.NormalCounter = 0
	c.AddStatus(plungeAvailableKey, 26, true)
}

func (c *char) nightsoulPointReduceFunc(src int) func() {
	return func() {
		if c.nightsoulSrc != src {
			return
		}
		c.reduceNightsoulPoints(0.8)
		// reduce 0.8 point per 6, which is 8 per second
		c.Core.Tasks.Add(c.nightsoulPointReduceFunc(src), 6)
	}
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		if p["hold"] == 0 {
			c.exitNightsoul()
			return action.Info{
				Frames:          frames.NewAbilFunc(skillCancelFrames),
				AnimationLength: skillCancelFrames[action.InvalidAction],
				CanQueueAfter:   skillCancelFrames[action.ActionLowPlunge], // earliest cancel
				State:           action.SkillState,
			}, nil
		}

		c.LowPlungeAttack(p)
	}

	c.enterNightsoul()
	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionAttack], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) healCB(a info.AttackCB) {
	em := c.Stat(attributes.EM)
	healAmt := skill_heal[c.TalentLvlSkill()]*em + skill_heal_flat[c.TalentLvlSkill()]
	healBonus := c.Stat(attributes.Heal)

	hi := info.HealInfo{
		Caller:  c.Index(),
		Target:  -1,
		Message: "Tonicshot Healing",
		Src:     healAmt,
		Bonus:   healBonus,
	}

	c.Core.Player.Heal(hi)
}

func (c *char) particleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.skillParticleICD {
		return
	}
	c.skillParticleICD = true

	c.Core.QueueParticle(c.Base.Key.String(), 4, attributes.Anemo, c.ParticleDelay)

	if c.Core.Rand.Float64() < 0.3 {
		c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Anemo, c.ParticleDelay)
	}
}
