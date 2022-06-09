package eula

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var skillPressFrames []int
var skillHoldFrames []int
var icewhirlHitmarks = []int{79, 92}

const skillPressHitmark = 20
const skillHoldHitmark = 49

func (c *char) Skill(p map[string]int) action.ActionInfo {
	if p["hold"] != 0 {
		return c.holdSkill(p)
	}
	return c.pressSkill(p)
}

func (c *char) pressSkill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Icetide Vortex",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       skillPress[c.TalentLvlSkill()],
	}
	//add 1 to grim heart if not capped by icd
	cb := func(a combat.AttackCB) {
		if c.Core.F < c.grimheartICD {
			return
		}
		c.grimheartICD = c.Core.F + 18

		if c.Tags["grimheart"] < 2 {
			c.Tags["grimheart"]++
			c.Core.Log.NewEvent("eula: grimheart stack", glog.LogCharacterEvent, c.Index, "current count", c.Tags["grimheart"])
		}
		c.grimheartReset = 18 * 60
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(1.5, false, combat.TargettableEnemy), skillPressHitmark, skillPressHitmark, cb)

	var count float64 = 1
	if c.Core.Rand.Float64() < .5 {
		count = 2
	}
	c.Core.QueueParticle("eula", count, attributes.Cryo, 20+100)

	c.SetCDWithDelay(action.ActionSkill, 60*4, 16)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.ActionDash], // earliest cancel
		Post:            skillPressFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) holdSkill(p map[string]int) action.ActionInfo {
	//hold e
	//296 to 341, but cd starts at 322
	//60 fps = 108 frames cast, cd starts 62 frames in so need to + 62 frames to cd
	lvl := c.TalentLvlSkill()
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Icetide Vortex (Hold)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       skillHold[lvl],
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(1.5, false, combat.TargettableEnemy), skillHoldHitmark, skillHoldHitmark)

	//multiple brand hits
	ai.Abil = "Icetide Vortex (Icewhirl)"
	ai.ICDTag = combat.ICDTagElementalArt
	ai.StrikeType = combat.StrikeTypeDefault
	ai.Mult = icewhirl[lvl]

	v := c.Tags["grimheart"]

	//shred
	var shredCB combat.AttackCBFunc
	if v > 0 {
		done := false
		shredCB = func(a combat.AttackCB) {
			if done {
				return
			}
			e, ok := a.Target.(core.Enemy)
			if !ok {
				return
			}
			done = true
			e.AddResistMod("eula-icewhirl-shred-cryo", 7*v*60, attributes.Cryo, -resRed[lvl])
			e.AddResistMod("eula-icewhirl-shred-phys", 7*v*60, attributes.Physical, -resRed[lvl])
		}
	}

	// this shouldn't happen, but to be safe
	if v > 2 {
		v = 2
	}
	for i := 0; i < v; i++ {
		//spacing it out for stacks
		c.Core.QueueAttack(
			ai,
			combat.NewDefCircHit(1.5, false, combat.TargettableEnemy),
			icewhirlHitmarks[i],
			icewhirlHitmarks[i],
			shredCB,
		)
	}

	//A1
	if v == 2 {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Icetide (Lightfall)",
			AttackTag:  combat.AttackTagElementalBurst,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeBlunt,
			Element:    attributes.Physical,
			Durability: 25,
			Mult:       burstExplodeBase[c.TalentLvlBurst()] * 0.5,
		}
		c.Core.QueueAttack(ai, combat.NewDefCircHit(1.5, false, combat.TargettableEnemy), 108, 108)
	}

	var count float64 = 2
	if c.Core.Rand.Float64() < .5 {
		count = 3
	}
	c.Core.QueueParticle("eula", count, attributes.Cryo, skillHoldHitmark+100)

	//c1 add debuff
	if c.Base.Cons >= 1 && v > 0 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.PhyP] = 0.3
		//TODO: check if the duration is right
		c.AddStatMod("eula-c1", (6*v+6)*60, attributes.PhyP, func() ([]float64, bool) {
			return m, true
		})
	}

	c.Tags["grimheart"] = 0
	cd := 10
	if c.Base.Cons >= 2 {
		cd = 4 //press and hold have same cd TODO: check if this is right
	}
	c.SetCDWithDelay(action.ActionSkill, cd*60, 46)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillHoldFrames),
		AnimationLength: skillHoldFrames[action.InvalidAction],
		CanQueueAfter:   skillHoldFrames[action.ActionDash], // earliest cancel
		Post:            skillHoldFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}
