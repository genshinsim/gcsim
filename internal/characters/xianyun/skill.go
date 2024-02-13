package xianyun

import (
	"slices"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillLeapFrames [][]int
var skillRecastFrames []int

const (
	skillPressHitmark        = 1
	skillFirstRecastHitmark  = 41
	skillSecondRecastHitmark = 18
	skillStateDur            = 2 * 60
	skillStateKey            = "cloud-transmogrification"

	// assuming the skill hitbox is the same size as the plunge collision hitbox
	skillRadius = 1.5
	leapKey     = "xianyun-leap"

	particleCount  = 5
	particleICD    = 0.2 * 60
	particleICDKey = "xianyun-particle-icd"
)

func init() {
	skillLeapFrames = make([][]int, 3)
	// skill -> x (can only use skill, plunge or wait(?))
	skillLeapFrames[0] = frames.InitAbilSlice(41)
	skillLeapFrames[0][action.ActionHighPlunge] = 28
	skillLeapFrames[0][action.ActionSkill] = skillFirstRecastHitmark

	// skill (recast) -> x (can only use skill, plunge or wait(?))
	skillLeapFrames[1] = frames.InitAbilSlice(46)
	skillLeapFrames[1][action.ActionHighPlunge] = 10
	skillLeapFrames[1][action.ActionSkill] = skillSecondRecastHitmark

	// skill (recast) -> x (can only use skill, plunge or wait(?))
	skillLeapFrames[2] = frames.InitAbilSlice(30)
	skillLeapFrames[2][action.ActionHighPlunge] = 42
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	// check if first leap

	if !c.StatusIsActive(skillStateKey) {
		c.skillCounter = 0
		c.SetCD(action.ActionSkill, 10*60)
		c.skillEnemiesHit = nil
	}

	if c.skillCounter == 3 {
		// Didn't plunge after the previous triple skill
		c.skillCounter = 0
		c.SetCD(action.ActionSkill, 10*60)
		c.skillEnemiesHit = nil
	}

	//C2: After using White Clouds at Dawn, Xianyun's ATK will be increased by 20% for 15s.
	c.c2buff()

	// This should only hit enemies once at most
	// During each Cloud Transmogrification state Xianyun enters, Skyladder may be used up to 3 times and only 1 instance of Skyladder DMG can be dealt to any one opponent.
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Skyladder",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 0,
		Mult:       skillPress[c.TalentLvlSkill()],
	}

	aoe := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, skillRadius)
	targets := c.Core.Combat.EnemiesWithinArea(aoe, func(t combat.Enemy) bool {
		return slices.Contains[[]targets.TargetKey](c.skillEnemiesHit, t.Key())
	})

	for _, t := range targets {
		c.Core.QueueAttack(
			ai,
			combat.NewSingleTargetHit(t.Key()),
			skillPressHitmark,
			skillPressHitmark,
		)
		c.skillEnemiesHit = append(c.skillEnemiesHit, t.Key())
	}

	c.skillSrc = c.Core.F
	c.QueueCharTask(c.cooldownReduce(c.Core.F), skillStateDur)
	c.AddStatus(skillStateKey, skillStateDur, true)

	idx := c.skillCounter
	c.skillCounter++

	return action.Info{
		Frames:          frames.NewAbilFunc(skillLeapFrames[idx]),
		AnimationLength: skillLeapFrames[idx][action.InvalidAction],
		CanQueueAfter:   skillLeapFrames[idx][action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) cooldownReduce(src int) func() {
	return func() {
		if c.skillSrc != src {
			return
		}
		// If Xianyun does not use Driftcloud Wave while in this state, the next CD of White Clouds at Dawn will be decreased by 3s.
		c.ReduceActionCooldown(action.ActionSkill, 3*60)
	}
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, particleICD, true)

	c.Core.QueueParticle(c.Base.Key.String(), particleCount, attributes.Anemo, c.ParticleDelay)
}
