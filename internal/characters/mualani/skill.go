package mualani

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillFrames []int
var skillCancelFrames []int

const (
	momentumDelay = 7

	particleICD     = 9999 * 60
	particleICDKey  = "mualani-particle-icd"
	momentumIcdKey  = "mualani-momentum-icd"
	markedAsPreyKey = "marked-as-prey"
	momentumIcd     = 0.7 * 60
	markedAsPreyDur = 10 * 60
)

func init() {
	skillFrames = frames.InitAbilSlice(69) // E -> E
	skillFrames[action.ActionAttack] = 5
	skillFrames[action.ActionBurst] = 6
	skillFrames[action.ActionDash] = 15
	skillFrames[action.ActionJump] = 18
	skillFrames[action.ActionWalk] = 19
	skillFrames[action.ActionSwap] = 65

	skillCancelFrames = frames.InitAbilSlice(17) // E -> Charge
	skillCancelFrames[action.ActionAttack] = 7
	skillCancelFrames[action.ActionBurst] = 4
	skillCancelFrames[action.ActionDash] = 13
	skillCancelFrames[action.ActionJump] = 12
	skillCancelFrames[action.ActionWalk] = 12
	skillCancelFrames[action.ActionSwap] = 11
}

func (c *char) reduceNightsoulPoints(val float64) {
	c.nightsoulState.ConsumePoints(val)
	if c.nightsoulState.Points() <= 0.00001 {
		c.cancelNightsoul()
	}
}

func (c *char) cancelNightsoul() {
	c.nightsoulState.ExitBlessing()
	c.SetCD(action.ActionSkill, 6*60)
	c.ResetActionCooldown(action.ActionAttack)
	c.momentumStacks = 0
	c.momentumSrc = -1
	c.nightsoulSrc = -1
}

func (c *char) nightsoulPointReduceFunc(src int) func() {
	return func() {
		if c.nightsoulSrc != src {
			return
		}

		if !c.nightsoulState.HasBlessing() {
			return
		}

		c.reduceNightsoulPoints(1)

		// reduce 1 point per 6f
		c.QueueCharTask(c.nightsoulPointReduceFunc(src), 6)
	}
}

func (c *char) momentumStackGain(src int) func() {
	return func() {
		if c.momentumSrc != src {
			return
		}

		if !c.nightsoulState.HasBlessing() {
			return
		}

		switch c.Core.Player.CurrentState() {
		case action.DashState, action.JumpState, action.WalkState:
		default:
			return
		}

		c.QueueCharTask(c.momentumStackGain(src), 0.1*60) // TODO: correct interval?

		ap := combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 0.9}, 0, 1)
		enemies := c.Core.Combat.Enemies()
		enemiesCollided := 0
		for _, e := range enemies {
			enemy, ok := e.(combat.Enemy)
			if !ok {
				continue
			}

			willLand, _ := e.AttackWillLand(ap)
			if willLand && !enemy.StatusIsActive(momentumIcdKey) {
				enemy.AddStatus(markedAsPreyKey, markedAsPreyDur, true)
				enemy.AddStatus(momentumIcdKey, 0.7*60, true)
				enemiesCollided++
			}
		}

		c.momentumStacks = min(c.momentumStacks+enemiesCollided, 3)
	}
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		c.cancelNightsoul()
		return action.Info{
			Frames:          frames.NewAbilFunc(skillCancelFrames),
			AnimationLength: skillCancelFrames[action.InvalidAction],
			CanQueueAfter:   skillCancelFrames[action.ActionAttack], // earliest cancel
			State:           action.SkillState,
		}, nil
	}

	c.nightsoulState.EnterBlessing(60)
	c.DeleteStatus(particleICDKey)
	c.a1Count = 0
	c.c1Done = false
	c.c2()
	c.nightsoulSrc = c.Core.F
	c.QueueCharTask(c.nightsoulPointReduceFunc(c.nightsoulSrc), 6)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionAttack], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, particleICD, true)

	count := 4.0
	if c.Core.Rand.Float64() < .5 {
		count = 5
	}
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Hydro, c.ParticleDelay)
}
