package frames

import "github.com/genshinsim/gcsim/pkg/core/action"

func InitNormalCancelSlice(slice *[][]int, index int, hitmark int, animation int) {
	var t [][]int = make([][]int, len(*slice))
	t[index] = make([]int, action.EndActionType)
	for i := range t[index] {
		t[index][i] = animation
	}
	t[index][action.ActionSkill] = animation
	t[index][action.ActionBurst] = animation
	t[index][action.ActionDash] = animation
	t[index][action.ActionJump] = animation
	t[index][action.ActionSwap] = animation
	slice = &t
}

func InitAbilSlice(animation int) []int {
	var t []int = make([]int, action.EndActionType)
	for i := range t {
		t[i] = animation
	}
	return t
}

func AtkSpdAdjust(f int, atkspd float64) int {
	if atkspd > 0.6 {
		atkspd = 0.6
	}
	return f + int(-0.5*atkspd*float64(f))
}
