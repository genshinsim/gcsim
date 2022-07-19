package traveleranemo

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillPressFrames []int
var skillHoldFrames [7][]int

func init() {
	skillPressFrames = frames.InitAbilSlice(81) // default is walk frames
	skillPressFrames[action.ActionAttack] = 61
	skillPressFrames[action.ActionSkill] = 60 // uses burst frames
	skillPressFrames[action.ActionBurst] = 60
	skillPressFrames[action.ActionDash] = 28
	skillPressFrames[action.ActionJump] = 28
	skillPressFrames[action.ActionSwap] = 60
	for i := 1; i <= 6; i += 1 {
		max_dur := 31 + (i-1)*15 + 5
		if i >= 2 {
			max_dur += 5
		}
		skillHoldFrames[i] = frames.InitAbilSlice(max_dur + 103 - 55) // default is walk frames
		skillHoldFrames[i][action.ActionAttack] = max_dur + 82 - 55
		skillHoldFrames[i][action.ActionSkill] = max_dur + 83 - 55 // uses burst frames
		skillHoldFrames[i][action.ActionBurst] = max_dur + 83 - 55
		skillHoldFrames[i][action.ActionDash] = max_dur + 55 - 55
		skillHoldFrames[i][action.ActionJump] = max_dur + 55 - 55
		skillHoldFrames[i][action.ActionSwap] = max_dur + 82 - 55
	}
}

func (c *char) SkillPress() {
	hitmark := 34
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Palm Vortex (Tap)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       skillInitialStorm[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), hitmark, hitmark)

	c.Core.QueueParticle(c.Base.Key.String(), 2, attributes.Anemo, hitmark+c.Core.Flags.ParticleDelay)
	c.SetCDWithDelay(action.ActionSkill, 5*60, hitmark-5)
}
func (c *char) SkillHold(holdTicks int) {

	c.eInfuse = attributes.NoElement
	c.eICDTag = combat.ICDTagNone

	aiCut := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Palm Vortex Initial Cutting (Hold)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       skillInitialCutting[c.TalentLvlSkill()],
	}

	aiCutAbs := aiCut
	aiCutAbs.Abil = "Palm Vortex Initial Cutting Absorbed (Hold)"
	aiCutAbs.ICDTag = combat.ICDTagNone
	aiCutAbs.Element = attributes.NoElement
	aiCutAbs.Mult = skillInitialCuttingAbsorb[c.TalentLvlSkill()]

	aiMaxCutAbs := aiCutAbs
	aiMaxCutAbs.Abil = "Palm Vortex Max Cutting Absorbed (Hold)"
	aiMaxCutAbs.Mult = skillMaxCuttingAbsorb[c.TalentLvlSkill()]

	// first tick is at 31f, with 15f between ticks, and an extra 5 frame delay when transitioning from Initial to Max
	firstTick := 31
	hitmark := firstTick
	for i := 0; i < holdTicks; i += 1 {

		c.Core.QueueAttack(aiCut, combat.NewDefCircHit(1, false, combat.TargettableEnemy), hitmark, hitmark)
		if i > 1 {
			c.Core.Tasks.Add(func() {
				if c.eInfuse != attributes.NoElement {
					aiMaxCutAbs.Element = c.eInfuse
					aiMaxCutAbs.ICDTag = c.eICDTag
					c.Core.QueueAttack(aiMaxCutAbs, combat.NewDefCircHit(1.5, false, combat.TargettableEnemy), 0, 0)
				}
				//check if infused
			}, hitmark)
		} else {
			c.Core.Tasks.Add(func() {
				if c.eInfuse != attributes.NoElement {
					aiCutAbs.Element = c.eInfuse
					aiCutAbs.ICDTag = c.eICDTag
					c.Core.QueueAttack(aiCutAbs, combat.NewDefCircHit(1.5, false, combat.TargettableEnemy), 0, 0)
				}
				//check if infused
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
		c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Anemo, hitmark+90)
		c.SetCDWithDelay(action.ActionSkill, 8*60, hitmark-5)
	} else {
		c.Core.QueueParticle(c.Base.Key.String(), 2, attributes.Anemo, hitmark+90)
		c.SetCDWithDelay(action.ActionSkill, 5*60, hitmark-5)
	}

	c.Core.QueueAttack(aiStorm, combat.NewDefCircHit(2, false, combat.TargettableEnemy), hitmark, hitmark)
	c.Core.Tasks.Add(func() {
		if c.eInfuse != attributes.NoElement {
			aiStormAbs.Element = c.eInfuse
			aiStormAbs.ICDTag = c.eICDTag
			c.Core.QueueAttack(aiStormAbs, combat.NewDefCircHit(1.5, false, combat.TargettableEnemy), 0, 0)
		}
		//check if infused
	}, hitmark)

	// starts absorbing after the first tick?
	c.Core.Tasks.Add(c.absorbCheckE(c.Core.F, 0, int((hitmark)/18)), firstTick+1)
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
		c.SkillPress()
		return action.ActionInfo{
			Frames:          frames.NewAbilFunc(skillPressFrames),
			AnimationLength: skillPressFrames[action.InvalidAction],
			CanQueueAfter:   skillPressFrames[action.ActionDash], // earliest cancel
			State:           action.SkillState,
		}
	} else {
		c.SkillHold(holdTicks)
		return action.ActionInfo{
			Frames:          frames.NewAbilFunc(skillHoldFrames[holdTicks]),
			AnimationLength: skillHoldFrames[holdTicks][action.InvalidAction],
			CanQueueAfter:   skillHoldFrames[holdTicks][action.ActionDash], // earliest cancel
			State:           action.SkillState,
		}
	}
}

func (c *char) absorbCheckE(src, count, max int) func() {
	return func() {
		if count == max {
			return
		}
		c.eInfuse = c.Core.Combat.AbsorbCheck(c.infuseCheckLocation, attributes.Cryo, attributes.Pyro, attributes.Hydro, attributes.Electro)
		switch c.eInfuse {
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
