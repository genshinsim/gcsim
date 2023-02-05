package razor

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var skillPressCDStarts = []int{30, 31}
var skillHoldCDStarts = []int{52, 52}

var skillPressHitmarks = []int{32, 33}
var skillHoldHitmarks = []int{55, 55}

var skillPressFrames [][]int
var skillHoldFrames [][]int

const (
	skillSigilDurationKey = "razor-sigil-duration"
)

func init() {
	// Tap E
	skillPressFrames = make([][]int, 2)

	// outside of Q
	skillPressFrames[0] = frames.InitAbilSlice(74) // Tap E -> Swap
	skillPressFrames[0][action.ActionAttack] = 70  // Tap E -> N1
	skillPressFrames[0][action.ActionBurst] = 69   // Tap E -> Q
	skillPressFrames[0][action.ActionDash] = 31    // Tap E -> D
	skillPressFrames[0][action.ActionJump] = 31    // Tap E -> J

	// inside of Q
	skillPressFrames[1] = frames.InitAbilSlice(76) // Tap E -> Swap
	skillPressFrames[1][action.ActionSwap] = 75    // Tap E -> N1
	skillPressFrames[1][action.ActionDash] = 32    // Tap E -> D
	skillPressFrames[1][action.ActionJump] = 32    // Tap E -> J

	// Hold E
	skillHoldFrames = make([][]int, 2)

	// outside of Q
	skillHoldFrames[0] = frames.InitAbilSlice(103) // Hold E -> Q
	skillHoldFrames[0][action.ActionAttack] = 102  // Hold E -> N1
	skillHoldFrames[0][action.ActionDash] = 52     // Hold E -> D
	skillHoldFrames[0][action.ActionJump] = 52     // Hold E -> J
	skillHoldFrames[0][action.ActionSwap] = 91     // Hold E -> Swap

	// inside of Q
	skillHoldFrames[1] = frames.InitAbilSlice(96) // Hold E -> N1
	skillHoldFrames[1][action.ActionDash] = 53    // Hold E -> D
	skillHoldFrames[1][action.ActionJump] = 52    // Hold E -> J
	skillHoldFrames[1][action.ActionSwap] = 88    // Hold E -> Swap
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	// check if Q is up for different E frames
	burstActive := 0
	if c.StatusIsActive(burstBuffKey) {
		burstActive = 1
	}

	if p["hold"] > 0 {
		return c.SkillHold(burstActive)
	}
	return c.SkillPress(burstActive)
}

func (c *char) SkillPress(burstActive int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Claw and Thunder (Press)",
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeSlash,
		Element:            attributes.Electro,
		Durability:         25,
		Mult:               skillPress[c.TalentLvlSkill()],
		HitlagHaltFrames:   0.1 * 60,
		HitlagFactor:       0.03,
		CanBeDefenseHalted: true,
	}

	var c4cb func(a combat.AttackCB)
	if c.Base.Cons >= 4 {
		c4cb = c.c4cb
	}

	radius := 2.4
	if c.StatusIsActive(burstBuffKey) {
		radius = 3
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTargetFanAngle(
			c.Core.Combat.Player(),
			combat.Point{Y: 1},
			radius,
			240,
		),
		skillPressHitmarks[burstActive],
		skillPressHitmarks[burstActive],
		c4cb,
		c.addSigil(false),
	)

	c.SetCDWithDelay(action.ActionSkill, c.a1CDReduction(6*60), skillPressCDStarts[burstActive])

	if !c.StatusIsActive(burstBuffKey) {
		//TODO: this delay used to be 80?
		c.Core.QueueParticle("razor", 3, attributes.Electro, skillPressHitmarks[burstActive]+c.ParticleDelay)
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillPressFrames[burstActive]),
		AnimationLength: skillPressFrames[burstActive][action.InvalidAction],
		CanQueueAfter:   skillPressFrames[burstActive][action.ActionDash], // earliest cancel is 1f before skillPressHitmark
		State:           action.SkillState,
	}
}

func (c *char) SkillHold(burstActive int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Claw and Thunder (Hold)",
		AttackTag:  combat.AttackTagElementalArtHold,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skillHold[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5),
		skillHoldHitmarks[burstActive],
		skillHoldHitmarks[burstActive],
	)

	c.Core.Tasks.Add(c.clearSigil, skillHoldHitmarks[burstActive])

	c.SetCDWithDelay(action.ActionSkill, c.a1CDReduction(10*60), skillHoldCDStarts[burstActive])

	if !c.StatusIsActive(burstBuffKey) {
		//TODO: this delay used to be 80?
		c.Core.QueueParticle("razor", 4, attributes.Electro, skillHoldHitmarks[burstActive]+c.ParticleDelay)
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillHoldFrames[burstActive]),
		AnimationLength: skillHoldFrames[burstActive][action.InvalidAction],
		CanQueueAfter:   skillHoldFrames[burstActive][action.ActionJump], // earliest cancel is 3f before skillHoldHitmark
		State:           action.SkillState,
	}
}

func (c *char) addSigil(done bool) combat.AttackCBFunc {
	return func(a combat.AttackCB) {
		if a.Target.Type() != combat.TargettableEnemy {
			return
		}
		if done {
			return
		}
		done = true
		if !c.StatusIsActive(skillSigilDurationKey) {
			c.sigils = 0
		}

		if c.sigils < 3 {
			c.sigils++
		}
		c.AddStatus(skillSigilDurationKey, 1080, true) //18 seconds
	}
}

func (c *char) clearSigil() {
	if !c.StatusIsActive(skillSigilDurationKey) {
		c.sigils = 0
		return
	}

	if c.sigils > 0 {
		c.AddEnergy("razor", float64(c.sigils)*5)
		c.sigils = 0
		c.DeleteStatus(skillSigilDurationKey)
	}
}

func (c *char) energySigil() {
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("er-sigil", -1),
		AffectedStat: attributes.ER,
		Amount: func() ([]float64, bool) {
			if c.StatusIsActive(skillSigilDurationKey) {
				c.skillSigilBonus[attributes.ER] = float64(c.sigils) * 0.2
				return c.skillSigilBonus, true
			}
			return nil, false
		},
	})
}
