package hutao

import "github.com/genshinsim/gcsim/pkg/core"

/**
[11:32 PM] sakuno | yanfei is my new maid: @gimmeabreak
https://www.youtube.com/watch?v=3aCiH2U4BjY

framecounts for 7 attempts of N2CJ (no hitlag):
83, 85, 88, 89, 77, 82, 84

first 3 not from the uploaded recording (as a n1cd player i cud barely pull it off :monkaS: )
YouTube
**/

//var normalFrames = []int{13, 16, 25, 36, 44, 39}               // from kqm lib
var normalFrames = []int{10, 13, 22, 33, 41, 36} // from kqm lib, -3 for hit lag
//var dmgFrame = [][]int{{13}, {16}, {25}, {36}, {26, 44}, {39}} // from kqm lib
var dmgFrame = [][]int{{10}, {13}, {22}, {33}, {23, 41}, {36}} // from kqm lib - 3 for hit lag

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := normalFrames[c.NormalCounter]
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f, f
	case core.ActionCharge:
		return 9, 9 //rough.. 11, -2 for hit lag
	case core.ActionSkill:
		return 42, 42 // from kqm lib
	case core.ActionBurst:
		return 130, 130 // from kqm lib
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Key.String(), a)
		return 0, 0
	}
}
