package chasca

import (
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillFrames []int
var skillCancelFrames []int

const (
	skillHitmarks      = 5
	plungeAvailableKey = "chasca-plunge-available"
)

func init() {
	skillFrames = frames.InitAbilSlice(26)
	skillFrames[action.ActionAttack] = 5
	skillFrames[action.ActionAim] = 16
	skillFrames[action.ActionSkill] = 26
	skillFrames[action.ActionBurst] = 5 // TODO: not in frames sheet
	skillFrames[action.ActionDash] = 5
	skillFrames[action.ActionJump] = 22
	skillFrames[action.ActionSwap] = 586 + 37 // wait for nightsoul to run out and fall onto the ground

	skillCancelFrames = frames.InitAbilSlice(40)
	skillCancelFrames[action.ActionAttack] = 38
	skillCancelFrames[action.ActionAim] = 38
	skillCancelFrames[action.ActionLowPlunge] = 2
	skillCancelFrames[action.ActionBurst] = 39
	skillCancelFrames[action.ActionWalk] = 37
	skillCancelFrames[action.ActionSwap] = 37
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

// If NS expired gives the skillCancelFrames, otherwise gives the next frames as input
func (c *char) skillNextFrames(f func(next action.Action) int, extraDelay int) func(next action.Action) int {
	// this is used to calculate the hitlag effect time elapsed since action start
	actionStart := c.TimePassed
	actionEnd := -1
	return func(next action.Action) int {
		if c.nightsoulState.HasBlessing() {
			return f(next)
		}
		if actionEnd < 0 {
			actionEnd = c.TimePassed
		}
		// TODO: set fall down animation to be "falling/idle" when this occurs?
		return actionEnd - actionStart + skillCancelFrames[next] + extraDelay
	}
}

func (c *char) enterNightsoul() {
	c.Core.Player.SwapCD = math.MaxInt16 // block swapping while in the air
	c.nightsoulState.EnterBlessing(80)
	c.nightsoulSrc = c.Core.F
	c.Core.Tasks.Add(c.nightsoulPointReduceFunc(c.nightsoulSrc), 6)
	c.NormalHitNum = 1
	c.NormalCounter = 0
	c.skillParticleICD = false
}

func (c *char) exitNightsoul() {
	if !c.nightsoulState.HasBlessing() {
		return
	}
	if c.Core.Player.CurrentState() == action.AimState {
		// keep charging bullets for up to 10f after NS ends
		c.QueueCharTask(c.fireBullets, 10)
	}

	c.nightsoulState.ExitBlessing()
	c.nightsoulState.ClearPoints()
	c.nightsoulSrc = -1
	c.SetCD(action.ActionSkill, 6.5*60)
	c.NormalHitNum = normalHitNum
	c.NormalCounter = 0
	c.Core.Player.SwapCD = 37
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
		c.exitNightsoul()
		return action.Info{
			Frames:          frames.NewAbilFunc(skillCancelFrames),
			AnimationLength: skillFrames[action.InvalidAction],
			CanQueueAfter:   skillFrames[action.ActionLowPlunge], // earliest cancel
			State:           action.SkillState,
		}, nil
	}

	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Spirit Reins, Shadow Hunt",
		AttackTag:      attacks.AttackTagElementalArt,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypePierce,
		Element:        attributes.Anemo,
		Durability:     25,
		Mult:           skillResonance[c.TalentLvlSkill()],
	}
	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		nil,
		5.5,
	)
	c.Core.QueueAttack(ai, ap, skillHitmarks, skillHitmarks)
	c.enterNightsoul()

	return action.Info{
		Frames:          c.skillNextFrames(frames.NewAbilFunc(skillFrames), 0),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.skillParticleICD {
		return
	}
	c.skillParticleICD = true
	c.Core.QueueParticle(c.Base.Key.String(), 5, attributes.Anemo, c.ParticleDelay)
}
