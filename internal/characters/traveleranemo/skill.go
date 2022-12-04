package traveleranemo

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillPressFrames [][]int
var skillHoldDelayFrames [][]int

func init() {
	// Tap E
	skillPressFrames = make([][]int, 2)

	// Male
	skillPressFrames[0] = frames.InitAbilSlice(74) // Tap E -> N1
	skillPressFrames[0][action.ActionBurst] = 76   // Tap E -> Q
	skillPressFrames[0][action.ActionDash] = 30    // Tap E -> D
	skillPressFrames[0][action.ActionJump] = 31    // Tap E -> J
	skillPressFrames[0][action.ActionSwap] = 66    // Tap E -> Swap

	// Female
	skillPressFrames[1] = frames.InitAbilSlice(62) // Tap E -> Q
	skillPressFrames[1][action.ActionAttack] = 61  // Tap E -> N1
	skillPressFrames[1][action.ActionDash] = 31    // Tap E -> D
	skillPressFrames[1][action.ActionJump] = 31    // Tap E -> J
	skillPressFrames[1][action.ActionSwap] = 60    // Tap E -> Swap

	// Short Hold E as base for Hold E frames
	// "2 tick duration - 2 tick last hitmark"
	skillHoldDelayFrames = make([][]int, 2)

	// Male
	skillHoldDelayFrames[0] = frames.InitAbilSlice(98 - 54) // Short Hold E -> N1/Q - Short Hold E -> D
	skillHoldDelayFrames[0][action.ActionDash] = 0          // Short Hold E -> D - Short Hold E -> D
	skillHoldDelayFrames[0][action.ActionJump] = 0          // Short Hold E -> J - Short Hold E -> D
	skillHoldDelayFrames[0][action.ActionSwap] = 89 - 54    // Short Hold E -> Swap - Short Hold E -> D

	// Female
	skillHoldDelayFrames[1] = frames.InitAbilSlice(84 - 54) // Short Hold E -> Q - Short Hold E -> D
	skillHoldDelayFrames[1][action.ActionAttack] = 83 - 54  // Short Hold E -> N1 - Short Hold E -> D
	skillHoldDelayFrames[1][action.ActionDash] = 0          // Short Hold E -> D - Short Hold E -> D
	skillHoldDelayFrames[1][action.ActionJump] = 0          // Short Hold E -> J - Short Hold E -> D
	skillHoldDelayFrames[1][action.ActionSwap] = 83 - 54    // Short Hold E -> Swap - Short Hold E -> D
}

func (c *char) SkillPress() action.ActionInfo {
	hitmark := 34
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Palm Vortex (Tap)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       skillInitialStorm[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTargetFanAngle(c.Core.Combat.Player(), nil, 6, 100),
		hitmark,
		hitmark,
	)

	c.Core.QueueParticle(c.Base.Key.String(), 2, attributes.Anemo, hitmark+c.ParticleDelay)
	c.SetCDWithDelay(action.ActionSkill, 5*60, hitmark-5)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillPressFrames[c.gender]),
		AnimationLength: skillPressFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   skillPressFrames[c.gender][action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) SkillHold(holdTicks int) action.ActionInfo {

	c.eAbsorb = attributes.NoElement
	c.eICDTag = combat.ICDTagNone
	c.eAbsorbCheckLocation = combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 1.2}, 3)

	aiCut := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Palm Vortex Initial Cutting (Hold)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeSlash,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       skillInitialCutting[c.TalentLvlSkill()],
	}

	aiCutAbs := aiCut
	aiCutAbs.Abil = "Palm Vortex Initial Cutting Absorbed (Hold)"
	aiCutAbs.ICDTag = combat.ICDTagNone
	aiCutAbs.StrikeType = combat.StrikeTypeDefault
	aiCutAbs.Element = attributes.NoElement
	aiCutAbs.Mult = skillInitialCuttingAbsorb[c.TalentLvlSkill()]

	aiMaxCutAbs := aiCutAbs
	aiMaxCutAbs.Abil = "Palm Vortex Max Cutting Absorbed (Hold)"
	aiMaxCutAbs.Mult = skillMaxCuttingAbsorb[c.TalentLvlSkill()]

	// first tick is at 31f, with 15f between ticks, and an extra 5 frame delay when transitioning from Initial to Max
	firstTick := 31
	hitmark := firstTick
	for i := 0; i < holdTicks; i += 1 {

		c.Core.QueueAttack(
			aiCut,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 1.2}, 1.7),
			hitmark,
			hitmark,
		)
		if i > 1 {
			c.Core.Tasks.Add(func() {
				if c.eAbsorb != attributes.NoElement {
					aiMaxCutAbs.Element = c.eAbsorb
					aiMaxCutAbs.ICDTag = c.eICDTag
					c.Core.QueueAttack(
						aiMaxCutAbs,
						combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 1.2}, 3.6),
						0,
						0,
					)
				}
				//check if absorbed
			}, hitmark)
		} else {
			c.Core.Tasks.Add(func() {
				if c.eAbsorb != attributes.NoElement {
					aiCutAbs.Element = c.eAbsorb
					aiCutAbs.ICDTag = c.eICDTag
					c.Core.QueueAttack(
						aiCutAbs,
						combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 1.2}, 1.7),
						0,
						0,
					)
				}
				//check if absorbed
			}, hitmark)
		}

		// go to next tick
		hitmark += 15
		if i == 1 {
			aiCut.Mult = skillMaxCutting[c.TalentLvlSkill()]
			aiCut.Abil = "Palm Vortex Max Cutting (Hold)"

			// there is a 5 frame delay when it shifts from initial to max
			hitmark += 5
		}

	}
	// move the hitmark back by 1 tick (15f) then forward by 5f for the Storm damage
	hitmark = hitmark - 15 + 5
	aiStorm := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Palm Vortex Initial Storm (Hold)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       skillInitialStorm[c.TalentLvlSkill()],
	}

	aiStormAbs := aiStorm
	aiStormAbs.Abil = "Palm Vortex Initial Storm Absorbed (Hold)"
	aiStormAbs.ICDTag = combat.ICDTagNone
	aiStormAbs.Element = attributes.NoElement
	aiStormAbs.Mult = skillInitialStormAbsorb[c.TalentLvlSkill()]

	// it does max storm when there are 2 or more ticks
	if holdTicks >= 2 {
		aiStorm.Mult = skillMaxStorm[c.TalentLvlSkill()]
		aiStorm.Abil = "Palm Vortex Max Storm (Hold)"

		aiStormAbs.Mult = skillMaxStormAbsorb[c.TalentLvlSkill()]
		aiStormAbs.Abil = "Palm Vortex Max Storm Absorbed (Hold)"

		count := 3.0
		if c.Core.Rand.Float64() < 0.33 {
			count = 4
		}
		c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Anemo, hitmark+c.ParticleDelay)
		c.SetCDWithDelay(action.ActionSkill, 8*60, hitmark-5)
	} else {
		c.Core.QueueParticle(c.Base.Key.String(), 2, attributes.Anemo, hitmark+c.ParticleDelay)
		c.SetCDWithDelay(action.ActionSkill, 5*60, hitmark-5)
	}

	c.Core.QueueAttack(
		aiStorm,
		combat.NewCircleHitOnTargetFanAngle(c.Core.Combat.Player(), nil, 6, 100),
		hitmark,
		hitmark,
	)
	c.Core.Tasks.Add(func() {
		if c.eAbsorb != attributes.NoElement {
			aiStormAbs.Element = c.eAbsorb
			aiStormAbs.ICDTag = c.eICDTag
			c.Core.QueueAttack(
				aiStormAbs,
				combat.NewCircleHitOnTargetFanAngle(c.Core.Combat.Player(), nil, 6, 100),
				0,
				0,
			)
		}
		//check if absorbed
	}, hitmark)

	// starts absorbing after the first tick?
	c.Core.Tasks.Add(c.absorbCheckE(c.Core.F, 0, int((hitmark)/18)), firstTick+1)
	return action.ActionInfo{
		Frames:          func(next action.Action) int { return skillHoldDelayFrames[c.gender][next] + hitmark },
		AnimationLength: skillHoldDelayFrames[c.gender][action.InvalidAction] + hitmark,
		CanQueueAfter:   skillHoldDelayFrames[c.gender][action.ActionDash] + hitmark, // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	holdTicks := 0
	if p["hold"] == 1 {
		holdTicks = 6
	}
	if 0 < p["hold_ticks"] {
		holdTicks = p["hold_ticks"]
	}
	if holdTicks > 6 {
		holdTicks = 6
	}

	if holdTicks == 0 {
		return c.SkillPress()
	} else {
		return c.SkillHold(holdTicks)
	}
}

func (c *char) absorbCheckE(src, count, max int) func() {
	return func() {
		if count == max {
			return
		}
		c.eAbsorb = c.Core.Combat.AbsorbCheck(c.eAbsorbCheckLocation, attributes.Cryo, attributes.Pyro, attributes.Hydro, attributes.Electro)
		switch c.eAbsorb {
		case attributes.Cryo:
			c.eICDTag = combat.ICDTagElementalArtCryo
		case attributes.Pyro:
			c.eICDTag = combat.ICDTagElementalArtPyro
		case attributes.Electro:
			c.eICDTag = combat.ICDTagElementalArtElectro
		case attributes.Hydro:
			c.eICDTag = combat.ICDTagElementalArtHydro
		case attributes.NoElement:
			//otherwise queue up
			c.Core.Tasks.Add(c.absorbCheckE(src, count+1, max), 18)
		}
	}
}
