package bennett

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	skillFrames       [][]int
	skillHoldHitmarks = [][]int{{45, 57}, {112, 121}}
	skillHoldHitboxes = [][]float64{{2.5}, {3, 3}}
	skillHoldOffsets  = []float64{0.5, 0}
)

const skillPressHitmark = 16

func init() {
	skillFrames = make([][]int, 5)

	// skill (press) -> x
	skillFrames[0] = frames.InitAbilSlice(42)
	skillFrames[0][action.ActionDash] = 22
	skillFrames[0][action.ActionJump] = 23
	skillFrames[0][action.ActionSwap] = 41

	// skill (hold=1) -> x
	skillFrames[1] = frames.InitAbilSlice(98)
	skillFrames[1][action.ActionBurst] = 97
	skillFrames[1][action.ActionDash] = 65
	skillFrames[1][action.ActionJump] = 66
	skillFrames[1][action.ActionSwap] = 96

	// skill (hold=1,c4) -> x
	skillFrames[2] = frames.InitAbilSlice(107)
	skillFrames[2][action.ActionDash] = 95
	skillFrames[2][action.ActionJump] = 95
	skillFrames[2][action.ActionSwap] = 106

	// skill (hold=2) -> x
	skillFrames[3] = frames.InitAbilSlice(343)
	skillFrames[3][action.ActionSkill] = 339 // uses burst frames
	skillFrames[3][action.ActionBurst] = 339
	skillFrames[3][action.ActionDash] = 231
	skillFrames[3][action.ActionJump] = 340
	skillFrames[3][action.ActionSwap] = 337

	// skill (hold=2,a4) -> x
	skillFrames[4] = frames.InitAbilSlice(175)
	skillFrames[4][action.ActionDash] = 171
	skillFrames[4][action.ActionJump] = 174
	skillFrames[4][action.ActionSwap] = 175
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	level, ok := p["hold"]
	if !ok || level < 0 || level > 2 {
		level = 0
	}

	c4Active := false
	if p["hold_c4"] == 1 && c.Base.Cons >= 4 {
		level = 1
		c4Active = true
	}

	if level != 0 {
		return c.skillHold(level, c4Active)
	}
	return c.skillPress()
}

func (c *char) skillPress() action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Passion Overload (Press)",
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeSlash,
		Element:            attributes.Pyro,
		Durability:         50,
		HitlagHaltFrames:   0.09 * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
		Mult:               skill[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTargetFanAngle(
			c.Core.Combat.Player(),
			combat.Point{Y: 0.8},
			2.5,
			270,
		),
		skillPressHitmark,
		skillPressHitmark,
	)

	//25 % chance of 3 orbs
	var count float64 = 2
	if c.Core.Rand.Float64() < .25 {
		count++
	}
	c.Core.QueueParticle("bennett", count, attributes.Pyro, skillPressHitmark+c.ParticleDelay)

	// a4 reduce cd by 50%
	if c.StatModIsActive(burstFieldKey) {
		//a4 reduces it from 240 to 120
		c.SetCDWithDelay(action.ActionSkill, 240/2, 14)
	} else {
		//default is 300, a2 reduces it by 20% to 240
		c.SetCDWithDelay(action.ActionSkill, 240, 14)
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames[0]),
		AnimationLength: skillFrames[0][action.InvalidAction],
		CanQueueAfter:   skillFrames[0][action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillHold(level int, c4Active bool) action.ActionInfo {

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Passion Overload (Level %v)", level),
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeSlash,
		Element:            attributes.Pyro,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
		Durability:         25,
	}

	for i, v := range skillHold[level-1] {
		ax := ai
		ax.Mult = v[c.TalentLvlSkill()]
		ax.HitlagHaltFrames = 0.09 * 60
		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			combat.Point{Y: skillHoldOffsets[i]},
			skillHoldHitboxes[i][0],
		)
		if i == 1 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				combat.Point{Y: skillHoldOffsets[i]},
				skillHoldHitboxes[i][0],
				skillHoldHitboxes[i][1],
			)
		}
		c.QueueCharTask(func() {
			c.Core.QueueAttack(ax, ap, 0, 0)
		}, skillHoldHitmarks[level-1][i])
	}
	if level == 2 {
		ai.StrikeType = combat.StrikeTypeDefault
		ai.Mult = explosion[c.TalentLvlSkill()]
		ai.HitlagHaltFrames = 0
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 1}, 3.5),
			166,
			166,
		)
	}

	//user-specified c4 variant adds an additional attack that deals 135% of the second hit
	if level == 1 && c4Active {
		ai.Mult = skillHold[level-1][1][c.TalentLvlSkill()] * 1.35
		ai.Abil = "Passion Overload (C4)"
		ai.HitlagHaltFrames = 0.12 * 60
		c.Core.QueueAttack(
			ai,
			combat.NewBoxHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: -1}, 3, 4),
			94,
			94,
		)
	}

	// TODO: particle timing??
	//Bennett Hold E is guaranteed 3 orbs
	c.Core.QueueParticle("bennett", 3, attributes.Pyro,
		skillHoldHitmarks[level-1][len(skillHoldHitmarks[level-1])-1]+c.ParticleDelay,
	)

	applyA4 := c.StatModIsActive(burstFieldKey)

	// figure out which frames to return
	// 0: skill (press) -> x
	// 1: skill (hold=1) -> x
	// 2: skill (hold=1,c4) -> x
	// 3: skill (hold=2) -> x
	// 4: skill (hold=2,a4) -> x
	idx := -1
	var cd, cdDelay int
	switch level {
	case 1:
		idx = 1
		cd = 450 - 90 //-90 for a2
		cdDelay = 43
		if c4Active {
			idx = 2
		}
	case 2:
		idx = 3
		cd = 600 - 120 //-120 from a2
		cdDelay = 110
		if applyA4 {
			idx = 4
		}
	default:
		panic("bennett skill (hold) level can only be 1 or 2")
	}

	// reduce cd by 50%
	if applyA4 {
		cd /= 2
	}
	c.SetCDWithDelay(action.ActionSkill, cd, cdDelay)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames[idx]),
		AnimationLength: skillFrames[idx][action.InvalidAction],
		CanQueueAfter:   skillFrames[idx][action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}
