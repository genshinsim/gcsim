package frames

import "github.com/genshinsim/gcsim/pkg/core/action"

func InitNormalCancelSlice(slice *[][]int, index int, hitmark int, animation int) {
	t := *slice
	t[index] = make([]int, action.EndActionType)
	for i := range t[index] {
		t[index][i] = animation
	}
	t[index][action.ActionSkill] = animation
	t[index][action.ActionBurst] = animation
	t[index][action.ActionDash] = animation
	t[index][action.ActionJump] = animation
	t[index][action.ActionSwap] = animation
}

func InitAbilSlice(slice *[]int, animation int) {

}
