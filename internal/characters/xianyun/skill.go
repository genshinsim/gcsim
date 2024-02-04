package xianyun

import (
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
	leapKey                  = "xianyun-leap"

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
	}

	if c.skillCounter == 3 {
		// Didn't plunge after the previous triple skill
		c.skillCounter = 0
		c.SetCD(action.ActionSkill, 10*60)
	}

	//C2: After using White Clouds at Dawn, Xianyun's ATK will be increased by 20% for 15s.
	if c.Base.Cons >= 2 {
		c.c2buff()
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Adeptal Aspect Trail",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 0,
		Mult:       skillPress[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 1),
		0,
		skillPressHitmark,
	)

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
