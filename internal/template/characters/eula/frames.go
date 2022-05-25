package eula

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		var f, a int
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 30
			a = 34
		case 1:
			f = 19
			a = 36
		case 2:
			f = 42
			a = 56
		case 3:
			f = 17
			a = 44
		case 4:
			f = 56
			a = 106
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, a
	case core.ActionCharge:
		return 35, 35 //TODO: no idea
	case core.ActionSkill:
		if p["hold"] != 0 {
			return 76, 100
		}
		return 30, 48
	case core.ActionBurst:
		return 117, 122
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {
	// NA/CA cancels. TODO: CA cancels?
	c.SetNormalCancelFrames(0, core.ActionAttack, 34-30-1)  // n1 -> next attack
	c.SetNormalCancelFrames(1, core.ActionAttack, 36-19-1)  // n2 -> next attack
	c.SetNormalCancelFrames(2, core.ActionAttack, 56-42-1)  // n3 -> next attack
	c.SetNormalCancelFrames(3, core.ActionAttack, 44-17-1)  // n4 -> next attack
	c.SetNormalCancelFrames(4, core.ActionAttack, 106-56-1) // n5 -> next attack

	// skill -> x
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionWalk, 48-30)

	// skill (hold) -> x
	c.SetAbilCancelFrames(core.ActionSkillHoldFramesOnly, core.ActionSwap, 100-76)

	// burst -> x
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 122-117)
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionCharge, 122-117) // uses n1 frames
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 122-117)
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 122-117)
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 122-117)
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 122-117)
}
