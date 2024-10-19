package mualani

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var skillFrames []int
var skillCancelFrames []int

const (
	particleICDKey  = "mualani-particle-icd"
	momentumIcdKey  = "mualani-momentum-icd"
	surfingKey      = "mualani-surfing"
	markedAsPreyKey = "marked-as-prey"

	skillDelay      = 2
	particleICD     = 9999 * 60
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
	c.DeleteStatus(surfingKey)
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

	c.QueueCharTask(func() {
		c.nightsoulState.EnterBlessing(60)
		c.DeleteStatus(particleICDKey)
		c.a1Count = 0
		c.c1Done = false
		c.c2()
		c.nightsoulSrc = c.Core.F
		c.QueueCharTask(func() {
			c.AddStatus(surfingKey, -1, false)
			c.nightsoulPointReduceFunc(c.nightsoulSrc)()
		}, 6)
	}, skillDelay)

	canQueueAfter := skillFrames[action.ActionAttack] // earliest cancel
	// press skill "while" walking
	isWalking := c.Core.Player.AnimationHandler.CurrentState() == action.WalkState
	if isWalking {
		canQueueAfter = skillDelay
	}

	return action.Info{
		Frames: func(next action.Action) int {
			if next == action.ActionWalk && isWalking {
				// TODO: or 0f?
				return skillDelay
			}
			return skillFrames[next]
		},
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   canQueueAfter,
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

func (c *char) surfingTick() {
	// TODO: create a gadget?
	c.Core.Events.Subscribe(event.OnTick, func(args ...interface{}) bool {
		if c.Core.Player.Active() != c.Index {
			return false
		}
		if !c.nightsoulState.HasBlessing() {
			return false
		}
		if !c.StatusIsActive(surfingKey) {
			return false
		}

		switch c.Core.Player.CurrentState() {
		case action.DashState, action.JumpState, action.WalkState:
		default:
			return false
		}

		// to avoid spamming Surfing Hit logs
		useAttack := false
		ap := combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 0.9}, 0, 1)
		for _, e := range c.Core.Combat.Enemies() {
			enemy, ok := e.(combat.Enemy)
			if !ok {
				continue
			}

			willLand, _ := e.AttackWillLand(ap)
			if willLand && !enemy.StatusIsActive(momentumIcdKey) {
				useAttack = true
				break
			}
		}
		if !useAttack {
			return false
		}

		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               "Surfing Hit",
			AttackTag:          attacks.AttackTagNone,
			ICDTag:             attacks.ICDTagNone,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeSpear,
			Element:            attributes.Physical,
			Durability:         100,
			HitlagFactor:       0.01,
			CanBeDefenseHalted: true,
			IsDeployable:       true,
		}
		c.Core.QueueAttack(ai, ap, 0, 0, c.surfingCB)

		return false
	}, "mualani-surfing")
}

func (c *char) surfingCB(a combat.AttackCB) {
	enemy, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	if enemy.StatusIsActive(momentumIcdKey) {
		return
	}
	enemy.AddStatus(markedAsPreyKey, markedAsPreyDur, true)
	enemy.AddStatus(momentumIcdKey, momentumIcd, false)
	c.momentumStacks = min(c.momentumStacks+1, 3)
}
