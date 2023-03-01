package yaoyao

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var skillFrames []int

const (
	skillCDStart     = 15
	yueguiThrowSpawn = 48
)

func init() {
	skillFrames = frames.InitAbilSlice(52)
	skillFrames[action.ActionDash] = 49
	skillFrames[action.ActionJump] = 48
	skillFrames[action.ActionSwap] = 50
}

func (c *char) Skill(p map[string]int) action.ActionInfo {

	// yuegui spawns after 48f
	c.Core.Status.Add("yuegui", 600+yueguiThrowSpawn)

	procAI := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Radish (Skill)",
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagElementalArt,
		ICDGroup:           attacks.ICDGroupYaoyaoRadishSkill,
		StrikeType:         attacks.StrikeTypeDefault,
		Element:            attributes.Dendro,
		Durability:         25,
		Mult:               skillRadishDMG[c.TalentLvlSkill()],
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}

	c.Core.Tasks.Add(func() {
		yuegui := c.newYueguiThrow(procAI)
		c.Core.Combat.AddGadget(yuegui)
	}, skillCDStart+yueguiThrowSpawn)

	c.SetCDWithDelay(action.ActionSkill, 15*60, skillCDStart)

	if c.Base.Cons >= 4 {
		c.c4()
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionJump], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) getSkillHealInfo(snap *combat.Snapshot) player.HealInfo {
	maxhp := snap.BaseHP*(1+snap.Stats[attributes.HPP]) + snap.Stats[attributes.HP]
	heal := skillRadishHealing[0][c.TalentLvlSkill()]*maxhp + skillRadishHealing[1][c.TalentLvlSkill()]
	return player.HealInfo{
		Caller:  c.Index,
		Message: "Yuegui Skill",
		Src:     heal,
		Bonus:   snap.Stats[attributes.Heal],
	}
}
