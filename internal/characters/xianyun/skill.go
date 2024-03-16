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
var skillStateDur = []int{220, 238, 179}

const (
	skillPressHitmark = 3

	skillStateKey = "cloud-transmogrification"

	// TODO: Find skill hitbox. Currently assuming the skill hitbox is the same size as the plunge collision hitbox
	skillRadius = 1.5

	particleCount  = 5
	particleICD    = 0.2 * 60
	particleICDKey = "xianyun-particle-icd"
)

func init() {
	skillLeapFrames = make([][]int, 3)
	skillLeapFrames[0] = frames.InitAbilSlice(244)

	skillLeapFrames[0][action.ActionAttack] = 221
	skillLeapFrames[0][action.ActionSkill] = 14
	skillLeapFrames[0][action.ActionBurst] = 40
	skillLeapFrames[0][action.ActionDash] = 39
	skillLeapFrames[0][action.ActionJump] = 41
	skillLeapFrames[0][action.ActionWalk] = 43
	skillLeapFrames[0][action.ActionSwap] = 36
	skillLeapFrames[0][action.ActionHighPlunge] = 14
	skillLeapFrames[0][action.ActionLowPlunge] = 14

	skillLeapFrames[1] = frames.InitAbilSlice(243)
	skillLeapFrames[1][action.ActionSkill] = 15
	skillLeapFrames[1][action.ActionBurst] = 60
	skillLeapFrames[1][action.ActionDash] = 60
	skillLeapFrames[1][action.ActionJump] = 60
	skillLeapFrames[1][action.ActionWalk] = 66
	skillLeapFrames[1][action.ActionSwap] = 59
	skillLeapFrames[1][action.ActionHighPlunge] = 15
	skillLeapFrames[1][action.ActionLowPlunge] = 15

	skillLeapFrames[2] = frames.InitAbilSlice(178)
	skillLeapFrames[2][action.ActionSkill] = 128
	skillLeapFrames[2][action.ActionBurst] = 126
	skillLeapFrames[2][action.ActionDash] = 130
	skillLeapFrames[2][action.ActionJump] = 129
	skillLeapFrames[2][action.ActionWalk] = 125
	skillLeapFrames[2][action.ActionSwap] = 126
	skillLeapFrames[2][action.ActionHighPlunge] = 18
	skillLeapFrames[2][action.ActionLowPlunge] = 18
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	// Check for first leap
	if !c.StatusIsActive(skillStateKey) || c.skillCounter == 3 { // Didn't plunge after the previous triple skill
		c.skillCounter = 0
		if c.StatusIsActive(c6Key) {
			c.skillWasC6 = true
			c.SetTag(c6Key, c.Tag(c6Key)-1)
			if c.Tag(c6Key) <= 0 {
				c.DeleteStatus(c6Key)
			}
		} else {
			c.SetCD(action.ActionSkill, 12*60)
			c.skillWasC6 = false
		}
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
		return !slices.Contains[[]targets.TargetKey](c.skillEnemiesHit, t.Key())
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
	c.QueueCharTask(c.cooldownReduce(c.Core.F), skillStateDur[c.skillCounter])
	c.AddStatus(skillStateKey, skillStateDur[c.skillCounter], true)

	defer func() { c.skillCounter++ }()

	return action.Info{
		Frames:          frames.NewAbilFunc(skillLeapFrames[c.skillCounter]),
		AnimationLength: skillLeapFrames[c.skillCounter][action.InvalidAction],
		CanQueueAfter:   skillLeapFrames[c.skillCounter][action.ActionHighPlunge], // earliest cancel
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

func (c *char) particleCB() func(combat.AttackCB) {
	// Particles are not produced if the skill was from c6
	if c.skillWasC6 {
		return nil
	}
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if c.StatusIsActive(particleICDKey) {
			return
		}
		c.AddStatus(particleICDKey, particleICD, true)

		c.Core.QueueParticle(c.Base.Key.String(), particleCount, attributes.Anemo, c.ParticleDelay)
	}
}
