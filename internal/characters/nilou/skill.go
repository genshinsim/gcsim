package nilou

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

var skillFrames []int

var skillDanceFrames [][]int
var skillDanceHitMarks = []int{14, 12, 35}

var whirlingStepsFrames [][]int
var whirlingStepsHitMarks = []int{21, 29, 43}

type NilouSkillType int

const (
	NilouSkillTypeNone  NilouSkillType = iota
	NilouSkillTypeDance                // NA
	NilouSkillTypeSteps                // Skill
)

const (
	pirouetteStatus = "pirouette"
	// skillTypeTag          = "skill_type"
	skillStep             = "skill_step"
	lunarPrayerStatus     = "lunarprayer"
	tranquilityAuraStatus = "tranquilityaura"

	skillHitmark = 16 // init

	delayDance = 40 // Lunar Prayer (8s) / Tranquility (12/18s) / A1 (30s) timers all start here
	delaySteps = 30
)

// TODO: cancel frames
func init() {
	skillFrames = frames.InitAbilSlice(22) // E -> Q
	skillFrames[action.ActionAttack] = 19
	skillFrames[action.ActionSkill] = 19
	skillFrames[action.ActionDash] = 17
	skillFrames[action.ActionJump] = 17
	skillFrames[action.ActionSwap] = 21

	skillDanceFrames = make([][]int, normalHitNum)
	skillDanceFrames[0] = frames.InitNormalCancelSlice(skillDanceHitMarks[0], 20) // N1 -> E
	skillDanceFrames[0][action.ActionAttack] = 18
	skillDanceFrames[0][action.ActionBurst] = skillDanceHitMarks[0]
	skillDanceFrames[0][action.ActionDash] = skillDanceHitMarks[0]
	skillDanceFrames[0][action.ActionJump] = skillDanceHitMarks[0]
	skillDanceFrames[0][action.ActionSwap] = skillDanceHitMarks[0]

	skillDanceFrames[1] = frames.InitNormalCancelSlice(skillDanceHitMarks[1], 23) // N2 -> N3/E
	skillDanceFrames[1][action.ActionBurst] = skillDanceHitMarks[1]
	skillDanceFrames[1][action.ActionDash] = skillDanceHitMarks[1]
	skillDanceFrames[1][action.ActionJump] = skillDanceHitMarks[1]
	skillDanceFrames[1][action.ActionSwap] = skillDanceHitMarks[1]

	skillDanceFrames[2] = frames.InitNormalCancelSlice(skillDanceHitMarks[2], 60) // N3 -> E/Q
	skillDanceFrames[2][action.ActionAttack] = 55
	skillDanceFrames[2][action.ActionDash] = skillDanceHitMarks[2]
	skillDanceFrames[2][action.ActionJump] = skillDanceHitMarks[2]
	skillDanceFrames[2][action.ActionSwap] = skillDanceHitMarks[2]

	whirlingStepsFrames = make([][]int, normalHitNum)
	whirlingStepsFrames[0] = frames.InitNormalCancelSlice(whirlingStepsHitMarks[0], 33) // E1 -> Q
	whirlingStepsFrames[0][action.ActionAttack] = 27
	whirlingStepsFrames[0][action.ActionSkill] = 27
	whirlingStepsFrames[0][action.ActionDash] = 26
	whirlingStepsFrames[0][action.ActionJump] = 27
	whirlingStepsFrames[0][action.ActionSwap] = 31

	whirlingStepsFrames[1] = frames.InitNormalCancelSlice(whirlingStepsHitMarks[1], 62) // E2 -> Swap
	whirlingStepsFrames[1][action.ActionAttack] = 40
	whirlingStepsFrames[1][action.ActionSkill] = 32
	whirlingStepsFrames[1][action.ActionBurst] = 40
	whirlingStepsFrames[1][action.ActionDash] = 36
	whirlingStepsFrames[1][action.ActionJump] = 37

	whirlingStepsFrames[2] = frames.InitNormalCancelSlice(whirlingStepsHitMarks[2], 63) // E3 -> N1/E/Q
	whirlingStepsFrames[2][action.ActionDash] = 57
	whirlingStepsFrames[2][action.ActionJump] = 57
	whirlingStepsFrames[0][action.ActionSwap] = 61
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	if c.StatusIsActive(pirouetteStatus) {
		return c.Pirouette(p, NilouSkillTypeSteps)
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Dance of Haftkarsvar",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypePierce,
		Element:    attributes.Hydro,
		Durability: 25,
		FlatDmg:    skill[c.TalentLvlSkill()] * c.MaxHP(),
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 1, false, combat.TargettableEnemy), skillHitmark, skillHitmark)

	c.SetTag(skillStep, 0)
	// c.SetTag(skillTypeTag, int(NilouSkillTypeNone))
	c.AddStatus(pirouetteStatus, 10*60, true)
	c.SetCD(action.ActionSkill, 18*60)

	var orb float64 = 1
	if c.Core.Rand.Float64() < 0.5 {
		orb = 2
	}
	c.Core.QueueParticle("nilou", orb, attributes.Hydro, skillHitmark+c.ParticleDelay)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionJump], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) Pirouette(p map[string]int, srcType NilouSkillType) action.ActionInfo {
	// c.SetTag(skillTypeTag, int(srcType))

	ai := action.ActionInfo{}
	delay := 0
	switch srcType {
	case NilouSkillTypeDance:
		ai = c.SwordDance(p)
		delay = delayDance
	case NilouSkillTypeSteps:
		ai = c.WhirlingSteps(p)
		delay = delaySteps
	}

	if c.Tag(skillStep) == 0 {
		c.QueueCharTask(func() {
			c.a1()
			c.DeleteStatus(pirouetteStatus)

			switch srcType {
			case NilouSkillTypeDance:
				c.AddStatus(lunarPrayerStatus, 8*60, true)
			case NilouSkillTypeSteps:
				dur := 12 * 60
				if c.Base.Cons >= 1 {
					dur += 6 * 60
				}
				c.AddStatus(tranquilityAuraStatus, dur, true)
				c.QueueCharTask(c.TranquilityAura, 30) // every 0.25 sec
			}
		}, delay)
	}

	return ai
}

func (c *char) AdvanceSkillIndex() {
	s := c.Tag(skillStep) + 1
	if s == 3 {
		s = 0
	}
	c.SetTag(skillStep, s)
}

func (c *char) SwordDance(p map[string]int) action.ActionInfo {
	s := c.Tag(skillStep)
	travel := 0

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Sword Dance %v", s),
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypePierce,
		Element:    attributes.Hydro,
		Durability: 25,
		FlatDmg:    swordDance[s][c.TalentLvlSkill()] * c.MaxHP(),
	}
	if s == 2 {
		ai.Abil = "Luminous Illusion"

		if t, ok := p["travel"]; ok {
			travel = t
		}
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), 1, false, combat.TargettableEnemy),
		skillDanceHitMarks[s]+travel,
		skillDanceHitMarks[s]+travel,
		c.c4cb(),
	)

	if c.StatusIsActive(pirouetteStatus) {
		c.Core.QueueParticle("nilou", 1, attributes.Hydro, skillDanceHitMarks[s]+travel+c.ParticleDelay)
	}

	defer c.AdvanceSkillIndex()

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillDanceFrames[s]), // TODO: atkspd can it affect?
		AnimationLength: skillDanceFrames[s][action.InvalidAction],
		CanQueueAfter:   skillDanceFrames[s][action.ActionJump], // earliest cancel
		State:           action.NormalAttackState,
	}
}

func (c *char) WhirlingSteps(p map[string]int) action.ActionInfo {
	s := c.Tag(skillStep)

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Whirling Steps %v", s),
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypePierce,
		Element:    attributes.Hydro,
		Durability: 25,
		FlatDmg:    whirlingSteps[s][c.TalentLvlSkill()] * c.MaxHP(),
	}
	if s == 2 {
		ai.Abil = "Water Wheel"
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), 1, false, combat.TargettableEnemy),
		whirlingStepsHitMarks[s],
		whirlingStepsHitMarks[s],
		c.c4cb(),
	)

	if c.StatusIsActive(pirouetteStatus) {
		c.Core.QueueParticle("nilou", 1, attributes.Hydro, whirlingStepsHitMarks[s]+c.ParticleDelay)
	}

	defer c.AdvanceSkillIndex()

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(whirlingStepsFrames[s]), // TODO: atkspd can it affect?
		AnimationLength: whirlingStepsFrames[s][action.InvalidAction],
		CanQueueAfter:   whirlingStepsFrames[s][action.ActionJump], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) TranquilityAura() {
	if !c.StatusIsActive(tranquilityAuraStatus) {
		return
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Tranquility Aura",
		AttackTag:  combat.AttackTagNone,
		ICDTag:     combat.ICDTagNilouTranquilityAura,
		ICDGroup:   combat.ICDGroupNilou,
		Element:    attributes.Hydro,
		Durability: 25,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 3, false, combat.TargettableEnemy), -1, 1)

	c.QueueCharTask(c.TranquilityAura, 30)
}

// Clears Nilou skill when she leaves the field
func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		c.DeleteStatus(pirouetteStatus)
		c.DeleteStatus(lunarPrayerStatus)
		c.SetTag(skillStep, 0)
		return false
	}, "nilou-exit")
}
