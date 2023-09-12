package frames

import (
	"github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

func InitNormalCancelSlice(hitmark, animation int) []int {
	t := make([]int, action.EndActionType)
	for i := range t {
		t[i] = animation
	}
	t[action.ActionAim] = hitmark
	t[action.ActionSkill] = hitmark
	t[action.ActionBurst] = hitmark
	t[action.ActionDash] = hitmark
	t[action.ActionJump] = hitmark
	t[action.ActionSwap] = hitmark
	return t
}

func InitAbilSlice(animation int) []int {
	t := make([]int, action.EndActionType)
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

func NewAttackFunc(c *character.Character, slice [][]int) func(action.Action) int {
	return func(next action.Action) int {
		// back out what last attack was
		n := c.NormalCounter - 1
		if n < 0 {
			n = c.NormalHitNum - 1
		}
		return AtkSpdAdjust(slice[n][next], c.Stat(attributes.AtkSpd))
	}
}

func NewAbilFunc(slice []int) func(action.Action) int {
	return func(next action.Action) int {
		return slice[next]
	}
}
