package ifa

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	skillFrames             []int
	skillCancelFrames       []int
	skillCancelPlungeFrames []int
)

const (
	// skillHitmarks      = 3
	plungeAvailableKey       = "ifa-plunge-available"
	skillCancelPlungeHitmark = 39
	healTapICDKey            = "tonicshot-tap-healing-icd"
	healHoldICDKey           = "tonicshot-hold-healing-icd"
)

func init() {
	skillFrames = frames.InitAbilSlice(31) // E -> N
	skillFrames[action.ActionCharge] = 28
	skillFrames[action.ActionSkill] = 17
	skillFrames[action.ActionBurst] = 5
	skillFrames[action.ActionDash] = 28
	skillFrames[action.ActionSwap] = 589 + 42 // wait for nightsoul to run out and fall onto the ground

	skillCancelFrames = frames.InitAbilSlice(44) // E -> Jump
	skillCancelFrames[action.ActionAttack] = 43
	skillCancelFrames[action.ActionBurst] = 42
	skillCancelFrames[action.ActionJump] = 43
	skillCancelFrames[action.ActionWalk] = 41
	skillCancelFrames[action.ActionSwap] = 42
	skillCancelFrames[action.ActionLowPlunge] = 6

	skillCancelPlungeFrames = frames.InitAbilSlice(69) // E -> Walk
	skillCancelPlungeFrames[action.ActionAttack] = 50
	skillCancelPlungeFrames[action.ActionCharge] = 49
	skillCancelPlungeFrames[action.ActionSkill] = skillCancelPlungeHitmark
	skillCancelPlungeFrames[action.ActionBurst] = 49
	skillCancelPlungeFrames[action.ActionDash] = 63 - 19
	skillCancelPlungeFrames[action.ActionJump] = 58
	skillCancelPlungeFrames[action.ActionSwap] = 55
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
	c.skillParticleICD = false
	c.nightsoulState.EnterBlessing(80)
	c.nightsoulSrc = c.Core.F
	c.Core.Tasks.Add(c.nightsoulPointReduceFunc(c.nightsoulSrc), 4)
	c.skillLastStamF = c.Core.Player.LastStamUse
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
	c.Core.Player.LastStamUse = c.skillLastStamF
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
		if p["hold"] == 1 {
			return c.skillPlunge(p)
		}
		c.exitNightsoul()

		return action.Info{
			Frames:          frames.NewAbilFunc(skillCancelFrames),
			AnimationLength: skillCancelFrames[action.InvalidAction],
			CanQueueAfter:   skillCancelFrames[action.ActionLowPlunge], // earliest cancel
			State:           action.SkillState,
		}, nil
	}

	c.enterNightsoul()
	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionBurst], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) skillPlunge(p map[string]int) (action.Info, error) {
	c.DeleteStatus(plungeAvailableKey)

	collision, ok := p["collision"]
	if !ok {
		collision = 0
	}

	if collision > 0 {
		c.plungeCollision(collisionHitmark)
	}

	ai := info.AttackInfo{
		ActorIndex:     c.Index(),
		Abil:           "Low Plunge Attack",
		AttackTag:      attacks.AttackTagPlunge,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Anemo,
		Durability:     25,
		Mult:           plunge_low[c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3),
		skillCancelPlungeHitmark, skillCancelPlungeHitmark)

	c.Core.Tasks.Add(func() {
		c.exitNightsoul()
		c.DeleteStatus(plungeAvailableKey)
		c.Core.Player.SwapCD = 0
	}, skillCancelPlungeHitmark)

	if c.nightsoulState.HasBlessing() {
		ai.AdditionalTags = []attacks.AdditionalTag{attacks.AdditionalTagNightsoul}
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(skillCancelPlungeFrames),
		AnimationLength: skillCancelPlungeFrames[action.InvalidAction],
		CanQueueAfter:   skillCancelPlungeFrames[action.ActionDash], // earliest cancel
		State:           action.PlungeAttackState,
	}, nil
}

func (c *char) healTapCB(a info.AttackCB) {
	if c.StatusIsActive(healTapICDKey) {
		return
	}

	c.AddStatus(healTapICDKey, 0.1*60, true)
	c.healSkill()
}

func (c *char) healHoldCB(a info.AttackCB) {
	if c.StatusIsActive(healHoldICDKey) {
		return
	}

	c.AddStatus(healHoldICDKey, 0.1*60, true)
	c.healSkill()
}

func (c *char) healSkill() {
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

	particles := 4.0

	if c.Core.Rand.Float64() < 0.3 {
		particles = 5.0
	}

	c.Core.QueueParticle(c.Base.Key.String(), particles, attributes.Anemo, c.ParticleDelay)
}
