package raiden

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
)

func (c *char) InitCancelFrames() {
	c.initNormalCancels()
	c.initBurstAttackCancels()

	c.SetAbilCancelFrames(core.ActionCharge, core.ActionAttack, 37-22) //charge -> n1
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSkill, 37-22)  //charge -> skill
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionBurst, 37-22)  //charge -> burst
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSwap, 36-22)   //charge -> swap

	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 111-98) //burst -> n1
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 111-98)  //burst -> skill
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 110-98)   //burst -> dash
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 112-98)   //burst -> jump
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 110-98)   //burst -> swap

	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 37-17)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 37-17)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionDash, 17-17)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionJump, 17-17)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 36-17)
}

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		a := 0
		if c.Core.Status.Duration("raidenburst") == 0 {
			switch c.NormalCounter {
			//TODO: need to add atkspd mod
			case 0:
				f = 14
				a = 18
			case 1:
				f = 9
				a = 13
			case 2:
				f = 14
				a = 26
			case 3:
				f = 27
				a = 41
			case 4:
				f = 34
				a = 50
			}
		} else {
			switch c.NormalCounter {
			//TODO: need to add atkspd mod
			case 0:
				f = 12
				a = 19
			case 1:
				f = 13
				a = 16
			case 2:
				f = 11
				a = 16
			case 3:
				f = 33
				a = 44
			case 4:
				f = 33
				a = 59
			}
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, a
	case core.ActionCharge:
		if c.Core.Status.Duration("raidenburst") == 0 {
			return 22, 37
		}
		return 24, 56
	case core.ActionSkill:
		return 17, 37
	case core.ActionBurst:
		return 98, 111
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

var attackFrames [][]int
var hitmarks = [][]int{{14}, {9}, {14}, {14, 27}, {34}}

func (c *char) initNormalCancels() {

	//normal cancels
	attackFrames = make([][]int, c.NormalHitNum) //should be 5

	//n1 animations
	setNormCancel(attackFrames, 0, hitmarks[0][0], 24)
	attackFrames[0][action.ActionAttack] = 18
	attackFrames[0][action.ActionCharge] = 24

	//n2 animations
	setNormCancel(attackFrames, 1, hitmarks[1][0], 26)
	attackFrames[1][action.ActionAttack] = 13
	attackFrames[1][action.ActionCharge] = 26

	//n3 animations
	setNormCancel(attackFrames, 2, hitmarks[2][0], 36)
	attackFrames[2][action.ActionAttack] = 26
	attackFrames[2][action.ActionCharge] = 36

	//n4 animations
	setNormCancel(attackFrames, 3, hitmarks[3][1], 57)
	attackFrames[3][action.ActionAttack] = 41
	attackFrames[3][action.ActionCharge] = 57

	//n5 animations
	setNormCancel(attackFrames, 4, hitmarks[4][0], 50)
	attackFrames[4][action.ActionAttack] = 41
	attackFrames[4][action.ActionCharge] = 100 //TODO: this action is illegal; need better way to handle it

}

var swordFrames [][]int
var burstHitmarks = [][]int{{12}, {13}, {11}, {22, 33}, {33}}

func (c *char) initBurstAttackCancels() {

	//normal cancels
	swordFrames = make([][]int, c.NormalHitNum) //should be 5

	//n1 animations
	setNormCancel(swordFrames, 0, burstHitmarks[0][0], 24)
	swordFrames[0][action.ActionAttack] = 19
	swordFrames[0][action.ActionCharge] = 24

	//n2 animations
	setNormCancel(swordFrames, 1, burstHitmarks[1][0], 26)
	swordFrames[1][action.ActionAttack] = 16
	swordFrames[1][action.ActionCharge] = 26

	//n3 animations
	setNormCancel(swordFrames, 2, burstHitmarks[2][0], 34)
	swordFrames[2][action.ActionAttack] = 16
	swordFrames[2][action.ActionCharge] = 34

	//n4 animations
	setNormCancel(swordFrames, 3, burstHitmarks[3][1], 67)
	swordFrames[3][action.ActionAttack] = 44
	swordFrames[3][action.ActionCharge] = 67

	//n5 animations
	setNormCancel(swordFrames, 4, burstHitmarks[4][0], 83)
	swordFrames[4][action.ActionAttack] = 59
	swordFrames[4][action.ActionCharge] = 83
}

func setNormCancel(slice [][]int, index int, hitmark int, longest int) {
	slice[index] = make([]int, action.EndActionType)
	for i := range slice[index] {
		slice[index][i] = longest
	}
	slice[index][action.ActionSkill] = longest
	slice[index][action.ActionBurst] = longest
	slice[index][action.ActionDash] = longest
	slice[index][action.ActionJump] = longest
	slice[index][action.ActionSwap] = longest
}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	// Provide a custom override for Raiden's Q Normals and Charges
	if (c.Core.LastAction.Typ == core.ActionAttack ||
		c.Core.LastAction.Typ == core.ActionCharge) &&
		c.Core.Status.Duration("raidenburst") > 0 {
		f := 0
		switch c.Core.LastAction.Typ {
		case core.ActionAttack:
			f = burstNormalCancels(next, c.NormalCounter)
		case core.ActionCharge:
			f = burstChargeCancels(next)
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f
	}
	//otherwise use default implementation
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

func burstChargeCancels(next core.ActionType) int {
	switch next {
	case core.ActionAttack:
		return 56 - 24
	case core.ActionSkill:
		return 56 - 24
	case core.ActionDash:
		return 35 - 24
	case core.ActionJump:
		return 35 - 24
	case core.ActionSwap:
		return 55 - 24
	default:
		return 0
	}
}
