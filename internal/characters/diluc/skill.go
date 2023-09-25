package diluc

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var (
	skillFrames       [][]int
	skillHitmarks     = []int{24, 28, 46}
	skillHitlagStages = []float64{.12, .12, .16}
	skillHitboxes     = [][]float64{{3, 3.5}, {2.2}, {3.5, 4}}
	skillOffsets      = []float64{0, 1.2, -0.3}
	skillFanAngles    = []float64{360, 300, 360}
)

const particleICDKey = "diluc-particle-icd"

func init() {
	skillFrames = make([][]int, 3)

	// skill (1st) -> x
	skillFrames[0] = frames.InitAbilSlice(32)
	skillFrames[0][action.ActionSkill] = 31
	skillFrames[0][action.ActionDash] = 24
	skillFrames[0][action.ActionJump] = 24
	skillFrames[0][action.ActionSwap] = 30

	// skill (2nd) -> x
	skillFrames[1] = frames.InitAbilSlice(38)
	skillFrames[1][action.ActionSkill] = 37
	skillFrames[1][action.ActionBurst] = 37
	skillFrames[1][action.ActionDash] = 28
	skillFrames[1][action.ActionJump] = 31
	skillFrames[1][action.ActionSwap] = 36

	// skill (3rd) -> x
	// TODO: missing counts for skill -> skill
	skillFrames[2] = frames.InitAbilSlice(66)
	skillFrames[2][action.ActionAttack] = 58
	skillFrames[2][action.ActionSkill] = 57 // uses burst frames
	skillFrames[2][action.ActionBurst] = 57
	skillFrames[2][action.ActionDash] = 47
	skillFrames[2][action.ActionJump] = 48
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	// reset counter
	if !c.StatusIsActive(eWindowKey) {
		c.eCounter = 0
	}

	// C6: After casting Searing Onslaught, the next 2 Normal Attacks within the
	// next 6s will have their DMG and ATK SPD increased by 30%.
	if c.Base.Cons >= 6 {
		c.c6Count = 0
		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = 0.3
		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag("diluc-c6-dmg", 360),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if atk.Info.AttackTag != attacks.AttackTagNormal {
					return nil, false
				}
				if c.c6Count > 1 {
					return nil, false
				}
				c.c6Count++
				return m, true
			},
		})

		mAtkSpd := make([]float64, attributes.EndStatType)
		mAtkSpd[attributes.AtkSpd] = 0.3
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("diluc-c6-speed", 360),
			AffectedStat: attributes.AtkSpd,
			Amount: func() ([]float64, bool) {
				if c.Core.Player.CurrentState() != action.NormalAttackState {
					return nil, false
				}
				if c.c6Count > 1 {
					return nil, false
				}
				return mAtkSpd, true
			},
		})
	}

	hitmark := skillHitmarks[c.eCounter]

	// actual skill cd starts immediately on first cast
	// times out after 4 seconds of not using
	// every hit applies pyro
	// apply attack speed
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Searing Onslaught %v", c.eCounter),
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagNone,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		Element:            attributes.Pyro,
		Durability:         25,
		Mult:               skill[c.eCounter][c.TalentLvlSkill()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   skillHitlagStages[c.eCounter] * 60,
		CanBeDefenseHalted: true,
	}
	ap := combat.NewCircleHitOnTargetFanAngle(
		c.Core.Combat.Player(),
		geometry.Point{Y: skillOffsets[c.eCounter]},
		skillHitboxes[c.eCounter][0],
		skillFanAngles[c.eCounter],
	)
	if c.eCounter == 0 || c.eCounter == 2 {
		ap = combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: skillOffsets[c.eCounter]},
			skillHitboxes[c.eCounter][0],
			skillHitboxes[c.eCounter][1],
		)
	}
	c.Core.QueueAttack(ai, ap, hitmark, hitmark, c.particleCB)

	// add a timer to activate c4
	if c.Base.Cons >= 4 {
		c.Core.Tasks.Add(func() {
			c.c4()
		}, hitmark+120) // 2seconds after cast
	}

	// allow skill to be used again if 4s hasn't passed since last use
	c.AddStatus(eWindowKey, 4*60, true)

	// store skill counter so we can determine which frames to return
	idx := c.eCounter
	c.eCounter++
	switch c.eCounter {
	case 1:
		// TODO: cd delay?
		// set cd on first use
		c.SetCD(action.ActionSkill, 10*60)
	case 3:
		// reset window since we're at 3rd use
		c.DeleteStatus(eWindowKey)
		c.eCounter = 0
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames[idx]),
		AnimationLength: skillFrames[idx][action.InvalidAction],
		CanQueueAfter:   skillFrames[idx][action.ActionDash], // earliest cancel
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
	c.AddStatus(particleICDKey, 0.3*60, true)

	count := 1.0
	if c.Core.Rand.Float64() < 0.33 {
		count = 2
	}
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Pyro, c.ParticleDelay)
}
