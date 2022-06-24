package diona

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

var skillPressFrames []int
var skillHoldFrames []int

const skillPressAnimation = 15
const skillHoldAnimation = 24

func init() {
	// skill (press) -> x
	skillPressFrames = frames.InitAbilSlice(skillPressAnimation)

	// skill (hold) -> x
	skillHoldFrames = frames.InitAbilSlice(skillHoldAnimation)
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	if p["hold"] == 1 {
		return c.skillHold(travel)
	}
	return c.skillPress(travel)
}

func (c *char) skillPress(travel int) action.ActionInfo {
	c.pawsPewPew(skillPressAnimation, travel, 2)
	c.SetCD(action.ActionSkill, 360+skillPressAnimation)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressAnimation,
		CanQueueAfter:   skillPressAnimation,
		State:           action.SkillState,
	}
}

func (c *char) skillHold(travel int) action.ActionInfo {
	c.pawsPewPew(skillHoldAnimation, travel, 5)
	c.SetCD(action.ActionSkill, 900+skillHoldAnimation)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillHoldFrames),
		AnimationLength: skillHoldAnimation,
		CanQueueAfter:   skillHoldAnimation,
		State:           action.SkillState,
	}
}

func (c *char) pawsPewPew(f, travel, pawCount int) {
	bonus := 1.0
	if pawCount == 5 {
		bonus = 1.75 //bonus if firing off 5
	}
	shdHp := (pawShieldPer[c.TalentLvlSkill()]*c.MaxHP() + pawShieldFlat[c.TalentLvlSkill()]) * bonus
	if c.Base.Cons >= 2 {
		shdHp = shdHp * 1.15
	}
	//call back to generate shield on hit
	//note that each paw should only be able to trigger callback once (if hit multi target)
	//and that subsequent shield generation should increase duation only
	//TODO: need to look into maybe additional paw hits actually create "new" shields?
	pawCB := func(done bool) combat.AttackCBFunc {
		return func(acb combat.AttackCB) {
			if done {
				return
			}
			//make sure this is only triggered once
			done = true

			//trigger particles if prob < 0.8
			if c.Core.Rand.Float64() < 0.8 {
				c.Core.QueueParticle("Diona", 1, attributes.Cryo, 90) //90s travel time
			}

			//check if shield already exists, if so then just update duration
			exist := c.Core.Player.Shields.Get(shield.ShieldDionaSkill)
			var shd *shield.Tmpl
			if exist != nil {
				//update
				shd, _ = exist.(*shield.Tmpl)
				shd.Expires = shd.Expires + pawDur[c.TalentLvlSkill()]
			} else {
				shd = &shield.Tmpl{
					Src:        c.Core.F,
					ShieldType: shield.ShieldDionaSkill,
					Name:       "Diona Skill",
					HP:         shdHp,
					Ele:        attributes.Cryo,
					Expires:    c.Core.F + pawDur[c.TalentLvlSkill()], //15 sec
				}
			}
			//TODO: check that this is actually properly extending duration
			c.Core.Player.Shields.Add(shd)
		}
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Icy Paw",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       paw[c.TalentLvlSkill()],
	}

	for i := 0; i < pawCount; i++ {
		done := false
		cb := pawCB(done)
		c.Core.QueueAttack(ai, combat.NewDefSingleTarget(1, combat.TargettableEnemy), 0, travel+f-5+i, cb)
	}
}
